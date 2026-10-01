package game

import (
	"eraofarcane/model"
)

type Card3321108CallSpiritGoshawk struct{ AlwaysActive }

func (Card3321108CallSpiritGoshawk) ID() string { return "3321108" }

func (Card3321108CallSpiritGoshawk) Name() string { return "唤灵术 苍鹰" }

func (Card3321108CallSpiritGoshawk) OnEnter(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlySkills(ctx.PlayerID, func(skill *CardInstance) bool {
		return skill != nil && skill.Card != nil && skill.Card.Category == model.ElementAir
	})
	if len(candidates) == 0 {
		return nil
	}
	applyBuff := func(targetID string) {
		if targetID == "" {
			return
		}
		ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
			Type:             TempModSkillPowerBonus,
			SourceCardNumber: ctx.Source.Card.Number,
			SourceName:       ctx.Source.Card.Name,
			TargetInstanceID: targetID,
			Amount:           1,
			RemainingUses:    1,
		})
		ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
			Type:             TempModNextSkillUseAttackBonus,
			SourceCardNumber: ctx.Source.Card.Number,
			SourceName:       ctx.Source.Card.Name,
			TargetInstanceID: targetID,
			Amount:           1,
			RemainingUses:    1,
		})
	}
	if len(candidates) == 1 {
		id, _ := candidates[0]["instance_id"].(string)
		applyBuff(id)
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "goshawk_air_skill_buff",
		"唤灵术 苍鹰:选择1个友方大气法术下一次使用时+1攻+1威", candidates, 1, 1,
		func(selected []string) {
			applyBuff(firstSelected(selected))
		})
	return nil
}
