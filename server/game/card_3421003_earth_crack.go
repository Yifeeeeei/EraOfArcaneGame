package game

type Card3421003EarthCrack struct{ AlwaysActive }

func (Card3421003EarthCrack) ID() string      { return "3421003" }
func (Card3421003EarthCrack) Name() string    { return "裂地重击" }
func (Card3421003EarthCrack) MasteryMax() int { return 3 }
func (Card3421003EarthCrack) OnMastery(ctx *EffectContext, level int) error {
	return nil
}

func (Card3421003EarthCrack) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	mastery := ctx.Source.Statuses[StatusMastery]
	if mastery >= 1 {
		stats.PowerBonus++
		stats.DamageBonus++
	}
	if mastery >= 3 {
		stats.PowerBonus++
		stats.DamageBonus++
	}
}
