package game

import "eraofarcane/model"

const DevourLife = "血"

func addSkillToPool(ctx *EffectContext, cardNumber string) {
	card := getCardDB()[cardNumber]
	if card == nil {
		return
	}
	skill := NewCardInstance(card, ctx.PlayerID, ctx.Engine.State.TurnNumber)
	ctx.Engine.State.Players[ctx.PlayerID].SkillPool = append(ctx.Engine.State.Players[ctx.PlayerID].SkillPool, skill)
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(ctx.Source),
		"effect": "bind_skill",
		"card":   cardToInfo(skill),
	}})
}

func bindSkillToHost(ctx *EffectContext, cardNumber string) {
	if ctx == nil || ctx.Source == nil {
		return
	}
	card := getCardDB()[cardNumber]
	if card == nil {
		return
	}
	skill := NewCardInstance(card, ctx.PlayerID, ctx.Engine.State.TurnNumber)
	skill.SlotIndex = -1
	skill.IsHorizontal = false
	ctx.Source.BoundSkills = append(ctx.Source.BoundSkills, skill)
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(ctx.Source),
		"effect": "bind_skill",
		"card":   cardToInfo(skill),
	}})
}

func devourFriendlyCompanion(ctx *EffectContext, cost map[string]int, prompt string) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if !ps.CanPayCost(cost) {
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != ctx.Source && card.Card.IsCompanion()
	})
	if len(candidates) == 0 {
		return nil
	}
	ps.PayCost(cost)
	ctx.Engine.SetPendingAction(ctx.PlayerID, "devour_companion", prompt, candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			target := ctx.Engine.findFieldCardByInstance(ps, selected[0])
			if target != nil && target != ctx.Source && target.Card.IsCompanion() {
				ctx.Engine.destroyUnit(target, ctx.PlayerID)
			}
		})
	return nil
}

func summonDevourRequirement(card *CardInstance) map[string]int {
	if card == nil || card.Card == nil {
		return nil
	}
	behavior := GetEffectRegistry().GetBehavior(card.Card.Number)
	devour, ok := behavior.(SummonDevourRequirementBehavior)
	if !ok || !devour.HasActiveDevourRequirement(card) {
		return nil
	}
	return devour.DevourRequirement()
}

func cardSatisfiesDevourRequirement(card *CardInstance, requirement map[string]int) bool {
	if len(requirement) == 0 {
		return true
	}
	gain := effectiveElementsGain(card)
	for elem, amount := range requirement {
		if gain[elem] < amount {
			return false
		}
	}
	return true
}

func searchDeckToHandByPredicate(ctx *EffectContext, actionType string, prompt string, predicate func(*CardInstance) bool) {
	searchDeckToHandByPredicateWithResult(ctx, actionType, prompt, predicate, nil)
}

func searchDeckToHandByPredicateWithResult(ctx *EffectContext, actionType string, prompt string, predicate func(*CardInstance) bool, afterSearch func(*CardInstance)) {
	candidates := ctx.Engine.friendlyDeckCards(ctx.PlayerID, predicate)
	if len(candidates) == 0 {
		return
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, actionType, prompt, candidates, 1, 1,
		func(selected []string) {
			if len(selected) > 0 {
				card := ctx.Engine.searchDeckCardToHand(ctx.PlayerID, selected[0])
				if card != nil && afterSearch != nil {
					afterSearch(card)
				}
			}
		})
}

func summonCardFreeFromHandOrDeck(ctx *EffectContext, instanceID string) *CardInstance {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	pos := ps.FindEmptyPosition()
	if pos == nil || instanceID == "" {
		return nil
	}
	var card *CardInstance
	for i, candidate := range ps.Hand {
		if candidate != nil && candidate.InstanceID == instanceID {
			card = candidate
			ps.Hand = append(ps.Hand[:i], ps.Hand[i+1:]...)
			break
		}
	}
	if card == nil {
		for i, candidate := range ps.Deck {
			if candidate != nil && candidate.InstanceID == instanceID {
				card = candidate
				ps.Deck = append(ps.Deck[:i], ps.Deck[i+1:]...)
				ctx.Engine.shuffleDeck(ctx.PlayerID)
				break
			}
		}
	}
	if card == nil {
		return nil
	}
	card.OwnerID = ctx.PlayerID
	card.Position = pos
	card.IsHorizontal = true
	card.EnterTurn = ctx.Engine.State.TurnNumber
	ps.Units[pos.Col][pos.Row] = card
	ctx.Engine.triggerEffects(TriggerOnEnter, card, nil, nil)
	ctx.Engine.triggerFieldEffectsWithData(TriggerOnUnitEnter, ctx.PlayerID, card, map[string]any{"entered_player": ctx.PlayerID})
	ctx.Engine.triggerFieldEffectsWithData(TriggerOnUnitEnter, ctx.OpponentID, card, map[string]any{"entered_player": ctx.PlayerID})
	return card
}

func isWaterCompanion(card *CardInstance) bool {
	return card.Card.IsCompanion() && card.Card.Category == model.ElementWater
}

func isWaterCard(card *CardInstance) bool {
	return card.Card.Category == model.ElementWater
}

func totalCost(cost map[string]int) int {
	total := 0
	for _, amount := range cost {
		total += amount
	}
	return total
}
