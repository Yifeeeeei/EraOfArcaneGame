package game

import (
	"fmt"
)

func isCounterTrapCard(number string) bool {
	_, ok := globalRegistry.GetBehavior(number).(CounterBehavior)
	return ok
}

func hiddenCounterInfo(ci *CardInstance) map[string]any {
	if ci == nil {
		return nil
	}
	return map[string]any{
		"instance_id":    ci.InstanceID,
		"name":           "盖放的卡牌",
		"type":           "道具",
		"is_hidden":      true,
		"is_set_counter": true,
		"is_horizontal":  ci.IsHorizontal,
		"slot_index":     ci.SlotIndex,
	}
}

const iceSoulSealCancelledBoostStatus = "ice_soul_seal_cancelled_boost"

func (e *Engine) placeCounterTrap(playerID int, card *CardInstance, handIdx int) error {
	ps := e.State.Players[playerID]
	slotIdx := -1
	for i := 0; i < equipmentSlotCapacity(ps); i++ {
		if ps.Equipment[i] == nil {
			slotIdx = i
			break
		}
	}
	if slotIdx == -1 {
		return fmt.Errorf("equipment area is full")
	}

	ps.RemoveFromHand(handIdx)
	card.IsSetCounter = true
	card.IsHorizontal = true
	card.SlotIndex = slotIdx
	card.EnterTurn = e.State.TurnNumber
	ps.Equipment[slotIdx] = card

	e.emit(GameEvent{
		Type:   "set_counter",
		Player: -1,
		Data: map[string]any{
			"player": playerID,
			"slot":   slotIdx,
		},
	})
	e.notifyCardEntered(playerID, card, map[string]any{"entered_player": playerID, "set_counter": true})
	return nil
}

func (e *Engine) promptCounterTrapIfEligible(counter *CardInstance, trigger EffectTrigger, eventSource *CardInstance, extraData map[string]any, afterResolve func()) bool {
	if counter == nil || counter.Card == nil || !counter.IsSetCounter {
		return false
	}
	if !counterTrapHasTrigger(counter.Card.Number, trigger) || !e.counterTrapConditionMet(counter, trigger, eventSource, extraData) {
		return false
	}

	ownerID := counter.OwnerID
	cost := e.effectiveCardPlayCost(e.State.Players[ownerID], counter)
	candidate := cardToInfo(counter)
	candidate["zone"] = "equipment"
	candidate["side"] = "own"
	context := e.counterTrapPendingContext(trigger, eventSource, extraData)
	available := func() bool {
		return !counterWindowCancelled(extraData) && counter.IsSetCounter &&
			e.counterTrapConditionMet(counter, trigger, eventSource, extraData)
	}
	action := e.setPendingActionWithOptions(ownerID, "counter_trigger",
		fmt.Sprintf("是否发动盖放的「%s」？", counter.Card.Name),
		[]map[string]any{candidate}, 0, 1, cost, true, nil, nil,
		func(selected []string, data map[string]any) error {
			if len(selected) == 0 {
				return nil
			}
			if selected[0] != counter.InstanceID {
				return fmt.Errorf("invalid counter selection")
			}
			if err := e.payAndRevealCounterTrap(ownerID, counter, cost, data); err != nil {
				return err
			}
			e.executeCounterTrap(counter, trigger, eventSource, extraData)

			return nil
		}, context, available)
	if action == nil {
		return false
	}
	e.addActionContinuation(action, "next counter", afterResolve)
	return true
}

func (e *Engine) eligibleCounterTraps(playerID int, trigger EffectTrigger, eventSource *CardInstance, extraData map[string]any) []*CardInstance {
	ps := e.State.Players[playerID]
	counters := []*CardInstance{}
	for _, card := range e.getAllFieldCards(ps) {
		if card == nil || card.Card == nil || !card.IsSetCounter {
			continue
		}
		if counterTrapHasTrigger(card.Card.Number, trigger) && e.counterTrapConditionMet(card, trigger, eventSource, extraData) {
			counters = append(counters, card)
		}
	}
	return counters
}

