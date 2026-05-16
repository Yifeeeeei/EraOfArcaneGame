package game

type Card1021015EnergeticSenior struct{}

func (Card1021015EnergeticSenior) ID() string   { return "1021015" }
func (Card1021015EnergeticSenior) Name() string { return "精力充沛的大师兄" }

func (Card1021015EnergeticSenior) OnEnter(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModNextNoCooldown,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		RemainingUses:    1,
		ExpiresTurn:      ctx.Engine.State.TurnNumber,
	})
	return nil
}
