package game

type Card4411001Whitebeard struct{ AlwaysActive }

func (Card4411001Whitebeard) ID() string   { return "4411001" }
func (Card4411001Whitebeard) Name() string { return "森林隐士 白须" }
func (Card4411001Whitebeard) OnTurnStart(ctx *EffectContext) error {
	if ctx.Engine.State.TurnNumber != 1 {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	for i, c := range ps.Deck {
		if c.Card.Category == "地" && isBeastPlantOrSpirit(c) {
			ps.Hand = append(ps.Hand, c)
			ps.Deck = append(ps.Deck[:i], ps.Deck[i+1:]...)
			shuffleDeck(ps.Deck)
			ctx.Engine.emit(GameEvent{
				Type:   "effect_trigger",
				Player: ctx.PlayerID,
				Data: map[string]any{
					"source": cardToInfo(ctx.Source),
					"effect": "search",
					"card":   cardToInfo(c),
				},
			})
			return nil
		}
	}
	return nil
}
