package game

func (e *Engine) resolveSketchScrollSkill(playerID int, skillID string) {
	ps := e.State.Players[playerID]
	skill := e.findSkill(ps, skillID)
	if skill == nil || skill.Card == nil || !canUseSkillForPurpose(skill.Card, skillPurposeAttack) {
		return
	}
	targets := e.enemyUnits(playerID, true, nil)
	if skillNeedsTargetInstance(skill) && len(targets) > 0 {
		e.SetPendingAction(playerID, "sketch_scroll_target",
			"选择速写卷轴释放目标", targets, 1, 1,
			func(selected []string) {
				if len(selected) == 0 {
					return
				}
				target := e.findUnitByID(1-playerID, selected[0])
				if target == nil || target.Position == nil {
					return
				}
				e.castSkillFromSketchScroll(playerID, skill, SpellTarget{Type: "unit", Position: *target.Position})
			})
		return
	}
	e.castSkillFromSketchScroll(playerID, skill, SpellTarget{Type: "none"})
}

func (e *Engine) castSkillFromSketchScroll(playerID int, skill *CardInstance, target SpellTarget) {
	ps := e.State.Players[playerID]
	cost := e.effectiveSkillUseCost(ps, skill)
	if !ps.PayCost(cost) {
		return
	}
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
	e.triggerFieldEffectsWithData(TriggerOnSpellCast, playerID, skill, spellCastData)
	e.triggerFieldEffectsWithData(TriggerOnSpellCast, 1-playerID, skill, spellCastData)
	if isSorcery {
		e.resolveSpellHit(playerID, skill, target, nil, nil)
		return
	}
	e.State.PendingSpell = &SpellCast{
		AttackerID:  playerID,
		Skill:       skill,
		Target:      target,
		TotalPower:  totalPower,
		BoostSkills: nil,
	}
	e.State.Phase = PhaseDefenseWindow
	e.emit(GameEvent{Type: "defense_window", Player: 1 - playerID, Data: map[string]any{"timeout": 30}})
}

func (e *Engine) findUnitByID(playerID int, instanceID string) *CardInstance {
	ps := e.State.Players[playerID]
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := ps.Units[col][row]
			if unit != nil && unit.InstanceID == instanceID {
				return unit
			}
		}
	}
	return nil
}
