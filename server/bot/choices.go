package bot

import (
	"encoding/json"
	"eraofarcane/game"
	"sort"
)

func (p *Policy) choose(s state, opts map[string]game.DecisionCard) Decision {
	pa := s.Pending
	if pa == nil {
		return Decision{}
	}
	var candidates []map[string]any
	for _, c := range pa.Candidates {
		if c["can_select"] != false && text(c["instance_id"]) != "" {
			candidates = append(candidates, c)
		}
	}
	value := func(c map[string]any) int {
		var ci card
		raw, _ := json.Marshal(c)
		_ = json.Unmarshal(raw, &ci)
		score := keepValue(&ci, s)
		if text(c["side"]) == "enemy" || number(c["owner"]) == s.Opponent.ID && c["owner"] != nil {
			score = targetScore(&ci, s, 1, false, false)
		}
		if ci.Type == "人物" {
			score += 100
		}
		if pa.Type == "swordsmanship_teacher_buff" && ci.Position != nil {
			score += 100 - ci.Position.Row*40 + ci.Life*3
		}
		if pa.Type == "discard" || pa.Type == "discard_hand" {
			score = -keepValue(&ci, s)
		}
		if pa.Type == "insight_eye_destroy_equipment" {
			if !ci.SetCounter {
				score = -1000
			}
		}
		return score
	}
	sort.SliceStable(candidates, func(i, j int) bool { return value(candidates[i]) > value(candidates[j]) })
	count := pa.Min
	switch pa.Type {
	case "embers_free_cast_target", "fire_arrow_damage", "nur_eye_fire_damage", "swordsmanship_teacher_buff":
		count = max(count, 1)
	case "spark_moth_reveal":
		count = min(pa.Max, len(candidates))
	case "counter_trigger":
		source := asMap(pa.Context["source"])
		if number(pa.Context["source_player"]) == s.Opponent.ID && source != nil {
			life := number(source["current_life"])
			burn := number(asMap(source["statuses"])["点燃"])
			if text(source["type"]) == "人物" || life == 1 && burn == 0 {
				count = 1
			}
		}
	}
	count = min(count, min(pa.Max, len(candidates)))
	// Retry different candidates after authoritative resyncs, retaining the
	// number/uniqueness contracts rather than replaying the same failed payload.
	for variant := 0; variant < max(1, len(candidates)); variant++ {
		selected := []string{}
		for i := 0; i < count; i++ {
			selected = append(selected, text(candidates[(variant+i)%len(candidates)]["instance_id"]))
		}
		if text(pa.UI["unique_by"]) == "side" {
			seen := map[string]bool{}
			selected = nil
			for _, c := range candidates {
				side := text(c["side"])
				if !seen[side] {
					selected = append(selected, text(c["instance_id"]))
					seen[side] = true
				}
				if len(selected) == count {
					break
				}
			}
		}
		data := map[string]any{"selected": selected, "decision_id": pa.ID}
		switch text(pa.UI["kind"]) {
		case "reorder":
			ids := []string{}
			for _, c := range candidates {
				ids = append(ids, text(c["instance_id"]))
			}
			data["top_order"] = ids
			data["bottom_order"] = []string{}
		case "allocation":
			amounts := map[string]int{}
			remaining := number(pa.UI["total"])
			for i, id := range selected {
				n := 1
				if i == 0 {
					n = remaining - len(selected) + 1
				}
				amounts[id] = n
				remaining -= n
			}
			if len(selected) > 0 {
				data["allocations"] = amounts
			}
		case "load_distribution":
			distribution := map[string]any{}
			for _, id := range selected {
				distribution[id] = map[string]int{"火": 1}
			}
			data["load_distribution"] = distribution
		case "resource_distribution":
			data["payments"] = []any{}
			for _, id := range selected {
				for _, c := range candidates {
					if text(c["instance_id"]) != id {
						continue
					}
					load := elements(c["payable_load"])
					if sum(load) > 0 {
						for _, e := range []string{"火", "无", "地", "气", "水", "光", "暗"} {
							if load[e] > 0 {
								data["payments"] = append(data["payments"].([]any), map[string]any{"instance_id": id, "amount": 1, "mode": "load", "element": e})
								break
							}
						}
					} else {
						data["payments"] = append(data["payments"].([]any), map[string]any{"instance_id": id, "amount": 1, "mode": "life"})
					}
				}
			}
		}
		cost := pa.Cost
		candidateCost := map[string]int{}
		hasCost := false
		for _, id := range selected {
			for _, c := range candidates {
				if text(c["instance_id"]) == id && c["payment_cost"] != nil {
					hasCost = true
					candidateCost = plus(candidateCost, elements(c["payment_cost"]))
				}
			}
		}
		if hasCost {
			cost = candidateCost
		}
		if len(selected) == 0 {
			cost = nil
		}
		pay, ids, ok := defensivePayment(s, opts, cost, pa.Min > 0)
		if ok {
			data["payment"] = pay
			if pa.Overexert {
				data["overexert_ids"] = ids
			} else if len(ids) > 0 {
				continue
			}
		} else {
			continue
		}
		move := action("resolve_action", data)
		if !p.blocked[Fingerprint(move)] {
			return Decision{move, "处理选择：已知效果按策略，其他效果使用合法默认项"}
		}
	}
	if pa.Min == 0 {
		move := action("resolve_action", map[string]any{"selected": []string{}, "decision_id": pa.ID})
		if !p.blocked[Fingerprint(move)] {
			return Decision{move, "放弃无收益或无法支付的可选效果"}
		}
	}
	// A broken/unsupported mandatory callback must never spin forever. This is
	// a visible terminal failure, not a fabricated successful rules resolution.
	return Decision{action("surrender", nil), "AI 异常退出：强制选择的默认方案均被规则拒绝，请保留房间日志"}
}

