package bot

import (
	"encoding/json"
	"eraofarcane/game"
	"fmt"
	"sort"
)

type intent struct {
	move   game.ActionMessage
	cost   map[string]int
	score  int
	reason string
}

func (p *Policy) Next(observation game.DecisionObservation) (Decision, bool, error) {
	var s state
	if err := json.Unmarshal(observation.State, &s); err != nil {
		return Decision{}, false, err
	}
	s.IncomingBurn = observation.IncomingBurn
	if s.Phase == "game_over" {
		return Decision{}, false, nil
	}
	if s.DrawOffer == s.Opponent.ID {
		revision := s.Revision
		return Decision{game.ActionMessage{Action: "respond_draw_offer", Data: map[string]any{"accept": false}, ExpectedRevision: &revision}, "练习对手继续对战，婉拒和局请求"}, true, nil
	}
	if s.Number != p.turn {
		p.turn = s.Number
		p.steps = 0
		p.blocked = map[string]bool{}
	}
	var d Decision
	if s.Phase == "mulligan" {
		if s.Mulligan[s.You.ID] {
			return Decision{}, false, nil
		}
		good := 0
		cost := 0
		artifact := false
		for _, c := range s.You.Hand {
			if c == nil {
				continue
			}
			if c.Number == "1121001" {
				good += 2
				continue
			}
			if c.Number == "2121001" {
				if artifact {
					continue
				}
				artifact = true
			}
			if (sum(c.Gain) > 0 || c.Number == "2121001" || c.Number == "2121014") && sum(c.Cost) <= 2 {
				good++
				cost += sum(c.Cost)
			}
		}
		keep := good >= 2 && cost <= 5
		d = Decision{action("mulligan", map[string]any{"keep": keep}), fmt.Sprintf("起手资源评估：%d 个早期资源权重，保留=%t", good, keep)}
	} else {
		if s.Decision != s.You.ID {
			return Decision{}, false, nil
		}
		switch s.Phase {
		case "waiting_action":
			d = p.choose(s, observation.Cards)
		case "defense_window":
			d = p.defend(s, observation.Cards)
		case "main":
			d = p.main(s, observation.Cards)
		default:
			return Decision{}, false, nil
		}
	}
	if d.Action.Action == "" {
		return Decision{}, false, nil
	}
	p.steps++
	revision := s.Revision
	d.Action.ExpectedRevision = &revision
	return d, true, nil
}

