package game

// CardBehavior is the object-oriented behavior surface for a card definition.
// Each concrete card with custom rules implements this interface plus whichever
// trigger interfaces it needs.
type CardBehavior interface {
	ID() string
	Name() string
}

type PerTurnLabelBehavior interface {
	PerTurnLabel(*CardInstance) string
}

// AlwaysActive is embedded by normal card behavior structs. It provides the
// default answer for instance-aware ability predicates: this card has the
// behavior and that behavior is currently usable/effective.
//
// Cards whose own text or state can switch an ability off should override the
// matching HasActive... method on the concrete card type.
type AlwaysActive struct{}

func (AlwaysActive) HasActiveOnEnter(*CardInstance) bool                   { return true }
func (AlwaysActive) HasActiveGameStart(*CardInstance) bool                 { return true }
func (AlwaysActive) HasActiveDeathrattle(*CardInstance) bool               { return true }
func (AlwaysActive) HasActiveTurnStart(*CardInstance) bool                 { return true }
func (AlwaysActive) HasActiveTurnEnd(*CardInstance) bool                   { return true }
func (AlwaysActive) HasActiveUnitEnter(*CardInstance) bool                 { return true }
func (AlwaysActive) HasActiveFriendlyDeath(*CardInstance) bool             { return true }
func (AlwaysActive) HasActiveEnemyDeath(*CardInstance) bool                { return true }
func (AlwaysActive) HasActiveDamaged(*CardInstance) bool                   { return true }
func (AlwaysActive) HasActiveFriendlyDamagedFromHidden(*CardInstance) bool { return true }
func (AlwaysActive) HasActiveSpellCast(*CardInstance) bool                 { return true }
func (AlwaysActive) HasActiveSpellHitBeforeDamage(*CardInstance) bool      { return true }
func (AlwaysActive) HasActiveSpellHit(*CardInstance) bool                  { return true }
func (AlwaysActive) HasActiveSpellReaction(*CardInstance) bool             { return true }
func (AlwaysActive) HasActiveDefend(*CardInstance) bool                    { return true }
func (AlwaysActive) HasActiveEquip(*CardInstance) bool                     { return true }
func (AlwaysActive) HasActiveUseItem(*CardInstance) bool                   { return true }
func (AlwaysActive) HasActiveConsume(*CardInstance) bool                   { return true }
func (AlwaysActive) HasActivePerTurn(*CardInstance) bool                   { return true }
func (AlwaysActive) HasActiveUltimate(*CardInstance) bool                  { return true }
func (AlwaysActive) HasActiveDevourRequirement(*CardInstance) bool         { return true }
func (AlwaysActive) HasActiveRush(*CardInstance) bool                      { return true }
func (AlwaysActive) HasActivePierce(*CardInstance) bool                    { return true }
func (AlwaysActive) HasActiveTemporary(*CardInstance) bool                 { return true }
func (AlwaysActive) HasActiveTaunt(*CardInstance) bool                     { return true }
func (AlwaysActive) HasActiveStealth(*CardInstance) bool                   { return true }
func (AlwaysActive) HasActiveShield(*CardInstance) bool                    { return true }
func (AlwaysActive) HasActiveShielding(*CardInstance) bool                 { return true }
func (AlwaysActive) HasActiveGlobalSpellRange(*CardInstance) bool          { return true }
func (AlwaysActive) HasActiveCooldown(*CardInstance) bool                  { return true }
func (AlwaysActive) HasActivePerTurnLimit(*CardInstance) bool              { return true }
func (AlwaysActive) HasActiveOverload(*CardInstance) bool                  { return true }
func (AlwaysActive) HasActiveMastery(*CardInstance) bool                   { return true }
func (AlwaysActive) HasActiveSpellTargeting(*CardInstance) bool            { return true }
func (AlwaysActive) HasActiveFriendlySpellTarget(*CardInstance) bool       { return true }
func (AlwaysActive) HasActiveSpellArea(*CardInstance) bool                 { return true }
func (AlwaysActive) HasActiveSpellAreaModifier(*CardInstance) bool         { return true }
func (AlwaysActive) HasActiveDrawReveal(*CardInstance) bool                { return true }
func (AlwaysActive) HasActiveDraw(*CardInstance) bool                      { return true }
func (AlwaysActive) HasActiveLoadGain(*CardInstance) bool                  { return true }
func (AlwaysActive) HasActiveMasteryAchieved(*CardInstance) bool           { return true }
func (AlwaysActive) HasActivePrayer(*CardInstance) bool                    { return true }
func (AlwaysActive) HasActiveSkillUsability(*CardInstance) bool            { return true }
func (AlwaysActive) HasActiveDefenseOnlySkill(*CardInstance) bool          { return true }
func (AlwaysActive) HasActiveSorcerySkill(*CardInstance) bool              { return true }
func (AlwaysActive) HasActiveSpellHitStatus(*CardInstance) bool            { return true }
func (AlwaysActive) HasActiveSpellElementGain(*CardInstance) bool          { return true }
func (AlwaysActive) HasActiveSpellDamage(*CardInstance) bool               { return true }
func (AlwaysActive) HasActiveSkillUseCostModifier(*CardInstance) bool      { return true }
func (AlwaysActive) HasActiveCardPlayCostModifier(*CardInstance) bool      { return true }
func (AlwaysActive) HasActiveCardPlayCostPaid(*CardInstance) bool          { return true }
func (AlwaysActive) HasActiveSkillUsePermissionModifier(*CardInstance) bool {
	return true
}
func (AlwaysActive) HasActiveSpellStatModifier(*CardInstance) bool { return true }
func (AlwaysActive) HasActiveEnemySpellStatModifier(*CardInstance) bool {
	return true
}
func (AlwaysActive) HasActiveDamagePrevention(*CardInstance) bool { return true }
func (AlwaysActive) HasActiveFieldDamagePrevention(*CardInstance) bool {
	return true
}
func (AlwaysActive) HasActiveNegativeStatusImmunity(*CardInstance) bool {
	return true
}
func (AlwaysActive) HasActiveAdjacentNegativeStatusProtection(*CardInstance) bool {
	return true
}
func (AlwaysActive) HasActiveSkillContributionModifier(*CardInstance) bool {
	return true
}

