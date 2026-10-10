package game

import (
	"eraofarcane/model"
)

type Card4411003ProfessorMaggie struct{ AlwaysActive }

func (Card4411003ProfessorMaggie) ID() string { return "4411003" }

func (Card4411003ProfessorMaggie) Name() string { return "麦吉教授" }

func (Card4411003ProfessorMaggie) ModifyCardPlayCost(ctx *EffectContext, card *CardInstance, cost map[string]int) {
	if ctx.Source.Statuses["麦吉折扣"] == 0 && card != nil && card.Card.TotalCost() > 5 {
		reduceCost(cost, model.ElementEarth, 2)
	}
}

func (Card4411003ProfessorMaggie) OnCardPlayCostPaid(ctx *EffectContext, card *CardInstance) {
	if ctx.Source.Statuses["麦吉折扣"] == 0 && card != nil && card.Card.TotalCost() > 5 {
		ctx.Source.Statuses["麦吉折扣"] = 1
	}
}
