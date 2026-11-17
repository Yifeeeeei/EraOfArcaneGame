package game

// Explicit payment metadata shared by validation and the human action menu.
type AbilityTapCostBehavior interface{ AbilityRequiresTap(EffectTrigger) bool }

func abilityRequiresTap(card *CardInstance, trigger EffectTrigger) bool {
	if b, ok := cardBehavior(card).(AbilityTapCostBehavior); ok {
		return b.AbilityRequiresTap(trigger)
	}
	return false
}