func (e *Engine) promptCounterTrapQueue(counters []*CardInstance, trigger EffectTrigger, eventSource *CardInstance, extraData map[string]any, afterDone func()) bool {
	prompted := false
	var promptNext func(int, bool)
	promptNext = func(index int, continuing bool) {
		if counterWindowCancelled(extraData) {
			if continuing && afterDone != nil {
				afterDone()
			}
			return
		}
		for index < len(counters) {
			counter := counters[index]
			index++
			if e.promptCounterTrapIfEligible(counter, trigger, eventSource, extraData, func() {
				promptNext(index, true)
			}) {
				prompted = true
				return
			}
		}
		if continuing && afterDone != nil {
			afterDone()
		}
	}
	promptNext(0, false)
	return prompted
}

func (e *Engine) counterTrapPendingContext(trigger EffectTrigger, eventSource *CardInstance, extraData map[string]any) map[string]any {
	context := map[string]any{
		"trigger":       triggerKey(trigger),
		"trigger_label": triggerLabel(trigger),
	}
	if eventSource != nil {
		context["source"] = cardToInfo(eventSource)
		context["source_player"] = eventSource.OwnerID
	}
	if extraData == nil {
		return context
	}
	if target, ok := extraData["target"].(*CardInstance); ok && target != nil {
		context["target"] = cardToInfo(target)
		context["target_player"] = target.OwnerID
	}
	if target, ok := extraData["target"].(SpellTarget); ok {
		context["spell_target"] = target
	}
	if castPlayer, ok := extraData["cast_player"].(int); ok {
		context["source_player"] = castPlayer
	}
	if consumedPlayer, ok := extraData["consumed_player"].(int); ok {
		context["source_player"] = consumedPlayer
	}
	if usedPlayer, ok := extraData["used_player"].(int); ok {
		context["source_player"] = usedPlayer
	}
	if enteredPlayer, ok := extraData["entered_player"].(int); ok {
		context["source_player"] = enteredPlayer
	}
	if damagedPlayer, ok := extraData["damaged_player"].(int); ok {
		context["target_player"] = damagedPlayer
	}
	if drawnPlayer, ok := extraData["drawn_player"].(int); ok {
		context["source_player"] = drawnPlayer
	}
	if endedPlayer, ok := extraData["ended_player"].(int); ok {
		context["source_player"] = endedPlayer
	}
	if drawCount := drawCountFromData(extraData); drawCount > 0 {
		context["draw_count"] = drawCount
	}
	if power := spellPowerFromData(extraData); power > 0 {
		context["power"] = power
	}
	return context
}

func triggerLabel(trigger EffectTrigger) string {
	switch trigger {
	case TriggerOnConsume:
		return "消耗"
	case TriggerOnSpellCast:
		return "施放法术"
	case TriggerOnUseItem:
		return "使用道具"
	case TriggerOnUnitEnter:
		return "伙伴入场"
	case TriggerOnTurnEnd:
		return "回合结束"
	case TriggerOnDraw:
		return "抽牌"
	case TriggerBeforeDamage:
		return "将受到伤害"
	case TriggerOnDamaged:
		return "受到伤害"
	case TriggerOnFriendlyDeath:
		return "友方死亡"
	case TriggerOnEnemyDeath:
		return "敌方死亡"
	case TriggerOnSpellHitBeforeDamage:
		return "法术命中"
	case TriggerOnCardEnter:
		return "卡牌入场"
	default:
		return fmt.Sprintf("触发%d", trigger)
	}
}

func triggerKey(trigger EffectTrigger) string {
	switch trigger {
	case TriggerOnConsume:
		return "consume"
	case TriggerOnSpellCast:
		return "spell_cast"
	case TriggerOnUseItem:
		return "use_item"
	case TriggerOnUnitEnter:
		return "unit_enter"
	case TriggerOnTurnEnd:
		return "turn_end"
	case TriggerOnDraw:
		return "draw"
	case TriggerBeforeDamage:
		return "before_damage"
	case TriggerOnDamaged:
		return "damaged"
	case TriggerOnFriendlyDeath:
		return "friendly_death"
	case TriggerOnEnemyDeath:
		return "enemy_death"
	case TriggerOnSpellHitBeforeDamage:
		return "spell_hit_before_damage"
	case TriggerOnCardEnter:
		return "card_enter"
	default:
		return fmt.Sprintf("trigger_%d", trigger)
	}
}

