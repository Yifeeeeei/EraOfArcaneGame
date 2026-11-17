package game

import (
	"eraofarcane/model"
)

type Card1121108FireButterfly struct{ AlwaysActive }

func (Card1121108FireButterfly) ID() string { return "1121108" }

func (Card1121108FireButterfly) Name() string { return "火蝴蝶" }

func (Card1121108FireButterfly) PerTurnLabel(*CardInstance) string {
	return "主动"
}

func (Card1121108FireButterfly) OnPerTurn(ctx *EffectContext) error {
	clearFireButterflyStoredLoad(ctx.Source)
	if ctx.Source.ElementsGainSet != nil {
		ctx.Source.Statuses[fireButterflyPreviousLoadSetStatus] = 1
		for _, elem := range model.AllElements {
			if amount := ctx.Source.ElementsGainSet[elem]; amount != 0 {
				ctx.Source.Statuses[fireButterflyPreviousLoadValuePrefix+elem] = amount
			}
		}
	}
	ctx.Source.ElementsGainSet = map[string]int{model.ElementAir: 1}
	ctx.Source.Statuses[fireButterflyTemporaryLoadStatus] = 1
	return nil
}

func (Card1121108FireButterfly) OnTurnEnd(ctx *EffectContext) error {
	if ctx.Source.Statuses[fireButterflyTemporaryLoadStatus] <= 0 {
		return nil
	}
	if fireButterflyTemporaryLoadStillCurrent(ctx.Source) {
		if ctx.Source.Statuses[fireButterflyPreviousLoadSetStatus] > 0 {
			previous := make(map[string]int)
			for _, elem := range model.AllElements {
				if amount := ctx.Source.Statuses[fireButterflyPreviousLoadValuePrefix+elem]; amount != 0 {
					previous[elem] = amount
				}
			}
			setElementsGain(ctx.Source, previous)
		} else {
			clearElementsGainSet(ctx.Source)
		}
	}
	clearFireButterflyStoredLoad(ctx.Source)
	return nil
}

func fireButterflyTemporaryLoadStillCurrent(card *CardInstance) bool {
	if card == nil || card.ElementsGainSet == nil {
		return false
	}
	if card.ElementsGainSet[model.ElementAir] != 1 {
		return false
	}
	for _, elem := range model.AllElements {
		if elem == model.ElementAir {
			continue
		}
		if card.ElementsGainSet[elem] != 0 {
			return false
		}
	}
	return true
}

func clearFireButterflyStoredLoad(card *CardInstance) {
	if card == nil {
		return
	}
	delete(card.Statuses, fireButterflyTemporaryLoadStatus)
	delete(card.Statuses, fireButterflyPreviousLoadSetStatus)
	for _, elem := range model.AllElements {
		delete(card.Statuses, fireButterflyPreviousLoadValuePrefix+elem)
	}
}
