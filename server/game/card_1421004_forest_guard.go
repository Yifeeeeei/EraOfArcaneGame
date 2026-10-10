package game

import (
	"eraofarcane/model"
)

type Card1421004ForestGuard struct{ AlwaysActive }

func (Card1421004ForestGuard) ID() string { return "1421004" }

func (Card1421004ForestGuard) Name() string { return "森林守卫" }

func (Card1421004ForestGuard) MasteryMax() int { return 5 }

func (Card1421004ForestGuard) OnMastery(ctx *EffectContext, level int) error {
	switch level {
	case 1:
		ctx.Engine.gainLife(ctx.Source, 1, ctx.Source)
	case 3:
		ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementEarth, 1, ctx.Source)
	case 5:
		ctx.Source.AttackBonus += 2
	}
	return nil
}
