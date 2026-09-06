package game

// SpellAffectedUnitsBehavior overrides geometry when a spell selects units by
// a rule such as matching the primary target's element. Both resolution and
// the public affected-position projection use this same query.
type SpellAffectedUnitsBehavior interface {
	HasActiveAffectedSpellUnits(*CardInstance) bool
	AffectedSpellUnits(*EffectContext, int, SpellTarget) ([]*CardInstance, bool)
}

type SpellTargetExpansionBehavior interface {
	HasActiveSpellTargetExpansion(*CardInstance) bool
	ExpandSpellTargets(*EffectContext, SpellTarget, []SpellTarget) []SpellTarget
}

type DefenseSpellCancellationBehavior interface {
	CanCancelDefenseSpell(*CardInstance) bool
}

func (AlwaysActive) HasActiveAffectedSpellUnits(*CardInstance) bool   { return true }
func (AlwaysActive) HasActiveSpellTargetExpansion(*CardInstance) bool { return true }

func (e *Engine) customSpellAffectedUnits(defenderID int, skill *CardInstance, target SpellTarget) ([]*CardInstance, bool) {
	if b, ok := cardBehavior(skill).(SpellAffectedUnitsBehavior); ok && b.HasActiveAffectedSpellUnits(skill) {
		return b.AffectedSpellUnits(e.skillContext(skill.OwnerID, skill), defenderID, target)
	}
	return nil, false
}

func (e *Engine) expandSpellTargets(playerID int, target SpellTarget, extra []SpellTarget) []SpellTarget {
	for _, source := range e.getAllFieldCards(e.State.Players[playerID]) {
		if source == nil || e.hasEffectiveStatus(source, StatusPetrify) {
			continue
		}
		if b, ok := cardBehavior(source).(SpellTargetExpansionBehavior); ok && b.HasActiveSpellTargetExpansion(source) {
			extra = b.ExpandSpellTargets(e.skillContext(playerID, source), target, extra)
		}
	}
	return extra
}

// Target grants combine monotonically. Card code determines whether its own
// or another friendly spell receives the grant; the engine combines grants.
type SpellTargetGrant struct{ Pierce, IgnoreRange, AllowSameExtraTarget bool }
type SpellTargetGrantBehavior interface {
	HasActiveSpellTargetGrant(*CardInstance) bool
	SpellTargetGrant(*EffectContext, *CardInstance, SpellTarget) SpellTargetGrant
}

func (AlwaysActive) HasActiveSpellTargetGrant(*CardInstance) bool { return true }
func (e *Engine) spellTargetGrants(playerID int, skill *CardInstance, target SpellTarget) SpellTargetGrant {
	var result SpellTargetGrant
	if skill == nil || skill.Card == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return result
	}
	sources := append([]*CardInstance{skill}, e.getAllFieldCards(e.State.Players[playerID])...)
	seen := map[*CardInstance]bool{}
	for _, source := range sources {
		if source == nil || seen[source] || e.hasEffectiveStatus(source, StatusPetrify) {
			continue
		}
		seen[source] = true
		b, ok := cardBehavior(source).(SpellTargetGrantBehavior)
		if !ok || !b.HasActiveSpellTargetGrant(source) {
			continue
		}
		grant := b.SpellTargetGrant(e.skillContext(playerID, source), skill, target)
		result.Pierce = result.Pierce || grant.Pierce
		result.IgnoreRange = result.IgnoreRange || grant.IgnoreRange
		result.AllowSameExtraTarget = result.AllowSameExtraTarget || grant.AllowSameExtraTarget
	}
	return result
}
