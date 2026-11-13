package game

import (
	"fmt"
)

// handleEquip handles equipping an item
func (e *Engine) handleEquip(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseMain {
		return fmt.Errorf("not in main phase")
	}
	if e.State.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}
	if e.actionRestricted(RuleEquip, nil) {
		return fmt.Errorf("a field rule prevents playing cards")
	}

	instanceID, _ := action.Data["instance_id"].(string)
	replaceID, _ := action.Data["replace_id"].(string)
	ps := e.State.Players[playerID]

	card, handIdx := ps.FindHandCard(instanceID)
	if card == nil {
		return fmt.Errorf("card not found in hand")
	}
	if !card.Card.IsItem() {
		return fmt.Errorf("card is not an item")
	}
	if !isEquipmentCard(card.Card) {
		return fmt.Errorf("card is not equipment")
	}
	cost := e.effectiveCardPlayCost(ps, card)
	if !e.canPayCost(ps, cost) {
		return fmt.Errorf("not enough elements")
	}

	slotIdx, replacedEquipment, err := prepareEquipmentPlacement(ps, card, replaceID)
	if err != nil {
		return err
	}

	if !e.payCostForCardAction(ps, card, cost, cost, paymentPurposePlay, action) {
		return fmt.Errorf("invalid payment")
	}
	e.notifyCardPlayCostPaid(ps, card)
	ps.RemoveFromHand(handIdx)
	if replacedEquipment != nil {
		e.moveEquipmentToGraveyard(playerID, slotIdx, replacedEquipment)
	}
	card.IsHorizontal = true
	card.SlotIndex = slotIdx
	card.EnterTurn = e.State.TurnNumber
	ps.Equipment[slotIdx] = card

	e.emit(GameEvent{
		Type:   "equip",
		Player: -1,
		Data: map[string]any{
			"player":   playerID,
			"card":     cardToInfo(card),
			"slot":     slotIdx,
			"elements": ps.Elements,
		},
	})

	e.triggerEffects(TriggerOnEquip, card, nil, nil)
	e.triggerEffects(TriggerOnEnter, card, nil, nil)
	e.notifyCardEntered(playerID, card, map[string]any{"entered_player": playerID, "equipped": true})

	return nil
}
