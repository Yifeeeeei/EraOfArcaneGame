package game

type Card3021012MindTempering struct{ AlwaysActive }

func (Card3021012MindTempering) ID() string { return "3021012" }

func (Card3021012MindTempering) Name() string { return "心炼" }

func (Card3021012MindTempering) OnSpellCast(ctx *EffectContext) error {
	if !isSpellBeingCast(ctx) {
		return nil
	}
	promptPermanentSkillBuff(ctx, "选择1个你的法术")
	return nil
}
