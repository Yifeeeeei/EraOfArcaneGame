package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card2311102LampusSword struct{ AlwaysActive }

func (Card2311102LampusSword) ID() string { return "2311102" }

func (Card2311102LampusSword) Name() string { return "兰普斯之剑" }

func (Card2311102LampusSword) PerTurnLabel(*CardInstance) string {
	return "主动"
}

func (Card2311102LampusSword) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	if !ctx.Engine.equipmentInOwnerSlot(ctx.PlayerID, ctx.Source) {
		return fmt.Errorf("兰普斯之剑需要在装备区")
	}
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.Category == model.ElementAir
	})
	sword := ctx.Source
	ctx.Engine.SetPendingAction(ctx.PlayerID, "lampus_sword_discard_air",
		"兰普斯之剑:弃置任意数量大气手牌", candidates, 0, len(candidates),
		func(selected []string) {
			if slot := equipmentSlotOf(ctx.Engine.State.Players[ctx.PlayerID], sword); slot >= 0 {
				ctx.Engine.moveEquipmentToGraveyard(ctx.PlayerID, slot, sword)
			}
			discarded := ctx.Engine.discardSelectedHandCardsMatching(ctx.PlayerID, selected, len(candidates), func(card *CardInstance) bool {
				return card != nil && card.Card != nil && card.Card.Category == model.ElementAir
			})
			damage := len(discarded)
			if damage > 0 {
				ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
					Type:             TempModLampusSwordDelayedDamage,
					SourceCardNumber: "2311102",
					SourceName:       "兰普斯之剑",
					Amount:           damage,
					RemainingUses:    1,
				})
			}
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source":    cardToInfo(sword),
				"effect":    "lampus_sword_prepare_damage",
				"discarded": damage,
			}})
		})
	return nil
}

var _ PerTurnAbility = Card2311102LampusSword{}
