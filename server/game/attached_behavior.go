package game

import ()

// AttachedBehavior is a runtime-granted behavior on one card instance.
//
// Use this for effects such as "give a unit deathrattle" where the target card
// did not have that behavior in its printed Go struct. Do not model these as
// free-form status strings; attached behaviors participate in the same
// interface-based rules as printed card behaviors.
type AttachedBehavior interface {
	AttachedID() string
	AttachedName() string
	AttachedInfo() AttachedBehaviorInfo
}

type AttachedBehaviorInfo struct {
	ID   string         `json:"id"`
	Name string         `json:"name"`
	Data map[string]any `json:"data,omitempty"`
}

type AttachedDeathrattleBehavior interface {
	AttachedBehavior
	HasActiveDeathrattle(*CardInstance) bool
	OnDeath(*EffectContext) error
}

func (ci *CardInstance) AddAttachedBehavior(behavior AttachedBehavior) {
	if ci == nil || behavior == nil {
		return
	}
	ci.AttachedBehaviors = append(ci.AttachedBehaviors, behavior)
}

func attachedBehaviorsInfo(card *CardInstance) []AttachedBehaviorInfo {
	if card == nil || len(card.AttachedBehaviors) == 0 {
		return nil
	}
	info := make([]AttachedBehaviorInfo, 0, len(card.AttachedBehaviors))
	for _, behavior := range card.AttachedBehaviors {
		if behavior == nil {
			continue
		}
		info = append(info, behavior.AttachedInfo())
	}
	return info
}

func attachedDeathrattles(card *CardInstance) []AttachedDeathrattleBehavior {
	if card == nil || len(card.AttachedBehaviors) == 0 {
		return nil
	}
	var deathrattles []AttachedDeathrattleBehavior
	for _, behavior := range card.AttachedBehaviors {
		deathrattle, ok := behavior.(AttachedDeathrattleBehavior)
		if ok && deathrattle.HasActiveDeathrattle(card) {
			deathrattles = append(deathrattles, deathrattle)
		}
	}
	return deathrattles
}

type AttachedDeathrattleDamageEnemyHero struct {
	Amount int
}

func (a AttachedDeathrattleDamageEnemyHero) AttachedID() string {
	return "deathrattle_damage_enemy_hero"
}

func (a AttachedDeathrattleDamageEnemyHero) AttachedName() string {
	return "遗言: 对敌方英雄造成伤害"
}

func (a AttachedDeathrattleDamageEnemyHero) AttachedInfo() AttachedBehaviorInfo {
	return AttachedBehaviorInfo{
		ID:   a.AttachedID(),
		Name: a.AttachedName(),
		Data: map[string]any{"amount": a.Amount},
	}
}

func (a AttachedDeathrattleDamageEnemyHero) HasActiveDeathrattle(card *CardInstance) bool {
	return card != nil && a.Amount > 0
}

func (a AttachedDeathrattleDamageEnemyHero) OnDeath(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || a.Amount <= 0 {
		return nil
	}
	opponent := ctx.Engine.State.Players[ctx.OpponentID]
	if opponent == nil || opponent.Hero == nil {
		return nil
	}
	ctx.DealDamage(opponent.Hero, a.Amount)
	return nil
}
