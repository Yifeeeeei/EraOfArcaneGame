package game

import (
	"eraofarcane/model"
	"fmt"
)

func emitBatchEffect(ctx *EffectContext, effect string) {
	if ctx == nil || ctx.Engine == nil {
		return
	}
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(ctx.Source),
		"effect": effect,
	}})
}

func selectedUnitFromCandidates(e *Engine, selected []string, candidates []map[string]any) *CardInstance {
	if e == nil || len(selected) == 0 {
		return nil
	}
	allowed := make(map[string]bool, len(candidates))
	for _, candidate := range candidates {
		if id, _ := candidate["instance_id"].(string); id != "" {
			allowed[id] = true
		}
	}
	for _, id := range selected {
		if allowed[id] {
			return e.findUnitByInstanceID(id)
		}
	}
	return nil
}

func healUnit(card *CardInstance, amount int) {
	if card == nil || amount <= 0 {
		return
	}
	wasFull := maxLife(card) > 0 && card.CurrentLife >= maxLife(card)
	card.CurrentLife += amount
	if life := maxLife(card); life > 0 && card.CurrentLife > life {
		card.CurrentLife = life
	}
	if wasFull && card.Card != nil && card.Card.Number == "1521016" {
		card.Statuses["max_life_bonus"]++
		card.CurrentLife++
	}
}

func (e *Engine) healUnit(card *CardInstance, amount int, source *CardInstance) int {
	if e == nil || card == nil || amount <= 0 {
		return 0
	}
	before := card.CurrentLife
	convertsToLifeGain := maxLife(card) > 0 && card.CurrentLife >= maxLife(card) && card.Card != nil && card.Card.Number == "1521016"
	healUnit(card, amount)
	recovered := card.CurrentLife - before
	if convertsToLifeGain && recovered > 0 {
		e.notifyLifeGain(card, source, before)
	}
	return recovered
}

func (e *Engine) gainLife(card *CardInstance, amount int, source *CardInstance) int {
	if e == nil || card == nil || amount <= 0 {
		return 0
	}
	before := card.CurrentLife
	card.CurrentLife += amount
	return e.notifyLifeGain(card, source, before)
}

func (e *Engine) notifyLifeGain(card *CardInstance, source *CardInstance, before int) int {
	if e == nil || card == nil || card.OwnerID < 0 || card.OwnerID >= len(e.State.Players) {
		return 0
	}
	gained := card.CurrentLife - before
	if gained <= 0 {
		return 0
	}
	data := map[string]any{
		"life_gain_player": card.OwnerID,
		"life_gain_source": source,
		"life_gain_target": card,
		"amount":           gained,
	}
	e.triggerEffects(TriggerOnLifeGain, card, card, data)
	e.triggerFieldEffectsWithData(TriggerOnLifeGain, card.OwnerID, card, data)
	e.triggerFieldEffectsWithData(TriggerOnLifeGain, 1-card.OwnerID, card, data)
	return gained
}

func maxLife(card *CardInstance) int {
	if card == nil || card.Card == nil {
		return 0
	}
	return card.Card.Life + card.Statuses["max_life_bonus"]
}

func addGeneratedCardToHand(ctx *EffectContext, cardNumber string) {
	card := getCardDB()[cardNumber]
	if card == nil {
		return
	}
	ctx.Engine.addCardToHand(ctx.PlayerID, ctx.Engine.newCardInstance(card, ctx.PlayerID, ctx.Engine.State.TurnNumber))
	emitBatchEffect(ctx, "add_generated_card_to_hand")
}

func resetInstance(card *CardInstance) {
	if card == nil {
		return
	}
	card.IsHorizontal = false
	card.Statuses[StatusCooldown] = 0
	card.UsedThisTurn = 0
}

func reduceCost(cost map[string]int, elem string, amount int) {
	if amount <= 0 {
		return
	}
	if cost[elem] > 0 {
		cost[elem] -= amount
		if cost[elem] < 0 {
			cost[elem] = 0
		}
	}
}

func reduceGenericCost(cost map[string]int, preferredElem string, amount int) {
	if amount <= 0 {
		return
	}
	if preferredElem != "" && cost[preferredElem] > 0 {
		reduceCost(cost, preferredElem, amount)
		return
	}
	for _, elem := range model.AllElements {
		if cost[elem] > 0 {
			reduceCost(cost, elem, amount)
			return
		}
	}
}

func summonHandCompanionFree(ctx *EffectContext, predicate func(*CardInstance) bool) *CardInstance {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	pos := ps.FindEmptyPosition()
	if pos == nil {
		return nil
	}
	for i, card := range ps.Hand {
		if card == nil || !card.Card.IsCompanion() || (predicate != nil && !predicate(card)) {
			continue
		}
		ps.Hand = append(ps.Hand[:i], ps.Hand[i+1:]...)
		card.Position = pos
		card.IsHorizontal = true
		card.EnterTurn = ctx.Engine.State.TurnNumber
		ps.Units[pos.Col][pos.Row] = card
		ctx.Engine.ApplySummonModifiersOnEnter(card)
		ctx.Engine.triggerEffects(TriggerOnEnter, card, nil, nil)
		ctx.Engine.triggerFieldEffectsWithData(TriggerOnUnitEnter, ctx.PlayerID, card, map[string]any{"entered_player": ctx.PlayerID})
		return card
	}
	return nil
}

