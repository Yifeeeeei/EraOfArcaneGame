package game

import (
	"fmt"
	"log"
)

// EffectTrigger represents when an effect fires
type EffectTrigger int

const (
	TriggerOnEnter       EffectTrigger = iota // 入场: card enters the field
	TriggerOnDeath                            // 遗言: card dies / leaves field to graveyard
	TriggerOnTurnStart                        // 回合开始: owner's turn starts
	TriggerOnTurnEnd                          // 回合结束: owner's turn ends
	TriggerOnConsume                          // 消耗时: card is consumed (横置)
	TriggerOnAttack                           // 攻击时: unit declares attack
	TriggerOnHit                              // 命中: attack hits (not defended)
	TriggerOnDamaged                          // 受伤时: card takes damage
	TriggerOnSpellCast                        // 施法时: spell is cast
	TriggerOnSpellHit                         // 法术命中: spell hits target
	TriggerOnDefend                           // 防御时: spell is defended
	TriggerOnDraw                             // 抽牌时: card is drawn
	TriggerOnSummon                           // 召唤时: any friendly unit is summoned
	TriggerOnFriendlyDeath                    // 友方死亡: any friendly unit dies
	TriggerOnEnemyDeath                       // 敌方死亡: any enemy unit dies
	TriggerOnUnitEnter                        // 任意入场: any unit enters field (for passive listeners)
	TriggerPerTurn                            // 回合技: active ability (per-turn)
	TriggerUltimate                           // 绝技: one-time active ability
	TriggerOnEquip                            // 装备时: item is equipped
	TriggerPassive                            // 被动: always active while on field
)

// EffectContext provides context for effect execution
type EffectContext struct {
	Engine     *Engine
	Source     *CardInstance // the card triggering the effect
	Target     *CardInstance // optional target
	TargetPos  *Position     // optional target position
	PlayerID   int           // source card's owner
	OpponentID int           // opponent of source card's owner
	DamageAmount int         // for damage-related triggers
	ExtraData  map[string]any // additional context data
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

// EffectRegistry holds all registered card effects
type EffectRegistry struct {
	effects map[string][]*CardEffect // cardNumber -> effects list
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
		effects: make(map[string][]*CardEffect),
	}
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

// GetEffects returns all effects for a card number and trigger type
func (r *EffectRegistry) GetEffects(cardNumber string, trigger EffectTrigger) []*CardEffect {
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
	for _, eff := range r.effects[cardNumber] {
		if eff.Trigger == trigger {
			return true
		}
	}
	return false
}

// GetAllEffects returns all effects for a card number
func (r *EffectRegistry) GetAllEffects(cardNumber string) []*CardEffect {
	return r.effects[cardNumber]
}

// ---------- Engine integration methods ----------

// triggerEffects fires all effects for a card at a specific trigger point
func (e *Engine) triggerEffects(trigger EffectTrigger, source *CardInstance, target *CardInstance, extraData map[string]any) {
	if source == nil || source.Card == nil {
		return
	}

	// Check for petrify - petrified cards have no effects
	if source.Statuses[StatusPetrify] > 0 && trigger != TriggerOnDeath {
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
}

// triggerFieldEffects fires effects for all cards on a player's field
func (e *Engine) triggerFieldEffects(trigger EffectTrigger, playerID int, eventSource *CardInstance) {
	ps := e.State.Players[playerID]
	allCards := e.getAllFieldCards(ps)
	for _, card := range allCards {
		if card == eventSource {
			continue // skip the source itself to avoid loops
		}
		e.triggerEffects(trigger, card, eventSource, nil)
	}
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
	case TriggerPerTurn:
		return "per_turn"
	case TriggerUltimate:
		return "ultimate"
	default:
		return fmt.Sprintf("trigger_%d", t)
	}
}
