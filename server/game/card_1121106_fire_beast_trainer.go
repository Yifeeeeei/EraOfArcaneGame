package game

import (
	"eraofarcane/model"
	"strings"
)

type Card1121106FireBeastTrainer struct{ AlwaysActive }

func (Card1121106FireBeastTrainer) ID() string { return "1121106" }

func (Card1121106FireBeastTrainer) Name() string { return "弗卡莱诺皇家驯兽师" }

func (Card1121106FireBeastTrainer) OnEnter(ctx *EffectContext) error {
	if ctx.Source == nil {
		return nil
	}
	ctx.Source.Statuses[fireBeastTrainerDiscountStatus] = 1
	return nil
}

func (Card1121106FireBeastTrainer) ModifyCardPlayCost(ctx *EffectContext, card *CardInstance, cost map[string]int) {
	if ctx.Source == nil || ctx.Source.Statuses[fireBeastTrainerDiscountStatus] <= 0 || !isFireBeastOrMonsterCompanion(card) {
		return
	}
	reduceGenericCost(cost, model.ElementFire, 2)
}

func (Card1121106FireBeastTrainer) OnCardPlayCostPaid(ctx *EffectContext, card *CardInstance) {
	if ctx.Source == nil || ctx.Source.Statuses[fireBeastTrainerDiscountStatus] <= 0 || !isFireBeastOrMonsterCompanion(card) {
		return
	}
	ctx.Source.Statuses[fireBeastTrainerDiscountStatus]--
	if ctx.Source.Statuses[fireBeastTrainerDiscountStatus] <= 0 {
		delete(ctx.Source.Statuses, fireBeastTrainerDiscountStatus)
	}
}

const fireBeastTrainerDiscountStatus = "弗卡莱诺皇家驯兽师下个火焰野兽异兽减费"

func isFireBeastOrMonsterCompanion(card *CardInstance) bool {
	if card == nil || card.Card == nil || !card.Card.IsCompanion() || card.Card.Category != model.ElementFire {
		return false
	}
	return strings.Contains(card.Card.Tag, "野兽") || strings.Contains(card.Card.Tag, "异兽")
}
