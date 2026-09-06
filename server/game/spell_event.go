package game

// SpellEvent separates the spell being observed from the observing card.
// The registry's map-shaped trigger envelope remains private plumbing; card
// predicates use this projection instead of guessing ownership from Target.
type SpellEvent struct {
	Spell            *CardInstance
	Caster           int
	CasterKnown      bool
	Purpose          skillPurpose
	Power, Damage    int
	BoostCount       int
	AffectedUnits    []*CardInstance
	DefenseSucceeded bool
	Timing           string
	Element          string
}

func (ctx *EffectContext) SpellEvent() SpellEvent {
	event := SpellEvent{Caster: -1, Purpose: skillPurposeAttack}
	if ctx == nil {
		return event
	}
	data := ctx.ExtraData
	event.Spell, _ = data["spell_source"].(*CardInstance)
	event.Caster, event.CasterKnown = data["attacker"].(int)
	if !event.CasterKnown {
		event.Caster, event.CasterKnown = data["cast_player"].(int)
	}
	if !event.CasterKnown {
		event.Caster = -1
	}
	if purpose, ok := data["purpose"].(string); ok && purpose != "" {
		event.Purpose = skillPurpose(purpose)
	}
	event.Power = spellPowerFromData(data)
	event.Damage = damageFromData(data)
	event.BoostCount = intFromData(data, "boost_count", -1)
	event.DefenseSucceeded = boolFromData(data, "defense_success")
	event.Timing, _ = data["timing"].(string)
	event.AffectedUnits, _ = data["affected_units"].([]*CardInstance)
	if ctx.Target != nil && ctx.Target.Card != nil {
		event.Element = ctx.Target.Card.Category
	} else if skill, ok := data["skill"].(map[string]any); ok {
		event.Element, _ = skill["category"].(string)
	}
	return event
}
