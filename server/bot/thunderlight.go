package bot

import (
	"eraofarcane/game"
	"sort"
	"strings"
)

func (p *Policy) isThunderlight() bool { return p.profile == ThunderlightProfileID }

func deckCount(s state, number string) int {
	for _, c := range s.You.Deck {
		if c.Number == number {
			return c.Count
		}
	}
	return 0
}

func thunderResources(s state) (air, light, arcane, independent int) {
	for _, c := range s.You.field() {
		if c == nil || c.Statuses[game.StatusStun] > 0 || c.Statuses[game.StatusFreeze] > 0 || c.Life == 1 && c.Statuses[game.StatusBurn] > 0 {
			continue
		}
		air += c.Gain["气"]
		light += c.Gain["光"]
		arcane += c.Gain["无"]
		if c.Type != "人物" && sum(c.Gain) > 0 {
			independent++
		}
	}
	return
}

func thunderItems(s state) (count, armor int) {
	for _, c := range s.You.Equipment {
		if c == nil {
			continue
		}
		if c.Gain["气"] > 0 && c.Gain["光"] > 0 {
			count++
		}
		if c.Number == "2321105" && c.Statuses[game.StatusPetrify] == 0 {
			armor++
		}
	}
	return
}

func ownCount(s state, n string) int {
	count := 0
	for _, c := range s.You.field() {
		if c != nil && c.Number == n && c.Statuses[game.StatusPetrify] == 0 {
			count++
		}
	}
	return count
}

func (p *Policy) thunderMulligan(s state) bool {
	stable, cheap, draw := 0, 0, 0
	for _, c := range s.You.Hand {
		if sum(c.Cost) <= 2 && sum(c.Gain) > 0 {
			stable++
			cheap++
		}
		if c.Number == "1321110" || c.Number == "2021101" {
			draw++
		}
		if c.Number == "1021016" {
			stable++
		}
	}
	return cheap >= 2 || cheap >= 1 && stable+draw >= 3
}

type thunderIntent struct {
	intent
	reserve map[string]int
	devour  []string
}

