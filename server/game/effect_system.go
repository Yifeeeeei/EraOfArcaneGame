package game

import (
	"fmt"
	"log"
)

// EffectTrigger represents when an effect fires
type EffectTrigger int

const (
	TriggerOnEnter                EffectTrigger = iota // 入场: card enters the field
	TriggerOnGameStart                                 // 游戏开始时: after initial hero enter, before initial draw
	TriggerOnDeath                                     // 遗言: card dies / leaves field to graveyard
	TriggerOnTurnStart                                 // 回合开始: owner's turn starts
	TriggerOnTurnEnd                                   // 回合结束: owner's turn ends
	TriggerOnConsume                                   // 消耗时: card is consumed (横置)
	TriggerOnAttack                                    // 攻击时: unit declares attack
	TriggerOnAttacked                                  // 受攻击时: unit becomes the target of a direct attack
	TriggerOnHit                                       // 命中: attack hits (not defended)
	TriggerOnDamaged                                   // 受伤时: card takes damage
	TriggerOnSpellCast                                 // 施法时: spell is cast
	TriggerOnSpellHitBeforeDamage                      // 法术命中时: before hit damage is dealt
	TriggerOnSpellHit                                  // 法术命中后: spell hit after damage
	TriggerOnDefend                                    // 防御时: spell is defended
	TriggerOnDraw                                      // 抽牌时: card is drawn
	TriggerOnLoadGain                                  // 获得负载时: a friendly card gains load
	TriggerOnMastery                                   // 达到精通时: a friendly card reaches a mastery level
	TriggerOnSummon                                    // 召唤时: any friendly unit is summoned
	TriggerOnFriendlyDeath                             // 友方死亡: any friendly unit dies
	TriggerOnEnemyDeath                                // 敌方死亡: any enemy unit dies
	TriggerOnUnitEnter                                 // 任意入场: any unit enters field (for passive listeners)
	TriggerPerTurn                                     // 回合技: active ability (per-turn)
	TriggerUltimate                                    // 绝技: one-time active ability
	TriggerOnEquip                                     // 装备时: item is equipped
	TriggerOnUseItem                                   // 使用消耗品道具时
	TriggerPassive                                     // 被动: always active while on field
)

// EffectContext provides context for effect execution
type EffectContext struct {
	Engine       *Engine
	Source       *CardInstance  // the card triggering the effect
	Target       *CardInstance  // optional target
	TargetPos    *Position      // optional target position
	PlayerID     int            // source card's owner
	OpponentID   int            // opponent of source card's owner
	DamageAmount int            // for damage-related triggers
	ExtraData    map[string]any // additional context data
}

// EffectHandler is a function that implements a card effect
type EffectHandler func(ctx *EffectContext) error

// CardEffect represents a registered effect for a specific card
type CardEffect struct {
	CardNumber string
	Trigger    EffectTrigger
	Handler    EffectHandler
	Priority   int  // higher = runs first
	IsActive   bool // true if this is an activated ability (回合技/绝技), not auto-trigger
}

// CardAbilityWithPrecondition 包含前提条件的卡牌能力
type CardAbilityWithPrecondition struct {
	Preconditions []Precondition // 前提条件列表
	Effect        EffectHandler  // 实际效果
	Trigger       EffectTrigger  // 触发时机
	IsActive      bool           // 是否为主动能力
}

// EffectRegistry holds all registered card effects
type EffectRegistry struct {
	effects              map[string][]*CardEffect                  // cardNumber -> effects list
	conditionalAbilities map[string][]*CardAbilityWithPrecondition // cardNumber -> conditional abilities
	behaviorFactories    map[string]func() CardBehavior            // cardNumber -> lazy behavior constructor
	loadedBehaviors      map[string]CardBehavior                   // cardNumber -> behavior registered into effects
	spellDamage          map[string]func(*EffectContext) int       // cardNumber -> spell damage override
}

// Global effect registry
var globalRegistry *EffectRegistry

func init() {
	globalRegistry = NewEffectRegistry()
}

// GetEffectRegistry returns the global effect registry
func GetEffectRegistry() *EffectRegistry {
	return globalRegistry
}

