package game

import "fmt"

func (e *Engine) sketchScrollSkillCandidates(playerID int) []map[string]any {
	ps := e.State.Players[playerID]
	return e.friendlySkills(playerID, func(skill *CardInstance) bool {
		cost := e.effectiveSkillUseCost(ps, skill)
		return canUseSkillForPurpose(skill.Card, skillPurposeAttack) &&
			!isSorcerySkill(skill.Card) &&
			e.canPayCostForCardAction(ps, skill, cost, cost, paymentPurposeUse, ActionMessage{})
	})
}

func (e *Engine) resolveSketchScrollSkill(playerID int, skillID string) error {
	ps := e.State.Players[playerID]
	skill := e.findSkill(ps, skillID)
	if skill == nil || skill.Card == nil || !canUseSkillForPurpose(skill.Card, skillPurposeAttack) {
		return fmt.Errorf("invalid sketch scroll skill")
	}
	if isSorcerySkill(skill.Card) {
		return fmt.Errorf("sketch scroll cannot copy sorceries")
	}
	cost := e.effectiveSkillUseCost(ps, skill)
	if !e.canPayCostForCardAction(ps, skill, cost, cost, paymentPurposeUse, ActionMessage{}) {
		return fmt.Errorf("not enough elements")
	}
	targets := e.spellTargetCandidates(playerID, skill)
	if skillNeedsTargetInstance(skill) {
		if len(targets) == 0 {
			return nil
		}
		e.SetPendingActionWithError(playerID, "sketch_scroll_target",
			"选择速写卷轴释放目标", targets, 1, 1,
			nil, false, func(selected []string, _ map[string]any) error {
				if len(selected) == 0 {
					return nil
				}
				target := selectedUnitFromCandidates(e, selected, targets)
				if target == nil || target.Position == nil {
					return fmt.Errorf("invalid sketch scroll target")
				}
				if target.Card != nil && target.Card.IsHero() {
					return fmt.Errorf("sketch scroll cannot target heroes as units")
				}
				spellTarget := SpellTarget{Type: "unit", Position: *target.Position}
				if err := e.validateSpellTarget(playerID, skill, spellTarget); err != nil {
					return err
				}
				return e.castSkillFromSketchScroll(playerID, skill, spellTarget)
			})
		return nil
	}
	return e.castSkillFromSketchScroll(playerID, skill, SpellTarget{Type: "none"})
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
			if unit.Card != nil && unit.Card.IsHero() {
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

func (e *Engine) castSkillFromSketchScroll(playerID int, skill *CardInstance, target SpellTarget) error {
	ps := e.State.Players[playerID]
	cost := e.effectiveSkillUseCost(ps, skill)
	if !e.payCostForCardAction(ps, skill, cost, cost, paymentPurposeUse, ActionMessage{}) {
		return fmt.Errorf("not enough elements")
	}
	e.applySkillUseCooldownModifiers(ps, skill)
	e.advanceMasteryForUsedSkills(playerID, skill)
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
	e.recordSpellCast(playerID, skill)
	e.emit(GameEvent{Type: "spell_cast", Player: -1, Data: spellCastData})
	e.triggerEffects(TriggerOnSpellCast, skill, nil, spellCastData)
	if isSorcery {
		resolveSorcery := func() {
			e.resolveSpellHit(playerID, skill, target, nil, nil)
		}
		if e.triggerSpellCastFieldEffectsWithContinuation(playerID, skill, spellCastData, resolveSorcery) {
			return nil
		}
		resolveSorcery()
		return nil
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
		return nil
	}
	openDefenseWindow()
	return nil
}