type OnEnterBehavior interface {
	HasActiveOnEnter(*CardInstance) bool
	OnEnter(*EffectContext) error
}

type OnGameStartBehavior interface {
	HasActiveGameStart(*CardInstance) bool
	OnGameStart(*EffectContext) error
}

type OnDeathBehavior interface {
	HasActiveDeathrattle(*CardInstance) bool
	OnDeath(*EffectContext) error
}

type OnTurnStartBehavior interface {
	HasActiveTurnStart(*CardInstance) bool
	OnTurnStart(*EffectContext) error
}

type OnTurnEndBehavior interface {
	HasActiveTurnEnd(*CardInstance) bool
	OnTurnEnd(*EffectContext) error
}

type OnUnitEnterBehavior interface {
	HasActiveUnitEnter(*CardInstance) bool
	OnUnitEnter(*EffectContext) error
}

type OnFriendlyDeathBehavior interface {
	HasActiveFriendlyDeath(*CardInstance) bool
	OnFriendlyDeath(*EffectContext) error
}

type OnEnemyDeathBehavior interface {
	HasActiveEnemyDeath(*CardInstance) bool
	OnEnemyDeath(*EffectContext) error
}

type OnDamagedBehavior interface {
	HasActiveDamaged(*CardInstance) bool
	OnDamaged(*EffectContext) error
}

type OnFriendlyDamagedFromHiddenBehavior interface {
	HasActiveFriendlyDamagedFromHidden(*CardInstance) bool
	OnFriendlyDamagedFromHidden(*EffectContext) error
}

type OnSpellCastBehavior interface {
	HasActiveSpellCast(*CardInstance) bool
	OnSpellCast(*EffectContext) error
}

type OnSpellHitBeforeDamageBehavior interface {
	HasActiveSpellHitBeforeDamage(*CardInstance) bool
	OnSpellHitBeforeDamage(*EffectContext) error
}

type OnSpellHitBehavior interface {
	HasActiveSpellHit(*CardInstance) bool
	OnSpellHit(*EffectContext) error
}

type SpellReactionBehavior interface {
	HasActiveSpellReaction(*CardInstance) bool
	CanReactToSpell(*EffectContext, *SpellCast) bool
	OnSpellReaction(*EffectContext, *SpellCast) error
}

type OnDefendBehavior interface {
	HasActiveDefend(*CardInstance) bool
	OnDefend(*EffectContext) error
}

type OnEquipBehavior interface {
	HasActiveEquip(*CardInstance) bool
	OnEquip(*EffectContext) error
}

type OnUseItemBehavior interface {
	HasActiveUseItem(*CardInstance) bool
	OnUseItem(*EffectContext) error
}

type OnConsumeBehavior interface {
	HasActiveConsume(*CardInstance) bool
	OnConsume(*EffectContext) error
}

type SummonDevourRequirementBehavior interface {
	HasActiveDevourRequirement(*CardInstance) bool
	DevourRequirement() map[string]int
}

