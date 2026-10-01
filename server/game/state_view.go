package game

func (e *Engine) pendingSpellAffectedPositions(spell *SpellCast) []map[string]int {
	if e == nil || spell == nil || spell.Skill == nil {
		return []map[string]int{}
	}
	defenderID := e.spellDefenderID(spell.AttackerID, spell.Skill, spell.Target)
	positions := make([]map[string]int, 0, 9+len(spell.ExtraTargets))
	seen := make(map[[3]int]bool)
	addPosition := func(ownerID, col, row int) {
		if ownerID < 0 || ownerID >= len(e.State.Players) || !validGridPosition(col, row) {
			return
		}
		key := [3]int{ownerID, col, row}
		if seen[key] {
			return
		}
		seen[key] = true
		positions = append(positions, map[string]int{"owner_id": ownerID, "col": col, "row": row})
	}

	if spell.Target.Type == "unit" && spell.Target.Position.Valid() {
		if affected, ok := e.customSpellAffectedUnits(defenderID, spell.Skill, spell.Target); ok {
			for _, unit := range affected {
				if unit != nil && unit.Position != nil {
					addPosition(unit.OwnerID, unit.Position.Col, unit.Position.Row)
				}
			}
		} else {
			col := spell.Target.Position.Col
			row := spell.Target.Position.Row
			switch e.effectiveSpellArea(spell.Skill) {
			case SpellAreaSquare:
				startCol := min(max(col, 0), 1)
				startRow := min(max(row, 0), 1)
				for areaCol := startCol; areaCol < startCol+2; areaCol++ {
					for areaRow := startRow; areaRow < startRow+2; areaRow++ {
						addPosition(defenderID, areaCol, areaRow)
					}
				}
			case SpellAreaAll:
				for areaCol := 0; areaCol < 3; areaCol++ {
					for areaRow := 0; areaRow < 3; areaRow++ {
						addPosition(defenderID, areaCol, areaRow)
					}
				}
			case SpellAreaColumn:
				for areaRow := 0; areaRow < 3; areaRow++ {
					addPosition(defenderID, col, areaRow)
				}
			case SpellAreaFrontRow:
				if frontRow := e.State.Players[defenderID].GetFrontRow(); frontRow >= 0 {
					for areaCol := 0; areaCol < 3; areaCol++ {
						addPosition(defenderID, areaCol, frontRow)
					}
				}
			case SpellAreaSplashCross:
				for _, delta := range []struct{ col, row int }{{0, 0}, {-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
					addPosition(defenderID, col+delta.col, row+delta.row)
				}
			default:
				addPosition(defenderID, col, row)
			}
		}
	}

	for _, extraTarget := range spell.ExtraTargets {
		if extraTarget.Type != "unit" || !extraTarget.Position.Valid() {
			continue
		}
		ownerID := defenderID
		if extraTarget.OwnerID != nil {
			ownerID = *extraTarget.OwnerID
		}
		addPosition(ownerID, extraTarget.Position.Col, extraTarget.Position.Row)
	}
	return positions
}

// GetStateForPlayer returns a filtered game state visible to the specified player
func (e *Engine) GetStateForPlayer(playerID int) map[string]any {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.State
	opponentID := 1 - playerID
	ps := state.Players[playerID]
	op := state.Players[opponentID]

	return map[string]any{
		"game_id":       state.GameID,
		"phase":         state.Phase.String(),
		"current_turn":  state.CurrentTurn,
		"first_player":  state.FirstPlayer,
		"turn_order":    map[string]string{"you": turnOrderLabel(playerID, state.FirstPlayer), "opponent": turnOrderLabel(opponentID, state.FirstPlayer)},
		"turn_number":   state.TurnNumber,
		"winner":        state.Winner,
		"draw_offer_by": state.DrawOfferBy,
		"you":           e.playerStateToInfo(ps, true),
		"opponent":      e.playerStateToInfo(op, false),
		"pending_spell": func() any {
			if state.PendingSpell != nil {
				return map[string]any{
					"attacker":           state.PendingSpell.AttackerID,
					"skill":              cardToInfo(state.PendingSpell.Skill),
					"target":             state.PendingSpell.Target,
					"power":              state.PendingSpell.TotalPower,
					"power_sources":      state.PendingSpell.PowerSources,
					"boost_skills":       cardsToInfo(state.PendingSpell.BoostSkills),
					"extra_targets":      state.PendingSpell.ExtraTargets,
					"affected_positions": e.pendingSpellAffectedPositions(state.PendingSpell),
				}
			}
			return nil
		}(),
		"pending_action": func() any {
			if state.PendingAction != nil && state.PendingAction.PlayerID == playerID {
				return map[string]any{
					"type":          state.PendingAction.Type,
					"player_id":     state.PendingAction.PlayerID,
					"prompt":        state.PendingAction.Prompt,
					"candidates":    state.PendingAction.Candidates,
					"min_select":    state.PendingAction.MinSelect,
					"max_select":    state.PendingAction.MaxSelect,
					"context":       state.PendingAction.Context,
					"cost":          state.PendingAction.Cost,
					"can_overexert": state.PendingAction.CanOverexert,
				}
			}
			return nil
		}(),
	}
}

// GetStateForSpectator returns the public game state without either player's
// hidden hand, skill-pool, or deck contents.
func (e *Engine) GetStateForSpectator() map[string]any {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.State
	return map[string]any{
		"game_id":       state.GameID,
		"phase":         state.Phase.String(),
		"current_turn":  state.CurrentTurn,
		"first_player":  state.FirstPlayer,
		"turn_order":    map[string]string{"you": turnOrderLabel(0, state.FirstPlayer), "opponent": turnOrderLabel(1, state.FirstPlayer)},
		"turn_number":   state.TurnNumber,
		"winner":        state.Winner,
		"draw_offer_by": state.DrawOfferBy,
		"is_spectator":  true,
		"you":           e.playerStateToInfo(state.Players[0], false),
		"opponent":      e.playerStateToInfo(state.Players[1], false),
		"pending_spell": func() any {
			if state.PendingSpell != nil {
				return map[string]any{
					"attacker":           state.PendingSpell.AttackerID,
					"skill":              cardToInfo(state.PendingSpell.Skill),
					"target":             state.PendingSpell.Target,
					"power":              state.PendingSpell.TotalPower,
					"power_sources":      state.PendingSpell.PowerSources,
					"boost_skills":       cardsToInfo(state.PendingSpell.BoostSkills),
					"extra_targets":      state.PendingSpell.ExtraTargets,
					"affected_positions": e.pendingSpellAffectedPositions(state.PendingSpell),
				}
			}
			return nil
		}(),
		"pending_action": func() any {
			if state.PendingAction != nil {
				return map[string]any{
					"type":      state.PendingAction.Type,
					"player_id": state.PendingAction.PlayerID,
					"prompt":    state.PendingAction.Prompt,
				}
			}
			return nil
		}(),
	}
}
