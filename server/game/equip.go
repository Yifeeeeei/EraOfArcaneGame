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

	slotIdx := -1
	var replacedEquipment *CardInstance
	newSubtype := restrictedEquipmentSubtype(card.Card)
	if replaceID != "" {
		for i := 0; i < equipmentSlotCapacity(ps); i++ {
			if ps.Equipment[i] != nil && ps.Equipment[i].InstanceID == replaceID {
				if ps.Equipment[i].IsHorizontal {
					return fmt.Errorf("can only replace vertical equipment")
				}
				if newSubtype != "" && restrictedEquipmentSubtype(ps.Equipment[i].Card) != newSubtype {
					return fmt.Errorf("restricted equipment can only replace same subtype")
				}
				replacedEquipment = ps.Equipment[i]
				slotIdx = i
				break
			}
		}
		if slotIdx == -1 {
			return fmt.Errorf("replacement equipment not found")
		}
	} else {
		if newSubtype != "" && !playerCanEquipDuplicateSubtypes(ps) {
			for _, equipment := range ps.Equipment {
				if equipment != nil && restrictedEquipmentSubtype(equipment.Card) == newSubtype {
					if equipment.IsHorizontal {
						return fmt.Errorf("same subtype equipment is horizontal and cannot be replaced")
					}
					return fmt.Errorf("same subtype equipment must be replaced")
				}
			}
		}
		// Find empty equipment slot
		for i := 0; i < equipmentSlotCapacity(ps); i++ {
			if ps.Equipment[i] == nil {
				slotIdx = i
				break
			}
		}
		if slotIdx == -1 {
			return fmt.Errorf("equipment area is full")
		}
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