func (p *Policy) thunderMain(s state, opts map[string]game.DecisionCard) Decision {
	var moves []thunderIntent
	add := func(move game.ActionMessage, cost map[string]int, score int, reason string, reserve map[string]int, devour []string) {
		moves = append(moves, thunderIntent{intent{move, cost, score, reason}, reserve, devour})
	}
	air, light, arcane, independent := thunderResources(s)
	totalAir := air + arcane
	grow := totalAir < 6 || light+arcane < 2 || independent < 2
	items, armor := thunderItems(s)
	beasts := ownCount(s, "1321006")
	known, airKnown, cheapAirKnown, ready := 0, 0, 0, 0
	var reserve map[string]int
	for _, c := range s.You.Skills {
		if c == nil {
			continue
		}
		known++
		if c.Category == "气" {
			airKnown++
			if sum(c.LearnCost) <= 3 {
				cheapAirKnown++
			}
		}
		if opts[c.ID].CanCast && !c.Sorcery {
			ready++
			cost := opts[c.ID].UseCost
			if reserve == nil || sum(cost) < sum(reserve) {
				reserve = cost
			}
		}
	}
	waitArmorer := false
	if noEquipment(s) {
		for _, c := range s.You.Hand {
			if c.Number == "1021016" && deckEquipment(s) > 0 {
				waitArmorer = true
			}
		}
	}
	for _, c := range s.You.Hand {
		o := opts[c.ID]
		if c.Type == "伙伴" && o.CanSummon {
			score := 230
			if grow {
				score = 420
			}
			if ownTurn(s) == 1 {
				score += 70
			}
			var devour []string
			var keep map[string]int
			switch c.Number {
			case "1321002":
				score = 450
				if len(s.You.units()) >= 7 {
					score = 120
				}
			case "1021001":
				if grow {
					score = 465
				}
			case "1321003":
				if c.Drawn {
					score = 510
				}
			case "1321011":
				if light+arcane < 2 {
					score = 510
				}
			case "1021016":
				if noEquipment(s) && deckEquipment(s) > 0 {
					score = 520
				} else {
					score = 140
				}
			case "1321110":
				if deckCount(s, "2021101") > 0 {
					score = 435
					if len(s.You.Hand) <= 3 {
						score = 490
					}
				}
			case "1321006":
				score = 220
				if airKnown >= 2 {
					score = 620
					keep = reserve
				}
				if beasts >= 2 {
					score = 100
				}
			case "1321010":
				score = 0
				if airKnown >= 2 && airKnown+ready >= 3 {
					devour = p.thunderDevour(s, opts, c)
					if len(devour) > 0 {
						score = 610
						keep = reserve
					}
				}
			case "1321111":
				score = 155
				if items >= 2 {
					score = 430
				}
				if items >= 3 {
					score = 475
				}
			case "1311001":
				score = 0
				if len(s.You.Hand) <= 3 && len(s.You.units()) <= 6 && len(s.You.Deck) > 0 && totalAir >= 7 {
					score = 435
					keep = map[string]int{"无": 2}
				}
			}
			pos, ok := thunderPlacement(s, c, devour)
			if !ok || score <= 0 {
				continue
			}
			if c.Temporary {
				score += 35
			}
			if !grow && sum(c.Gain) > 0 && sum(c.Cost) <= 3 {
				// After the initial spell package, cheap load still takes priority
				// over buying a fourth expensive skill after a costly defense.
				score = max(score, 350)
			}
			if grow && sum(c.Gain) > 0 && sum(c.Cost) <= 3 {
				// Build renewable load before spending it on a growing skill pool.
				score += 110
			}
			add(action("summon", map[string]any{"instance_id": c.ID, "col": pos.Col, "row": pos.Row, "devour_ids": devour}), c.Cost, score, "雷光发育：部署伙伴并保护气光来源", keep, devour)
		} else if c.Equipment {
			score := 365
			if grow {
				score = 425
			}
			keep := map[string]int(nil)
			if c.Number == "2021104" && light+arcane < 2 {
				score = 500
			}
			if c.Number == "2321105" {
				score = 155
				if items >= 2 {
					score = 615
					keep = reserve
				} else if items == 1 && len(s.You.Hand) >= 4 {
					score = 340
				}
			}
			if c.Number == "2021105" {
				score = 90
				for _, h := range s.You.Hand {
					if h.Number == "2321105" && armor > 0 && !s.You.Duplicates {
						score = 560
					}
				}
			}
			if waitArmorer && !(c.Number == "2021104" && light+arcane == 0) && !(c.Number == "2321105" && items >= 2) {
				score = 0
			}
			replace := ""
			if !c.CanPlace {
				for _, id := range c.ReplacementIDs {
					old := s.You.find(id)
					if old != nil && p.keepValue(old, s)+45 < p.keepValue(c, s) && !(old.Gain["光"] > 0 && old.Gain["气"] > 0 && items <= 3) {
						replace = id
						break
					}
				}
				if replace == "" {
					continue
				}
			}
			if score > 0 {
				add(action("equip", map[string]any{"instance_id": c.ID, "replace_id": replace}), c.Cost, score, "部署雷光组件：先检索，按缺件与光负载安排", keep, nil)
			}
		} else if c.Number == "2021101" && len(s.You.Deck) > 0 {
			score := 480
			if len(s.You.Hand) >= 6 {
				score = 230
			}
			add(action("use_item", map[string]any{"instance_id": c.ID}), c.Cost, score, "银叶花寻找资源与组件，利用弃牌收益", nil, nil)
		}
	}
	for _, c := range s.You.spells() {
		if c == nil || !opts[c.ID].CanCast || c.Sorcery {
			continue
		}
		for _, target := range opts[c.ID].Targets {
			if c.Number == "3021101" && !thunderSpareLight(s, opts, target.Cost) {
				continue
			}
			score := p.thunderAttackValue(s, c, target)
			if score <= 0 {
				continue
			}
			data := map[string]any{"instance_id": c.ID, "target_type": target.Type, "target_col": target.Position.Col, "target_row": target.Position.Row, "boost_ids": []string{}}
			if target.AreaOrigin != nil {
				data["area_col"] = target.AreaOrigin.Col
				data["area_row"] = target.AreaOrigin.Row
			}
			if target.ExtraTarget != nil {
				data["extra_target_col"] = target.ExtraTarget.Col
				data["extra_target_row"] = target.ExtraTarget.Row
			}
			plainScore := score
			if score >= 500 && target.Power <= publicDefensePower(s) {
				plainScore -= 90
			}
			add(action("cast_spell", data), target.Cost, plainScore, "雷光进攻：常态清前排与高威胁目标，低血时追击人物", nil, nil)
			if score >= 500 {
				p.thunderBoostMoves(s, opts, c, target, data, score, add)
			}
		}
	}
	for _, c := range s.You.field() {
		if c == nil {
			continue
		}
		for _, pos := range opts[c.ID].AttackTargets {
			t := s.Opponent.Units[pos.Col][pos.Row]
			if t == nil {
				continue
			}
			score := thunderUnitValue(t, s, c.Attack) + 45
			if t.Statuses["防止致命"] > 0 && c.Attack >= t.Life {
				continue
			}
			add(action("attack", map[string]any{"attacker_id": c.ID, "target_col": pos.Col, "target_row": pos.Row}), nil, score, "直接攻击前排，保留有攻击价值的伙伴朝向", nil, nil)
		}
		if c.Number == "2321101" && opts[c.ID].CanPerTurn && len(s.Opponent.units()) >= 2 {
			for _, skill := range s.You.spells() {
				if skill != nil && strings.Contains(skill.Tag, "驱动") && opts[skill.ID].CanCast && len(opts[skill.ID].Targets) > 0 {
					copy := s
					copy.You.Elements = thunderAvailable(s, opts, c.ID)
					if _, ok := payment(copy, copy.You.Elements, opts[skill.ID].UseCost); ok {
						add(action("use_ability", map[string]any{"instance_id": c.ID, "ability_type": "per_turn"}), nil, 545, "雷之链充能：已有可支付驱动和第二目标", nil, nil)
						break
					}
				}
			}
		}
	}
	for _, c := range s.You.Pool {
		if c == nil || !opts[c.ID].CanLearn || p.cycled[c.ID] {
			continue
		}
		score := thunderLearnValue(s, c, air, light, arcane, known, beasts)
		if c.Number == "3021101" && (airKnown < 2 || grow || !thunderSpareLight(s, opts, c.LearnCost)) {
			continue
		}
		// Establish several efficient air spells before buying a large spell.
		if cheapAirKnown < 2 && sum(c.LearnCost) >= 4 {
			score -= 240
		}
		if known >= 3 {
			score -= 160
			if c.Number == "3021101" {
				score = max(score, 60)
			}
		}
		if ownTurn(s) == 1 {
			score -= 180
		}
		if known == 0 && c.Category == "气" && sum(c.LearnCost) <= 3 {
			score += 30
		}
		replace := ""
		if known >= s.You.SkillSlots {
			lowest := 10000
			for _, id := range opts[c.ID].LearnReplacementIDs {
				old := s.You.find(id)
				if old == nil || p.cycled[id] {
					continue
				}
				value := thunderLearnValue(s, old, air, light, arcane, known, beasts)
				if value < lowest && score > value+55 {
					lowest = value
					replace = id
				}
			}
			if replace == "" {
				continue
			}
		}
		if score > 0 {
			add(action("learn_skill", map[string]any{"instance_id": c.ID, "replace_id": replace}), c.LearnCost, score, "雷光学习：缺发展时以攻代守，按场面选择五槽", nil, nil)
		}
	}
	sort.SliceStable(moves, func(i, j int) bool { return moves[i].score > moves[j].score })
	if p.steps < 100 {
		for _, m := range moves {
			if m.score <= 0 {
				continue
			}
			budget := m.reserve
			if !p.attacked && !s.Suppressed && reserve != nil && m.move.Action != "cast_spell" && m.move.Action != "attack" && sum(reserve) > sum(budget) {
				budget = reserve
			}
			if _, ok := thunderPay(s, thunderAvailable(s, opts, ""), plus(m.cost, budget)); !ok {
				continue
			}
			for _, id := range m.devour {
				if opts[id].CanConsume && sum(opts[id].ConsumeGain) > 0 {
					move := action("consume", map[string]any{"instance_id": id})
					if !p.blocked[Fingerprint(move)] {
						return Decision{move, "奇美拉入场前先取得材料负载"}
					}
				}
			}
			if move, ok := p.thunderFund(s, opts, m.move, m.cost); ok {
				return Decision{move, m.reason}
			}
		}
	}
	return Decision{action("end_turn", nil), "雷光计划完成，保留资源来源并处理回合末选择"}
}

