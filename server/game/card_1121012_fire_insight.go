package game

import (
	"eraofarcane/model"
)

type Card1121012FireInsight struct{ AlwaysActive }

func (Card1121012FireInsight) ID() string { return "1121012" }

func (Card1121012FireInsight) Name() string { return "火焰洞察者" }

func (Card1121012FireInsight) DamageScope() DamageScope { return DamageAny }

func (Card1121012FireInsight) OnDamaged(ctx *EffectContext, event DamageEvent) error {
	if event.Element != model.ElementFire && event.Status != StatusBurn {
		return nil
	}
	if useTriggeredTurn(ctx.Source) {
		ctx.Engine.drawCards(ctx.PlayerID, 1)
	}
	return nil
}
