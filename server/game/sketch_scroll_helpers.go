package game

func (e *Engine) resolveSketchScrollSkill(playerID int, skillID string) {
	ps := e.State.Players[playerID]
	skill := e.findSkill(ps, skillID)
	if skill == nil || skill.Card == nil || !canUseSkillForPurpose(skill.Card, skillPurposeAttack) {
		return
	}
	targets := e.spellTargetCandidates(playerID, skill)
	if skillNeedsTargetInstance(skill) {
		if len(targets) == 0 {
			return
		}
		e.SetPendingAction(playerID, "sketch_scroll_target",
			"选择速写卷轴释放目标", targets, 1, 1,
			func(selected []string) {
				if len(selected) == 0 {
					return
				}
				target := selectedUnitFromCandidates(e, selected, targets)
				if target == nil || target.Position == nil {
					return
				}
				spellTarget := SpellTarget{Type: "unit", Position: *target.Position}
				if err := e.validateSpellTarget(playerID, skill, spellTarget); err != nil {
					return
				}
				e.castSkillFromSketchScroll(playerID, skill, spellTarget)
			})
		return
	}
	e.castSkillFromSketchScroll(playerID, skill, SpellTarget{Type: "none"})
}

func (e *Engine) spellTargetCandidates(playerID int, skill *CardInstance) []map[string]any {
	candidates := make([]map[string]any, 0)
	ps := e.State.Players[1-playerID]
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := ps.Units[col][row]
			if unit == nil || unit.Position == nil {
				continue
			}
			target := SpellTarget{Type: "unit", Position: *unit.Position}
			if err := e.validateSpellTarget(playerID, skill, target); err != nil {
				continue
			}
			candidates = append(candidates, candidateInfo(unit, "unit", "enemy"))
		}
	}
	return candidates
}

func (e *Engine) castSkillFromSketchScroll(playerID int, skill *CardInstance, target SpellTarget) {
	ps := e.State.Players[playerID]
	cost := e.effectiveSkillUseCost(ps, skill)
	if !e.payCostForAction(ps, cost, ActionMessage{}) {
		return
	}
	e.applySkillUseCooldownModifiers(ps, skill)
	totalPower := e.effectiveSpellPower(playerID, skill, nil, target)
	isSorcery := isSorcerySkill(skill.Card)
	spellCastData := map[string]any{
		"cast_player": playerID,
		"attacker":    playerID,
		"skill":       cardToInfo(skill),
		"target":      target,
		"power":       totalPower,
		"boost_count": 0,
		"is_sorcery":  isSorcery,
		"via":         "sketch_scroll",
	}
	e.emit(GameEvent{Type: "spell_cast", Player: -1, Data: spellCastData})
	e.triggerEffects(TriggerOnSpellCast, skill, nil, spellCastData)
	if isSorcery {
		resolveSorcery := func() {
			e.resolveSpellHit(playerID, skill, target, nil, nil)
		}
		if e.triggerSpellCastFieldEffectsWithContinuation(playerID, skill, spellCastData, resolveSorcery) {
			return
		}
		resolveSorcery()
		return
	}
	e.State.PendingSpell = &SpellCast{
		AttackerID:  playerID,
		Skill:       skill,
		Target:      target,
		TotalPower:  totalPower,
		BoostSkills: nil,
	}
	openDefenseWindow := func() {
		if e.State.PendingSpell == nil {
			return
		}
		e.State.ResumePhase = PhaseDefenseWindow
		e.State.Phase = PhaseDefenseWindow
		e.emit(GameEvent{Type: "defense_window", Player: 1 - playerID, Data: map[string]any{"timeout": 30}})
	}
	if e.triggerSpellCastFieldEffectsWithContinuation(playerID, skill, spellCastData, openDefenseWindow) {
		e.State.ResumePhase = PhaseDefenseWindow
		return
	}
	openDefenseWindow()
}
