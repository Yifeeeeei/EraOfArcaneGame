package game

// dispatchPreparedSpell is shared by paid, free and virtual casts. Payment,
// readiness and target preparation have already succeeded. Every cast now
// uses the same counter -> defense/hit -> continuation protocol.
func (e *Engine) dispatchPreparedSpell(spell *SpellCast, data map[string]any, virtual bool) {
	playerID, skill, target := spell.AttackerID, spell.Skill, spell.Target
	if isSorcerySkill(skill.Card) {
		resolve := func() {
			e.runResolution("sorcery", func() {
				if e.shouldResolveSorceryHit(skill) {
					e.resolveSpellHit(playerID, skill, target, spell.BoostSkills, spell.ExtraTargets)
				}
			}, func() {
				if !virtual {
					e.removeStoredArchmageStaffSkillAfterUse(playerID, skill)
					e.notifySorceryResolved(playerID, skill)
				}
			})
		}
		if !e.triggerSpellCastFieldEffectsWithContinuation(playerID, skill, data, resolve) {
			resolve()
		}
		return
	}
	e.State.PendingSpell = spell
	continueSpell := func() {
		if e.State.PendingSpell == nil {
			return
		}
		if !e.spellAllowsDefense(playerID, skill, target) {
			e.resolvePendingSpellHit()
			return
		}
		e.State.ResumePhase, e.State.Phase = PhaseDefenseWindow, PhaseDefenseWindow
		e.emit(GameEvent{Type: "defense_window", Player: 1 - playerID, Data: map[string]any{"timeout": 30}})
	}
	resumeDefense := func() {
		if e.spellAllowsDefense(playerID, skill, target) {
			e.State.ResumePhase = PhaseDefenseWindow
		}
	}
	continueAfterCounters := func() {
		if e.promptAttackBoostSpellCastCounters(playerID, spell.BoostSkills, continueSpell) {
			resumeDefense()
			return
		}
		continueSpell()
	}
	if e.triggerSpellCastFieldEffectsWithContinuation(playerID, skill, data, continueAfterCounters) {
		resumeDefense()
		return
	}
	continueAfterCounters()
}
