package game

type Card3321102KillingWind struct{ AlwaysActive }

func (Card3321102KillingWind) ID() string   { return "3321102" }
func (Card3321102KillingWind) Name() string { return "肃杀之风" }
func (Card3321102KillingWind) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Engine == nil {
		return
	}
	ownHand := len(ctx.Engine.State.Players[ctx.PlayerID].Hand)
	enemyHand := len(ctx.Engine.State.Players[ctx.OpponentID].Hand)
	stats.PowerBonus += abs(ownHand - enemyHand)
}
