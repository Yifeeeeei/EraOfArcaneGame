package game

import (
	"eraofarcane/model"
)

type Card2411001AncientTreeHeart struct{ AlwaysActive }

func (Card2411001AncientTreeHeart) ID() string { return "2411001" }

func (Card2411001AncientTreeHeart) Name() string { return "古树之心" }

func (Card2411001AncientTreeHeart) OnLoadGain(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil {
		return nil
	}
	target, _ := ctx.ExtraData["load_gain_target"].(*CardInstance)
	if target == nil || target.OwnerID != ctx.PlayerID || target.Position == nil {
		return nil
	}
	ctx.Engine.SetTriggeredTurnAction(ctx.Source, ctx.PlayerID, "ancient_tree_heart_life",
		"古树之心:是否使获得负载的友方单位获得+1血", []map[string]any{candidateInfo(target, "unit", "own")}, 0, 1,
		func(selected []string) {
			accepted := len(selected) > 0 && selected[0] == target.InstanceID && target.OwnerID == ctx.PlayerID && target.Position != nil
			if accepted && useTriggeredTurn(ctx.Source) {
				ctx.Engine.gainLife(target, 1, ctx.Source)
			}
		})
	return nil
}

func (Card2411001AncientTreeHeart) OnLifeGain(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil {
		return nil
	}
	target, _ := ctx.ExtraData["life_gain_target"].(*CardInstance)
	if target == nil || target.OwnerID != ctx.PlayerID || target.Position == nil {
		return nil
	}
	ctx.Engine.SetTriggeredTurnAction(ctx.Source, ctx.PlayerID, "ancient_tree_heart_load",
		"古树之心:是否使获得生命的友方单位负载+1地", []map[string]any{candidateInfo(target, "unit", "own")}, 0, 1,
		func(selected []string) {
			accepted := len(selected) > 0 && selected[0] == target.InstanceID && target.OwnerID == ctx.PlayerID && target.Position != nil
			if accepted && useTriggeredTurn(ctx.Source) {
				ctx.Engine.addElementsGainBonus(target, ctx.PlayerID, model.ElementEarth, 1, ctx.Source)
			}
		})
	return nil
}
