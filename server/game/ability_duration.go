package game

const StatusAbilityDuration = "异能持续"

func abilityDurationActive(card *CardInstance) bool {
	return card != nil && card.Statuses[StatusAbilityDuration] > 0
}
