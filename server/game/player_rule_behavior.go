package game

// These interfaces describe rule changes at explicit engine boundaries. Their
// instance predicates are evaluated at use time, not inferred from card IDs.
type BeforeDrawBehavior interface {
	HasActiveBeforeDraw(*CardInstance) bool
	OnBeforeDraw(*EffectContext) error
}

type OpeningHandBehavior interface {
	HasActiveOpeningHand(*CardInstance) bool
	OpeningHandBonus(*CardInstance) int
}

type EmptyDeckDrawBehavior interface {
	HasActiveEmptyDeckDraw(*CardInstance) bool
	DrawFromEmptyDeck(*EffectContext, int) []*CardInstance
}

type HandLimitEnforcementBehavior interface {
	EnforcesHandLimit(*CardInstance, int) bool
}

type PaymentRules struct {
	LightPaysAny bool
	AnyPaysLight bool
}

type PaymentRulesBehavior interface {
	HasActivePaymentRules(*CardInstance) bool
	PaymentRules(*CardInstance) PaymentRules
}

type DefenseRequirementBehavior interface {
	HasActiveDefenseRequirement(*CardInstance) bool
	RequiredDefensePower(*CardInstance, int) int
}

func (AlwaysActive) HasActiveBeforeDraw(*CardInstance) bool         { return true }
func (AlwaysActive) HasActiveOpeningHand(*CardInstance) bool        { return true }
func (AlwaysActive) HasActiveEmptyDeckDraw(*CardInstance) bool      { return true }
func (AlwaysActive) HasActivePaymentRules(*CardInstance) bool       { return true }
func (AlwaysActive) HasActiveDefenseRequirement(*CardInstance) bool { return true }

func (e *Engine) playerPaymentRules(ps *PlayerState) PaymentRules {
	var rules PaymentRules
	if e == nil || ps == nil {
		return rules
	}
	for _, source := range e.getAllFieldCards(ps) {
		if source == nil || e.hasEffectiveStatus(source, StatusPetrify) {
			continue
		}
		behavior, ok := cardBehavior(source).(PaymentRulesBehavior)
		if !ok || !behavior.HasActivePaymentRules(source) {
			continue
		}
		grant := behavior.PaymentRules(source)
		rules.LightPaysAny = rules.LightPaysAny || grant.LightPaysAny
		rules.AnyPaysLight = rules.AnyPaysLight || grant.AnyPaysLight
	}
	return rules
}
