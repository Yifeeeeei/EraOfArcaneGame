package game

import (
	"fmt"
)

// handleResolveAction handles the player's response to a pending action
func (e *Engine) handleResolveAction(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseWaitingAction {
		return fmt.Errorf("no pending action")
	}
	pa := e.State.PendingAction
	if pa == nil {
		return fmt.Errorf("no pending action")
	}
	if pa.PlayerID != playerID {
		return fmt.Errorf("not your pending action")
	}

	selectedRaw, _ := action.Data["selected"].([]any)
	var selected []string
	for _, s := range selectedRaw {
		if str, ok := s.(string); ok {
			selected = append(selected, str)
		}
	}

	if len(selected) < pa.MinSelect {
		return fmt.Errorf("must select at least %d", pa.MinSelect)
	}
	if len(selected) > pa.MaxSelect {
		return fmt.Errorf("can select at most %d", pa.MaxSelect)
	}
	allowed := make(map[string]bool, len(pa.Candidates))
	selectable := make(map[string]bool, len(pa.Candidates))
	for _, candidate := range pa.Candidates {
		if id, ok := candidate["instance_id"].(string); ok && id != "" {
			allowed[id] = true
			selectable[id] = candidate["can_select"] != false
		}
	}
	seen := make(map[string]bool, len(selected))
	for _, id := range selected {
		if !allowed[id] {
			return fmt.Errorf("invalid selection")
		}
		if !selectable[id] {
			return fmt.Errorf("invalid selection")
		}
		if seen[id] {
			return fmt.Errorf("duplicate selection")
		}
		seen[id] = true
	}

	// Execute callback
	callback := pa.Callback
	callbackData := pa.CallbackData
	callbackErr := pa.CallbackErr
	data := action.Data
	e.State.PendingAction = nil
	e.State.Phase = e.State.ResumePhase

	if callbackErr != nil {
		if err := callbackErr(selected, data); err != nil {
			e.State.PendingAction = pa
			e.State.Phase = PhaseWaitingAction
			return err
		}
	} else if callbackData != nil {
		callbackData(selected, data)
	} else if callback != nil {
		callback(selected)
	}

	e.emitPendingActionCleared(pa)
	e.advancePendingActionQueue()
	e.completeActionResolutions(pa)
	if e.State.PendingAction == nil && e.State.Phase == PhaseDefenseWindow && e.State.PendingSpell == nil {
		e.State.Phase = PhaseMain
	}

	e.checkWinCondition()
	return nil
}

// SetPendingAction sets a pending player action and pauses the game
func (e *Engine) SetPendingAction(playerID int, actionType string, prompt string, candidates []map[string]any, minSelect, maxSelect int, callback func([]string)) {
	e.setPendingAction(playerID, actionType, prompt, candidates, minSelect, maxSelect, callback, nil)
}

func (e *Engine) SetPendingActionWithData(playerID int, actionType string, prompt string, candidates []map[string]any, minSelect, maxSelect int, callback func([]string, map[string]any)) {
	e.setPendingAction(playerID, actionType, prompt, candidates, minSelect, maxSelect, nil, callback)
}

func (e *Engine) SetPendingActionWithError(playerID int, actionType string, prompt string, candidates []map[string]any, minSelect, maxSelect int, cost map[string]int, canOverexert bool, callback func([]string, map[string]any) error) {
	e.setPendingActionWithOptions(playerID, actionType, prompt, candidates, minSelect, maxSelect, cost, canOverexert, nil, nil, callback, nil, nil)
}

func (e *Engine) SetPendingActionWithErrorAndContext(playerID int, actionType string, prompt string, candidates []map[string]any, minSelect, maxSelect int, cost map[string]int, canOverexert bool, context map[string]any, callback func([]string, map[string]any) error) {
	e.setPendingActionWithOptions(playerID, actionType, prompt, candidates, minSelect, maxSelect, cost, canOverexert, nil, nil, callback, context, nil)
}

func (e *Engine) setPendingAction(playerID int, actionType string, prompt string, candidates []map[string]any, minSelect, maxSelect int, callback func([]string), callbackData func([]string, map[string]any)) {
	e.setPendingActionWithOptions(playerID, actionType, prompt, candidates, minSelect, maxSelect, nil, false, callback, callbackData, nil, nil, nil)
}

func (e *Engine) setPendingActionWithOptions(playerID int, actionType string, prompt string, candidates []map[string]any, minSelect, maxSelect int, cost map[string]int, canOverexert bool, callback func([]string), callbackData func([]string, map[string]any), callbackErr func([]string, map[string]any) error, context map[string]any, available func() bool) *PendingAction {
	if minSelect > 0 && len(candidates) == 0 {
		return nil
	}
	if available != nil && !available() {
		return nil
	}
	resumePhase := e.State.Phase
	if e.State.PendingAction != nil {
		resumePhase = e.State.ResumePhase
	}
	action := &PendingAction{
		Type:         actionType,
		PlayerID:     playerID,
		Prompt:       prompt,
		Candidates:   candidates,
		MinSelect:    minSelect,
		MaxSelect:    maxSelect,
		Context:      context,
		Cost:         cost,
		CanOverexert: canOverexert,
		Callback:     callback,
		CallbackData: callbackData,
		CallbackErr:  callbackErr,
		Available:    available,
	}
	if e.State.PendingAction != nil {
		e.State.PendingActionQueue = append(e.State.PendingActionQueue, action)
		return action
	}
	e.activatePendingAction(action, resumePhase)
	return action
}

func (e *Engine) activatePendingAction(action *PendingAction, resumePhase GamePhase) {
	if action == nil {
		return
	}
	e.State.ResumePhase = resumePhase
	e.State.Phase = PhaseWaitingAction
	e.State.PendingAction = action
	data := map[string]any{
		"type":       action.Type,
		"player_id":  action.PlayerID,
		"prompt":     action.Prompt,
		"candidates": action.Candidates,
		"min_select": action.MinSelect,
		"max_select": action.MaxSelect,
	}
	if action.Context != nil {
		data["context"] = action.Context
	}
	if action.Cost != nil {
		data["cost"] = action.Cost
	}
	if action.CanOverexert {
		data["can_overexert"] = true
	}
	e.emit(GameEvent{Type: "pending_action", Player: action.PlayerID, Data: data})
}

func (e *Engine) advancePendingActionQueue() bool {
	if e.State.PendingAction != nil || len(e.State.PendingActionQueue) == 0 {
		return false
	}
	var skipped []*PendingAction
	defer func() {
		// First expose the next valid sibling. Parents of skipped choices
		// must wait for it just like parents of successfully resolved ones.
		for _, action := range skipped {
			e.completeActionResolutions(action)
		}
	}()
	for len(e.State.PendingActionQueue) > 0 {
		next := e.State.PendingActionQueue[0]
		e.State.PendingActionQueue = e.State.PendingActionQueue[1:]
		if next.Available != nil && !next.Available() {
			skipped = append(skipped, next)
			continue
		}
		e.activatePendingAction(next, e.State.Phase)
		return true
	}
	return false
}

func (e *Engine) emitPendingActionCleared(action *PendingAction) {
	if action == nil {
		return
	}
	e.emit(GameEvent{Type: "pending_action_cleared", Player: action.PlayerID, Data: map[string]any{
		"type":      action.Type,
		"player_id": action.PlayerID,
	}})
}
