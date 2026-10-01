package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card1221112WaterMage struct{ AlwaysActive }

func (Card1221112WaterMage) ID() string { return "1221112" }

func (Card1221112WaterMage) Name() string { return "水魔导师" }

func (Card1221112WaterMage) OnUltimate(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlySkillsIncludingBound(ctx.PlayerID, func(skill *CardInstance) bool {
		return skill != nil && skill.Card != nil && skill.Card.IsSkill() &&
			skill.Card.Category == model.ElementWater &&
			totalElementCost(skill.Card.ElementsExpense) < 3 &&
			skill.IsHorizontal
	})
	if len(candidates) == 0 {
		return nil
	}
	allowed := make(map[string]bool, len(candidates))
	for _, candidate := range candidates {
		if id, _ := candidate["instance_id"].(string); id != "" {
			allowed[id] = true
		}
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "water_mage_reset_skill",
		"水魔导师:选择1个使用花费小于3的水纹法术重置", candidates, 1, 1,
		func(selected []string) {
			id := firstSelected(selected)
			if !allowed[id] {
				return
			}
			skill := ctx.Engine.findSkill(ctx.Engine.State.Players[ctx.PlayerID], id)
			if skill != nil && skill.Card != nil && skill.Card.IsSkill() && skill.Card.Category == model.ElementWater && totalElementCost(skill.Card.ElementsExpense) < 3 {
				skill.IsHorizontal = false
			}
		})
	return nil
}

func (Card1221112WaterMage) ValidateAbility(ctx *EffectContext, trigger EffectTrigger) error {
	if trigger != TriggerUltimate {
		return nil
	}
	if !ctx.Engine.hasResettableWaterSpell(ctx.PlayerID) {
		return fmt.Errorf("水魔导师需要1个已横置且使用花费小于3的水纹法术")
	}
	return nil
}

func (e *Engine) hasResettableWaterSpell(playerID int) bool {
	return len(e.friendlySkillsIncludingBound(playerID, func(skill *CardInstance) bool {
		return skill != nil && skill.Card != nil && skill.Card.IsSkill() &&
			skill.Card.Category == model.ElementWater &&
			totalElementCost(skill.Card.ElementsExpense) < 3 &&
			skill.IsHorizontal
	})) > 0
}
