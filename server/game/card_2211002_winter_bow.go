package game

import (
	"eraofarcane/model"
)

type Card2211002WinterBow struct{ AlwaysActive }

func (Card2211002WinterBow) ID() string { return "2211002" }

func (Card2211002WinterBow) Name() string { return "嗜魔弓 凛冬" }

func (Card2211002WinterBow) OnEnter(ctx *EffectContext) error {
	bindSkillToHost(ctx, "3201002")
	return nil
}

func (Card2211002WinterBow) OnSpellCast(ctx *EffectContext) error {
	winterBowPlayer := ctx.Engine.State.Players[ctx.PlayerID]
	if ctx.ExtraData == nil {
		return nil
	}
	if _, ok := ctx.ExtraData["cast_player"].(int); !ok {
		return nil
	}
	if !ctx.Engine.canPayCost(winterBowPlayer, map[string]int{model.ElementWater: 1}) {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "winter_bow_water_mark", "嗜魔弓 凛冬:是否支付1水放置1个水纹标记物", []map[string]any{candidateInfo(ctx.Source, "equipment", "own")}, 0, 1, func(selected []string) {
		if len(selected) == 0 {
			return
		}
		if ctx.Engine.payCostForAction(winterBowPlayer, map[string]int{model.ElementWater: 1}, ActionMessage{}) {
			ctx.Source.Statuses[winterBowWaterMark]++
		}
	})
	return nil
}
