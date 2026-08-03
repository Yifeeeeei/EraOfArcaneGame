package game

import "eraofarcane/model"

type Card4311003Muling struct{ AlwaysActive }

func (Card4311003Muling) ID() string   { return "4311003" }
func (Card4311003Muling) Name() string { return "掌门 穆伶" }

func (Card4311003Muling) OnUltimate(ctx *EffectContext) error {
	ownCandidates := mulingCandidates(ctx.Engine, ctx.PlayerID, ctx.PlayerID, "own")
	enemyCandidates := mulingCandidates(ctx.Engine, ctx.PlayerID, ctx.OpponentID, "enemy")
	if len(ownCandidates) == 0 || len(enemyCandidates) == 0 {
		return nil
	}

	candidates := append(ownCandidates, enemyCandidates...)
	ctx.Engine.SetPendingAction(ctx.PlayerID, "muling_return",
		"选择双方各1个伙伴，支付入场费用差值的大气，将它们移回手牌",
		candidates, 2, 2,
		func(selected []string) {
			own := selectedUnitFromCandidates(ctx.Engine, selected, ownCandidates)
			enemy := selectedUnitFromCandidates(ctx.Engine, selected, enemyCandidates)
			if own == nil || enemy == nil {
				return
			}
			diff := absInt(own.Card.TotalCost() - enemy.Card.TotalCost())
			cost := map[string]int{model.ElementAir: diff}
			ps := ctx.Engine.State.Players[ctx.PlayerID]
			if !ctx.Engine.payCostForAction(ps, cost, ActionMessage{}) {
				ctx.Engine.emit(GameEvent{
					Type:   "error",
					Player: ctx.PlayerID,
					Data:   map[string]any{"message": "穆伶绝技费用不足"},
				})
				return
			}
			ctx.Engine.returnUnitToHand(own, ctx.PlayerID)
			ctx.Engine.returnUnitToHand(enemy, ctx.OpponentID)
		})
	return nil
}

func mulingCandidates(e *Engine, playerID int, ownerID int, side string) []map[string]any {
	ps := e.State.Players[ownerID]
	candidates := make([]map[string]any, 0)
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := ps.Units[col][row]
			if unit == nil || unit.Card.IsHero() || !unit.Card.IsCompanion() {
				continue
			}
			if ownerID != playerID && !e.IsInSpellRange(playerID, col, row, false) {
				continue
			}
			info := cardToInfo(unit)
			info["side"] = side
			candidates = append(candidates, info)
		}
	}
	return candidates
}

func findSelectedUnit(ps *PlayerState, selected []string) *CardInstance {
	selectedIDs := make(map[string]bool, len(selected))
	for _, id := range selected {
		selectedIDs[id] = true
	}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := ps.Units[col][row]
			if unit != nil && selectedIDs[unit.InstanceID] {
				return unit
			}
		}
	}
	return nil
}

func (e *Engine) returnUnitToHand(unit *CardInstance, ownerID int) {
	ps := e.State.Players[ownerID]
	if unit.Position != nil {
		ps.Units[unit.Position.Col][unit.Position.Row] = nil
	}
	unit.Position = nil
	unit.IsHorizontal = true
	unit.Statuses = make(map[string]int)
	e.addCardToHand(ownerID, unit)
	e.emit(GameEvent{
		Type:   "unit_returned",
		Player: -1,
		Data: map[string]any{
			"player": ownerID,
			"card":   cardToInfo(unit),
		},
	})
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
