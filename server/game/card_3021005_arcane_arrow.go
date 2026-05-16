package game

type Card3021005ArcaneArrow struct{}

func (Card3021005ArcaneArrow) ID() string   { return "3021005" }
func (Card3021005ArcaneArrow) Name() string { return "奥术箭矢" }
func (Card3021005ArcaneArrow) SpellDamage(*EffectContext) int {
	return 1
}
