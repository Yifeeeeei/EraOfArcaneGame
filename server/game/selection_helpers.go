package game

import (
	"math/rand"
	"strings"

	"eraofarcane/cards"
	"eraofarcane/model"
)

func firstSelected(selected []string) string {
	if len(selected) == 0 {
		return ""
	}
	return selected[0]
}

func candidateInfo(card *CardInstance, zone string, side string) map[string]any {
	info := cardToInfo(card)
	info["zone"] = zone
	info["side"] = side
	return info
}

func (e *Engine) friendlyUnits(playerID int, includeHero bool, predicate func(*CardInstance) bool) []map[string]any {
	ps := e.State.Players[playerID]
	candidates := make([]map[string]any, 0)
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := ps.Units[col][row]
			if unit == nil {
				continue
			}
			if !includeHero && unit.Card.IsHero() {
				continue
			}
			if predicate != nil && !predicate(unit) {
				continue
			}
			candidates = append(candidates, candidateInfo(unit, "unit", "own"))
		}
	}
	return candidates
}

func (e *Engine) enemyUnits(playerID int, includeHero bool, predicate func(*CardInstance) bool) []map[string]any {
	opponentID := 1 - playerID
	ps := e.State.Players[opponentID]
	candidates := make([]map[string]any, 0)
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := ps.Units[col][row]
			if unit == nil {
				continue
			}
			if !includeHero && unit.Card.IsHero() {
				continue
			}
			if predicate != nil && !predicate(unit) {
				continue
			}
			candidates = append(candidates, candidateInfo(unit, "unit", "enemy"))
		}
	}
	return candidates
}

func (e *Engine) friendlyHandCards(playerID int, predicate func(*CardInstance) bool) []map[string]any {
	ps := e.State.Players[playerID]
	candidates := make([]map[string]any, 0, len(ps.Hand))
	for _, card := range ps.Hand {
		if card == nil {
			continue
		}
		if predicate != nil && !predicate(card) {
			continue
		}
		candidates = append(candidates, candidateInfo(card, "hand", "own"))
	}
	return candidates
}

func (e *Engine) friendlyEquipment(playerID int, predicate func(*CardInstance) bool) []map[string]any {
	ps := e.State.Players[playerID]
	candidates := make([]map[string]any, 0)
	for _, card := range ps.Equipment {
		if card == nil {
			continue
		}
		if predicate != nil && !predicate(card) {
			continue
		}
		candidates = append(candidates, candidateInfo(card, "equipment", "own"))
	}
	return candidates
}

func (e *Engine) enemyEquipment(playerID int, predicate func(*CardInstance) bool) []map[string]any {
	ps := e.State.Players[1-playerID]
	candidates := make([]map[string]any, 0)
	for _, card := range ps.Equipment {
		if card == nil {
			continue
		}
		if predicate != nil && !predicate(card) {
			continue
		}
		candidates = append(candidates, candidateInfo(card, "equipment", "enemy"))
	}
	return candidates
}

func (e *Engine) friendlySkills(playerID int, predicate func(*CardInstance) bool) []map[string]any {
	ps := e.State.Players[playerID]
	candidates := make([]map[string]any, 0)
	for _, skill := range ps.Skills {
		if skill == nil {
			continue
		}
		if predicate != nil && !predicate(skill) {
			continue
		}
		candidates = append(candidates, candidateInfo(skill, "skill", "own"))
	}
	return candidates
}

func (e *Engine) friendlySkillsIncludingBound(playerID int, predicate func(*CardInstance) bool) []map[string]any {
	ps := e.State.Players[playerID]
	candidates := e.friendlySkills(playerID, predicate)
	for _, card := range e.getAllFieldCards(ps) {
		if card == nil {
			continue
		}
		for _, skill := range card.BoundSkills {
			if skill == nil {
				continue
			}
			if predicate != nil && !predicate(skill) {
				continue
			}
			candidates = append(candidates, candidateInfo(skill, "bound_skill", "own"))
		}
	}
	return candidates
}

func (e *Engine) enemySkills(playerID int, predicate func(*CardInstance) bool) []map[string]any {
	ps := e.State.Players[1-playerID]
	candidates := make([]map[string]any, 0)
	for _, skill := range ps.Skills {
		if skill == nil {
			continue
		}
		if predicate != nil && !predicate(skill) {
			continue
		}
		candidates = append(candidates, candidateInfo(skill, "skill", "enemy"))
	}
	return candidates
}

func (e *Engine) friendlyDeckCards(playerID int, predicate func(*CardInstance) bool) []map[string]any {
	ps := e.State.Players[playerID]
	candidates := make([]map[string]any, 0)
	for _, card := range ps.Deck {
		if card == nil {
			continue
		}
		if !canFlipOrSearchCard(card) {
			continue
		}
		if predicate != nil && !predicate(card) {
			continue
		}
		candidates = append(candidates, candidateInfo(card, "deck", "own"))
	}
	return candidates
}

