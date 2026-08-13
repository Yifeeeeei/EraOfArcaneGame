package game

type Card2021108PsychicDisk struct{ AlwaysActive }

func (Card2021108PsychicDisk) ID() string   { return "2021108" }
func (Card2021108PsychicDisk) Name() string { return "通灵盘" }

func (Card2021108PsychicDisk) ModifySkillUseCost(ctx *EffectContext, cost map[string]int) {
	if ctx == nil || ctx.Source == nil || ctx.Source.Card == nil || !hasCardTag(ctx.Source.Card, "灵媒") {
		return
	}
	reduceGenericCost(cost, ctx.Source.Card.Category, 1)
}

var _ SkillUseCostModifier = Card2021108PsychicDisk{}
