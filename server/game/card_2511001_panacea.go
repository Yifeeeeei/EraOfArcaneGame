package game

import "eraofarcane/model"

type Card2511001Panacea struct{}

func (Card2511001Panacea) ID() string   { return "2511001" }
func (Card2511001Panacea) Name() string { return "万灵药" }

func (Card2511001Panacea) OnUseItem(ctx *EffectContext) error {
	choices := []map[string]any{
		{"instance_id": "heal", "number": "2511001", "name": "回复1个友方单位所有生命", "type": "选择", "zone": "choice", "side": "own"},
		{"instance_id": "draw", "number": "2511001", "name": "抽4张牌", "type": "选择", "zone": "choice", "side": "own"},
		{"instance_id": "gain", "number": "2511001", "name": "获得5无", "type": "选择", "zone": "choice", "side": "own"},
		{"instance_id": "reset", "number": "2511001", "name": "重置1个技能", "type": "选择", "zone": "choice", "side": "own"},
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "panacea_mode", "选择万灵药效果", choices, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			switch selected[0] {
			case "draw":
				drawn := ctx.Engine.State.Players[ctx.PlayerID].DrawCards(4)
				for _, card := range drawn {
					ctx.Engine.emit(GameEvent{Type: "draw_card", Player: ctx.PlayerID, Data: map[string]any{"card": cardToInfo(card)}})
				}
			case "gain":
				ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementArcane: 5})
			case "heal":
				targets := ctx.Engine.friendlyUnits(ctx.PlayerID, true, nil)
				ctx.Engine.SetPendingAction(ctx.PlayerID, "panacea_heal", "选择1个友方单位回复全部生命", targets, 1, 1,
					func(selected []string) {
						if len(selected) == 0 {
							return
						}
						target := ctx.Engine.findFieldCardByInstance(ctx.Engine.State.Players[ctx.PlayerID], selected[0])
						if target != nil {
							target.CurrentLife = target.Card.Life
						}
					})
			case "reset":
				targets := ctx.Engine.friendlySkills(ctx.PlayerID, nil)
				ctx.Engine.SetPendingAction(ctx.PlayerID, "panacea_reset", "选择1个技能重置", targets, 1, 1,
					func(selected []string) {
						if len(selected) == 0 {
							return
						}
						target := ctx.Engine.findFieldCardByInstance(ctx.Engine.State.Players[ctx.PlayerID], selected[0])
						if target != nil {
							target.IsHorizontal = false
						}
					})
			}
		})
	return nil
}
