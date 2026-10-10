package game

type Card2621106PainScreamScroll struct{ AlwaysActive }

func (Card2621106PainScreamScroll) ID() string { return "2621106" }

func (Card2621106PainScreamScroll) Name() string { return "苦痛尖啸卷轴" }

func (Card2621106PainScreamScroll) OnUseItem(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModPainScreamWeakenOnDamage,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		ExpiresTurn:      ctx.Engine.State.TurnNumber,
	})
	return nil
}