// NewEffectRegistry creates a new effect registry
func NewEffectRegistry() *EffectRegistry {
	return &EffectRegistry{
		effects:              make(map[string][]*CardEffect),
		conditionalAbilities: make(map[string][]*CardAbilityWithPrecondition),
		behaviorFactories:    make(map[string]func() CardBehavior),
		loadedBehaviors:      make(map[string]CardBehavior),
		spellDamage:          make(map[string]func(*EffectContext) int),
	}
}

// RegisterBehaviorFactory records how to construct a card behavior without
// allocating it at server startup. The behavior is materialized only when that
// card number is queried by the engine.
func (r *EffectRegistry) RegisterBehaviorFactory(cardNumber string, factory func() CardBehavior) {
	r.behaviorFactories[cardNumber] = factory
}

func (r *EffectRegistry) ensureBehaviorLoaded(cardNumber string) {
	if r.loadedBehaviors[cardNumber] != nil {
		return
	}
	factory, ok := r.behaviorFactories[cardNumber]
	if !ok {
		return
	}
	behavior := factory()
	r.loadedBehaviors[cardNumber] = behavior
	registerBehavior(r, behavior)
}

func (r *EffectRegistry) GetBehavior(cardNumber string) CardBehavior {
	r.ensureBehaviorLoaded(cardNumber)
	return r.loadedBehaviors[cardNumber]
}

// Register adds an effect for a card
func (r *EffectRegistry) Register(cardNumber string, trigger EffectTrigger, handler EffectHandler) {
	r.effects[cardNumber] = append(r.effects[cardNumber], &CardEffect{
		CardNumber: cardNumber,
		Trigger:    trigger,
		Handler:    handler,
	})
}

// RegisterActive adds an activated ability for a card
func (r *EffectRegistry) RegisterActive(cardNumber string, trigger EffectTrigger, handler EffectHandler) {
	r.effects[cardNumber] = append(r.effects[cardNumber], &CardEffect{
		CardNumber: cardNumber,
		Trigger:    trigger,
		Handler:    handler,
		IsActive:   true,
	})
}

func (r *EffectRegistry) RegisterSpellDamage(cardNumber string, handler func(*EffectContext) int) {
	r.spellDamage[cardNumber] = handler
}

func (r *EffectRegistry) SpellDamage(cardNumber string, ctx *EffectContext) (int, bool) {
	r.ensureBehaviorLoaded(cardNumber)
	handler, ok := r.spellDamage[cardNumber]
	if !ok {
		return 0, false
	}
	return handler(ctx), true
}

// RegisterWithPrecondition 注册带前提条件的能力
func (r *EffectRegistry) RegisterWithPrecondition(
	cardNumber string,
	trigger EffectTrigger,
	preconditions []Precondition,
	handler EffectHandler,
	isActive bool,
) {
	r.conditionalAbilities[cardNumber] = append(r.conditionalAbilities[cardNumber], &CardAbilityWithPrecondition{
		Preconditions: preconditions,
		Effect:        handler,
		Trigger:       trigger,
		IsActive:      isActive,
	})
}

// CheckPreconditions 检查卡牌的所有前提条件
func (r *EffectRegistry) CheckPreconditions(
	cardNumber string,
	trigger EffectTrigger,
	ctx *ConditionContext,
) (bool, string) {
	abilities, exists := r.conditionalAbilities[cardNumber]
	if !exists {
		return true, "" // 没有条件要求，默认通过
	}

	for _, ability := range abilities {
		if ability.Trigger != trigger {
			continue
		}

		// 检查所有前提条件
		for _, precond := range ability.Preconditions {
			result := precond.Check(ctx)
			if !result.Passed {
				return false, result.Reason
			}
		}
	}

	return true, ""
}

// GetConditionalAbilities 获取卡牌的所有条件能力
func (r *EffectRegistry) GetConditionalAbilities(cardNumber string) []*CardAbilityWithPrecondition {
	return r.conditionalAbilities[cardNumber]
}

// GetEffects returns all effects for a card number and trigger type
func (r *EffectRegistry) GetEffects(cardNumber string, trigger EffectTrigger) []*CardEffect {
	r.ensureBehaviorLoaded(cardNumber)
	var result []*CardEffect
	for _, eff := range r.effects[cardNumber] {
		if eff.Trigger == trigger {
			result = append(result, eff)
		}
	}
	return result
}

// HasEffect checks if a card has any registered effect for a trigger
func (r *EffectRegistry) HasEffect(cardNumber string, trigger EffectTrigger) bool {
	r.ensureBehaviorLoaded(cardNumber)
	for _, eff := range r.effects[cardNumber] {
		if eff.Trigger == trigger {
			return true
		}
	}
	return false
}

