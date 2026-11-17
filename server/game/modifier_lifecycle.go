package game

// Lifetime boundaries are evaluated by turn cleanup, not by UI/stat queries.
// ExpiresTurn is inclusive: the modifier participates in that turn's effects,
// then leaves during cleanup. Player-specific expiry is a separate boundary.
func (m TemporaryModifier) expiresAfterTurn(turn int) bool {
	return m.ExpiresTurn > 0 && m.ExpiresTurn <= turn
}

func (m TemporaryModifier) expiresAfterPlayerTurn(playerID int) bool {
	return m.ExpiresAtTurnEnd && m.ExpiresOnPlayerID == playerID
}

// spendUse only runs when an operation commits. Zero means no uses remain;
// unlimited duration is expressed by a continuous modifier kind, never by
// decrementing a counter below zero. Queries must not call this method.
func (m *TemporaryModifier) spendUse() bool {
	if m == nil || m.RemainingUses <= 0 {
		return false
	}
	m.RemainingUses--
	return true
}
