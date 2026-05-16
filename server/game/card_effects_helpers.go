package game

func findAnyUnit(ps *PlayerState) *CardInstance {
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if ps.Units[col][row] != nil {
				return ps.Units[col][row]
			}
		}
	}
	return nil
}

func findFrontRowUnit(ps *PlayerState) *CardInstance {
	frontRow := ps.GetFrontRow()
	if frontRow < 0 {
		return nil
	}
	for col := 0; col < 3; col++ {
		if ps.Units[col][frontRow] != nil {
			return ps.Units[col][frontRow]
		}
	}
	return nil
}
