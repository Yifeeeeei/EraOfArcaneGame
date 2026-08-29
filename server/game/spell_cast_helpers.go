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

func spellUsePurpose(ctx *EffectContext) skillPurpose {
	if ctx != nil && ctx.ExtraData != nil {
		if purpose, ok := ctx.ExtraData["purpose"].(string); ok && purpose != "" {
			return skillPurpose(purpose)
		}
	}
	return skillPurposeAttack
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

func spellCasterFromData(ctx *EffectContext) (int, bool) {
	if ctx == nil || ctx.ExtraData == nil {
		return 0, false
	}
	if attacker, ok := ctx.ExtraData["attacker"].(int); ok {
		return attacker, true
	}
	if castPlayer, ok := ctx.ExtraData["cast_player"].(int); ok {
		return castPlayer, true
	}
	return 0, false
}

func isOwnSpellHit(ctx *EffectContext) bool {
	return ctx != nil && ctx.ExtraData != nil && ctx.ExtraData["spell_source"] == ctx.Source
}

func (e *Engine) triggerSpellUseFieldEffectsWithContinuation(casterID int, source *CardInstance, extraData map[string]any, afterDone func()) bool {
	runFieldEffects := func() bool {
		fieldData := cloneExtraData(extraData)
		fieldData["skip_counter_traps"] = true
		e.triggerFieldEffectsWithData(TriggerOnSpellCast, casterID, source, fieldData)
		if !spellSuppressesOpponentResponses(source) {
			e.triggerFieldEffectsWithData(TriggerOnSpellCast, 1-casterID, source, fieldData)
		}
		if e.State.PendingAction == nil {
			return false
		}
		e.wrapPendingActionContinuation(afterDone)
		return true
	}

	if e.State.PendingAction != nil {
		e.wrapPendingActionContinuation(func() {
			if !runFieldEffects() && afterDone != nil {
				afterDone()
			}
		})
		return true
	}
	return runFieldEffects()
}

func (e *Engine) triggerSpellCastFieldEffectsWithContinuation(casterID int, source *CardInstance, extraData map[string]any, afterDone func()) bool {
	runFieldEffectsAndCounters := func() bool {
		suppressesOpponentResponses := spellSuppressesOpponentResponses(source)
		fieldData := cloneExtraData(extraData)
		fieldData["skip_counter_traps"] = true

		e.triggerFieldEffectsWithData(TriggerOnSpellCast, casterID, source, fieldData)
		if !suppressesOpponentResponses {
			e.triggerFieldEffectsWithData(TriggerOnSpellCast, 1-casterID, source, fieldData)
		}

		continueAfterFieldEffects := func() {
			counters := e.eligibleCounterTraps(casterID, TriggerOnSpellCast, source, extraData)
			if !suppressesOpponentResponses {
				counters = append(counters, e.eligibleCounterTraps(1-casterID, TriggerOnSpellCast, source, extraData)...)
			}
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
			spellCastCounterCandidates(e, casterID, source, extraData, suppressesOpponentResponses),
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

func spellCastCounterCandidates(e *Engine, casterID int, source *CardInstance, extraData map[string]any, suppressOpponent bool) []*CardInstance {
	counters := e.eligibleCounterTraps(casterID, TriggerOnSpellCast, source, extraData)
	if !suppressOpponent {
		counters = append(counters, e.eligibleCounterTraps(1-casterID, TriggerOnSpellCast, source, extraData)...)
	}
	return counters
}

func cloneExtraData(data map[string]any) map[string]any {
	cloned := make(map[string]any, len(data)+1)
	for key, value := range data {
		cloned[key] = value
	}
	return cloned
}
