package game

import "eraofarcane/model"

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
	if !ok {
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
	candidates := ctx.Engine.friendlyDeckCards(ctx.PlayerID, predicate)
	if len(candidates) == 0 {
		return
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, actionType, prompt, candidates, 1, 1,
		func(selected []string) {
			if len(selected) > 0 {
				ctx.Engine.searchDeckToHand(ctx.PlayerID, selected[0])
			}
		})
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
