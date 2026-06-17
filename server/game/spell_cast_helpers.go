package game

func isFriendlySpellCast(ctx *EffectContext) bool {
	if ctx == nil || ctx.ExtraData == nil {
		return true
	}
	castPlayer, ok := ctx.ExtraData["cast_player"].(int)
	return !ok || castPlayer == ctx.PlayerID
}

func isEnemySpellCast(ctx *EffectContext) bool {
	if ctx == nil || ctx.ExtraData == nil {
		return false
	}
	castPlayer, ok := ctx.ExtraData["cast_player"].(int)
	return ok && castPlayer != ctx.PlayerID
}

func isSpellBeingCast(ctx *EffectContext) bool {
	return ctx != nil && ctx.Target == nil
}

func spellCastSourceElement(ctx *EffectContext) string {
	if ctx == nil {
		return ""
	}
	if ctx.Target != nil && ctx.Target.Card != nil {
		return ctx.Target.Card.Category
	}
	if ctx.ExtraData == nil {
		return ""
	}
	if skill, ok := ctx.ExtraData["skill"].(map[string]any); ok {
		if category, ok := skill["category"].(string); ok {
			return category
		}
	}
	return ""
}

func isFriendlySpellHit(ctx *EffectContext) bool {
	if ctx == nil || ctx.ExtraData == nil {
		return true
	}
	attacker, ok := ctx.ExtraData["attacker"].(int)
	return !ok || attacker == ctx.PlayerID
}

func (e *Engine) triggerSpellCastFieldEffectsWithContinuation(casterID int, source *CardInstance, extraData map[string]any, afterDone func()) bool {
	runFieldEffectsAndCounters := func() bool {
		fieldData := cloneExtraData(extraData)
		fieldData["skip_counter_traps"] = true

		e.triggerFieldEffectsWithData(TriggerOnSpellCast, casterID, source, fieldData)
		e.triggerFieldEffectsWithData(TriggerOnSpellCast, 1-casterID, source, fieldData)

		continueAfterFieldEffects := func() {
			counters := e.eligibleCounterTraps(casterID, TriggerOnSpellCast, source, extraData)
			counters = append(counters, e.eligibleCounterTraps(1-casterID, TriggerOnSpellCast, source, extraData)...)
			if e.promptCounterTrapQueue(counters, TriggerOnSpellCast, source, extraData, afterDone) {
				return
			}
			if afterDone != nil {
				afterDone()
			}
		}

		if e.State.PendingAction != nil {
			e.wrapPendingActionContinuation(continueAfterFieldEffects)
			return true
		}

		if e.promptCounterTrapQueue(
			append(
				e.eligibleCounterTraps(casterID, TriggerOnSpellCast, source, extraData),
				e.eligibleCounterTraps(1-casterID, TriggerOnSpellCast, source, extraData)...,
			),
			TriggerOnSpellCast,
			source,
			extraData,
			afterDone,
		) {
			return true
		}
		return false
	}

	if e.State.PendingAction != nil {
		e.wrapPendingActionContinuation(func() {
			if !runFieldEffectsAndCounters() && afterDone != nil {
				afterDone()
			}
		})
		return true
	}

	return runFieldEffectsAndCounters()
}

func cloneExtraData(data map[string]any) map[string]any {
	cloned := make(map[string]any, len(data)+1)
	for key, value := range data {
		cloned[key] = value
	}
	return cloned
}
