package game

import (
	"fmt"
	"math/rand"

	"eraofarcane/model"
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

type Card1021018ArcaneBulwark struct{ AlwaysActive }

func (Card1021018ArcaneBulwark) ID() string   { return "1021018" }
func (Card1021018ArcaneBulwark) Name() string { return "奥术壁垒" }
func (Card1021018ArcaneBulwark) OnDeath(ctx *EffectContext) error {
	ctx.Engine.State.Players[ctx.OpponentID].GainElements(map[string]int{model.ElementArcane: 2})
	return nil
}

type Card1121016FireDancer struct{ AlwaysActive }

func (Card1121016FireDancer) ID() string   { return "1121016" }
func (Card1121016FireDancer) Name() string { return "舞火者" }
func (Card1121016FireDancer) OnEnter(ctx *EffectContext) error {
	grantFireNegativeStatusImmunity(ctx)
	return nil
}
func (Card1121016FireDancer) OnDeath(ctx *EffectContext) error {
	grantFireNegativeStatusImmunity(ctx)
	return nil
}

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

type Card1221016IceSpikeFortress struct{ AlwaysActive }

func (Card1221016IceSpikeFortress) ID() string   { return "1221016" }
func (Card1221016IceSpikeFortress) Name() string { return "冰刺堡垒" }
func (Card1221016IceSpikeFortress) OnDamaged(ctx *EffectContext) error {
	if ctx.ExtraData == nil || ctx.ExtraData["attacker"] == ctx.PlayerID {
		return nil
	}
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card.Position != nil && ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "ice_spike_fortress",
		"冰刺堡垒:选择法力范围内1个敌人冻结1，若已冻结则造成1点伤害", candidates, 1, 1,
		func(selected []string) {
			target := selectedUnitFromCandidates(ctx.Engine, selected, candidates)
			if target == nil {
				return
			}
			if target.Statuses[StatusFreeze] > 0 {
				ctx.Engine.dealDamageWithExtra(target, 1, target.OwnerID, map[string]any{"damage_source": "effect", "attacker": ctx.PlayerID})
				return
			}
			ctx.Engine.addStatus(target, StatusFreeze, 1)
		})
	return nil
}

type Card1321016ThunderGolem struct{ AlwaysActive }

func (Card1321016ThunderGolem) ID() string   { return "1321016" }
func (Card1321016ThunderGolem) Name() string { return "雷傀儡" }
func (Card1321016ThunderGolem) OnDeath(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.OpponentID]
	if len(ps.Hand) == 0 {
		return nil
	}
	idx := rand.Intn(len(ps.Hand))
	card := ps.Hand[idx]
	ps.Hand = append(ps.Hand[:idx], ps.Hand[idx+1:]...)
	ps.Graveyard = append(ps.Graveyard, card)
	ctx.Engine.emit(GameEvent{Type: "discard", Player: ctx.OpponentID, Data: map[string]any{"card": cardToInfo(card)}})
	return nil
}

type Card1421016Scavenger struct{ AlwaysActive }

func (Card1421016Scavenger) ID() string   { return "1421016" }
func (Card1421016Scavenger) Name() string { return "食腐者" }
func (Card1421016Scavenger) OnDamaged(ctx *EffectContext) error {
	if ctx.ExtraData == nil || ctx.Target == nil || ctx.Target == ctx.Source {
		return nil
	}
	damagedPlayer, _ := ctx.ExtraData["damaged_player"].(int)
	attacker, hasAttacker := ctx.ExtraData["attacker"].(int)
	if damagedPlayer == ctx.PlayerID && hasAttacker && attacker != ctx.PlayerID {
		ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementEarth: 2})
	}
	return nil
}

type Card1521016SoloCityDefender struct{ AlwaysActive }

func (Card1521016SoloCityDefender) ID() string   { return "1521016" }
func (Card1521016SoloCityDefender) Name() string { return "索洛城的坚守者" }

type Card1621016VengefulDead struct{ AlwaysActive }

func (Card1621016VengefulDead) ID() string   { return "1621016" }
func (Card1621016VengefulDead) Name() string { return "复仇死者" }
func (Card1621016VengefulDead) OnDeath(ctx *EffectContext) error {
	sourcePlayer := ctx.Source.Statuses["lethal_source_player"] - 1
	if sourcePlayer < 0 || sourcePlayer >= len(ctx.Engine.State.Players) {
		return nil
	}
	hero := ctx.Engine.State.Players[sourcePlayer].Hero
	if hero != nil {
		ctx.Engine.dealDamageWithExtra(hero, 2, sourcePlayer, map[string]any{"damage_source": "effect", "attacker": ctx.PlayerID})
	}
	return nil
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

func (m markerEquipment) ID() string   { return m.id }
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
		return nil
	}
	if ctx.Source.Statuses[m.counter] <= 0 {
		return nil
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
				target.CurrentLife++
				target.Statuses["max_life_bonus"]++
			}
			ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementLight: 2})
		})
		return nil
	}}
}

func newBurier() CardBehavior {
	canUse := func(ctx *EffectContext) bool {
		return len(ctx.Engine.friendlyDeckCards(ctx.PlayerID, func(card *CardInstance) bool {
			return card.Card.IsCompanion() && card.Card.Category == model.ElementShadow
		})) > 0
	}
	return markerEquipment{id: "2621014", name: "埋葬者", counter: burierCounter, counters: 3, canUse: canUse, effect: func(ctx *EffectContext) error {
		candidates := ctx.Engine.friendlyDeckCards(ctx.PlayerID, func(card *CardInstance) bool {
			return card.Card.IsCompanion() && card.Card.Category == model.ElementShadow
		})
		ctx.Engine.SetPendingAction(ctx.PlayerID, "burier", "埋葬者:选择1张暗影伙伴送去弃牌堆", candidates, 1, 1, func(selected []string) {
			card := removeSelectedDeckCard(ctx.Engine, ctx.PlayerID, selected, candidates)
			if card != nil {
				ps := ctx.Engine.State.Players[ctx.PlayerID]
				ps.Graveyard = append(ps.Graveyard, card)
				ctx.Engine.emit(GameEvent{Type: "discard", Player: ctx.PlayerID, Data: map[string]any{"card": cardToInfo(card)}})
			}
			ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementShadow: 2})
		})
		return nil
	}}
}

func removeSelectedDeckCard(e *Engine, playerID int, selected []string, candidates []map[string]any) *CardInstance {
	id := firstSelected(selected)
	if id == "" || !candidateContains(candidates, id) {
		return nil
	}
	ps := e.State.Players[playerID]
	for i, card := range ps.Deck {
		if card != nil && card.InstanceID == id {
			ps.Deck = append(ps.Deck[:i], ps.Deck[i+1:]...)
			e.shuffleDeck(playerID)
			return card
		}
	}
	return nil
}

func candidateContains(candidates []map[string]any, instanceID string) bool {
	for _, candidate := range candidates {
		if candidate["instance_id"] == instanceID {
			return true
		}
	}
	return false
}
