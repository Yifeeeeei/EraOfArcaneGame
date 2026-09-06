package game

import (
	"strings"
)

func findEnemySkillIncludingBound(e *Engine, playerID int, instanceID string) *CardInstance {
	for _, skill := range enemySpellInstancesIncludingBound(e, playerID) {
		if skill != nil && skill.InstanceID == instanceID {
			return skill
		}
	}
	return nil
}

func friendlyFieldCardsIncludingBound(e *Engine, playerID int, predicate func(*CardInstance) bool) []map[string]any {
	ps := e.State.Players[playerID]
	candidates := make([]map[string]any, 0)
	for _, card := range e.getAllFieldCards(ps) {
		if card == nil || predicate != nil && !predicate(card) {
			continue
		}
		candidates = append(candidates, candidateInfo(card, "field", "own"))
		for _, skill := range card.BoundSkills {
			if skill == nil || predicate != nil && !predicate(skill) {
				continue
			}
			candidates = append(candidates, candidateInfo(skill, "bound_skill", "own"))
		}
	}
	return candidates
}

func isRaiderCompanion(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.IsCompanion() && strings.Contains(card.Card.Name, "掠夺者")
}

func (e *Engine) destroyAndisGiftDoomedUnits(ps *PlayerState) {
	if e == nil || ps == nil {
		return
	}
	for _, card := range append([]*CardInstance(nil), e.getAllFieldCards(ps)...) {
		if card == nil || card.Card == nil || card.Statuses[andisGiftDoomedStatus] <= 0 {
			continue
		}
		delete(card.Statuses, andisGiftDoomedStatus)
		if card.Card.IsHero() {
			card.CurrentLife = 0
			e.emit(GameEvent{Type: "unit_destroyed", Player: -1, Data: map[string]any{
				"player": ps.PlayerID,
				"card":   cardToInfo(card),
				"reason": "andis_gift",
			}})
			e.checkWinCondition()
			continue
		}
		e.destroyUnitWithCause(card, ps.PlayerID, "andis_gift")
	}
}

func royalFriendlyUnits(ctx *EffectContext) []*CardInstance {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	units := make([]*CardInstance, 0, 9)
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if unit := ps.Units[col][row]; unit != nil {
				units = append(units, unit)
			}
		}
	}
	return units
}

func (e *Engine) findFriendlyCardIncludingBound(playerID int, instanceID string) *CardInstance {
	if card, _ := e.findFriendlyCandidate(playerID, instanceID); card != nil {
		return card
	}
	ps := e.State.Players[playerID]
	if ps == nil {
		return nil
	}
	for _, host := range e.getAllFieldCards(ps) {
		if host == nil {
			continue
		}
		for _, skill := range host.BoundSkills {
			if skill != nil && skill.InstanceID == instanceID {
				return skill
			}
		}
	}
	return nil
}

func (e *Engine) placeExistingCompanionAtPosition(playerID int, card *CardInstance, pos Position, triggerEnter bool) bool {
	if card == nil || card.Card == nil || !card.Card.IsCompanion() || !pos.Valid() || playerID < 0 || playerID >= len(e.State.Players) {
		return false
	}
	ps := e.State.Players[playerID]
	if ps.Units[pos.Col][pos.Row] != nil {
		return false
	}
	card.OwnerID = playerID
	card.Position = &Position{Col: pos.Col, Row: pos.Row}
	card.EnterTurn = e.State.TurnNumber
	ps.Units[pos.Col][pos.Row] = card
	e.ApplySummonModifiersOnEnter(card)
	if triggerEnter {
		e.triggerEffects(TriggerOnEnter, card, nil, nil)
		e.notifyCardEntered(playerID, card, map[string]any{"entered_player": playerID})
		e.triggerFieldEffectsWithData(TriggerOnUnitEnter, playerID, card, map[string]any{"entered_player": playerID})
		e.triggerFieldEffectsWithData(TriggerOnUnitEnter, 1-playerID, card, map[string]any{"entered_player": playerID})
	}
	return true
}

func (e *Engine) summonFreshCardAtPosition(playerID int, cardNumber string, pos Position, triggerEnter bool) *CardInstance {
	cardDef := getCardDB()[cardNumber]
	if cardDef == nil {
		return nil
	}
	instance := e.newCardInstance(cardDef, playerID, e.State.TurnNumber)
	if !e.placeExistingCompanionAtPosition(playerID, instance, pos, triggerEnter) {
		return nil
	}
	return instance
}

func resetCardForResummon(card *CardInstance) {
	if card == nil || card.Card == nil {
		return
	}
	card.CurrentLife = card.Card.Life
	card.CurrentAttack = card.Card.Attack
	card.DamageTakenThisTurn = 0
	card.IsHorizontal = true
	card.Position = nil
	card.Statuses = make(map[string]int)
	card.ElementsGainBonus = make(map[string]int)
	card.ElementsGainSet = nil
	card.PowerBonus = 0
	card.AttackBonus = 0
	card.UsedThisTurn = 0
	card.UltimateUsed = false
	card.BoundSkills = nil
	card.UnderCards = nil
	card.AttachedBehaviors = nil
}

func adjacentFriendlyCompanions(ctx *EffectContext) []map[string]any {
	if ctx.Source == nil || ctx.Source.Position == nil {
		return nil
	}
	candidates := make([]map[string]any, 0, 4)
	for _, unit := range adjacentUnits(ctx.Engine.State.Players[ctx.PlayerID], ctx.Source.Position) {
		if unit != nil && unit.Card != nil && unit.Card.IsCompanion() {
			candidates = append(candidates, candidateInfo(unit, "unit", "own"))
		}
	}
	return candidates
}
