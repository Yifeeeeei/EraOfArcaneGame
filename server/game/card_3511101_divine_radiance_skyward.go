package game

import (
	"fmt"
)

type Card3511101DivineRadianceSkyward struct{ AlwaysActive }

func (Card3511101DivineRadianceSkyward) ID() string { return "3511101" }

func (Card3511101DivineRadianceSkyward) Name() string { return "神辉驭空" }

func (Card3511101DivineRadianceSkyward) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.Source.Card.Number != "3511101" {
		return
	}
	stats.PowerBonus += len(ctx.Engine.State.Players[ctx.OpponentID].Hand)
}

func (Card3511101DivineRadianceSkyward) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.Source.Card.Number != "3511101" {
		return nil
	}
	candidates := []map[string]any{
		{"instance_id": "player:own", "name": fmt.Sprintf("玩家%d", ctx.PlayerID+1), "zone": "player", "side": "own", "player_id": ctx.PlayerID},
		{"instance_id": "player:opponent", "name": fmt.Sprintf("玩家%d", ctx.OpponentID+1), "zone": "player", "side": "enemy", "player_id": ctx.OpponentID},
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "divine_radiance_reset_hand",
		"神辉驭空:选择1名玩家弃掉全部手牌并抽至手牌上限", candidates, 0, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			targetPlayer := ctx.PlayerID
			if firstSelected(selected) == "player:opponent" {
				targetPlayer = ctx.OpponentID
			}
			count := ctx.Engine.discardAllHandCards(targetPlayer)
			limit := ctx.Engine.handLimitForPlayer(ctx.Engine.State.Players[targetPlayer])
			if limit > 0 {
				ctx.Engine.drawCards(targetPlayer, limit)
			}
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source":    cardToInfo(ctx.Source),
				"effect":    "divine_radiance_reset_hand",
				"target":    targetPlayer,
				"discarded": count,
				"drawn":     limit,
			}})
		})
	return nil
}
