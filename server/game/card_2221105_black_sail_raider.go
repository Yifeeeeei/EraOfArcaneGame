package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card2221105BlackSailRaider struct{ AlwaysActive }

func (Card2221105BlackSailRaider) ID() string { return "2221105" }

func (Card2221105BlackSailRaider) Name() string { return "掠夺者黑帆" }

func (Card2221105BlackSailRaider) OnUseItem(ctx *EffectContext) error {
	hasRaiderOnField := len(ctx.Engine.friendlyUnits(ctx.PlayerID, false, isRaiderCompanion)) > 0
	searchDeckToHandByPredicateWithResult(ctx,
		"black_sail_raider_search",
		"掠夺者黑帆:检索1个掠夺者伙伴",
		isRaiderCompanion,
		func(card *CardInstance) {
			if !hasRaiderOnField || card == nil {
				return
			}
			choices := []map[string]any{
				{"instance_id": model.ElementWater, "number": "2221105", "name": "入场花费-1水", "type": "选择", "zone": "choice", "side": "own"},
				{"instance_id": model.ElementShadow, "number": "2221105", "name": "入场花费-1暗", "type": "选择", "zone": "choice", "side": "own"},
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "black_sail_raider_discount",
				"掠夺者黑帆:选择检索卡牌的入场花费减免元素", choices, 1, 1,
				func(selected []string) {
					elem := firstSelected(selected)
					if elem != model.ElementWater && elem != model.ElementShadow {
						return
					}
					if !cardInstanceInSlice(ctx.Engine.State.Players[ctx.PlayerID].Hand, card) {
						return
					}
					card.Statuses["入场费用"+elem+"-1"]++
				})
		})
	return nil
}

func (Card2221105BlackSailRaider) ValidateItemUse(ctx *EffectContext) error {
	e, playerID := ctx.Engine, ctx.PlayerID
	if len(e.friendlyDeckCards(playerID, isRaiderCompanion)) == 0 {
		return fmt.Errorf("Black Sail Raider requires a searchable raider companion")
	}
	return nil
}

func cardInstanceInSlice(cards []*CardInstance, target *CardInstance) bool {
	for _, card := range cards {
		if card == target {
			return true
		}
	}
	return false
}
