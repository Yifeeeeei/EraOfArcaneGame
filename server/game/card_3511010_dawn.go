package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card3511010Dawn struct{ AlwaysActive }

func (Card3511010Dawn) ID() string { return "3511010" }

func (Card3511010Dawn) Name() string { return "破晓" }

func (Card3511010Dawn) OnUnitEnter(ctx *EffectContext) error {
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

func (Card3511010Dawn) ValidateSkillUse(ctx *EffectContext, skill *CardInstance, purpose skillPurpose) error {
	if skill == nil || skill.Card == nil || skill.Card.Number != "3511010" || purpose != skillPurposeAttack {
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

func (Card3511010Dawn) AffectedSpellUnits(ctx *EffectContext, defenderID int, target SpellTarget) ([]*CardInstance, bool) {
	if defenderID == ctx.PlayerID || !target.Position.Valid() {
		return nil, false
	}
	defender := ctx.Engine.State.Players[defenderID]
	targetUnit := defender.Units[target.Position.Col][target.Position.Row]
	if targetUnit == nil || targetUnit.Card == nil || !targetUnit.Card.IsCompanion() {
		return nil, false
	}
	units := make([]*CardInstance, 0, 9)
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := defender.Units[col][row]
			if unit != nil && unit.Card != nil && unit.Card.Category == targetUnit.Card.Category {
				units = append(units, unit)
			}
		}
	}
	return units, true
}
