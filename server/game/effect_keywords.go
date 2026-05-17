package game

// ApplyKeywordOnEnter applies keyword effects when a card enters the field
func (e *Engine) ApplyKeywordOnEnter(card *CardInstance) {
	if cardHasRush(card) {
		card.IsHorizontal = false
	}

	if stealthLayers := cardStealthLayers(card); stealthLayers > 0 {
		card.Statuses["隐蔽"] = stealthLayers
	}

	if cardHasTaunt(card) {
		card.Statuses["引魔"] = 1
	}

	if shieldVal := cardShieldLayers(card); shieldVal > 0 {
		card.Statuses["护盾"] = shieldVal
	}
}

// ApplyKeywordOnSkillUse applies keyword effects when a skill is used
func (e *Engine) ApplyKeywordOnSkillUse(skill *CardInstance) {
	if cd := skillCooldown(skill); cd > 0 {
		skill.Statuses[StatusCooldown] = cd
	}
}

// IsInSpellRange checks if a target position is in the caster's spell range
// Default spell range: all friendly units + enemy front row
// 穿透 (Pierce): can target any unit
// 隐蔽 (Stealth): can't be targeted unless caster has 穿透
// 引魔 (Taunt): always in range
// 屏蔽 (Shielding): blocks enemy pierce
func (e *Engine) IsInSpellRange(casterID int, targetCol, targetRow int, hasPierce bool) bool {
	opponent := e.State.Players[1-casterID]
	caster := e.State.Players[casterID]

	for _, card := range e.getAllFieldCards(caster) {
		if card != nil && card.Card != nil && !e.hasEffectiveStatus(card, StatusPetrify) {
			if h, ok := behaviorForNumber(card.Card.Number).(GlobalSpellRangeBehavior); ok && h.HasActiveGlobalSpellRange(card) && h.HasGlobalSpellRange() {
				return true
			}
		}
	}

	// Check if opponent has shielding active (屏蔽)
	if hasPierce {
		opponentCards := e.getAllFieldCards(opponent)
		for _, c := range opponentCards {
			if cardHasShielding(c) && !e.hasEffectiveStatus(c, StatusPetrify) {
				hasPierce = false // shielding blocks pierce
				break
			}
		}
	}

	if hasPierce {
		return true // pierce can target anything
	}

	// Default: enemy front row only
	frontRow := opponent.GetFrontRow()
	if frontRow == -1 {
		return true // no enemy units, any position is valid
	}

	target := opponent.Units[targetCol][targetRow]
	if target == nil {
		return false
	}

	// 引魔 (Taunt) units are always in range
	if target.Statuses["引魔"] > 0 {
		return true
	}

	// 隐蔽 (Stealth) units can't be targeted
	if target.Statuses["隐蔽"] > 0 {
		return false
	}

	return targetRow == frontRow
}

// IsInAttackRange checks if a unit can attack a target
// Only front row units can attack, and only enemy front row
func (e *Engine) IsInAttackRange(attackerID int, attacker *CardInstance, targetCol, targetRow int) bool {
	ps := e.State.Players[attackerID]
	opponent := e.State.Players[1-attackerID]

	// Attacker must be in front row
	frontRow := ps.GetFrontRow()
	if attacker.Position == nil || attacker.Position.Row != frontRow {
		return false
	}

	// Target must be in opponent's front row
	enemyFront := opponent.GetFrontRow()
	if enemyFront == -1 {
		return false
	}

	target := opponent.Units[targetCol][targetRow]
	if target == nil {
		return false
	}

	// 隐蔽 units can't be attacked unless attacker has 穿透
	if target.Statuses["隐蔽"] > 0 {
		if !cardHasPierce(attacker) {
			return false
		}
	}

	return targetRow == enemyFront
}

// HandleTemporaryUnits removes 临时 (temporary) units at end of turn
func (e *Engine) HandleTemporaryUnits(ps *PlayerState) {
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := ps.Units[col][row]
			if unit != nil && cardIsTemporary(unit) {
				e.destroyUnit(unit, ps.PlayerID)
			}
		}
	}
}

// HandleShieldDecay reduces shield by 1 at end of turn
func (e *Engine) HandleShieldDecay(ps *PlayerState) {
	allCards := e.getAllFieldCards(ps)
	for _, card := range allCards {
		if card.Statuses["护盾"] > 0 {
			card.Statuses["护盾"]--
		}
		// Stealth decays too
		if card.Statuses["隐蔽"] > 0 {
			card.Statuses["隐蔽"]--
		}
	}
}

// ApplyShieldDamage reduces damage by shield amount
// Returns remaining damage after shield
func ApplyShieldDamage(target *CardInstance, damage int) int {
	shield := target.Statuses["护盾"]
	if shield <= 0 {
		return damage
	}
	if damage <= shield {
		target.Statuses["护盾"] -= 1 // shield absorbs hit, lose 1 layer
		return 0
	}
	target.Statuses["护盾"] = 0
	return damage - shield
}
