package game

type Card2621101BlackPineWand struct{ AlwaysActive }

func (Card2621101BlackPineWand) ID() string { return "2621101" }

func (Card2621101BlackPineWand) Name() string { return "黑松木魔杖" }

func (Card2621101BlackPineWand) ModifySkillUseCost(ctx *EffectContext, cost map[string]int) {
	if ctx == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.ExtraData == nil {
		return
	}
	target, _ := ctx.ExtraData["spell_target"].(SpellTarget)
	targetUnit, _ := ctx.ExtraData["spell_target_unit"].(*CardInstance)
	if target.Type != "unit" || targetUnit == nil || targetUnit.OwnerID != ctx.PlayerID || ctx.Source.Card.Category == "" {
		return
	}
	reduceCost(cost, ctx.Source.Card.Category, 1)
}

var _ SkillUseCostModifier = Card2621101BlackPineWand{}
