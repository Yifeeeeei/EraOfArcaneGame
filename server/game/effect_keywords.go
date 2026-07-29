package game

// ApplyKeywordOnEnter applies keyword effects when a card enters the field
func (e *Engine) ApplyKeywordOnEnter(card *CardInstance) {
	if cardHasRush(card) {
		card.IsHorizontal = false
	}

	if stealthLayers := cardStealthLayers(card); stealthLayers > 0 {
		card.Statuses[StatusStealth] = stealthLayers
	}

	if shieldLayers := cardShieldLayers(card); shieldLayers > 0 {
		e.gainPlayerShield(card.OwnerID, shieldLayers)
	}

	if cardHasTaunt(card) {
		card.Statuses["引魔"] = 1
	}
}

func (e *Engine) ApplySummonModifiersOnEnter(card *CardInstance) {
	if card == nil || card.Card == nil || !card.Card.IsCompanion() || card.OwnerID < 0 || card.OwnerID >= len(e.State.Players) {
		return
	}
	ps := e.State.Players[card.OwnerID]
	if ps.NextCompanionStealth <= 0 {
		return
	}
	e.grantStealth(card, ps.NextCompanionStealth)
	ps.NextCompanionStealth = 0
}

func (e *Engine) grantStealth(card *CardInstance, amount int) bool {
	return e.addStatus(card, StatusStealth, amount)
}

// ApplyKeywordOnSkillUse applies keyword effects when a skill is used
func (e *Engine) ApplyKeywordOnSkillUse(skill *CardInstance) {
	if cd := skillCooldown(skill); cd > 0 {
		skill.Statuses[StatusCooldown] = cd
	}
	if skill != nil && skill.Card != nil && skill.Card.Duration > 0 {
		skill.Statuses[StatusAbilityDuration] = skill.Card.Duration
	}
}

// IsInSpellRange checks if a target position is in the caster's spell range
// Default spell range: all friendly units + enemy front row
// 穿透 (Pierce): can target any unit
// 隐蔽 (Stealth): can't be targeted by opponents unless explicitly allowed
// 引魔 (Taunt): always in range
// 屏蔽 (Shielding): blocks enemy pierce
func (e *Engine) IsInSpellRange(casterID int, targetCol, targetRow int, hasPierce bool) bool {
	opponent := e.State.Players[1-casterID]
	caster := e.State.Players[casterID]

	target := opponent.Units[targetCol][targetRow]
	if target != nil && e.hasStealthFromOpponent(casterID, target) {
		return false
	}

	for _, card := range e.getAllFieldCards(caster) {
		if card != nil && card.Card != nil && !e.hasEffectiveStatus(card, StatusPetrify) {
			if cardHasActiveGlobalSpellRange(card) {
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
		return true // pierce can target any non-stealth unit or area
	}

	// Default: enemy front row only. Stealth units do not block spell range.
	frontRow := e.spellRangeFrontRowAgainst(opponent, casterID)
	if frontRow == -1 {
		return true // no enemy units, any position is valid
	}

	if target == nil {
		return false
	}

	// 引魔 (Taunt) units are always in range
	if target.Statuses["引魔"] > 0 || cardHasTaunt(target) {
		return true
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

	// 隐蔽 has priority over pierce and blocks opposing direct attacks.
	if e.hasStealthFromOpponent(attackerID, target) {
		return false
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

// HandleShieldDecay reduces player shield and stealth by 1 at end of turn.
func (e *Engine) HandleShieldDecay(ps *PlayerState) {
	if ps.Shield > 0 && !e.playerShieldDecayPrevented(ps) {
		ps.Shield--
	}
	allCards := e.getAllFieldCards(ps)
	for _, card := range allCards {
		if card.Statuses[StatusStealth] > 0 {
			card.Statuses[StatusStealth]--
		}
	}
}

func (e *Engine) playerShieldDecayPrevented(ps *PlayerState) bool {
	if ps == nil {
		return false
	}
	if e.playerHasActiveCard(ps, "2011101") {
		return true
	}
	if ps.Shield >= 3 {
		return false
	}
	return e.playerHasActiveCard(ps, "4411101")
}

func (e *Engine) playerHasActiveCard(ps *PlayerState, number string) bool {
	if ps == nil || number == "" {
		return false
	}
	for _, card := range e.getAllFieldCards(ps) {
		if card != nil && card.Card != nil && card.Card.Number == number && !e.hasEffectiveStatus(card, StatusPetrify) {
			return true
		}
	}
	return false
}

func (e *Engine) hasStealthFromOpponent(playerID int, target *CardInstance) bool {
	if target == nil || target.OwnerID == playerID {
		return false
	}
	return e.hasActiveStealth(target)
}

func (e *Engine) hasActiveStealth(card *CardInstance) bool {
	return card != nil && !e.hasEffectiveStatus(card, StatusPetrify) && e.hasEffectiveStatus(card, StatusStealth)
}

func (e *Engine) spellRangeFrontRowAgainst(ps *PlayerState, casterID int) int {
	if ps == nil {
		return -1
	}
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			unit := ps.Units[col][row]
			if unit == nil {
				continue
			}
			if unit.OwnerID != casterID && e.hasActiveStealth(unit) {
				continue
			}
			return row
		}
	}
	return -1
}

func (e *Engine) gainPlayerShield(playerID int, amount int) {
	if amount <= 0 || playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	ps := e.State.Players[playerID]
	if ps == nil || ps.CannotGainShield {
		return
	}
	ps.Shield += amount
	e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
		"effect": "gain_shield",
		"amount": amount,
		"shield": ps.Shield,
	}})
}

