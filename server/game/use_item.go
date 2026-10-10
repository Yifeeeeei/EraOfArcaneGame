package game

import (
	"eraofarcane/cards"
	"fmt"
)

// handleUseItem handles using a consumable item from hand
func (e *Engine) handleUseItem(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseMain {
		return fmt.Errorf("not in main phase")
	}
	if e.State.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}
	if e.actionRestricted(RuleUseItem, nil) {
		return fmt.Errorf("a field rule prevents playing cards")
	}

	instanceID, _ := action.Data["instance_id"].(string)
	ps := e.State.Players[playerID]

	card, handIdx := ps.FindHandCard(instanceID)
	if card == nil {
		return fmt.Errorf("card not found in hand")
	}
	if !card.Card.IsItem() {
		return fmt.Errorf("card is not an item")
	}
	if isCounterTrapCard(card.Card.Number) {
		return e.placeCounterTrap(playerID, card, handIdx)
	}

	// Check if this is a terrain card - terrain cards go to battlefield
	if cards.IsTerrain(card.Card.Number) {
		return e.handlePlaceTerrain(playerID, action)
	}
	if spellScrollUsesGenericCast(card.Card) {
		return e.handleUseSpellScrollItem(playerID, action, card, handIdx)
	}
	if cards.IsEquipment(card.Card.Number) {
		return fmt.Errorf("equipment cannot be used as a consumable item")
	}
	if !cards.IsConsumable(card.Card.Number) {
		if behavior, ok := globalRegistry.GetBehavior(card.Card.Number).(OnUseItemBehavior); !ok || !behavior.HasActiveUseItem(card) {
			return fmt.Errorf("item is not consumable")
		}
	}
	if err := e.validateConsumableItemUse(playerID, card); err != nil {
		return err
	}

	// Regular consumable item
	cost := e.effectiveCardPlayCost(ps, card)
	if !e.canPayCost(ps, cost) {
		return fmt.Errorf("not enough elements")
	}

	// Pay and use
	if !e.payCostForAction(ps, cost, action) {
		return fmt.Errorf("invalid payment")
	}
	e.notifyCardPlayCostPaid(ps, card)
	ps.RemoveFromHand(handIdx)
	e.addToGraveyard(playerID, card)

	e.emit(GameEvent{
		Type:   "use_item",
		Player: -1,
		Data: map[string]any{
			"player":   playerID,
			"card":     cardToInfo(card),
			"elements": ps.Elements,
		},
	})

	cancelled := false
	useData := map[string]any{"used_player": playerID, "cancel_item": &cancelled}
	resolveItem := func() {
		if cancelled {
			return
		}
		e.triggerEffects(TriggerOnUseItem, card, nil, nil)
		e.triggerFieldEffectsWithData(TriggerOnUseItem, playerID, card, useData)
		e.triggerFieldEffectsWithData(TriggerOnUseItem, 1-playerID, card, useData)
	}
	if e.promptOpponentCounterTrap(playerID, TriggerOnUseItem, card, useData, resolveItem) {
		return nil
	}
	resolveItem()

	return nil
}

func (e *Engine) validateConsumableItemUse(playerID int, card *CardInstance) error {
	behavior, ok := cardBehavior(card).(ItemUseValidationBehavior)
	if !ok {
		return nil
	}
	return behavior.ValidateItemUse(&EffectContext{Engine: e, Source: card, PlayerID: playerID, OpponentID: 1 - playerID})
}

