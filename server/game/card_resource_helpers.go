package game

import (
	"eraofarcane/model"
	"fmt"
)

const (
	fireNegativeStatusImmunityUntil = "fire_negative_status_immunity_until"
	arcaneCylinderCounter           = "奥术魔法筒标记物"
	fireBoxCounter                  = "火匣子标记物"
	waterAriaCounter                = "水之咏叹标记物"
	windQuillCounter                = "聆风羽毛笔标记物"
	forestStorageCounter            = "森之贮藏标记物"
	blessingStaffCounter            = "祝福之杖标记物"
	burierCounter                   = "埋葬者标记物"
)

func grantFireNegativeStatusImmunity(ctx *EffectContext) {
	if ctx == nil || ctx.Engine == nil {
		return
	}
	until := ctx.Engine.State.TurnNumber + 2
	for _, card := range ctx.Engine.getAllFieldCards(ctx.Engine.State.Players[ctx.PlayerID]) {
		if card != nil && card.Card != nil && card.Card.Category == model.ElementFire {
			card.Statuses[fireNegativeStatusImmunityUntil] = until
		}
	}
}

type markerEquipment struct {
	AlwaysActive
	id       string
	name     string
	counter  string
	counters int
	limit    int
	canUse   func(*EffectContext) bool
	effect   func(*EffectContext) error
}

func (m markerEquipment) ID() string { return m.id }

func (m markerEquipment) Name() string { return m.name }

func (m markerEquipment) PerTurnLimit() int {
	if m.limit > 0 {
		return m.limit
	}
	return 1
}

func (m markerEquipment) OnEnter(ctx *EffectContext) error {
	ctx.Source.Statuses[m.counter] += m.counters
	return nil
}

func (m markerEquipment) OnPerTurn(ctx *EffectContext) error {
	if ctx.Source.IsHorizontal {
		return fmt.Errorf("%s is horizontal", m.name)
	}
	if ctx.Source.Statuses[m.counter] <= 0 {
		return fmt.Errorf("%s has no markers", m.name)
	}
	if m.canUse != nil && !m.canUse(ctx) {
		return fmt.Errorf("%s has no valid target", m.name)
	}
	ctx.Source.IsHorizontal = true
	ctx.Source.Statuses[m.counter]--
	if m.effect != nil {
		return m.effect(ctx)
	}
	return nil
}

func newArcaneCylinder() CardBehavior {
	return markerEquipment{id: "2021023", name: "奥术魔法筒", counter: arcaneCylinderCounter, counters: 3, effect: func(ctx *EffectContext) error {
		ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementArcane: 2})
		return nil
	}}
}

func newFireBox() CardBehavior {
	return markerEquipment{id: "2121014", name: "火匣子", counter: fireBoxCounter, counters: 3, effect: func(ctx *EffectContext) error {
		ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementFire: 2})
		return nil
	}}
}

func newWaterAria() CardBehavior {
	return markerEquipment{id: "2221014", name: "水之咏叹", counter: waterAriaCounter, counters: 4, effect: func(ctx *EffectContext) error {
		ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementWater: 3})
		return nil
	}}
}

func newWindQuill() CardBehavior {
	return markerEquipment{id: "2321014", name: "聆风羽毛笔", counter: windQuillCounter, counters: 3, effect: func(ctx *EffectContext) error {
		ctx.Engine.drawCards(ctx.PlayerID, 1)
		ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementAir: 1})
		return nil
	}}
}

func newForestStorage() CardBehavior {
	return markerEquipment{id: "2421014", name: "森之贮藏", counter: forestStorageCounter, counters: 4, limit: 99, effect: func(ctx *EffectContext) error {
		ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementEarth: 4})
		return nil
	}}
}

func newBlessingStaff() CardBehavior {
	canUse := func(ctx *EffectContext) bool {
		return len(ctx.Engine.friendlyUnits(ctx.PlayerID, true, nil)) > 0
	}
	return markerEquipment{id: "2521014", name: "祝福之杖", counter: blessingStaffCounter, counters: 3, canUse: canUse, effect: func(ctx *EffectContext) error {
		candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, nil)
		ctx.Engine.SetPendingAction(ctx.PlayerID, "blessing_staff", "祝福之杖:选择1个友方单位+1血", candidates, 1, 1, func(selected []string) {
			target := selectedUnitFromCandidates(ctx.Engine, selected, candidates)
			if target != nil {
				target.Statuses["max_life_bonus"]++
				ctx.Engine.gainLife(target, 1, ctx.Source)
			}
			ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementLight: 2})
		})
		return nil
	}}
}

func newBurier() CardBehavior {
	return markerEquipment{id: "2621014", name: "埋葬者", counter: burierCounter, counters: 3, effect: func(ctx *EffectContext) error {
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		count := min(2, len(ps.Deck))
		for i := 0; i < count; i++ {
			card := ps.Deck[0]
			ps.Deck = ps.Deck[1:]
			ctx.Engine.addToGraveyard(ctx.PlayerID, card)
			ctx.Engine.emit(GameEvent{Type: "discard", Player: ctx.PlayerID, Data: map[string]any{"card": cardToInfo(card)}})
		}
		ps.GainElements(map[string]int{model.ElementShadow: 2})
		return nil
	}}
}

func candidateContains(candidates []map[string]any, instanceID string) bool {
	for _, candidate := range candidates {
		if candidate["instance_id"] == instanceID {
			return true
		}
	}
	return false
}
