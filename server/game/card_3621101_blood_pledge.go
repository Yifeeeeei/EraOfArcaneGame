package game

import (
	"eraofarcane/model"
)

type Card3621101BloodPledge struct{ AlwaysActive }

func (Card3621101BloodPledge) ID() string { return "3621101" }

func (Card3621101BloodPledge) Name() string { return "歃血" }

func (Card3621101BloodPledge) OnSpellHit(ctx *EffectContext) error {
	if ctx.Source == nil || ctx.ExtraData == nil {
		return nil
	}
	friendlyDamage, _ := ctx.ExtraData["actual_friendly_damage_by_instance"].(map[string]int)
	totalDamage := 0
	for _, amount := range friendlyDamage {
		totalDamage += amount
	}
	if totalDamage <= 0 {
		return nil
	}
	ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementShadow: 2})
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModNextAttackSpellPowerBonus,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		Amount:           2,
		RemainingUses:    1,
	})
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModNextSkillUseAttackBonus,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		Amount:           1,
		RemainingUses:    1,
	})
	return nil
}
