package game

import (
	"eraofarcane/model"
)

type Card1221013Raincaller struct{ AlwaysActive }

func (Card1221013Raincaller) ID() string { return "1221013" }

func (Card1221013Raincaller) Name() string { return "唤雨师" }

func (Card1221013Raincaller) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if ctx.Target == nil || ctx.Target.Card == nil {
		return
	}
	if ctx.Target.Card.Category == model.ElementWater || ctx.Target.Card.Category == model.ElementAir {
		stats.PowerBonus++
	}
}
