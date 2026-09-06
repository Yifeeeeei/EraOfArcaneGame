package game

type Card3221105CorrosiveFlow struct{ AlwaysActive }

func (Card3221105CorrosiveFlow) ID() string { return "3221105" }

func (Card3221105CorrosiveFlow) Name() string { return "腐蚀之流" }

func (Card3221105CorrosiveFlow) OnSpellHit(ctx *EffectContext) error {
	ctx.Engine.discardRandomHandCard(ctx.OpponentID)
	return nil
}
