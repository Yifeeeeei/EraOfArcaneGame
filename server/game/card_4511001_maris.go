package game

import "eraofarcane/model"

type Card4511001Maris struct{}

func (Card4511001Maris) ID() string   { return "4511001" }
func (Card4511001Maris) Name() string { return "圣使 玛丽斯 南森埃尔" }

func (Card4511001Maris) OnUltimate(ctx *EffectContext) error {
	ctx.Source.Statuses["玛丽斯反伤赐光"] = 1
	return nil
}

func (Card4511001Maris) OnDamaged(ctx *EffectContext) error {
	if ctx.Source.Statuses["玛丽斯反伤赐光"] <= 0 || ctx.ExtraData == nil {
		return nil
	}
	damagedPlayer, _ := ctx.ExtraData["damaged_player"].(int)
	if damagedPlayer != ctx.PlayerID {
		return nil
	}
	if attacker, ok := ctx.ExtraData["attacker"].(int); ok && attacker == ctx.PlayerID {
		return nil
	}
	ctx.Engine.State.Players[ctx.PlayerID].Elements[model.ElementLight] += 2
	return nil
}
