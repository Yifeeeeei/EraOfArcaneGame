package game

type Card3021107ArcaneShield struct{ AlwaysActive }

func (Card3021107ArcaneShield) ID() string { return "3021107" }

func (Card3021107ArcaneShield) Name() string { return "奥能护盾" }

func (Card3021107ArcaneShield) OnSpellCast(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModDelayedShieldGain,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		Amount:           1,
	})
	return nil
}
