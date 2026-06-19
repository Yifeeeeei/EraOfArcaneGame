package game

import "fmt"

const illusionScrollDoneID = "illusion_scroll_done"

func (e *Engine) startIllusionScrollRearrange(playerID int, source *CardInstance, extraData map[string]any) {
	if e.State.PendingSpell == nil {
		return
	}
	e.promptIllusionScrollUnit(playerID, source, extraData)
}

func (e *Engine) promptIllusionScrollUnit(playerID int, source *CardInstance, extraData map[string]any) {
	if e.State.PendingSpell == nil {
		return
	}
	candidates := e.friendlyUnits(playerID, true, nil)
	candidates = append(candidates, map[string]any{
		"instance_id": illusionScrollDoneID,
		"name":        "完成重排",
		"zone":        "choice",
		"side":        "own",
	})
	e.SetPendingAction(playerID, "illusion_scroll_unit",
		"幻术卷轴:选择1个己方单位移动/交换，或完成重排", candidates, 1, 1,
		func(selected []string) {
			if firstSelected(selected) == illusionScrollDoneID {
				e.promptIllusionScrollRetarget(playerID, source, extraData)
				return
			}
			unit := selectedUnitFromCandidates(e, selected, candidates)
			if unit == nil {
				return
			}
			e.promptIllusionScrollPosition(playerID, source, extraData, unit.InstanceID)
		})
}

func (e *Engine) promptIllusionScrollPosition(playerID int, source *CardInstance, extraData map[string]any, unitID string) {
	if e.State.PendingSpell == nil {
		return
	}
	positions := e.allUnitPositionsForPlayer(playerID, playerID)
	e.SetPendingAction(playerID, "illusion_scroll_position",
		"幻术卷轴:选择移动/交换到的位置", positions, 1, 1,
		func(selected []string) {
			pos, ok := positionFromSelectionID(firstSelected(selected))
			if ok {
				e.moveOrSwapUnitToPosition(playerID, unitID, pos)
			}
			e.promptIllusionScrollUnit(playerID, source, extraData)
		})
}

func (e *Engine) promptIllusionScrollRetarget(playerID int, source *CardInstance, extraData map[string]any) {
	if e.State.PendingSpell == nil {
		return
	}
	spell := e.State.PendingSpell
	attackerID := spell.AttackerID
	candidates := e.spellTargetCandidates(attackerID, spell.Skill)
	if len(candidates) == 0 {
		markSpellCastCancelled(extraData)
		e.cancelPendingSpell(playerID, source, "illusion_scroll_no_legal_target")
		return
	}
	e.SetPendingAction(attackerID, "illusion_scroll_retarget",
		"幻术卷轴:重新选择法术攻击目标", candidates, 1, 1,
		func(selected []string) {
			target := selectedUnitFromCandidates(e, selected, candidates)
			if target == nil || target.Position == nil || e.State.PendingSpell == nil {
				return
			}
			spellTarget := SpellTarget{Type: "unit", Position: *target.Position}
			if err := e.validateSpellTarget(attackerID, spell.Skill, spellTarget); err != nil {
				markSpellCastCancelled(extraData)
				e.cancelPendingSpell(playerID, source, "illusion_scroll_invalid_retarget")
				return
			}
			e.State.PendingSpell.Target = spellTarget
			e.State.PendingSpell.ExtraTargets = nil
		})
}

func (e *Engine) allUnitPositionsForPlayer(boardPlayerID int, viewerID int) []map[string]any {
	if boardPlayerID < 0 || boardPlayerID >= len(e.State.Players) {
		return nil
	}
	side := "own"
	if boardPlayerID != viewerID {
		side = "enemy"
	}
	candidates := make([]map[string]any, 0, 9)
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			pos := Position{Col: col, Row: row}
			candidates = append(candidates, map[string]any{
				"instance_id": positionSelectionID(pos),
				"name":        fmt.Sprintf("位置 (%d,%d)", col, row),
				"zone":        "unit_position",
				"side":        side,
				"position":    pos,
			})
		}
	}
	return candidates
}

func (e *Engine) moveOrSwapUnitToPosition(playerID int, instanceID string, pos Position) bool {
	if !pos.Valid() || playerID < 0 || playerID >= len(e.State.Players) {
		return false
	}
	ps := e.State.Players[playerID]
	unit := e.findUnitOnGrid(ps, instanceID)
	if unit == nil || unit.Position == nil {
		return false
	}
	from := *unit.Position
	other := ps.Units[pos.Col][pos.Row]
	ps.Units[pos.Col][pos.Row] = unit
	unit.Position = &Position{Col: pos.Col, Row: pos.Row}
	ps.Units[from.Col][from.Row] = other
	if other != nil {
		other.Position = &Position{Col: from.Col, Row: from.Row}
	}
	return true
}

func markSpellCastCancelled(data map[string]any) {
	if data == nil {
		return
	}
	if cancelled, ok := data["cancel_spell_cast"].(*bool); ok && cancelled != nil {
		*cancelled = true
		return
	}
	cancelled := true
	data["cancel_spell_cast"] = &cancelled
}
