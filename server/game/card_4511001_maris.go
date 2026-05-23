package game

import "eraofarcane/model"

type Card4511001Maris struct{ AlwaysActive }

func (Card4511001Maris) ID() string   { return "4511001" }
func (Card4511001Maris) Name() string { return "圣使 玛丽斯 南森埃尔" }

func (Card4511001Maris) OnDamaged(ctx *EffectContext) error {
	if ctx == nil || ctx.Source == nil || ctx.ExtraData == nil {
		return nil
	}
	damagedPlayer, _ := ctx.ExtraData["damaged_player"].(int)
	if damagedPlayer != ctx.PlayerID {
		return nil
	}
	if attacker, ok := ctx.ExtraData["attacker"].(int); ok && attacker == ctx.PlayerID {
		return nil
	}
	choice := cardToInfo(ctx.Source)
	choice["name"] = "使用玛丽斯: 获得2光辉元素"
	choice["zone"] = "choice"
	candidates := []map[string]any{choice}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "maris_gain_light",
		"玛丽斯: 是否获得2光辉元素?", candidates, 0, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			ctx.Engine.State.Players[ctx.PlayerID].Elements[model.ElementLight] += 2
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source": cardToInfo(ctx.Source),
				"effect": "gain_element",
				"elem":   model.ElementLight,
				"amount": 2,
			}})
		})
	return nil
}
