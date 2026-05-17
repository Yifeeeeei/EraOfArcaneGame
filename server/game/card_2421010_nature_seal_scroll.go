package game

type Card2421010NatureSealScroll struct{}

func (Card2421010NatureSealScroll) ID() string   { return "2421010" }
func (Card2421010NatureSealScroll) Name() string { return "自然封印卷轴" }
func (Card2421010NatureSealScroll) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:          TempModAllSpellDamageZero,
		RemainingUses: 1,
		ExpiresTurn:   ctx.Engine.State.TurnNumber + 2,
	})
	return nil
}
