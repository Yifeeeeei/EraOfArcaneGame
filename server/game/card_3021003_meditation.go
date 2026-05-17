package game

import "eraofarcane/model"

type Card3021003Meditation struct{ AlwaysActive }

func (Card3021003Meditation) ID() string   { return "3021003" }
func (Card3021003Meditation) Name() string { return "冥想" }

func (Card3021003Meditation) SpellElementGains(*EffectContext) map[string]int {
	return map[string]int{model.ElementArcane: 1}
}
