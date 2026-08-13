package game

import "eraofarcane/model"

func effectiveElementsGain(card *CardInstance) map[string]int {
	return effectiveElementsGainWithStealth(card, card != nil && card.Statuses[StatusStealth] > 0 && card.Statuses[StatusPetrify] <= 0)
}

func (e *Engine) effectiveElementsGain(card *CardInstance) map[string]int {
	gains := effectiveElementsGainWithStealth(card, e != nil && e.hasActiveStealth(card))
	if e == nil || card == nil || card.Card == nil || card.OwnerID < 0 || card.OwnerID >= len(e.State.Players) {
		return gains
	}
	ps := e.State.Players[card.OwnerID]
	if ps == nil {
		return gains
	}
	ctx := &EffectContext{
		Engine:     e,
		Source:     card,
		PlayerID:   card.OwnerID,
		OpponentID: 1 - card.OwnerID,
	}
	for _, fieldCard := range e.getAllFieldCards(ps) {
		if fieldCard == nil || fieldCard.Card == nil || e.hasEffectiveStatus(fieldCard, StatusPetrify) {
			continue
		}
		behavior := globalRegistry.GetBehavior(fieldCard.Card.Number)
		modifier, ok := behavior.(ElementsGainModifier)
		if !ok || !modifier.HasActiveElementsGainModifier(fieldCard) {
			continue
		}
		ctx.Target = fieldCard
		modifier.ModifyElementsGain(ctx, card, gains)
	}
	return gains
}

func effectiveElementsGainWithStealth(card *CardInstance, hasEffectiveStealth bool) map[string]int {
	gains := make(map[string]int)
	if card == nil || card.Card == nil {
		return gains
	}
	base := card.Card.ElementsGain
	if card.ElementsGainSet != nil {
		base = card.ElementsGainSet
	}
	for elem, amount := range base {
		if amount != 0 {
			gains[elem] += amount
		}
	}
	for elem, amount := range card.ElementsGainBonus {
		if amount != 0 {
			gains[elem] += amount
		}
	}
	if card.Card.Number == "1221109" && hasEffectiveStealth {
		gains[model.ElementWater] += 2
	}
	return gains
}

func setElementsGain(card *CardInstance, gains map[string]int) {
	if card == nil {
		return
	}
	clearFireButterflyStoredLoad(card)
	card.ElementsGainSet = make(map[string]int)
	for elem, amount := range gains {
		if amount != 0 {
			card.ElementsGainSet[elem] = amount
		}
	}
}

func clearElementsGainSet(card *CardInstance) {
	if card == nil {
		return
	}
	clearFireButterflyStoredLoad(card)
	card.ElementsGainSet = nil
}

func addElementsGainBonus(card *CardInstance, elem string, amount int) {
	if card == nil || amount == 0 {
		return
	}
	if card.ElementsGainBonus == nil {
		card.ElementsGainBonus = make(map[string]int)
	}
	card.ElementsGainBonus[elem] += amount
}

func (e *Engine) addElementsGainBonus(card *CardInstance, ownerID int, elem string, amount int, source *CardInstance) {
	if card == nil || amount == 0 {
		return
	}
	addElementsGainBonus(card, elem, amount)
	e.emit(GameEvent{Type: "effect_trigger", Player: ownerID, Data: map[string]any{
		"source":  cardToInfo(source),
		"target":  cardToInfo(card),
		"effect":  "load_gain",
		"element": elem,
		"amount":  amount,
	}})
	e.triggerEffects(TriggerOnLoadGain, card, card, map[string]any{
		"load_gain_player": ownerID,
		"load_gain_source": source,
		"load_gain_target": card,
		"element":          elem,
		"amount":           amount,
	})
	e.triggerFieldEffectsWithData(TriggerOnLoadGain, ownerID, card, map[string]any{
		"load_gain_player": ownerID,
		"load_gain_source": source,
		"load_gain_target": card,
		"element":          elem,
		"amount":           amount,
	})
}

func (e *Engine) totalLoad(card *CardInstance) int {
	total := 0
	for _, amount := range e.effectiveElementsGain(card) {
		if amount > 0 {
			total += amount
		}
	}
	return total
}

func totalLoad(card *CardInstance) int {
	total := 0
	for _, amount := range effectiveElementsGain(card) {
		if amount > 0 {
			total += amount
		}
	}
	return total
}

func convertLoad(card *CardInstance, from string, to string, maxAmount int) int {
	if card == nil || maxAmount <= 0 {
		return 0
	}
	current := effectiveElementsGain(card)
	amount := min(current[from], maxAmount)
	if amount <= 0 {
		return 0
	}
	current[from] -= amount
	current[to] += amount
	setElementsGain(card, current)
	card.ElementsGainBonus = make(map[string]int)
	return amount
}
