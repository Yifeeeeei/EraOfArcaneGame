package game

import "fmt"

var counterTrapTriggers = map[string][]EffectTrigger{
	"2021018": {TriggerOnSpellCast},
	"2021022": {TriggerOnUseItem},
	"2121002": {TriggerOnConsume},
	"2121012": {TriggerOnUnitEnter},
	"2221002": {TriggerOnConsume},
	"2221005": {TriggerOnTurnEnd},
	"2221010": {TriggerOnDraw},
	"2221011": {TriggerOnDamaged},
	"2321002": {TriggerOnConsume},
	"2321010": {TriggerOnSpellCast},
	"2321011": {TriggerOnUnitEnter, TriggerOnConsume},
	"2521002": {TriggerOnSpellHitBeforeDamage},
	"2521004": {TriggerOnSpellCast},
	"2521011": {TriggerOnSpellCast},
	"2621003": {TriggerOnUnitEnter},
	"2621005": {TriggerOnFriendlyDeath, TriggerOnEnemyDeath},
	"2621010": {TriggerOnFriendlyDeath},
	"2621011": {TriggerOnConsume},
}

func isCounterTrapCard(number string) bool {
	_, ok := counterTrapTriggers[number]
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
	return nil
}

func (e *Engine) promptCounterTrapIfEligible(counter *CardInstance, trigger EffectTrigger, eventSource *CardInstance, extraData map[string]any, afterResolve func()) bool {
	if counter == nil || counter.Card == nil || !counter.IsSetCounter || e.State.PendingAction != nil {
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
	e.SetPendingActionWithErrorAndContext(ownerID, "counter_trigger",
		fmt.Sprintf("是否发动盖放的「%s」？", counter.Card.Name),
		[]map[string]any{candidate}, 0, 1, cost, true, context,
		func(selected []string, data map[string]any) error {
			if len(selected) == 0 {
				if afterResolve != nil {
					afterResolve()
				}
				return nil
			}
			if selected[0] != counter.InstanceID {
				return fmt.Errorf("invalid counter selection")
			}
			if err := e.payAndRevealCounterTrap(ownerID, counter, cost, data); err != nil {
				return err
			}
			e.executeCounterTrap(counter, trigger, eventSource, extraData)
			if e.State.PendingAction != nil && e.State.PendingAction.Type != "counter_trigger" {
				e.wrapPendingActionContinuation(afterResolve)
				return nil
			}
			if afterResolve != nil {
				afterResolve()
			}
			return nil
		})
	return e.State.PendingAction != nil && e.State.PendingAction.Type == "counter_trigger"
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
	if e.State.PendingAction != nil {
		return false
	}
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
				return
			}
		}
		if continuing && afterDone != nil {
			afterDone()
		}
	}
	promptNext(0, false)
	return e.State.PendingAction != nil && e.State.PendingAction.Type == "counter_trigger"
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
	case TriggerOnDamaged:
		return "受到伤害"
	case TriggerOnFriendlyDeath:
		return "友方死亡"
	case TriggerOnEnemyDeath:
		return "敌方死亡"
	case TriggerOnSpellHitBeforeDamage:
		return "法术命中"
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
	case TriggerOnDamaged:
		return "damaged"
	case TriggerOnFriendlyDeath:
		return "friendly_death"
	case TriggerOnEnemyDeath:
		return "enemy_death"
	case TriggerOnSpellHitBeforeDamage:
		return "spell_hit_before_damage"
	default:
		return fmt.Sprintf("trigger_%d", trigger)
	}
}

func counterTrapHasTrigger(number string, trigger EffectTrigger) bool {
	for _, candidate := range counterTrapTriggers[number] {
		if candidate == trigger {
			return true
		}
	}
	return false
}

