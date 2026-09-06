package game

import (
	"fmt"
)

type Card1021106SkyCityTycoon struct{ AlwaysActive }

func (Card1021106SkyCityTycoon) ID() string { return "1021106" }

func (Card1021106SkyCityTycoon) Name() string { return "云霄城富豪" }

func (Card1021106SkyCityTycoon) PerTurnLabel(*CardInstance) string {
	return "主动"
}

func (Card1021106SkyCityTycoon) OnPerTurn(ctx *EffectContext) error {
	if !ctx.Engine.canConsumeCard(ctx.Source) {
		return fmt.Errorf("云霄城富豪不能被消耗")
	}
	ctx.Source.IsHorizontal = true
	choices := []map[string]any{
		{"instance_id": "self_first", "number": "1021106", "name": "你先抽", "type": "选择", "zone": "choice", "side": "own"},
		{"instance_id": "opponent_first", "number": "1021106", "name": "对手先抽", "type": "选择", "zone": "choice", "side": "own"},
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "sky_city_tycoon_draw_order",
		"云霄城富豪:选择双方抽牌次序", choices, 1, 1,
		func(selected []string) {
			if firstSelected(selected) == "opponent_first" {
				ctx.Engine.drawCards(ctx.OpponentID, 1)
				ctx.Engine.drawCards(ctx.PlayerID, 1)
				return
			}
			ctx.Engine.drawCards(ctx.PlayerID, 1)
			ctx.Engine.drawCards(ctx.OpponentID, 1)
		})
	return nil
}
