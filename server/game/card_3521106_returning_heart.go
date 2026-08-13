package game

import "eraofarcane/model"

type Card3521106ReturningHeart struct{ AlwaysActive }

func (Card3521106ReturningHeart) ID() string   { return "3521106" }
func (Card3521106ReturningHeart) Name() string { return "归心" }

func (Card3521106ReturningHeart) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Engine == nil {
		return
	}
	for _, col := range ctx.Engine.State.Players[ctx.PlayerID].Units {
		for _, unit := range col {
			if unit == nil || unit.Card == nil || !unit.Card.IsCompanion() {
				continue
			}
			if unit.Card.Category == model.ElementLight {
				stats.PowerBonus++
			} else {
				stats.PowerBonus--
			}
		}
	}
}
