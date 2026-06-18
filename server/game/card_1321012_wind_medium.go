package game

import "eraofarcane/model"

type Card1321012WindMedium struct{ AlwaysActive }

func (Card1321012WindMedium) ID() string   { return "1321012" }
func (Card1321012WindMedium) Name() string { return "风灵媒师" }

func (Card1321012WindMedium) OnSpellCast(ctx *EffectContext) error {
	if !isFriendlySpellCast(ctx) {
		return nil
	}
	if ctx.Source == nil || ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) {
		return nil
	}
	skill, ok := ctx.ExtraData["skill"].(map[string]any)
	if !ok {
		if ctx.Target == nil || ctx.Target.Card == nil || ctx.Target.Card.Category != model.ElementAir {
			return nil
		}
		ctx.Source.UsedThisTurn++
		return DrawCards(1)(ctx)
	}
	if skill["category"] != model.ElementAir {
		return nil
	}
	ctx.Source.UsedThisTurn++
	return DrawCards(1)(ctx)
}
