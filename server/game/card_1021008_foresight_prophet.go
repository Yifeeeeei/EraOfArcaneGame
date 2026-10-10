package game

type Card1021008ForesightProphet struct{ AlwaysActive }

func (Card1021008ForesightProphet) ID() string { return "1021008" }

func (Card1021008ForesightProphet) Name() string { return "预见先知" }

func (Card1021008ForesightProphet) OnBeforeDraw(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if len(ps.Deck) == 0 || ctx.Engine.State.PendingAction != nil {
		return nil
	}
	card := ps.Deck[0]
	candidates := []map[string]any{candidateInfo(card, "deck", "own")}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "foresight_prophet_top_card",
		"预见先知:查看牌堆顶1张牌,选择它则置于牌堆底,不选择则放回牌堆顶",
		candidates, 0, 1, func(selected []string) {
			if len(selected) == 0 || selected[0] != card.InstanceID || len(ps.Deck) == 0 || ps.Deck[0] != card {
				return
			}
			ps.Deck = append(ps.Deck[1:], card)
			emitBatchEffect(ctx, "peek_top_to_bottom")
		})
	return nil
}
