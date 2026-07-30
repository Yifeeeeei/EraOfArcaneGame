package game

import (
	"fmt"
	"math/rand"
	"strings"

	"eraofarcane/model"
)

type Card1321108EmeraldHummingbird struct{ AlwaysActive }

func (Card1321108EmeraldHummingbird) ID() string   { return "1321108" }
func (Card1321108EmeraldHummingbird) Name() string { return "翡翠蜂鸟" }
func (Card1321108EmeraldHummingbird) OnEnter(ctx *EffectContext) error {
	if len(ctx.Engine.State.Players[ctx.PlayerID].Hand) < 2 {
		ctx.Engine.drawCards(ctx.PlayerID, 2)
	}
	return nil
}

type Card1001101AbandonedPawn struct{ AlwaysActive }

func (Card1001101AbandonedPawn) ID() string   { return "1001101" }
func (Card1001101AbandonedPawn) Name() string { return "弃子" }
func (Card1001101AbandonedPawn) OnDeath(ctx *EffectContext) error {
	if ctx.Source == nil || ctx.Source.Position == nil {
		return nil
	}
	pos := *ctx.Source.Position
	damaged := adjacentUnits(ctx.Engine.State.Players[ctx.PlayerID], &pos)
	damaged = append(damaged, adjacentUnits(ctx.Engine.State.Players[ctx.OpponentID], &pos)...)
	for _, target := range damaged {
		if target == nil || target.CurrentLife <= 0 {
			continue
		}
		targetPos := Position{}
		if target.Position != nil {
			targetPos = *target.Position
		}
		ctx.Engine.dealDamageWithExtra(target, 1, target.OwnerID, map[string]any{
			"damage_source": "effect",
			"attacker":      ctx.PlayerID,
		})
		if target.CurrentLife <= 0 && !target.Card.IsHero() {
			ownerID := target.OwnerID
			if ctx.Engine.unitInOwnerGrid(target, ownerID) {
				ctx.Engine.destroyUnitWithData(target, ownerID, map[string]any{
					"death_cause": "abandoned_pawn",
					"attacker":    ctx.PlayerID,
				})
			}
			if ctx.Engine.State.Players[ownerID].Units[targetPos.Col][targetPos.Row] == nil {
				ctx.Engine.summonFreshCardAtPosition(ownerID, "1001101", targetPos, true)
			}
		}
	}
	return nil
}

type Card1021105RoyalTaxCollector struct{ AlwaysActive }

func (Card1021105RoyalTaxCollector) ID() string   { return "1021105" }
func (Card1021105RoyalTaxCollector) Name() string { return "皇城征税员" }

const royalTaxCollectorUntilOpponentTurnEndStatus = "皇城征税员征税至对手回合结束"

func (Card1021105RoyalTaxCollector) OnEnter(ctx *EffectContext) error {
	if ctx.Source == nil {
		return nil
	}
	ctx.Source.Statuses[royalTaxCollectorUntilOpponentTurnEndStatus] = ctx.Engine.State.TurnNumber
	return nil
}

func (Card1021105RoyalTaxCollector) OnDraw(ctx *EffectContext) error {
	if ctx.Source == nil || ctx.Source.Statuses[royalTaxCollectorUntilOpponentTurnEndStatus] <= 0 || ctx.ExtraData == nil {
		return nil
	}
	drawnPlayer, _ := ctx.ExtraData["drawn_player"].(int)
	if drawnPlayer != ctx.OpponentID {
		return nil
	}
	ctx.Engine.State.Players[ctx.PlayerID].Elements[model.ElementArcane]++
	return nil
}

func (Card1021105RoyalTaxCollector) OnTurnEnd(ctx *EffectContext) error {
	if ctx.Source == nil || ctx.Source.Statuses[royalTaxCollectorUntilOpponentTurnEndStatus] <= 0 || ctx.ExtraData == nil {
		return nil
	}
	endedPlayer, _ := ctx.ExtraData["ended_player"].(int)
	if endedPlayer == ctx.OpponentID && ctx.Engine.State.TurnNumber >= ctx.Source.Statuses[royalTaxCollectorUntilOpponentTurnEndStatus] {
		delete(ctx.Source.Statuses, royalTaxCollectorUntilOpponentTurnEndStatus)
	}
	return nil
}

type Card1121106FireBeastTrainer struct{ AlwaysActive }

func (Card1121106FireBeastTrainer) ID() string   { return "1121106" }
func (Card1121106FireBeastTrainer) Name() string { return "弗卡莱诺皇家驯兽师" }

const fireBeastTrainerDiscountStatus = "弗卡莱诺皇家驯兽师下个火焰野兽异兽减费"

func (Card1121106FireBeastTrainer) OnEnter(ctx *EffectContext) error {
	if ctx.Source == nil {
		return nil
	}
	ctx.Source.Statuses[fireBeastTrainerDiscountStatus] = 1
	return nil
}

func (Card1121106FireBeastTrainer) ModifyCardPlayCost(ctx *EffectContext, card *CardInstance, cost map[string]int) {
	if ctx.Source == nil || ctx.Source.Statuses[fireBeastTrainerDiscountStatus] <= 0 || !isFireBeastOrMonsterCompanion(card) {
		return
	}
	reduceGenericCost(cost, model.ElementFire, 2)
}

func (Card1121106FireBeastTrainer) OnCardPlayCostPaid(ctx *EffectContext, card *CardInstance) {
	if ctx.Source == nil || ctx.Source.Statuses[fireBeastTrainerDiscountStatus] <= 0 || !isFireBeastOrMonsterCompanion(card) {
		return
	}
	ctx.Source.Statuses[fireBeastTrainerDiscountStatus]--
	if ctx.Source.Statuses[fireBeastTrainerDiscountStatus] <= 0 {
		delete(ctx.Source.Statuses, fireBeastTrainerDiscountStatus)
	}
}

func isFireBeastOrMonsterCompanion(card *CardInstance) bool {
	if card == nil || card.Card == nil || !card.Card.IsCompanion() || card.Card.Category != model.ElementFire {
		return false
	}
	return strings.Contains(card.Card.Tag, "野兽") || strings.Contains(card.Card.Tag, "异兽")
}

type Card1221106MirrorLotus struct{ AlwaysActive }

func (Card1221106MirrorLotus) ID() string            { return "1221106" }
func (Card1221106MirrorLotus) Name() string          { return "镜花海之莲" }
func (Card1221106MirrorLotus) IsPrayerAbility() bool { return true }
func (Card1221106MirrorLotus) OnPerTurn(ctx *EffectContext) error {
	ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementWater, 1, ctx.Source)
	return nil
}

type Card1021102SwordsmanshipTeacher struct{ AlwaysActive }

func (Card1021102SwordsmanshipTeacher) ID() string   { return "1021102" }
func (Card1021102SwordsmanshipTeacher) Name() string { return "剑术师傅" }
func (Card1021102SwordsmanshipTeacher) OnEnter(ctx *EffectContext) error {
	candidates := adjacentFriendlyCompanions(ctx)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "swordsmanship_teacher_buff",
		"剑术师傅:选择1个相邻友方伙伴获得+1攻", candidates, 1, 1,
		func(selected []string) {
			target := ctx.Engine.findFieldCardByInstance(ctx.Engine.State.Players[ctx.PlayerID], firstSelected(selected))
			if target != nil && target.Card != nil && target.Card.IsCompanion() {
				target.AttackBonus++
			}
		})
	return nil
}

type Card1021101PrivateTeacher struct{ AlwaysActive }