func (p *Policy) main(s state, opts map[string]game.DecisionCard) Decision {
	var moves []intent
	add := func(move game.ActionMessage, cost map[string]int, score int, reason string) {
		moves = append(moves, intent{move, cost, score, reason})
	}
	fire, independent := resources(s)
	grow := fire < 6 || independent < 2
	first := ownTurn(s) <= 1
	known := 0
	learned := map[string]bool{}
	for _, c := range s.You.Skills {
		if c != nil {
			known++
			learned[c.Number] = true
		}
	}
	readySpells := 0
	for _, c := range s.You.spells() {
		if c != nil && opts[c.ID].CanCast && !c.Sorcery {
			readySpells++
		}
	}
	// Attack policy has one explicit opposing-card exception requested by the
	// author: remove Pegasus before wasting spells on its protected allies.
	pegasus := false
	for _, c := range s.Opponent.units() {
		if c.Number == "1511003" && c.Statuses["石化"] == 0 {
			pegasus = true
		}
	}
	for _, c := range s.You.field() {
		if c == nil {
			continue
		}
		for _, pos := range opts[c.ID].AttackTargets {
			target := s.Opponent.Units[pos.Col][pos.Row]
			if target == nil {
				continue
			}
			score := targetScore(target, s, c.Attack, false, false) + 35
			add(action("attack", map[string]any{"attacker_id": c.ID, "target_col": pos.Col, "target_row": pos.Row}), nil, score, "直接攻击：保留攻击机会，不提前消耗攻击者")
		}
	}
	for _, c := range append(s.You.spells(), s.You.Hand...) {
		if c == nil || !opts[c.ID].CanCast || (!c.Scroll && c.Type != "技能") {
			continue
		}
		o := opts[c.ID]
		if c.Sorcery {
			score := -1000
			switch c.Number {
			case "3121007":
				if readySpells >= 2 {
					score = 450
				}
			case "3121008":
				if readySpells >= 2 && fire >= 7 {
					score = 455
				}
			case "3021006":
				for _, e := range s.Opponent.Equipment {
					if e != nil && e.SetCounter {
						score = 290
					}
				}
			}
			if score > 0 {
				add(action("cast_spell", map[string]any{"instance_id": c.ID, "target_type": "none"}), o.UseCost, score, "先使用有收益的辅助咒术")
			}
			continue
		}
		for _, target := range o.Targets {
			unit := s.Opponent.find(target.InstanceID)
			if unit == nil {
				continue
			}
			score := targetScore(unit, s, target.Damage, pegasus, c.Number == "3121015")
			switch c.Number {
			case "3121109":
				score += 48
			case "3121003":
				score += 30
			case "3121002":
				score += 22
				if unit.Statuses["点燃"] > 0 {
					score += 12
				}
			case "3101001":
				score += 28
			case "3121105":
				score -= 100
				if s.You.Elements["火"] == 0 {
					continue
				}
			case "2121111":
				score -= 20
				if unit.Type != "人物" && unit.Life < 3 {
					score -= 80
				}
			}
			if c.Area == "front_row" {
				n := 0
				kills := 0
				for _, u := range s.Opponent.units() {
					if u.Position.Row == s.Opponent.front() {
						n++
						if u.Life <= target.Damage {
							kills++
						}
					}
				}
				score += kills * 35
				if n > 1 && kills == 0 {
					score -= 80
				}
			}
			if target.Damage == 0 && unit.Statuses["点燃"] > 0 {
				score -= 80
			}
			add(action(map[bool]string{true: "use_item", false: "cast_spell"}[c.Scroll], map[string]any{"instance_id": c.ID, "target_type": target.Type, "target_col": target.Position.Col, "target_row": target.Position.Row, "boost_ids": []string{}}), target.Cost, score, "进攻：人物优先，受保护时清理阻挡／珀伽索斯")
		}
	}
	for _, c := range s.You.Hand {
		if c == nil {
			continue
		}
		if c.Type == "伙伴" {
			score := 210
			if grow {
				score = 390
			}
			if first {
				score += 100
			}
			score -= sum(c.Cost) * 8
			if c.Number == "1121001" {
				score += 65
			}
			if c.Number == "1121111" {
				score += 12
			}
			if c.Number == "1121012" {
				score = 265
				if readySpells > 0 {
					score = 405
				}
			}
			if c.Number == "1121110" {
				score = 260
				if readySpells > 0 {
					score = 365
				}
			}
			if c.Number == "1111103" {
				score = 140
				if known >= 3 && len(s.You.Hand) <= 2 {
					score = 275
				}
			}
			if c.Number == "1121112" && sum(c.Cost) == 0 {
				score += 50
			}
			var devour []string
			if c.Number == "1111001" {
				score = 160
				var fodder []*card
				for _, u := range s.You.field() {
					if u != nil && u.Type != "人物" && u.Gain["火"] > 0 && u.Number != "1121012" && u.Attack <= 0 {
						fodder = append(fodder, u)
					}
				}
				sort.SliceStable(fodder, func(i, j int) bool { return fodder[i].Life < fodder[j].Life })
				load := 0
				for _, u := range fodder {
					devour = append(devour, u.ID)
					load += u.Gain["火"]
					if load >= 3 {
						break
					}
				}
				if load < 3 || fire-load+3 < 5 {
					continue
				}
				if readySpells >= 2 {
					score = 295
				}
			}
			pos, ok := placement(s, c, devour)
			if !ok {
				continue
			}
			if c.Number == "1021102" {
				score = 155
				if buffTarget(s, pos) != nil {
					score = 320
				} else {
					continue
				}
			}
			add(action("summon", map[string]any{"instance_id": c.ID, "col": pos.Col, "row": pos.Row, "devour_ids": devour}), c.Cost, score, "部署伙伴：建立独立负载／准备特殊效果")
		} else if c.Equipment {
			score := 180
			switch c.Number {
			case "2121006":
				score = 385
			case "2121001":
				score = 405
			case "2121014":
				score = 380
			case "2111002":
				score = 165
				if readySpells >= 2 {
					score = 270
				}
			case "2121004":
				score = 230
			}
			if !grow && c.Number != "2121004" && c.Number != "2111002" {
				score -= 100
			}
			if first {
				score += 100
			}
			replace := ""
			if !c.CanPlace {
				for _, id := range c.ReplacementIDs {
					old := s.You.find(id)
					if old == nil {
						continue
					}
					empty := old.Number == "2121001" && old.Statuses["火焰标记物"] == 0 || old.Number == "2121014" && old.Statuses["火匣子标记物"] == 0
					if empty {
						replace = id
						break
					}
				}
				if replace == "" {
					continue
				}
			}
			add(action("equip", map[string]any{"instance_id": c.ID, "replace_id": replace}), c.Cost, score, "部署装备资源，遵守神器与其他装备互斥")
		} else if c.Number == "2021107" && len(s.You.Hand) <= 2 {
			add(action("use_item", map[string]any{"instance_id": c.ID}), c.Cost, 170, "重塑：消化有用手牌后补充两张牌")
		} else if c.Counter && hasEquipmentSpace(s.You) {
			add(action("use_item", map[string]any{"instance_id": c.ID}), nil, 100, "免费盖放火焰符文，触发时另作决定")
		}
	}
	for _, c := range s.You.field() {
		if c != nil && c.Number == "2121004" && opts[c.ID].CanUltimate {
			add(action("use_ability", map[string]any{"instance_id": c.ID, "ability_type": "ultimate"}), nil, 430, "火焰箭：使用效果伤害补人物或清阻挡")
		}
	}
	if !first {
		for _, c := range s.You.Pool {
			if c == nil {
				continue
			}
			score := 0
			switch c.Number {
			case "3121001":
				score = 400
			case "3121109":
				score = 355
			case "3121105":
				score = 335
			case "3121003":
				score = 330
			case "3121002":
				score = 300
			case "3121007":
				if readySpells >= 2 {
					score = 365
				}
			case "3121008":
				if fire >= 7 && known >= 3 {
					score = 290
				}
			case "3121015":
				for _, u := range s.Opponent.units() {
					if u.Position.Row > s.Opponent.front() && u.Life == 1 {
						score = 385
					}
				}
				if pegasus {
					score = 420
				}
			case "3121107":
				if fire >= 7 {
					score = 200
				}
			case "3021006":
				for _, e := range s.Opponent.Equipment {
					if e != nil && e.SetCounter && known >= 3 {
						score = 190
					}
				}
			}
			if score == 0 {
				continue
			}
			if grow && known > 0 {
				score -= 45
			}
			if known == 0 && c.Number == "3121001" {
				score = 460
			}
			replace := ""
			if known >= s.You.SkillSlots {
				for _, old := range s.You.Skills {
					oldRank := map[string]int{"3021006": 190, "3121008": 290, "3121107": 200}[oldNumber(old)]
					if old != nil && score > oldRank+30 && !old.Horizontal && old.Statuses["冷却"] == 0 && (old.Number == "3021006" || old.Number == "3121008" || old.Number == "3121107") {
						replace = old.ID
						break
					}
				}
				if replace == "" {
					continue
				}
			}
			add(action("learn_skill", map[string]any{"instance_id": c.ID, "replace_id": replace}), c.LearnCost, score, "学习进攻技能：第二个自己的回合起，根据场面补齐法术")
		}
	}
	sort.SliceStable(moves, func(i, j int) bool { return moves[i].score > moves[j].score })
	if p.steps < 100 {
		for _, m := range moves {
			if m.score < 0 {
				continue
			}
			if move, ok := p.fund(s, opts, m.move, m.cost); ok {
				return Decision{move, m.reason}
			}
		}
	}
	return Decision{action("end_turn", nil), "结束主阶段，继续响应余火与强制选择"}
}

