package game

func (e *Engine) resetCards(ps *PlayerState) {
	e.resetCardOrientations(ps)
	e.refreshTurnAbilityUses(ps)
}

func (e *Engine) resetCardOrientations(ps *PlayerState) {
	for _, card := range e.getAllFieldCards(ps) {
		e.resetCard(card)
		for _, skill := range card.BoundSkills {
			e.resetCard(skill)
		}
	}
}

func (e *Engine) refreshTurnAbilityUses(ps *PlayerState) {
	for _, card := range e.getAllFieldCards(ps) {
		card.UsedThisTurn = 0
		for _, skill := range card.BoundSkills {
			skill.UsedThisTurn = 0
		}
	}
}

func (e *Engine) refreshElements(ps *PlayerState) {
	for elem := range ps.Elements {
		ps.Elements[elem] = 0
	}
}