func counterTrapHasTrigger(number string, trigger EffectTrigger) bool {
	behavior, ok := globalRegistry.GetBehavior(number).(CounterBehavior)
	if !ok {
		return false
	}
	for _, candidate := range behavior.CounterTriggers() {
		if candidate == trigger {
			return true
		}
	}
	return false
}

func (e *Engine) counterTrapConditionMet(counter *CardInstance, trigger EffectTrigger, eventSource *CardInstance, extraData map[string]any) bool {
	if counter == nil || counter.Card == nil {
		return false
	}
	behavior, ok := globalRegistry.GetBehavior(counter.Card.Number).(CounterBehavior)
	return ok && behavior.HasActiveCounter(counter) && behavior.CanTriggerCounter(e.counterContext(counter, trigger, eventSource, extraData))
}

func (e *Engine) adjacentUnitCandidatesForCounter(counterOwnerID int, source *CardInstance) []map[string]any {
	if source == nil || source.Position == nil || source.OwnerID < 0 || source.OwnerID >= len(e.State.Players) {
		return nil
	}
	ps := e.State.Players[source.OwnerID]
	candidates := make([]map[string]any, 0, 4)
	for _, delta := range []struct{ col, row int }{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
		col := source.Position.Col + delta.col
		row := source.Position.Row + delta.row
		if col < 0 || col >= 3 || row < 0 || row >= 3 {
			continue
		}
		unit := ps.Units[col][row]
		if unit == nil || unit == source {
			continue
		}
		side := "enemy"
		if unit.OwnerID == counterOwnerID {
			side = "own"
		}
		candidates = append(candidates, candidateInfo(unit, "unit", side))
	}
	return candidates
}

func (e *Engine) payAndRevealCounterTrap(playerID int, counter *CardInstance, cost map[string]int, data map[string]any) error {
	overexertIDs := stringsFromAnySlice(anySliceFromData(data, "overexert_ids"))
	overexertUnits, err := e.collectOverexertUnits(e.State.Players[playerID], overexertIDs)
	if err != nil {
		return err
	}
	if !e.canPayCostWithOverexertOptions(e.State.Players[playerID], cost, overexertUnits, false) {
		return fmt.Errorf("not enough elements for counter")
	}
	if !e.payDefenseCostWithOptions(e.State.Players[playerID], cost, ActionMessage{Data: data}, overexertUnits, false) {
		return fmt.Errorf("invalid payment")
	}
	e.destroyFuyeDoomedAfterExert(overexertUnits)
	counter.IsSetCounter = false
	e.emit(GameEvent{Type: "counter_revealed", Player: -1, Data: map[string]any{
		"player": playerID,
		"card":   cardToInfo(counter),
	}})
	return nil
}

func (e *Engine) executeCounterTrap(counter *CardInstance, trigger EffectTrigger, eventSource *CardInstance, extraData map[string]any) {
	if behavior, ok := globalRegistry.GetBehavior(counter.Card.Number).(CounterResolutionBehavior); ok {
		behavior.ResolveCounter(e.counterContext(counter, trigger, eventSource, extraData))
	} else if len(globalRegistry.GetEffects(counter.Card.Number, trigger)) > 0 {
		e.triggerEffects(trigger, counter, eventSource, extraData)
	} else {
		e.triggerEffects(TriggerOnUseItem, counter, eventSource, extraData)
	}
	e.discardCounterTrap(counter.OwnerID, counter)
}

