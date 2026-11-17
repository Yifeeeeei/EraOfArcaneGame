package game

import (
	"eraofarcane/model"
)

func (Card2221010TideRune) CounterTriggers() []EffectTrigger { return []EffectTrigger{TriggerOnDraw} }

func (Card2221010TideRune) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Trigger == TriggerOnDraw && ctx.Event.PlayerID != ctx.Source.OwnerID && ctx.Event.DrawCount >= 3 &&
		len(ctx.Engine.friendlyUnits(ctx.Source.OwnerID, false, isWaterCompanion)) > 0
}

type Card2221010TideRune struct{ AlwaysActive }

func (Card2221010TideRune) ID() string { return "2221010" }

func (Card2221010TideRune) Name() string { return "潮涌符文" }

func (Card2221010TideRune) OnUseItem(ctx *EffectContext) error {
	targets := ctx.Engine.friendlyUnits(ctx.PlayerID, false, isWaterCompanion)
	if len(targets) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "tide_rune_buff",
		"潮涌符文:选择你的1个水纹伙伴获得负载+2水", targets, 1, 1,
		func(selected []string) {
			target := selectedUnitFromCandidates(ctx.Engine, selected, targets)
			if target != nil {
				ctx.Engine.addElementsGainBonus(target, ctx.PlayerID, model.ElementWater, 2, ctx.Source)
			}
		})
	return nil
}
