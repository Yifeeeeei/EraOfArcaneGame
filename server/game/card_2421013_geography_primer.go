package game

import (
	"eraofarcane/model"
)

type Card2421013GeographyPrimer struct{ AlwaysActive }

func (Card2421013GeographyPrimer) ID() string { return "2421013" }

func (Card2421013GeographyPrimer) Name() string { return "《地理学入门》" }

func (Card2421013GeographyPrimer) ModifyCardPlayCost(ctx *EffectContext, card *CardInstance, cost map[string]int) {
	if card == nil || card.Card.TotalCost() <= 5 {
		return
	}
	reduceCost(cost, model.ElementEarth, 2)
}
