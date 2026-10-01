package game

import (
	"eraofarcane/model"
)

type Card1221108HeartLotusMirrorMage struct{ AlwaysActive }

func (Card1221108HeartLotusMirrorMage) ID() string { return "1221108" }

func (Card1221108HeartLotusMirrorMage) Name() string { return "心莲镜魔师" }

func (Card1221108HeartLotusMirrorMage) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.Card == nil {
		return nil
	}
	if ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) || !isFriendlySpellCast(ctx) || !hasCardTag(ctx.Target.Card, "创造") {
		return nil
	}
	if !deckHasMatch(ctx.Engine.State.Players[ctx.PlayerID], isWaterItem) {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "heart_lotus_mirror_mage_flip_water_item",
		"心莲镜魔师:是否翻取1张水纹道具", []map[string]any{candidateInfo(ctx.Source, "unit", "own")}, 0, 1,
		func(selected []string) {
			if len(selected) == 0 || ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) {
				return
			}
			drawn := ctx.Engine.flipDeckMatchesToHandThen(ctx.PlayerID, 1, 0, isWaterItem, func(drawn []*CardInstance) {
				if len(drawn) == 0 {
					return
				}
				card := drawn[0]
				if isCounterTrapCard(card.Card.Number) && ctx.Engine.freeEquipmentSlots(ctx.PlayerID) > 0 {
					makeEntryCostZero(card)
					ctx.Engine.setHandCounterTrapFree(ctx.PlayerID, card)
				}
			})
			if len(drawn) == 0 {
				return
			}
			ctx.Source.UsedThisTurn++
		})
	return nil
}

func isWaterItem(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.IsItem() && card.Card.Category == model.ElementWater
}

func (e *Engine) freeEquipmentSlots(playerID int) int {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return 0
	}
	ps := e.State.Players[playerID]
	free := 0
	for i := 0; i < equipmentSlotCapacity(ps); i++ {
		if ps.Equipment[i] == nil {
			free++
		}
	}
	return free
}

func (e *Engine) setHandCounterTrapFree(playerID int, card *CardInstance) bool {
	if e == nil || card == nil || card.Card == nil || playerID < 0 || playerID >= len(e.State.Players) || !isCounterTrapCard(card.Card.Number) {
		return false
	}
	ps := e.State.Players[playerID]
	_, handIdx := ps.FindHandCard(card.InstanceID)
	if handIdx < 0 || e.firstFreeEquipmentSlot(playerID) < 0 {
		return false
	}
	return e.placeCounterTrap(playerID, card, handIdx) == nil
}
