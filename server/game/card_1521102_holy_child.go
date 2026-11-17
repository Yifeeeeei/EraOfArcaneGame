package game

import (
	"eraofarcane/model"
)

type Card1521102HolyChild struct{ AlwaysActive }

func (Card1521102HolyChild) ID() string { return "1521102" }

func (Card1521102HolyChild) Name() string { return "神圣之子" }

func (Card1521102HolyChild) OnLoadGain(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target != ctx.Source {
		return nil
	}
	ctx.Engine.triggerHolyChildBonus(ctx.PlayerID, ctx.Source)
	return nil
}

func (Card1521102HolyChild) OnLifeGain(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target != ctx.Source {
		return nil
	}
	ctx.Engine.triggerHolyChildBonus(ctx.PlayerID, ctx.Source)
	return nil
}

var _ OnLoadGainBehavior = Card1521102HolyChild{}

var _ OnLifeGainBehavior = Card1521102HolyChild{}

const holyChildResolvingStatus = "holy_child_resolving"

func (e *Engine) triggerHolyChildBonus(playerID int, child *CardInstance) {
	if e == nil || child == nil || child.Card == nil || child.Card.Number != "1521102" || child.UltimateUsed || child.Statuses[holyChildResolvingStatus] > 0 {
		return
	}
	choices := []map[string]any{
		{"instance_id": "gain_light_load", "name": "额外获得负载+1光", "zone": "choice", "side": "own"},
		{"instance_id": "gain_life", "name": "额外获得+1血", "zone": "choice", "side": "own"},
	}
	e.SetPendingAction(playerID, "holy_child_bonus",
		"神圣之子:选择额外获得负载或生命", choices, 1, 1,
		func(selected []string) {
			if child == nil || child.CurrentLife <= 0 || child.UltimateUsed || child.Statuses[holyChildResolvingStatus] > 0 {
				return
			}
			child.UltimateUsed = true
			child.Statuses[holyChildResolvingStatus] = 1
			defer delete(child.Statuses, holyChildResolvingStatus)
			switch firstSelected(selected) {
			case "gain_life":
				e.gainLife(child, 1, child)
				e.emit(GameEvent{
					Type:   "life_gain",
					Player: -1,
					Data: map[string]any{
						"player": playerID,
						"target": cardToInfo(child),
						"amount": 1,
					},
				})
			default:
				e.addElementsGainBonus(child, playerID, model.ElementLight, 1, child)
			}
		})
}
