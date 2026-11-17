package game

// CardStateChange informs a rule provider after a batch of state changes. This
// structural notification also runs when a card becomes inactive or leaves a
// zone, so it can withdraw a continuous rule. It is not a player-triggered ability.
type CardStateChange struct {
	Card          *CardInstance
	Status        string
	Before, After int
	LeftField     bool
}

type CardStateChangeBehavior interface {
	OnCardStateChange(*EffectContext, CardStateChange)
}

func (e *Engine) notifyCardStateChanges(changes ...CardStateChange) {
	for _, change := range changes {
		if change.Card == nil {
			continue
		}
		if b, ok := cardBehavior(change.Card).(CardStateChangeBehavior); ok {
			b.OnCardStateChange(e.skillContext(change.Card.OwnerID, change.Card), change)
		}
	}
}