func (Card1021101PrivateTeacher) ID() string   { return "1021101" }
func (Card1021101PrivateTeacher) Name() string { return "私家教师" }
func (Card1021101PrivateTeacher) OnEnter(ctx *EffectContext) error {
	candidates := make([]map[string]any, 0)
	allowed := make(map[string]bool)
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	hasEmptySlot := false
	for i := 0; i < skillSlotCapacity(ps); i++ {
		if ps.Skills[i] == nil {
			hasEmptySlot = true
			break
		}
	}
	for _, skill := range ps.SkillPool {
		if skill == nil || skill.Card == nil || !skill.Card.IsSkill() || totalElementCost(skill.Card.ElementsCost) >= 4 {
			continue
		}
		if hasEmptySlot {
			candidates = append(candidates, candidateInfo(skill, "skill_pool", "own"))
			allowed[skill.InstanceID] = true
			continue
		}
		for _, learned := range ps.Skills {
			if learned == nil || learned.IsHorizontal {
				continue
			}
			id := skill.InstanceID + "|" + learned.InstanceID
			candidate := candidateInfo(skill, "skill_pool", "own")
			candidate["instance_id"] = id
			candidate["name"] = fmt.Sprintf("学习%s，替换%s", skill.Card.Name, learned.Card.Name)
			candidate["replace_id"] = learned.InstanceID
			allowed[id] = true
			candidates = append(candidates, candidate)
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "private_teacher_learn_skill",
		"私家教师:选择1个学习花费小于4的法术免费学习", candidates, 1, 1,
		func(selected []string) {
			id := firstSelected(selected)
			if !allowed[id] {
				return
			}
			skillID := id
			replaceID := ""
			if before, after, ok := strings.Cut(id, "|"); ok {
				skillID = before
				replaceID = after
			}
			ctx.Engine.learnSkillFromPoolWithoutCost(ctx.PlayerID, skillID, replaceID)
		})
	return nil
}

type Card1021104DimensionalRiftBeast struct{ AlwaysActive }

func (Card1021104DimensionalRiftBeast) ID() string   { return "1021104" }
func (Card1021104DimensionalRiftBeast) Name() string { return "次元撕裂兽" }
func (Card1021104DimensionalRiftBeast) OnEnter(ctx *EffectContext) error {
	candidates := companionSpellRangeCandidates(ctx, false)
	filtered := candidates[:0]
	for _, candidate := range candidates {
		if candidate["side"] == "enemy" {
			filtered = append(filtered, candidate)
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "dimensional_rift_beast_exile",
		"次元撕裂兽:选择法力范围内1个敌方伙伴移出游戏", filtered, 1, 1,
		func(selected []string) {
			target := selectedUnitFromCandidates(ctx.Engine, selected, filtered)
			if target != nil && target.Card != nil && target.Card.IsCompanion() {
				ctx.Engine.exileCard(target.OwnerID, target)
			}
		})
	return nil
}

type Card1021106SkyCityTycoon struct{ AlwaysActive }

func (Card1021106SkyCityTycoon) ID() string   { return "1021106" }
func (Card1021106SkyCityTycoon) Name() string { return "云霄城富豪" }
func (Card1021106SkyCityTycoon) PerTurnLabel(*CardInstance) string {
	return "主动"
}
func (Card1021106SkyCityTycoon) OnPerTurn(ctx *EffectContext) error {
	if !ctx.Engine.canConsumeCard(ctx.Source) {
		return fmt.Errorf("云霄城富豪不能被消耗")
	}
	ctx.Source.IsHorizontal = true
	choices := []map[string]any{
		{"instance_id": "self_first", "number": "1021106", "name": "你先抽", "type": "选择", "zone": "choice", "side": "own"},
		{"instance_id": "opponent_first", "number": "1021106", "name": "对手先抽", "type": "选择", "zone": "choice", "side": "own"},
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "sky_city_tycoon_draw_order",
		"云霄城富豪:选择双方抽牌次序", choices, 1, 1,
		func(selected []string) {
			if firstSelected(selected) == "opponent_first" {
				ctx.Engine.drawCards(ctx.OpponentID, 1)
				ctx.Engine.drawCards(ctx.PlayerID, 1)
				return
			}
			ctx.Engine.drawCards(ctx.PlayerID, 1)
			ctx.Engine.drawCards(ctx.OpponentID, 1)
		})
	return nil
}

type Card1021108AlchemyApprentice struct{ AlwaysActive }

func (Card1021108AlchemyApprentice) ID() string   { return "1021108" }
func (Card1021108AlchemyApprentice) Name() string { return "炼金术学徒" }
func (Card1021108AlchemyApprentice) PerTurnLabel(*CardInstance) string {
	return "主动"
}
func (Card1021108AlchemyApprentice) OnPerTurn(ctx *EffectContext) error {
	if !ctx.Engine.canConsumeCard(ctx.Source) {
		return fmt.Errorf("炼金术学徒不能被消耗")
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ps.Elements[model.ElementArcane] < 1 {
		return fmt.Errorf("炼金术学徒需要1点奥术元素")
	}
	choices := make([]map[string]any, 0, 12)
	for _, elem := range []string{model.ElementFire, model.ElementWater, model.ElementEarth, model.ElementAir, model.ElementLight, model.ElementShadow} {
		for i := 1; i <= 2; i++ {
			choices = append(choices, map[string]any{"instance_id": fmt.Sprintf("%s#%d", elem, i), "number": "1021108", "name": elem, "type": "选择", "zone": "choice", "side": "own", "element": elem})
		}
	}
	ctx.Source.IsHorizontal = true
	ps.Elements[model.ElementArcane]--
	ctx.Engine.SetPendingAction(ctx.PlayerID, "alchemy_apprentice_elements",
		"炼金术学徒:选择2点非奥术元素", choices, 2, 2,
		func(selected []string) {
			gain := make(map[string]int)
			seen := make(map[string]bool, len(selected))
			for _, id := range selected {
				if seen[id] {
					continue
				}
				seen[id] = true
				elem, _, ok := strings.Cut(id, "#")
				if ok && isNonArcaneElement(elem) {
					gain[elem]++
				}
			}
			if len(gain) > 0 {
				ps.GainElements(gain)
			}
		})
	return nil
}

type Card1021109ChurchEnvoy struct{ AlwaysActive }

func (Card1021109ChurchEnvoy) ID() string   { return "1021109" }
func (Card1021109ChurchEnvoy) Name() string { return "教廷特使" }
func (Card1021109ChurchEnvoy) OnUltimate(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, hasAnyNegativeStatus)
	candidates = append(candidates, ctx.Engine.friendlyEquipment(ctx.PlayerID, hasAnyNegativeStatus)...)
	candidates = append(candidates, ctx.Engine.friendlySkillsIncludingBound(ctx.PlayerID, hasAnyNegativeStatus)...)
	if len(candidates) == 0 {
		return nil
	}
	allowed := make(map[string]bool, len(candidates))
	for _, candidate := range candidates {
		if id, _ := candidate["instance_id"].(string); id != "" {
			allowed[id] = true
		}
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "church_envoy_purify",
		"教廷特使:选择1张友方卡牌移除全部负面效果", candidates, 1, 1,
		func(selected []string) {
			id := firstSelected(selected)
			if !allowed[id] {
				return
			}
			clearNegativeStatuses(ctx.Engine.findFriendlyCardIncludingBound(ctx.PlayerID, id))
		})
	return nil
}

type Card1121103BeaconGuard struct{ AlwaysActive }

func (Card1121103BeaconGuard) ID() string   { return "1121103" }
func (Card1121103BeaconGuard) Name() string { return "烽火台守卫" }
func (Card1121103BeaconGuard) OnEnter(ctx *EffectContext) error {
	if royalCompanionCount(ctx.Engine.State.Players[ctx.PlayerID]) < royalCompanionCount(ctx.Engine.State.Players[ctx.OpponentID]) {
		ctx.Engine.gainPlayerShield(ctx.PlayerID, 3)
	}
	return nil
}

type Card1121108FireButterfly struct{ AlwaysActive }

const (
	fireButterflyTemporaryLoadStatus     = "火蝴蝶临时负载"
	fireButterflyPreviousLoadSetStatus   = "火蝴蝶原负载覆盖"
	fireButterflyPreviousLoadValuePrefix = "火蝴蝶原负载:"
)

func (Card1121108FireButterfly) ID() string   { return "1121108" }
func (Card1121108FireButterfly) Name() string { return "火蝴蝶" }
func (Card1121108FireButterfly) PerTurnLabel(*CardInstance) string {
	return "主动"
}
func (Card1121108FireButterfly) OnPerTurn(ctx *EffectContext) error {
	clearFireButterflyStoredLoad(ctx.Source)
	if ctx.Source.ElementsGainSet != nil {
		ctx.Source.Statuses[fireButterflyPreviousLoadSetStatus] = 1
		for _, elem := range model.AllElements {
			if amount := ctx.Source.ElementsGainSet[elem]; amount != 0 {
				ctx.Source.Statuses[fireButterflyPreviousLoadValuePrefix+elem] = amount
			}
		}
	}
	ctx.Source.ElementsGainSet = map[string]int{model.ElementAir: 1}
	ctx.Source.Statuses[fireButterflyTemporaryLoadStatus] = 1
	return nil
}
func (Card1121108FireButterfly) OnTurnEnd(ctx *EffectContext) error {
	if ctx.Source.Statuses[fireButterflyTemporaryLoadStatus] <= 0 {
		return nil
	}
	if fireButterflyTemporaryLoadStillCurrent(ctx.Source) {
		if ctx.Source.Statuses[fireButterflyPreviousLoadSetStatus] > 0 {
			previous := make(map[string]int)
			for _, elem := range model.AllElements {
				if amount := ctx.Source.Statuses[fireButterflyPreviousLoadValuePrefix+elem]; amount != 0 {
					previous[elem] = amount
				}
			}
			setElementsGain(ctx.Source, previous)
		} else {
			clearElementsGainSet(ctx.Source)
		}
	}
	clearFireButterflyStoredLoad(ctx.Source)
	return nil
}

type Card1421115Geomancer struct{ AlwaysActive }

func (Card1421115Geomancer) ID() string   { return "1421115" }
func (Card1421115Geomancer) Name() string { return "地卜行者" }
func (Card1421115Geomancer) OnEnter(ctx *EffectContext) error {
	ctx.Engine.drawCards(ctx.PlayerID, 1)
	return nil
}

type Card1221112WaterMage struct{ AlwaysActive }

func (Card1221112WaterMage) ID() string   { return "1221112" }
func (Card1221112WaterMage) Name() string { return "水魔导师" }
func (Card1221112WaterMage) OnUltimate(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlySkillsIncludingBound(ctx.PlayerID, func(skill *CardInstance) bool {
		return skill != nil && skill.Card != nil && skill.Card.IsSkill() &&
			skill.Card.Category == model.ElementWater &&
			totalElementCost(skill.Card.ElementsExpense) < 3 &&
			skill.IsHorizontal
	})
	if len(candidates) == 0 {
		return nil
	}
	allowed := make(map[string]bool, len(candidates))
	for _, candidate := range candidates {
		if id, _ := candidate["instance_id"].(string); id != "" {
			allowed[id] = true
		}
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "water_mage_reset_skill",
		"水魔导师:选择1个使用花费小于3的水纹法术重置", candidates, 1, 1,
		func(selected []string) {
			id := firstSelected(selected)
			if !allowed[id] {
				return
			}
			skill := ctx.Engine.findSkill(ctx.Engine.State.Players[ctx.PlayerID], id)
			if skill != nil && skill.Card != nil && skill.Card.IsSkill() && skill.Card.Category == model.ElementWater && totalElementCost(skill.Card.ElementsExpense) < 3 {
				skill.IsHorizontal = false
			}
		})
	return nil
}

type Card1421107DragonBloodTreant struct{ AlwaysActive }

func (Card1421107DragonBloodTreant) ID() string   { return "1421107" }
func (Card1421107DragonBloodTreant) Name() string { return "龙血树精" }
func (Card1421107DragonBloodTreant) OnEnter(ctx *EffectContext) error {
	candidates := make([]map[string]any, 0)
	for _, card := range ctx.Engine.getAllFieldCards(ctx.Engine.State.Players[ctx.PlayerID]) {
		if card == nil || card.Card == nil {
			continue
		}
		load := dragonBloodTreantReducibleLoad(card)
		for _, elem := range model.AllElements {
			if load[elem] <= 0 {
				continue
			}
			candidate := candidateInfo(card, "field", "own")
			candidate["instance_id"] = card.InstanceID + "|" + elem
			candidate["name"] = card.Card.Name + " - " + elem
			candidate["element"] = elem
			candidates = append(candidates, candidate)
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	removeLoad := func(selection string) {
		instanceID, elem, ok := strings.Cut(selection, "|")
		if !ok || elem == "" {
			return
		}
		target, _ := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, instanceID)
		if target == nil || dragonBloodTreantReducibleLoad(target)[elem] <= 0 {
			return
		}
		dragonBloodTreantRemoveLoad(target, elem)
		ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementShadow, 1, ctx.Source)
	}
	if len(candidates) == 1 {
		id, _ := candidates[0]["instance_id"].(string)
		removeLoad(id)
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "dragon_blood_treant_remove_load",
		"龙血树精:选择1个友方卡牌失去1点负载", candidates, 1, 1,
		func(selected []string) {
			removeLoad(firstSelected(selected))
		})
	return nil
}

func dragonBloodTreantReducibleLoad(card *CardInstance) map[string]int {
	load := make(map[string]int)
	if card == nil || card.Card == nil {
		return load
	}
	base := card.Card.ElementsGain
	if card.ElementsGainSet != nil {
		base = card.ElementsGainSet
	}
	for elem, amount := range base {
		if amount > 0 {
			load[elem] += amount
		}
	}
	for elem, amount := range card.ElementsGainBonus {
		if amount > 0 {
			load[elem] += amount
		}
	}
	return load
}

func dragonBloodTreantRemoveLoad(card *CardInstance, elem string) {
	if card == nil || elem == "" {
		return
	}
	if card.ElementsGainBonus != nil && card.ElementsGainBonus[elem] > 0 {
		card.ElementsGainBonus[elem]--
		if card.ElementsGainBonus[elem] == 0 {
			delete(card.ElementsGainBonus, elem)
		}
		return
	}
	base := copyElementCost(card.Card.ElementsGain)
	if card.ElementsGainSet != nil {
		base = copyElementCost(card.ElementsGainSet)
	}
	if base[elem] <= 0 {
		return
	}
	base[elem]--
	setElementsGain(card, base)
}

type Card1321110SilverleafMessenger struct{ AlwaysActive }

func (Card1321110SilverleafMessenger) ID() string   { return "1321110" }
func (Card1321110SilverleafMessenger) Name() string { return "银叶信使" }
func (Card1321110SilverleafMessenger) OnEnter(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyDeckCards(ctx.PlayerID, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.Number == "2021101"
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "silverleaf_messenger_search",
		"银叶信使:检索1张失落的银叶花", candidates, 1, 1,
		func(selected []string) {
			ctx.Engine.searchDeckToHand(ctx.PlayerID, firstSelected(selected))
		})
	return nil
}

type Card1321111ThunderlightWarrior struct{ AlwaysActive }

func (Card1321111ThunderlightWarrior) ID() string   { return "1321111" }
func (Card1321111ThunderlightWarrior) Name() string { return "雷光战士" }
func (Card1321111ThunderlightWarrior) OnEnter(ctx *EffectContext) error {
	count := 0
	for _, item := range ctx.Engine.State.Players[ctx.PlayerID].Equipment {
		if isThunderlightItem(item) {
			count++
		}
	}
	if count == 0 {
		return nil
	}
	choices := make([]map[string]any, 0, count*4)
	for i := 0; i < count; i++ {
		for _, choice := range []struct {
			id   string
			name string
		}{
			{id: "life", name: "+2血"},
			{id: "attack", name: "+1攻"},
			{id: "air", name: "负载+1气"},
			{id: "light", name: "负载+1光"},
		} {
			choices = append(choices, map[string]any{
				"instance_id": fmt.Sprintf("%s#%d", choice.id, i),
				"number":      "1321111",
				"name":        choice.name,
				"type":        "选择",
				"zone":        "choice",
				"side":        "own",
			})
		}
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "thunderlight_warrior_rewards",
		"雷光战士:每件雷光道具选择1项奖励", choices, count, count,
		func(selected []string) {
			for _, id := range selected {
				reward, _, _ := strings.Cut(id, "#")
				switch reward {
				case "life":
					ctx.Source.CurrentLife += 2
				case "attack":
					ctx.Source.AttackBonus++
				case "air":
					ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementAir, 1, ctx.Source)
				case "light":
					ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementLight, 1, ctx.Source)
				}
			}
		})
	return nil
}

func isThunderlightItem(card *CardInstance) bool {
	if card == nil || card.Card == nil || !card.Card.IsItem() {
		return false
	}
	gain := effectiveElementsGain(card)
	return gain[model.ElementAir] > 0 && gain[model.ElementLight] > 0
}