type RushBehavior interface {
	HasActiveRush(*CardInstance) bool
	HasRush() bool
}

type PierceBehavior interface {
	HasActivePierce(*CardInstance) bool
	HasPierce() bool
}

type TemporaryBehavior interface {
	HasActiveTemporary(*CardInstance) bool
	IsTemporary() bool
}

type TauntBehavior interface {
	HasActiveTaunt(*CardInstance) bool
	HasTaunt() bool
}

type StealthBehavior interface {
	HasActiveStealth(*CardInstance) bool
	StealthLayers() int
}

type ShieldBehavior interface {
	HasActiveShield(*CardInstance) bool
	ShieldLayers() int
}

type ShieldingBehavior interface {
	HasActiveShielding(*CardInstance) bool
	HasShielding() bool
}

type GlobalSpellRangeBehavior interface {
	HasActiveGlobalSpellRange(*CardInstance) bool
	HasGlobalSpellRange() bool
}

type CooldownBehavior interface {
	HasActiveCooldown(*CardInstance) bool
	Cooldown() int
}

type PerTurnLimitBehavior interface {
	HasActivePerTurnLimit(*CardInstance) bool
	PerTurnLimit() int
}

type OverloadBehavior interface {
	HasActiveOverload(*CardInstance) bool
	OverloadThreshold() int
}

type MasteryBehavior interface {
	HasActiveMastery(*CardInstance) bool
	MasteryMax() int
	OnMastery(*EffectContext, int) error
}

type SpellTargetingBehavior interface {
	HasActiveSpellTargeting(*CardInstance) bool
	NeedsSpellTarget() bool
}

type FriendlySpellTargetBehavior interface {
	HasActiveFriendlySpellTarget(*CardInstance) bool
	AllowsFriendlySpellTarget() bool
}

type SpellAreaBehavior interface {
	HasActiveSpellArea(*CardInstance) bool
	SpellArea() SpellArea
}

type SpellAreaModifier interface {
	HasActiveSpellAreaModifier(*CardInstance) bool
	ModifySpellArea(*EffectContext, *SpellArea)
}

type DrawRevealBehavior interface {
	HasActiveDrawReveal(*CardInstance) bool
	RevealsOnDraw() bool
}

type OnDrawBehavior interface {
	HasActiveDraw(*CardInstance) bool
	OnDraw(*EffectContext) error
}

type OnSelfDrawBehavior interface {
	HasActiveDraw(*CardInstance) bool
	OnSelfDraw(*EffectContext) error
}

type OnLoadGainBehavior interface {
	HasActiveLoadGain(*CardInstance) bool
	OnLoadGain(*EffectContext) error
}

type OnMasteryAchievedBehavior interface {
	HasActiveMasteryAchieved(*CardInstance) bool
	OnMasteryAchieved(*EffectContext, int) error
}

type PrayerAbility interface {
	HasActivePrayer(*CardInstance) bool
	IsPrayerAbility() bool
}

type OptionalPrayerAbility interface {
	IsPrayerOptional(*CardInstance) bool
}

type SkillUsabilityBehavior interface {
	HasActiveSkillUsability(*CardInstance) bool
	CanUseForSkillPurpose(skillPurpose) bool
}

type DefenseOnlySkillBehavior interface {
	HasActiveDefenseOnlySkill(*CardInstance) bool
	IsDefenseOnlySkill() bool
}

type SorcerySkillBehavior interface {
	HasActiveSorcerySkill(*CardInstance) bool
	IsSorcerySkill() bool
}

type SpellHitStatusBehavior interface {
	HasActiveSpellHitStatus(*CardInstance) bool
	SpellHitStatuses(*EffectContext) map[string]int
}

type SpellElementGainBehavior interface {
	HasActiveSpellElementGain(*CardInstance) bool
	SpellElementGains(*EffectContext) map[string]int
}

type PerTurnAbility interface {
	HasActivePerTurn(*CardInstance) bool
	OnPerTurn(*EffectContext) error
}

type UltimateAbility interface {
	HasActiveUltimate(*CardInstance) bool
	OnUltimate(*EffectContext) error
}

type SpellDamageBehavior interface {
	HasActiveSpellDamage(*CardInstance) bool
	SpellDamage(*EffectContext) int
}

type SkillUseCostModifier interface {
	HasActiveSkillUseCostModifier(*CardInstance) bool
	ModifySkillUseCost(*EffectContext, map[string]int)
}

type CardPlayCostModifier interface {
	HasActiveCardPlayCostModifier(*CardInstance) bool
	ModifyCardPlayCost(*EffectContext, *CardInstance, map[string]int)
}

