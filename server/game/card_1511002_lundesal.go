package game

type Card1511002Lundesal struct{ AlwaysActive }

func (Card1511002Lundesal) ID() string { return "1511002" }

func (Card1511002Lundesal) Name() string { return "大法师 伦德萨尔" }

func (Card1511002Lundesal) OnEnter(ctx *EffectContext) error {
	promptPermanentSkillBuff(ctx, "选择1个你的法术")
	return nil
}

func (Card1511002Lundesal) OnDeath(ctx *EffectContext) error {
	promptPermanentSkillBuff(ctx, "选择1个你的法术")
	return nil
}
