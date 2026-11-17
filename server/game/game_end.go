package game

import (
	"fmt"
)

// checkWinCondition checks if the game is over
func (e *Engine) checkWinCondition() {
	if e.State.Phase == PhaseGameOver {
		return
	}
	if e.resolutionDepth > 0 || e.resolvingDeaths || len(e.deathQueue) > 0 {
		return
	}

	p0Dead := e.State.Players[0].Hero != nil && e.State.Players[0].Hero.CurrentLife <= 0
	p1Dead := e.State.Players[1].Hero != nil && e.State.Players[1].Hero.CurrentLife <= 0
	switch {
	case p0Dead && p1Dead:
		e.State.Winner = -2
		e.State.Phase = PhaseGameOver
		e.clearPendingForGameOver()
		e.emit(GameEvent{
			Type:   "game_over",
			Player: -1,
			Data: map[string]any{
				"winner": e.State.Winner,
				"reason": "both_heroes_killed",
			},
		})
	case p0Dead:
		e.State.Winner = 1
		e.State.Phase = PhaseGameOver
		e.clearPendingForGameOver()
		e.emit(GameEvent{
			Type:   "game_over",
			Player: -1,
			Data: map[string]any{
				"winner": e.State.Winner,
				"reason": "hero_killed",
			},
		})
	case p1Dead:
		e.State.Winner = 0
		e.State.Phase = PhaseGameOver
		e.clearPendingForGameOver()
		e.emit(GameEvent{
			Type:   "game_over",
			Player: -1,
			Data: map[string]any{
				"winner": e.State.Winner,
				"reason": "hero_killed",
			},
		})
	}
}

func (e *Engine) finishGame(winner int, reason string, actor int) {
	if e.State.Phase == PhaseGameOver {
		return
	}
	e.State.Winner = winner
	e.State.Phase = PhaseGameOver
	e.clearPendingForGameOver()
	e.emit(GameEvent{
		Type:   "game_over",
		Player: -1,
		Data: map[string]any{
			"winner": e.State.Winner,
			"reason": reason,
			"actor":  actor,
		},
	})
}

func (e *Engine) clearPendingForGameOver() {
	e.State.PendingAction = nil
	e.State.PendingActionQueue = nil
	e.State.PendingSpell = nil
	e.State.DrawOfferBy = -1
	e.State.ResumePhase = PhaseGameOver
}

func (e *Engine) handleSurrender(playerID int) error {
	if err := e.ensureMatchActionPlayer(playerID); err != nil {
		return err
	}
	e.finishGame(1-playerID, "surrender", playerID)
	return nil
}

func (e *Engine) handleOfferDraw(playerID int) error {
	if err := e.ensureMatchActionPlayer(playerID); err != nil {
		return err
	}
	if e.State.DrawOfferBy == 1-playerID {
		e.finishGame(-2, "draw_agreement", playerID)
		return nil
	}
	if e.State.DrawOfferBy == playerID {
		return fmt.Errorf("draw offer already pending")
	}
	e.State.DrawOfferBy = playerID
	e.emit(GameEvent{
		Type:   "draw_offer",
		Player: -1,
		Data: map[string]any{
			"player": playerID,
		},
	})
	return nil
}

func (e *Engine) handleRespondDrawOffer(playerID int, action ActionMessage) error {
	if err := e.ensureMatchActionPlayer(playerID); err != nil {
		return err
	}
	if e.State.DrawOfferBy != 1-playerID {
		return fmt.Errorf("no draw offer to respond to")
	}
	accept, _ := action.Data["accept"].(bool)
	offerBy := e.State.DrawOfferBy
	if accept {
		e.finishGame(-2, "draw_agreement", playerID)
		return nil
	}
	e.State.DrawOfferBy = -1
	e.emit(GameEvent{
		Type:   "draw_offer_declined",
		Player: -1,
		Data: map[string]any{
			"player":   playerID,
			"offer_by": offerBy,
		},
	})
	return nil
}

func (e *Engine) ensureMatchActionPlayer(playerID int) error {
	if playerID < 0 || playerID >= len(e.State.Players) || e.State.Players[playerID] == nil {
		return fmt.Errorf("invalid player")
	}
	if e.State.Phase == PhaseWaitingPlayers {
		return fmt.Errorf("game has not started")
	}
	if e.State.Phase == PhaseGameOver {
		return fmt.Errorf("game is already over")
	}
	return nil
}
