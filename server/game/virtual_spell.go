package game

import (
	"fmt"
)

func (e *Engine) cloneVirtualSpell(source *CardInstance, ownerID int, turn int) *CardInstance {
	virtual := e.newCardInstance(source.Card, ownerID, turn)
	virtual.AttackBonus = source.AttackBonus
	virtual.PowerBonus = source.PowerBonus
	virtual.ElementsGainBonus = copyElements(source.ElementsGainBonus)
	virtual.Statuses = copyStatuses(source.Statuses)
	virtual.BoundSkills = append([]*CardInstance{}, source.BoundSkills...)
	return virtual
}

func copyElements(values map[string]int) map[string]int {
	if len(values) == 0 {
		return nil
	}
	copied := make(map[string]int, len(values))
	for key, value := range values {
		copied[key] = value
	}
	return copied
}

func copyStatuses(values map[string]int) map[string]int {
	if len(values) == 0 {
		return map[string]int{}
	}
	copied := make(map[string]int, len(values))
	for key, value := range values {
		copied[key] = value
	}
	return copied
}

func (e *Engine) startVirtualSpellCastNoBoost(playerID int, skill *CardInstance, target SpellTarget, extraData map[string]any) error {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) || skill == nil || skill.Card == nil {
		return fmt.Errorf("invalid virtual spell cast")
	}
	if !canUseSkillForPurpose(skill.Card, skillPurposeAttack) {
		return fmt.Errorf("virtual spell cannot attack")
	}
	if skillNeedsTargetInstance(skill) {
		if err := e.validateSpellTargetWithPierce(playerID, skill, target, e.skillHasPierce(playerID, skill)); err != nil {
			return err
		}
	}
	e.applyCoralBellyFirstSpellAttackBonus(playerID, skill)
	totalPower := e.effectiveSpellPower(playerID, skill, nil, target)
	powerSources := e.spellPowerSources(playerID, skill, nil, totalPower, target)
	e.consumeNextSpellPowerBonuses(e.State.Players[playerID], skill)

	isSorcery := isSorcerySkill(skill.Card)
	e.recordSpellCast(playerID, skill)
	e.triggerMagicMothAfterFocusSpellCast(playerID, skill)
	spellCastData := map[string]any{
		"cast_player":  playerID,
		"attacker":     playerID,
		"skill":        cardToInfo(skill),
		"target":       target,
		"power":        totalPower,
		"boost_count":  0,
		"is_sorcery":   isSorcery,
		"virtual_cast": true,
	}
	for key, value := range extraData {
		spellCastData[key] = value
	}
	e.emit(GameEvent{Type: "spell_cast", Player: -1, Data: spellCastData})
	e.triggerEffects(TriggerOnSpellCast, skill, nil, spellCastData)
	e.dispatchPreparedSpell(&SpellCast{AttackerID: playerID, Skill: skill, Target: target,
		TotalPower: totalPower, PowerSources: powerSources, BoostSkills: nil, ExtraTargets: nil}, spellCastData, true)
	return nil
}

func (e *Engine) cloneSpellInstances(skills []*CardInstance, ownerID int, turn int) []*CardInstance {
	clones := make([]*CardInstance, 0, len(skills))
	for _, skill := range skills {
		if skill == nil || skill.Card == nil {
			continue
		}
		clones = append(clones, e.cloneVirtualSpell(skill, ownerID, turn))
	}
	return clones
}

func (e *Engine) startVirtualSpellCastWithBoosts(playerID int, skill *CardInstance, target SpellTarget, boostSkills []*CardInstance, extraData map[string]any) error {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) || skill == nil || skill.Card == nil {
		return fmt.Errorf("invalid virtual boosted spell cast")
	}
	if !canUseSkillForPurpose(skill.Card, skillPurposeAttack) {
		return fmt.Errorf("virtual spell cannot attack")
	}
	if target.Type == "unit" && !target.Position.Valid() {
		return fmt.Errorf("invalid virtual spell target")
	}
	e.applyCoralBellyFirstSpellAttackBonus(playerID, skill)
	totalPower := e.effectiveSpellPower(playerID, skill, boostSkills, target)
	powerSources := e.spellPowerSources(playerID, skill, boostSkills, totalPower, target)
	isSorcery := isSorcerySkill(skill.Card)
	e.recordSpellCast(playerID, skill)
	e.triggerMagicMothAfterFocusSpellCast(playerID, skill)
	spellCastData := map[string]any{
		"cast_player":  playerID,
		"attacker":     playerID,
		"skill":        cardToInfo(skill),
		"target":       target,
		"power":        totalPower,
		"boost_count":  len(boostSkills),
		"boost_skills": boostSkills,
		"is_sorcery":   isSorcery,
		"virtual_cast": true,
	}
	for key, value := range extraData {
		spellCastData[key] = value
	}
	e.emit(GameEvent{Type: "spell_cast", Player: -1, Data: spellCastData})
	e.triggerEffects(TriggerOnSpellCast, skill, nil, spellCastData)
	e.dispatchPreparedSpell(&SpellCast{AttackerID: playerID, Skill: skill, Target: target,
		TotalPower: totalPower, PowerSources: powerSources, BoostSkills: boostSkills, ExtraTargets: nil}, spellCastData, true)
	return nil
}
