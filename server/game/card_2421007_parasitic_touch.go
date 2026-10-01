package game

import (
	"eraofarcane/model"
)

type Card2421007ParasiticTouch struct{ AlwaysActive }

func (Card2421007ParasiticTouch) ID() string { return "2421007" }

func (Card2421007ParasiticTouch) Name() string { return "寄生之触" }

func (Card2421007ParasiticTouch) MasteryMax() int { return 1 }

func (Card2421007ParasiticTouch) OnMastery(ctx *EffectContext, level int) error {
	if level == 1 {
		ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementEarth, 1, ctx.Source)
	}
	return nil
}