type Card2321104ThunderlightCrown struct{ AlwaysActive }

func (Card2321104ThunderlightCrown) ID() string            { return "2321104" }
func (Card2321104ThunderlightCrown) Name() string          { return "雷光头冠" }
func (Card2321104ThunderlightCrown) IsPrayerAbility() bool { return true }
func (Card2321104ThunderlightCrown) OnPerTurn(ctx *EffectContext) error {
	ctx.Engine.addNextTaggedSpellPowerBonus(ctx.PlayerID, "聚能", 1)
	return nil
}

type Card2321105ThunderlightArmor struct{ AlwaysActive }

func (Card2321105ThunderlightArmor) ID() string   { return "2321105" }
func (Card2321105ThunderlightArmor) Name() string { return "雷光战铠" }
func (Card2321105ThunderlightArmor) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	skill := ctx.Target
	if skill == nil || skill.Card == nil || !isSpellLikeCard(skill.Card) {
		return
	}
	if !hasCardTag(skill.Card, "驱动") && !hasCardTag(skill.Card, "聚能") {
		return
	}
	count := 0
	for _, item := range ctx.Engine.State.Players[ctx.PlayerID].Equipment {
		if isThunderlightItem(item) {
			count++
		}
	}
	if count >= 3 {
		stats.PowerBonus += 2
	}
}

type Card2321110PigeonRaidOrder struct{ AlwaysActive }

func (Card2321110PigeonRaidOrder) ID() string   { return "2321110" }
func (Card2321110PigeonRaidOrder) Name() string { return "飞鸽急袭令" }
func (Card2321110PigeonRaidOrder) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlySkills(ctx.PlayerID, func(skill *CardInstance) bool {
		return isLearnedRushSkillThisTurn(ctx.Engine, skill)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "pigeon_raid_order_skill",
		"飞鸽急袭令:选择本回合学习的速攻法术", candidates, 1, 1,
		func(selected []string) {
			skill := findSkillSlotCard(ctx.Engine.State.Players[ctx.PlayerID], firstSelected(selected))
			if !isLearnedRushSkillThisTurn(ctx.Engine, skill) {
				return
			}
			ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
				Type:             TempModSkillPowerBonus,
				TargetInstanceID: skill.InstanceID,
				Amount:           1,
				RemainingUses:    1,
				ExpiresTurn:      ctx.Engine.State.TurnNumber + 2,
			})
			ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
				Type:             TempModNextSkillUseAttackBonus,
				TargetInstanceID: skill.InstanceID,
				Amount:           1,
				RemainingUses:    1,
				ExpiresTurn:      ctx.Engine.State.TurnNumber + 2,
			})
		})
	return nil
}

type Card2321107PigeonArrestOrder struct{ AlwaysActive }

func (Card2321107PigeonArrestOrder) ID() string   { return "2321107" }
func (Card2321107PigeonArrestOrder) Name() string { return "飞鸽拘捕令" }
func (Card2321107PigeonArrestOrder) OnSpellHit(ctx *EffectContext) error {
	if !isFriendlySpellHit(ctx) || ctx.Source == nil || ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) {
		return nil
	}
	addGeneratedCardToPlayerHand(ctx, ctx.OpponentID, "2001102")
	ctx.Source.UsedThisTurn++
	return nil
}

func isLearnedRushSkillThisTurn(e *Engine, skill *CardInstance) bool {
	return e != nil &&
		skill != nil &&
		skill.Card != nil &&
		skill.Card.IsSkill() &&
		skill.EnterTurn == e.State.TurnNumber &&
		cardHasRush(skill)
}

type Card3021107ArcaneShield struct{ AlwaysActive }

func (Card3021107ArcaneShield) ID() string   { return "3021107" }
func (Card3021107ArcaneShield) Name() string { return "奥能护盾" }
func (Card3021107ArcaneShield) OnSpellCast(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModDelayedShieldGain,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		Amount:           1,
	})
	return nil
}

type Card3121109FlameFlash struct{ AlwaysActive }

func (Card3121109FlameFlash) ID() string   { return "3121109" }
func (Card3121109FlameFlash) Name() string { return "烈焰闪" }
func (Card3121109FlameFlash) OnSpellHit(ctx *EffectContext) error {
	ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementFire: 3})
	return nil
}

type Card3221103WaterMirrorWall struct{ AlwaysActive }

func (Card3221103WaterMirrorWall) ID() string   { return "3221103" }
func (Card3221103WaterMirrorWall) Name() string { return "水镜壁" }
func (Card3221103WaterMirrorWall) OnDefend(ctx *EffectContext) error {
	if ctx.ExtraData == nil {
		return nil
	}
	success, _ := ctx.ExtraData["defense_success"].(bool)
	if !success {
		return nil
	}
	ctx.Engine.gainPlayerShield(ctx.PlayerID, 1)
	return nil
}

type Card3221105CorrosiveFlow struct{ AlwaysActive }

func (Card3221105CorrosiveFlow) ID() string   { return "3221105" }
func (Card3221105CorrosiveFlow) Name() string { return "腐蚀之流" }
func (Card3221105CorrosiveFlow) OnSpellHit(ctx *EffectContext) error {
	ctx.Engine.discardRandomHandCard(ctx.OpponentID)
	return nil
}

type Card3221110PlunderingTide struct{ AlwaysActive }

func (Card3221110PlunderingTide) ID() string   { return "3221110" }
func (Card3221110PlunderingTide) Name() string { return "劫掠之潮" }
func (Card3221110PlunderingTide) OnSpellHitBeforeDamage(ctx *EffectContext) error {
	if ctx.ExtraData == nil {
		return nil
	}
	affected, _ := ctx.ExtraData["affected_units"].([]*CardInstance)
	hitUnits := 0
	for _, unit := range affected {
		if unit != nil {
			hitUnits++
		}
	}
	if hitUnits == 0 {
		return nil
	}
	for i := 0; i < hitUnits; i++ {
		ctx.Engine.discardRandomHandCard(ctx.OpponentID)
	}
	ctx.Engine.drawCards(ctx.PlayerID, hitUnits)
	return nil
}

type Card3321108CallSpiritGoshawk struct{ AlwaysActive }

func (Card3321108CallSpiritGoshawk) ID() string   { return "3321108" }
func (Card3321108CallSpiritGoshawk) Name() string { return "唤灵术 苍鹰" }
func (Card3321108CallSpiritGoshawk) OnEnter(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlySkills(ctx.PlayerID, func(skill *CardInstance) bool {
		return skill != nil && skill.Card != nil && skill.Card.Category == model.ElementAir
	})
	if len(candidates) == 0 {
		return nil
	}
	applyBuff := func(targetID string) {
		if targetID == "" {
			return
		}
		ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
			Type:             TempModSkillPowerBonus,
			SourceCardNumber: ctx.Source.Card.Number,
			SourceName:       ctx.Source.Card.Name,
			TargetInstanceID: targetID,
			Amount:           1,
			RemainingUses:    1,
		})
		ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
			Type:             TempModNextSkillUseAttackBonus,
			SourceCardNumber: ctx.Source.Card.Number,
			SourceName:       ctx.Source.Card.Name,
			TargetInstanceID: targetID,
			Amount:           1,
			RemainingUses:    1,
		})
	}
	if len(candidates) == 1 {
		id, _ := candidates[0]["instance_id"].(string)
		applyBuff(id)
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "goshawk_air_skill_buff",
		"唤灵术 苍鹰:选择1个友方大气法术下一次使用时+1攻+1威", candidates, 1, 1,
		func(selected []string) {
			applyBuff(firstSelected(selected))
		})
	return nil
}

type Card3321110AirFlow struct{ AlwaysActive }

func (Card3321110AirFlow) ID() string   { return "3321110" }
func (Card3321110AirFlow) Name() string { return "气蕴成流" }
func (Card3321110AirFlow) OnEnter(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModNextLearnedSkillHaste,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		Element:          model.ElementAir,
		RemainingUses:    1,
	})
	return nil
}

type Card3421105AgingTouch struct{ AlwaysActive }

func (Card3421105AgingTouch) ID() string   { return "3421105" }
func (Card3421105AgingTouch) Name() string { return "苍老之触" }
func (Card3421105AgingTouch) OnSpellHit(ctx *EffectContext) error {
	if ctx.Target == nil || ctx.Target.Card == nil || !ctx.Target.Card.IsCompanion() {
		return nil
	}
	setElementsGain(ctx.Target, map[string]int{})
	ctx.Target.ElementsGainBonus = make(map[string]int)
	return nil
}

type Card3521110LightSpiritDrain struct{ AlwaysActive }

func (Card3521110LightSpiritDrain) ID() string   { return "3521110" }
func (Card3521110LightSpiritDrain) Name() string { return "光灵汲取" }
func (Card3521110LightSpiritDrain) OnSpellHit(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementLight
	})
	if len(candidates) == 0 {
		return nil
	}
	applyLoad := func(instanceID string) {
		target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, instanceID)
		if target == nil || zone != "unit" || target.Card == nil || !target.Card.IsCompanion() || target.Card.Category != model.ElementLight {
			return
		}
		ctx.Engine.addElementsGainBonus(target, ctx.PlayerID, model.ElementLight, 1, ctx.Source)
	}
	if len(candidates) == 1 {
		id, _ := candidates[0]["instance_id"].(string)
		applyLoad(id)
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "light_spirit_drain_load",
		"光灵汲取:选择1个友方光辉伙伴获得负载+1光", candidates, 1, 1,
		func(selected []string) {
			applyLoad(firstSelected(selected))
		})
	return nil
}

type Card3621103BloodSoulSlash struct{ AlwaysActive }

func (Card3621103BloodSoulSlash) ID() string   { return "3621103" }
func (Card3621103BloodSoulSlash) Name() string { return "血魂斩" }
func (Card3621103BloodSoulSlash) OnSpellCast(ctx *EffectContext) error {
	hero := ctx.Engine.State.Players[ctx.PlayerID].Hero
	if hero != nil {
		ctx.Engine.dealDamageWithExtra(hero, 1, ctx.PlayerID, map[string]any{
			"damage_source": "blood_soul_slash",
			"attacker":      ctx.PlayerID,
		})
	}
	return nil
}
func (Card3621103BloodSoulSlash) OnSpellHit(ctx *EffectContext) error {
	hero := ctx.Engine.State.Players[ctx.PlayerID].Hero
	if hero != nil {
		healUnit(hero, 2)
	}
	return nil
}

type Card3621101BloodPledge struct{ AlwaysActive }

func (Card3621101BloodPledge) ID() string   { return "3621101" }
func (Card3621101BloodPledge) Name() string { return "歃血" }
func (Card3621101BloodPledge) OnSpellHit(ctx *EffectContext) error {
	if ctx.Source == nil || ctx.ExtraData == nil {
		return nil
	}
	friendlyDamage, _ := ctx.ExtraData["actual_friendly_damage_by_instance"].(map[string]int)
	totalDamage := 0
	for _, amount := range friendlyDamage {
		totalDamage += amount
	}
	if totalDamage <= 0 {
		return nil
	}
	ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementShadow: 2})
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModSkillPowerBonus,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		TargetInstanceID: ctx.Source.InstanceID,
		Amount:           2,
		RemainingUses:    1,
	})
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModNextSkillUseAttackBonus,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		TargetInstanceID: ctx.Source.InstanceID,
		Amount:           1,
		RemainingUses:    1,
	})
	return nil
}

func findSkillSlotCard(ps *PlayerState, instanceID string) *CardInstance {
	if ps == nil {
		return nil
	}
	for _, skill := range ps.Skills {
		if skill != nil && skill.InstanceID == instanceID {
			return skill
		}
	}
	return nil
}

type Card1321106SilverleafRanger struct{ AlwaysActive }

func (Card1321106SilverleafRanger) ID() string   { return "1321106" }
func (Card1321106SilverleafRanger) Name() string { return "银叶游侠" }
func (Card1321106SilverleafRanger) PerTurnLabel(*CardInstance) string {
	return "主动"
}
func (Card1321106SilverleafRanger) OnPerTurn(ctx *EffectContext) error {
	if !ctx.Engine.canConsumeCard(ctx.Source) {
		return fmt.Errorf("银叶游侠不能被消耗")
	}
	ctx.Source.IsHorizontal = true
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:          TempModSkillAttackBonus,
		Amount:        1,
		RemainingUses: 1,
		ExpiresTurn:   ctx.Engine.State.TurnNumber + 2,
	})
	return nil
}

