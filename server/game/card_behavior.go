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

type OnDamagedBehavior interface {
	OnDamaged(*EffectContext) error
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

type OnUseItemBehavior interface {
	OnUseItem(*EffectContext) error
}

type OnConsumeBehavior interface {
	OnConsume(*EffectContext) error
}

type SummonDevourRequirementBehavior interface {
	DevourRequirement() map[string]int
}

type RushBehavior interface {
	HasRush() bool
}

type PierceBehavior interface {
	HasPierce() bool
}

type TemporaryBehavior interface {
	IsTemporary() bool
}

type TauntBehavior interface {
	HasTaunt() bool
}

type StealthBehavior interface {
	StealthLayers() int
}

type ShieldBehavior interface {
	ShieldLayers() int
}

type ShieldingBehavior interface {
	HasShielding() bool
}

type GlobalSpellRangeBehavior interface {
	HasGlobalSpellRange() bool
}

type CooldownBehavior interface {
	Cooldown() int
}

type PerTurnLimitBehavior interface {
	PerTurnLimit() int
}

type OverloadBehavior interface {
	OverloadThreshold() int
}

type SpellTargetingBehavior interface {
	NeedsSpellTarget() bool
}

type SpellAreaBehavior interface {
	SpellArea() SpellArea
}

type SkillUsabilityBehavior interface {
	CanUseForSkillPurpose(skillPurpose) bool
}

type DefenseOnlySkillBehavior interface {
	IsDefenseOnlySkill() bool
}

type SorcerySkillBehavior interface {
	IsSorcerySkill() bool
}

type SpellHitStatusBehavior interface {
	SpellHitStatuses(*EffectContext) map[string]int
}

type SpellElementGainBehavior interface {
	SpellElementGains(*EffectContext) map[string]int
}

type PerTurnAbility interface {
	OnPerTurn(*EffectContext) error
}

type UltimateAbility interface {
	OnUltimate(*EffectContext) error
}

type SpellDamageBehavior interface {
	SpellDamage(*EffectContext) int
}

type SkillUseCostModifier interface {
	ModifySkillUseCost(*EffectContext, map[string]int)
}

type CardPlayCostModifier interface {
	ModifyCardPlayCost(*EffectContext, *CardInstance, map[string]int)
}

type SkillUsePermissionModifier interface {
	ValidateSkillUse(*EffectContext, *CardInstance, skillPurpose) error
}

type SpellStats struct {
	PowerBonus  int
	DamageBonus int
}

type SpellStatModifier interface {
	ModifySpellStats(*EffectContext, *SpellStats)
}

// SkillContributionModifier lets a concrete skill card alter its own
// contribution when it is used to attack, defend, or boost another spell.
type SkillContributionModifier interface {
	ModifySkillContribution(*EffectContext, *SpellStats)
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
	if h, ok := behavior.(OnDamagedBehavior); ok {
		r.Register(id, TriggerOnDamaged, h.OnDamaged)
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
	if h, ok := behavior.(OnUseItemBehavior); ok {
		r.Register(id, TriggerOnUseItem, h.OnUseItem)
	}
	if h, ok := behavior.(OnConsumeBehavior); ok {
		r.Register(id, TriggerOnConsume, h.OnConsume)
	}
	if h, ok := behavior.(PerTurnAbility); ok {
		r.RegisterActive(id, TriggerPerTurn, h.OnPerTurn)
	}
	if h, ok := behavior.(UltimateAbility); ok {
		r.RegisterActive(id, TriggerUltimate, h.OnUltimate)
	}
	if h, ok := behavior.(SpellDamageBehavior); ok {
		r.RegisterSpellDamage(id, h.SpellDamage)
	}
}

type noopPerTurn struct{}

func (noopPerTurn) OnPerTurn(*EffectContext) error { return nil }

type noopUltimate struct{}

func (noopUltimate) OnUltimate(*EffectContext) error { return nil }