func (e *Engine) losePlayerShield(playerID int, amount int) {
	if amount <= 0 || playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	ps := e.State.Players[playerID]
	if ps == nil || ps.Shield <= 0 {
		return
	}
	lost := min(amount, ps.Shield)
	ps.Shield -= lost
	e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
		"effect": "lose_shield",
		"amount": lost,
		"shield": ps.Shield,
	}})
}

func (e *Engine) gainStrictArcane(playerID int, amount int) {
	if amount <= 0 || playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	e.State.Players[playerID].StrictArcane += amount
	e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
		"effect":        "gain_strict_arcane",
		"amount":        amount,
		"strict_arcane": e.State.Players[playerID].StrictArcane,
	}})
}

func (e *Engine) spendStrictArcane(playerID int, amount int) bool {
	if amount < 0 || playerID < 0 || playerID >= len(e.State.Players) {
		return false
	}
	ps := e.State.Players[playerID]
	if ps.StrictArcane < amount {
		return false
	}
	ps.StrictArcane -= amount
	e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
		"effect":        "spend_strict_arcane",
		"amount":        amount,
		"strict_arcane": ps.StrictArcane,
	}})
	return true
}

func (e *Engine) applyPlayerShieldDamage(target *CardInstance, damage int, damageData map[string]any) int {
	if target == nil || damage <= 0 || damageData == nil || damageData["damage_source"] != "spell" {
		return damage
	}
	if skip, _ := damageData["skip_player_shield"].(bool); skip {
		return damage
	}
	attacker, ok := damageData["attacker"].(int)
	if !ok || attacker == target.OwnerID {
		return damage
	}
	ps := e.State.Players[target.OwnerID]
	if ps == nil || ps.Shield <= 0 {
		return damage
	}
	prevented := min(damage, ps.Shield)
	ps.Shield -= prevented
	remaining := damage - prevented
	e.emit(GameEvent{
		Type:   "shield_block",
		Player: -1,
		Data: map[string]any{
			"target":    cardToInfo(target),
			"owner":     target.OwnerID,
			"attacker":  attacker,
			"prevented": prevented,
			"shield":    ps.Shield,
			"broken":    ps.Shield == 0,
		},
	})
	return remaining
}
