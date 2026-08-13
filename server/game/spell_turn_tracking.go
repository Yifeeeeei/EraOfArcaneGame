package game

import "eraofarcane/model"

func (e *Engine) recordSpellCast(playerID int, skill *CardInstance) {
	if e == nil || skill == nil || skill.Card == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	ps := e.State.Players[playerID]
	if ps == nil {
		return
	}
	if ps.SpellsCastThisTurn == nil {
		ps.SpellsCastThisTurn = make(map[string]int)
	}
	if ps.SpellsCastByNumberThisTurn == nil {
		ps.SpellsCastByNumberThisTurn = make(map[string]int)
	}
	ps.SpellsCastThisTurn[skill.Card.Category]++
	ps.SpellsCastByNumberThisTurn[skill.Card.Number]++
	if skill.Card.Category == model.ElementWater && isSpellLikeCard(skill.Card) && totalElementCost(skillUseCost(skill.Card)) < 3 {
		ps.LastLowCostWaterSpell = cloneVirtualSpell(skill, playerID, e.State.TurnNumber)
	}
}

func clearSpellCastTracking(ps *PlayerState) {
	if ps == nil {
		return
	}
	ps.SpellsCastThisTurn = make(map[string]int)
	ps.SpellsCastByNumberThisTurn = make(map[string]int)
}

func (e *Engine) recordSpellHitStats(playerID int, targetCount, damage int) {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	ps := e.State.Players[playerID]
	if ps == nil {
		return
	}
	ps.SpellHitsThisTurn++
	ps.SpellHitTargetsThisTurn += max(targetCount, 0)
	ps.SpellDamageThisTurn += max(damage, 0)
}

func rollSpellHitTracking(ps *PlayerState) {
	if ps == nil {
		return
	}
	ps.SpellHitsLastTurn = ps.SpellHitsThisTurn
	ps.SpellHitTargetsLastTurn = ps.SpellHitTargetsThisTurn
	ps.SpellDamageLastTurn = ps.SpellDamageThisTurn
	ps.SpellHitsThisTurn = 0
	ps.SpellHitTargetsThisTurn = 0
	ps.SpellDamageThisTurn = 0
}

func spellCastByNumberThisTurn(ps *PlayerState, number string) int {
	if ps == nil || ps.SpellsCastByNumberThisTurn == nil || number == "" {
		return 0
	}
	return ps.SpellsCastByNumberThisTurn[number]
}

func totalSpellsCastThisTurn(ps *PlayerState) int {
	if ps == nil || ps.SpellsCastThisTurn == nil {
		return 0
	}
	total := 0
	for _, count := range ps.SpellsCastThisTurn {
		total += count
	}
	return total
}
