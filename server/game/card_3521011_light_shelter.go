package game

const lightShelterUntilStatus = "light_shelter_until"

type Card3521011LightShelter struct{ AlwaysActive }

func (Card3521011LightShelter) ID() string { return "3521011" }

func (Card3521011LightShelter) Name() string { return "光之庇护" }

func (Card3521011LightShelter) AllowsFriendlySpellTarget() bool {
	return true
}

func (Card3521011LightShelter) OnSpellHit(ctx *EffectContext) error {
	if !isOwnSpellHit(ctx) {
		return nil
	}
	target := ctx.Target
	if target != nil && target.Card.IsCompanion() {
		target.Statuses["防止致命"] = 1
		target.Statuses[lightShelterUntilStatus] = ctx.Engine.turnOrdinal() + 1
	}
	return nil
}

func (e *Engine) lightShelterActive(target *CardInstance) bool {
	return target != nil && target.Statuses["防止致命"] > 0 && target.Statuses[lightShelterUntilStatus] >= e.turnOrdinal()
}

func (e *Engine) expireTimedDamageProtection() {
	for _, player := range e.State.Players {
		for _, target := range e.getAllFieldCards(player) {
			if until := target.Statuses[lightShelterUntilStatus]; until > 0 && until <= e.turnOrdinal() {
				delete(target.Statuses, lightShelterUntilStatus)
				delete(target.Statuses, "防止致命")
			}
		}
	}
}