// GetAllEffects returns all effects for a card number
func (r *EffectRegistry) GetAllEffects(cardNumber string) []*CardEffect {
	r.ensureBehaviorLoaded(cardNumber)
	return r.effects[cardNumber]
}

// ---------- Engine integration methods ----------

// triggerEffects fires all effects for a card at a specific trigger point
func (e *Engine) triggerEffects(trigger EffectTrigger, source *CardInstance, target *CardInstance, extraData map[string]any) {
	if source == nil || source.Card == nil {
		return
	}

	// Check for petrify - petrified cards have no effects
	if e.hasEffectiveStatus(source, StatusPetrify) && trigger != TriggerOnDeath {
		return
	}

	ctx := &EffectContext{
		Engine:     e,
		Source:     source,
		Target:     target,
		PlayerID:   source.OwnerID,
		OpponentID: 1 - source.OwnerID,
		ExtraData:  extraData,
	}

	effects := globalRegistry.GetEffects(source.Card.Number, trigger)
	for _, eff := range effects {
		if eff.IsActive {
			continue // active abilities are triggered by player action, not auto
		}
		if err := eff.Handler(ctx); err != nil {
			log.Printf("[Effect] Error executing %s effect for %s: %v",
				triggerName(trigger), source.Card.Name, err)
		}
	}

	if trigger == TriggerOnDeath {
		for _, deathrattle := range attachedDeathrattles(source) {
			if err := deathrattle.OnDeath(ctx); err != nil {
				log.Printf("[Effect] Error executing attached deathrattle %s for %s: %v",
					deathrattle.AttachedID(), source.Card.Name, err)
			}
		}
	}
}

// triggerFieldEffects fires effects for all cards on a player's field
func (e *Engine) triggerFieldEffects(trigger EffectTrigger, playerID int, eventSource *CardInstance) bool {
	return e.triggerFieldEffectsWithData(trigger, playerID, eventSource, nil)
}

func (e *Engine) triggerFieldEffectsWithData(trigger EffectTrigger, playerID int, eventSource *CardInstance, extraData map[string]any) bool {
	ps := e.State.Players[playerID]
	allCards := e.getAllFieldCards(ps)
	counterCandidates := []*CardInstance{}
	skipCounters := false
	if extraData != nil {
		skipCounters, _ = extraData["skip_counter_traps"].(bool)
	}
	for _, card := range allCards {
		if card == eventSource {
			continue // skip the source itself to avoid loops
		}
		if card.IsSetCounter {
			if !skipCounters && counterTrapHasTrigger(card.Card.Number, trigger) && e.counterTrapConditionMet(card, trigger, eventSource, extraData) {
				counterCandidates = append(counterCandidates, card)
			}
			continue
		}
		if trigger == TriggerOnDefend && isSpellLikeCard(card.Card) {
			continue
		}
		e.triggerEffects(trigger, card, eventSource, extraData)
	}
	return e.promptCounterTrapQueue(counterCandidates, trigger, eventSource, extraData, nil)
}

func triggerName(t EffectTrigger) string {
	switch t {
	case TriggerOnEnter:
		return "on_enter"
	case TriggerOnDeath:
		return "on_death"
	case TriggerOnTurnStart:
		return "on_turn_start"
	case TriggerOnTurnEnd:
		return "on_turn_end"
	case TriggerOnConsume:
		return "on_consume"
	case TriggerOnAttack:
		return "on_attack"
	case TriggerOnAttacked:
		return "on_attacked"
	case TriggerOnHit:
		return "on_hit"
	case TriggerOnDamaged:
		return "on_damaged"
	case TriggerOnSpellCast:
		return "on_spell_cast"
	case TriggerOnSpellHit:
		return "on_spell_hit"
	case TriggerOnDefend:
		return "on_defend"
	case TriggerOnDraw:
		return "on_draw"
	case TriggerOnLoadGain:
		return "on_load_gain"
	case TriggerOnMastery:
		return "on_mastery"
	case TriggerPerTurn:
		return "per_turn"
	case TriggerUltimate:
		return "ultimate"
	case TriggerOnUseItem:
		return "on_use_item"
	default:
		return fmt.Sprintf("trigger_%d", t)
	}
}
