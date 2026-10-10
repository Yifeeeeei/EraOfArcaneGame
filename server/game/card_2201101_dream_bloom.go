package game

type Card2201101DreamBloom struct{ AlwaysActive }

func (Card2201101DreamBloom) ID() string { return "2201101" }

func (Card2201101DreamBloom) Name() string { return "幻创之梦-绽放" }

func (Card2201101DreamBloom) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.drawCards(ctx.PlayerID, 3)
	return nil
}
