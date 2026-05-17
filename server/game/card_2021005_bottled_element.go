package game

import "eraofarcane/model"

type Card2021005BottledElement struct{ AlwaysActive }

func (Card2021005BottledElement) ID() string   { return "2021005" }
func (Card2021005BottledElement) Name() string { return "瓶装元素" }

func (Card2021005BottledElement) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.State.Players[ctx.PlayerID].Elements[model.ElementArcane]++
	return nil
}
