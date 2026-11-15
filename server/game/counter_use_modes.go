package game

import "fmt"

type CounterNormalUseBehavior interface{ AllowsNormalCounterUse() bool }

func counterAllowsNormalUse(card *CardInstance) bool {
	if card == nil || card.Card == nil {
		return false
	}
	b, ok := cardBehavior(card).(CounterNormalUseBehavior)
	return ok && b.AllowsNormalCounterUse()
}
func (e *Engine) promptCounterUseMode(playerID int, card *CardInstance) error {
	ps := e.State.Players[playerID]
	choices := []map[string]any{}
	if e.canPayCost(ps, e.effectiveCardPlayCost(ps, card)) {
		choices = append(choices, map[string]any{"instance_id": "normal", "name": "支付费用，立即使用", "zone": "choice"})
	}
	if len(e.counterSetChoices(playerID, card)) > 0 {
		choices = append(choices, map[string]any{"instance_id": "set", "name": "免费盖放，等待反制时机", "zone": "choice"})
	}
	if len(choices) == 0 {
		return fmt.Errorf("no legal item use")
	}
	e.SetPendingActionWithError(playerID, "counter_use_mode", "选择"+card.Card.Name+"的使用方式", choices, 1, 1, nil, false, func(selected []string, _ map[string]any) error {
		if firstSelected(selected) == "set" {
			e.promptSetCounterFromHand(playerID, card, false)
			return nil
		}
		if firstSelected(selected) != "normal" {
			return fmt.Errorf("invalid item use mode")
		}
		cost := e.effectiveCardPlayCost(ps, card)
		e.SetPendingActionWithError(playerID, "counter_normal_use_payment", "支付费用使用"+card.Card.Name, []map[string]any{candidateInfo(card, "hand", "own")}, 1, 1, cost, false, func(_ []string, data map[string]any) error {
			data["instance_id"] = card.InstanceID
			data["normal_use"] = true
			return e.handleUseItem(playerID, ActionMessage{Data: data})
		})
		return nil
	})
	return nil
}
func (e *Engine) counterSetChoices(playerID int, card *CardInstance) []map[string]any {
	ps := e.State.Players[playerID]
	choices := []map[string]any{}
	if _, _, err := prepareEquipmentPlacement(ps, card, ""); err == nil {
		return []map[string]any{{"instance_id": "set", "name": "盖放", "zone": "choice"}}
	}
	for _, equipment := range ps.Equipment {
		if equipment == nil {
			continue
		}
		if _, _, err := prepareEquipmentPlacement(ps, card, equipment.InstanceID); err == nil {
			choices = append(choices, candidateInfo(equipment, "equipment", "own"))
		}
	}
	return choices
}
func (e *Engine) promptSetCounterFromHand(playerID int, card *CardInstance, optional bool) {
	choices := e.counterSetChoices(playerID, card)
	if len(choices) == 0 {
		return
	}
	minSelect := 1
	if optional {
		minSelect = 0
	}
	e.SetPendingActionWithError(playerID, "set_counter_from_effect", "选择盖放位置（选择装备则将其替换）", choices, minSelect, 1, nil, false, func(selected []string, _ map[string]any) error {
		if len(selected) == 0 {
			return nil
		}
		ps := e.State.Players[playerID]
		current, index := ps.FindHandCard(card.InstanceID)
		if current != card {
			return fmt.Errorf("card is no longer in hand")
		}
		replace := firstSelected(selected)
		if replace == "set" {
			replace = ""
		}
		return e.placeCounterTrap(playerID, card, index, replace)
	})
}