func (e *Engine) findFriendlyCandidate(playerID int, instanceID string) (*CardInstance, string) {
	ps := e.State.Players[playerID]
	for _, card := range ps.Hand {
		if card != nil && card.InstanceID == instanceID {
			return card, "hand"
		}
	}
	for _, card := range ps.Equipment {
		if card != nil && card.InstanceID == instanceID {
			return card, "equipment"
		}
	}
	for _, card := range ps.Skills {
		if card != nil && card.InstanceID == instanceID {
			return card, "skill"
		}
	}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := ps.Units[col][row]
			if unit != nil && unit.InstanceID == instanceID {
				return unit, "unit"
			}
		}
	}
	return nil, ""
}

func (e *Engine) findUnitByInstanceID(instanceID string) *CardInstance {
	for _, ps := range e.State.Players {
		if unit := e.findUnitOnGrid(ps, instanceID); unit != nil {
			return unit
		}
	}
	return nil
}

func (e *Engine) discardFriendlyCandidate(playerID int, instanceID string) bool {
	ps := e.State.Players[playerID]
	for i, card := range ps.Hand {
		if card != nil && card.InstanceID == instanceID {
			ps.Graveyard = append(ps.Graveyard, card)
			ps.Hand = append(ps.Hand[:i], ps.Hand[i+1:]...)
			delete(ps.RevealedHand, card.InstanceID)
			e.emit(GameEvent{Type: "discard", Player: playerID, Data: map[string]any{"card": cardToInfo(card)}})
			e.resolveDiscardedCardEffects(playerID, card)
			return true
		}
	}
	for i, card := range ps.Equipment {
		if card != nil && card.InstanceID == instanceID {
			e.moveEquipmentToGraveyard(playerID, i, card)
			return true
		}
	}
	return false
}

func (e *Engine) destroyEnemyEquipment(playerID int, instanceID string) bool {
	ps := e.State.Players[1-playerID]
	for i, card := range ps.Equipment {
		if card == nil || card.InstanceID != instanceID {
			continue
		}
		e.moveEquipmentToGraveyard(1-playerID, i, card)
		return true
	}
	return false
}

func (e *Engine) searchDeckToHand(playerID int, instanceID string) bool {
	return e.searchDeckCardToHand(playerID, instanceID) != nil
}

func (e *Engine) searchDeckCardToHand(playerID int, instanceID string) *CardInstance {
	return e.searchDeckCardToHandThen(playerID, instanceID, nil)
}

func (e *Engine) searchDeckCardToHandThen(playerID int, instanceID string, afterSearch func(*CardInstance)) *CardInstance {
	ps := e.State.Players[playerID]
	for i, card := range ps.Deck {
		if card != nil && card.InstanceID == instanceID && canFlipOrSearchCard(card) {
			ps.Hand = append(ps.Hand, card)
			ps.Deck = append(ps.Deck[:i], ps.Deck[i+1:]...)
			e.shuffleDeck(playerID)
			e.emit(GameEvent{Type: "search_card", Player: playerID, Data: map[string]any{"card": cardToInfo(card)}})
			e.notifyCardSearchedThen(playerID, card, func() {
				if afterSearch != nil {
					afterSearch(card)
				}
			})
			return card
		}
	}
	return nil
}

func (e *Engine) flipDeckMatchesToHand(playerID int, count int, limit int, predicate func(*CardInstance) bool) []*CardInstance {
	if count <= 0 {
		return nil
	}
	ps := e.State.Players[playerID]
	lookLimit := len(ps.Deck)
	if limit > 0 && limit < lookLimit {
		lookLimit = limit
	}
	drawn := make([]*CardInstance, 0, count)
	remaining := make([]*CardInstance, 0, len(ps.Deck))
	for i, card := range ps.Deck {
		withinLook := i < lookLimit
		if withinLook && len(drawn) < count && canFlipOrSearchCard(card) && (predicate == nil || predicate(card)) {
			drawn = append(drawn, card)
			continue
		}
		remaining = append(remaining, card)
	}
	ps.Deck = remaining
	for _, card := range drawn {
		ps.Hand = append(ps.Hand, card)
		e.notifyCardDrawn(playerID, card)
	}
	e.shuffleDeck(playerID)
	e.emit(GameEvent{
		Type:   "flip_deck",
		Player: playerID,
		Data: map[string]any{
			"cards": cardsToInfo(drawn),
			"count": len(drawn),
			"limit": limit,
		},
	})
	return drawn
}

func canFlipOrSearchCard(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.Number != "2211101"
}

