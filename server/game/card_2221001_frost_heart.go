package game

type Card2221001FrostHeart struct{ AlwaysActive }

func (Card2221001FrostHeart) ID() string { return "2221001" }

func (Card2221001FrostHeart) Name() string { return "冰霜之心" }

func (Card2221001FrostHeart) CanReactToSpell(ctx *EffectContext, spell *SpellCast) bool {
	return ctx != nil && spell != nil && spell.AttackerID != ctx.PlayerID
}

func (Card2221001FrostHeart) OnSpellReaction(ctx *EffectContext, spell *SpellCast) error {
	targetInstanceID := ""
	if spell != nil && spell.Skill != nil {
		targetInstanceID = spell.Skill.InstanceID
	}
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModAllSpellDamageZero,
		TargetInstanceID: targetInstanceID,
		RemainingUses:    1,
		ExpiresTurn:      ctx.Engine.State.TurnNumber + 1,
	})
	ctx.Engine.discardFriendlyCandidate(ctx.PlayerID, ctx.Source.InstanceID)
	ctx.Engine.emit(GameEvent{Type: "spell_reaction", Player: -1, Data: map[string]any{
		"player": ctx.PlayerID,
		"card":   cardToInfo(ctx.Source),
		"effect": "damage_zero",
	}})
	return nil
}
