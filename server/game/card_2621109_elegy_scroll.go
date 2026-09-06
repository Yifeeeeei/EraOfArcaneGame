package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card2621109ElegyScroll struct{ AlwaysActive }

func (Card2621109ElegyScroll) ID() string { return "2621109" }

func (Card2621109ElegyScroll) Name() string { return "哀歌卷轴" }

func (Card2621109ElegyScroll) OnUseItem(ctx *EffectContext) error {
	hasShadowGrave := countShadowCompanionsInGraveyard(ctx.Engine.State.Players[ctx.PlayerID]) > 0
	drawn := ctx.Engine.flipDeckMatchesToHand(ctx.PlayerID, 1, 0, isShadowCompanionWithDeathrattle)
	if hasShadowGrave && len(drawn) > 0 {
		drawn[0].Statuses["入场费用"+model.ElementShadow+"-1"]++
	}
	return nil
}

func (Card2621109ElegyScroll) ValidateItemUse(ctx *EffectContext) error {
	e, playerID := ctx.Engine, ctx.PlayerID
	if len(e.friendlyDeckCards(playerID, isShadowCompanionWithDeathrattle)) == 0 {
		return fmt.Errorf("Elegy Scroll requires a searchable shadow companion with deathrattle")
	}
	return nil
}
