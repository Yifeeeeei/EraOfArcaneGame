package game

type Card2411002EarthsplitterSword struct{ AlwaysActive }

func (Card2411002EarthsplitterSword) ID() string { return "2411002" }

func (Card2411002EarthsplitterSword) Name() string { return "裂地巨剑 阿托比斯" }

func (Card2411002EarthsplitterSword) OnConsume(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{Type: TempModSkillPowerBonus, Amount: 4, RemainingUses: 1, ExpiresTurn: ctx.Engine.State.TurnNumber + 1})
	return nil
}
