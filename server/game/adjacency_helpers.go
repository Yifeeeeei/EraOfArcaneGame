package game

func adjacentUnits(ps *PlayerState, pos *Position) []*CardInstance {
	if ps == nil || pos == nil {
		return nil
	}
	units := make([]*CardInstance, 0, 4)
	for _, delta := range []struct{ col, row int }{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
		col := pos.Col + delta.col
		row := pos.Row + delta.row
		if col < 0 || col >= 3 || row < 0 || row >= 3 {
			continue
		}
		if ps.Units[col][row] != nil {
			units = append(units, ps.Units[col][row])
		}
	}
	return units
}
