package game

type Card2021107Reshape struct{ AlwaysActive }

func (Card2021107Reshape) ID() string { return "2021107" }

func (Card2021107Reshape) Name() string { return "重塑" }

func (Card2021107Reshape) OnUseItem(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	for _, card := range ps.Hand {
		if card == nil {
			continue
		}
		ctx.Engine.discardHandCardToGraveyard(ctx.PlayerID, card)
	}
	ps.Hand = nil
	ctx.Engine.drawCards(ctx.PlayerID, 2)
	return nil
}
