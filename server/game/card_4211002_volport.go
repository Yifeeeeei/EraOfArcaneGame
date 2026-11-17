package game

type Card4211002Volport struct{ AlwaysActive }

func (Card4211002Volport) ID() string { return "4211002" }

func (Card4211002Volport) Name() string { return "大贤者 沃尔波特" }

func (Card4211002Volport) OnEnter(ctx *EffectContext) error {
	addSkillToPool(ctx, "3201001")
	return nil
}
