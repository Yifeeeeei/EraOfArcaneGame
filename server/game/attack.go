package game

import (
	"fmt"
)

// handleAttack handles a unit attacking another unit
func (e *Engine) handleAttack(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseMain {
		return fmt.Errorf("not in main phase")
	}
	if e.State.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}
	if e.actionRestricted(RuleAttack, nil) {
		return fmt.Errorf("a field rule prevents card attacks")
	}

	attackerID, _ := action.Data["attacker_id"].(string)
	targetPos, err := requiredBoardPosition(action.Data, "target_col", "target_row")
	if err != nil {
		return err
	}
	targetCol := targetPos.Col
	targetRow := targetPos.Row

	ps := e.State.Players[playerID]
	opponent := e.State.Players[1-playerID]

	// Find attacker
	attacker := e.findUnitOnGrid(ps, attackerID)
	attackerIsEquipment := false
	if attacker == nil {
		attacker = e.findEquipment(ps, attackerID)
		attackerIsEquipment = attacker != nil
	}
	if attacker == nil {
		return fmt.Errorf("attacker not found")
	}
	if effectiveCurrentAttack(attacker) <= 0 {
		return fmt.Errorf("attacker has no attack")
	}
	if attacker.IsHorizontal {
		return fmt.Errorf("attacker is horizontal")
	}
	if e.hasEffectiveStatus(attacker, StatusStun) {
		return fmt.Errorf("attacker is stunned")
	}

	if !attackerIsEquipment {
		// Check attacker is in front row unless its own rule says otherwise.
		frontRow := ps.GetFrontRow()
		if attacker.Position == nil || (attacker.Position.Row != frontRow && (e.hasEffectiveStatus(attacker, StatusPetrify) || !cardCanAttackFromNonFront(attacker))) {
			return fmt.Errorf("attacker is not in front row")
		}
	}

	target := opponent.Units[targetCol][targetRow]
	if target == nil {
		return fmt.Errorf("no unit at target position")
	}
	if !e.isInDirectAttackRange(playerID, attacker, attackerIsEquipment, targetCol, targetRow) {
		return fmt.Errorf("target is not in attack range")
	}
	if attackCost := e.effectiveAttackCost(ps, attacker); totalElementCost(attackCost) > 0 {
		if !e.payCostForCardAction(ps, attacker, attackCost, attackCost, paymentPurposeAttack, action) {
			return fmt.Errorf("invalid attack payment")
		}
	}

	// Consume attacker (横置)
	attacker.IsHorizontal = true

	attackData := map[string]any{
		"attacker_player": playerID,
		"attacker":        attacker,
		"attack_source":   attackSourceKind(attackerIsEquipment),
		"target":          target,
		"target_pos":      targetPos,
	}

	// Trigger 攻击时 effects
	e.triggerEffects(TriggerOnAttack, attacker, target, attackData)

	// Trigger 受攻击时 effects before damage is dealt.
	e.triggerFieldEffectsWithData(TriggerOnAttacked, 1-playerID, attacker, attackData)
	e.triggerFieldEffectsWithData(TriggerOnAttacked, playerID, attacker, attackData)

	dmg := effectiveCurrentAttack(attacker)

	e.emit(GameEvent{
		Type:   "unit_attack",
		Player: -1,
		Data: map[string]any{
			"attacker_player": playerID,
			"attacker":        cardToInfo(attacker),
			"attack_source":   attackSourceKind(attackerIsEquipment),
			"target":          cardToInfo(target),
			"target_pos":      targetPos,
			"damage":          dmg,
		},
	})

	// Deal damage (unit attacks cannot be defended)
	if dmg > 0 {
		e.ApplyDamage(DamageRequest{Target: target, Amount: dmg, Kind: "attack", SourcePlayer: playerID, SourceKnown: true})
		// Trigger 命中 effects
	}

	e.checkWinCondition()
	return nil
}

func (e *Engine) resolveForcedUnitAttack(attackerOwnerID int, attacker *CardInstance, target *CardInstance, reason string) {
	if attacker == nil || target == nil || effectiveCurrentAttack(attacker) <= 0 {
		return
	}
	attackData := map[string]any{
		"attacker_player": attackerOwnerID,
		"attacker":        attacker,
		"attack_source":   "unit",
		"target":          target,
		"target_pos":      target.Position,
		"forced":          true,
		"reason":          reason,
	}
	e.triggerEffects(TriggerOnAttack, attacker, target, attackData)
	triggered := map[int]bool{}
	for _, ownerID := range []int{target.OwnerID, attackerOwnerID} {
		if ownerID < 0 || ownerID >= len(e.State.Players) || triggered[ownerID] {
			continue
		}
		triggered[ownerID] = true
		e.triggerFieldEffectsWithData(TriggerOnAttacked, ownerID, attacker, attackData)
	}

	dmg := effectiveCurrentAttack(attacker)
	e.emit(GameEvent{
		Type:   "unit_attack",
		Player: -1,
		Data: map[string]any{
			"attacker_player": attackerOwnerID,
			"attacker":        cardToInfo(attacker),
			"attack_source":   "unit",
			"target":          cardToInfo(target),
			"target_pos":      target.Position,
			"damage":          dmg,
			"forced":          true,
			"reason":          reason,
		},
	})
	e.ApplyDamage(DamageRequest{Target: target, Amount: dmg, Kind: "attack", SourcePlayer: attackerOwnerID, SourceKnown: true, ForcedAttack: true, Reason: reason})
}

func attackSourceKind(isEquipment bool) string {
	if isEquipment {
		return "equipment"
	}
	return "unit"
}

func (e *Engine) isInDirectAttackRange(playerID int, attacker *CardInstance, attackerIsEquipment bool, targetCol, targetRow int) bool {
	if attackerIsEquipment {
		return e.isEnemyFrontRowAttackTarget(playerID, attacker, targetCol, targetRow)
	}
	return e.IsInAttackRange(playerID, attacker, targetCol, targetRow)
}

func (e *Engine) isEnemyFrontRowAttackTarget(playerID int, attacker *CardInstance, targetCol, targetRow int) bool {
	opponent := e.State.Players[1-playerID]
	enemyFront := opponent.GetFrontRow()
	if enemyFront == -1 || targetRow != enemyFront {
		return false
	}
	target := opponent.Units[targetCol][targetRow]
	if target == nil {
		return false
	}
	if e.hasStealthFromOpponent(playerID, target) {
		return false
	}
	return true
}
