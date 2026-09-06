package game

// resetCards resets all cards to vertical state
func (e *Engine) resetCards(ps *PlayerState) {
	// Reset units
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if ps.Units[col][row] != nil {
				e.resetCard(ps.Units[col][row])
				ps.Units[col][row].UsedThisTurn = 0
				for _, skill := range ps.Units[col][row].BoundSkills {
					e.resetCard(skill)
					skill.UsedThisTurn = 0
				}
			}
		}
	}
	// Reset skills
	for i := range ps.Skills {
		if ps.Skills[i] != nil {
			e.resetCard(ps.Skills[i])
			ps.Skills[i].UsedThisTurn = 0
		}
	}
	// Reset equipment
	for i := range ps.Equipment {
		if ps.Equipment[i] != nil {
			e.resetCard(ps.Equipment[i])
			ps.Equipment[i].UsedThisTurn = 0
			for _, skill := range ps.Equipment[i].BoundSkills {
				e.resetCard(skill)
				skill.UsedThisTurn = 0
			}
		}
	}
}

// refreshElements calculates available elements from all vertical (竖置) cards
func (e *Engine) refreshElements(ps *PlayerState) {
	// Clear elements
	for elem := range ps.Elements {
		ps.Elements[elem] = 0
	}
	// Don't auto-gain. Elements are gained by consuming (横置) cards.
	// At turn start, elements reset to 0. Player must consume cards to gain elements.
}
