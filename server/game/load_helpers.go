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
