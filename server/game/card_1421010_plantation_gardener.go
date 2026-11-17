package game

type Card1421010PlantationGardener struct{ AlwaysActive }

func (Card1421010PlantationGardener) ID() string { return "1421010" }

func (Card1421010PlantationGardener) Name() string { return "种植园丁" }

func (Card1421010PlantationGardener) OnTurnStart(ctx *EffectContext) error {
	ctx.Source.Statuses["地脉标记"]++
	return nil
}

func (Card1421010PlantationGardener) OnLoadGain(ctx *EffectContext) error {
	if ctx.ExtraData == nil {
		return nil
	}
	if player, ok := ctx.ExtraData["load_gain_player"].(int); ok && player == ctx.PlayerID {
		ctx.Source.Statuses["地脉标记"]++
	}
	return nil
}

func (Card1421010PlantationGardener) OnPerTurn(ctx *EffectContext) error {
	if ctx.Source.Statuses["地脉标记"] >= 2 {
		ctx.Source.Statuses["地脉标记"] -= 2
		ctx.Engine.drawCards(ctx.PlayerID, 1)
	}
	return nil
}