type CardPlayCostPaidBehavior interface {
	HasActiveCardPlayCostPaid(*CardInstance) bool
	OnCardPlayCostPaid(*EffectContext, *CardInstance)
}

type SkillUsePermissionModifier interface {
	HasActiveSkillUsePermissionModifier(*CardInstance) bool
	ValidateSkillUse(*EffectContext, *CardInstance, skillPurpose) error
}

type SpellStats struct {
	PowerBonus  int
	DamageBonus int
	Pierce      bool
}

type SpellStatModifier interface {
	HasActiveSpellStatModifier(*CardInstance) bool
	ModifySpellStats(*EffectContext, *SpellStats)
}

// EnemySpellStatModifier is evaluated from the non-casting player's field when
// an enemy spell is resolving. Use this for cards whose text says
// "敌方法术..." rather than making normal friendly spell bonuses inspect text.
type EnemySpellStatModifier interface {
	HasActiveEnemySpellStatModifier(*CardInstance) bool
	ModifyEnemySpellStats(*EffectContext, *SpellStats)
}

type DamagePreventionBehavior interface {
	HasActiveDamagePrevention(*CardInstance) bool
	PreventsDamage(*EffectContext) bool
}

type FieldDamagePreventionBehavior interface {
	HasActiveFieldDamagePrevention(*CardInstance) bool
	PreventsFieldDamage(*EffectContext) bool
}

type NegativeStatusImmunityBehavior interface {
	HasActiveNegativeStatusImmunity(*CardInstance) bool
	HasNegativeStatusImmunity() bool
}

type AdjacentNegativeStatusProtectionBehavior interface {
	HasActiveAdjacentNegativeStatusProtection(*CardInstance) bool
	ProtectsAdjacentFromNegativeStatus() bool
}

// SkillContributionModifier lets a concrete skill card alter its own
// contribution when it is used to attack, defend, or boost another spell.
type SkillContributionModifier interface {
	HasActiveSkillContributionModifier(*CardInstance) bool
	ModifySkillContribution(*EffectContext, *SpellStats)
}

