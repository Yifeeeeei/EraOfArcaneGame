package game

import (
	"eraofarcane/model"
	"fmt"
)

func (e *Engine) findFriendlyHandCard(playerID int, instanceID string) *CardInstance {
	if e == nil || instanceID == "" || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	for _, card := range e.State.Players[playerID].Hand {
		if card != nil && card.InstanceID == instanceID {
			return card
		}
	}
	return nil
}

func (e *Engine) offerDiscardedSpeckledSparrowSummon(playerID int, card *CardInstance) {
	if e == nil || card == nil || card.Card == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	ps := e.State.Players[playerID]
	cost := map[string]int{model.ElementAir: 1}
	if !e.canPayCost(ps, cost) {
		return
	}
	positions := e.friendlyEmptyUnitPositions(playerID)
	if len(positions) == 0 {
		return
	}
	cardID := card.InstanceID
	e.SetPendingActionWithError(playerID, "speckled_sparrow_discard_summon",
		"花斑麻雀:是否支付1气召唤被弃置的此卡", positions, 0, 1, cost, false,
		func(selected []string, data map[string]any) error {
			if len(selected) == 0 {
				return nil
			}
			pos, ok := positionFromSelectionID(firstSelected(selected))
			if !ok || !pos.Valid() || e.State.Players[playerID].Units[pos.Col][pos.Row] != nil {
				return fmt.Errorf("invalid speckled sparrow position")
			}
			if !e.payCostForAction(ps, cost, ActionMessage{Data: data}) {
				return fmt.Errorf("invalid speckled sparrow payment")
			}
			if !e.reviveCompanionFromGraveyardWithLifeAtPosition(playerID, cardID, 0, false, pos) {
				return fmt.Errorf("invalid speckled sparrow summon")
			}
			return nil
		})
}

func deckHasMatch(ps *PlayerState, predicate func(*CardInstance) bool) bool {
	if ps == nil {
		return false
	}
	for _, card := range ps.Deck {
		if card != nil && (predicate == nil || predicate(card)) {
			return true
		}
	}
	return false
}

func countShadowCompanionsInGraveyard(ps *PlayerState) int {
	if ps == nil {
		return 0
	}
	count := 0
	for _, card := range ps.Graveyard {
		if card != nil && card.Card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementShadow {
			count++
		}
	}
	return count
}

func addGeneratedCardToPlayerHand(ctx *EffectContext, playerID int, cardNumber string) *CardInstance {
	card := getCardDB()[cardNumber]
	if card == nil {
		return nil
	}
	instance := ctx.Engine.newCardInstance(card, playerID, ctx.Engine.State.TurnNumber)
	ctx.Engine.addCardToHand(playerID, instance)
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(ctx.Source),
		"card":   cardToInfo(instance),
		"effect": "add_generated_card_to_hand",
	}})
	return instance
}

func addGeneratedCardsToPlayerDeck(ctx *EffectContext, playerID int, cardNumber string, count int) []*CardInstance {
	card := getCardDB()[cardNumber]
	if card == nil || count <= 0 {
		return nil
	}
	ps := ctx.Engine.State.Players[playerID]
	added := make([]*CardInstance, 0, count)
	for i := 0; i < count; i++ {
		instance := ctx.Engine.newCardInstance(card, playerID, ctx.Engine.State.TurnNumber)
		ps.Deck = append(ps.Deck, instance)
		added = append(added, instance)
	}
	ctx.Engine.shuffleDeck(playerID)
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(ctx.Source),
		"cards":  cardsToInfo(added),
		"effect": "add_generated_cards_to_deck",
	}})
	return added
}

func (e *Engine) discardRandomHandCard(playerID int) *CardInstance {
	ps := e.State.Players[playerID]
	if ps == nil || len(ps.Hand) == 0 {
		return nil
	}
	idx := e.randomIntn(len(ps.Hand))
	return e.discardHandCardAt(playerID, idx)
}

func (e *Engine) discardSelectedHandCards(playerID int, selected []string, limit int) int {
	return len(e.discardSelectedHandCardsMatching(playerID, selected, limit, nil))
}

