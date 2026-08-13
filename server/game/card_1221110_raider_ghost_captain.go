package game

import "eraofarcane/model"

type Card1221110RaiderGhostCaptain struct{ AlwaysActive }

func (Card1221110RaiderGhostCaptain) ID() string   { return "1221110" }
func (Card1221110RaiderGhostCaptain) Name() string { return "掠夺者幽灵船长" }

func (Card1221110RaiderGhostCaptain) ModifyElementsGain(ctx *EffectContext, target *CardInstance, gains map[string]int) {
	if ctx == nil || ctx.Target == nil || target == nil || target == ctx.Target {
		return
	}
	if target.OwnerID != ctx.PlayerID || !isRaiderCompanion(target) {
		return
	}
	gains[model.ElementWater]++
}
