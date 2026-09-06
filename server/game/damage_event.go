package game

import "eraofarcane/model"

// DamageEvent describes the damaged unit independently of the card observing
// the event. An unknown source is -1; it must never silently become player 0.
type DamageEvent struct {
	Target       *CardInstance
	SourcePlayer int
	Amount       int
	Kind         string
	Element      string
	Status       string
	Spell        *model.Card
	BoostCount   int
	Prevent      *bool
}

type DamageScope int

const (
	DamageSelf DamageScope = iota
	DamageOtherFriendly
	DamageFriendly
	DamageEnemy
	DamageAny
)

func (event DamageEvent) Matches(observer *CardInstance, scope DamageScope) bool {
	if observer == nil || event.Target == nil || event.Target.Card == nil {
		return false
	}
	switch scope {
	case DamageSelf:
		return observer == event.Target
	case DamageOtherFriendly:
		return observer != event.Target && observer.OwnerID == event.Target.OwnerID
	case DamageFriendly:
		return observer.OwnerID == event.Target.OwnerID
	case DamageEnemy:
		return observer.OwnerID != event.Target.OwnerID
	case DamageAny:
		return true
	default:
		return false
	}
}

func (event DamageEvent) IsEnemyDamage(playerID int) bool {
	return event.SourcePlayer >= 0 && event.SourcePlayer != playerID
}

func (event DamageEvent) IsFire() bool {
	return event.Element == model.ElementFire || event.Status == StatusBurn
}

// damageEventFromContext is the adapter for legacy trigger producers. Rules
// consume DamageEvent, never the transport/legacy map. Producers are migrated
// independently of card predicates.
func damageEventFromContext(ctx *EffectContext) DamageEvent {
	target := ctx.Target
	if target == nil {
		target = ctx.Source
	}
	event := DamageEvent{Target: target, SourcePlayer: intFromData(ctx.ExtraData, "attacker", -1),
		Amount: damageFromData(ctx.ExtraData), BoostCount: intFromData(ctx.ExtraData, "boost_count", -1)}
	event.Kind, _ = ctx.ExtraData["damage_source"].(string)
	event.Element, _ = ctx.ExtraData["damage_element"].(string)
	event.Status, _ = ctx.ExtraData["status_damage"].(string)
	event.Prevent, _ = ctx.ExtraData["prevent_damage"].(*bool)
	if number, ok := ctx.ExtraData["skill"].(string); ok {
		event.Spell = getCardDB()[number]
	}
	return event
}

// DealDamage attributes card-effect damage to its source and derives the
// affected player from the target. A caller cannot accidentally swap these.
func (ctx *EffectContext) DealDamage(target *CardInstance, amount int) {
	if ctx == nil || ctx.Engine == nil || target == nil {
		return
	}
	ctx.Engine.ApplyDamage(DamageRequest{Target: target, Amount: amount, Kind: "effect", SourcePlayer: ctx.PlayerID, SourceKnown: true, Source: ctx.Source})
}