type Card1321103LoneStarTowerWatcher struct{ AlwaysActive }

func (Card1321103LoneStarTowerWatcher) ID() string   { return "1321103" }
func (Card1321103LoneStarTowerWatcher) Name() string { return "孤星塔守望者" }
func (Card1321103LoneStarTowerWatcher) OnUltimate(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "lone_star_tower_watcher_discard",
		"孤星塔守望者:丢弃至多3张手牌并获得等量护盾", candidates, 0, min(3, len(candidates)),
		func(selected []string) {
			discarded := ctx.Engine.discardSelectedHandCards(ctx.PlayerID, selected, 3)
			if discarded > 0 {
				ctx.Engine.gainPlayerShield(ctx.PlayerID, discarded)
			}
		})
	return nil
}

type Card1321109StormHorn struct{ AlwaysActive }

func (Card1321109StormHorn) ID() string   { return "1321109" }
func (Card1321109StormHorn) Name() string { return "风暴之角" }
func (Card1321109StormHorn) OnUltimate(ctx *EffectContext) error {
	handCandidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, nil)
	if len(handCandidates) == 0 || !ctx.Engine.hasAirEquipmentInDeck(ctx.PlayerID) {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "storm_horn_discard",
		"风暴之角:丢弃1张手牌", handCandidates, 1, 1,
		func(selected []string) {
			if ctx.Engine.discardSelectedHandCards(ctx.PlayerID, selected, 1) != 1 {
				return
			}
			searchDeckToHandByPredicate(ctx, "storm_horn_search_air_equipment", "风暴之角:翻取1张大气装备", isAirEquipment)
		})
	return nil
}

type Card1321113CouncilMessenger struct{ AlwaysActive }

func (Card1321113CouncilMessenger) ID() string   { return "1321113" }
func (Card1321113CouncilMessenger) Name() string { return "议庭传信鸽" }
func (Card1321113CouncilMessenger) OnEnter(ctx *EffectContext) error {
	addGeneratedCardToPlayerHand(ctx, ctx.OpponentID, "2001102")
	return nil
}

type Card1021115JiuxiaoAssassin struct{ AlwaysActive }

func (Card1021115JiuxiaoAssassin) ID() string   { return "1021115" }
func (Card1021115JiuxiaoAssassin) Name() string { return "九霄刺客" }
func (Card1021115JiuxiaoAssassin) OnEnter(ctx *EffectContext) error {
	addGeneratedCardToPlayerHand(ctx, ctx.OpponentID, "2001102")
	return nil
}
func (Card1021115JiuxiaoAssassin) OnDeath(ctx *EffectContext) error {
	addGeneratedCardsToPlayerDeck(ctx, ctx.OpponentID, "2001102", 4)
	return nil
}

type Card1321112JiuxiaoContact struct{ AlwaysActive }

func (Card1321112JiuxiaoContact) ID() string            { return "1321112" }
func (Card1321112JiuxiaoContact) Name() string          { return "九霄接头人" }
func (Card1321112JiuxiaoContact) IsPrayerAbility() bool { return true }
func (Card1321112JiuxiaoContact) OnPerTurn(ctx *EffectContext) error {
	opponent := ctx.Engine.State.Players[ctx.OpponentID]
	if len(opponent.Hand) < ctx.Engine.handLimitForPlayer(opponent) {
		addGeneratedCardToPlayerHand(ctx, ctx.OpponentID, "2001102")
	}
	return nil
}

type Card1321114CouncilExecutor struct{ AlwaysActive }

func (Card1321114CouncilExecutor) ID() string   { return "1321114" }
func (Card1321114CouncilExecutor) Name() string { return "议庭执行者" }
func (Card1321114CouncilExecutor) OnEnter(ctx *EffectContext) error {
	first := ctx.Engine.discardRandomHandCard(ctx.OpponentID)
	if first != nil && first.Card != nil && first.Card.Number == "2001102" {
		ctx.Engine.discardRandomHandCard(ctx.OpponentID)
	}
	return nil
}

type Card1521110CouncilSpeaker struct{ AlwaysActive }

func (Card1521110CouncilSpeaker) ID() string   { return "1521110" }
func (Card1521110CouncilSpeaker) Name() string { return "议庭言客" }
func (Card1521110CouncilSpeaker) OnEnter(ctx *EffectContext) error {
	addGeneratedCardsToPlayerDeck(ctx, ctx.OpponentID, "2001102", 4)
	return nil
}
func (Card1521110CouncilSpeaker) OnDeath(ctx *EffectContext) error {
	ctx.Engine.moveDeckCardToTop(ctx.OpponentID, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.Number == "2001102"
	})
	return nil
}

type Card1521115LoneStarIronKnight struct{ AlwaysActive }

func (Card1521115LoneStarIronKnight) ID() string   { return "1521115" }
func (Card1521115LoneStarIronKnight) Name() string { return "孤星铁骑士" }
func (Card1521115LoneStarIronKnight) OnEnter(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ctx.Source == nil || ctx.Source.Position == nil || ps == nil || ctx.Source.Position.Row != ps.GetFrontRow() || len(adjacentFriendlyCompanions(ctx)) > 0 {
		return nil
	}
	ctx.Source.CurrentLife++
	ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementLight, 1, ctx.Source)
	return nil
}

type Card1511102LoneStarSoul struct{ AlwaysActive }

func (Card1511102LoneStarSoul) ID() string   { return "1511102" }
func (Card1511102LoneStarSoul) Name() string { return "孤星之魂 凯拉莫将军" }
func (Card1511102LoneStarSoul) OnDamaged(ctx *EffectContext) error {
	if ctx == nil || ctx.Source == nil || ctx.Target != nil || ctx.ExtraData == nil {
		return nil
	}
	attacker, hasAttacker := ctx.ExtraData["attacker"].(int)
	if !hasAttacker || attacker == ctx.PlayerID {
		return nil
	}
	if len(adjacentFriendlyCompanions(ctx)) > 0 {
		return nil
	}
	ctx.Engine.gainPlayerShield(ctx.PlayerID, 1)
	ctx.Source.CurrentAttack++
	return nil
}

type Card1421105InactiveRoot struct{ AlwaysActive }

func (Card1421105InactiveRoot) ID() string            { return "1421105" }
func (Card1421105InactiveRoot) Name() string          { return "失活的根须" }
func (Card1421105InactiveRoot) IsPrayerAbility() bool { return true }
func (Card1421105InactiveRoot) OnPerTurn(ctx *EffectContext) error {
	if totalLoad(ctx.Source) == 0 {
		ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementEarth, 1, ctx.Source)
	}
	return nil
}

type Card1121113LavaFortHellhound struct{ AlwaysActive }

func (Card1121113LavaFortHellhound) ID() string   { return "1121113" }
func (Card1121113LavaFortHellhound) Name() string { return "熔岩堡地狱犬" }
func (Card1121113LavaFortHellhound) OnConsume(ctx *EffectContext) error {
	if ctx == nil || ctx.Source == nil || ctx.Engine == nil || ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) || ctx.ExtraData == nil {
		return nil
	}
	if ctx.Target != nil && ctx.Target != ctx.Source {
		return nil
	}
	if source, _ := ctx.ExtraData["consume_source"].(string); source == "" || source == ctx.Source.Card.Number {
		return nil
	}
	candidates := companionSpellRangeCandidates(ctx, false)
	if len(candidates) < 2 {
		return nil
	}
	ctx.Source.UsedThisTurn++
	ctx.Engine.SetPendingAction(ctx.PlayerID, "lava_fort_hellhound_damage",
		"熔岩堡地狱犬:选择法力范围内2个不同单位各造成1点伤害", candidates, 2, 2,
		func(selected []string) {
			seen := map[string]bool{}
			for _, id := range selected {
				if seen[id] {
					continue
				}
				seen[id] = true
				target := ctx.Engine.findUnitByInstanceID(id)
				if target == nil || target.Card == nil || !target.Card.IsCompanion() || target.Position == nil {
					continue
				}
				if target.OwnerID != ctx.PlayerID && !ctx.Engine.IsInSpellRange(ctx.PlayerID, target.Position.Col, target.Position.Row, false) {
					continue
				}
				ctx.Engine.dealDamageWithExtra(target, 1, target.OwnerID, map[string]any{"damage_source": "effect", "attacker": ctx.PlayerID})
			}
		})
	return nil
}

type Card1421108CelticDeer struct{ AlwaysActive }

func (Card1421108CelticDeer) ID() string   { return "1421108" }
func (Card1421108CelticDeer) Name() string { return "凯尔特灵鹿" }
func (Card1421108CelticDeer) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Source == nil || ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) || ctx.Target == nil || ctx.Target.Card == nil {
		return nil
	}
	if !ctx.Target.Card.IsSkill() || !hasCardTag(ctx.Target.Card, "灵媒") {
		return nil
	}
	resetCard(ctx.Source)
	ctx.Source.UsedThisTurn++
	return nil
}

type Card1521114HuiPrayer struct{ AlwaysActive }

func (Card1521114HuiPrayer) ID() string   { return "1521114" }
func (Card1521114HuiPrayer) Name() string { return "辉之都祈祷者" }
func (Card1521114HuiPrayer) OnEnter(ctx *EffectContext) error {
	wounded := 0
	for _, unit := range royalFriendlyUnits(ctx) {
		if unit != nil && unit.CurrentLife < maxLife(unit) {
			wounded++
		}
	}
	if wounded > 0 {
		ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementLight: wounded})
	}
	return nil
}

type Card1521106ChurchExorcist struct{ AlwaysActive }

func (Card1521106ChurchExorcist) ID() string   { return "1521106" }
func (Card1521106ChurchExorcist) Name() string { return "教廷驱魔师" }
func (Card1521106ChurchExorcist) OnEnter(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, hasAnyNegativeStatus)
	candidates = append(candidates, ctx.Engine.friendlyEquipment(ctx.PlayerID, hasAnyNegativeStatus)...)
	candidates = append(candidates, ctx.Engine.friendlySkillsIncludingBound(ctx.PlayerID, hasAnyNegativeStatus)...)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "church_exorcist_purify",
		"教廷驱魔师:选择1张友方卡牌移除全部负面状态", candidates, 1, 1,
		func(selected []string) {
			target := ctx.Engine.findFriendlyCardIncludingBound(ctx.PlayerID, firstSelected(selected))
			removed := countNegativeStatusLayers(target)
			if removed <= 0 {
				return
			}
			clearNegativeStatuses(target)
			ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementLight: removed})
		})
	return nil
}

type Card2021107Reshape struct{ AlwaysActive }

func (Card2021107Reshape) ID() string   { return "2021107" }
func (Card2021107Reshape) Name() string { return "重塑" }
func (Card2021107Reshape) OnUseItem(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	for _, card := range ps.Hand {
		if card == nil {
			continue
		}
		ctx.Engine.discardHandCardToGraveyard(ctx.PlayerID, card)
	}
	ps.Hand = nil
	ctx.Engine.drawCards(ctx.PlayerID, 2)
	return nil
}

type Card2021104FiveColorCoral struct{ AlwaysActive }

func (Card2021104FiveColorCoral) ID() string   { return "2021104" }
func (Card2021104FiveColorCoral) Name() string { return "五色珊瑚" }
func (Card2021104FiveColorCoral) OnEnter(ctx *EffectContext) error {
	choices := elementChoiceCandidates("2021104", model.ElementFire, model.ElementWater, model.ElementEarth, model.ElementAir, model.ElementLight, model.ElementShadow)
	ctx.Engine.SetPendingAction(ctx.PlayerID, "five_color_coral_load",
		"五色珊瑚:选择2种不同的非奥术元素各获得1点负载", choices, 2, 2,
		func(selected []string) {
			seen := make(map[string]bool, len(selected))
			for _, elem := range selected {
				if isNonArcaneElement(elem) && !seen[elem] {
					ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, elem, 1, ctx.Source)
					seen[elem] = true
				}
			}
		})
	return nil
}

type Card2121108BurnoutScroll struct{ AlwaysActive }

