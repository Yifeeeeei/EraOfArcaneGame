package game

import (
	"eraofarcane/cards"
	"fmt"
)

type Card1011102ConductorLos struct{ AlwaysActive }

func (Card1011102ConductorLos) ID() string { return "1011102" }

func (Card1011102ConductorLos) Name() string { return "\"指挥家\" 洛斯" }

func (Card1011102ConductorLos) OnConsume(ctx *EffectContext) error {
	if ctx.Source == nil || ctx.Target == nil {
		return nil
	}
	ctx.Source.Statuses[conductorConsumedCountStatus]++
	if ctx.Source.Statuses[conductorConsumedCountStatus] < 4 {
		return nil
	}
	if ctx.Engine.firstFreeEquipmentSlot(ctx.PlayerID) < 0 {
		return nil
	}
	ctx.Source.Statuses[conductorConsumedCountStatus] -= 4
	if ctx.Source.Statuses[conductorConsumedCountStatus] <= 0 {
		delete(ctx.Source.Statuses, conductorConsumedCountStatus)
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "conductor_los_equip_finale_violin",
		"\"指挥家\" 洛斯:是否装备1个落幕提琴", []map[string]any{candidateInfo(ctx.Source, "unit", "own")}, 0, 1,
		func(selected []string) {
			if len(selected) == 0 || !ctx.Engine.cardStillOnField(ctx.Source) || ctx.Engine.firstFreeEquipmentSlot(ctx.PlayerID) < 0 {
				return
			}
			ctx.Engine.equipGeneratedCard(ctx.PlayerID, "2001101")
		})
	return nil
}

func (Card1011102ConductorLos) PerTurnLabel(*CardInstance) string {
	return "消耗"
}

func (Card1011102ConductorLos) OnPerTurn(ctx *EffectContext) error {
	if ctx.Source == nil || !ctx.Engine.canConsumeCard(ctx.Source) || !ctx.Engine.cardStillOnField(ctx.Source) {
		return fmt.Errorf("\"指挥家\" 洛斯需要在场且竖置才能消耗")
	}
	ctx.Engine.consumeCardForEffectWithTriggers(ctx.PlayerID, ctx.Source, ctx.Engine.effectiveElementsGain(ctx.Source), "1011102")
	for _, equipment := range ctx.Engine.State.Players[ctx.PlayerID].Equipment {
		if equipment != nil && equipment.Card != nil && equipment.Card.Number == "2001101" {
			equipment.IsHorizontal = false
		}
	}
	return nil
}

const conductorConsumedCountStatus = "指挥家消耗计数"

func (e *Engine) equipGeneratedCard(playerID int, number string) *CardInstance {
	if e == nil || number == "" || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	cardDef, ok := cards.PlayableCardDB[number]
	if !ok || !cardDef.IsItem() || !isEquipmentCard(cardDef) {
		return nil
	}
	slot := e.firstFreeEquipmentSlot(playerID)
	if slot < 0 {
		return nil
	}
	card := e.newCardInstance(cardDef, playerID, e.State.TurnNumber)
	card.IsHorizontal = true
	card.SlotIndex = slot
	ps := e.State.Players[playerID]
	ps.Equipment[slot] = card
	e.emit(GameEvent{
		Type:   "equip",
		Player: -1,
		Data: map[string]any{
			"player": playerID,
			"card":   cardToInfo(card),
			"slot":   slot,
		},
	})
	e.triggerEffects(TriggerOnEquip, card, nil, nil)
	e.triggerEffects(TriggerOnEnter, card, nil, nil)
	e.notifyCardEntered(playerID, card, map[string]any{"entered_player": playerID, "equipped": true})
	return card
}
