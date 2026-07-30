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

type Card1421115Geomancer struct{ AlwaysActive }

func (Card1421115Geomancer) ID() string   { return "1421115" }
func (Card1421115Geomancer) Name() string { return "地卜行者" }
func (Card1421115Geomancer) OnEnter(ctx *EffectContext) error {
	ctx.Engine.drawCards(ctx.PlayerID, 1)
	return nil
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