func (Card2121108BurnoutScroll) ID() string   { return "2121108" }
func (Card2121108BurnoutScroll) Name() string { return "燃烬卷轴" }
func (Card2121108BurnoutScroll) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return isFireCompanion(card) && ctx.Engine.canConsumeCard(card)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "burnout_scroll_consume_fire_companion",
		"燃烬卷轴:消耗1个友方火焰伙伴并获得其入场花费的元素", candidates, 1, 1,
		func(selected []string) {
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target == nil || zone != "unit" || !isFireCompanion(target) || !ctx.Engine.canConsumeCard(target) {
				return
			}
			gains := copyElementCost(target.Card.ElementsCost)
			target.IsHorizontal = true
			if len(gains) > 0 {
				ctx.Engine.State.Players[ctx.PlayerID].GainElements(gains)
			}
			ctx.Engine.emit(GameEvent{
				Type:   "consume",
				Player: -1,
				Data: map[string]any{
					"player":      ctx.PlayerID,
					"instance_id": target.InstanceID,
					"elements":    ctx.Engine.State.Players[ctx.PlayerID].Elements,
					"gained":      gains,
				},
			})
			consumeData := map[string]any{
				"consumed_player": ctx.PlayerID,
				"consume_source":  "2121108",
				"gained":          gains,
			}
			ctx.Engine.triggerEffects(TriggerOnConsume, target, nil, consumeData)
			ctx.Engine.triggerFieldEffectsWithData(TriggerOnConsume, ctx.PlayerID, target, consumeData)
			ctx.Engine.triggerFieldEffectsWithData(TriggerOnConsume, ctx.OpponentID, target, consumeData)
			ctx.Engine.advanceMastery(target, ctx.PlayerID, 1)
			ctx.Engine.destroyFuyeDoomedCardAfterExert(target)
		})
	return nil
}

type royalInfusionRune struct {
	AlwaysActive
	id          string
	name        string
	powerBonus  int
	attackBonus int
}

func (r royalInfusionRune) ID() string   { return r.id }
func (r royalInfusionRune) Name() string { return r.name }
func (r royalInfusionRune) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlySkillsIncludingBound(ctx.PlayerID, func(skill *CardInstance) bool {
		return skill != nil && skill.Card != nil && isSpellLikeCard(skill.Card)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "royal_infusion_rune_skill",
		r.name+":选择你的1个法术永久强化", candidates, 1, 1,
		func(selected []string) {
			skill := ctx.Engine.findSkill(ctx.Engine.State.Players[ctx.PlayerID], firstSelected(selected))
			if skill == nil || skill.Card == nil || !isSpellLikeCard(skill.Card) {
				return
			}
			skill.PowerBonus += r.powerBonus
			skill.AttackBonus += r.attackBonus
			ctx.Engine.refreshPendingSpellPowerForModifiedSkill(ctx.PlayerID, skill)
		})
	return nil
}

type Card2021101LostSilverleaf struct{ AlwaysActive }

func (Card2021101LostSilverleaf) ID() string   { return "2021101" }
func (Card2021101LostSilverleaf) Name() string { return "失落的银叶花" }
func (Card2021101LostSilverleaf) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.drawCards(ctx.PlayerID, 2)
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "lost_silverleaf_discard",
		"失落的银叶花:弃1张手牌", candidates, 1, 1,
		func(selected []string) {
			if len(selected) > 0 {
				ctx.Engine.discardFriendlyCandidate(ctx.PlayerID, selected[0])
			}
		})
	return nil
}

type Card2221105BlackSailRaider struct{ AlwaysActive }

func (Card2221105BlackSailRaider) ID() string   { return "2221105" }
func (Card2221105BlackSailRaider) Name() string { return "掠夺者黑帆" }
func (Card2221105BlackSailRaider) OnUseItem(ctx *EffectContext) error {
	hasRaiderOnField := len(ctx.Engine.friendlyUnits(ctx.PlayerID, false, isRaiderCompanion)) > 0
	searchDeckToHandByPredicateWithResult(ctx,
		"black_sail_raider_search",
		"掠夺者黑帆:检索1个掠夺者伙伴",
		isRaiderCompanion,
		func(card *CardInstance) {
			if !hasRaiderOnField || card == nil {
				return
			}
			choices := []map[string]any{
				{"instance_id": model.ElementWater, "number": "2221105", "name": "入场花费-1水", "type": "选择", "zone": "choice", "side": "own"},
				{"instance_id": model.ElementShadow, "number": "2221105", "name": "入场花费-1暗", "type": "选择", "zone": "choice", "side": "own"},
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "black_sail_raider_discount",
				"掠夺者黑帆:选择检索卡牌的入场花费减免元素", choices, 1, 1,
				func(selected []string) {
					elem := firstSelected(selected)
					if elem != model.ElementWater && elem != model.ElementShadow {
						return
					}
					if !cardInstanceInSlice(ctx.Engine.State.Players[ctx.PlayerID].Hand, card) {
						return
					}
					card.Statuses["入场费用"+elem+"-1"]++
				})
		})
	return nil
}

func isRaiderCompanion(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.IsCompanion() && strings.Contains(card.Card.Name, "掠夺者")
}

func cardInstanceInSlice(cards []*CardInstance, target *CardInstance) bool {
	for _, card := range cards {
		if card == target {
			return true
		}
	}
	return false
}

func (e *Engine) equipmentInOwnerSlot(playerID int, target *CardInstance) bool {
	if e == nil || target == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return false
	}
	for _, card := range e.State.Players[playerID].Equipment {
		if card == target {
			return true
		}
	}
	return false
}

type royalWaterUseCostReduction struct {
	AlwaysActive
	id             string
	name           string
	requireWater   bool
	triggerOnEnter bool
}

func (r royalWaterUseCostReduction) ID() string   { return r.id }
func (r royalWaterUseCostReduction) Name() string { return r.name }
func (r royalWaterUseCostReduction) HasActiveUseItem(*CardInstance) bool {
	return !r.triggerOnEnter
}
func (r royalWaterUseCostReduction) HasActiveOnEnter(*CardInstance) bool {
	return r.triggerOnEnter
}
func (r royalWaterUseCostReduction) OnUseItem(ctx *EffectContext) error {
	return r.prompt(ctx)
}
func (r royalWaterUseCostReduction) OnEnter(ctx *EffectContext) error {
	if !r.triggerOnEnter {
		return nil
	}
	return r.prompt(ctx)
}
func (r royalWaterUseCostReduction) prompt(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlySkillsIncludingBound(ctx.PlayerID, func(skill *CardInstance) bool {
		return skill != nil && skill.Card != nil && isSpellLikeCard(skill.Card) && (!r.requireWater || skill.Card.Category == model.ElementWater)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "royal_water_use_cost_reduction",
		r.name+":选择你的1个法术使用花费-1水", candidates, 1, 1,
		func(selected []string) {
			skill := ctx.Engine.findSkill(ctx.Engine.State.Players[ctx.PlayerID], firstSelected(selected))
			if skill == nil || skill.Card == nil || !isSpellLikeCard(skill.Card) || (r.requireWater && skill.Card.Category != model.ElementWater) {
				return
			}
			skill.Statuses["使用费用"+model.ElementWater+"-1"]++
		})
	return nil
}

type Card2421103Dreamcatcher struct{ AlwaysActive }

func (Card2421103Dreamcatcher) ID() string   { return "2421103" }
func (Card2421103Dreamcatcher) Name() string { return "捕梦网" }
func (Card2421103Dreamcatcher) OnEnter(ctx *EffectContext) error {
	for _, skill := range ctx.Engine.State.Players[ctx.PlayerID].Skills {
		if skill != nil && skill.Card != nil && isSpellLikeCard(skill.Card) && hasCardTag(skill.Card, "灵媒") {
			skill.PowerBonus += 2
			ctx.Engine.refreshPendingSpellPowerForModifiedSkill(ctx.PlayerID, skill)
		}
	}
	return nil
}

type Card2421109CaveElfPickaxe struct{ AlwaysActive }

func (Card2421109CaveElfPickaxe) ID() string   { return "2421109" }
func (Card2421109CaveElfPickaxe) Name() string { return "地穴精灵矿镐" }
func (Card2421109CaveElfPickaxe) PerTurnLabel(*CardInstance) string {
	return "消耗"
}
func (Card2421109CaveElfPickaxe) OnPerTurn(ctx *EffectContext) error {
	if !ctx.Engine.canConsumeCard(ctx.Source) {
		return fmt.Errorf("地穴精灵矿镐不能被消耗")
	}
	if !ctx.Engine.equipmentInOwnerSlot(ctx.PlayerID, ctx.Source) {
		return fmt.Errorf("地穴精灵矿镐必须从装备区发动")
	}
	choices := []map[string]any{
		{"instance_id": "companion", "number": "2421109", "name": "伙伴", "type": "选择", "zone": "choice", "side": "own"},
		{"instance_id": "item", "number": "2421109", "name": "道具", "type": "选择", "zone": "choice", "side": "own"},
	}
	ctx.Source.IsHorizontal = true
	ctx.Engine.SetPendingAction(ctx.PlayerID, "cave_elf_pickaxe_kind",
		"地穴精灵矿镐:选择翻取伙伴或道具", choices, 1, 1,
		func(selected []string) {
			switch firstSelected(selected) {
			case "companion":
				ctx.Engine.flipDeckMatchesToHand(ctx.PlayerID, 1, 5, func(card *CardInstance) bool {
					return card != nil && card.Card != nil && card.Card.IsCompanion()
				})
			case "item":
				ctx.Engine.flipDeckMatchesToHand(ctx.PlayerID, 1, 5, func(card *CardInstance) bool {
					return card != nil && card.Card != nil && card.Card.IsItem()
				})
			}
		})
	return nil
}

type Card2621111DarkBurstScroll struct{ AlwaysActive }

func (Card2621111DarkBurstScroll) ID() string   { return "2621111" }
func (Card2621111DarkBurstScroll) Name() string { return "暗黑爆发卷轴" }
func (Card2621111DarkBurstScroll) OnUseItem(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	targets := make([]*CardInstance, 0)
	for _, card := range ps.Graveyard {
		if card != nil && card.Card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementShadow {
			targets = append(targets, card)
		}
	}
	if len(targets) < 5 {
		return nil
	}
	exiled := 0
	for _, card := range targets {
		if ctx.Engine.exileCard(ctx.PlayerID, card) {
			exiled++
		}
	}
	if exiled > 0 {
		ps.GainElements(map[string]int{model.ElementShadow: exiled * 2})
	}
	return nil
}

type Card2621109ElegyScroll struct{ AlwaysActive }

func (Card2621109ElegyScroll) ID() string   { return "2621109" }
func (Card2621109ElegyScroll) Name() string { return "哀歌卷轴" }
func (Card2621109ElegyScroll) OnUseItem(ctx *EffectContext) error {
	hasShadowGrave := countShadowCompanionsInGraveyard(ctx.Engine.State.Players[ctx.PlayerID]) > 0
	drawn := ctx.Engine.flipDeckMatchesToHand(ctx.PlayerID, 1, 0, isShadowCompanionWithDeathrattle)
	if hasShadowGrave && len(drawn) > 0 {
		drawn[0].Statuses["入场费用"+model.ElementShadow+"-1"]++
	}
	return nil
}

func countShadowCompanionsInGraveyard(ps *PlayerState) int {
	if ps == nil {
		return 0
	}
	count := 0
	for _, card := range ps.Graveyard {
		if card != nil && card.Card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementShadow {
			count++
		}
	}
	return count
}

type Card2321102WindCycle struct{ AlwaysActive }

func (Card2321102WindCycle) ID() string   { return "2321102" }
func (Card2321102WindCycle) Name() string { return "风之轮回" }
func (Card2321102WindCycle) PerTurnLabel(*CardInstance) string {
	return "主动"
}
func (Card2321102WindCycle) OnPerTurn(ctx *EffectContext) error {
	if !ctx.Engine.canConsumeCard(ctx.Source) {
		return fmt.Errorf("风之轮回不能被消耗")
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	inEquipment := false
	for _, card := range ps.Equipment {
		if card == ctx.Source {
			inEquipment = true
			break
		}
	}
	if !inEquipment {
		return fmt.Errorf("风之轮回必须从装备区献祭")
	}
	candidates := make([]map[string]any, 0)
	allowed := make(map[string]bool)
	for _, card := range ps.Graveyard {
		if card != nil && card.Card != nil && card.Card.Category == model.ElementAir {
			candidates = append(candidates, candidateInfo(card, "graveyard", "own"))
			allowed[card.InstanceID] = true
		}
	}
	ctx.Source.IsHorizontal = true
	if !ctx.Engine.sacrificeEquipment(ctx.PlayerID, ctx.Source.InstanceID) {
		return fmt.Errorf("风之轮回必须从装备区献祭")
	}
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "wind_cycle_shuffle_air",
		"风之轮回:选择任意数量的大气弃牌洗回卡组", candidates, 0, len(candidates),
		func(selected []string) {
			selectedSet := make(map[string]bool, len(selected))
			for _, id := range selected {
				if allowed[id] {
					selectedSet[id] = true
				}
			}
			if len(selectedSet) == 0 {
				return
			}
			for i := 0; i < len(ps.Graveyard); {
				card := ps.Graveyard[i]
				if card != nil && selectedSet[card.InstanceID] && card.Card != nil && card.Card.Category == model.ElementAir {
					ps.Graveyard = append(ps.Graveyard[:i], ps.Graveyard[i+1:]...)
					ps.Deck = append(ps.Deck, card)
					continue
				}
				i++
			}
			ctx.Engine.shuffleDeck(ctx.PlayerID)
		})
	return nil
}

type Card2321103ThunderBreath struct{ AlwaysActive }

func (Card2321103ThunderBreath) ID() string   { return "2321103" }
func (Card2321103ThunderBreath) Name() string { return "雷鸣之息" }
func (Card2321103ThunderBreath) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementAir: 1})
	return nil
}

