package game

type Card3021105ArcanePurification struct{ AlwaysActive }

func (Card3021105ArcanePurification) ID() string   { return "3021105" }
func (Card3021105ArcanePurification) Name() string { return "奥能净化" }

func (Card3021105ArcanePurification) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil {
		return nil
	}
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModFriendlyNegativeStatusIgnore,
		SourceCardNumber: "3021105",
		SourceName:       "奥能净化",
		ExpiresTurn:      ctx.Engine.State.TurnNumber,
	})
	return nil
}

var _ OnSpellCastBehavior = Card3021105ArcanePurification{}
