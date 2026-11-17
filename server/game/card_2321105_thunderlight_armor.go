package game

type Card2321105ThunderlightArmor struct{ AlwaysActive }

func (Card2321105ThunderlightArmor) ID() string { return "2321105" }

func (Card2321105ThunderlightArmor) Name() string { return "雷光战铠" }

func (Card2321105ThunderlightArmor) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	skill := ctx.Target
	if skill == nil || skill.Card == nil || !isSpellLikeCard(skill.Card) {
		return
	}
	if !hasCardTag(skill.Card, "驱动") && !hasCardTag(skill.Card, "聚能") {
		return
	}
	count := 0
	for _, item := range ctx.Engine.State.Players[ctx.PlayerID].Equipment {
		if isThunderlightItem(item) {
			count++
		}
	}
	if count >= 3 {
		stats.PowerBonus += 2
	}
}