func totalFieldLoad(ps *PlayerState) int {
	total := 0
	for _, card := range ps.Units {
		for _, unit := range card {
			for _, amount := range effectiveElementsGain(unit) {
				total += amount
			}
		}
	}
	for _, card := range ps.Equipment {
		for _, amount := range effectiveElementsGain(card) {
			total += amount
		}
	}
	return total
}

const leviathanCooldownStatus = "leviathan_cooldown"

func inheritLifeSeedBonuses(e *Engine, source *CardInstance, target *CardInstance, playerID int) {
	if source == nil || target == nil {
		return
	}
	if bonusLife := source.CurrentLife - source.Card.Life; bonusLife > 0 {
		target.CurrentLife += bonusLife
	}
	target.CurrentAttack += max(source.CurrentAttack-source.Card.Attack, 0)
	target.PowerBonus += source.PowerBonus
	target.AttackBonus += source.AttackBonus
	for elem, amount := range source.ElementsGainBonus {
		if amount != 0 {
			e.addElementsGainBonus(target, playerID, elem, amount, source)
		}
	}
	if len(source.ElementsGainSet) > 0 {
		target.ElementsGainSet = copyElementCost(source.ElementsGainSet)
	}
	for status, amount := range source.Statuses {
		if amount <= 0 || status == StatusMastery {
			continue
		}
		target.Statuses[status] += amount
	}
}

func summonGreatDruidLifeSeedAtPosition(ctx *EffectContext, pos Position) {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	seedCard := getCardDB()["1401001"]
	if !pos.Valid() || ps.Units[pos.Col][pos.Row] != nil || seedCard == nil || ctx.Target == nil {
		return
	}
	seed := ctx.Engine.newCardInstance(seedCard, ctx.PlayerID, ctx.Engine.State.TurnNumber)
	seed.Position = &Position{Col: pos.Col, Row: pos.Row}
	seed.EnterTurn = ctx.Engine.State.TurnNumber
	seed.IsHorizontal = true
	if bonusLife := ctx.Target.CurrentLife - ctx.Target.Card.Life; bonusLife > 0 {
		seed.CurrentLife += bonusLife
	}
	for elem, amount := range ctx.Target.ElementsGainBonus {
		ctx.Engine.addElementsGainBonus(seed, ctx.PlayerID, elem, amount, ctx.Source)
	}
	if len(ctx.Target.ElementsGainSet) > 0 {
		seed.ElementsGainSet = copyElementCost(ctx.Target.ElementsGainSet)
	}
	ps.Units[pos.Col][pos.Row] = seed
	ctx.Engine.ApplySummonModifiersOnEnter(seed)
	ctx.Engine.triggerEffects(TriggerOnEnter, seed, nil, nil)
	ctx.Engine.notifyCardEntered(ctx.PlayerID, seed, map[string]any{"entered_player": ctx.PlayerID})
	ctx.Engine.triggerFieldEffectsWithData(TriggerOnUnitEnter, ctx.PlayerID, seed, map[string]any{"entered_player": ctx.PlayerID})
	ctx.Engine.triggerFieldEffectsWithData(TriggerOnUnitEnter, ctx.OpponentID, seed, map[string]any{"entered_player": ctx.PlayerID})
}

func resolveWhiteRobeSageControl(ctx *EffectContext, target *CardInstance, pos Position, cost map[string]int, data map[string]any) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if target == nil || target.Position == nil || !target.Card.IsCompanion() || target.OwnerID != ctx.OpponentID {
		return fmt.Errorf("invalid control target")
	}
	if !ctx.Engine.IsInSpellRange(ctx.PlayerID, target.Position.Col, target.Position.Row, cardHasPierce(ctx.Source)) {
		return fmt.Errorf("target out of range")
	}
	if !pos.Valid() || ps.Units[pos.Col][pos.Row] != nil {
		return fmt.Errorf("no empty position")
	}
	if !ctx.Engine.payCostForAction(ps, cost, ActionMessage{Data: data}) {
		return fmt.Errorf("invalid payment")
	}
	op := ctx.Engine.State.Players[ctx.OpponentID]
	op.Units[target.Position.Col][target.Position.Row] = nil
	target.OwnerID = ctx.PlayerID
	target.Position = &Position{Col: pos.Col, Row: pos.Row}
	ps.Units[pos.Col][pos.Row] = target
	return nil
}

const archmageStaffStoredSkillStatus = "archmage_staff_stored_skill"

const kingRobeReductionStatus = "君王法袍绝技减攻"

const nurEyeFireMark = "nur_eye_fire_mark"

const winterBowWaterMark = "winter_bow_water_mark"

func thunderDrumMarks(source *CardInstance) int {
	if source == nil {
		return 0
	}
	return source.Statuses["雷鼓标记"]
}

func spendThunderDrumMarks(source *CardInstance, amount int) {
	if source == nil || amount <= 0 {
		return
	}
	spend := min(source.Statuses["雷鼓标记"], amount)
	source.Statuses["雷鼓标记"] -= spend
}

const (
	windbreathCompassTemporaryAirStatus = "风息罗盘临时气负载"
)