func (p *Policy) defend(s state, opts map[string]game.DecisionCard) Decision {
	decline := Decision{action("no_defend", nil), "快攻保留进攻技能：此次不防御"}
	if s.Spell == nil || s.You.Hero == nil {
		return decline
	}
	damage := 0
	if s.Spell.Skill != nil {
		damage = max(0, s.Spell.Skill.Attack)
	}
	// Pending spell cards expose printed attack as well as current_attack;
	// state.card uses the latter. Inspect the permitted envelope for that field.
	if damage == 0 && s.Spell.Skill != nil {
		damage = max(0, s.Spell.Skill.PrintedAttack)
	}
	threatensHero := false
	for _, pos := range s.Spell.Affected {
		if pos.Owner == s.You.ID && s.You.Hero.Position != nil && pos.Col == s.You.Hero.Position.Col && pos.Row == s.You.Hero.Position.Row {
			threatensHero = true
		}
	}
	if !threatensHero {
		return decline
	}
	remaining := s.You.Hero.Life - max(0, damage-s.You.Shield)
	if remaining > 1 {
		return decline
	}
	var spells []*card
	for _, c := range append(s.You.spells(), s.You.Hand...) {
		if c != nil && opts[c.ID].CanDefend && !c.Sorcery {
			spells = append(spells, c)
		}
	}
	if len(spells) > 12 {
		spells = spells[:12]
	}
	best := 100000
	var result Decision
	for mask := 1; mask < (1 << len(spells)); mask++ {
		power := 0
		cost := map[string]int{}
		ids := []string{}
		scrolls := []string{}
		penalty := 0
		for i, c := range spells {
			if mask&(1<<i) == 0 {
				continue
			}
			power += c.DefensePower
			cost = plus(cost, opts[c.ID].DefenseCost)
			penalty += 5
			if c.Number == "3121109" {
				penalty += 8
			}
			if c.Scroll {
				scrolls = append(scrolls, c.ID)
				penalty += 8
			} else {
				ids = append(ids, c.ID)
			}
		}
		if power < s.Spell.Power {
			continue
		}
		pay, overexert, ok := defensivePayment(s, opts, cost, remaining <= 0)
		if !ok {
			continue
		}
		penalty += sum(cost)*2 + len(overexert)*8
		for _, id := range overexert {
			if id == s.You.Hero.ID {
				penalty += 30
			}
		}
		if remaining > 0 && penalty > 18 {
			continue
		}
		move := action("defend", map[string]any{"skill_ids": ids, "scroll_ids": scrolls, "boost_ids": []string{}, "overexert_ids": overexert, "payment": pay})
		if penalty < best && !p.blocked[Fingerprint(move)] {
			best = penalty
			result = Decision{move, "防御人物致命／濒死伤害，选择最低机会成本方案"}
		}
	}
	if result.Action.Action != "" {
		return result
	}
	return decline
}
