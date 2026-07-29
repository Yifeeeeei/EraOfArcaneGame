package game

import (
	"fmt"

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

type Card1421115Geomancer struct{ AlwaysActive }

func (Card1421115Geomancer) ID() string   { return "1421115" }
func (Card1421115Geomancer) Name() string { return "地卜行者" }
func (Card1421115Geomancer) OnEnter(ctx *EffectContext) error {
	ctx.Engine.drawCards(ctx.PlayerID, 1)
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

type Card2021107Reshape struct{ AlwaysActive }

func (Card2021107Reshape) ID() string   { return "2021107" }
func (Card2021107Reshape) Name() string { return "重塑" }
func (Card2021107Reshape) OnUseItem(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	for _, card := range ps.Hand {
		if card == nil {
			continue
		}
		delete(ps.RevealedHand, card.InstanceID)
		ps.Graveyard = append(ps.Graveyard, card)
		ctx.Engine.emit(GameEvent{Type: "discard", Player: ctx.PlayerID, Data: map[string]any{"card": cardToInfo(card)}})
	}
	ps.Hand = nil
	ctx.Engine.drawCards(ctx.PlayerID, 2)
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
