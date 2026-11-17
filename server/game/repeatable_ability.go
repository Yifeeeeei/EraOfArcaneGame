package game

// A plain active effect is limited by its printed payment, rather than the
// default once-per-round budget of a 回合技. The legacy action name is shared.
type RepeatableActiveAbilityBehavior interface {
	RepeatableActiveAbility() bool
}

func cardHasRepeatableActiveAbility(card *CardInstance) bool {
	if card == nil || card.Card == nil {
		return false
	}
	b, ok := cardBehavior(card).(RepeatableActiveAbilityBehavior)
	return ok && b.RepeatableActiveAbility()
}
