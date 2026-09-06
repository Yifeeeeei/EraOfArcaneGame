package game

import (
	"eraofarcane/model"
)

type Card2421004DruidTest struct{ AlwaysActive }

func (Card2421004DruidTest) ID() string { return "2421004" }

func (Card2421004DruidTest) Name() string { return "德鲁伊水平测试" }

func (Card2421004DruidTest) OnUseItem(ctx *EffectContext) error {
	for _, unit := range ctx.Engine.getAllFieldCards(ctx.Engine.State.Players[ctx.PlayerID]) {
		if unit.Card.IsCompanion() && totalLoad(unit) > 2 {
			ctx.Engine.addElementsGainBonus(unit, ctx.PlayerID, model.ElementEarth, 1, ctx.Source)
		}
	}
	return nil
}
