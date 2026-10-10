package game

import (
	"eraofarcane/model"
	"fmt"
)

const arcaneSealExtraUseCostStatus = "奥术封印额外使用费用"

type Card3021108ArcaneSeal struct{ AlwaysActive }

func (Card3021108ArcaneSeal) ID() string { return "3021108" }

func (Card3021108ArcaneSeal) Name() string { return "奥术封印" }

func (Card3021108ArcaneSeal) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	candidates := ctx.Engine.enemySkills(ctx.PlayerID, nil)
	if len(candidates) == 0 {
		return nil
	}
	source := ctx.Source
	ctx.Engine.SetPendingAction(ctx.PlayerID, "arcane_seal_skill",
		"奥术封印:选择敌方场上1个技能封印到其下个回合结束", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			target := findSkillSlotCard(ctx.Engine.State.Players[ctx.OpponentID], selected[0])
			if target == nil {
				return
			}
			target.Statuses[StatusSeal] = 1
			source.Statuses[arcaneSealExtraUseCostStatus]++
		})
	return nil
}

func (Card3021108ArcaneSeal) ModifySkillUseCost(ctx *EffectContext, cost map[string]int) {
	if ctx == nil || ctx.Source == nil || ctx.Target != ctx.Source {
		return
	}
	if extra := ctx.Source.Statuses[arcaneSealExtraUseCostStatus]; extra > 0 {
		cost[model.ElementArcane] += 2 * extra
	}
}

var _ OnSpellCastBehavior = Card3021108ArcaneSeal{}

var _ SkillUseCostModifier = Card3021108ArcaneSeal{}

func (Card3021108ArcaneSeal) PrepareSpellCast(ctx *EffectContext, _ SpellTarget, _ ActionMessage) (SpellCastOptions, error) {
	if len(ctx.Engine.enemySkills(ctx.PlayerID, nil)) == 0 {
		return SpellCastOptions{}, fmt.Errorf("arcane seal requires an enemy skill")
	}
	return SpellCastOptions{}, nil
}
