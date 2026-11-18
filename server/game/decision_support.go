package game

import (
	"encoding/json"
	"fmt"
)

// DecisionObservation is a detached, player-permitted snapshot. Controllers
// receive neither engine pointers nor the other player's hidden zones.
type DecisionObservation struct {
	State             json.RawMessage         `json:"state"`
	Cards             map[string]DecisionCard `json:"cards"`
	IncomingBurn      int                     `json:"incoming_burn,omitempty"`
	OffenseSuppressed bool                    `json:"offense_suppressed"`
}

type DecisionCard struct {
	CanLearn            bool             `json:"can_learn"`
	CanReplace          bool             `json:"can_replace"`
	CanSummon           bool             `json:"can_summon"`
	DrawnThisTurn       bool             `json:"drawn_this_turn"`
	DiscardAtEnd        bool             `json:"discard_at_end"`
	HasAura             bool             `json:"has_aura"`
	LearnReplacementIDs []string         `json:"learn_replacement_ids"`
	DevourCandidates    []string         `json:"devour_candidates"`
	CanAttackBoost      bool             `json:"can_attack_boost"`
	AttackBoostPower    int              `json:"attack_boost_power"`
	AttackBoostCost     map[string]int   `json:"attack_boost_cost"`
	Incoming            *DecisionImpact  `json:"incoming,omitempty"`
	ConsumeGain         map[string]int   `json:"consume_gain,omitempty"`
	CanConsume          bool             `json:"can_consume"`
	CanOverexert        bool             `json:"can_overexert"`
	CanPerTurn          bool             `json:"can_per_turn"`
	CanUltimate         bool             `json:"can_ultimate"`
	CanCast             bool             `json:"can_cast"`
	CanDefend           bool             `json:"can_defend"`
	UseCost             map[string]int   `json:"use_cost,omitempty"`
	DefenseCost         map[string]int   `json:"defense_cost,omitempty"`
	Targets             []DecisionTarget `json:"targets"`
	AttackTargets       []Position       `json:"attack_targets"`
}

type DecisionTarget struct {
	AreaOrigin  *Position        `json:"area_origin,omitempty"`
	ExtraTarget *Position        `json:"extra_target,omitempty"`
	Impacts     []DecisionImpact `json:"impacts"`
	Type        string           `json:"type"`
	Position    Position         `json:"position"`
	InstanceID  string           `json:"instance_id,omitempty"`
	Cost        map[string]int   `json:"cost"`
	Power       int              `json:"power"`
	Damage      int              `json:"damage"`
}

type DecisionImpact struct {
	InstanceID string `json:"instance_id"`
	Damage     int    `json:"damage"`
	Burn       int    `json:"burn"`
	Stun       int    `json:"stun"`
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
	for _, owner := range e.State.Players {
		for _, modifier := range owner.TempModifiers {
			if modifier.Type == TempModAllSpellDamageZeroUntilEnd {
				result.OffenseSuppressed = true
			}
		}
	}
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
		opt.HasAura = e.decisionHasAura(card)
		opt.CanReplace = onField && card.Card.IsSkill() && cardCanBeReplaced(card)
		if hand, _ := ps.FindHandCard(card.InstanceID); hand != nil {
			opt.CanSummon = card.Card.IsCompanion() && !e.actionRestricted(RuleSummon, nil)
			opt.DrawnThisTurn = ps.DrawnTurn[card.InstanceID] == e.State.TurnNumber
			opt.DiscardAtEnd = ps.DiscardAtTurnEnd[card.InstanceID]
			if len(summonDevourRequirement(card)) > 0 {
				for _, candidate := range field {
					if isValidSummonDevourTarget(candidate, card) {
						opt.DevourCandidates = append(opt.DevourCandidates, candidate.InstanceID)
					}
				}
			}
		}
		for _, skill := range ps.SkillPool {
			if skill != card {
				continue
			}
			opt.CanLearn = e.validateSkillLearnPermissionModifiers(player, card) == nil
			for i := 0; i < skillSlotCapacity(ps); i++ {
				if ps.Skills[i] != nil && cardCanBeReplaced(ps.Skills[i]) && skillAllowedInSlot(ps, card, i) {
					opt.LearnReplacementIDs = append(opt.LearnReplacementIDs, ps.Skills[i].InstanceID)
				}
			}
		}
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
			opt.CanAttackBoost = onField && e.validateSkillForPurpose(card, skillPurposeAttackBoost) == nil
			opt.AttackBoostPower = e.effectiveSkillPowerForPurposeWithData(player, card, nil, skillPurposeAttackBoost, map[string]any{"stat": "power"})
			opt.AttackBoostCost = e.effectiveSkillUseCostForPurpose(ps, card, skillPurposeAttackBoost)
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
						for _, variant := range e.decisionTargetVariants(player, card, target) {
							affected := e.spellAffectedUnits(1-player, card, variant)
							base := DecisionTarget{Type: "unit", Position: variant.Position, AreaOrigin: variant.AreaOrigin, InstanceID: unit.InstanceID, Cost: cost, Power: e.effectiveSpellPower(player, card, nil, variant), Damage: e.effectiveSpellDamage(player, card, max(0, card.Card.Attack+card.AttackBonus), nil, affected)}
							base.Impacts = e.decisionSpellImpacts(player, card, base.Damage, affected)
							opt.Targets = append(opt.Targets, base)
							if e.hasNextDriveSpellExtraTarget(ps, card) {
								for _, extra := range e.getAllFieldCards(e.State.Players[1-player]) {
									if extra.Position == nil {
										continue
									}
									t := SpellTarget{Type: "unit", Position: *extra.Position}
									if e.validateSpellExtraTargetForSkill(player, card, variant, t, e.skillHasPierce(player, card)) != nil {
										continue
									}
									copy := base
									pos := *extra.Position
									copy.ExtraTarget = &pos
									copy.Impacts = e.decisionSpellImpacts(player, card, base.Damage, append(append([]*CardInstance{}, affected...), extra))
									opt.Targets = append(opt.Targets, copy)
								}
							}
						}
					}
				}
			}
		}
		result.Cards[card.InstanceID] = opt
	}
	for _, unit := range e.getAllFieldCards(e.State.Players[1-player]) {
		if unit.IsSetCounter {
			continue
		}
		result.Cards[unit.InstanceID] = DecisionCard{HasAura: e.decisionHasAura(unit)}
	}
	if spell := e.State.PendingSpell; spell != nil && e.State.Phase == PhaseDefenseWindow && e.decisionPlayer() == player {
		affected := e.spellAffectedUnits(player, spell.Skill, spell.Target)
		for _, t := range spell.ExtraTargets {
			if unit := e.spellTargetUnitForCaster(spell.AttackerID, t); unit != nil {
				affected = append(affected, unit)
			}
		}
		damage := e.effectiveSpellDamage(spell.AttackerID, spell.Skill, max(0, spell.Skill.Card.Attack+spell.Skill.AttackBonus), spell.BoostSkills, affected)
		for _, impact := range e.decisionSpellImpacts(spell.AttackerID, spell.Skill, damage, affected) {
			opt := result.Cards[impact.InstanceID]
			copy := impact
			opt.Incoming = &copy
			result.Cards[impact.InstanceID] = opt
		}
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
