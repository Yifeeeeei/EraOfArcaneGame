package game

import (
	"eraofarcane/cards"
	"fmt"
)

// handlePlaceTerrain handles placing a terrain card (地形牌) on the battlefield
func (e *Engine) handlePlaceTerrain(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseMain {
		return fmt.Errorf("not in main phase")
	}
	if e.State.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}
	if e.actionRestricted(RulePlaceTerrain, nil) {
		return fmt.Errorf("a field rule prevents playing cards")
	}

	instanceID, _ := action.Data["instance_id"].(string)
	pos, err := requiredBoardPosition(action.Data, "col", "row")
	if err != nil {
		return err
	}
	col := pos.Col
	row := pos.Row

	ps := e.State.Players[playerID]

	// Find card in hand
	card, handIdx := ps.FindHandCard(instanceID)
	if card == nil {
		return fmt.Errorf("card not found in hand")
	}

	// Must be an item with terrain keyword
	if !card.Card.IsItem() {
		return fmt.Errorf("card is not an item")
	}
	if !cards.IsTerrain(card.Card.Number) {
		return fmt.Errorf("card is not a terrain")
	}

	if ps.Terrain[col][row] != nil {
		return fmt.Errorf("position already has terrain")
	}

	// Check cost
	cost := e.effectiveCardPlayCost(ps, card)
	if !e.canPayCost(ps, cost) {
		return fmt.Errorf("not enough elements")
	}

	// Pay cost and place
	if !e.payCostForAction(ps, cost, action) {
		return fmt.Errorf("invalid payment")
	}
	e.notifyCardPlayCostPaid(ps, card)
	ps.RemoveFromHand(handIdx)
	card.Position = &Position{Col: col, Row: row}
	card.EnterTurn = e.State.TurnNumber
	ps.Terrain[col][row] = card

	e.emit(GameEvent{
		Type:   "place_terrain",
		Player: -1,
		Data: map[string]any{
			"player":   playerID,
			"card":     cardToInfo(card),
			"position": pos,
			"elements": ps.Elements,
		},
	})

	// Trigger 入场 (on enter) effects for the terrain
	e.triggerEffects(TriggerOnEnter, card, nil, nil)
	e.notifyCardEntered(playerID, card, map[string]any{"entered_player": playerID, "terrain": true})

	e.checkWinCondition()
	return nil
}