type Card2521102MoonlightDust struct{ AlwaysActive }

func (Card2521102MoonlightDust) ID() string   { return "2521102" }
func (Card2521102MoonlightDust) Name() string { return "月霞之尘" }
func (Card2521102MoonlightDust) OnUseItem(ctx *EffectContext) error {
	choices := make([]map[string]any, 0, 2)
	if ctx.Engine.hasEnemySetCounter(ctx.PlayerID) {
		choices = append(choices, map[string]any{"instance_id": "destroy_counters", "name": "摧毁敌方盖放的所有卡牌", "zone": "choice", "side": "own"})
	}
	if ctx.Engine.hasEnemyFrontStealth(ctx.PlayerID) {
		choices = append(choices, map[string]any{"instance_id": "remove_front_stealth", "name": "使前排敌人失去隐蔽", "zone": "choice", "side": "own"})
	}
	if len(choices) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "moonlight_dust_mode",
		"月霞之尘:选择1项效果", choices, 1, 1,
		func(selected []string) {
			switch firstSelected(selected) {
			case "destroy_counters":
				ctx.Engine.destroyEnemySetCounters(ctx.PlayerID)
			case "remove_front_stealth":
				ctx.Engine.removeEnemyFrontStealth(ctx.PlayerID)
			}
		})
	return nil
}

type Card4611101BloodCountHubert struct{ AlwaysActive }

func (Card4611101BloodCountHubert) ID() string   { return "4611101" }
func (Card4611101BloodCountHubert) Name() string { return "鲜血伯爵 休伯特 黑松" }
func (Card4611101BloodCountHubert) OnEnter(ctx *EffectContext) error {
	addSkillToPool(ctx, "3601101")
	return nil
}

type Card4611102CalamityRoseDom struct{ AlwaysActive }

func (Card4611102CalamityRoseDom) ID() string   { return "4611102" }
func (Card4611102CalamityRoseDom) Name() string { return "灾厄玫瑰 多姆" }
func (Card4611102CalamityRoseDom) OnEnter(ctx *EffectContext) error {
	ctx.Engine.millTopDeckCards(ctx.PlayerID, 4)
	ctx.Engine.millTopDeckCards(ctx.OpponentID, 4)
	return nil
}

type Card1321107SkyCityThief struct{ AlwaysActive }

func (Card1321107SkyCityThief) ID() string   { return "1321107" }
func (Card1321107SkyCityThief) Name() string { return "云霄城大盗" }
func (Card1321107SkyCityThief) OnEnter(ctx *EffectContext) error {
	ctx.Engine.discardRandomHandCard(ctx.PlayerID)
	ctx.Engine.discardRandomHandCard(ctx.OpponentID)
	return nil
}

type Card1621103BloodPuppet struct{ AlwaysActive }

func (Card1621103BloodPuppet) ID() string   { return "1621103" }
func (Card1621103BloodPuppet) Name() string { return "鲜血傀儡" }
func (Card1621103BloodPuppet) OnEnter(ctx *EffectContext) error {
	ctx.Engine.dealDamage(ctx.Engine.State.Players[ctx.PlayerID].Hero, 2, ctx.PlayerID)
	return nil
}

type Card1521103LoneStarGuardianSpirit struct{ AlwaysActive }

func (Card1521103LoneStarGuardianSpirit) ID() string   { return "1521103" }
func (Card1521103LoneStarGuardianSpirit) Name() string { return "孤星城的守护灵" }
func (Card1521103LoneStarGuardianSpirit) OnEnter(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion()
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "lone_star_guardian_life",
		"孤星城的守护灵:选择1个友方伙伴+1血", candidates, 1, 1,
		func(selected []string) {
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target != nil && zone == "unit" && target.Card != nil && target.Card.IsCompanion() {
				target.CurrentLife++
			}
		})
	return nil
}
func (Card1521103LoneStarGuardianSpirit) OnDeath(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion()
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "lone_star_guardian_load",
		"孤星城的守护灵:选择1个友方伙伴负载+1光", candidates, 1, 1,
		func(selected []string) {
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target != nil && zone == "unit" && target.Card != nil && target.Card.IsCompanion() {
				ctx.Engine.addElementsGainBonus(target, ctx.PlayerID, model.ElementLight, 1, ctx.Source)
			}
		})
	return nil
}

type Card1521108ContradictoryKnight struct{ AlwaysActive }

func (Card1521108ContradictoryKnight) ID() string   { return "1521108" }
func (Card1521108ContradictoryKnight) Name() string { return "矛盾的骑士" }
func (Card1521108ContradictoryKnight) OnDeath(ctx *EffectContext) error {
	opponentID := ctx.OpponentID
	candidates := ctx.Engine.friendlyEmptyUnitPositions(opponentID)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(opponentID, "contradictory_knight_summon",
		"矛盾的骑士:选择位置为你召唤此卡", candidates, 1, 1,
		func(selected []string) {
			pos, ok := positionFromSelectionID(firstSelected(selected))
			if !ok {
				return
			}
			opponent := ctx.Engine.State.Players[opponentID]
			if opponent.Units[pos.Col][pos.Row] != nil {
				return
			}
			if !ctx.Engine.removeCardFromGraveyard(ctx.PlayerID, ctx.Source) {
				return
			}
			if ctx.Source.Card.Life > 1 {
				cardCopy := *ctx.Source.Card
				cardCopy.Life--
				ctx.Source.Card = &cardCopy
			}
			ctx.Source.OwnerID = opponentID
			ctx.Source.CurrentLife = ctx.Source.Card.Life
			ctx.Source.CurrentAttack = ctx.Source.Card.Attack
			ctx.Source.DamageTakenThisTurn = 0
			ctx.Source.IsHorizontal = true
			ctx.Source.Position = nil
			ctx.Source.Statuses = make(map[string]int)
			ctx.Source.ElementsGainBonus = make(map[string]int)
			ctx.Source.ElementsGainSet = nil
			ctx.Source.PowerBonus = 0
			ctx.Source.AttackBonus = 0
			ctx.Source.UsedThisTurn = 0
			ctx.Source.UltimateUsed = false
			ctx.Source.BoundSkills = nil
			ctx.Source.UnderCards = nil
			ctx.Source.AttachedBehaviors = nil
			if !ctx.Engine.placeExistingCompanionAtPosition(opponentID, ctx.Source, pos, true) {
				ctx.Engine.addToGraveyard(ctx.PlayerID, ctx.Source)
			}
		})
	return nil
}

type Card1521113RadiantWatchdog struct{ AlwaysActive }

func (Card1521113RadiantWatchdog) ID() string   { return "1521113" }
func (Card1521113RadiantWatchdog) Name() string { return "辉之都戒卫犬" }
func (Card1521113RadiantWatchdog) OnDeath(ctx *EffectContext) error {
	if ctx.ExtraData == nil {
		return nil
	}
	attacker, ok := ctx.ExtraData["attacker"].(int)
	if !ok || attacker == ctx.PlayerID {
		return nil
	}
	candidates := ctx.Engine.friendlyDeckCards(ctx.PlayerID, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion()
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "radiant_watchdog_search",
		"辉之都戒卫犬:翻取1个伙伴牌并使其入场花费-1光", candidates, 1, 1,
		func(selected []string) {
			card := ctx.Engine.searchDeckCardToHand(ctx.PlayerID, firstSelected(selected))
			if card != nil {
				card.Statuses["入场费用"+model.ElementLight+"-1"]++
			}
		})
	return nil
}

type Card1621112WhisperElfHunter struct{ AlwaysActive }

func (Card1621112WhisperElfHunter) ID() string   { return "1621112" }
func (Card1621112WhisperElfHunter) Name() string { return "谧语精灵猎手" }
func (Card1621112WhisperElfHunter) OnDeath(ctx *EffectContext) error {
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, true, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "whisper_elf_hunter_damage",
		"谧语精灵猎手:选择1个敌人造成1点伤害", candidates, 1, 1,
		func(selected []string) {
			target := ctx.Engine.findUnitByInstanceID(firstSelected(selected))
			if target != nil {
				ctx.Engine.dealDamage(target, 1, target.OwnerID)
			}
		})
	return nil
}

type Card1621113WhisperElfPriest struct{ AlwaysActive }

func (Card1621113WhisperElfPriest) ID() string   { return "1621113" }
func (Card1621113WhisperElfPriest) Name() string { return "谧语精灵祭司" }
func (Card1621113WhisperElfPriest) OnDeath(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion()
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "whisper_elf_priest_load",
		"谧语精灵祭司:选择1个友方伙伴负载+1暗", candidates, 1, 1,
		func(selected []string) {
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target != nil && zone == "unit" && target.Card != nil && target.Card.IsCompanion() {
				ctx.Engine.addElementsGainBonus(target, ctx.PlayerID, model.ElementShadow, 1, ctx.Source)
			}
		})
	return nil
}

type Card1621114SoulSymbiote struct{ AlwaysActive }

const soulMarkerStatus = "灵魂标记物"

type Card1621106SoulHunter struct{ AlwaysActive }

func (Card1621106SoulHunter) ID() string   { return "1621106" }
func (Card1621106SoulHunter) Name() string { return "猎魂者" }
func (Card1621106SoulHunter) OnSpellHit(ctx *EffectContext) error {
	if !isFriendlySpellHit(ctx) || ctx.Source == nil || ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) {
		return nil
	}
	skill := ctx.Target
	if ctx.ExtraData != nil {
		if source, ok := ctx.ExtraData["spell_source"].(*CardInstance); ok && source != nil {
			skill = source
		}
	}
	if skill == nil || skill.Card == nil || !isSpellLikeCard(skill.Card) {
		return nil
	}
	skill.Statuses[soulMarkerStatus]++
	skill.PowerBonus += 2
	ctx.Source.UsedThisTurn++
	return nil
}

func (Card1621114SoulSymbiote) ID() string   { return "1621114" }
func (Card1621114SoulSymbiote) Name() string { return "灵魂共生体" }
func (Card1621114SoulSymbiote) OnDeath(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlySkillsIncludingBound(ctx.PlayerID, func(skill *CardInstance) bool {
		return skill != nil && skill.Card != nil && skill.Card.IsSkill()
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "soul_symbiote_mark_skills",
		"灵魂共生体:选择最多2个法术放置灵魂标记物", candidates, 0, 2,
		func(selected []string) {
			seen := map[string]bool{}
			for _, id := range selected {
				if seen[id] {
					continue
				}
				seen[id] = true
				skill := ctx.Engine.findSkill(ctx.Engine.State.Players[ctx.PlayerID], id)
				if skill == nil || skill.Card == nil || !skill.Card.IsSkill() {
					continue
				}
				skill.Statuses[soulMarkerStatus]++
				skill.PowerBonus += 2
			}
		})
	return nil
}

type Card1621101PainSoul struct{ AlwaysActive }

func (Card1621101PainSoul) ID() string   { return "1621101" }
func (Card1621101PainSoul) Name() string { return "苦痛之魂" }
func (Card1621101PainSoul) OnDamaged(ctx *EffectContext) error {
	if ctx.Source == nil || ctx.Target != nil || ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) {
		return nil
	}
	ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementShadow, 1, ctx.Source)
	ctx.Source.UsedThisTurn++
	return nil
}