func (e *Engine) discardSelectedHandCardsMatching(playerID int, selected []string, limit int, predicate func(*CardInstance) bool) []*CardInstance {
	if playerID < 0 || playerID >= len(e.State.Players) || limit <= 0 {
		return nil
	}
	ps := e.State.Players[playerID]
	selectedSet := map[string]bool{}
	for _, id := range selected {
		if id != "" {
			selectedSet[id] = true
		}
	}
	discarded := make([]*CardInstance, 0, limit)
	for i := len(ps.Hand) - 1; i >= 0 && len(discarded) < limit; i-- {
		card := ps.Hand[i]
		if card == nil || !selectedSet[card.InstanceID] {
			continue
		}
		if predicate != nil && !predicate(card) {
			continue
		}
		if discardedCard := e.discardHandCardAt(playerID, i); discardedCard != nil {
			discarded = append(discarded, discardedCard)
		}
	}
	return discarded
}

func (e *Engine) discardAllHandCards(playerID int) int {
	if playerID < 0 || playerID >= len(e.State.Players) {
		return 0
	}
	ps := e.State.Players[playerID]
	discarded := 0
	for len(ps.Hand) > 0 {
		if e.discardHandCardAt(playerID, len(ps.Hand)-1) != nil {
			discarded++
		}
	}
	return discarded
}

func (e *Engine) resolveDiscardedCardEffects(playerID int, card *CardInstance) {
	if card == nil || card.Card == nil {
		return
	}
	if card.Card.Number == "2001102" {
		if hero := e.playerHeroCard(playerID); hero != nil {
			e.dealDamage(hero, 2, playerID)
		}
	}
	if card.Card.Number == "2321103" {
		e.State.Players[playerID].GainElements(map[string]int{model.ElementAir: 1})
	}
	if card.Card.Number == "1321102" {
		e.offerDiscardedSpeckledSparrowSummon(playerID, card)
	}
}

func (e *Engine) removeCardFromGraveyard(playerID int, card *CardInstance) bool {
	if card == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return false
	}
	ps := e.State.Players[playerID]
	for i, candidate := range ps.Graveyard {
		if candidate == card {
			ps.Graveyard = append(ps.Graveyard[:i], ps.Graveyard[i+1:]...)
			return true
		}
	}
	return false
}

func resolveTopDeckReorder(e *Engine, playerID int, deckOwnerID int, looked []*CardInstance, data map[string]any, effect string) {
	if e == nil {
		return
	}
	ps := e.State.Players[deckOwnerID]
	lookCount := min(len(looked), len(ps.Deck))
	rest := append([]*CardInstance(nil), ps.Deck[lookCount:]...)
	pool := make(map[string]*CardInstance, len(looked))
	for _, card := range looked {
		if card != nil {
			pool[card.InstanceID] = card
		}
	}
	used := make(map[string]bool, len(pool))
	top := orderedWaterDivinationCards(stringsFromActionData(data, "top_order"), pool, used)
	bottom := orderedWaterDivinationCards(stringsFromActionData(data, "bottom_order"), pool, used)
	for _, card := range looked {
		if card == nil || used[card.InstanceID] {
			continue
		}
		top = append(top, card)
		used[card.InstanceID] = true
	}
	ps.Deck = append(append(top, rest...), bottom...)
	e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
		"effect":       effect,
		"deck_owner":   deckOwnerID,
		"top_order":    cardsToInfo(top),
		"bottom_order": cardsToInfo(bottom),
	}})
}

func (e *Engine) triggerRoseProphetAfterOpponentShuffle(shuffledPlayerID int) {
	if e == nil || shuffledPlayerID < 0 || shuffledPlayerID >= len(e.State.Players) {
		return
	}
	viewerID := 1 - shuffledPlayerID
	shuffled := e.State.Players[shuffledPlayerID]
	if len(shuffled.Deck) == 0 {
		return
	}
	for _, card := range e.getAllFieldCards(e.State.Players[viewerID]) {
		if card == nil || card.Card == nil || card.Card.Number != "1511103" || e.hasEffectiveStatus(card, StatusPetrify) {
			continue
		}
		lookCount := min(3, len(shuffled.Deck))
		looked := append([]*CardInstance(nil), shuffled.Deck[:lookCount]...)
		candidates := make([]map[string]any, 0, len(looked))
		for i, deckCard := range looked {
			info := candidateInfo(deckCard, "deck", "enemy")
			info["deck_index"] = i
			info["can_select"] = false
			candidates = append(candidates, info)
		}
		e.SetPendingActionWithData(viewerID, "rose_prophet_reorder",
			"玫瑰先知:调整对手牌库顶3张的顺序并放回牌库顶或牌库底", candidates, 0, 0,
			func(selected []string, data map[string]any) {
				resolveTopDeckReorder(e, viewerID, shuffledPlayerID, looked, data, "rose_prophet_reorder")
			})
	}
}
