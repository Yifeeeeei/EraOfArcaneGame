package bot

import (
	"eraofarcane/game"
	"sort"
)

func thunderPay(s state, available, cost map[string]int) (map[string]int, bool) {
	pay, ok := payment(s, available, cost)
	if !ok {
		return nil, false
	}
	if !s.You.Light && !s.You.LightCost {
		n := min(cost["无"], min(max(0, pay["光"]-cost["光"]), max(0, available["气"]-pay["气"])))
		pay["光"] -= n
		pay["气"] += n
	}
	return pay, true
}

func thunderIncomes(s state, opts map[string]game.DecisionCard, cost map[string]int) []income {
	choices := incomes(s, opts, cost)
	for _, c := range s.You.Hand {
		if c.Number == "2321103" {
			choices = append(choices, income{c, map[string]int{"气": 1}, action("use_item", map[string]any{"instance_id": c.ID}), 1})
		}
	}
	return choices
}

func thunderAvailable(s state, opts map[string]game.DecisionCard, exclude string) map[string]int {
	result := plus(s.You.Elements, nil)
	for _, i := range thunderIncomes(s, opts, nil) {
		if i.source.ID != exclude {
			result = plus(result, i.gain)
		}
	}
	return result
}

func (p *Policy) thunderFund(s state, opts map[string]game.DecisionCard, move game.ActionMessage, cost map[string]int) (game.ActionMessage, bool) {
	if c := s.You.find(text(move.Data["instance_id"])); c != nil && c.Number == "3021101" && sum(cost) == cost["无"] {
		if s.You.Elements["无"] >= sum(cost) {
			move.Data["payment"] = map[string]int{"无": sum(cost)}
			move = action(move.Action, move.Data)
			return move, !p.blocked[Fingerprint(move)]
		}
		if move.Action == "cast_spell" && thunderAvailable(s, opts, "")["光"] >= sum(cost) {
			cost = map[string]int{"光": sum(cost)}
		}
	}
	if pay, ok := thunderPay(s, s.You.Elements, cost); ok {
		move.Data["payment"] = pay
		move = action(move.Action, move.Data)
		return move, !p.blocked[Fingerprint(move)]
	}
	choices := thunderIncomes(s, opts, cost)
	score := func(i income) int {
		v := -i.penalty
		if _, ok := thunderPay(s, plus(s.You.Elements, i.gain), cost); ok {
			v += 200
		}
		for _, e := range []string{"光", "气", "无"} {
			v += min(max(0, cost[e]-s.You.Elements[e]), i.gain[e]) * 12
		}
		return v
	}
	sort.SliceStable(choices, func(i, j int) bool { return score(choices[i]) > score(choices[j]) })
	for _, i := range choices {
		if !p.blocked[Fingerprint(i.move)] {
			return i.move, true
		}
	}
	return game.ActionMessage{}, false
}

func thunderPlacement(s state, c *card, devour []string) (game.Position, bool) {
	removed := map[string]bool{}
	for _, id := range devour {
		removed[id] = true
	}
	front := c.Taunt || c.Number == "1321002" || c.Number == "1321111" || c.Number == "1021016" || c.MaxLife >= 2 && sum(c.Gain) <= 1
	best := -10000
	var pos game.Position
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if old := s.You.Units[col][row]; old != nil && !removed[old.ID] {
				continue
			}
			score := row * 30
			if front {
				score = (2 - row) * 45
			}
			if c.Number == "1321006" {
				for _, old := range s.You.units() {
					if old.Number == "1321006" {
						score += (abs(col-old.Position.Col) + abs(row-old.Position.Row)) * 12
					}
				}
			}
			if score > best {
				best = score
				pos = game.Position{Col: col, Row: row}
			}
		}
	}
	return pos, best > -10000
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (p *Policy) thunderDevour(s state, opts map[string]game.DecisionCard, c *card) []string {
	var candidates []*card
	for _, id := range opts[c.ID].DevourCandidates {
		u := s.You.find(id)
		if u != nil && u.Gain["气"] > 0 && u.Number != "1321006" && u.Number != "1321010" && u.Number != "1321111" {
			candidates = append(candidates, u)
		}
	}
	sort.SliceStable(candidates, func(i, j int) bool { return p.keepValue(candidates[i], s) < p.keepValue(candidates[j], s) })
	if len(candidates) > 12 {
		candidates = candidates[:12]
	}
	items, _ := thunderItems(s)
	air, light, arcane, _ := thunderResources(s)
	best := 10000
	var result []string
	for mask := 1; mask < 1<<len(candidates); mask++ {
		var ids []string
		load, lostLight, lostItems, value := 0, 0, 0, 0
		for i, u := range candidates {
			if mask&(1<<i) == 0 {
				continue
			}
			ids = append(ids, u.ID)
			load += u.Gain["气"]
			lostLight += u.Gain["光"]
			value += p.keepValue(u, s)
			if u.Equipment && u.Gain["光"] > 0 {
				lostItems++
			}
		}
		if len(ids) > 4 || load < 3 || air+arcane-load+2 < 4 || light+arcane-lostLight < 1 || items >= 3 && items-lostItems < 3 || value > 220 {
			continue
		}
		if value < best {
			best = value
			result = ids
		}
	}
	return result
}

func (p *Policy) keepValue(c *card, s state) int {
	if !p.isThunderlight() {
		return keepValue(c, s)
	}
	if c == nil {
		return -100
	}
	air, light, arcane, independent := thunderResources(s)
	v := 35 + sum(c.Gain)*9
	if independent < 2 || air+arcane < 6 {
		v += 15
	}
	if c.Gain["光"] > 0 && light+arcane < 2 {
		v += 40
	}
	switch c.Number {
	case "1321006":
		v = 155
	case "1321010":
		v = 110
	case "2321105":
		v = 130
	case "2021104", "2321101", "2321104":
		v = 95
	case "2021105":
		v = 45
		if ownCount(s, "2321105") > 0 {
			v = 100
		}
	case "1321111":
		v = 90
	case "1021016":
		v = 30
		if noEquipment(s) {
			v = 100
		}
	case "1311001":
		v = 25
		if air+arcane >= 7 {
			v = 65
		}
	case "1321002":
		v = 25
	case "1321102":
		v = 28
	case "2321103":
		v = 10
	case "2021101":
		v = 75
	case "1321110":
		if deckCount(s, "2021101") > 0 {
			v += 15
		}
	}
	if c.Temporary {
		v -= 35
	}
	return v
}
