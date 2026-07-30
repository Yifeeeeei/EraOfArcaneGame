package game

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
}

func spellCastByNumberThisTurn(ps *PlayerState, number string) int {
	if ps == nil || ps.SpellsCastByNumberThisTurn == nil || number == "" {
		return 0
	}
	return ps.SpellsCastByNumberThisTurn[number]
}
