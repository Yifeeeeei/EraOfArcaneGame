package game

// CounterBehavior describes a set card's response windows. Adding a counter
// does not require edits to the engine's trigger or card-number dispatch lists.
type CounterBehavior interface {
	HasActiveCounter(*CardInstance) bool
	CounterTriggers() []EffectTrigger
	CanTriggerCounter(*CounterContext) bool
}

// CounterResolutionBehavior is needed only when revealing the counter changes
// the pending operation itself instead of invoking its ordinary card effect.
type CounterResolutionBehavior interface {
	ResolveCounter(*CounterContext)
}

func (AlwaysActive) HasActiveCounter(*CardInstance) bool { return true }

type CounterEvent struct {
	Trigger       EffectTrigger
	Card          *CardInstance
	PlayerID      int
	Power         int
	Damage        int
	DrawCount     int
	DefenderID    int
	LearnedSkill  bool
	BoostUse      bool
	DefenseSkills []*CardInstance
	DefenseBoosts []*CardInstance
	CancelItem    *bool
	CancelHit     *bool
	PreventDamage *bool
}

func (event CounterEvent) IsLethal() bool {
	return event.Card != nil && event.Card.Card != nil && event.Damage > 0 && event.Card.CurrentLife-event.Damage <= 0
}

type CounterContext struct {
	Engine *Engine
	Source *CardInstance
	Event  CounterEvent
	// Adapter payload belongs to engine dispatch, never to card predicates.
	payload map[string]any
}

func (e *Engine) counterContext(source *CardInstance, trigger EffectTrigger, card *CardInstance, data map[string]any) *CounterContext {
	actor := -1
	if card != nil {
		actor = card.OwnerID
	}
	if trigger != TriggerOnFriendlyDeath && trigger != TriggerOnEnemyDeath {
		for _, key := range []string{"cast_player", "attacker", "drawn_player", "damaged_player", "entered_player", "consumed_player", "used_player", "ended_player"} {
			if value, ok := data[key].(int); ok {
				actor = value
			}
		}
	}
	cancelItem, _ := data["cancel_item"].(*bool)
	preventDamage, _ := data["prevent_damage"].(*bool)
	cancelHit, _ := data["cancel_spell_hit"].(*bool)
	return &CounterContext{Engine: e, Source: source, payload: data, Event: CounterEvent{
		Trigger: trigger, Card: card, PlayerID: actor,
		Power: spellPowerFromData(data), Damage: damageFromData(data), DrawCount: drawCountFromData(data),
		DefenderID: intFromData(data, "defender", -1), LearnedSkill: boolFromData(data, "learned_skill"),
		BoostUse: boolFromData(data, "boost_use"), DefenseSkills: spellInstancesFromData(data, "defense_skills"),
		DefenseBoosts: spellInstancesFromData(data, "defense_boosts"), CancelItem: cancelItem, CancelHit: cancelHit, PreventDamage: preventDamage,
	}}
}

func (ctx *CounterContext) CancelItem() {
	if cancelled := ctx.Event.CancelItem; cancelled != nil {
		*cancelled = true
		ctx.Engine.emit(GameEvent{Type: "item_cancelled", Player: -1, Data: map[string]any{
			"player": ctx.Source.OwnerID, "card": cardToInfo(ctx.Event.Card), "source": cardToInfo(ctx.Source),
		}})
	}
}

func (ctx *CounterContext) CancelHit() {
	if cancelled := ctx.Event.CancelHit; cancelled != nil {
		*cancelled = true
		ctx.Engine.emit(GameEvent{Type: "spell_hit_cancelled", Player: -1, Data: map[string]any{
			"player": ctx.Source.OwnerID, "skill": cardToInfo(ctx.Event.Card), "source": cardToInfo(ctx.Source),
		}})
	}
}

func (ctx *CounterContext) CancelBoost() {
	ctx.Engine.cancelBoostSpellWithIceSoulSeal(ctx.Event.Card, ctx.payload)
	ctx.Engine.emit(GameEvent{Type: "boost_spell_cancelled", Player: -1, Data: map[string]any{
		"player": ctx.Source.OwnerID, "skill": cardToInfo(ctx.Event.Card), "source": cardToInfo(ctx.Source),
	}})
}

func (ctx *CounterContext) ReflectSpell() {
	ctx.Engine.promptCounterWindHoleScroll(ctx.Source, ctx.Event.Card, ctx.payload)
}
