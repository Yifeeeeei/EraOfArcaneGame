package game

import "fmt"

// SpellTapCostBehavior declares an additional tap payment for every use of a
// spell, including defense and boosts. Eligibility and bonus queries are pure.
type SpellTapCostBehavior interface {
	HasActiveSpellTapCost(*CardInstance) bool
	SpellTapCostBonus(*EffectContext, *CardInstance) int
}

type spellTapCost struct {
	source, target *CardInstance
	bonus          int
}

func (e *Engine) prepareSpellTapCosts(player int, sources []*CardInstance, action ActionMessage, retry func(ActionMessage) error) ([]spellTapCost, bool, error) {
	var required []*CardInstance
	for _, source := range sources {
		if b, ok := cardBehavior(source).(SpellTapCostBehavior); ok && b.HasActiveSpellTapCost(source) {
			required = append(required, source)
		}
	}
	if len(required) == 0 {
		return nil, false, nil
	}
	selected := make(map[string]string)
	if raw, ok := action.Data["consume_ids"].(map[string]any); ok {
		for id, value := range raw {
			selected[id], _ = value.(string)
		}
	}
	if len(required) == 1 && selected[required[0].InstanceID] == "" {
		selected[required[0].InstanceID], _ = action.Data["consume_id"].(string)
		if selected[required[0].InstanceID] == "" && action.Action == "use_item" {
			selected[required[0].InstanceID], _ = action.Data["sacrifice_id"].(string)
		}
	}
	reserved := map[string]bool{}
	for _, id := range stringsFromAnySlice(anySliceFromData(action.Data, "overexert_ids")) {
		reserved[id] = true
	}
	// A unit cannot also pay an independent sacrifice cost in the same action.
	if action.Action != "use_item" || len(required) != 1 || action.Data["consume_id"] != nil || action.Data["consume_ids"] != nil {
		if id, _ := action.Data["sacrifice_id"].(string); id != "" {
			reserved[id] = true
		}
	}
	var plans []spellTapCost
	for _, source := range required {
		id := selected[source.InstanceID]
		if id == "" {
			continue
		}
		target := e.findCardOnField(e.State.Players[player], id)
		bonus := cardBehavior(source).(SpellTapCostBehavior).SpellTapCostBonus(e.skillContext(player, source), target)
		if reserved[id] || bonus <= 0 {
			return nil, false, fmt.Errorf("%s 的消耗费用需要另选可横置的友方伙伴，不能重复透支或消耗同一张牌", source.Card.Name)
		}
		reserved[id] = true
		plans = append(plans, spellTapCost{source, target, bonus})
	}
	for _, source := range required {
		if selected[source.InstanceID] != "" {
			continue
		}
		var candidates []map[string]any
		for _, target := range e.getAllFieldCards(e.State.Players[player]) {
			if !reserved[target.InstanceID] && cardBehavior(source).(SpellTapCostBehavior).SpellTapCostBonus(e.skillContext(player, source), target) > 0 {
				candidates = append(candidates, candidateInfo(target, "unit", "own"))
			}
		}
		if len(candidates) < len(required)-len(plans) {
			return nil, false, fmt.Errorf("%s 没有足够的伙伴支付消耗费用", source.Card.Name)
		}
		e.SetPendingActionWithError(player, "spell_tap_cost", source.Card.Name+":选择要消耗的友方伙伴（保留在场）", candidates, 1, 1, nil, false,
			func(ids []string, _ map[string]any) error {
				data := cloneExtraData(action.Data)
				choices := make(map[string]any)
				for id, target := range selected {
					if target != "" {
						choices[id] = target
					}
				}
				choices[source.InstanceID] = firstSelected(ids)
				data["consume_ids"] = choices
				return retry(ActionMessage{Action: action.Action, Data: data})
			})
		return nil, true, nil
	}
	return plans, false, nil
}

// Commit all validated tap payments without emitting events or running effects.
// Grant bonuses only after the whole batch has been turned horizontal.
func (e *Engine) commitSpellTapCosts(plans []spellTapCost) {
	for _, plan := range plans {
		plan.target.IsHorizontal = true
	}
	for _, plan := range plans {
		plan.source.spellTapPowerBonus = plan.bonus
	}
}

func (e *Engine) resolveSpellTapCosts(player int, plans []spellTapCost, after func()) {
	if len(plans) == 0 {
		after()
		return
	}
	steps := make([]func(), 0, len(plans)+1)
	for _, plan := range plans {
		steps = append(steps, func() {
			e.resolveConsumedCardForEffect(player, plan.target, nil, plan.source.Card.Number)
		})
	}
	steps = append(steps, after)
	e.runResolution("spell tap payment", steps...)
}
