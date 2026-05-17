package game

type Card1221011WinterfellWarlock struct{ AlwaysActive }

func (Card1221011WinterfellWarlock) ID() string   { return "1221011" }
func (Card1221011WinterfellWarlock) Name() string { return "凛冬城术士" }

func (Card1221011WinterfellWarlock) OnUltimate(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModNextSpellHitStatus,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		Status:           StatusFreeze,
		Amount:           1,
		RemainingUses:    1,
		ExpiresTurn:      ctx.Engine.State.TurnNumber,
	})
	return nil
}
