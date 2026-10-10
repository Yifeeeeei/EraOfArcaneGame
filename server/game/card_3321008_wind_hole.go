package game

type Card3321008WindHole struct{ AlwaysActive }

func (Card3321008WindHole) ID() string { return "3321008" }

func (Card3321008WindHole) Name() string { return "风洞" }

func (Card3321008WindHole) CanReactToSpell(ctx *EffectContext, spell *SpellCast) bool {
	return ctx != nil && spell != nil && spell.AttackerID != ctx.PlayerID && spellArea(spell.Skill) == SpellAreaSingle
}

func (Card3321008WindHole) OnSpellReaction(ctx *EffectContext, spell *SpellCast) error {
	ctx.Engine.cancelPendingSpell(ctx.PlayerID, ctx.Source, "wind_hole")
	return nil
}
