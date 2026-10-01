package game

// CardSearchedEvent is emitted for a completed search, independent of whether
// another effect currently owns the player's choice window.
type CardSearchedEvent struct {
	PlayerID int
	Card     *CardInstance
}
type CardSearchedBehavior interface {
	HasActiveCardSearched(*CardInstance) bool
	OnCardSearched(*EffectContext, CardSearchedEvent)
}

func (AlwaysActive) HasActiveCardSearched(*CardInstance) bool { return true }

func (e *Engine) notifyCardSearched(playerID int, card *CardInstance) {
	e.notifyCardSearchedThen(playerID, card, nil)
}
func (e *Engine) notifyCardSearchedThen(playerID int, card *CardInstance, after func()) {
	if e == nil || card == nil || card.Card == nil {
		if after != nil {
			after()
		}
		return
	}
	func() {
		ps := e.State.Players[playerID]
		if ps == nil {
			return
		}
		event := CardSearchedEvent{PlayerID: playerID, Card: card}
		for _, source := range e.getAllFieldCards(ps) {
			if source == nil || e.hasEffectiveStatus(source, StatusPetrify) {
				continue
			}
			if b, ok := cardBehavior(source).(CardSearchedBehavior); ok && b.HasActiveCardSearched(source) {
				b.OnCardSearched(e.skillContext(playerID, source), event)
			}
		}
	}()
	if after != nil {
		e.runResolution("after card searched", after)
	}
}
