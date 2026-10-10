package game

type Card2021015ManaBoosterC struct{ AlwaysActive }

func (Card2021015ManaBoosterC) ID() string { return "2021015" }

func (Card2021015ManaBoosterC) Name() string { return "法力增强剂C型" }

func (Card2021015ManaBoosterC) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModCurrentTurnSkillCostZero,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		RemainingUses:    99,
		ExpiresTurn:      ctx.Engine.State.TurnNumber + 1,
	})
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModSkillUseCooldownAdd,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		Amount:           2,
		ExpiresTurn:      ctx.Engine.State.TurnNumber + 1,
	})
	return nil
}
