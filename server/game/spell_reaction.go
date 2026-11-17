package game

import (
	"fmt"
)

func (e *Engine) handleReactSpell(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseDefenseWindow {
		return fmt.Errorf("not in spell reaction window")
	}
	if e.State.PendingSpell == nil {
		return fmt.Errorf("no pending spell")
	}
	if playerID == e.State.PendingSpell.AttackerID {
		return fmt.Errorf("attacker cannot react to their own spell this way")
	}
	if spellSuppressesOpponentResponses(e.State.PendingSpell.Skill) {
		return fmt.Errorf("opponent cannot react to this spell")
	}

	instanceID, _ := action.Data["instance_id"].(string)
	ps := e.State.Players[playerID]
	skill := e.findReactionCard(ps, instanceID)
	if skill == nil {
		return fmt.Errorf("reaction skill not found")
	}
	if err := e.validateSkillForPurpose(skill, skillPurposeReaction); err != nil {
		return err
	}
	cost := map[string]int{}
	if skill.Card.IsSkill() {
		cost = e.effectiveSkillUseCost(ps, skill)
	}
	overexertIDsRaw, _ := action.Data["overexert_ids"].([]any)
	overexertIDs := stringsFromAnySlice(overexertIDsRaw)
	overexertUnits, err := e.collectOverexertUnits(ps, overexertIDs)
	if err != nil {
		return err
	}
	if !e.canPayCostWithOverexertOptions(ps, cost, overexertUnits, e.playerHasLightWildcard(ps)) {
		return fmt.Errorf("not enough elements")
	}
	if !e.payDefenseCostWithOptions(ps, cost, action, overexertUnits, e.playerHasLightWildcard(ps)) {
		return fmt.Errorf("invalid payment")
	}
	e.destroyFuyeDoomedAfterExert(overexertUnits)

	if skill.Card.IsSkill() {
		skill.IsHorizontal = true
		if !e.shouldSkipCooldown(ps, skill) {
			e.ApplyKeywordOnSkillUse(skill)
		}
		e.applySkillUseCooldownModifiers(ps, skill)
	}
	if skill.Card.IsSkill() {
		e.consumeNextSkillUseModifiers(ps, skill)
	}
	e.advanceMasteryForUsedSkills(playerID, skill)

	behavior := behaviorForNumber(skill.Card.Number).(SpellReactionBehavior)
	if !behavior.HasActiveSpellReaction(skill) {
		return fmt.Errorf("skill cannot react to spells")
	}
	ctx := &EffectContext{
		Engine:     e,
		Source:     skill,
		PlayerID:   playerID,
		OpponentID: 1 - playerID,
		ExtraData: map[string]any{
			"react_player": playerID,
			"spell":        e.State.PendingSpell,
		},
	}
	return behavior.OnSpellReaction(ctx, e.State.PendingSpell)
}
