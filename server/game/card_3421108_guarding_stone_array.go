package game

import "eraofarcane/model"

type Card3421108GuardingStoneArray struct{ AlwaysActive }

func (Card3421108GuardingStoneArray) ID() string   { return "3421108" }
func (Card3421108GuardingStoneArray) Name() string { return "御守石阵" }

func (Card3421108GuardingStoneArray) ModifySkillUseCost(ctx *EffectContext, cost map[string]int) {
	if ctx == nil || ctx.Source == nil || ctx.Target != ctx.Source {
		return
	}
	if ctx.ExtraData == nil || ctx.ExtraData["purpose"] != string(skillPurposeDefend) {
		return
	}
	reduceCost(cost, model.ElementEarth, 1)
}
