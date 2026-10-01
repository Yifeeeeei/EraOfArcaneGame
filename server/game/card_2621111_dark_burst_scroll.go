package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card2621111DarkBurstScroll struct{ AlwaysActive }

func (Card2621111DarkBurstScroll) ID() string { return "2621111" }

func (Card2621111DarkBurstScroll) Name() string { return "暗黑爆发卷轴" }

func (Card2621111DarkBurstScroll) OnUseItem(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	targets := make([]*CardInstance, 0)
	for _, card := range ps.Graveyard {
		if card != nil && card.Card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementShadow {
			targets = append(targets, card)
		}
	}
	if len(targets) < 5 {
		return nil
	}
	exiled := 0
	for _, card := range targets {
		if ctx.Engine.exileCard(ctx.PlayerID, card) {
			exiled++
		}
	}
	if exiled > 0 {
		ps.GainElements(map[string]int{model.ElementShadow: exiled * 2})
	}
	return nil
}

func (Card2621111DarkBurstScroll) ValidateItemUse(ctx *EffectContext) error {
	e, playerID := ctx.Engine, ctx.PlayerID
	if countShadowCompanionsInGraveyard(e.State.Players[playerID]) < 5 {
		return fmt.Errorf("Dark Burst Scroll requires at least five shadow companions in graveyard")
	}
	return nil
}
