package game

type Card2021014ManaBoosterA struct{}

func (Card2021014ManaBoosterA) ID() string   { return "2021014" }
func (Card2021014ManaBoosterA) Name() string { return "法力增强剂A型" }

func (Card2021014ManaBoosterA) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModNextSkillCostZero,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		RemainingUses:    1,
		ExpiresTurn:      ctx.Engine.State.TurnNumber,
	})
	return nil
}
