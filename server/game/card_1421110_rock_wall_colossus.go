package game

import (
	"eraofarcane/model"
)

type Card1421110RockWallColossus struct{ AlwaysActive }

func (Card1421110RockWallColossus) ID() string { return "1421110" }

func (Card1421110RockWallColossus) Name() string { return "岩壁巨像" }

func (Card1421110RockWallColossus) OnUnitEnter(ctx *EffectContext) error {
	if ctx == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.Card == nil {
		return nil
	}
	if ctx.Target.OwnerID != ctx.PlayerID || ctx.Target.Card.Category != model.ElementEarth || !ctx.Target.Card.IsCompanion() {
		return nil
	}
	if playerHasLearnedSpell(ctx.Engine.State.Players[ctx.PlayerID]) {
		return nil
	}
	ctx.Target.Statuses["max_life_bonus"]++
	ctx.Engine.gainLife(ctx.Target, 1, ctx.Source)
	return nil
}

var _ OnUnitEnterBehavior = Card1421110RockWallColossus{}
