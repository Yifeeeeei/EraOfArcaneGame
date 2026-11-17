package game

import (
	"encoding/json"
	"fmt"
)

// DecisionObservation is a detached, player-permitted snapshot. Controllers
// receive neither engine pointers nor the other player's hidden zones.
type DecisionObservation struct {
	State        json.RawMessage         `json:"state"`
	Cards        map[string]DecisionCard `json:"cards"`
	IncomingBurn int                     `json:"incoming_burn,omitempty"`
}

type DecisionCard struct {
	ConsumeGain   map[string]int   `json:"consume_gain,omitempty"`
	CanConsume    bool             `json:"can_consume"`
	CanOverexert  bool             `json:"can_overexert"`
	CanPerTurn    bool             `json:"can_per_turn"`
	CanUltimate   bool             `json:"can_ultimate"`
	CanCast       bool             `json:"can_cast"`
	CanDefend     bool             `json:"can_defend"`
	UseCost       map[string]int   `json:"use_cost,omitempty"`
	DefenseCost   map[string]int   `json:"defense_cost,omitempty"`
	Targets       []DecisionTarget `json:"targets"`
	AttackTargets []Position       `json:"attack_targets"`
}

type DecisionTarget struct {
	Type       string         `json:"type"`
	Position   Position       `json:"position"`
	InstanceID string         `json:"instance_id,omitempty"`
	Cost       map[string]int `json:"cost"`
	Power      int            `json:"power"`
	Damage     int            `json:"damage"`
}

// ObserveDecision performs only pure rules queries. JSON detaches legacy maps
// while the engine lock is held, before an asynchronous controller uses them.
func (e *Engine) ObserveDecision(player int) (DecisionObservation, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if player < 0 || player > 1 || e.State == nil || e.State.Players[player] == nil {
		return DecisionObservation{}, fmt.Errorf("game not started")
	}
	view, err := json.Marshal(e.stateForPlayerLocked(player))
	if err != nil {
		return DecisionObservation{}, err
	}
	result := DecisionObservation{State: view, Cards: map[string]DecisionCard{}}
	ps := e.State.Players[player]
	if spell := e.State.PendingSpell; spell != nil && e.State.Phase == PhaseDefenseWindow && e.decisionPlayer() == player && !e.rejectsNegativeStatusApplication(ps.Hero, StatusBurn) {
		result.IncomingBurn = max(0, e.explicitSpellHitStatuses(spell.Skill, ps.Hero)[StatusBurn])
	}
	field := e.getAllFieldCards(ps)
	sources := append(append(append([]*CardInstance{}, field...), ps.Hand...), ps.SkillPool...)
	for _, host := range field {
		if host != nil {
			sources = append(sources, host.BoundSkills...)
		}
	}
	for _, card := range sources {
		if card == nil || card.Card == nil {
			continue
		}
		onField := e.findFieldCardByInstance(ps, card.InstanceID) == card || e.findSkill(ps, card.InstanceID) == card
		opt := DecisionCard{}
		if onField {
			opt.CanConsume = e.canConsumeCard(card)
			if opt.CanConsume {
				opt.ConsumeGain, _ = e.elementsGainedFromConsume(player, card, ActionMessage{})
			}
			opt.CanOverexert = e.canConsumeCard(card) && (card.Card.IsHero() || card.Card.IsCompanion() || isEquipmentCard(card.Card)) && totalElementCost(e.effectiveElementsGain(card)) > 0
			active := !e.hasEffectiveStatus(card, StatusPetrify)
			opt.CanPerTurn = active && cardHasActivePerTurn(card) && (cardHasRepeatableActiveAbility(card) || card.UsedThisTurn < perTurnLimit(card)) && e.validateAbility(card, TriggerPerTurn) == nil
			opt.CanUltimate = active && cardHasActiveUltimate(card) && !card.UltimateUsed && e.validateAbility(card, TriggerUltimate) == nil
			if effectiveCurrentAttack(card) > 0 && e.canConsumeCard(card) && !e.actionRestricted(RuleAttack, card) {
				isEquip := isEquipmentCard(card.Card)
				if isEquip || card.Position != nil && (card.Position.Row == ps.GetFrontRow() || cardCanAttackFromNonFront(card)) {
					for row := 0; row < 3; row++ {
						for col := 0; col < 3; col++ {
							if e.State.Players[1-player].Units[col][row] != nil && e.isInDirectAttackRange(player, card, isEquip, col, row) {
								opt.AttackTargets = append(opt.AttackTargets, Position{Col: col, Row: row})
							}
						}
					}
				}
			}
		}
		if isSpellLikeCard(card.Card) {
			opt.CanCast = onField && e.validateSkillForPurpose(card, skillPurposeAttack) == nil
			opt.CanDefend = (onField || isSpellScrollCard(card.Card)) && e.validateSkillForPurpose(card, skillPurposeDefend) == nil
			opt.UseCost = e.effectiveSkillUseCostForPurpose(ps, card, skillPurposeAttack)
			opt.DefenseCost = e.effectiveSkillUseCostForPurpose(ps, card, skillPurposeDefend)
			if isSpellScrollCard(card.Card) {
				opt.CanCast = true
				opt.CanDefend = canUseSkillForPurpose(card.Card, skillPurposeDefend) && e.validateSkillUsePermissionModifiers(card, skillPurposeDefend) == nil
				opt.UseCost = e.effectiveCardPlayCost(ps, card)
				opt.DefenseCost = opt.UseCost
			}
			if !skillNeedsTargetInstance(card) {
				opt.Targets = append(opt.Targets, DecisionTarget{Type: "none", Cost: opt.UseCost})
			} else {
				for row := 0; row < 3; row++ {
					for col := 0; col < 3; col++ {
						unit := e.State.Players[1-player].Units[col][row]
						if unit == nil {
							continue
						}
						target := SpellTarget{Type: "unit", Position: Position{Col: col, Row: row}}
						if e.validateSpellTarget(player, card, target) != nil {
							continue
						}
						cost := e.effectiveSkillUseCostForPurposeWithData(ps, card, skillPurposeAttack, map[string]any{"spell_target": target, "spell_target_unit": unit})
						if isSpellScrollCard(card.Card) {
							cost = opt.UseCost
						}
						affected := e.spellAffectedUnits(1-player, card, target)
						opt.Targets = append(opt.Targets, DecisionTarget{Type: "unit", Position: target.Position, InstanceID: unit.InstanceID, Cost: cost, Power: e.effectiveSpellPower(player, card, nil, target), Damage: e.effectiveSpellDamage(player, card, max(0, card.Card.Attack+card.AttackBonus), nil, affected)})
					}
				}
			}
		}
		result.Cards[card.InstanceID] = opt
	}
	// No rule-owned map or slice escapes the lock.
	encoded, err := json.Marshal(result)
	if err != nil {
		return DecisionObservation{}, err
	}
	var detached DecisionObservation
	err = json.Unmarshal(encoded, &detached)
	return detached, err
}

// DecisionPayment shares the engine's element matching with clients that plan
// ordinary actions. Strict arcane can satisfy only the wildcard requirement.
func DecisionPayment(available, cost map[string]int, strict int, lightWildcard, lightCostWildcard bool) (map[string]int, bool) {
	for n := min(strict, cost["无"]); n >= 0; n-- {
		remaining := cloneElements(cost)
		remaining["无"] -= n
		payment, ok := calculateElementPaymentWithOptions(available, remaining, lightWildcard, lightCostWildcard)
		if ok {
			payment["无"] += n
			return payment, true
		}
	}
	return nil, false
}
