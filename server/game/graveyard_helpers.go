package game

func (e *Engine) reviveCompanionFromGraveyard(playerID int, instanceID string) bool {
	ps := e.State.Players[playerID]
	pos := ps.FindEmptyPosition()
	if pos == nil {
		return false
	}
	for i, card := range ps.Graveyard {
		if card == nil || card.InstanceID != instanceID || !card.Card.IsCompanion() {
			continue
		}
		if !ps.CanPayCost(card.Card.ElementsCost) {
			return false
		}
		ps.PayCost(card.Card.ElementsCost)
		ps.Graveyard = append(ps.Graveyard[:i], ps.Graveyard[i+1:]...)
		card.CurrentLife = max(card.CurrentLife, card.Card.Life)
		card.IsHorizontal = true
		card.Position = pos
		ps.Units[pos.Col][pos.Row] = card
		e.emit(GameEvent{Type: "summon", Player: -1, Data: map[string]any{
			"player":   playerID,
			"card":     cardToInfo(card),
			"position": pos,
			"elements": ps.Elements,
		}})
		e.triggerEffects(TriggerOnEnter, card, nil, nil)
		return true
	}
	return false
}