func targetScore(c *card, s state, damage int, pegasus, wind bool) int {
	score := 300
	if c.Type == "人物" {
		score = 420
		if damage >= c.Life+s.Opponent.Shield {
			score = 1100
		}
	}
	if c.Position != nil && c.Position.Row == s.Opponent.front() {
		score += 20
	}
	if c.Type != "人物" && damage >= c.Life {
		score += 45
	}
	if wind && c.Type != "人物" && c.Life == 1 && c.Position.Row > s.Opponent.front() {
		score = 475
	}
	if pegasus {
		if c.Number == "1511003" {
			score = 650
		} else {
			// A nominal lethal attack is still prevented by the public aura.
			score = min(score, 100)
		}
	}
	return score
}

func hasEquipmentSpace(p player) bool {
	for _, c := range p.Equipment {
		if c == nil {
			return true
		}
	}
	return false
}

func oldNumber(c *card) string {
	if c == nil {
		return ""
	}
	return c.Number
}
func placement(s state, c *card, devour []string) (game.Position, bool) {
	removed := map[string]bool{}
	for _, id := range devour {
		removed[id] = true
	}
	var positions []game.Position
	if c.Number == "1111001" || c.Attack > 0 {
		positions = []game.Position{{Col: 1, Row: 0}, {Col: 0, Row: 0}, {Col: 2, Row: 0}}
	}
	if c.Number == "1121014" {
		positions = append(positions, game.Position{Col: 0, Row: 0}, game.Position{Col: 2, Row: 0})
	}
	for _, pos := range []game.Position{{Col: 0, Row: 2}, {Col: 2, Row: 2}, {Col: 1, Row: 2}, {Col: 0, Row: 1}, {Col: 2, Row: 1}, {Col: 0, Row: 0}, {Col: 2, Row: 0}, {Col: 1, Row: 0}} {
		positions = append(positions, pos)
	}
	for _, pos := range positions {
		old := s.You.Units[pos.Col][pos.Row]
		if old != nil && !removed[old.ID] {
			continue
		}
		if c.Number == "1021102" && buffTarget(s, pos) == nil {
			continue
		}
		return pos, true
	}
	return game.Position{}, false
}
func buffTarget(s state, pos game.Position) *card {
	var best *card
	for _, c := range s.You.units() {
		if c.Type == "人物" {
			continue
		}
		dx := c.Position.Col - pos.Col
		if dx < 0 {
			dx = -dx
		}
		dy := c.Position.Row - pos.Row
		if dy < 0 {
			dy = -dy
		}
		if dx+dy != 1 {
			continue
		}
		if best == nil || c.Position.Row < best.Position.Row {
			best = c
		}
	}
	return best
}
