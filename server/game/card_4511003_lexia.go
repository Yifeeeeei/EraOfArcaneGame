package game

type Card4511003Lexia struct{ AlwaysActive }

func (Card4511003Lexia) ID() string { return "4511003" }

func (Card4511003Lexia) Name() string { return "骑士团长 蕾曦娅" }

func (Card4511003Lexia) OnTurnStart(ctx *EffectContext) error {
	if ctx.Source.Statuses["团结希望"] == 0 {
		if !replaceSkillInPool(ctx, "3521007", "3501001") {
			addSkillToPool(ctx, "3501001")
		}
		ctx.Source.Statuses["团结希望"] = 1
	}
	return nil
}