func (e *Engine) counterTrapConditionMet(counter *CardInstance, trigger EffectTrigger, eventSource *CardInstance, extraData map[string]any) bool {
	ownerID := counter.OwnerID
	sourceOwner := -1
	if eventSource != nil {
		sourceOwner = eventSource.OwnerID
	}
	if trigger != TriggerOnFriendlyDeath && trigger != TriggerOnEnemyDeath {
		if castPlayer, ok := extraData["cast_player"].(int); ok {
			sourceOwner = castPlayer
		}
		if attacker, ok := extraData["attacker"].(int); ok {
			sourceOwner = attacker
		}
		if drawnPlayer, ok := extraData["drawn_player"].(int); ok {
			sourceOwner = drawnPlayer
		}
		if damagedPlayer, ok := extraData["damaged_player"].(int); ok {
			sourceOwner = damagedPlayer
		}
		if enteredPlayer, ok := extraData["entered_player"].(int); ok {
			sourceOwner = enteredPlayer
		}
		if consumedPlayer, ok := extraData["consumed_player"].(int); ok {
			sourceOwner = consumedPlayer
		}
		if usedPlayer, ok := extraData["used_player"].(int); ok {
			sourceOwner = usedPlayer
		}
		if endedPlayer, ok := extraData["ended_player"].(int); ok {
			sourceOwner = endedPlayer
		}
	}

	switch counter.Card.Number {
	case "2021018":
		return trigger == TriggerOnSpellCast && sourceOwner != ownerID && len(e.friendlySkillsIncludingBound(ownerID, nil)) > 0
	case "2021022":
		return trigger == TriggerOnUseItem && sourceOwner != ownerID && eventSource != nil && counterRuneCanCancel(eventSource.Card.Number)
	case "2121002":
		return trigger == TriggerOnConsume && eventSource != nil && (eventSource.Card.IsHero() || eventSource.Card.IsCompanion())
	case "2121012", "2621003":
		return trigger == TriggerOnUnitEnter && sourceOwner != ownerID && eventSource != nil && eventSource.Card.IsCompanion()
	case "2221002":
		return trigger == TriggerOnConsume && sourceOwner != ownerID && eventSource != nil && eventSource.Card.IsCompanion()
	case "2221005":
		return trigger == TriggerOnTurnEnd && sourceOwner != ownerID
	case "2221010":
		return trigger == TriggerOnDraw && sourceOwner != ownerID && drawCountFromData(extraData) >= 3 &&
			len(e.friendlyUnits(ownerID, false, isWaterCompanion)) > 0
	case "2221011":
		return trigger == TriggerOnDamaged && sourceOwner == ownerID
	case "2321002":
		return trigger == TriggerOnConsume && sourceOwner != ownerID && eventSource != nil && (eventSource.Card.IsHero() || eventSource.Card.IsCompanion())
	case "2321010":
		return trigger == TriggerOnSpellCast && sourceOwner != ownerID && eventSource != nil && !isSorcerySkill(eventSource.Card) && e.State.PendingSpell != nil
	case "2321011":
		return eventSource != nil && eventSource.Card.IsCompanion() && eventSource.Position != nil &&
			len(e.emptyUnitPositionsForPlayer(eventSource.OwnerID, ownerID)) > 0
	case "2521002":
		return trigger == TriggerOnSpellHitBeforeDamage && sourceOwner != ownerID && eventSource != nil && !isSorcerySkill(eventSource.Card) && spellPowerFromData(extraData) < 10
	case "2521004":
		return trigger == TriggerOnSpellCast && sourceOwner != ownerID && eventSource != nil && isSorcerySkill(eventSource.Card)
	case "2521011":
		return trigger == TriggerOnSpellCast && sourceOwner != ownerID
	case "2621005":
		return eventSource != nil && eventSource.Card.IsCompanion()
	case "2621010":
		return trigger == TriggerOnFriendlyDeath && sourceOwner == ownerID && eventSource != nil && eventSource.DamageTakenThisTurn > 0
	case "2621011":
		return trigger == TriggerOnConsume && sourceOwner != ownerID && eventSource != nil && eventSource.Card.IsCompanion() && eventSource.CurrentAttack > 0 && len(e.adjacentUnitCandidatesForCounter(ownerID, eventSource)) > 0
	default:
		return false
	}
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
	if !canPayCostWithOverexert(e.State.Players[playerID], cost, overexertUnits) {
		return fmt.Errorf("not enough elements for counter")
	}
	if !payDefenseCost(e.State.Players[playerID], cost, ActionMessage{Data: data}, overexertUnits) {
		return fmt.Errorf("invalid payment")
	}
	counter.IsSetCounter = false
	e.emit(GameEvent{Type: "counter_revealed", Player: -1, Data: map[string]any{
		"player": playerID,
		"card":   cardToInfo(counter),
	}})
	return nil
}

func (e *Engine) executeCounterTrap(counter *CardInstance, trigger EffectTrigger, eventSource *CardInstance, extraData map[string]any) {
	if counter.Card.Number == "2021022" && trigger == TriggerOnUseItem {
		if cancelled, ok := extraData["cancel_item"].(*bool); ok && cancelled != nil {
			*cancelled = true
			e.emit(GameEvent{Type: "item_cancelled", Player: -1, Data: map[string]any{
				"player": counter.OwnerID,
				"card":   cardToInfo(eventSource),
				"source": cardToInfo(counter),
			}})
		}
		e.discardCounterTrap(counter.OwnerID, counter)
		return
	}
	if counter.Card.Number == "2521002" && trigger == TriggerOnSpellHitBeforeDamage {
		if cancelled, ok := extraData["cancel_spell_hit"].(*bool); ok && cancelled != nil {
			*cancelled = true
			e.emit(GameEvent{Type: "spell_hit_cancelled", Player: -1, Data: map[string]any{
				"player": counter.OwnerID,
				"skill":  cardToInfo(eventSource),
				"source": cardToInfo(counter),
			}})
		}
		e.discardCounterTrap(counter.OwnerID, counter)
		return
	}
	if len(globalRegistry.GetEffects(counter.Card.Number, trigger)) > 0 {
		e.triggerEffects(trigger, counter, eventSource, extraData)
	} else {
		e.triggerEffects(TriggerOnUseItem, counter, eventSource, extraData)
	}
	e.discardCounterTrap(counter.OwnerID, counter)
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
			ps.Graveyard = append(ps.Graveyard, counter)
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

func counterRuneCanCancel(number string) bool {
	switch number {
	case "2021010", "2021012", "2021018", "2021019", "2021021", "2021022",
		"2121002", "2121003", "2121008", "2121009", "2121011", "2121012",
		"2221002", "2221003", "2221008", "2221009", "2221010", "2221011", "2221013",
		"2321002", "2321003", "2321005", "2321008", "2321009", "2321010", "2321011",
		"2421002", "2421003", "2421004", "2421005", "2421008", "2421009", "2421010",
		"2521002", "2521003", "2521004", "2521005", "2521008", "2521009", "2521011", "2521012", "2521013",
		"2621001", "2621003", "2621005", "2621008", "2621009", "2621010", "2621011":
		return true
	default:
		return false
	}
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