func noEquipment(s state) bool {
	for _, c := range s.You.Equipment {
		if c != nil {
			return false
		}
	}
	return true
}
func deckEquipment(s state) int {
	n := 0
	for _, c := range s.You.Deck {
		if c.Equipment {
			n += c.Count
		}
	}
	return n
}

func thunderLearnValue(s state, c *card, air, light, arcane, known, beasts int) int {
	switch c.Number {
	case "3321002":
		return 440
	case "3321003":
		if beasts > 0 || ownCount(s, "1321010") > 0 {
			return 460
		}
		return 230
	case "3321013":
		if s.Opponent.Hero != nil && s.Opponent.Hero.Life <= 3 {
			return 550
		}
		return 340
	case "3321004":
		if air+arcane >= 6 && light+arcane > 0 {
			return 470
		}
		return 315
	case "3321106":
		if air+arcane >= 6 {
			return 465
		}
		return 300
	case "3321006":
		if len(s.Opponent.units()) >= 3 {
			return 480
		}
		return 285
	case "3321103":
		if air+arcane >= 8 && light+arcane >= 2 {
			return 485
		}
		return 200
	case "3521105":
		if light+arcane >= 3 {
			return 410
		}
		return 150
	case "3021101":
		return 140
	case "3021008":
		for _, item := range s.Opponent.Equipment {
			if item != nil && !item.SetCounter {
				return 505
			}
		}
		return 90
	}
	return 0
}

