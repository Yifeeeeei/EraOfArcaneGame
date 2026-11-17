package bot

import (
	"eraofarcane/game"
	"sort"
)

type income struct {
	source  *card
	gain    map[string]int
	move    game.ActionMessage
	penalty int
}

func incomes(s state, opts map[string]game.DecisionCard, cost map[string]int) []income {
	var result []income
	needAir := cost["气"] > s.You.Elements["气"]+s.You.Elements["无"]
	for _, c := range s.You.field() {
		if c == nil {
			continue
		}
		o := opts[c.ID]
		gain := o.ConsumeGain
		if gain == nil && o.CanConsume {
			gain = c.Gain
		}
		move := action("consume", map[string]any{"instance_id": c.ID})
		penalty := 2
		if c.Type == "人物" {
			penalty = 1
		}
		if c.Attack > 0 && c.Position != nil && c.Position.Row == s.You.front() {
			penalty += 12
		}
		if c.Number == "1121001" {
			penalty += 4
			if c.Life <= 1 {
				penalty += 6
			}
		}
		if c.Number == "1121108" && needAir && o.CanConsume && o.CanPerTurn {
			gain = map[string]int{"气": 1}
			move = action("use_ability", map[string]any{"instance_id": c.ID, "ability_type": "per_turn"})
		}
		safeConversion := c.Life >= 4 || c.Life >= 2 && s.Opponent.Hero != nil && s.Opponent.Hero.Life <= 2
		if c.Number == "4111002" && needAir && o.CanConsume && o.CanPerTurn && safeConversion && c.Statuses["点燃"] == 0 {
			gain = plus(c.Gain, map[string]int{"火": -1, "无": 1})
			move = action("use_ability", map[string]any{"instance_id": c.ID, "ability_type": "per_turn"})
			penalty += 5
		}
		if o.CanConsume && sum(gain) > 0 {
			if c.Type == "人物" && s.You.ID == s.First && s.Number == 1 {
				limit := (sum(c.Gain) + 1) / 2
				reduced := map[string]int{}
				if needAir && gain["无"] > 0 {
					reduced["无"] = 1
					limit--
				}
				for _, e := range []string{"火", "无", "气", "地", "光", "水", "暗"} {
					n := min(limit, gain[e]-reduced[e])
					if n > 0 {
						reduced[e] += n
						limit -= n
					}
				}
				gain = reduced
				if move.Action == "consume" {
					move.Data["gain"] = gain
					move = action(move.Action, move.Data)
				}
			}
			result = append(result, income{c, gain, move, penalty})
		}
		if o.CanPerTurn && c.Number == "2121001" && c.Statuses["火焰标记物"] > 0 {
			result = append(result, income{c, map[string]int{"火": 1}, action("use_ability", map[string]any{"instance_id": c.ID, "ability_type": "per_turn"}), 4})
		}
		if o.CanPerTurn && c.Number == "2121014" && !c.Horizontal && c.Statuses["眩晕"] == 0 && c.Statuses["火匣子标记物"] > 0 {
			result = append(result, income{c, map[string]int{"火": 2}, action("use_ability", map[string]any{"instance_id": c.ID, "ability_type": "per_turn"}), 4})
		}
	}
	return result
}

func (p *Policy) fund(s state, opts map[string]game.DecisionCard, move game.ActionMessage, cost map[string]int) (game.ActionMessage, bool) {
	if pay, ok := payment(s, s.You.Elements, cost); ok {
		move.Data["payment"] = pay
		move = action(move.Action, move.Data)
		return move, !p.blocked[Fingerprint(move)]
	}
	choices := incomes(s, opts, cost)
	sort.SliceStable(choices, func(i, j int) bool {
		_, a := payment(s, plus(s.You.Elements, choices[i].gain), cost)
		_, b := payment(s, plus(s.You.Elements, choices[j].gain), cost)
		if a != b {
			return a
		}
		return choices[i].penalty < choices[j].penalty
	})
	available := plus(s.You.Elements, nil)
	var first *income
	for i := range choices {
		c := &choices[i]
		if p.blocked[Fingerprint(c.move)] {
			continue
		}
		if first == nil {
			first = c
		}
		available = plus(available, c.gain)
	}
	if _, ok := payment(s, available, cost); !ok || first == nil {
		return game.ActionMessage{}, false
	}
	return first.move, true
}

func resources(s state) (fire, independent int) {
	for _, c := range s.You.field() {
		if c == nil {
			continue
		}
		if c.Type != "人物" && c.Life > 0 && c.Statuses["点燃"] > 0 && c.Life <= 1 {
			continue
		}
		n := c.Gain["火"] + c.Gain["无"]
		if c.Number == "2121001" && c.Statuses["火焰标记物"] >= 2 {
			n++
		}
		if c.Number == "2121014" && c.Statuses["火匣子标记物"] >= 2 {
			n += 2
		}
		if c.Statuses["冻结"] > 0 || c.Statuses["眩晕"] > 0 {
			n = 0
		}
		fire += n
		if c.Type != "人物" && n > 0 {
			independent++
		}
	}
	return
}

func keepValue(c *card, s state) int {
	if c == nil {
		return -100
	}
	fire, _ := resources(s)
	if c.Number == "1121001" {
		return 80
	}
	if c.Number == "2121001" {
		for _, old := range s.You.Equipment {
			if old != nil && old.Subtype == "神器" && old.Statuses["火焰标记物"] > 0 {
				return 12
			}
		}
		return 65
	}
	if c.Number == "1111001" {
		return 24
	}
	if c.Number == "2121111" {
		if s.Opponent.Hero != nil && s.Opponent.Hero.Life <= 4 {
			return 95
		}
		return 28
	}
	if c.Number == "1111103" {
		return 15
	}
	if c.Number == "2021107" {
		return 48
	}
	if sum(c.Gain) > 0 {
		if fire < 6 {
			return 70 - sum(c.Cost)*5
		}
		return 25
	}
	if c.Scroll {
		return 55
	}
	return 30
}

// Defensive payment never invokes a tap ability. Excess overexerted load is
// consumed by the engine for this payment only.
func defensivePayment(s state, opts map[string]game.DecisionCard, cost map[string]int, allowHero bool) (map[string]int, []string, bool) {
	if pay, ok := payment(s, s.You.Elements, cost); ok {
		return pay, nil, true
	}
	var sources []*card
	for _, c := range s.You.field() {
		if c != nil && opts[c.ID].CanOverexert && (allowHero || c.Type != "人物") {
			sources = append(sources, c)
		}
	}
	sort.SliceStable(sources, func(i, j int) bool {
		if (sources[i].Type == "人物") != (sources[j].Type == "人物") {
			return sources[j].Type == "人物"
		}
		return sum(sources[i].Gain) < sum(sources[j].Gain)
	})
	available := plus(s.You.Elements, nil)
	var ids []string
	for _, c := range sources {
		available = plus(available, c.Gain)
		ids = append(ids, c.ID)
		if pay, ok := payment(s, available, cost); ok {
			return pay, ids, true
		}
	}
	return nil, nil, false
}
