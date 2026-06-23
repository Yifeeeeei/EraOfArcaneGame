package game

func (e *Engine) applyGenericElementGain(playerID int, skill *CardInstance) {
	if skill == nil || skill.Card == nil {
		return
	}
	behavior, ok := globalRegistry.GetBehavior(skill.Card.Number).(SpellElementGainBehavior)
	if !ok || !behavior.HasActiveSpellElementGain(skill) {
		return
	}
	ctx := &EffectContext{Engine: e, Source: skill, PlayerID: playerID, OpponentID: 1 - playerID}
	gains := behavior.SpellElementGains(ctx)
	ps := e.State.Players[playerID]
	for elem, amount := range gains {
		if amount <= 0 {
			continue
		}
		ps.Elements[elem] += amount
		e.emit(GameEvent{
			Type:   "effect_trigger",
			Player: -1,
			Data: map[string]any{
				"source":  cardToInfo(skill),
				"effect":  "gain_element",
				"element": elem,
				"amount":  amount,
			},
		})
	}
}

func (e *Engine) applyExplicitSpellHitStatuses(skill *CardInstance, target *CardInstance) {
	if skill == nil || skill.Card == nil || target == nil {
		return
	}
	behavior, ok := globalRegistry.GetBehavior(skill.Card.Number).(SpellHitStatusBehavior)
	statuses := traitsForCardNumber(skill.Card.Number).statuses
	if ok && behavior.HasActiveSpellHitStatus(skill) {
		ctx := &EffectContext{Engine: e, Source: skill, Target: target, PlayerID: skill.OwnerID, OpponentID: 1 - skill.OwnerID}
		statuses = behavior.SpellHitStatuses(ctx)
	}
	for status, amount := range statuses {
		if amount <= 0 {
			continue
		}
		if !e.addStatus(target, status, amount) {
			continue
		}
		e.emit(GameEvent{
			Type:   "effect_trigger",
			Player: -1,
			Data: map[string]any{
				"source": cardToInfo(skill),
				"effect": "apply_status",
				"status": status,
				"amount": amount,
				"target": cardToInfo(target),
			},
		})
	}
}

func (e *Engine) genericSpellBonus(playerID int, skill *CardInstance, bonusKind string) int {
	return 0
}
