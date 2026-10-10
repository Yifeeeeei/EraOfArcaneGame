package game

import (
	"eraofarcane/model"
)

type Card1621101PainSoul struct{ AlwaysActive }

func (Card1621101PainSoul) ID() string { return "1621101" }

func (Card1621101PainSoul) Name() string { return "苦痛之魂" }

func (Card1621101PainSoul) DamageScope() DamageScope { return DamageSelf }

func (Card1621101PainSoul) OnDamaged(ctx *EffectContext, event DamageEvent) error {
	if !useTriggeredTurn(ctx.Source) {
		return nil
	}
	ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementShadow, 1, ctx.Source)
	return nil
}
