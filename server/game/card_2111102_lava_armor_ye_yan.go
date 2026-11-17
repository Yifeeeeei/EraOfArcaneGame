package game

import "fmt"

type Card2111102LavaArmorYeYan struct{ AlwaysActive }

func (Card2111102LavaArmorYeYan) ID() string { return "2111102" }

func (Card2111102LavaArmorYeYan) Name() string { return "熔岩魔甲 业炎" }

func (Card2111102LavaArmorYeYan) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil {
		return nil
	}
	attacker, _ := ctx.ExtraData["attacker"].(int)
	if attacker == ctx.PlayerID || !ctx.Engine.cardStillOnField(ctx.Source) {
		return nil
	}
	armor := ctx.Source
	ctx.Engine.SetPendingAction(ctx.PlayerID, "lava_armor_yeyan_sacrifice",
		"熔岩魔甲 业炎:是否献祭此卡获得护盾2", []map[string]any{candidateInfo(armor, "equipment", "own")}, 0, 1,
		func(selected []string) {
			if len(selected) == 0 || !ctx.Engine.cardStillOnField(armor) {
				return
			}
			slot := equipmentSlotOf(ctx.Engine.State.Players[ctx.PlayerID], armor)
			if slot < 0 {
				return
			}
			ctx.Engine.moveEquipmentToGraveyard(ctx.PlayerID, slot, armor)
			ctx.Engine.gainPlayerShield(ctx.PlayerID, 2)
		})
	return nil
}

func (Card2111102LavaArmorYeYan) OnDeath(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.PlayerID < 0 || ctx.PlayerID >= len(ctx.Engine.State.Players) {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ps == nil {
		return nil
	}
	if ps.ShieldBrokenThisTurn {
		ctx.Engine.equipMoltenArmorForLavaArmorYeYan(ctx.PlayerID, ctx.Source)
		return nil
	}
	sourceNumber, sourceName, sourceID := "", "", ""
	if ctx.Source != nil && ctx.Source.Card != nil {
		sourceNumber = ctx.Source.Card.Number
		sourceName = ctx.Source.Card.Name
	}
	if ctx.Source != nil {
		sourceID = ctx.Source.InstanceID
	}
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModLavaArmorYeYanShieldBreak,
		SourceCardNumber: sourceNumber,
		SourceName:       sourceName,
		TargetInstanceID: sourceID,
		RemainingUses:    1,
		ExpiresTurn:      ctx.Engine.State.TurnNumber + 1,
	})
	return nil
}

var _ OnSpellHitBehavior = Card2111102LavaArmorYeYan{}

var _ OnDeathBehavior = Card2111102LavaArmorYeYan{}

func (e *Engine) equipMoltenArmorForLavaArmorYeYan(playerID int, source *CardInstance) *CardInstance {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	equipped := e.equipCardFromHandOrDeckFree(playerID, "2121013")
	if equipped != nil {
		e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
			"source": cardToInfo(source),
			"effect": "lava_armor_yeyan_equip_molten_armor",
			"card":   cardToInfo(equipped),
		}})
	}
	return equipped
}

func (e *Engine) equipCardFromHandOrDeckFree(playerID int, number string) *CardInstance {
	if e == nil || number == "" || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	ps := e.State.Players[playerID]
	var card *CardInstance
	for _, candidate := range append(append([]*CardInstance(nil), ps.Hand...), ps.Deck...) {
		if candidate != nil && candidate.Card != nil && candidate.Card.Number == number && isEquipmentCard(candidate.Card) {
			card = candidate
			break
		}
	}
	if card == nil {
		return nil
	}
	equip := func(replaceID string) error {
		slot, replaced, err := prepareEquipmentPlacement(ps, card, replaceID)
		if err != nil {
			return err
		}
		fromDeck := false
		if current, index := ps.FindHandCard(card.InstanceID); current == card {
			ps.RemoveFromHand(index)
		} else {
			found := false
			for i, current := range ps.Deck {
				if current == card {
					ps.Deck = append(ps.Deck[:i], ps.Deck[i+1:]...)
					found = true
					fromDeck = true
					break
				}
			}
			if !found {
				return fmt.Errorf("equipment is no longer in hand or deck")
			}
		}
		if replaced != nil {
			e.moveEquipmentToGraveyard(playerID, slot, replaced)
		}
		card.OwnerID = playerID
		card.Position = nil
		card.IsSetCounter = false
		card.IsHorizontal = true
		card.SlotIndex = slot
		card.EnterTurn = e.State.TurnNumber
		ps.Equipment[slot] = card
		if fromDeck {
			e.shuffleDeck(playerID)
		}
		e.emit(GameEvent{Type: "equip", Player: -1, Data: map[string]any{"player": playerID, "card": cardToInfo(card), "slot": slot, "free": true}})
		e.triggerEffects(TriggerOnEquip, card, nil, nil)
		e.triggerEffects(TriggerOnEnter, card, nil, nil)
		e.notifyCardEntered(playerID, card, map[string]any{"entered_player": playerID, "equipped": true, "free": true})
		return nil
	}
	if _, _, err := prepareEquipmentPlacement(ps, card, ""); err == nil {
		if equip("") == nil {
			return card
		}
		return nil
	}
	candidates := []map[string]any{}
	for _, equipment := range ps.Equipment {
		if equipment == nil {
			continue
		}
		if _, _, err := prepareEquipmentPlacement(ps, card, equipment.InstanceID); err == nil {
			candidates = append(candidates, candidateInfo(equipment, "equipment", "own"))
		}
	}
	if len(candidates) > 0 {
		e.SetPendingActionWithError(playerID, "free_equip_replacement", "免费装备：选择要替换的竖置装备", candidates, 1, 1, nil, false, func(selected []string, _ map[string]any) error { return equip(firstSelected(selected)) })
	}
	return nil
}
