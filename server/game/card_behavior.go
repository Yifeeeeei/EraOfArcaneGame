package game

// CardBehavior is the object-oriented behavior surface for a card definition.
// Each concrete card with custom rules implements this interface plus whichever
// trigger interfaces it needs.
type CardBehavior interface {
	ID() string
	Name() string
}

type OnEnterBehavior interface {
	OnEnter(*EffectContext) error
}

type OnDeathBehavior interface {
	OnDeath(*EffectContext) error
}

type OnTurnStartBehavior interface {
	OnTurnStart(*EffectContext) error
}

type OnUnitEnterBehavior interface {
	OnUnitEnter(*EffectContext) error
}

type OnFriendlyDeathBehavior interface {
	OnFriendlyDeath(*EffectContext) error
}

type OnEnemyDeathBehavior interface {
	OnEnemyDeath(*EffectContext) error
}

type OnSpellCastBehavior interface {
	OnSpellCast(*EffectContext) error
}

type OnSpellHitBehavior interface {
	OnSpellHit(*EffectContext) error
}

type OnDefendBehavior interface {
	OnDefend(*EffectContext) error
}

type OnEquipBehavior interface {
	OnEquip(*EffectContext) error
}

type PerTurnAbility interface {
	OnPerTurn(*EffectContext) error
}

type UltimateAbility interface {
	OnUltimate(*EffectContext) error
}

func registerBehavior(r *EffectRegistry, behavior CardBehavior) {
	id := behavior.ID()
	if h, ok := behavior.(OnEnterBehavior); ok {
		r.Register(id, TriggerOnEnter, h.OnEnter)
	}
	if h, ok := behavior.(OnDeathBehavior); ok {
		r.Register(id, TriggerOnDeath, h.OnDeath)
	}
	if h, ok := behavior.(OnTurnStartBehavior); ok {
		r.Register(id, TriggerOnTurnStart, h.OnTurnStart)
	}
	if h, ok := behavior.(OnUnitEnterBehavior); ok {
		r.Register(id, TriggerOnUnitEnter, h.OnUnitEnter)
	}
	if h, ok := behavior.(OnFriendlyDeathBehavior); ok {
		r.Register(id, TriggerOnFriendlyDeath, h.OnFriendlyDeath)
	}
	if h, ok := behavior.(OnEnemyDeathBehavior); ok {
		r.Register(id, TriggerOnEnemyDeath, h.OnEnemyDeath)
	}
	if h, ok := behavior.(OnSpellCastBehavior); ok {
		r.Register(id, TriggerOnSpellCast, h.OnSpellCast)
	}
	if h, ok := behavior.(OnSpellHitBehavior); ok {
		r.Register(id, TriggerOnSpellHit, h.OnSpellHit)
	}
	if h, ok := behavior.(OnDefendBehavior); ok {
		r.Register(id, TriggerOnDefend, h.OnDefend)
	}
	if h, ok := behavior.(OnEquipBehavior); ok {
		r.Register(id, TriggerOnEquip, h.OnEquip)
	}
	if h, ok := behavior.(PerTurnAbility); ok {
		r.RegisterActive(id, TriggerPerTurn, h.OnPerTurn)
	}
	if h, ok := behavior.(UltimateAbility); ok {
		r.RegisterActive(id, TriggerUltimate, h.OnUltimate)
	}
}

type noopPerTurn struct{}

func (noopPerTurn) OnPerTurn(*EffectContext) error { return nil }

type noopUltimate struct{}

func (noopUltimate) OnUltimate(*EffectContext) error { return nil }
