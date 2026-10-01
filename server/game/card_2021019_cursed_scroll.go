package game

type Card2021019CursedScroll struct{ AlwaysActive }

func (Card2021019CursedScroll) ID() string { return "2021019" }

func (Card2021019CursedScroll) Name() string { return "诅咒卷轴" }

func (Card2021019CursedScroll) OnUseItem(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	drawn := ctx.Engine.drawCards(ctx.PlayerID, 2)
	if ps.DiscardAtTurnEnd == nil {
		ps.DiscardAtTurnEnd = make(map[string]bool)
	}
	for _, card := range drawn {
		ps.DiscardAtTurnEnd[card.InstanceID] = true
	}
	return nil
}
