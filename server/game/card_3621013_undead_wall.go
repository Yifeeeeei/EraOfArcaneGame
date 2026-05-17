package game

const recentFriendlyDeathStatus = "recent_friendly_death"

type Card3621013UndeadWall struct{ AlwaysActive }

func (Card3621013UndeadWall) ID() string   { return "3621013" }
func (Card3621013UndeadWall) Name() string { return "亡灵护壁" }

func (Card3621013UndeadWall) OnFriendlyDeath(ctx *EffectContext) error {
	if ctx.Source == nil {
		return nil
	}
	ctx.Source.Statuses[recentFriendlyDeathStatus] = 2
	return nil
}

func (Card3621013UndeadWall) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx.ExtraData["purpose"] != string(skillPurposeDefend) {
		return
	}
	if ctx.Source == nil || ctx.Source.Statuses[recentFriendlyDeathStatus] <= 0 {
		return
	}
	stats.PowerBonus += 2
}
