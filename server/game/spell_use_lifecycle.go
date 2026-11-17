package game

// Finishing a use may itself open a choice. Callers include it in a resolution
// frame, so the next spell/turn step cannot overtake that choice.
type SkillUseFinishedBehavior interface {
	HasActiveSkillUseFinished(*CardInstance) bool
	OnSkillUseFinished(*EffectContext)
}

func (AlwaysActive) HasActiveSkillUseFinished(*CardInstance) bool { return true }

// Cleanup is lifecycle bookkeeping, and still runs when effects are petrified.
type SkillUseCleanupBehavior interface{ ClearSkillUseState(*CardInstance) }

func (e *Engine) finishSkillUses(skills ...*CardInstance) {
	seen := map[*CardInstance]bool{}
	for _, skill := range skills {
		if skill == nil || seen[skill] {
			continue
		}
		seen[skill] = true
		if b, ok := cardBehavior(skill).(SkillUseCleanupBehavior); ok {
			b.ClearSkillUseState(skill)
		}
		if b, ok := cardBehavior(skill).(SkillUseFinishedBehavior); ok && b.HasActiveSkillUseFinished(skill) && !e.hasEffectiveStatus(skill, StatusPetrify) {
			b.OnSkillUseFinished(e.skillContext(skill.OwnerID, skill))
		}
		skill.usePayment = nil
	}
}
