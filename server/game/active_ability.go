package game

import (
	"fmt"
)

// handleUseAbility handles using a card's activated ability (回合技/绝技)
func (e *Engine) handleUseAbility(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseMain {
		return fmt.Errorf("not in main phase")
	}
	if e.State.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}

	instanceID, _ := action.Data["instance_id"].(string)
	abilityType, _ := action.Data["ability_type"].(string) // "per_turn" or "ultimate"
	targetID, _ := action.Data["target_id"].(string)

	ps := e.State.Players[playerID]

	// Find the card with the ability
	card := e.findCardOnField(ps, instanceID)
	if card == nil {
		card = e.findSkill(ps, instanceID)
	}
	if card == nil {
		return fmt.Errorf("card not found on field or skill area")
	}

	if e.hasEffectiveStatus(card, StatusPetrify) {
		return fmt.Errorf("card is petrified")
	}
	if e.hasEffectiveStatus(card, StatusStun) {
		return fmt.Errorf("card is stunned")
	}

	var trigger EffectTrigger
	if abilityType == "ultimate" {
		trigger = TriggerUltimate
		if !cardHasActiveUltimate(card) {
			return fmt.Errorf("card has no active ultimate ability")
		}
		if card.UltimateUsed {
			return fmt.Errorf("ultimate already used")
		}
		if err := e.validateAbility(card, TriggerUltimate); err != nil {
			return err
		}
	} else {
		trigger = TriggerPerTurn
		if !cardHasActivePerTurn(card) {
			return fmt.Errorf("card has no active per-turn ability")
		}
		// Check if already used this turn (回合技 limit)
		maxUses := perTurnLimit(card)
		if card.UsedThisTurn >= maxUses {
			return fmt.Errorf("ability already used this turn")
		}
		if err := e.validateAbility(card, TriggerPerTurn); err != nil {
			return err
		}
	}

	// Find target if specified
	var target *CardInstance
	if targetID != "" {
		// Search both player fields for target
		for i := 0; i < 2; i++ {
			t := e.findCardOnField(e.State.Players[i], targetID)
			if t == nil {
				t = e.findUnitOnGrid(e.State.Players[i], targetID)
			}
			if t != nil {
				target = t
				break
			}
		}
	}

	// Execute the ability
	effects := globalRegistry.GetEffects(card.Card.Number, trigger)
	if len(effects) == 0 {
		return fmt.Errorf("card has no %s ability", abilityType)
	}

	ctx := &EffectContext{
		Engine:     e,
		Source:     card,
		Target:     target,
		PlayerID:   playerID,
		OpponentID: 1 - playerID,
	}

	for _, eff := range effects {
		if eff.IsActive {
			if err := eff.Handler(ctx); err != nil {
				return err
			}
		}
	}

	if abilityType == "ultimate" {
		card.UltimateUsed = true
	} else {
		card.UsedThisTurn++
	}

	e.emit(GameEvent{
		Type:   "ability_used",
		Player: -1,
		Data: map[string]any{
			"player":  playerID,
			"card":    cardToInfo(card),
			"ability": abilityType,
		},
	})

	e.checkWinCondition()
	return nil
}
