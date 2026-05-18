package game

func effectiveElementsGain(card *CardInstance) map[string]int {
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
	return gains
}

func setElementsGain(card *CardInstance, gains map[string]int) {
	if card == nil {
		return
	}
	card.ElementsGainSet = make(map[string]int)
	for elem, amount := range gains {
		if amount != 0 {
			card.ElementsGainSet[elem] = amount
		}
	}
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
