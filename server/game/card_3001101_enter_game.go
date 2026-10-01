package game

import (
	"fmt"
)

type Card3001101EnterGame struct{ AlwaysActive }

func (Card3001101EnterGame) ID() string { return "3001101" }

func (Card3001101EnterGame) Name() string { return "入局" }

func (Card3001101EnterGame) OnSpellCast(ctx *EffectContext) error {
	candidates := enterGamePlayerCandidates(ctx)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "enter_game_player",
		"入局:选择召唤弃子的玩家", candidates, 1, 1,
		func(selected []string) {
			targetPlayerID, ok := enterGamePlayerIDFromSelection(firstSelected(selected))
			if !ok || targetPlayerID < 0 || targetPlayerID >= len(ctx.Engine.State.Players) {
				return
			}
			positions := ctx.Engine.emptyUnitPositionsForPlayer(targetPlayerID, ctx.PlayerID)
			if len(positions) == 0 {
				return
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "enter_game_position",
				"入局:选择弃子的召唤位置", positions, 1, 1,
				func(posSelected []string) {
					pos, ok := positionFromSelectionID(firstSelected(posSelected))
					if !ok || targetPlayerID < 0 || targetPlayerID >= len(ctx.Engine.State.Players) {
						return
					}
					ps := ctx.Engine.State.Players[targetPlayerID]
					if ps == nil || ps.Units[pos.Col][pos.Row] != nil {
						return
					}
					ctx.Engine.summonFreshCardAtPosition(targetPlayerID, "1001101", pos, true)
				})
		})
	return nil
}

func enterGamePlayerCandidates(ctx *EffectContext) []map[string]any {
	if ctx == nil || ctx.Engine == nil {
		return nil
	}
	candidates := make([]map[string]any, 0, len(ctx.Engine.State.Players))
	for playerID := range ctx.Engine.State.Players {
		if len(ctx.Engine.emptyUnitPositionsForPlayer(playerID, ctx.PlayerID)) == 0 {
			continue
		}
		side := "enemy"
		if playerID == ctx.PlayerID {
			side = "own"
		}
		candidates = append(candidates, map[string]any{
			"instance_id": fmt.Sprintf("player:%d", playerID),
			"name":        fmt.Sprintf("玩家%d", playerID+1),
			"zone":        "player",
			"side":        side,
			"player_id":   playerID,
		})
	}
	return candidates
}

func enterGamePlayerIDFromSelection(id string) (int, bool) {
	var playerID int
	if _, err := fmt.Sscanf(id, "player:%d", &playerID); err != nil {
		return 0, false
	}
	return playerID, true
}
