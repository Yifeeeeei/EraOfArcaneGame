package game

type Card3621105Intimidation struct{ AlwaysActive }

func (Card3621105Intimidation) ID() string { return "3621105" }

func (Card3621105Intimidation) Name() string { return "恐吓" }

func (Card3621105Intimidation) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Engine == nil {
		return
	}
	count := 0
	for _, skill := range ctx.Engine.ruleSources(ctx.Engine.State.Players[ctx.OpponentID]) {
		if skill == nil || skill.Card == nil || !isCombatSpellCard(skill.Card) {
			continue
		}
		if skill.Statuses[StatusWeaken] > 0 && ctx.Engine.hasEffectiveStatus(skill, StatusWeaken) {
			count++
		}
	}
	bonus := min(count, 2)
	stats.PowerBonus += bonus
	stats.DamageBonus += bonus
}
