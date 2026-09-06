package game

import (
	"eraofarcane/model"
)

type Card1521115LoneStarIronKnight struct{ AlwaysActive }

func (Card1521115LoneStarIronKnight) ID() string { return "1521115" }

func (Card1521115LoneStarIronKnight) Name() string { return "孤星铁骑士" }

func (Card1521115LoneStarIronKnight) OnEnter(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ctx.Source == nil || ctx.Source.Position == nil || ps == nil || ctx.Source.Position.Row != ps.GetFrontRow() || len(adjacentFriendlyCompanions(ctx)) > 0 {
		return nil
	}
	ctx.Engine.gainLife(ctx.Source, 1, ctx.Source)
	ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementLight, 1, ctx.Source)
	return nil
}