type Card1621102PainAvenger struct{ AlwaysActive }

func (Card1621102PainAvenger) ID() string   { return "1621102" }
func (Card1621102PainAvenger) Name() string { return "苦痛复仇者" }
func (Card1621102PainAvenger) OnDamaged(ctx *EffectContext) error {
	if ctx.Source == nil || ctx.Target != nil || ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) {
		return nil
	}
	ctx.Source.CurrentAttack++
	ctx.Source.UsedThisTurn++
	return nil
}

type Card1621104RoseGardenGardener struct{ AlwaysActive }

func (Card1621104RoseGardenGardener) ID() string   { return "1621104" }
func (Card1621104RoseGardenGardener) Name() string { return "蔷薇花园园丁" }
func (Card1621104RoseGardenGardener) OnFriendlyDeath(ctx *EffectContext) error {
	if ctx.Source == nil || ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) {
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card != nil && card.CurrentLife < maxLife(card)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "rose_garden_gardener_heal",
		"蔷薇花园园丁:选择1个友方单位回复2血", candidates, 1, 1,
		func(selected []string) {
			target := ctx.Engine.findUnitByInstanceID(firstSelected(selected))
			if target != nil && target.OwnerID == ctx.PlayerID && target.CurrentLife < maxLife(target) && ctx.Source.UsedThisTurn < perTurnLimit(ctx.Source) {
				healUnit(target, 2)
				ctx.Source.UsedThisTurn++
			}
		})
	return nil
}

const curseBoxMarkerStatus = "诅咒魔盒标记物"

type Card2621107CurseBox struct{ AlwaysActive }

func (Card2621107CurseBox) ID() string   { return "2621107" }
func (Card2621107CurseBox) Name() string { return "诅咒魔盒" }
func (Card2621107CurseBox) OnFriendlyDeath(ctx *EffectContext) error {
	return addCurseBoxMarker(ctx)
}
func (Card2621107CurseBox) OnEnemyDeath(ctx *EffectContext) error {
	return addCurseBoxMarker(ctx)
}
func (Card2621107CurseBox) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Source == nil || ctx.Engine == nil {
		return nil
	}
	markers := ctx.Source.Statuses[curseBoxMarkerStatus]
	if markers <= 0 {
		return fmt.Errorf("诅咒魔盒没有标记物")
	}
	candidates := ctx.Engine.enemySkills(ctx.PlayerID, canInstanceBeWeakened)
	maxSelect := min(3, min(markers, len(candidates)))
	if maxSelect <= 0 {
		return fmt.Errorf("没有可虚弱的敌方法术")
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "curse_box_weaken",
		"诅咒魔盒:移除最多3个标记物，使等量敌方法术虚弱1", candidates, 1, maxSelect,
		func(selected []string) {
			removed := 0
			seen := map[string]bool{}
			for _, id := range selected {
				if seen[id] || ctx.Source.Statuses[curseBoxMarkerStatus] <= 0 {
					continue
				}
				seen[id] = true
				for _, skill := range ctx.Engine.State.Players[ctx.OpponentID].Skills {
					if skill != nil && skill.InstanceID == id && canInstanceBeWeakened(skill) {
						ctx.Engine.addStatus(skill, StatusWeaken, 1)
						ctx.Source.Statuses[curseBoxMarkerStatus]--
						removed++
						break
					}
				}
			}
			if ctx.Source.Statuses[curseBoxMarkerStatus] <= 0 {
				delete(ctx.Source.Statuses, curseBoxMarkerStatus)
			}
			if removed == 0 && ctx.Source.Statuses[curseBoxMarkerStatus] <= 0 {
				delete(ctx.Source.Statuses, curseBoxMarkerStatus)
			}
		})
	return nil
}

func addCurseBoxMarker(ctx *EffectContext) error {
	if ctx == nil || ctx.Source == nil {
		return nil
	}
	ctx.Source.Statuses[curseBoxMarkerStatus]++
	return nil
}

type Card2201101DreamBloom struct{ AlwaysActive }

func (Card2201101DreamBloom) ID() string   { return "2201101" }
func (Card2201101DreamBloom) Name() string { return "幻创之梦-绽放" }
func (Card2201101DreamBloom) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.drawCards(ctx.PlayerID, 3)
	return nil
}

type Card2201102DreamMana struct{ AlwaysActive }

func (Card2201102DreamMana) ID() string   { return "2201102" }
func (Card2201102DreamMana) Name() string { return "幻创之梦-幻能" }
func (Card2201102DreamMana) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementArcane: 3})
	return nil
}

type Card3621110BloodNourish struct{ AlwaysActive }

func (Card3621110BloodNourish) ID() string   { return "3621110" }
func (Card3621110BloodNourish) Name() string { return "鲜血滋养" }
func (Card3621110BloodNourish) OnSpellCast(ctx *EffectContext) error {
	if !isFriendlySpellCast(ctx) || !isSpellBeingCast(ctx) {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	candidates := make([]map[string]any, 0)
	for _, card := range ps.Graveyard {
		if card != nil && card.Card != nil && card.Card.Category == model.ElementShadow {
			candidates = append(candidates, candidateInfo(card, "graveyard", "own"))
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "blood_nourish_exile",
		"鲜血滋养:选择弃牌堆1张暗影卡牌移出游戏，获得2暗", candidates, 1, 1,
		func(selected []string) {
			for _, card := range ps.Graveyard {
				if card != nil && card.InstanceID == firstSelected(selected) && card.Card != nil && card.Card.Category == model.ElementShadow {
					if ctx.Engine.exileCard(ctx.PlayerID, card) {
						ps.GainElements(map[string]int{model.ElementShadow: 2})
					}
					return
				}
			}
		})
	return nil
}

type Card2021116ArcaneBomb struct{ AlwaysActive }

func (Card2021116ArcaneBomb) ID() string   { return "2021116" }
func (Card2021116ArcaneBomb) Name() string { return "奥能炸弹" }
func (Card2021116ArcaneBomb) OnUseItem(ctx *EffectContext) error {
	candidates := companionSpellRangeCandidates(ctx, false)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "arcane_bomb_target",
		"奥能炸弹:选择法力范围内1个伙伴造成2点伤害", candidates, 1, 1,
		func(selected []string) {
			target := ctx.Engine.findUnitByInstanceID(firstSelected(selected))
			if target != nil && target.Card != nil && target.Card.IsCompanion() {
				ctx.Engine.dealDamage(target, 2, target.OwnerID)
			}
		})
	return nil
}

type Card2511101JiuxiaoRadiance struct{ AlwaysActive }

func (Card2511101JiuxiaoRadiance) ID() string   { return "2511101" }
func (Card2511101JiuxiaoRadiance) Name() string { return "九霄辉迹" }
func (Card2511101JiuxiaoRadiance) OnUltimate(ctx *EffectContext) error {
	counts := make([]int, len(ctx.Engine.State.Players))
	for playerID, ps := range ctx.Engine.State.Players {
		counts[playerID] = len(ps.Hand)
	}
	for playerID := range ctx.Engine.State.Players {
		ctx.Engine.discardAllHandCards(playerID)
	}
	for playerID, count := range counts {
		if count > 0 {
			ctx.Engine.drawCards(playerID, count)
		}
	}
	return nil
}

type Card2521104GoldenDragonbone struct{ AlwaysActive }

func (Card2521104GoldenDragonbone) ID() string   { return "2521104" }
func (Card2521104GoldenDragonbone) Name() string { return "黄金龙骨" }
func (Card2521104GoldenDragonbone) PerTurnLabel(*CardInstance) string {
	return "主动"
}
func (Card2521104GoldenDragonbone) OnPerTurn(ctx *EffectContext) error {
	if !ctx.Engine.sacrificeEquipment(ctx.PlayerID, ctx.Source.InstanceID) {
		return fmt.Errorf("golden dragonbone must be sacrificed from equipment")
	}
	ctx.Engine.drawCards(ctx.PlayerID, 2)
	return nil
}

type Card2421112AutumnMapleGem struct{ AlwaysActive }

const autumnMapleGemCounter = "秋枫宝钻标记物"

func (Card2421112AutumnMapleGem) ID() string   { return "2421112" }
func (Card2421112AutumnMapleGem) Name() string { return "秋枫宝钻" }
func (Card2421112AutumnMapleGem) OnEnter(ctx *EffectContext) error {
	ctx.Source.Statuses[autumnMapleGemCounter] += 2
	return nil
}
func (Card2421112AutumnMapleGem) PerTurnLabel(*CardInstance) string {
	return "回合技"
}
func (Card2421112AutumnMapleGem) OnPerTurn(ctx *EffectContext) error {
	if ctx.Source.Statuses[autumnMapleGemCounter] <= 0 {
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementEarth && card.IsHorizontal
	})
	if len(candidates) == 0 {
		return nil
	}
	allowed := make(map[string]bool, len(candidates))
	for _, candidate := range candidates {
		if id, _ := candidate["instance_id"].(string); id != "" {
			allowed[id] = true
		}
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "autumn_maple_gem_reset",
		"秋枫宝钻:选择1个地脉伙伴重置", candidates, 1, 1,
		func(selected []string) {
			id := firstSelected(selected)
			if !allowed[id] || ctx.Source.Statuses[autumnMapleGemCounter] <= 0 {
				return
			}
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, id)
			if target == nil || zone != "unit" || target.Card == nil || !target.Card.IsCompanion() || target.Card.Category != model.ElementEarth {
				return
			}
			ctx.Source.Statuses[autumnMapleGemCounter]--
			target.IsHorizontal = false
		})
	return nil
}

type Card2521101BlessedLoneStar struct{ AlwaysActive }

func (Card2521101BlessedLoneStar) ID() string   { return "2521101" }
func (Card2521101BlessedLoneStar) Name() string { return "赐福之孤星" }
func (Card2521101BlessedLoneStar) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion()
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "blessed_lone_star_target",
		"赐福之孤星:选择1个友方伙伴获得负载+1光和+1血", candidates, 1, 1,
		func(selected []string) {
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target == nil || zone != "unit" || target.Card == nil || !target.Card.IsCompanion() {
				return
			}
			target.CurrentLife++
			ctx.Engine.addElementsGainBonus(target, ctx.PlayerID, model.ElementLight, 1, ctx.Source)
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source":  cardToInfo(ctx.Source),
				"target":  cardToInfo(target),
				"effect":  "life_and_load",
				"life":    1,
				"element": model.ElementLight,
				"amount":  1,
			}})
		})
	return nil
}

type Card2521106MoonlightScroll struct{ AlwaysActive }

func (Card2521106MoonlightScroll) ID() string   { return "2521106" }
func (Card2521106MoonlightScroll) Name() string { return "沐光卷轴" }
func (Card2521106MoonlightScroll) OnUseItem(ctx *EffectContext) error {
	for _, unit := range royalFriendlyUnits(ctx) {
		healUnit(unit, 2)
	}
	return nil
}

type Card2421108EmeraldFruit struct{ AlwaysActive }

func (Card2421108EmeraldFruit) ID() string   { return "2421108" }
func (Card2421108EmeraldFruit) Name() string { return "翡翠果" }
func (Card2421108EmeraldFruit) OnEnter(ctx *EffectContext) error {
	targets := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion()
	})
	if len(targets) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "emerald_fruit_target",
		"翡翠果:选择1个友方伙伴获得负载", targets, 1, 1,
		func(selected []string) {
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target == nil || zone != "unit" || target.Card == nil || !target.Card.IsCompanion() {
				return
			}
			choices := elementChoiceCandidates("2421108", model.ElementFire, model.ElementWater, model.ElementAir, model.ElementLight, model.ElementShadow)
			ctx.Engine.SetPendingAction(ctx.PlayerID, "emerald_fruit_element",
				"翡翠果:选择除地与奥术外的1点负载", choices, 1, 1,
				func(selected []string) {
					elem := firstSelected(selected)
					if elem != model.ElementEarth && isNonArcaneElement(elem) {
						ctx.Engine.addElementsGainBonus(target, ctx.PlayerID, elem, 1, ctx.Source)
					}
				})
		})
	return nil
}

