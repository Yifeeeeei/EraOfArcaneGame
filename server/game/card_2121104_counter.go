package game

import (
	"eraofarcane/model"
)

func (Card2121104FireRebirthScroll) CounterTriggers() []EffectTrigger {
	return []EffectTrigger{TriggerOnTurnEnd}
}

func (Card2121104FireRebirthScroll) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Trigger == TriggerOnTurnEnd && ctx.Event.PlayerID != ctx.Source.OwnerID &&
		len(ctx.Engine.rebirthScrollReviveCandidates(ctx.Source.OwnerID)) > 0 &&
		len(ctx.Engine.friendlyEmptyUnitPositions(ctx.Source.OwnerID)) > 0
}

type Card2121104FireRebirthScroll struct{ AlwaysActive }

func (Card2121104FireRebirthScroll) ID() string { return "2121104" }

func (Card2121104FireRebirthScroll) Name() string { return "浴火重生卷轴" }

func (Card2121104FireRebirthScroll) OnTurnEnd(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	candidates := ctx.Engine.rebirthScrollReviveCandidates(ctx.PlayerID)
	if len(candidates) == 0 {
		return nil
	}
	if len(ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)) == 0 {
		return nil
	}
	if len(candidates) == 1 {
		ctx.Engine.promptFireRebirthScrollPosition(ctx.PlayerID, ctx.Source, candidates[0].InstanceID)
		return nil
	}
	choices := make([]map[string]any, 0, len(candidates))
	for _, card := range candidates {
		info := cardToInfo(card)
		info["zone"] = "graveyard"
		info["side"] = "own"
		choices = append(choices, info)
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "fire_rebirth_scroll",
		"浴火重生卷轴:选择1个本回合死亡的友方火焰伙伴复活", choices, 1, 1,
		func(selected []string) {
			ctx.Engine.promptFireRebirthScrollPosition(ctx.PlayerID, ctx.Source, firstSelected(selected))
		})
	return nil
}

var _ OnTurnEndBehavior = Card2121104FireRebirthScroll{}

func (e *Engine) promptFireRebirthScrollPosition(playerID int, source *CardInstance, instanceID string) {
	positions := e.friendlyEmptyUnitPositions(playerID)
	if len(positions) == 0 || e.findRecentFireRebirthCandidate(playerID, instanceID) == nil {
		return
	}
	e.SetPendingAction(playerID, "fire_rebirth_scroll_position",
		"浴火重生卷轴:选择复活位置", positions, 1, 1,
		func(selected []string) {
			pos, ok := positionFromSelectionID(firstSelected(selected))
			if !ok {
				return
			}
			if e.reviveRecentFireCompanionAtPosition(playerID, instanceID, pos) {
				e.emit(GameEvent{
					Type:   "fire_rebirth_scroll_revive",
					Player: -1,
					Data: map[string]any{
						"player":   playerID,
						"source":   cardToInfo(source),
						"revived":  instanceID,
						"position": pos,
						"count":    1,
					},
				})
			}
		})
}

func (e *Engine) findRecentFireRebirthCandidate(playerID int, instanceID string) *CardInstance {
	for _, card := range e.rebirthScrollReviveCandidates(playerID) {
		if card != nil && card.InstanceID == instanceID {
			return card
		}
	}
	return nil
}

func (e *Engine) reviveRecentFireCompanionAtPosition(playerID int, instanceID string, pos Position) bool {
	card := e.findRecentFireRebirthCandidate(playerID, instanceID)
	if card == nil || !pos.Valid() || e.State.Players[playerID].Units[pos.Col][pos.Row] != nil {
		return false
	}
	if e.reviveCompanionFromGraveyardWithLifeAtPosition(playerID, card.InstanceID, 0, false, pos) {
		card.IsHorizontal = false
		if card.Statuses != nil {
			delete(card.Statuses, enteredGraveyardTurnStatus)
		}
		return true
	}
	return false
}

func (e *Engine) rebirthScrollReviveCandidates(playerID int) []*CardInstance {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	ps := e.State.Players[playerID]
	candidates := make([]*CardInstance, 0)
	for _, card := range ps.Graveyard {
		if card == nil || card.Card == nil || !card.Card.IsCompanion() || card.Card.Category != model.ElementFire {
			continue
		}
		if card.Statuses[enteredGraveyardTurnStatus] == e.State.TurnNumber {
			candidates = append(candidates, card)
		}
	}
	return candidates
}
