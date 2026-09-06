package game

type Card1021014ImpatientJunior struct{ AlwaysActive }

func (Card1021014ImpatientJunior) ID() string { return "1021014" }

func (Card1021014ImpatientJunior) Name() string { return "急不可耐的小师弟" }

func (Card1021014ImpatientJunior) OnEnter(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{Type: TempModNextLearnedSkillHaste, RemainingUses: 1, ExpiresTurn: ctx.Engine.State.TurnNumber + 1})
	return nil
}
