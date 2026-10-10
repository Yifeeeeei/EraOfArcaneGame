package game

type Card4611101BloodCountHubert struct{ AlwaysActive }

func (Card4611101BloodCountHubert) ID() string { return "4611101" }

func (Card4611101BloodCountHubert) Name() string { return "鲜血伯爵 休伯特 黑松" }

func (Card4611101BloodCountHubert) OnEnter(ctx *EffectContext) error {
	addSkillToPool(ctx, "3601101")
	return nil
}