func (e *Engine) handleUseSpellScrollItem(playerID int, action ActionMessage, card *CardInstance, handIdx int) error {
	ps := e.State.Players[playerID]
	if isDefenseOnlySkill(card.Card) {
		return fmt.Errorf("defense spell scroll can only be used during a defense window")
	}

	target := SpellTarget{Type: "none"}
	if skillNeedsTargetInstance(card) {
		targetType, _ := action.Data["target_type"].(string)
		if targetType == "hero" {
			target = SpellTarget{Type: "hero"}
			if err := e.validateSpellTarget(playerID, card, target); err != nil {
				return err
			}
		} else {
			targetPos, err := requiredBoardPosition(action.Data, "target_col", "target_row")
			if err != nil {
				return fmt.Errorf("spell scroll requires a target: %w", err)
			}
			target = SpellTarget{Type: "unit", Position: targetPos}
			if ownerF, ok := action.Data["target_owner"].(float64); ok {
				ownerID := int(ownerF)
				target.OwnerID = &ownerID
			}
			if err := e.validateSpellTarget(playerID, card, target); err != nil {
				return err
			}
		}
	}
	oracleGloryBonus, err := e.validateOracleGlorySupport(playerID, card, action)
	if err != nil {
		return err
	}
	flameArraySacrifice, flameArrayBonus, err := e.validateFlameArrayScrollSacrifice(playerID, card, action)
	if err != nil {
		return err
	}

	cost := e.effectiveCardPlayCost(ps, card)
	if !e.canPayCost(ps, cost) {
		return fmt.Errorf("not enough elements")
	}
	if !e.payCostForAction(ps, cost, action) {
		return fmt.Errorf("invalid payment")
	}
	e.notifyCardPlayCostPaid(ps, card)
	ps.RemoveFromHand(handIdx)
	e.addToGraveyard(playerID, card)

	e.emit(GameEvent{
		Type:   "use_item",
		Player: -1,
		Data: map[string]any{
			"player":   playerID,
			"card":     cardToInfo(card),
			"elements": ps.Elements,
		},
	})

	cancelled := false
	useData := map[string]any{"used_player": playerID, "cancel_item": &cancelled, "spell_scroll": true}
	resolveItem := func() {
		if cancelled {
			return
		}
		if oracleGloryBonus > 0 {
			e.addTemporaryModifier(playerID, TemporaryModifier{
				Type:             TempModNextAttackSpellPowerBonus,
				SourceCardNumber: card.Card.Number,
				SourceName:       card.Card.Name,
				TargetInstanceID: card.InstanceID,
				Amount:           oracleGloryBonus,
				RemainingUses:    1,
			})
		}
		if flameArraySacrifice != nil {
			e.destroyUnitWithCause(flameArraySacrifice, playerID, DeathCauseSacrifice)
			e.addTemporaryModifier(playerID, TemporaryModifier{
				Type:             TempModNextAttackSpellPowerBonus,
				SourceCardNumber: card.Card.Number,
				SourceName:       card.Card.Name,
				TargetInstanceID: card.InstanceID,
				Amount:           flameArrayBonus,
				RemainingUses:    1,
			})
		}
		e.startSpellScrollCast(playerID, card, target)
	}
	if e.promptOpponentCounterTrap(playerID, TriggerOnUseItem, card, useData, resolveItem) {
		return nil
	}
	resolveItem()
	return nil
}

func (e *Engine) startSpellScrollCast(playerID int, scroll *CardInstance, target SpellTarget) {
	ps := e.State.Players[playerID]
	boostSkills := []*CardInstance{}
	e.applyCoralBellyFirstSpellAttackBonus(playerID, scroll)
	totalPower := e.effectiveSpellPower(playerID, scroll, boostSkills, target)
	powerSources := e.spellPowerSources(playerID, scroll, boostSkills, totalPower, target)
	e.consumeNextSpellPowerBonuses(ps, scroll)

	e.recordSpellCast(playerID, scroll)
	spellCastData := map[string]any{
		"cast_player":  playerID,
		"attacker":     playerID,
		"purpose":      string(skillPurposeAttack),
		"skill":        cardToInfo(scroll),
		"target":       target,
		"power":        totalPower,
		"boost_count":  0,
		"is_sorcery":   false,
		"spell_scroll": true,
	}
	e.emit(GameEvent{Type: "spell_cast", Player: -1, Data: spellCastData})
	e.triggerEffects(TriggerOnSpellCast, scroll, nil, spellCastData)

	e.State.PendingSpell = &SpellCast{
		AttackerID:   playerID,
		Skill:        scroll,
		Target:       target,
		TotalPower:   totalPower,
		PowerSources: powerSources,
		BoostSkills:  boostSkills,
	}
	openDefenseWindow := func() {
		if e.State.PendingSpell == nil {
			return
		}
		e.State.ResumePhase = PhaseDefenseWindow
		e.State.Phase = PhaseDefenseWindow
		e.emit(GameEvent{Type: "defense_window", Player: 1 - playerID, Data: map[string]any{"timeout": 30}})
	}
	continueSpell := openDefenseWindow
	if !e.spellAllowsDefense(playerID, scroll, target) {
		continueSpell = func() {
			e.resolvePendingSpellHit()
		}
	}
	if e.triggerSpellCastFieldEffectsWithContinuation(playerID, scroll, spellCastData, continueSpell) {
		if e.spellAllowsDefense(playerID, scroll, target) {
			e.State.ResumePhase = PhaseDefenseWindow
		}
		return
	}
	continueSpell()
}
