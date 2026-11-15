package game

import (
	"eraofarcane/model"
	"strings"
)

type Card1121106FireBeastTrainer struct{ AlwaysActive }

func (Card1121106FireBeastTrainer) ID() string { return "1121106" }

func (Card1121106FireBeastTrainer) Name() string { return "弗卡莱诺皇家驯兽师" }

func (Card1121106FireBeastTrainer) OnEnter(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{Type: TempModNextFireBeastCostMinus, SourceCardNumber: "1121106", SourceName: "弗卡莱诺皇家驯兽师", Element: model.ElementFire, Amount: 2, RemainingUses: 1})
	return nil
}

func isFireBeastOrMonsterCompanion(card *CardInstance) bool {
	if card == nil || card.Card == nil || !card.Card.IsCompanion() || card.Card.Category != model.ElementFire {
		return false
	}
	return strings.Contains(card.Card.Tag, "野兽") || strings.Contains(card.Card.Tag, "异兽")
}