func (e *Engine) drawFirstDeckMatch(playerID int, predicate func(*CardInstance) bool) *CardInstance {
	ps := e.State.Players[playerID]
	for i, card := range ps.Deck {
		if card == nil || !canFlipOrSearchCard(card) || (predicate != nil && !predicate(card)) {
			continue
		}
		ps.Hand = append(ps.Hand, card)
		ps.Deck = append(ps.Deck[:i], ps.Deck[i+1:]...)
		e.emit(GameEvent{Type: "search_card", Player: playerID, Data: map[string]any{"card": cardToInfo(card)}})
		e.notifyCardSearched(playerID, card)
		return card
	}
	return nil
}

func (e *Engine) moveGraveyardCardToDeckTop(playerID int, instanceID string) bool {
	ps := e.State.Players[playerID]
	for i, card := range ps.Graveyard {
		if card == nil || card.InstanceID != instanceID {
			continue
		}
		ps.Graveyard = append(ps.Graveyard[:i], ps.Graveyard[i+1:]...)
		resetCardForHiddenZone(card)
		ps.Deck = append([]*CardInstance{card}, ps.Deck...)
		e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
			"effect": "graveyard_to_deck_top",
			"card":   cardToInfo(card),
		}})
		return true
	}
	return false
}

func (e *Engine) moveGraveyardCardToHand(playerID int, instanceID string) bool {
	ps := e.State.Players[playerID]
	for i, card := range ps.Graveyard {
		if card == nil || card.InstanceID != instanceID {
			continue
		}
		ps.Graveyard = append(ps.Graveyard[:i], ps.Graveyard[i+1:]...)
		resetCardForHiddenZone(card)
		ps.Hand = append(ps.Hand, card)
		e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
			"effect": "graveyard_to_hand",
			"card":   cardToInfo(card),
		}})
		return true
	}
	return false
}

func resetCardForHiddenZone(card *CardInstance) {
	if card == nil || card.Card == nil {
		return
	}
	resetCardForPublicSpecialZone(card)
	card.IsHorizontal = true
}

func resetCardForPublicSpecialZone(card *CardInstance) {
	if card == nil || card.Card == nil {
		return
	}
	card.CurrentLife = card.Card.Life
	card.CurrentAttack = card.Card.Attack
	card.DamageTakenThisTurn = 0
	card.IsHorizontal = false
	card.Statuses = make(map[string]int)
	card.ElementsGainBonus = make(map[string]int)
	card.ElementsGainSet = nil
	card.PowerBonus = 0
	card.AttackBonus = 0
	card.IsSetCounter = false
	card.Position = nil
	card.SlotIndex = -1
	card.EnterTurn = 0
	card.BoundSkills = nil
	card.UnderCards = nil
	card.AttachedBehaviors = nil
	card.UsedThisTurn = 0
	card.UltimateUsed = false
	card.UsesRemaining = 0
}

func (e *Engine) shuffleDeck(playerID int) {
	deck := e.State.Players[playerID].Deck
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
}

func (e *Engine) hasAnyEquipment(playerID int) bool {
	for _, card := range e.State.Players[playerID].Equipment {
		if card != nil {
			return true
		}
	}
	return false
}

func resetCard(card *CardInstance) {
	if card != nil {
		card.IsHorizontal = false
	}
}

func isNonHeroFireCard(card *CardInstance) bool {
	return card != nil && !card.Card.IsHero() && card.Card.Category == model.ElementFire
}

func isLightSkill(card *CardInstance) bool {
	return card != nil && card.Card.IsSkill() && card.Card.Category == model.ElementLight
}

func isFireCompanionWithCostAboveFour(card *CardInstance) bool {
	return card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementFire && totalElementCost(card.Card.ElementsCost) >= 4
}

func isFireCompanion(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementFire
}

func isEarthCompanionWithCostAboveFive(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementEarth && totalElementCost(card.Card.ElementsCost) > 5
}

func isLightSpirit(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementLight && hasCardTag(card.Card, "精灵")
}

func isShadowCompanionWithDeathrattle(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementShadow && cardHasActiveDeathrattle(card)
}

func isAirEquipment(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.IsItem() && card.Card.Category == model.ElementAir && cards.IsEquipment(card.Card.Number)
}

func hasCardTag(card *model.Card, tag string) bool {
	return card != nil && strings.Contains(card.Tag, tag)
}

func lowCostSkill(card *CardInstance) bool {
	if card == nil || !card.Card.IsSkill() {
		return false
	}
	return totalElementCost(skillUseCost(card.Card)) < 3
}

func totalElementCost(cost map[string]int) int {
	total := 0
	for _, amount := range cost {
		if amount > 0 {
			total += amount
		}
	}
	return total
}
