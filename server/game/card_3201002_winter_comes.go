package game

import (
	"fmt"
)

type Card3201002WinterComes struct{ AlwaysActive }

func (Card3201002WinterComes) ID() string { return "3201002" }

func (Card3201002WinterComes) Name() string { return "凛冬将至" }

func (Card3201002WinterComes) ValidateSkillUse(ctx *EffectContext, skill *CardInstance, purpose skillPurpose) error {
	for _, equipment := range ctx.Engine.State.Players[ctx.PlayerID].Equipment {
		if equipment != nil && equipment.Card.Number == "2211002" && equipment.Statuses[winterBowWaterMark] >= 5 {
			return nil
		}
	}
	return fmt.Errorf("winter comes requires winter bow with 5 water counters")
}

func (Card3201002WinterComes) OnSpellCast(ctx *EffectContext) error {
	for _, equipment := range ctx.Engine.State.Players[ctx.PlayerID].Equipment {
		if equipment != nil && equipment.Card.Number == "2211002" && equipment.Statuses[winterBowWaterMark] >= 5 {
			equipment.Statuses[winterBowWaterMark] -= 5
			return nil
		}
	}
	return nil
}
