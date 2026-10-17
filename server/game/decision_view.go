package game

// Public turn ownership and decision ownership are separate. No queued frame,
// counter identity, candidate, or private prompt is exposed by this projection.
func (e *Engine) decisionPlayer() int {
	if e.State.Phase == PhaseGameOver || e.State.Phase == PhaseMulligan {
		return -1
	}
	if p := e.State.PendingAction; p != nil {
		return p.PlayerID
	}
	if e.State.Phase == PhaseDefenseWindow && e.State.PendingSpell != nil {
		s := e.State.PendingSpell
		return e.spellDefenderID(s.AttackerID, s.Skill, s.Target)
	}
	return e.State.CurrentTurn
}

// Explicit projection of legacy choice contracts. The frontend renders these
// primitives; it never matches card numbers, names or prompt text.
func pendingActionUI(p *PendingAction) map[string]any {
	ui := map[string]any{"kind": "selection"}
	if p.AllocationTotal > 0 {
		return map[string]any{"kind": "allocation", "total": p.AllocationTotal}
	}
	switch p.Type {
	case "water_divination_search", "okoru_observe", "foresight_orb_reorder", "council_guard_reorder", "rose_prophet_reorder":
		ui["kind"] = "reorder"
	case "natural_communion_distribute":
		ui["kind"] = "load_distribution"
	case "blood_sand_array_pay", "blood_sand_array_pay_opponent":
		ui["kind"] = "resource_distribution"
		ui["max_total"] = 3
		ui["modes"] = []map[string]string{{"value": "auto", "label": "先负载，后生命"}, {"value": "life", "label": "仅生命"}}
	case "muling_return":
		ui["unique_by"] = "side"
	}
	return ui
}

// Existing legacy preparation requirements are projected at one boundary.
// New custom preparations should extend this contract, never inspect UI text.
func cardPreparationUI(card *CardInstance) string {
	switch card.Card.Number {
	case "2121105":
		return "fire_sacrifice"
	case "2521111":
		return "companion_support"
	case "3121104":
		return "optional_fire_sacrifice"
	}
	return ""
}

// Allocation keeps distinct candidate selection separate from integer amounts.
// Older clients may omit amounts and retain the original first-target remainder.
func (e *Engine) SetAllocatedPendingAction(player int, kind, prompt string, candidates []map[string]any, minSelect, maxSelect, total int, callback func([]string)) {
	action := e.setPendingActionWithOptions(player, kind, prompt, candidates, minSelect, maxSelect, nil, false, callback, nil, nil, nil, nil)
	if action != nil {
		action.AllocationTotal = total
	}
}
