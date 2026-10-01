package game

type Card4111001Longjuanhuo struct{ AlwaysActive }

func (Card4111001Longjuanhuo) ID() string { return "4111001" }

func (Card4111001Longjuanhuo) Name() string { return "掌门 龙卷火" }

func (Card4111001Longjuanhuo) OnTurnStart(ctx *EffectContext) error {
	if ctx.Source.Statuses["开局技能"] == 0 {
		addSkillToPool(ctx, "3101002")
		ctx.Source.Statuses["开局技能"] = 1
	}
	return nil
}
