package game

// TurnNumber counts rounds. This ordinal advances on either player's turn and
// is only used for effects whose text refers to a single, unqualified turn.
func (e *Engine) turnOrdinal() int {
	ordinal := max(0, e.State.TurnNumber-1) * 2
	if e.State.CurrentTurn != e.State.FirstPlayer {
		ordinal++
	}
	return ordinal + 1
}