func registerBehavior(r *EffectRegistry, behavior CardBehavior) {
	id := behavior.ID()
	if h, ok := behavior.(OnEnterBehavior); ok {
		r.Register(id, TriggerOnEnter, func(ctx *EffectContext) error {
			if !h.HasActiveOnEnter(ctx.Source) {
				return nil
			}
			return h.OnEnter(ctx)
		})
	}
	if h, ok := behavior.(OnGameStartBehavior); ok {
		r.Register(id, TriggerOnGameStart, func(ctx *EffectContext) error {
			if !h.HasActiveGameStart(ctx.Source) {
				return nil
			}
			return h.OnGameStart(ctx)
		})
	}
	if h, ok := behavior.(OnDeathBehavior); ok {
		r.Register(id, TriggerOnDeath, func(ctx *EffectContext) error {
			if !h.HasActiveDeathrattle(ctx.Source) {
				return nil
			}
			return h.OnDeath(ctx)
		})
	}
	if h, ok := behavior.(OnTurnStartBehavior); ok {
		r.Register(id, TriggerOnTurnStart, func(ctx *EffectContext) error {
			if !h.HasActiveTurnStart(ctx.Source) {
				return nil
			}
			return h.OnTurnStart(ctx)
		})
	}
	if h, ok := behavior.(OnTurnEndBehavior); ok {
		r.Register(id, TriggerOnTurnEnd, func(ctx *EffectContext) error {
			if !h.HasActiveTurnEnd(ctx.Source) {
				return nil
			}
			return h.OnTurnEnd(ctx)
		})
	}
	if h, ok := behavior.(OnUnitEnterBehavior); ok {
		r.Register(id, TriggerOnUnitEnter, func(ctx *EffectContext) error {
			if !h.HasActiveUnitEnter(ctx.Source) {
				return nil
			}
			return h.OnUnitEnter(ctx)
		})
	}
	if h, ok := behavior.(OnFriendlyDeathBehavior); ok {
		r.Register(id, TriggerOnFriendlyDeath, func(ctx *EffectContext) error {
			if !h.HasActiveFriendlyDeath(ctx.Source) {
				return nil
			}
			return h.OnFriendlyDeath(ctx)
		})
	}
	if h, ok := behavior.(OnEnemyDeathBehavior); ok {
		r.Register(id, TriggerOnEnemyDeath, func(ctx *EffectContext) error {
			if !h.HasActiveEnemyDeath(ctx.Source) {
				return nil
			}
			return h.OnEnemyDeath(ctx)
		})
	}
	if h, ok := behavior.(OnDamagedBehavior); ok {
		r.Register(id, TriggerOnDamaged, func(ctx *EffectContext) error {
			if !h.HasActiveDamaged(ctx.Source) {
				return nil
			}
			return h.OnDamaged(ctx)
		})
	}
	if h, ok := behavior.(OnSpellCastBehavior); ok {
		r.Register(id, TriggerOnSpellCast, func(ctx *EffectContext) error {
			if !h.HasActiveSpellCast(ctx.Source) {
				return nil
			}
			return h.OnSpellCast(ctx)
		})
	}
	if h, ok := behavior.(OnSpellHitBeforeDamageBehavior); ok {
		r.Register(id, TriggerOnSpellHitBeforeDamage, func(ctx *EffectContext) error {
			if !h.HasActiveSpellHitBeforeDamage(ctx.Source) {
				return nil
			}
			return h.OnSpellHitBeforeDamage(ctx)
		})
	}
	if h, ok := behavior.(OnSpellHitBehavior); ok {
		r.Register(id, TriggerOnSpellHit, func(ctx *EffectContext) error {
			if !h.HasActiveSpellHit(ctx.Source) {
				return nil
			}
			return h.OnSpellHit(ctx)
		})
	}
	if h, ok := behavior.(OnDefendBehavior); ok {
		r.Register(id, TriggerOnDefend, func(ctx *EffectContext) error {
			if !h.HasActiveDefend(ctx.Source) {
				return nil
			}
			return h.OnDefend(ctx)
		})
	}
	if h, ok := behavior.(OnEquipBehavior); ok {
		r.Register(id, TriggerOnEquip, func(ctx *EffectContext) error {
			if !h.HasActiveEquip(ctx.Source) {
				return nil
			}
			return h.OnEquip(ctx)
		})
	}
	if h, ok := behavior.(OnUseItemBehavior); ok {
		r.Register(id, TriggerOnUseItem, func(ctx *EffectContext) error {
			if !h.HasActiveUseItem(ctx.Source) {
				return nil
			}
			return h.OnUseItem(ctx)
		})
	}
	if h, ok := behavior.(OnConsumeBehavior); ok {
		r.Register(id, TriggerOnConsume, func(ctx *EffectContext) error {
			if !h.HasActiveConsume(ctx.Source) {
				return nil
			}
			return h.OnConsume(ctx)
		})
	}
	if h, ok := behavior.(OnDrawBehavior); ok {
		r.Register(id, TriggerOnDraw, func(ctx *EffectContext) error {
			if !h.HasActiveDraw(ctx.Source) {
				return nil
			}
			return h.OnDraw(ctx)
		})
	}
	if h, ok := behavior.(OnLoadGainBehavior); ok {
		r.Register(id, TriggerOnLoadGain, func(ctx *EffectContext) error {
			if !h.HasActiveLoadGain(ctx.Source) {
				return nil
			}
			return h.OnLoadGain(ctx)
		})
	}
	if h, ok := behavior.(OnMasteryAchievedBehavior); ok {
		r.Register(id, TriggerOnMastery, func(ctx *EffectContext) error {
			if !h.HasActiveMasteryAchieved(ctx.Source) {
				return nil
			}
			level, _ := ctx.ExtraData["mastery"].(int)
			return h.OnMasteryAchieved(ctx, level)
		})
	}
	if h, ok := behavior.(PerTurnAbility); ok {
		if prayer, ok := behavior.(PrayerAbility); !ok || !prayer.IsPrayerAbility() {
			r.RegisterActive(id, TriggerPerTurn, func(ctx *EffectContext) error {
				if !h.HasActivePerTurn(ctx.Source) {
					return nil
				}
				return h.OnPerTurn(ctx)
			})
		}
	}
	if h, ok := behavior.(UltimateAbility); ok {
		r.RegisterActive(id, TriggerUltimate, func(ctx *EffectContext) error {
			if !h.HasActiveUltimate(ctx.Source) {
				return nil
			}
			return h.OnUltimate(ctx)
		})
	}
	if h, ok := behavior.(SpellDamageBehavior); ok {
		r.RegisterSpellDamage(id, func(ctx *EffectContext) int {
			if !h.HasActiveSpellDamage(ctx.Source) {
				return 0
			}
			return h.SpellDamage(ctx)
		})
	}
}

type noopPerTurn struct{ AlwaysActive }

func (noopPerTurn) OnPerTurn(*EffectContext) error { return nil }

type noopUltimate struct{ AlwaysActive }

func (noopUltimate) OnUltimate(*EffectContext) error { return nil }