func thunderUnitValue(c *card, s state, damage int) int {
	if c == nil {
		return 0
	}
	if c.Type == "人物" {
		if damage >= c.Life {
			return 1300
		}
		if c.Life <= 3 {
			return 720
		}
		return 390
	}
	score := 330
	if c.Position != nil && c.Position.Row == s.Opponent.front() {
		score += 45
	}
	if c.Statuses[game.StatusPetrify] == 0 {
		if c.Taunt {
			score += 35
		}
		if c.Aura {
			score += 65
		}
		if c.HasPerTurn || c.HasPrayer {
			score += 35
		}
	}
	if damage >= c.Life {
		score += 80
	}
	score += min(sum(c.Gain)*10, 40)
	return score
}

func (p *Policy) thunderAttackValue(s state, c *card, t game.DecisionTarget) int {
	score, hits := 0, 0
	for _, hit := range t.Impacts {
		u := s.Opponent.find(hit.InstanceID)
		if u == nil {
			continue
		}
		if hit.Damage > 0 {
			v := thunderUnitValue(u, s, hit.Damage)
			if hits == 0 || v > score {
				score = max(score, v)
			}
			hits++
		}
		if hit.Stun > 0 && !u.Horizontal && u.Statuses[game.StatusStun] == 0 && sum(u.Gain) >= 2 {
			score = max(score, 350+sum(u.Gain)*12)
		}
	}
	if hits > 1 {
		score += (hits - 1) * 85
	}
	if c.Number == "3021008" {
		for _, u := range s.Opponent.Equipment {
			if u != nil && !u.SetCounter {
				score = max(score, 380+sum(u.Gain)*20)
			}
		}
	}
	if score == 0 {
		return 0
	}
	if s.Suppressed {
		score -= 120
	}
	return score - sum(t.Cost)*5
}
