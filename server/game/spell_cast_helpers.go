package game

func isFriendlySpellCast(ctx *EffectContext) bool {
	event := ctx.SpellEvent()
	return !event.CasterKnown || ctx != nil && event.Caster == ctx.PlayerID
}
func isEnemySpellCast(ctx *EffectContext) bool {
	event := ctx.SpellEvent()
	return event.CasterKnown && ctx != nil && event.Caster != ctx.PlayerID
}
func isSpellBeingCast(ctx *EffectContext) bool         { return ctx != nil && ctx.Target == nil }
func spellUsePurpose(ctx *EffectContext) skillPurpose  { return ctx.SpellEvent().Purpose }
func spellCastSourceElement(ctx *EffectContext) string { return ctx.SpellEvent().Element }
func isFriendlySpellHit(ctx *EffectContext) bool       { return isFriendlySpellCast(ctx) }
func spellCasterFromData(ctx *EffectContext) (int, bool) {
	event := ctx.SpellEvent()
	return event.Caster, event.CasterKnown
}
func isOwnSpellHit(ctx *EffectContext) bool {
	return ctx != nil && ctx.SpellEvent().Spell != nil && ctx.SpellEvent().Spell == ctx.Source
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
		e.continueAfterPendingAction(afterDone)
		return true
	}

	if e.State.PendingAction != nil {
		e.continueAfterPendingAction(func() {
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
			e.continueAfterPendingAction(continueAfterFieldEffects)
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
		e.continueAfterPendingAction(func() {
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
