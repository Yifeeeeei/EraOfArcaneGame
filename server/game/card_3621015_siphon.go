package game

type Card3621015Siphon struct{ AlwaysActive }

func (Card3621015Siphon) ID() string { return "3621015" }

func (Card3621015Siphon) Name() string { return "虹吸" }

func (Card3621015Siphon) CanReactToSpell(ctx *EffectContext, spell *SpellCast) bool {
	return ctx != nil && spell != nil && spell.AttackerID != ctx.PlayerID && spell.Skill != nil
}

func (Card3621015Siphon) OnSpellReaction(ctx *EffectContext, spell *SpellCast) error {
	defenderID := ctx.PlayerID
	affected := ctx.Engine.spellAffectedUnitsWithExtraTargets(defenderID, spell.Skill, spell.Target, spell.ExtraTargets)
	if len(affected) == 0 {
		return nil
	}
	baseDamage := max(spell.Skill.Card.Attack+spell.Skill.AttackBonus, 0)
	if override, ok := globalRegistry.SpellDamage(spell.Skill.Card.Number, &EffectContext{
		Engine:     ctx.Engine,
		Source:     spell.Skill,
		Target:     affected[0],
		PlayerID:   spell.AttackerID,
		OpponentID: defenderID,
		ExtraData:  map[string]any{"target": spell.Target},
	}); ok {
		baseDamage = max(override, 0)
	}
	damage := ctx.Engine.effectiveSpellDamage(spell.AttackerID, spell.Skill, baseDamage, spell.BoostSkills)
	for _, unit := range affected {
		ctx.Engine.healUnit(unit, damage, ctx.Source)
	}
	ctx.Engine.cancelPendingSpell(ctx.PlayerID, ctx.Source, "siphon")
	return nil
}
