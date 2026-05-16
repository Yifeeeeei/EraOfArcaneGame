package game

type Card1121004CelticLion struct{}

func (Card1121004CelticLion) ID() string   { return "1121004" }
func (Card1121004CelticLion) Name() string { return "凯尔特雄狮" }

func (Card1121004CelticLion) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	stats.PowerBonus++
}
