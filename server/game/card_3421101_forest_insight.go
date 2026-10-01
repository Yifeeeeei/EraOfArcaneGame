package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card3421101ForestInsight struct{ AlwaysActive }

func (Card3421101ForestInsight) ID() string { return "3421101" }

func (Card3421101ForestInsight) Name() string { return "森之洞察" }

func (Card3421101ForestInsight) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || !isFriendlySpellCast(ctx) || !isSpellBeingCast(ctx) {
		return nil
	}
	if ctx.Source.Card.Number != "3421101" {
		return nil
	}
	drawCount := min(5, countFriendlyEarthCompanions(ctx.Engine, ctx.PlayerID))
	if drawCount <= 0 {
		return nil
	}
	drawn := ctx.Engine.drawCards(ctx.PlayerID, drawCount)
	shuffleBackCount := len(drawn)
	if shuffleBackCount <= 0 {
		return nil
	}
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, nil)
	if len(candidates) < shuffleBackCount {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "forest_insight_shuffle_hand",
		fmt.Sprintf("森之洞察:选择%d张手牌洗回卡组", shuffleBackCount), candidates, shuffleBackCount, shuffleBackCount,
		func(selected []string) {
			if moveSelectedHandCardsToDeck(ctx.Engine, ctx.PlayerID, selected, shuffleBackCount) > 0 {
				ctx.Engine.shuffleDeck(ctx.PlayerID)
			}
		})
	return nil
}

func countFriendlyEarthCompanions(e *Engine, playerID int) int {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return 0
	}
	count := 0
	for _, card := range e.getAllFieldCards(e.State.Players[playerID]) {
		if card != nil && card.Card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementEarth {
			count++
		}
	}
	return count
}

func moveSelectedHandCardsToDeck(e *Engine, playerID int, selected []string, maxCount int) int {
	if e == nil || maxCount <= 0 || playerID < 0 || playerID >= len(e.State.Players) {
		return 0
	}
	ps := e.State.Players[playerID]
	selectedSet := make(map[string]bool, len(selected))
	for _, id := range selected {
		selectedSet[id] = true
	}
	moved := 0
	for i := 0; i < len(ps.Hand) && moved < maxCount; {
		card := ps.Hand[i]
		if card != nil && selectedSet[card.InstanceID] {
			ps.Hand = append(ps.Hand[:i], ps.Hand[i+1:]...)
			resetCardForHiddenZone(card)
			ps.Deck = append(ps.Deck, card)
			moved++
			continue
		}
		i++
	}
	return moved
}
