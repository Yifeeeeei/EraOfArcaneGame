package game

type Card2421010NatureSealScroll struct{ AlwaysActive }

func (Card2421010NatureSealScroll) ID() string   { return "2421010" }
func (Card2421010NatureSealScroll) Name() string { return "自然封印卷轴" }
func (Card2421010NatureSealScroll) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:          TempModAllSpellDamageZero,
		RemainingUses: 1,
		ExpiresTurn:   ctx.Engine.State.TurnNumber + 1,
	})
	return nil
}
