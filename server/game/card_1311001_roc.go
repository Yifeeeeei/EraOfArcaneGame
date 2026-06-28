package game

type Card1311001Roc struct{ AlwaysActive }

func (Card1311001Roc) ID() string   { return "1311001" }
func (Card1311001Roc) Name() string { return "大鹏" }

func (Card1311001Roc) OnEnter(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	drawn := make([]*CardInstance, 0)
	remaining := make([]*CardInstance, 0, len(ps.Deck))
	for i, card := range ps.Deck {
		if i < 8 && card != nil && totalCost(card.Card.ElementsCost) < 3 {
			drawn = append(drawn, card)
			continue
		}
		remaining = append(remaining, card)
	}
	ps.Deck = remaining
	for _, card := range drawn {
		ps.Hand = append(ps.Hand, card)
		ps.DrawCountThisTurn++
		if ps.DiscardAtTurnEnd == nil {
			ps.DiscardAtTurnEnd = make(map[string]bool)
		}
		ps.DiscardAtTurnEnd[card.InstanceID] = true
		ctx.Engine.emit(GameEvent{Type: "draw_card", Player: ctx.PlayerID, Data: map[string]any{"card": cardToInfo(card)}})
		ctx.Engine.triggerFieldEffectsWithData(TriggerOnDraw, ctx.PlayerID, card, map[string]any{
			"drawn_card":           card,
			"drawn_player":         ctx.PlayerID,
			"draw_count_this_turn": ps.DrawCountThisTurn,
		})
	}
	ctx.Engine.shuffleDeck(ctx.PlayerID)
	return nil
}
