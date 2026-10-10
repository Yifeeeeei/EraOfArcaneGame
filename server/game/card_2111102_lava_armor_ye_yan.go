package game

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
			_ = (Card2111102LavaArmorYeYan{}).OnDeath(&EffectContext{
				Engine:     ctx.Engine,
				Source:     armor,
				PlayerID:   ctx.PlayerID,
				OpponentID: 1 - ctx.PlayerID,
				ExtraData:  ctx.ExtraData,
			})
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
	slot := e.firstFreeEquipmentSlot(playerID)
	if slot < 0 {
		return nil
	}
	ps := e.State.Players[playerID]
	var card *CardInstance
	for i, handCard := range ps.Hand {
		if handCard != nil && handCard.Card != nil && handCard.Card.Number == number {
			card = handCard
			ps.RemoveFromHand(i)
			break
		}
	}
	if card == nil {
		for i, deckCard := range ps.Deck {
			if deckCard != nil && deckCard.Card != nil && deckCard.Card.Number == number {
				card = deckCard
				ps.Deck = append(ps.Deck[:i], ps.Deck[i+1:]...)
				break
			}
		}
	}
	if card == nil || card.Card == nil || !card.Card.IsItem() || !isEquipmentCard(card.Card) {
		return nil
	}
	card.OwnerID = playerID
	card.Position = nil
	card.IsSetCounter = false
	card.IsHorizontal = true
	card.SlotIndex = slot
	card.EnterTurn = e.State.TurnNumber
	ps.Equipment[slot] = card
	e.emit(GameEvent{
		Type:   "equip",
		Player: -1,
		Data: map[string]any{
			"player": playerID,
			"card":   cardToInfo(card),
			"slot":   slot,
			"free":   true,
		},
	})
	e.triggerEffects(TriggerOnEquip, card, nil, nil)
	e.triggerEffects(TriggerOnEnter, card, nil, nil)
	e.notifyCardEntered(playerID, card, map[string]any{"entered_player": playerID, "equipped": true, "free": true})
	return card
}
