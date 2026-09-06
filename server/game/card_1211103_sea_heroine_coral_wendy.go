package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card1211103SeaHeroineCoralWendy struct{ AlwaysActive }

func (Card1211103SeaHeroineCoralWendy) ID() string { return "1211103" }

func (Card1211103SeaHeroineCoralWendy) Name() string { return "海上巾帼 珊瑚 雯迪" }

func (Card1211103SeaHeroineCoralWendy) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.Card == nil {
		return nil
	}
	if ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) || !isFriendlySpellCast(ctx) || totalElementCost(ctx.Target.Card.ElementsExpense) >= 3 {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	cost := map[string]int{model.ElementWater: 2}
	if !ctx.Engine.canPayCost(ps, cost) {
		return nil
	}
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "coral_wendy_reset_spell",
		"海上巾帼 珊瑚 雯迪:是否支付2水重置刚使用的法术", []map[string]any{candidateInfo(ctx.Target, "skill", "own")}, 0, 1, cost, false,
		func(selected []string, data map[string]any) error {
			if len(selected) == 0 {
				return nil
			}
			if selected[0] != ctx.Target.InstanceID || ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) {
				return fmt.Errorf("invalid coral wendy reset")
			}
			target := ctx.Engine.findSkill(ctx.Engine.State.Players[ctx.PlayerID], selected[0])
			if target == nil || target.Card == nil || !target.IsHorizontal || totalElementCost(target.Card.ElementsExpense) >= 3 {
				return fmt.Errorf("invalid coral wendy reset")
			}
			if !ctx.Engine.payCostForAction(ps, cost, ActionMessage{Data: data}) {
				return fmt.Errorf("invalid coral wendy payment")
			}
			target.IsHorizontal = false
			ctx.Source.UsedThisTurn++
			return nil
		})
	return nil
}
