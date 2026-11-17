package game

type Card2501001Shackle struct{ AlwaysActive }

func (Card2501001Shackle) ID() string { return "2501001" }

func (Card2501001Shackle) Name() string { return "桎梏" }

func (Card2501001Shackle) RevealsOnDraw() bool {
	return true
}

func (Card2501001Shackle) OnSelfDraw(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	if ctx.ExtraData != nil {
		if initial, _ := ctx.ExtraData["initial_hand"].(bool); initial {
			return nil
		}
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	for i, card := range ps.Hand {
		if card == nil || card.InstanceID != ctx.Source.InstanceID {
			continue
		}
		ctx.Engine.discardHandCardAt(ctx.PlayerID, i)
		ctx.Engine.drawCards(ctx.PlayerID, 1)
		return nil
	}
	return nil
}

func (Card2501001Shackle) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.drawCards(ctx.PlayerID, 1)
	return nil
}
