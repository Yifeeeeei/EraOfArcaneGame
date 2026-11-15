package game

import (
	"eraofarcane/model"
)

type Card1321115SkyPainter struct{ AlwaysActive }

func (Card1321115SkyPainter) ID() string { return "1321115" }

func (Card1321115SkyPainter) Name() string { return "苍穹描摹者" }

func (Card1321115SkyPainter) OnEnter(ctx *EffectContext) error {
	candidates := friendlyFieldCardsIncludingBound(ctx.Engine, ctx.PlayerID, func(card *CardInstance) bool {
		return canSkyPainterCopy(ctx.Engine, ctx.Source, card)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "sky_painter_copy_enter",
		"苍穹描摹者:复制另一个入场花费小于6的大气卡牌入场效果", candidates, 1, 1, nil, false,
		func(selected []string, _ map[string]any) error {
			target := ctx.Engine.findFriendlyCardIncludingBound(ctx.PlayerID, firstSelected(selected))
			if !canSkyPainterCopy(ctx.Engine, ctx.Source, target) {
				return nil
			}
			copied := *ctx
			copied.ExtraData = map[string]any{"copied_from": cardToInfo(target)}
			return behaviorForNumber(target.Card.Number).(OnEnterBehavior).OnEnter(&copied)
		})
	return nil
}

func canSkyPainterCopy(e *Engine, source *CardInstance, card *CardInstance) bool {
	if e == nil || card == nil || card == source || card.Card == nil ||
		card.Card.Category != model.ElementAir ||
		totalElementCost(card.Card.ElementsCost) >= 6 {
		return false
	}
	behavior, ok := behaviorForNumber(card.Card.Number).(OnEnterBehavior)
	return ok && behavior.HasActiveOnEnter(card)
}

func (r royalWaterUseCostReduction) HasActiveOnEnter(*CardInstance) bool {
	return r.triggerOnEnter
}
