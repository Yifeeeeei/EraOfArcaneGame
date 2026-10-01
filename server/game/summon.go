package game

import (
	"fmt"
)

// summonPlan is prepared and committed under the same engine lock. It contains
// the exact payment and devour costs; validation never mutates game state.
// Commit cannot fail validation or choose a different payment after costs fire.
type summonPlan struct {
	player       *PlayerState
	card         *CardInstance
	handIndex    int
	position     Position
	payment      map[string]int
	strictArcane int
	devour       []*CardInstance
}

func (e *Engine) handleSummon(playerID int, action ActionMessage) error {
	plan, err := e.prepareSummon(playerID, action)
	if err != nil {
		return err
	}
	e.commitSummon(plan)
	return nil
}

func (e *Engine) prepareSummon(playerID int, action ActionMessage) (*summonPlan, error) {
	if e.State.Phase != PhaseMain {
		return nil, fmt.Errorf("not in main phase")
	}
	if e.State.CurrentTurn != playerID {
		return nil, fmt.Errorf("not your turn")
	}
	if e.actionRestricted(RuleSummon, nil) {
		return nil, fmt.Errorf("a field rule prevents summoning cards")
	}

	instanceID, _ := action.Data["instance_id"].(string)
	pos, err := requiredBoardPosition(action.Data, "col", "row")
	if err != nil {
		return nil, err
	}
	col := pos.Col
	row := pos.Row

	ps := e.State.Players[playerID]

	// Find card in hand
	card, handIdx := ps.FindHandCard(instanceID)
	if card == nil {
		return nil, fmt.Errorf("card not found in hand")
	}

	// Must be a companion or item
	if !card.Card.IsCompanion() {
		return nil, fmt.Errorf("can only summon companions to unit area")
	}

	cost := e.effectiveCardPlayCost(ps, card)
	if !e.canPayCost(ps, cost) {
		return nil, fmt.Errorf("not enough elements")
	}

	payment, strictArcane, ok := e.cardPaymentPlanForAction(ps, card, cost, cost, paymentPurposePlay, action)
	if !ok {
		return nil, fmt.Errorf("invalid payment")
	}
	devour, err := e.prepareSummonDevour(playerID, card, action)
	if err != nil {
		return nil, err
	}
	removed := make(map[*CardInstance]bool, len(devour))
	unitCount := ps.CountUnits()
	for _, target := range devour {
		removed[target] = true
		if e.unitInOwnerGrid(target, playerID) {
			unitCount--
		}
	}
	if occupant := ps.Units[col][row]; occupant != nil && !removed[occupant] {
		return nil, fmt.Errorf("position already occupied")
	}

	// Validate the projected field after devour, without paying the cost yet.
	if unitCount >= 9 {
		return nil, fmt.Errorf("unit area is full")
	}

	return &summonPlan{player: ps, card: card, handIndex: handIdx, position: pos, payment: payment, strictArcane: strictArcane, devour: devour}, nil
}

func (e *Engine) commitSummon(plan *summonPlan) {
	e.beginResolution()
	defer e.endResolution()
	ps, card, handIdx, pos := plan.player, plan.card, plan.handIndex, plan.position
	playerID, col, row := ps.PlayerID, pos.Col, pos.Row
	e.spendPaymentPlan(ps, plan.payment, plan.strictArcane)
	for _, target := range plan.devour {
		if target.Card.IsCompanion() {
			e.destroyUnitWithCause(target, playerID, DeathCauseDevour)
		} else {
			e.discardFriendlyCandidate(playerID, target.InstanceID)
		}
	}
	e.notifyCardPlayCostPaid(ps, card)
	ps.RemoveFromHand(handIdx)
	card.Position = &Position{Col: col, Row: row}
	card.IsHorizontal = true // enters horizontal by default
	card.EnterTurn = e.State.TurnNumber
	ps.Units[col][row] = card

	// Apply keyword effects (速攻 makes it enter vertical, etc.)
	e.ApplyKeywordOnEnter(card)
	e.ApplySummonModifiersOnEnter(card)

	e.emit(GameEvent{
		Type:   "summon",
		Player: -1,
		Data: map[string]any{
			"player":   playerID,
			"card":     cardToInfo(card),
			"position": pos,
			"elements": ps.Elements,
		},
	})

	// Trigger 入场 (on enter) effects for the summoned card
	e.triggerEffects(TriggerOnEnter, card, nil, nil)

	enterData := map[string]any{"entered_player": playerID}
	e.notifyCardEntered(playerID, card, enterData)
	// Notify both sides about the new unit entering; individual card behaviors
	// decide whether they care about friendly or enemy units.
	e.triggerFieldEffectsWithData(TriggerOnUnitEnter, playerID, card, enterData)
	e.triggerFieldEffectsWithData(TriggerOnUnitEnter, 1-playerID, card, enterData)

	e.checkWinCondition()
}

func (e *Engine) prepareSummonDevour(playerID int, card *CardInstance, action ActionMessage) ([]*CardInstance, error) {
	requirement := summonDevourRequirement(card)
	cardRequirement := summonDevourCardRequirement(card)
	if len(requirement) == 0 && cardRequirement.Count <= 0 {
		return nil, nil
	}

	devourIDsRaw, _ := action.Data["devour_ids"].([]any)
	devourIDs := stringsFromAnySlice(devourIDsRaw)
	if legacyID, _ := action.Data["devour_id"].(string); legacyID != "" {
		devourIDs = append(devourIDs, legacyID)
	}
	if len(devourIDs) == 0 {
		return nil, fmt.Errorf("%s requires devour before summon", card.Card.Name)
	}

	ps := e.State.Players[playerID]
	targets := make([]*CardInstance, 0, len(devourIDs))
	total := make(map[string]int)
	cardRequirementCount := 0
	seen := make(map[string]bool, len(devourIDs))
	for _, devourID := range devourIDs {
		if seen[devourID] {
			return nil, fmt.Errorf("duplicate devour target")
		}
		seen[devourID] = true
		target := e.findFieldCardByInstance(ps, devourID)
		if target == nil {
			target = e.findUnitOnGrid(ps, devourID)
		}
		if !isValidSummonDevourTarget(target, card) {
			return nil, fmt.Errorf("invalid devour target")
		}
		if target.CurrentLife > 0 {
			total[DevourLife] += target.CurrentLife
		}
		for elem, amount := range e.effectiveElementsGain(target) {
			if amount > 0 {
				total[elem] += amount
			}
		}
		if cardSatisfiesDevourCardRequirement(target, cardRequirement) {
			cardRequirementCount++
		}
		targets = append(targets, target)
	}

	for elem, amount := range requirement {
		if total[elem] < amount {
			return nil, fmt.Errorf("devour targets load does not satisfy requirement")
		}
	}
	if cardRequirement.Count > 0 && cardRequirementCount < cardRequirement.Count {
		return nil, fmt.Errorf("devour targets do not satisfy card requirement")
	}
	return targets, nil
}

// handleSummon handles summoning a companion to the field
func isValidSummonDevourTarget(target *CardInstance, summoned *CardInstance) bool {
	if target == nil || target.Card == nil || target == summoned || target.Card.IsHero() {
		return false
	}
	return target.Card.IsCompanion() || isEquipmentItem(target)
}
