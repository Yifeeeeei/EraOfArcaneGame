package game

import (
	"eraofarcane/model"
)

type Card3421012Earthshaker struct{ AlwaysActive }

func (Card3421012Earthshaker) ID() string { return "3421012" }

func (Card3421012Earthshaker) Name() string { return "石破天惊" }

func (Card3421012Earthshaker) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx.ExtraData["purpose"] != string(skillPurposeAttack) {
		return
	}
	player := ctx.Engine.State.Players[ctx.PlayerID]
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := player.Units[col][row]
			if unit == nil || unit.Card == nil || !unit.Card.IsCompanion() {
				continue
			}
			stats.PowerBonus += effectiveElementsGain(unit)[model.ElementEarth]
		}
	}
}
