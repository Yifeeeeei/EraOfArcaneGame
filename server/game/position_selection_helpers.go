package game

import (
	"fmt"
	"strconv"
	"strings"
)

func positionSelectionID(pos Position) string {
	return fmt.Sprintf("pos:%d:%d", pos.Col, pos.Row)
}

func positionFromSelectionID(id string) (Position, bool) {
	parts := strings.Split(id, ":")
	if len(parts) != 3 || parts[0] != "pos" {
		return Position{}, false
	}
	col, err := strconv.Atoi(parts[1])
	if err != nil {
		return Position{}, false
	}
	row, err := strconv.Atoi(parts[2])
	if err != nil {
		return Position{}, false
	}
	pos := Position{Col: col, Row: row}
	return pos, pos.Valid()
}

func (e *Engine) friendlyEmptyUnitPositions(playerID int) []map[string]any {
	return e.emptyUnitPositionsForPlayer(playerID, playerID)
}

func (e *Engine) emptyUnitPositionsForPlayer(boardPlayerID int, viewerID int) []map[string]any {
	if boardPlayerID < 0 || boardPlayerID >= len(e.State.Players) {
		return nil
	}
	ps := e.State.Players[boardPlayerID]
	candidates := make([]map[string]any, 0, 9)
	side := "own"
	if boardPlayerID != viewerID {
		side = "enemy"
	}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if ps.Units[col][row] != nil {
				continue
			}
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