func (e *Engine) cancelBoostSpellWithIceSoulSeal(boost *CardInstance, extraData map[string]any) {
	if boost == nil {
		return
	}
	if boost.Statuses == nil {
		boost.Statuses = map[string]int{}
	}
	boost.Statuses[iceSoulSealCancelledBoostStatus] = 1
	if cancelled, ok := extraData["cancel_boost"].(*bool); ok && cancelled != nil {
		*cancelled = true
	}
	if e.State.PendingSpell == nil {
		return
	}
	removed := false
	kept := e.State.PendingSpell.BoostSkills[:0]
	for _, skill := range e.State.PendingSpell.BoostSkills {
		if skill == boost {
			removed = true
			continue
		}
		kept = append(kept, skill)
	}
	if removed {
		e.State.PendingSpell.BoostSkills = kept
		e.refreshPendingSpellPowerForModifiedSkill(e.State.PendingSpell.AttackerID, e.State.PendingSpell.Skill)
	}
}

func (e *Engine) promptOpponentCounterTrap(playerID int, trigger EffectTrigger, eventSource *CardInstance, extraData map[string]any, afterDecline func()) bool {
	opponentID := 1 - playerID
	return e.promptCounterTrapQueue(e.eligibleCounterTraps(opponentID, trigger, eventSource, extraData), trigger, eventSource, extraData, afterDecline)
}

func (e *Engine) discardCounterTrap(playerID int, counter *CardInstance) {
	ps := e.State.Players[playerID]
	for i, card := range ps.Equipment {
		if card != nil && card.InstanceID == counter.InstanceID {
			ps.Equipment[i] = nil
			e.addToGraveyard(playerID, counter)
			e.emit(GameEvent{Type: "discard", Player: playerID, Data: map[string]any{"card": cardToInfo(counter)}})
			return
		}
	}
}

func anySliceFromData(data map[string]any, key string) []any {
	if data == nil {
		return nil
	}
	values, _ := data[key].([]any)
	return values
}

func drawCountFromData(data map[string]any) int {
	if data == nil {
		return 0
	}
	if n, ok := data["draw_count_this_turn"].(int); ok {
		return n
	}
	if f, ok := data["draw_count_this_turn"].(float64); ok {
		return int(f)
	}
	return 0
}

func lethalDamageFromData(data map[string]any, target *CardInstance) bool {
	if data == nil || target == nil || target.Card == nil {
		return false
	}
	damage := damageFromData(data)
	return damage > 0 && target.CurrentLife-damage <= 0
}

func damageFromData(data map[string]any) int {
	if data == nil {
		return 0
	}
	if n, ok := data["damage"].(int); ok {
		return n
	}
	if f, ok := data["damage"].(float64); ok {
		return int(f)
	}
	return 0
}

func intFromData(data map[string]any, key string, fallback int) int {
	if data == nil {
		return fallback
	}
	if n, ok := data[key].(int); ok {
		return n
	}
	if f, ok := data[key].(float64); ok {
		return int(f)
	}
	return fallback
}

func boolFromData(data map[string]any, key string) bool {
	if data == nil {
		return false
	}
	value, _ := data[key].(bool)
	return value
}

func spellInstancesFromData(data map[string]any, key string) []*CardInstance {
	if data == nil {
		return nil
	}
	if cards, ok := data[key].([]*CardInstance); ok {
		return cards
	}
	return nil
}

func spellPowerFromData(data map[string]any) int {
	if data == nil {
		return 0
	}
	if n, ok := data["power"].(int); ok {
		return n
	}
	if f, ok := data["power"].(float64); ok {
		return int(f)
	}
	return 0
}

func counterWindowCancelled(data map[string]any) bool {
	if data == nil {
		return false
	}
	if cancelled, ok := data["cancel_item"].(*bool); ok && cancelled != nil && *cancelled {
		return true
	}
	if cancelled, ok := data["cancel_spell_hit"].(*bool); ok && cancelled != nil && *cancelled {
		return true
	}
	if cancelled, ok := data["cancel_spell_cast"].(*bool); ok && cancelled != nil && *cancelled {
		return true
	}
	return false
}
