package game

type Card2421003SturdyScroll struct{}

func (Card2421003SturdyScroll) ID() string   { return "2421003" }
func (Card2421003SturdyScroll) Name() string { return "坚固卷轴" }
func (Card2421003SturdyScroll) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:          TempModFriendlySpellDamageMinus,
		Amount:        1,
		RemainingUses: 1,
		ExpiresTurn:   ctx.Engine.State.TurnNumber + 2,
	})
	return nil
}
