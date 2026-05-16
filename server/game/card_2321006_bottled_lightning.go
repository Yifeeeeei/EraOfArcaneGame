package game

import "eraofarcane/model"

type Card2321006BottledLightning struct{}

func (Card2321006BottledLightning) ID() string   { return "2321006" }
func (Card2321006BottledLightning) Name() string { return "瓶中闪电" }

func (Card2321006BottledLightning) OnUseItem(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	ps.Elements[model.ElementAir] += 3
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source":  cardToInfo(ctx.Source),
		"effect":  "gain_element",
		"element": model.ElementAir,
		"amount":  3,
	}})
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "bottled_lightning_stun",
		"选择1个友方单位获得晕眩2", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, selected[0])
			if target == nil || zone != "unit" {
				return
			}
			target.Statuses[StatusStun] += 2
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source": cardToInfo(ctx.Source),
				"target": cardToInfo(target),
				"effect": "apply_status",
				"status": StatusStun,
				"amount": 2,
			}})
		})
	return nil
}
