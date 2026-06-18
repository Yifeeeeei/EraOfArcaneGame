package game

import (
	"fmt"

	"eraofarcane/model"
)

type Card3521010Dawn struct{ AlwaysActive }

func (Card3521010Dawn) ID() string   { return "3521010" }
func (Card3521010Dawn) Name() string { return "破晓" }

func (Card3521010Dawn) OnUnitEnter(ctx *EffectContext) error {
	if ctx.Target == nil || ctx.Target.Card == nil || !ctx.Target.Card.IsCompanion() {
		return nil
	}
	if ctx.ExtraData != nil {
		if enteredPlayer, ok := ctx.ExtraData["entered_player"].(int); ok && enteredPlayer != ctx.PlayerID {
			return nil
		}
	}
	if effectiveElementsGain(ctx.Target)[model.ElementLight] <= 0 {
		return nil
	}
	ctx.Source.PowerBonus++
	return nil
}

func (Card3521010Dawn) ValidateSkillUse(ctx *EffectContext, skill *CardInstance, purpose skillPurpose) error {
	if skill == nil || skill.Card == nil || skill.Card.Number != "3521010" || purpose != skillPurposeAttack {
		return nil
	}
	power := skill.Card.Power + skill.PowerBonus
	if ctx != nil && ctx.Engine != nil {
		power += ctx.Engine.temporarySpellPowerBonus(skill.OwnerID, skill)
	}
	if power <= 8 {
		return fmt.Errorf("Dawn can attack only when its power is greater than 8")
	}
	return nil
}
