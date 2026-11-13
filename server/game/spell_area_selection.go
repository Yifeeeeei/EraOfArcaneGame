package game

import "fmt"

// Effects that cast another spell choose its square before the shared cast
// pipeline commits resources, bonuses, or a new defense window.
func (e *Engine) chooseEffectSpellSquare(playerID int, skill *CardInstance, target SpellTarget, cast func(SpellTarget) error) (bool, error) {
	if err := e.validateSpellAreaSelection(skill, target); err != nil {
		return true, err
	}
	if target.Type != "unit" || e.effectiveSpellArea(skill) != SpellAreaSquare || target.AreaOrigin != nil {
		return false, nil
	}
	origins := map[string]Position{}
	candidates := []map[string]any{}
	for col := 0; col < 2; col++ {
		for row := 0; row < 2; row++ {
			if target.Position.Col < col || target.Position.Col > col+1 || target.Position.Row < row || target.Position.Row > row+1 {
				continue
			}
			id := fmt.Sprintf("square:%d:%d", col, row)
			origins[id] = Position{Col: col, Row: row}
			candidates = append(candidates, map[string]any{"instance_id": id, "name": fmt.Sprintf("%d–%d列 · %d–%d排方阵", col+1, col+2, row+1, row+2), "zone": "area"})
		}
	}
	if len(candidates) == 1 {
		origin := origins[candidates[0]["instance_id"].(string)]
		target.AreaOrigin = &origin
		return true, cast(target)
	}
	if len(candidates) == 0 {
		return true, fmt.Errorf("invalid square center")
	}
	e.SetPendingActionWithErrorAndContext(playerID, "spell_square_area", "选择"+skill.Card.Name+"的方阵", candidates, 1, 1, nil, false, map[string]any{"source": cardToInfo(skill), "center": target.Position}, func(selected []string, _ map[string]any) error {
		origin, ok := origins[firstSelected(selected)]
		if !ok {
			return fmt.Errorf("invalid square area")
		}
		target.AreaOrigin = &origin
		return cast(target)
	})
	return true, nil
}

func squareAreaOrigin(target SpellTarget) Position {
	if target.AreaOrigin != nil {
		return *target.AreaOrigin
	}
	// Legacy/headless callers retain the previous deterministic default.
	return Position{Col: min(target.Position.Col, 1), Row: min(target.Position.Row, 1)}
}

func (e *Engine) spellAreaSelection(skill *CardInstance, target SpellTarget, action ActionMessage) (SpellTarget, error) {
	origin, supplied, err := optionalBoardPosition(action.Data, "area_col", "area_row")
	if err != nil {
		return target, err
	}
	if supplied {
		target.AreaOrigin = &origin
	}
	if err := e.validateSpellAreaSelection(skill, target); err != nil {
		return target, err
	}
	return target, nil
}

func (e *Engine) validateSpellAreaSelection(skill *CardInstance, target SpellTarget) error {
	if target.AreaOrigin == nil {
		return nil
	}
	pos := *target.AreaOrigin
	if e.effectiveSpellArea(skill) != SpellAreaSquare || target.Type != "unit" || pos.Col < 0 || pos.Col > 1 || pos.Row < 0 || pos.Row > 1 || target.Position.Col < pos.Col || target.Position.Col > pos.Col+1 || target.Position.Row < pos.Row || target.Position.Row > pos.Row+1 {
		return fmt.Errorf("square area must contain its selected center")
	}
	return nil
}

func (e *Engine) isAreaCenterInRange(playerID int, skill *CardInstance, target SpellTarget, pierce bool) bool {
	ps, op := e.State.Players[playerID], e.State.Players[1-playerID]
	for _, source := range e.getAllFieldCards(ps) {
		if !e.hasEffectiveStatus(source, StatusPetrify) && cardHasActiveGlobalSpellRange(source) {
			return true
		}
	}
	for _, source := range e.getAllFieldCards(op) {
		if cardHasShielding(source) && !e.hasEffectiveStatus(source, StatusPetrify) {
			pierce = false
		}
	}
	if pierce {
		return true
	}
	unit := op.Units[target.Position.Col][target.Position.Row]
	if unit != nil && (unit.Statuses["引魔"] > 0 || cardHasTaunt(unit)) {
		return true
	}
	front := op.GetFrontRow()
	return front < 0 || target.Position.Row <= front
}
