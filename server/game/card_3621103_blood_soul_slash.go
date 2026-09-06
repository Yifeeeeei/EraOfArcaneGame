package game

type Card3621103BloodSoulSlash struct{ AlwaysActive }

func (Card3621103BloodSoulSlash) ID() string { return "3621103" }

func (Card3621103BloodSoulSlash) Name() string { return "血魂斩" }

func (Card3621103BloodSoulSlash) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Source == nil || !isSpellBeingCast(ctx) {
		return nil
	}
	hero := ctx.Engine.State.Players[ctx.PlayerID].Hero
	if hero != nil {
		ctx.Engine.ApplyDamage(DamageRequest{Target: hero, Amount: 1, Kind: "blood_soul_slash", SourcePlayer: ctx.PlayerID, SourceKnown: true, Source: ctx.Source})
	}
	return nil
}

func (Card3621103BloodSoulSlash) OnSpellHit(ctx *EffectContext) error {
	if !isOwnSpellHit(ctx) {
		return nil
	}
	hero := ctx.Engine.State.Players[ctx.PlayerID].Hero
	if hero != nil {
		ctx.Engine.healUnit(hero, 2, ctx.Source)
	}
	return nil
}
