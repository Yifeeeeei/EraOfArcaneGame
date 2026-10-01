package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card1111101SupremeQueenDailinCeltic struct{ AlwaysActive }

func (Card1111101SupremeQueenDailinCeltic) ID() string { return "1111101" }

func (Card1111101SupremeQueenDailinCeltic) Name() string { return "无上女王 黛琳 凯尔特" }

func (Card1111101SupremeQueenDailinCeltic) OnEnter(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Position == nil {
		return nil
	}
	markTemporaryDamageAndNegativeImmunity(ctx.Engine, ctx.Source)
	maxSummons := len(ctx.Engine.adjacentEmptyUnitPositions(ctx.PlayerID, *ctx.Source.Position))
	if maxSummons == 0 {
		return nil
	}
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementFire
	})
	if len(candidates) == 0 {
		return nil
	}
	if maxSummons > len(candidates) {
		maxSummons = len(candidates)
	}
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "supreme_queen_summon_cards",
		"无上女王 黛琳 凯尔特:选择火焰伙伴召唤到相邻位置", candidates, 0, maxSummons,
		nil, false, func(selected []string, _ map[string]any) error {
			ctx.Engine.continueSupremeQueenSummons(ctx, selected, 0)
			return nil
		})
	return nil
}

func markTemporaryDamageAndNegativeImmunity(e *Engine, card *CardInstance) {
	if e == nil || card == nil {
		return
	}
	card.Statuses[temporaryDamageAndNegativeImmunityUntilStatus] = e.State.TurnNumber + 1
}

func (e *Engine) adjacentEmptyUnitPositions(playerID int, center Position) []map[string]any {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) || !center.Valid() {
		return nil
	}
	ps := e.State.Players[playerID]
	candidates := make([]map[string]any, 0, 4)
	for _, delta := range []struct{ col, row int }{{0, -1}, {-1, 0}, {1, 0}, {0, 1}} {
		pos := Position{Col: center.Col + delta.col, Row: center.Row + delta.row}
		if !pos.Valid() || ps.Units[pos.Col][pos.Row] != nil {
			continue
		}
		candidates = append(candidates, map[string]any{
			"instance_id": positionSelectionID(pos),
			"name":        fmt.Sprintf("位置 (%d,%d)", pos.Col, pos.Row),
			"zone":        "unit_position",
			"side":        "own",
			"position":    pos,
		})
	}
	return candidates
}

func (e *Engine) continueSupremeQueenSummons(ctx *EffectContext, selected []string, index int) {
	if e == nil || ctx == nil || ctx.Source == nil || index >= len(selected) {
		return
	}
	card := e.findFriendlyHandCard(ctx.PlayerID, selected[index])
	if card == nil || card.Card == nil || !card.Card.IsCompanion() || card.Card.Category != model.ElementFire {
		e.continueSupremeQueenSummons(ctx, selected, index+1)
		return
	}
	if ctx.Source.Position == nil {
		return
	}
	positions := e.adjacentEmptyUnitPositions(ctx.PlayerID, *ctx.Source.Position)
	if len(positions) == 0 {
		return
	}
	e.SetPendingActionWithError(ctx.PlayerID, "supreme_queen_summon_position",
		fmt.Sprintf("无上女王 黛琳 凯尔特:选择%s的召唤位置", card.Card.Name), positions, 1, 1,
		nil, false, func(posSelected []string, _ map[string]any) error {
			pos, ok := positionFromSelectionID(firstSelected(posSelected))
			if !ok || ctx.Source.Position == nil || abs(pos.Col-ctx.Source.Position.Col)+abs(pos.Row-ctx.Source.Position.Row) != 1 {
				return fmt.Errorf("invalid queen summon position")
			}
			if e.State.Players[ctx.PlayerID].Units[pos.Col][pos.Row] != nil {
				return fmt.Errorf("queen summon position occupied")
			}
			card := e.removeFriendlyHandCard(ctx.PlayerID, selected[index])
			if card == nil || card.Card == nil || !card.Card.IsCompanion() || card.Card.Category != model.ElementFire {
				return fmt.Errorf("invalid queen summon card")
			}
			if !e.placeExistingCompanionAtPosition(ctx.PlayerID, card, pos, true) {
				return fmt.Errorf("queen summon failed")
			}
			markTemporaryDamageAndNegativeImmunity(e, card)
			if e.State.PendingAction != nil {
				e.continueAfterPendingAction(func() {
					e.continueSupremeQueenSummons(ctx, selected, index+1)
				})
				return nil
			}
			e.continueSupremeQueenSummons(ctx, selected, index+1)
			return nil
		})
}

func (e *Engine) removeFriendlyHandCard(playerID int, instanceID string) *CardInstance {
	if e == nil || instanceID == "" || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	ps := e.State.Players[playerID]
	for i, card := range ps.Hand {
		if card != nil && card.InstanceID == instanceID {
			ps.Hand = append(ps.Hand[:i], ps.Hand[i+1:]...)
			delete(ps.RevealedHand, card.InstanceID)
			return card
		}
	}
	return nil
}
