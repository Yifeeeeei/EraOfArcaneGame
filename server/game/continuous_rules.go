package game

// Continuous grants are pure queries. Nonempty groups explicitly prevent
// stacking copies of the same rule; different groups compose in field order.
type HandLimitGrant struct {
	Group string
	Delta int
}
type HandLimitBehavior interface {
	HasActiveHandLimit(*CardInstance) bool
	HandLimitGrant(*CardInstance, int) HandLimitGrant
}
type SpellPowerScale struct {
	Group                  string
	Numerator, Denominator int
}
type SpellPowerScaleBehavior interface {
	HasActiveSpellPowerScale(*CardInstance) bool
	SpellPowerScale(*CardInstance, *CardInstance) SpellPowerScale
}

func (AlwaysActive) HasActiveHandLimit(*CardInstance) bool       { return true }
func (AlwaysActive) HasActiveSpellPowerScale(*CardInstance) bool { return true }

func (e *Engine) handLimitForPlayer(ps *PlayerState) int {
	limit := e.State.HandLimit
	if ps == nil {
		return limit
	}
	seen := map[string]bool{}
	for _, owner := range e.State.Players {
		if owner == nil {
			continue
		}
		for _, source := range e.getAllFieldCards(owner) {
			if source == nil || e.hasEffectiveStatus(source, StatusPetrify) {
				continue
			}
			b, ok := cardBehavior(source).(HandLimitBehavior)
			if !ok || !b.HasActiveHandLimit(source) {
				continue
			}
			grant := b.HandLimitGrant(source, ps.PlayerID)
			if grant.Group != "" {
				if seen[grant.Group] {
					continue
				}
				seen[grant.Group] = true
			}
			limit += grant.Delta
		}
	}
	return max(0, limit)
}

func (e *Engine) scaleSpellPower(playerID int, skill *CardInstance, power int) int {
	seen := map[string]bool{}
	for _, source := range e.getAllFieldCards(e.State.Players[playerID]) {
		if source == nil || e.hasEffectiveStatus(source, StatusPetrify) {
			continue
		}
		b, ok := cardBehavior(source).(SpellPowerScaleBehavior)
		if !ok || !b.HasActiveSpellPowerScale(source) {
			continue
		}
		scale := b.SpellPowerScale(source, skill)
		if scale.Denominator <= 0 {
			continue
		}
		if scale.Group != "" {
			if seen[scale.Group] {
				continue
			}
			seen[scale.Group] = true
		}
		power = (power*scale.Numerator + scale.Denominator - 1) / scale.Denominator
	}
	return max(0, power)
}
