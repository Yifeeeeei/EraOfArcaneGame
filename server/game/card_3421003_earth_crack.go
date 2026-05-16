package game

type Card3421003EarthCrack struct{}

func (Card3421003EarthCrack) ID() string   { return "3421003" }
func (Card3421003EarthCrack) Name() string { return "裂地重击" }

func (Card3421003EarthCrack) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	load := totalLoad(ctx.Source)
	if load >= 1 {
		stats.PowerBonus++
		stats.DamageBonus++
	}
	if load >= 3 {
		stats.PowerBonus++
		stats.DamageBonus++
	}
}
