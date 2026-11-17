package game

// SpellCastOptions are prepared before any payment or state mutation. Optional
// hooks customize common casting operations without adding card-number branches
// to the cast handler. Commit and ResolveInstead must not fail validation.
type SpellCastOptions struct {
	Pierce           bool
	AllowExtraTarget bool
	ExtraTargets     func(hasPierce bool) ([]SpellTarget, error)
	ModifyCost       func(map[string]int)
	Commit           func()
	ResolveInstead   func()
}

type SpellPreparationBehavior interface {
	HasActiveSpellPreparation(*CardInstance) bool
	PrepareSpellCast(*EffectContext, SpellTarget, ActionMessage) (SpellCastOptions, error)
}

type OwnCostValidationBehavior interface {
	ValidateOwnCost(*EffectContext, map[string]int, ActionMessage) error
}

type SkillUseCommittedBehavior interface {
	HasActiveSkillUseCommitted(*CardInstance) bool
	OnSkillUseCommitted(*EffectContext)
}

type SorceryResolvedBehavior interface {
	HasActiveSorceryResolved(*CardInstance) bool
	OnSorceryResolved(*EffectContext)
}

type SorceryHitPolicyBehavior interface {
	HasActiveSorceryHitPolicy(*CardInstance) bool
	ResolvesSorceryHit(*CardInstance) bool
}

func (AlwaysActive) HasActiveSpellPreparation(*CardInstance) bool  { return true }
func (AlwaysActive) HasActiveSkillUseCommitted(*CardInstance) bool { return true }
func (AlwaysActive) HasActiveSorceryResolved(*CardInstance) bool   { return true }
func (AlwaysActive) HasActiveSorceryHitPolicy(*CardInstance) bool  { return true }

func (e *Engine) skillContext(playerID int, skill *CardInstance) *EffectContext {
	return &EffectContext{Engine: e, Source: skill, PlayerID: playerID, OpponentID: 1 - playerID}
}

func (e *Engine) prepareSpellOptions(playerID int, skill *CardInstance, target SpellTarget, action ActionMessage) (SpellCastOptions, error) {
	if b, ok := cardBehavior(skill).(SpellPreparationBehavior); ok && b.HasActiveSpellPreparation(skill) {
		return b.PrepareSpellCast(e.skillContext(playerID, skill), target, action)
	}
	return SpellCastOptions{}, nil
}

func (e *Engine) validateOwnCost(playerID int, skill *CardInstance, cost map[string]int, action ActionMessage) error {
	if b, ok := cardBehavior(skill).(OwnCostValidationBehavior); ok {
		return b.ValidateOwnCost(e.skillContext(playerID, skill), cost, action)
	}
	return nil
}

func (e *Engine) notifySkillUseCommitted(playerID int, skill *CardInstance) {
	if b, ok := cardBehavior(skill).(SkillUseCommittedBehavior); ok && b.HasActiveSkillUseCommitted(skill) {
		b.OnSkillUseCommitted(e.skillContext(playerID, skill))
	}
}

func (e *Engine) notifySorceryResolved(playerID int, skill *CardInstance) {
	if b, ok := cardBehavior(skill).(SorceryResolvedBehavior); ok && b.HasActiveSorceryResolved(skill) {
		b.OnSorceryResolved(e.skillContext(playerID, skill))
	}
}
