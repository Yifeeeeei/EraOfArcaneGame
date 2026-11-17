package game

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
	}
	return nil
}