func royalFriendlyUnits(ctx *EffectContext) []*CardInstance {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	units := make([]*CardInstance, 0, 9)
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if unit := ps.Units[col][row]; unit != nil {
				units = append(units, unit)
			}
		}
	}
	return units
}

func royalCompanionCount(ps *PlayerState) int {
	if ps == nil {
		return 0
	}
	count := 0
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := ps.Units[col][row]
			if unit != nil && unit.Card != nil && unit.Card.IsCompanion() {
				count++
			}
		}
	}
	return count
}

func addGeneratedCardToPlayerHand(ctx *EffectContext, playerID int, cardNumber string) *CardInstance {
	card := getCardDB()[cardNumber]
	if card == nil {
		return nil
	}
	instance := NewCardInstance(card, playerID, ctx.Engine.State.TurnNumber)
	ctx.Engine.State.Players[playerID].Hand = append(ctx.Engine.State.Players[playerID].Hand, instance)
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(ctx.Source),
		"card":   cardToInfo(instance),
		"effect": "add_generated_card_to_hand",
	}})
	return instance
}

func addGeneratedCardsToPlayerDeck(ctx *EffectContext, playerID int, cardNumber string, count int) []*CardInstance {
	card := getCardDB()[cardNumber]
	if card == nil || count <= 0 {
		return nil
	}
	ps := ctx.Engine.State.Players[playerID]
	added := make([]*CardInstance, 0, count)
	for i := 0; i < count; i++ {
		instance := NewCardInstance(card, playerID, ctx.Engine.State.TurnNumber)
		ps.Deck = append(ps.Deck, instance)
		added = append(added, instance)
	}
	ctx.Engine.shuffleDeck(playerID)
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(ctx.Source),
		"cards":  cardsToInfo(added),
		"effect": "add_generated_cards_to_deck",
	}})
	return added
}

func (e *Engine) discardRandomHandCard(playerID int) *CardInstance {
	ps := e.State.Players[playerID]
	if ps == nil || len(ps.Hand) == 0 {
		return nil
	}
	idx := rand.Intn(len(ps.Hand))
	return e.discardHandCardAt(playerID, idx)
}

func (e *Engine) discardSelectedHandCards(playerID int, selected []string, limit int) int {
	if playerID < 0 || playerID >= len(e.State.Players) || limit <= 0 {
		return 0
	}
	ps := e.State.Players[playerID]
	selectedSet := map[string]bool{}
	for _, id := range selected {
		if id != "" {
			selectedSet[id] = true
		}
	}
	discarded := 0
	for i := len(ps.Hand) - 1; i >= 0 && discarded < limit; i-- {
		card := ps.Hand[i]
		if card == nil || !selectedSet[card.InstanceID] {
			continue
		}
		if e.discardHandCardAt(playerID, i) != nil {
			discarded++
		}
	}
	return discarded
}

func (e *Engine) discardAllHandCards(playerID int) int {
	if playerID < 0 || playerID >= len(e.State.Players) {
		return 0
	}
	ps := e.State.Players[playerID]
	discarded := 0
	for len(ps.Hand) > 0 {
		if e.discardHandCardAt(playerID, len(ps.Hand)-1) != nil {
			discarded++
		}
	}
	return discarded
}

func (e *Engine) hasAirEquipmentInDeck(playerID int) bool {
	if playerID < 0 || playerID >= len(e.State.Players) {
		return false
	}
	for _, card := range e.State.Players[playerID].Deck {
		if isAirEquipment(card) && canFlipOrSearchCard(card) {
			return true
		}
	}
	return false
}

func (e *Engine) resolveDiscardedCardEffects(playerID int, card *CardInstance) {
	if card == nil || card.Card == nil {
		return
	}
	if card.Card.Number == "2001102" {
		if hero := e.playerHeroCard(playerID); hero != nil {
			e.dealDamage(hero, 2, playerID)
		}
	}
	if card.Card.Number == "2321103" {
		e.State.Players[playerID].GainElements(map[string]int{model.ElementAir: 1})
	}
}

func (e *Engine) playerHeroCard(playerID int) *CardInstance {
	ps := e.State.Players[playerID]
	if ps == nil {
		return nil
	}
	if ps.Hero != nil {
		return ps.Hero
	}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			card := ps.Units[col][row]
			if card != nil && card.Card != nil && card.Card.IsHero() {
				return card
			}
		}
	}
	return nil
}

func (e *Engine) moveDeckCardToTop(playerID int, predicate func(*CardInstance) bool) *CardInstance {
	ps := e.State.Players[playerID]
	if ps == nil {
		return nil
	}
	for i, card := range ps.Deck {
		if card == nil || (predicate != nil && !predicate(card)) {
			continue
		}
		ps.Deck = append(ps.Deck[:i], ps.Deck[i+1:]...)
		ps.Deck = append([]*CardInstance{card}, ps.Deck...)
		e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
			"card":   cardToInfo(card),
			"effect": "deck_card_to_top",
		}})
		return card
	}
	return nil
}

func (e *Engine) hasEnemySetCounter(playerID int) bool {
	opponent := e.State.Players[1-playerID]
	if opponent == nil {
		return false
	}
	for _, card := range opponent.Equipment {
		if card != nil && card.IsSetCounter {
			return true
		}
	}
	return false
}

func (e *Engine) destroyEnemySetCounters(playerID int) {
	opponentID := 1 - playerID
	opponent := e.State.Players[opponentID]
	if opponent == nil {
		return
	}
	for i := range opponent.Equipment {
		card := opponent.Equipment[i]
		if card != nil && card.IsSetCounter {
			e.moveEquipmentToGraveyard(opponentID, i, card)
		}
	}
}

func (e *Engine) hasEnemyFrontStealth(playerID int) bool {
	opponent := e.State.Players[1-playerID]
	if opponent == nil {
		return false
	}
	frontRow := opponent.GetFrontRow()
	if frontRow < 0 || frontRow >= 3 {
		return false
	}
	for col := 0; col < 3; col++ {
		unit := opponent.Units[col][frontRow]
		if unit != nil && unit.Statuses[StatusStealth] > 0 {
			return true
		}
	}
	return false
}

func (e *Engine) removeEnemyFrontStealth(playerID int) {
	opponent := e.State.Players[1-playerID]
	if opponent == nil {
		return
	}
	frontRow := opponent.GetFrontRow()
	if frontRow < 0 || frontRow >= 3 {
		return
	}
	for col := 0; col < 3; col++ {
		unit := opponent.Units[col][frontRow]
		if unit != nil {
			delete(unit.Statuses, StatusStealth)
		}
	}
}

func countNegativeStatusLayers(card *CardInstance) int {
	if card == nil {
		return 0
	}
	total := 0
	for _, status := range negativeStatuses {
		if card.Statuses[status] > 0 {
			total += card.Statuses[status]
		}
	}
	return total
}

func (e *Engine) findFriendlyCardIncludingBound(playerID int, instanceID string) *CardInstance {
	if card, _ := e.findFriendlyCandidate(playerID, instanceID); card != nil {
		return card
	}
	ps := e.State.Players[playerID]
	if ps == nil {
		return nil
	}
	for _, host := range e.getAllFieldCards(ps) {
		if host == nil {
			continue
		}
		for _, skill := range host.BoundSkills {
			if skill != nil && skill.InstanceID == instanceID {
				return skill
			}
		}
	}
	return nil
}

func elementChoiceCandidates(sourceNumber string, elements ...string) []map[string]any {
	candidates := make([]map[string]any, 0, len(elements))
	for _, elem := range elements {
		candidates = append(candidates, map[string]any{
			"instance_id": elem,
			"number":      sourceNumber,
			"name":        elem,
			"type":        "元素",
			"zone":        "choice",
			"side":        "own",
		})
	}
	return candidates
}

func isNonArcaneElement(elem string) bool {
	return elem == model.ElementFire || elem == model.ElementWater || elem == model.ElementEarth || elem == model.ElementAir || elem == model.ElementLight || elem == model.ElementShadow
}

func fireButterflyTemporaryLoadStillCurrent(card *CardInstance) bool {
	if card == nil || card.ElementsGainSet == nil {
		return false
	}
	if card.ElementsGainSet[model.ElementAir] != 1 {
		return false
	}
	for _, elem := range model.AllElements {
		if elem == model.ElementAir {
			continue
		}
		if card.ElementsGainSet[elem] != 0 {
			return false
		}
	}
	return true
}

func clearFireButterflyStoredLoad(card *CardInstance) {
	if card == nil {
		return
	}
	delete(card.Statuses, fireButterflyTemporaryLoadStatus)
	delete(card.Statuses, fireButterflyPreviousLoadSetStatus)
	for _, elem := range model.AllElements {
		delete(card.Statuses, fireButterflyPreviousLoadValuePrefix+elem)
	}
}

func (e *Engine) hasResettableWaterSpell(playerID int) bool {
	return len(e.friendlySkillsIncludingBound(playerID, func(skill *CardInstance) bool {
		return skill != nil && skill.Card != nil && skill.Card.IsSkill() &&
			skill.Card.Category == model.ElementWater &&
			totalElementCost(skill.Card.ElementsExpense) < 3 &&
			skill.IsHorizontal
	})) > 0
}

func (e *Engine) hasResettableEarthCompanion(playerID int) bool {
	return len(e.friendlyUnits(playerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() &&
			card.Card.Category == model.ElementEarth &&
			card.IsHorizontal
	})) > 0
}

func (e *Engine) removeCardFromGraveyard(playerID int, card *CardInstance) bool {
	if card == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return false
	}
	ps := e.State.Players[playerID]
	for i, candidate := range ps.Graveyard {
		if candidate == card {
			ps.Graveyard = append(ps.Graveyard[:i], ps.Graveyard[i+1:]...)
			return true
		}
	}
	return false
}

func (e *Engine) placeExistingCompanionAtPosition(playerID int, card *CardInstance, pos Position, triggerEnter bool) bool {
	if card == nil || card.Card == nil || !card.Card.IsCompanion() || !pos.Valid() || playerID < 0 || playerID >= len(e.State.Players) {
		return false
	}
	ps := e.State.Players[playerID]
	if ps.Units[pos.Col][pos.Row] != nil {
		return false
	}
	card.OwnerID = playerID
	card.Position = &Position{Col: pos.Col, Row: pos.Row}
	card.EnterTurn = e.State.TurnNumber
	ps.Units[pos.Col][pos.Row] = card
	e.ApplySummonModifiersOnEnter(card)
	if triggerEnter {
		e.triggerEffects(TriggerOnEnter, card, nil, nil)
		e.triggerFieldEffectsWithData(TriggerOnUnitEnter, playerID, card, map[string]any{"entered_player": playerID})
		e.triggerFieldEffectsWithData(TriggerOnUnitEnter, 1-playerID, card, map[string]any{"entered_player": playerID})
	}
	return true
}

func (e *Engine) summonFreshCardAtPosition(playerID int, cardNumber string, pos Position, triggerEnter bool) *CardInstance {
	cardDef := getCardDB()[cardNumber]
	if cardDef == nil {
		return nil
	}
	instance := NewCardInstance(cardDef, playerID, e.State.TurnNumber)
	if !e.placeExistingCompanionAtPosition(playerID, instance, pos, triggerEnter) {
		return nil
	}
	return instance
}

func adjacentFriendlyCompanions(ctx *EffectContext) []map[string]any {
	if ctx.Source == nil || ctx.Source.Position == nil {
		return nil
	}
	candidates := make([]map[string]any, 0, 4)
	for _, unit := range adjacentUnits(ctx.Engine.State.Players[ctx.PlayerID], ctx.Source.Position) {
		if unit != nil && unit.Card != nil && unit.Card.IsCompanion() {
			candidates = append(candidates, candidateInfo(unit, "unit", "own"))
		}
	}
	return candidates
}
