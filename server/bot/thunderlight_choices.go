package bot

import (
	"eraofarcane/game"
	"fmt"
)

func (p *Policy) thunderChoiceValue(raw map[string]any, c card, s state, score int) int {
	if own := s.You.find(c.ID); own != nil {
		c = *own
	}
	switch s.Pending.Type {
	case "five_color_coral_load":
		if c.ID == "气" || c.ID == "光" {
			return 1000
		}
		return -1000
	case "search_deck":
		items, armor := thunderItems(s)
		_, light, arcane, _ := thunderResources(s)
		switch c.Number {
		case "2321105":
			if armor == 0 && items >= 2 {
				return 1000
			}
			if armor == 0 {
				return 650
			}
			if s.You.Duplicates {
				return 550
			}
			return 120
		case "2021104":
			if light+arcane < 2 {
				return 900
			}
			return 630
		case "2321104", "2321101":
			return 620
		case "2021105":
			if armor > 0 {
				return 700
			}
			return 100
		}
	case "lost_silverleaf_discard", "discard", "discard_hand":
		v := -p.keepValue(&c, s)
		if c.Number == "2321103" {
			v += 150
		}
		if c.Number == "1321102" && s.You.Elements["气"] > 0 && len(s.You.units()) < 9 {
			v += 80
		}
		return v
	case "speckled_sparrow_discard_summon":
		pos := asMap(raw["position"])
		return 100 - number(pos["row"])*30
	case "spell_square_area":
		var col, row int
		if _, err := fmt.Sscanf(c.ID, "square:%d:%d", &col, &row); err == nil {
			value := 0
			for _, u := range s.Opponent.units() {
				if u.Position.Col >= col && u.Position.Col <= col+1 && u.Position.Row >= row && u.Position.Row <= row+1 {
					value += thunderUnitValue(u, s, 1)
				}
			}
			return value
		}
	}
	return score
}

func (p *Policy) thunderWarriorChoice(s state) Decision {
	front := 0
	for _, c := range s.You.units() {
		if c.Type != "人物" && c.Position.Row <= 1 && c.Life >= 2 {
			front++
		}
	}
	_, light, arcane, _ := thunderResources(s)
	selected := []string{}
	for i := 0; i < s.Pending.Min; i++ {
		reward := "air"
		if i == 0 && front < 2 {
			reward = "life"
		} else if light+arcane < 2 {
			reward = "light"
			light++
		}
		id := fmt.Sprintf("%s#%d", reward, i)
		selected = append(selected, id)
	}
	return Decision{action("resolve_action", map[string]any{"decision_id": s.Pending.ID, "selected": selected}), "雷光战士按前排需求和气光缺口选择奖励"}
}

// Only public ready spells and visible resource sources contribute to this
// estimate. Hidden scrolls, reactions and future draws are not inspected.
func publicDefensePower(s state) int {
	available := plus(s.Opponent.Elements, nil)
	for _, c := range s.Opponent.field() {
		if c != nil && !c.Horizontal && c.Statuses["眩晕"] == 0 {
			available = plus(available, c.Gain)
		}
	}
	var spells []*card
	for _, c := range s.Opponent.spells() {
		if c != nil && c.CanDefend && !c.Horizontal && c.Statuses["冷却"] == 0 && c.Statuses["眩晕"] == 0 {
			spells = append(spells, c)
		}
	}
	if len(spells) > 10 {
		spells = spells[:10]
	}
	best := 0
	for mask := 1; mask < 1<<len(spells); mask++ {
		power := 0
		cost := map[string]int{}
		for i, c := range spells {
			if mask&(1<<i) != 0 {
				power += c.DefensePower
				cost = plus(cost, c.DefenseCost)
			}
		}
		if power <= best {
			continue
		}
		if _, ok := game.DecisionPayment(available, cost, s.Opponent.Strict, s.Opponent.Light, s.Opponent.LightCost); ok {
			best = power
		}
	}
	return best
}

func (p *Policy) thunderBoostMoves(s state, opts map[string]game.DecisionCard, main *card, target game.DecisionTarget, data map[string]any, score int, add func(game.ActionMessage, map[string]int, int, string, map[string]int, []string)) {
	defense := publicDefensePower(s)
	if target.Power > defense || s.Suppressed {
		return
	}
	var boosts []*card
	for _, c := range s.You.spells() {
		if c != nil && c.ID != main.ID && opts[c.ID].CanAttackBoost {
			boosts = append(boosts, c)
		}
	}
	if len(boosts) > 8 {
		boosts = boosts[:8]
	}
	for i, a := range boosts {
		for j := -1; j < i; j++ {
			ids := []string{a.ID}
			cost := plus(target.Cost, opts[a.ID].AttackBoostCost)
			power := target.Power + opts[a.ID].AttackBoostPower
			if j >= 0 {
				b := boosts[j]
				ids = append(ids, b.ID)
				cost = plus(cost, opts[b.ID].AttackBoostCost)
				power += opts[b.ID].AttackBoostPower
			}
			if power <= defense {
				continue
			}
			copy := map[string]any{}
			for k, v := range data {
				copy[k] = v
			}
			copy["boost_ids"] = ids
			add(action("cast_spell", copy), cost, score-len(ids)*15-max(0, sum(cost)-sum(target.Cost))*3, "有限强化：提升关键一发的威力，控制强化数量", nil, nil)
		}
	}
}
