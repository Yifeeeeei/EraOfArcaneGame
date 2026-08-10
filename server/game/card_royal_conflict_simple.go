package game

import (
	"fmt"
	"math/rand"
	"strings"

	"eraofarcane/cards"
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

type Card1111101SupremeQueenDailinCeltic struct{ AlwaysActive }

func (Card1111101SupremeQueenDailinCeltic) ID() string   { return "1111101" }
func (Card1111101SupremeQueenDailinCeltic) Name() string { return "无上女王 黛琳 凯尔特" }
func (Card1111101SupremeQueenDailinCeltic) OnEnter(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Position == nil {
		return nil
	}
	markTemporaryDamageAndNegativeImmunity(ctx.Engine, ctx.Source)
	maxSummons := len(ctx.Engine.adjacentEmptyUnitPositions(ctx.PlayerID, *ctx.Source.Position))
	if maxSummons == 0 {
		return nil
	}
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementFire
	})
	if len(candidates) == 0 {
		return nil
	}
	if maxSummons > len(candidates) {
		maxSummons = len(candidates)
	}
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "supreme_queen_summon_cards",
		"无上女王 黛琳 凯尔特:选择火焰伙伴召唤到相邻位置", candidates, 0, maxSummons,
		nil, false, func(selected []string, _ map[string]any) error {
			ctx.Engine.continueSupremeQueenSummons(ctx, selected, 0)
			return nil
		})
	return nil
}

func markTemporaryDamageAndNegativeImmunity(e *Engine, card *CardInstance) {
	if e == nil || card == nil {
		return
	}
	card.Statuses[temporaryDamageAndNegativeImmunityUntilStatus] = e.State.TurnNumber + 1
}

func (e *Engine) adjacentEmptyUnitPositions(playerID int, center Position) []map[string]any {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) || !center.Valid() {
		return nil
	}
	ps := e.State.Players[playerID]
	candidates := make([]map[string]any, 0, 4)
	for _, delta := range []struct{ col, row int }{{0, -1}, {-1, 0}, {1, 0}, {0, 1}} {
		pos := Position{Col: center.Col + delta.col, Row: center.Row + delta.row}
		if !pos.Valid() || ps.Units[pos.Col][pos.Row] != nil {
			continue
		}
		candidates = append(candidates, map[string]any{
			"instance_id": positionSelectionID(pos),
			"name":        fmt.Sprintf("位置 (%d,%d)", pos.Col, pos.Row),
			"zone":        "unit_position",
			"side":        "own",
			"position":    pos,
		})
	}
	return candidates
}

func (e *Engine) continueSupremeQueenSummons(ctx *EffectContext, selected []string, index int) {
	if e == nil || ctx == nil || ctx.Source == nil || index >= len(selected) {
		return
	}
	card := e.findFriendlyHandCard(ctx.PlayerID, selected[index])
	if card == nil || card.Card == nil || !card.Card.IsCompanion() || card.Card.Category != model.ElementFire {
		e.continueSupremeQueenSummons(ctx, selected, index+1)
		return
	}
	if ctx.Source.Position == nil {
		return
	}
	positions := e.adjacentEmptyUnitPositions(ctx.PlayerID, *ctx.Source.Position)
	if len(positions) == 0 {
		return
	}
	e.SetPendingActionWithError(ctx.PlayerID, "supreme_queen_summon_position",
		fmt.Sprintf("无上女王 黛琳 凯尔特:选择%s的召唤位置", card.Card.Name), positions, 1, 1,
		nil, false, func(posSelected []string, _ map[string]any) error {
			pos, ok := positionFromSelectionID(firstSelected(posSelected))
			if !ok || ctx.Source.Position == nil || abs(pos.Col-ctx.Source.Position.Col)+abs(pos.Row-ctx.Source.Position.Row) != 1 {
				return fmt.Errorf("invalid queen summon position")
			}
			if e.State.Players[ctx.PlayerID].Units[pos.Col][pos.Row] != nil {
				return fmt.Errorf("queen summon position occupied")
			}
			card := e.removeFriendlyHandCard(ctx.PlayerID, selected[index])
			if card == nil || card.Card == nil || !card.Card.IsCompanion() || card.Card.Category != model.ElementFire {
				return fmt.Errorf("invalid queen summon card")
			}
			if !e.placeExistingCompanionAtPosition(ctx.PlayerID, card, pos, true) {
				return fmt.Errorf("queen summon failed")
			}
			markTemporaryDamageAndNegativeImmunity(e, card)
			if e.State.PendingAction != nil {
				e.wrapPendingActionContinuation(func() {
					e.continueSupremeQueenSummons(ctx, selected, index+1)
				})
				return nil
			}
			e.continueSupremeQueenSummons(ctx, selected, index+1)
			return nil
		})
}

func (e *Engine) findFriendlyHandCard(playerID int, instanceID string) *CardInstance {
	if e == nil || instanceID == "" || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	for _, card := range e.State.Players[playerID].Hand {
		if card != nil && card.InstanceID == instanceID {
			return card
		}
	}
	return nil
}

func (e *Engine) removeFriendlyHandCard(playerID int, instanceID string) *CardInstance {
	if e == nil || instanceID == "" || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	ps := e.State.Players[playerID]
	for i, card := range ps.Hand {
		if card != nil && card.InstanceID == instanceID {
			ps.Hand = append(ps.Hand[:i], ps.Hand[i+1:]...)
			delete(ps.RevealedHand, card.InstanceID)
			return card
		}
	}
	return nil
}

type Card2211102ManesArbitration struct{ AlwaysActive }

func (Card2211102ManesArbitration) ID() string   { return "2211102" }
func (Card2211102ManesArbitration) Name() string { return "玛涅斯之予夺" }
func (Card2211102ManesArbitration) OnEquip(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil {
		return nil
	}
	choices := []map[string]any{
		{"instance_id": "all_water_power", "name": "此后学习的所有水纹法术+2威", "zone": "choice", "side": "own"},
	}
	waterSkills := ctx.Engine.friendlySkills(ctx.PlayerID, func(skill *CardInstance) bool {
		return skill != nil && skill.Card != nil && skill.Card.IsSkill() && skill.Card.Category == model.ElementWater
	})
	if len(waterSkills) > 0 {
		choices = append(choices, map[string]any{
			"instance_id": "one_water_power_attack",
			"name":        "选择1个水纹法术+3威+1攻并禁止再学习水纹",
			"zone":        "choice",
			"side":        "own",
		})
	}
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "manes_arbitration_choice",
		"玛涅斯之予夺:选择本局游戏持续的效果", choices, 1, 1,
		nil, false, func(selected []string, _ map[string]any) error {
			switch firstSelected(selected) {
			case "all_water_power":
				ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
					Type:             TempModSkillPowerBonus,
					SourceCardNumber: "2211102",
					SourceName:       "玛涅斯之予夺",
					Element:          model.ElementWater,
					Amount:           2,
				})
			case "one_water_power_attack":
				currentWaterSkills := ctx.Engine.friendlySkills(ctx.PlayerID, func(skill *CardInstance) bool {
					return skill != nil && skill.Card != nil && skill.Card.IsSkill() && skill.Card.Category == model.ElementWater
				})
				if len(currentWaterSkills) == 0 {
					return fmt.Errorf("no water skill to empower")
				}
				ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "manes_arbitration_skill",
					"玛涅斯之予夺:选择1个水纹法术", currentWaterSkills, 1, 1,
					nil, false, func(skillSelected []string, _ map[string]any) error {
						skill := findFriendlySkillIncludingBound(ctx.Engine, ctx.PlayerID, firstSelected(skillSelected))
						if skill == nil || skill.Card == nil || skill.Card.Category != model.ElementWater {
							return fmt.Errorf("invalid Manes target")
						}
						skill.PowerBonus += 3
						skill.AttackBonus++
						ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
							Type:             TempModCannotLearnElementSkill,
							SourceCardNumber: "2211102",
							SourceName:       "玛涅斯之予夺",
							Element:          model.ElementWater,
						})
						return nil
					})
			default:
				return fmt.Errorf("invalid Manes choice")
			}
			return nil
		})
	return nil
}

type Card2211101DeepSword struct{ AlwaysActive }

func (Card2211101DeepSword) ID() string   { return "2211101" }
func (Card2211101DeepSword) Name() string { return "珊瑚秘宝 深邃之剑" }
func (Card2211101DeepSword) OnDraw(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	if drawnPlayer, ok := ctx.ExtraData["drawn_player"].(int); !ok || drawnPlayer != ctx.PlayerID {
		return nil
	}
	if ctx.Engine.currentLearnedSpellPower(ctx.OpponentID) <= ctx.Engine.currentLearnedSpellPower(ctx.PlayerID) {
		return nil
	}
	targets := ctx.Engine.enemyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() && card.Position != nil &&
			ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false)
	})
	if len(targets) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "deep_sword_reveal_damage",
		"珊瑚秘宝 深邃之剑:是否展示并对法力范围内所有敌人造成2点伤害",
		[]map[string]any{candidateInfo(ctx.Source, "hand", "own")}, 0, 1, func(selected []string) {
			if len(selected) == 0 || ctx.Engine.findFriendlyHandCard(ctx.PlayerID, ctx.Source.InstanceID) == nil {
				return
			}
			ps := ctx.Engine.State.Players[ctx.PlayerID]
			if ps.RevealedHand == nil {
				ps.RevealedHand = make(map[string]bool)
			}
			ps.RevealedHand[ctx.Source.InstanceID] = true
			currentTargets := ctx.Engine.enemyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
				return card != nil && card.Card != nil && card.Card.IsCompanion() && card.Position != nil &&
					ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false)
			})
			for _, info := range currentTargets {
				id, _ := info["instance_id"].(string)
				target := findEnemyCardCandidate(ctx.Engine, ctx.PlayerID, id, currentTargets)
				if target == nil || !ctx.Engine.unitStillOnField(target) {
					continue
				}
				ctx.Engine.dealDamageWithExtra(target, 2, target.OwnerID, map[string]any{
					"damage_source":  "effect",
					"damage_element": model.ElementWater,
					"source_card":    ctx.Source,
					"attacker":       ctx.PlayerID,
				})
				if target.CurrentLife <= 0 && !target.Card.IsHero() && ctx.Engine.unitInOwnerGrid(target, target.OwnerID) {
					ctx.Engine.destroyUnitWithData(target, target.OwnerID, map[string]any{
						"death_cause": "deep_sword",
						"source_card": ctx.Source,
						"attacker":    ctx.PlayerID,
					})
				}
			}
		})
	return nil
}

func (e *Engine) currentLearnedSpellPower(playerID int) int {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return 0
	}
	total := 0
	for _, skill := range e.State.Players[playerID].Skills {
		if skill == nil || skill.Card == nil || !isSpellLikeCard(skill.Card) {
			continue
		}
		total += e.effectiveSkillPowerForPurpose(playerID, skill, skillPurposeAttack)
	}
	return total
}

const flowerSeaDreamWhaleCreationCountStatus = "flower_sea_dream_whale_creation_count"

type Card1211102FlowerSeaDreamWhale struct{ AlwaysActive }

func (Card1211102FlowerSeaDreamWhale) ID() string   { return "1211102" }
func (Card1211102FlowerSeaDreamWhale) Name() string { return "花海梦鲸" }
func (Card1211102FlowerSeaDreamWhale) OnEnter(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	for _, number := range []string{"2201101", "2201102", "2201103"} {
		card := getCardDB()[number]
		if card == nil {
			continue
		}
		ps.Deck = append(ps.Deck, NewCardInstance(card, ctx.PlayerID, ctx.Engine.State.TurnNumber))
	}
	ctx.Engine.shuffleDeck(ctx.PlayerID)
	ctx.Engine.emit(GameEvent{
		Type:   "flower_sea_dream_whale_shuffle_dreams",
		Player: -1,
		Data: map[string]any{
			"player": ctx.PlayerID,
			"source": cardToInfo(ctx.Source),
			"count":  3,
		},
	})
	return nil
}
func (Card1211102FlowerSeaDreamWhale) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.ExtraData == nil {
		return nil
	}
	castPlayer, _ := ctx.ExtraData["cast_player"].(int)
	if castPlayer != ctx.PlayerID || ctx.Target.OwnerID != ctx.PlayerID || ctx.Target.Card == nil || !hasCardTag(ctx.Target.Card, "创造") {
		return nil
	}
	ctx.Source.Statuses[flowerSeaDreamWhaleCreationCountStatus]++
	if ctx.Source.Statuses[flowerSeaDreamWhaleCreationCountStatus] < 2 {
		return nil
	}
	searchDeckToHandByPredicateWithResult(ctx, "flower_sea_dream_whale_search",
		"花海梦鲸:检索1张幻创之梦", isDreamCreationCardInstance,
		func(*CardInstance) {
			ctx.Source.Statuses[flowerSeaDreamWhaleCreationCountStatus] -= 2
			if ctx.Source.Statuses[flowerSeaDreamWhaleCreationCountStatus] < 0 {
				ctx.Source.Statuses[flowerSeaDreamWhaleCreationCountStatus] = 0
			}
		})
	return nil
}

func isDreamCreationCardInstance(card *CardInstance) bool {
	return card != nil && card.Card != nil && (card.Card.Number == "2201101" || card.Card.Number == "2201102" || card.Card.Number == "2201103")
}

var _ OnEnterBehavior = Card1211102FlowerSeaDreamWhale{}
var _ OnSpellCastBehavior = Card1211102FlowerSeaDreamWhale{}

const rottenAncientTreeHeartSpellCountPrefix = "rotten_ancient_tree_heart_spell_count_p"

type Card2411102RottenAncientTreeHeart struct{ AlwaysActive }

func (Card2411102RottenAncientTreeHeart) ID() string   { return "2411102" }
func (Card2411102RottenAncientTreeHeart) Name() string { return "腐朽的古树之心" }
func (Card2411102RottenAncientTreeHeart) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil {
		return nil
	}
	castPlayer, ok := ctx.ExtraData["cast_player"].(int)
	if !ok || castPlayer < 0 || castPlayer >= len(ctx.Engine.State.Players) {
		return nil
	}
	key := fmt.Sprintf("%s%d", rottenAncientTreeHeartSpellCountPrefix, castPlayer)
	ctx.Source.Statuses[key]++
	if ctx.Source.Statuses[key] < 2 {
		return nil
	}
	candidates := loadRemovalCandidates(ctx.Engine, castPlayer, "own")
	if len(candidates) == 0 {
		return nil
	}
	removeLoad := func(selection string) {
		targetID, elem, ok := strings.Cut(selection, "|")
		if !ok || elem == "" {
			return
		}
		target := ctx.Engine.findCardOnField(ctx.Engine.State.Players[castPlayer], targetID)
		if target == nil || reducibleElementLoad(target, elem) <= 0 {
			return
		}
		ctx.Engine.reduceCardElementLoadWithTriggers(castPlayer, target, elem, 1, ctx.Source)
		ctx.Source.Statuses[key] -= 2
		if ctx.Source.Statuses[key] < 0 {
			ctx.Source.Statuses[key] = 0
		}
		ctx.Engine.emit(GameEvent{
			Type:   "rotten_ancient_tree_heart_remove_load",
			Player: -1,
			Data: map[string]any{
				"player":  castPlayer,
				"source":  cardToInfo(ctx.Source),
				"target":  cardToInfo(target),
				"element": elem,
			},
		})
		if target == ctx.Source && ctx.Engine.totalLoad(ctx.Source) <= 0 {
			ctx.Engine.sacrificeEquipment(ctx.Source.OwnerID, ctx.Source.InstanceID)
		}
	}
	if len(candidates) == 1 {
		id, _ := candidates[0]["instance_id"].(string)
		removeLoad(id)
		return nil
	}
	ctx.Engine.SetPendingAction(castPlayer, "rotten_ancient_tree_heart_remove_load",
		"腐朽的古树之心:选择自己场上1点负载移除", candidates, 1, 1,
		func(selected []string) {
			removeLoad(firstSelected(selected))
		})
	return nil
}

func loadRemovalCandidates(e *Engine, playerID int, side string) []map[string]any {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	candidates := make([]map[string]any, 0)
	for _, card := range e.getAllFieldCards(e.State.Players[playerID]) {
		if card == nil || card.Card == nil || e.hasEffectiveStatus(card, StatusPetrify) {
			continue
		}
		for elem, amount := range dragonBloodTreantReducibleLoad(card) {
			if amount <= 0 {
				continue
			}
			info := candidateInfo(card, "field", side)
			info["instance_id"] = card.InstanceID + "|" + elem
			info["load_element"] = elem
			info["name"] = fmt.Sprintf("%s - 移除%s负载", card.Card.Name, elem)
			candidates = append(candidates, info)
		}
	}
	return candidates
}

var _ OnSpellCastBehavior = Card2411102RottenAncientTreeHeart{}

const conductorConsumedCountStatus = "指挥家消耗计数"

type Card1011102ConductorLos struct{ AlwaysActive }

func (Card1011102ConductorLos) ID() string   { return "1011102" }
func (Card1011102ConductorLos) Name() string { return "\"指挥家\" 洛斯" }
func (Card1011102ConductorLos) OnConsume(ctx *EffectContext) error {
	if ctx.Source == nil || ctx.Target == nil {
		return nil
	}
	ctx.Source.Statuses[conductorConsumedCountStatus]++
	if ctx.Source.Statuses[conductorConsumedCountStatus] < 4 {
		return nil
	}
	if ctx.Engine.firstFreeEquipmentSlot(ctx.PlayerID) < 0 {
		return nil
	}
	ctx.Source.Statuses[conductorConsumedCountStatus] -= 4
	if ctx.Source.Statuses[conductorConsumedCountStatus] <= 0 {
		delete(ctx.Source.Statuses, conductorConsumedCountStatus)
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "conductor_los_equip_finale_violin",
		"\"指挥家\" 洛斯:是否装备1个落幕提琴", []map[string]any{candidateInfo(ctx.Source, "unit", "own")}, 0, 1,
		func(selected []string) {
			if len(selected) == 0 || !ctx.Engine.cardStillOnField(ctx.Source) || ctx.Engine.firstFreeEquipmentSlot(ctx.PlayerID) < 0 {
				return
			}
			ctx.Engine.equipGeneratedCard(ctx.PlayerID, "2001101")
		})
	return nil
}
func (Card1011102ConductorLos) PerTurnLabel(*CardInstance) string {
	return "消耗"
}
func (Card1011102ConductorLos) OnPerTurn(ctx *EffectContext) error {
	if ctx.Source == nil || !ctx.Engine.canConsumeCard(ctx.Source) || !ctx.Engine.cardStillOnField(ctx.Source) {
		return fmt.Errorf("\"指挥家\" 洛斯需要在场且竖置才能消耗")
	}
	ctx.Engine.consumeCardForEffectWithTriggers(ctx.PlayerID, ctx.Source, ctx.Engine.effectiveElementsGain(ctx.Source), "1011102")
	for _, equipment := range ctx.Engine.State.Players[ctx.PlayerID].Equipment {
		if equipment != nil && equipment.Card != nil && equipment.Card.Number == "2001101" {
			equipment.IsHorizontal = false
		}
	}
	return nil
}

type Card2001101FinaleViolin struct{ AlwaysActive }

func (Card2001101FinaleViolin) ID() string   { return "2001101" }
func (Card2001101FinaleViolin) Name() string { return "落幕提琴" }
func (Card2001101FinaleViolin) AttackCost(*EffectContext) map[string]int {
	return map[string]int{model.ElementArcane: 2}
}

type Card3001101EnterGame struct{ AlwaysActive }

func (Card3001101EnterGame) ID() string   { return "3001101" }
func (Card3001101EnterGame) Name() string { return "入局" }
func (Card3001101EnterGame) OnSpellCast(ctx *EffectContext) error {
	candidates := enterGamePlayerCandidates(ctx)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "enter_game_player",
		"入局:选择召唤弃子的玩家", candidates, 1, 1,
		func(selected []string) {
			targetPlayerID, ok := enterGamePlayerIDFromSelection(firstSelected(selected))
			if !ok || targetPlayerID < 0 || targetPlayerID >= len(ctx.Engine.State.Players) {
				return
			}
			positions := ctx.Engine.emptyUnitPositionsForPlayer(targetPlayerID, ctx.PlayerID)
			if len(positions) == 0 {
				return
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "enter_game_position",
				"入局:选择弃子的召唤位置", positions, 1, 1,
				func(posSelected []string) {
					pos, ok := positionFromSelectionID(firstSelected(posSelected))
					if !ok || targetPlayerID < 0 || targetPlayerID >= len(ctx.Engine.State.Players) {
						return
					}
					ps := ctx.Engine.State.Players[targetPlayerID]
					if ps == nil || ps.Units[pos.Col][pos.Row] != nil {
						return
					}
					ctx.Engine.summonFreshCardAtPosition(targetPlayerID, "1001101", pos, true)
				})
		})
	return nil
}

func enterGamePlayerCandidates(ctx *EffectContext) []map[string]any {
	if ctx == nil || ctx.Engine == nil {
		return nil
	}
	candidates := make([]map[string]any, 0, len(ctx.Engine.State.Players))
	for playerID := range ctx.Engine.State.Players {
		if len(ctx.Engine.emptyUnitPositionsForPlayer(playerID, ctx.PlayerID)) == 0 {
			continue
		}
		side := "enemy"
		if playerID == ctx.PlayerID {
			side = "own"
		}
		candidates = append(candidates, map[string]any{
			"instance_id": fmt.Sprintf("player:%d", playerID),
			"name":        fmt.Sprintf("玩家%d", playerID+1),
			"zone":        "player",
			"side":        side,
			"player_id":   playerID,
		})
	}
	return candidates
}

func enterGamePlayerIDFromSelection(id string) (int, bool) {
	var playerID int
	if _, err := fmt.Sscanf(id, "player:%d", &playerID); err != nil {
		return 0, false
	}
	return playerID, true
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

type Card1021107GiftedYouth struct{ AlwaysActive }

func (Card1021107GiftedYouth) ID() string      { return "1021107" }
func (Card1021107GiftedYouth) Name() string    { return "天才少年" }
func (Card1021107GiftedYouth) MasteryMax() int { return 2 }
func (Card1021107GiftedYouth) OnMastery(ctx *EffectContext, level int) error {
	if level != 2 {
		return nil
	}
	source := ctx.Source
	choices := elementChoiceCandidates("1021107", model.ElementFire, model.ElementWater, model.ElementEarth, model.ElementAir, model.ElementLight, model.ElementShadow)
	ctx.Engine.SetPendingAction(ctx.PlayerID, "gifted_youth_mastery_load",
		"天才少年:选择获得的非奥术负载", choices, 1, 1,
		func(selected []string) {
			elem := firstSelected(selected)
			if !isNonArcaneElement(elem) || !ctx.Engine.cardStillOnField(source) {
				return
			}
			ctx.Engine.addElementsGainBonus(source, ctx.PlayerID, elem, 1, source)
		})
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

type Card4111101ChiefAdvisorFelin struct{ AlwaysActive }

func (Card4111101ChiefAdvisorFelin) ID() string   { return "4111101" }
func (Card4111101ChiefAdvisorFelin) Name() string { return "首席顾问 费林" }
func (Card4111101ChiefAdvisorFelin) OnUltimate(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, isFireCompanion)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "felin_sacrifice_fire_companion",
		"首席顾问 费林:献祭1个友方火焰伙伴", candidates, 1, 1,
		func(selected []string) {
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target == nil || zone != "unit" || !isFireCompanion(target) {
				return
			}
			cost := copyElementCost(target.Card.ElementsCost)
			if totalElementCost(cost) <= 0 {
				return
			}
			ctx.Engine.destroyUnitWithCause(target, ctx.PlayerID, DeathCauseSacrifice)
			for _, elem := range model.AllElements {
				if cost[elem] <= 0 {
					continue
				}
				ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
					Type:          TempModNextFireCardPlayCostMinus,
					Element:       elem,
					Amount:        cost[elem],
					RemainingUses: 1,
					ExpiresTurn:   ctx.Engine.State.TurnNumber + 1,
				})
			}
		})
	return nil
}

type Card4111102GeneralKelan struct{ AlwaysActive }

func (Card4111102GeneralKelan) ID() string   { return "4111102" }
func (Card4111102GeneralKelan) Name() string { return "大将军 克兰" }
func (Card4111102GeneralKelan) OnDefend(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil {
		return nil
	}
	success, _ := ctx.ExtraData["defense_success"].(bool)
	if !success || ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) || !deckHasMatch(ctx.Engine.State.Players[ctx.PlayerID], isFireCard) {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "general_kelan_flip_fire_card",
		"大将军 克兰:是否翻取1张火焰卡牌并弃1张手牌", []map[string]any{candidateInfo(ctx.Source, "hero", "own")}, 0, 1,
		func(selected []string) {
			if len(selected) == 0 || ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) {
				return
			}
			drawn := ctx.Engine.flipDeckMatchesToHandThen(ctx.PlayerID, 1, 0, isFireCard, func(drawn []*CardInstance) {
				if len(drawn) == 0 {
					return
				}
				candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, nil)
				if len(candidates) == 0 {
					return
				}
				ctx.Engine.SetPendingAction(ctx.PlayerID, "general_kelan_discard",
					"大将军 克兰:弃1张手牌", candidates, 1, 1,
					func(discardSelected []string) {
						ctx.Engine.discardFriendlyCandidate(ctx.PlayerID, firstSelected(discardSelected))
					})
			})
			if len(drawn) == 0 {
				return
			}
			ctx.Source.UsedThisTurn++
		})
	return nil
}

func isFireCard(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.Category == model.ElementFire
}

type Card1311101SparrowSilverleaf struct{ AlwaysActive }

func (Card1311101SparrowSilverleaf) ID() string   { return "1311101" }
func (Card1311101SparrowSilverleaf) Name() string { return "斯帕罗 银叶" }
func (Card1311101SparrowSilverleaf) OnEnter(ctx *EffectContext) error {
	damage := min(ctx.Engine.State.Players[ctx.PlayerID].DiscardedHandCountThisTurn, 3)
	if damage <= 0 {
		return nil
	}
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card.Position != nil && ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "sparrow_silverleaf_entry_damage",
		"斯帕罗 银叶:选择法力范围内1名敌人造成弃牌数量伤害", candidates, 1, 1,
		func(selected []string) {
			target := findEnemyCardCandidate(ctx.Engine, ctx.PlayerID, firstSelected(selected), candidates)
			if target == nil || target.Position == nil || !ctx.Engine.IsInSpellRange(ctx.PlayerID, target.Position.Col, target.Position.Row, false) {
				return
			}
			ctx.Engine.dealDamage(target, damage, ctx.PlayerID)
		})
	return nil
}

type Card1321102SpeckledSparrow struct{ AlwaysActive }

func (Card1321102SpeckledSparrow) ID() string   { return "1321102" }
func (Card1321102SpeckledSparrow) Name() string { return "花斑麻雀" }

func (e *Engine) offerDiscardedSpeckledSparrowSummon(playerID int, card *CardInstance) {
	if e == nil || card == nil || card.Card == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	ps := e.State.Players[playerID]
	cost := map[string]int{model.ElementAir: 1}
	if !e.canPayCost(ps, cost) {
		return
	}
	positions := e.friendlyEmptyUnitPositions(playerID)
	if len(positions) == 0 {
		return
	}
	cardID := card.InstanceID
	e.SetPendingActionWithError(playerID, "speckled_sparrow_discard_summon",
		"花斑麻雀:是否支付1气召唤被弃置的此卡", positions, 0, 1, cost, false,
		func(selected []string, data map[string]any) error {
			if len(selected) == 0 {
				return nil
			}
			pos, ok := positionFromSelectionID(firstSelected(selected))
			if !ok || !pos.Valid() || e.State.Players[playerID].Units[pos.Col][pos.Row] != nil {
				return fmt.Errorf("invalid speckled sparrow position")
			}
			if !e.payCostForAction(ps, cost, ActionMessage{Data: data}) {
				return fmt.Errorf("invalid speckled sparrow payment")
			}
			if !e.reviveCompanionFromGraveyardWithLifeAtPosition(playerID, cardID, 0, false, pos) {
				return fmt.Errorf("invalid speckled sparrow summon")
			}
			return nil
		})
}

type Card1321115SkyPainter struct{ AlwaysActive }

func (Card1321115SkyPainter) ID() string   { return "1321115" }
func (Card1321115SkyPainter) Name() string { return "苍穹描摹者" }
func (Card1321115SkyPainter) OnEnter(ctx *EffectContext) error {
	candidates := friendlyFieldCardsIncludingBound(ctx.Engine, ctx.PlayerID, func(card *CardInstance) bool {
		return canSkyPainterCopy(ctx.Engine, ctx.Source, card)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "sky_painter_copy_enter",
		"苍穹描摹者:复制另一个低费大气卡牌的入场效果", candidates, 1, 1,
		func(selected []string) {
			target := ctx.Engine.findFriendlyCardIncludingBound(ctx.PlayerID, firstSelected(selected))
			if !canSkyPainterCopy(ctx.Engine, ctx.Source, target) {
				return
			}
			ctx.Engine.triggerEffects(TriggerOnEnter, target, nil, map[string]any{
				"copied_by": cardToInfo(ctx.Source),
			})
		})
	return nil
}

func canSkyPainterCopy(e *Engine, source *CardInstance, card *CardInstance) bool {
	if e == nil || card == nil || card == source || card.Card == nil ||
		card.Card.Category != model.ElementAir ||
		totalElementCost(card.Card.ElementsCost) >= 6 {
		return false
	}
	behavior, ok := behaviorForNumber(card.Card.Number).(OnEnterBehavior)
	return ok && behavior.HasActiveOnEnter(card)
}

type Card1521105RadiantCityPriest struct{ AlwaysActive }

func (Card1521105RadiantCityPriest) ID() string   { return "1521105" }
func (Card1521105RadiantCityPriest) Name() string { return "辉之都祭司" }
func (Card1521105RadiantCityPriest) OnSpellHitBeforeDamage(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil || ctx.Source.UltimateUsed {
		return nil
	}
	attacker, ok := ctx.ExtraData["attacker"].(int)
	damagePtr, hasDamage := ctx.ExtraData["damage_ptr"].(*int)
	if !ok || attacker == ctx.PlayerID || !hasDamage || damagePtr == nil || *damagePtr <= 0 {
		return nil
	}
	source := ctx.Source
	ctx.Engine.SetPendingAction(ctx.PlayerID, "radiant_city_priest_convert_damage_to_burn",
		"辉之都祭司:是否将该法术伤害转为点燃", []map[string]any{candidateInfo(source, "unit", "own")}, 0, 1,
		func(selected []string) {
			if len(selected) == 0 || source.UltimateUsed || !ctx.Engine.cardStillOnField(source) {
				return
			}
			burn := *damagePtr
			if burn <= 0 {
				return
			}
			for _, target := range spellHitAffectedUnitsFromData(ctx) {
				if target != nil && ctx.Engine.unitStillOnField(target) {
					ctx.Engine.addStatus(target, StatusBurn, burn)
				}
			}
			*damagePtr = 0
			ctx.ExtraData["damage"] = 0
			source.UltimateUsed = true
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source": cardToInfo(source),
				"effect": "convert_spell_damage_to_burn",
				"burn":   burn,
			}})
		})
	return nil
}

type Card2221106FrostRobe struct{ AlwaysActive }

func (Card2221106FrostRobe) ID() string   { return "2221106" }
func (Card2221106FrostRobe) Name() string { return "凛冰法袍" }
func (Card2221106FrostRobe) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil || ctx.Source.UltimateUsed {
		return nil
	}
	attacker, ok := ctx.ExtraData["attacker"].(int)
	if !ok || attacker == ctx.PlayerID || !friendlyWaterUnitTookSpellDamage(ctx) {
		return nil
	}
	frozen := 0
	for _, candidate := range ctx.Engine.enemyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card != nil && card.Position != nil && ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false)
	}) {
		id, _ := candidate["instance_id"].(string)
		target := findEnemyCardCandidate(ctx.Engine, ctx.PlayerID, id, []map[string]any{candidate})
		if target != nil && ctx.Engine.addStatus(target, StatusFreeze, 1) {
			frozen++
		}
	}
	if frozen == 0 {
		return nil
	}
	ctx.Source.UltimateUsed = true
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(ctx.Source),
		"effect": "frost_robe_freeze_enemies",
		"count":  frozen,
	}})
	return nil
}

func friendlyWaterUnitTookSpellDamage(ctx *EffectContext) bool {
	if ctx == nil || ctx.ExtraData == nil {
		return false
	}
	actualDamage, _ := ctx.ExtraData["actual_friendly_damage_by_instance"].(map[string]int)
	if len(actualDamage) == 0 {
		return false
	}
	for _, unit := range spellHitAffectedUnitsFromData(ctx) {
		if unit == nil || unit.OwnerID != ctx.PlayerID || unit.Card == nil || unit.Card.Category != model.ElementWater {
			continue
		}
		if actualDamage[unit.InstanceID] > 0 {
			return true
		}
	}
	return false
}

func spellHitAffectedUnitsFromData(ctx *EffectContext) []*CardInstance {
	if ctx == nil {
		return nil
	}
	if ctx.ExtraData != nil {
		if affected, ok := ctx.ExtraData["affected_units"].([]*CardInstance); ok && len(affected) > 0 {
			return affected
		}
	}
	if ctx.Target != nil {
		return []*CardInstance{ctx.Target}
	}
	return nil
}

type Card2621108BlackPineCoffin struct{ AlwaysActive }

func (Card2621108BlackPineCoffin) ID() string   { return "2621108" }
func (Card2621108BlackPineCoffin) Name() string { return "黑松棺木" }
func (Card2621108BlackPineCoffin) OnEnter(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, isLowCostShadowCompanion)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "black_pine_coffin_discard_shadow_companions",
		"黑松棺木:丢弃最多2张低费暗影伙伴并结算它们的遗言", candidates, 0, min(2, len(candidates)),
		func(selected []string) {
			discarded := ctx.Engine.discardSelectedHandCardsMatching(ctx.PlayerID, selected, 2, isLowCostShadowCompanion)
			for _, card := range discarded {
				ctx.Engine.triggerEffects(TriggerOnDeath, card, nil, map[string]any{
					"death_cause": "black_pine_coffin",
					"from_zone":   "hand",
				})
			}
		})
	return nil
}

func isLowCostShadowCompanion(card *CardInstance) bool {
	return card != nil && card.Card != nil &&
		card.Card.IsCompanion() &&
		card.Card.Category == model.ElementShadow &&
		totalElementCost(card.Card.ElementsCost) < 5
}

const permanentPierceStatus = "永久穿透"

type Card2111101DivineFireStaffCrimsonSky struct{ AlwaysActive }

func (Card2111101DivineFireStaffCrimsonSky) ID() string   { return "2111101" }
func (Card2111101DivineFireStaffCrimsonSky) Name() string { return "神火杖 赤空" }
func (Card2111101DivineFireStaffCrimsonSky) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.Card == nil {
		return nil
	}
	attacker, _ := ctx.ExtraData["attacker"].(int)
	if attacker != ctx.PlayerID || ctx.Target.Card.Category != model.ElementFire || !isSpellLikeCard(ctx.Target.Card) || !ctx.Engine.canConsumeCard(ctx.Source) {
		return nil
	}
	staff := ctx.Source
	skill := ctx.Target
	ctx.Engine.SetPendingAction(ctx.PlayerID, "divine_fire_staff_empower_spell",
		"神火杖 赤空:是否消耗此卡永久强化命中的火焰法术", []map[string]any{candidateInfo(staff, "equipment", "own")}, 0, 1,
		func(selected []string) {
			if len(selected) == 0 || !ctx.Engine.canConsumeCard(staff) || !ctx.Engine.cardStillOnField(staff) {
				return
			}
			if skill == nil || skill.Card == nil || !isSpellLikeCard(skill.Card) || skill.Card.Category != model.ElementFire {
				return
			}
			ctx.Engine.consumeCardForEffectWithTriggers(ctx.PlayerID, staff, ctx.Engine.effectiveElementsGain(staff), "2111101")
			skill.PowerBonus++
			skill.Statuses[permanentPierceStatus] = 1
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source":  cardToInfo(staff),
				"target":  cardToInfo(skill),
				"effect":  "permanent_fire_spell_empower",
				"power":   1,
				"pierce":  true,
				"consume": true,
			}})
		})
	return nil
}

type Card2111102LavaArmorYeYan struct{ AlwaysActive }

func (Card2111102LavaArmorYeYan) ID() string   { return "2111102" }
func (Card2111102LavaArmorYeYan) Name() string { return "熔岩魔甲 业炎" }

func (Card2111102LavaArmorYeYan) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil {
		return nil
	}
	attacker, _ := ctx.ExtraData["attacker"].(int)
	if attacker == ctx.PlayerID || !ctx.Engine.cardStillOnField(ctx.Source) {
		return nil
	}
	armor := ctx.Source
	ctx.Engine.SetPendingAction(ctx.PlayerID, "lava_armor_yeyan_sacrifice",
		"熔岩魔甲 业炎:是否献祭此卡获得护盾2", []map[string]any{candidateInfo(armor, "equipment", "own")}, 0, 1,
		func(selected []string) {
			if len(selected) == 0 || !ctx.Engine.cardStillOnField(armor) {
				return
			}
			slot := equipmentSlotOf(ctx.Engine.State.Players[ctx.PlayerID], armor)
			if slot < 0 {
				return
			}
			ctx.Engine.moveEquipmentToGraveyard(ctx.PlayerID, slot, armor)
			_ = (Card2111102LavaArmorYeYan{}).OnDeath(&EffectContext{
				Engine:     ctx.Engine,
				Source:     armor,
				PlayerID:   ctx.PlayerID,
				OpponentID: 1 - ctx.PlayerID,
				ExtraData:  ctx.ExtraData,
			})
			ctx.Engine.gainPlayerShield(ctx.PlayerID, 2)
		})
	return nil
}

func (Card2111102LavaArmorYeYan) OnDeath(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.PlayerID < 0 || ctx.PlayerID >= len(ctx.Engine.State.Players) {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ps == nil || !ps.ShieldBrokenThisTurn {
		return nil
	}
	equipped := ctx.Engine.equipCardFromHandOrDeckFree(ctx.PlayerID, "2121013")
	if equipped != nil {
		ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
			"source": cardToInfo(ctx.Source),
			"effect": "lava_armor_yeyan_equip_molten_armor",
			"card":   cardToInfo(equipped),
		}})
	}
	return nil
}

var _ OnSpellHitBehavior = Card2111102LavaArmorYeYan{}
var _ OnDeathBehavior = Card2111102LavaArmorYeYan{}

const (
	erebosSoulChainMarkedUnitStatus  = "erebos_soul_chain_marked_unit"
	erebosSoulChainMarkedSpellStatus = "erebos_soul_chain_marked_spell"
)

type Card2611101ErebosSoulChain struct{ AlwaysActive }

func (Card2611101ErebosSoulChain) ID() string   { return "2611101" }
func (Card2611101ErebosSoulChain) Name() string { return "厄瑞波斯的魂链" }

func (Card2611101ErebosSoulChain) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.UltimateUsed || ctx.Target == nil || ctx.Target.Card == nil || ctx.ExtraData == nil {
		return nil
	}
	castPlayer, _ := ctx.ExtraData["cast_player"].(int)
	if castPlayer == ctx.PlayerID {
		return nil
	}
	overexertUnits := spellInstancesFromData(ctx.ExtraData, "overexert_units")
	if len(overexertUnits) == 0 {
		return nil
	}
	markedUnits := 0
	for _, unit := range overexertUnits {
		if unit == nil || unit.Card == nil || unit.OwnerID != castPlayer || !unit.Card.IsCompanion() {
			continue
		}
		if unit.Statuses == nil {
			unit.Statuses = make(map[string]int)
		}
		unit.Statuses[erebosSoulChainMarkedUnitStatus] = 1
		markedUnits++
	}
	if markedUnits == 0 {
		return nil
	}
	markedSpells := 0
	for _, spell := range append([]*CardInstance{ctx.Target}, spellInstancesFromData(ctx.ExtraData, "boost_skills")...) {
		if spell == nil || spell.Card == nil || spell.OwnerID != castPlayer || !isSpellLikeCard(spell.Card) {
			continue
		}
		if spell.Statuses == nil {
			spell.Statuses = make(map[string]int)
		}
		spell.Statuses[erebosSoulChainMarkedSpellStatus] = 1
		markedSpells++
	}
	if markedSpells > 0 {
		ctx.Source.UltimateUsed = true
		ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
			"source":        cardToInfo(ctx.Source),
			"effect":        "erebos_soul_chain_mark",
			"cast_player":   castPlayer,
			"marked_units":  markedUnits,
			"marked_spells": markedSpells,
		}})
	}
	return nil
}

func (Card2611101ErebosSoulChain) OnConsume(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Target == nil {
		return nil
	}
	ctx.Engine.weakenErebosSoulChainMarkedSpellsForUnit(ctx.Target)
	return nil
}

func (e *Engine) triggerErebosSoulChainMarkedOverexert(playerID int, units []*CardInstance) {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) || len(units) == 0 {
		return
	}
	opponentID := 1 - playerID
	if opponentID < 0 || opponentID >= len(e.State.Players) || !e.playerHasActiveCard(e.State.Players[opponentID], "2611101") {
		return
	}
	for _, unit := range units {
		e.weakenErebosSoulChainMarkedSpellsForUnit(unit)
	}
}

func (e *Engine) weakenErebosSoulChainMarkedSpellsForUnit(unit *CardInstance) int {
	if e == nil || unit == nil || unit.Statuses[erebosSoulChainMarkedUnitStatus] <= 0 || unit.OwnerID < 0 || unit.OwnerID >= len(e.State.Players) {
		return 0
	}
	weakened := 0
	for _, card := range e.getAllFieldCards(e.State.Players[unit.OwnerID]) {
		if card == nil || card.Statuses[erebosSoulChainMarkedSpellStatus] <= 0 || !canInstanceBeWeakened(card) {
			continue
		}
		if e.addStatus(card, StatusWeaken, 1) {
			weakened++
		}
	}
	if weakened > 0 {
		e.emit(GameEvent{Type: "effect_trigger", Player: unit.OwnerID, Data: map[string]any{
			"effect": "erebos_soul_chain_weaken",
			"unit":   cardToInfo(unit),
			"count":  weakened,
		}})
	}
	return weakened
}

var _ OnSpellCastBehavior = Card2611101ErebosSoulChain{}
var _ OnConsumeBehavior = Card2611101ErebosSoulChain{}

func (e *Engine) triggerAshKeltAfterOpponentShieldBreak(playerID int, data map[string]any) {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) || data == nil {
		return
	}
	attacker, _ := data["attacker"].(int)
	if attacker == playerID {
		return
	}
	ps := e.State.Players[playerID]
	for _, source := range e.getAllFieldCards(ps) {
		if source == nil || source.Card == nil || source.Card.Number != "1511101" || e.hasEffectiveStatus(source, StatusPetrify) {
			continue
		}
		e.drawCards(playerID, 2)
		ps.GainElements(map[string]int{model.ElementLight: 2})
		e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
			"source": cardToInfo(source),
			"effect": "ash_kelt_shield_break",
			"amount": 2,
		}})
	}
}

type Card2311102LampusSword struct{ AlwaysActive }

func (Card2311102LampusSword) ID() string   { return "2311102" }
func (Card2311102LampusSword) Name() string { return "兰普斯之剑" }
func (Card2311102LampusSword) PerTurnLabel(*CardInstance) string {
	return "主动"
}
func (Card2311102LampusSword) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	if !ctx.Engine.equipmentInOwnerSlot(ctx.PlayerID, ctx.Source) {
		return fmt.Errorf("兰普斯之剑需要在装备区")
	}
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.Category == model.ElementAir
	})
	sword := ctx.Source
	ctx.Engine.SetPendingAction(ctx.PlayerID, "lampus_sword_discard_air",
		"兰普斯之剑:弃置任意数量大气手牌", candidates, 0, len(candidates),
		func(selected []string) {
			if slot := equipmentSlotOf(ctx.Engine.State.Players[ctx.PlayerID], sword); slot >= 0 {
				ctx.Engine.moveEquipmentToGraveyard(ctx.PlayerID, slot, sword)
			}
			discarded := ctx.Engine.discardSelectedHandCardsMatching(ctx.PlayerID, selected, len(candidates), func(card *CardInstance) bool {
				return card != nil && card.Card != nil && card.Card.Category == model.ElementAir
			})
			damage := len(discarded)
			if damage > 0 {
				ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
					Type:             TempModLampusSwordDelayedDamage,
					SourceCardNumber: "2311102",
					SourceName:       "兰普斯之剑",
					Amount:           damage,
					RemainingUses:    1,
				})
			}
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source":    cardToInfo(sword),
				"effect":    "lampus_sword_prepare_damage",
				"discarded": damage,
			}})
		})
	return nil
}

func (e *Engine) promptLampusSwordDelayedDamage(playerID int, modifier TemporaryModifier) {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) || modifier.Amount <= 0 {
		return
	}
	candidates := e.enemyCompanionsInSpellRange(playerID)
	if len(candidates) == 0 {
		e.removeTemporaryModifier(playerID, modifier.ID)
		return
	}
	maxSelect := min(modifier.Amount, len(candidates))
	e.SetPendingAction(playerID, "lampus_sword_distribute_damage",
		"兰普斯之剑:分配延迟伤害", candidates, 1, maxSelect,
		func(selected []string) {
			allocations := map[string]int{}
			order := make([]string, 0, len(selected))
			for _, id := range selected {
				if allocations[id] == 0 {
					order = append(order, id)
				}
				allocations[id]++
			}
			for remaining := modifier.Amount - len(selected); remaining > 0 && len(order) > 0; remaining-- {
				allocations[order[0]]++
			}
			for _, id := range order {
				target := findEnemyCardCandidate(e, playerID, id, candidates)
				if target == nil || !e.unitStillOnField(target) {
					continue
				}
				e.dealDamageWithExtra(target, allocations[id], target.OwnerID, map[string]any{
					"damage_source": "effect",
					"attacker":      playerID,
				})
			}
			e.removeTemporaryModifier(playerID, modifier.ID)
		})
}

func (e *Engine) enemyCompanionsInSpellRange(playerID int) []map[string]any {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	opponent := e.State.Players[1-playerID]
	candidates := make([]map[string]any, 0)
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := opponent.Units[col][row]
			if unit == nil || unit.Card == nil || !unit.Card.IsCompanion() {
				continue
			}
			if !e.IsInSpellRange(playerID, col, row, false) {
				continue
			}
			candidates = append(candidates, candidateInfo(unit, "unit", "enemy"))
		}
	}
	return candidates
}

var _ PerTurnAbility = Card2311102LampusSword{}

const bloodSandArrayMarkerStatus = "蔽天阵 血沙标记物"

type Card3411102BloodSandArray struct{ AlwaysActive }

func (Card3411102BloodSandArray) ID() string   { return "3411102" }
func (Card3411102BloodSandArray) Name() string { return "蔽天阵 血沙" }
func (Card3411102BloodSandArray) PerTurnLabel(*CardInstance) string {
	return "主动"
}
func (Card3411102BloodSandArray) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	source := ctx.Source
	firstCandidates := ctx.Engine.bloodSandPaymentCandidates(ctx.PlayerID)
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "blood_sand_array_pay",
		"蔽天阵 血沙:选择最多3点己方单位负载或生命支付", firstCandidates, 0, min(3, len(firstCandidates)), nil, false,
		func(selected []string, data map[string]any) error {
			if !ctx.Engine.cardStillOnField(source) {
				return nil
			}
			firstPaid, err := ctx.Engine.applyBloodSandPayments(ctx.PlayerID, selected, data, source)
			if err != nil {
				return err
			}
			opponentID := 1 - ctx.PlayerID
			secondCandidates := ctx.Engine.bloodSandPaymentCandidates(opponentID)
			ctx.Engine.SetPendingActionWithError(opponentID, "blood_sand_array_pay_opponent",
				"蔽天阵 血沙:选择最多3点己方单位负载或生命支付", secondCandidates, 0, min(3, len(secondCandidates)), nil, false,
				func(opponentSelected []string, opponentData map[string]any) error {
					if !ctx.Engine.cardStillOnField(source) {
						return nil
					}
					secondPaid, err := ctx.Engine.applyBloodSandPayments(opponentID, opponentSelected, opponentData, source)
					if err != nil {
						return err
					}
					diff := bloodSandAbsDiff(firstPaid - secondPaid)
					if diff > 0 {
						source.Statuses[bloodSandArrayMarkerStatus] += diff
					}
					ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
						"source":      cardToInfo(source),
						"effect":      "blood_sand_array_markers",
						"paid_owner":  firstPaid,
						"paid_enemy":  secondPaid,
						"markers_add": diff,
						"markers":     source.Statuses[bloodSandArrayMarkerStatus],
					}})
					return nil
				})
			return nil
		})
	return nil
}

func (Card3411102BloodSandArray) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Source == nil {
		return
	}
	markers := ctx.Source.Statuses[bloodSandArrayMarkerStatus]
	if markers <= 0 {
		return
	}
	stats.PowerBonus += markers * 3
	stats.DamageBonus += markers
}

func (e *Engine) bloodSandPaymentCandidates(playerID int) []map[string]any {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	candidates := make([]map[string]any, 0)
	for _, unit := range e.getAllFieldCards(e.State.Players[playerID]) {
		if unit == nil || unit.Card == nil || unit.Position == nil {
			continue
		}
		if e.bloodSandPayablePoints(unit) <= 0 {
			continue
		}
		candidates = append(candidates, candidateInfo(unit, "unit", "own"))
	}
	return candidates
}

func (e *Engine) applyBloodSandPayments(playerID int, selected []string, data map[string]any, source *CardInstance) (int, error) {
	requests := bloodSandPaymentRequests(selected, data)
	if len(requests) == 0 {
		return 0, nil
	}
	total := 0
	for _, request := range requests {
		if total >= 3 {
			break
		}
		unit := e.findOwnUnitByInstanceID(playerID, request.instanceID)
		if unit == nil {
			return total, fmt.Errorf("invalid blood sand payment target")
		}
		amount := min(request.amount, 3-total)
		for i := 0; i < amount; i++ {
			if request.mode == "life" {
				if !bloodSandPayLife(unit) {
					return total, fmt.Errorf("cannot pay life from selected unit")
				}
				total++
				continue
			}
			if e.payOneBloodSandLoad(playerID, unit, source) || bloodSandPayLife(unit) {
				total++
				continue
			}
			return total, fmt.Errorf("selected unit cannot pay load or life")
		}
	}
	return total, nil
}

type bloodSandPaymentRequest struct {
	instanceID string
	amount     int
	mode       string
}

func bloodSandPaymentRequests(selected []string, data map[string]any) []bloodSandPaymentRequest {
	if data != nil {
		if raw, ok := data["payments"].([]any); ok {
			requests := make([]bloodSandPaymentRequest, 0, len(raw))
			for _, entry := range raw {
				m, ok := entry.(map[string]any)
				if !ok {
					continue
				}
				id, _ := m["instance_id"].(string)
				if id == "" {
					continue
				}
				amount := intFromData(m, "amount", 1)
				if amount <= 0 {
					continue
				}
				mode, _ := m["mode"].(string)
				requests = append(requests, bloodSandPaymentRequest{instanceID: id, amount: amount, mode: mode})
			}
			return requests
		}
	}
	requests := make([]bloodSandPaymentRequest, 0, len(selected))
	for _, id := range selected {
		if id != "" {
			requests = append(requests, bloodSandPaymentRequest{instanceID: id, amount: 1})
		}
	}
	return requests
}

func (e *Engine) findOwnUnitByInstanceID(playerID int, instanceID string) *CardInstance {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) || instanceID == "" {
		return nil
	}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := e.State.Players[playerID].Units[col][row]
			if unit != nil && unit.InstanceID == instanceID {
				return unit
			}
		}
	}
	return nil
}

func (e *Engine) bloodSandPayablePoints(unit *CardInstance) int {
	if unit == nil {
		return 0
	}
	return min(3, totalLoad(unit)+max(unit.CurrentLife-1, 0))
}

func (e *Engine) payOneBloodSandLoad(playerID int, unit *CardInstance, source *CardInstance) bool {
	if unit == nil {
		return false
	}
	for _, elem := range model.AllElements {
		if e.effectiveElementsGain(unit)[elem] <= 0 {
			continue
		}
		return e.reduceCardElementLoadWithTriggers(playerID, unit, elem, 1, source) == 1
	}
	return false
}

func bloodSandPayLife(unit *CardInstance) bool {
	if unit == nil || unit.CurrentLife <= 1 {
		return false
	}
	unit.CurrentLife--
	return true
}

func bloodSandAbsDiff(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

var _ PerTurnAbility = Card3411102BloodSandArray{}
var _ SkillContributionModifier = Card3411102BloodSandArray{}

type Card2321109MistMask struct{ AlwaysActive }

func (Card2321109MistMask) ID() string   { return "2321109" }
func (Card2321109MistMask) Name() string { return "幻雾面罩" }
func (Card2321109MistMask) OnSpellHitBeforeDamage(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil || ctx.Source.UltimateUsed {
		return nil
	}
	attacker, ok := ctx.ExtraData["attacker"].(int)
	damagePtr, hasDamage := ctx.ExtraData["damage_ptr"].(*int)
	if !ok || attacker == ctx.PlayerID || !hasDamage || damagePtr == nil || *damagePtr <= 0 {
		return nil
	}
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, nil)
	if len(candidates) == 0 {
		return nil
	}
	mask := ctx.Source
	ctx.Engine.SetPendingAction(ctx.PlayerID, "mist_mask_discard_reduce_spell_attack",
		"幻雾面罩:丢弃最多3张手牌降低该法术伤害", candidates, 0, min(3, len(candidates)),
		func(selected []string) {
			if len(selected) == 0 || mask.UltimateUsed || !ctx.Engine.cardStillOnField(mask) {
				return
			}
			discarded := ctx.Engine.discardSelectedHandCardsMatching(ctx.PlayerID, selected, 3, nil)
			reduction := min(len(discarded), *damagePtr)
			if reduction <= 0 {
				return
			}
			*damagePtr -= reduction
			ctx.ExtraData["damage"] = *damagePtr
			mask.UltimateUsed = true
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source":    cardToInfo(mask),
				"effect":    "mist_mask_reduce_spell_attack",
				"discarded": len(discarded),
				"reduction": reduction,
				"damage":    *damagePtr,
			}})
		})
	return nil
}

type Card3011101AbsolutePurityArcaneOneness struct{ AlwaysActive }

func (Card3011101AbsolutePurityArcaneOneness) ID() string   { return "3011101" }
func (Card3011101AbsolutePurityArcaneOneness) Name() string { return "绝对纯净 奥能一心" }
func (Card3011101AbsolutePurityArcaneOneness) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.Source.Card.Number != "3011101" {
		return
	}
	stats.PowerBonus += countTopConsecutiveArcaneCards(ctx.Engine.State.Players[ctx.PlayerID])
}
func (Card3011101AbsolutePurityArcaneOneness) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.Source.Card.Number != "3011101" || !isFriendlySpellCast(ctx) {
		return nil
	}
	count := countTopConsecutiveArcaneCards(ctx.Engine.State.Players[ctx.PlayerID])
	if count > 0 || len(ctx.Engine.State.Players[ctx.PlayerID].Deck) > 0 {
		ctx.Engine.shuffleDeck(ctx.PlayerID)
	}
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source":       cardToInfo(ctx.Source),
		"effect":       "arcane_oneness_reveal",
		"arcane_count": count,
	}})
	return nil
}

func countTopConsecutiveArcaneCards(ps *PlayerState) int {
	if ps == nil {
		return 0
	}
	count := 0
	for _, card := range ps.Deck {
		if card == nil || card.Card == nil || card.Card.Category != model.ElementArcane {
			break
		}
		count++
	}
	return count
}

type Card3311102StarfallSilverleaf struct{ AlwaysActive }

func (Card3311102StarfallSilverleaf) ID() string   { return "3311102" }
func (Card3311102StarfallSilverleaf) Name() string { return "星落之银叶" }
func (Card3311102StarfallSilverleaf) OnDiscard(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.ExtraData == nil {
		return nil
	}
	discardedPlayer, _ := ctx.ExtraData["discarded_player"].(int)
	if discardedPlayer != ctx.PlayerID || ctx.Target == ctx.Source {
		return nil
	}
	discarded := ctx.Target
	source := ctx.Source
	ctx.Engine.SetPendingAction(ctx.PlayerID, "starfall_silverleaf_store_discard",
		"星落之银叶:将弃置的卡牌放在此卡下方", []map[string]any{candidateInfo(discarded, "graveyard", "own")}, 1, 1,
		func(selected []string) {
			if !ctx.Engine.cardStillOnField(source) {
				return
			}
			ctx.Engine.placeCardUnder(source, discarded)
		})
	return nil
}
func (Card3311102StarfallSilverleaf) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.Source.Card.Number != "3311102" || len(ctx.Source.UnderCards) == 0 {
		return nil
	}
	candidates := make([]map[string]any, 0, len(ctx.Source.UnderCards))
	for _, card := range ctx.Source.UnderCards {
		if card != nil {
			candidates = append(candidates, candidateInfo(card, "under", "own"))
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	source := ctx.Source
	ctx.Engine.SetPendingAction(ctx.PlayerID, "starfall_silverleaf_recycle_under_card",
		"星落之银叶:选择下方1张牌洗回卡组并抽1张", candidates, 1, 1,
		func(selected []string) {
			if !ctx.Engine.cardStillOnField(source) {
				return
			}
			card := ctx.Engine.detachCardFromKnownZones(ctx.Engine.State.Players[ctx.PlayerID], firstSelected(selected))
			if card == nil {
				return
			}
			resetCardForPublicSpecialZone(card)
			ctx.Engine.State.Players[ctx.PlayerID].Deck = append(ctx.Engine.State.Players[ctx.PlayerID].Deck, card)
			ctx.Engine.shuffleDeck(ctx.PlayerID)
			ctx.Engine.drawCards(ctx.PlayerID, 1)
		})
	return nil
}

type Card3511101DivineRadianceSkyward struct{ AlwaysActive }

func (Card3511101DivineRadianceSkyward) ID() string   { return "3511101" }
func (Card3511101DivineRadianceSkyward) Name() string { return "神辉驭空" }
func (Card3511101DivineRadianceSkyward) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.Source.Card.Number != "3511101" {
		return
	}
	stats.PowerBonus += len(ctx.Engine.State.Players[ctx.OpponentID].Hand)
}
func (Card3511101DivineRadianceSkyward) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.Source.Card.Number != "3511101" {
		return nil
	}
	candidates := []map[string]any{
		{"instance_id": "player:own", "name": fmt.Sprintf("玩家%d", ctx.PlayerID+1), "zone": "player", "side": "own", "player_id": ctx.PlayerID},
		{"instance_id": "player:opponent", "name": fmt.Sprintf("玩家%d", ctx.OpponentID+1), "zone": "player", "side": "enemy", "player_id": ctx.OpponentID},
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "divine_radiance_reset_hand",
		"神辉驭空:选择1名玩家弃掉全部手牌并抽至手牌上限", candidates, 0, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			targetPlayer := ctx.PlayerID
			if firstSelected(selected) == "player:opponent" {
				targetPlayer = ctx.OpponentID
			}
			count := ctx.Engine.discardAllHandCards(targetPlayer)
			limit := ctx.Engine.handLimitForPlayer(ctx.Engine.State.Players[targetPlayer])
			if limit > 0 {
				ctx.Engine.drawCards(targetPlayer, limit)
			}
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source":    cardToInfo(ctx.Source),
				"effect":    "divine_radiance_reset_hand",
				"target":    targetPlayer,
				"discarded": count,
				"drawn":     limit,
			}})
		})
	return nil
}

type Card2321112RendingImpactScroll struct{ AlwaysActive }

func (Card2321112RendingImpactScroll) ID() string   { return "2321112" }
func (Card2321112RendingImpactScroll) Name() string { return "撕裂冲击卷轴" }
func (Card2321112RendingImpactScroll) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.Source.Card.Number != "2321112" || ctx.ExtraData == nil {
		return nil
	}
	candidates := make([]map[string]any, 0)
	for _, unit := range spellHitAffectedUnitsFromData(ctx) {
		if unit != nil && unit.OwnerID == ctx.OpponentID && ctx.Engine.unitStillOnField(unit) {
			candidates = append(candidates, candidateInfo(unit, "unit", "enemy"))
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "rending_impact_distribute_damage",
		"撕裂冲击卷轴:选择目标范围内单位分配共计3点伤害", candidates, 1, min(3, len(candidates)),
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			allocations := map[string]int{}
			order := make([]string, 0, len(selected))
			for _, id := range selected {
				if allocations[id] == 0 {
					order = append(order, id)
				}
				allocations[id]++
			}
			for remaining := 3 - len(selected); remaining > 0 && len(order) > 0; remaining-- {
				allocations[order[0]]++
			}
			for _, id := range order {
				target := findEnemyCardCandidate(ctx.Engine, ctx.PlayerID, id, candidates)
				if target == nil || !ctx.Engine.unitStillOnField(target) {
					continue
				}
				ctx.Engine.dealDamageWithExtra(target, allocations[id], target.OwnerID, map[string]any{
					"damage_source": "effect",
					"attacker":      ctx.PlayerID,
				})
			}
		})
	return nil
}

type Card3601101BloodFeast struct{ AlwaysActive }

func (Card3601101BloodFeast) ID() string   { return "3601101" }
func (Card3601101BloodFeast) Name() string { return "鲜血盛宴" }
func (Card3601101BloodFeast) AllowsFriendlySpellTarget() bool {
	return true
}
func (Card3601101BloodFeast) ValidateSpellTarget(ctx *EffectContext, target SpellTarget, targetUnit *CardInstance) error {
	if ctx == nil || target.Type != "unit" || targetUnit == nil || targetUnit.OwnerID != ctx.PlayerID {
		return fmt.Errorf("鲜血盛宴只能攻击友方单位")
	}
	return nil
}
func (Card3601101BloodFeast) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.Source.Card.Number != "3601101" || !isOwnSpellHit(ctx) {
		return nil
	}
	choices := []map[string]any{
		{"instance_id": "gain_shadow", "name": "获得2暗", "zone": "choice", "side": "own"},
		{"instance_id": "heal_hero", "name": "人物回复1血", "zone": "choice", "side": "own"},
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "blood_feast_reward",
		"鲜血盛宴:选择命中奖励", choices, 1, 1,
		func(selected []string) {
			switch firstSelected(selected) {
			case "gain_shadow":
				ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementShadow: 2})
			case "heal_hero":
				hero := ctx.Engine.playerHeroCard(ctx.PlayerID)
				if hero != nil && hero.CurrentLife < maxLife(hero) {
					hero.CurrentLife++
				}
			}
		})
	return nil
}
func (Card3601101BloodFeast) PerTurnLabel(*CardInstance) string { return "绑定" }
func (Card3601101BloodFeast) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ps.Elements[model.ElementShadow] < 1 {
		return fmt.Errorf("鲜血盛宴需要支付1暗绑定到人物")
	}
	hero := ctx.Engine.playerHeroCard(ctx.PlayerID)
	if hero == nil {
		return nil
	}
	for i, skill := range ps.Skills {
		if skill == ctx.Source {
			ps.Skills[i] = nil
			ps.Elements[model.ElementShadow]--
			ctx.Source.SlotIndex = -1
			markTransferredBoundSkill(ctx.Source)
			hero.BoundSkills = append(hero.BoundSkills, ctx.Source)
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source": cardToInfo(ctx.Source),
				"target": cardToInfo(hero),
				"effect": "bind_existing_skill",
			}})
			return nil
		}
	}
	return nil
}

type Card2621106PainScreamScroll struct{ AlwaysActive }

func (Card2621106PainScreamScroll) ID() string   { return "2621106" }
func (Card2621106PainScreamScroll) Name() string { return "苦痛尖啸卷轴" }
func (Card2621106PainScreamScroll) OnUseItem(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModPainScreamWeakenOnDamage,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		ExpiresTurn:      ctx.Engine.State.TurnNumber,
	})
	return nil
}

func (e *Engine) promptPainScreamWeakenAfterFriendlyDamage(playerID int, target *CardInstance, damage int) {
	if e == nil || target == nil || damage <= 0 || playerID < 0 || playerID >= len(e.State.Players) || e.State.PendingAction != nil {
		return
	}
	ps := e.State.Players[playerID]
	if ps == nil || !playerHasPainScreamModifier(ps) {
		return
	}
	candidates := enemySpellCandidatesWithoutWeaken(e, playerID)
	if len(candidates) == 0 {
		return
	}
	e.SetPendingAction(playerID, "pain_scream_weaken_enemy_spells",
		"苦痛尖啸卷轴:选择没有虚弱的敌方法术获得虚弱2", candidates, 1, min(damage, len(candidates)),
		func(selected []string) {
			weakened := 0
			for _, id := range selected {
				if weakened >= damage {
					break
				}
				skill := findEnemySkillIncludingBound(e, playerID, id)
				if skill == nil || !canInstanceBeWeakened(skill) || skill.Statuses[StatusWeaken] > 0 {
					continue
				}
				e.addStatus(skill, StatusWeaken, 2)
				weakened++
			}
		})
}

func playerHasPainScreamModifier(ps *PlayerState) bool {
	if ps == nil {
		return false
	}
	for _, modifier := range ps.TempModifiers {
		if modifier.Type == TempModPainScreamWeakenOnDamage {
			return true
		}
	}
	return false
}

func enemySpellCandidatesWithoutWeaken(e *Engine, playerID int) []map[string]any {
	candidates := make([]map[string]any, 0)
	for _, skill := range enemySpellInstancesIncludingBound(e, playerID) {
		if skill != nil && skill.Card != nil && canInstanceBeWeakened(skill) && skill.Statuses[StatusWeaken] <= 0 {
			candidates = append(candidates, candidateInfo(skill, "skill", "enemy"))
		}
	}
	return candidates
}

func findEnemySkillIncludingBound(e *Engine, playerID int, instanceID string) *CardInstance {
	for _, skill := range enemySpellInstancesIncludingBound(e, playerID) {
		if skill != nil && skill.InstanceID == instanceID {
			return skill
		}
	}
	return nil
}

const protectorSivalPreventionUntilStatus = "西瓦尔防止伤害至回合"

type Card4511101ProtectorSival struct{ AlwaysActive }

func (Card4511101ProtectorSival) ID() string   { return "4511101" }
func (Card4511101ProtectorSival) Name() string { return "庇护者 西瓦尔" }
func (Card4511101ProtectorSival) OnDamaged(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.ExtraData == nil || ctx.Source.UltimateUsed {
		return nil
	}
	damagedPlayer, _ := ctx.ExtraData["damaged_player"].(int)
	if damagedPlayer != ctx.PlayerID || friendlyDamageTakenThisTurn(ctx.Engine, ctx.PlayerID) < 3 {
		return nil
	}
	source := ctx.Source
	ctx.Engine.SetPendingAction(ctx.PlayerID, "protector_sival_prevent_all_damage",
		"庇护者 西瓦尔:是否发动绝技防止所有友方单位伤害", []map[string]any{candidateInfo(source, "hero", "own")}, 0, 1,
		func(selected []string) {
			if len(selected) == 0 || source.UltimateUsed || !ctx.Engine.cardStillOnField(source) {
				return
			}
			source.UltimateUsed = true
			source.Statuses[protectorSivalPreventionUntilStatus] = ctx.Engine.State.TurnNumber + 1
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source": cardToInfo(source),
				"effect": "prevent_friendly_damage",
			}})
		})
	return nil
}
func (Card4511101ProtectorSival) PreventsFieldDamage(ctx *EffectContext) bool {
	if ctx == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.OwnerID != ctx.PlayerID {
		return false
	}
	return ctx.Source.Statuses[protectorSivalPreventionUntilStatus] >= ctx.Engine.State.TurnNumber
}
func (Card4511101ProtectorSival) PreventsDamage(ctx *EffectContext) bool {
	if ctx == nil || ctx.Source == nil || ctx.Engine == nil {
		return false
	}
	return ctx.Source.Statuses[protectorSivalPreventionUntilStatus] >= ctx.Engine.State.TurnNumber
}

func friendlyDamageTakenThisTurn(e *Engine, playerID int) int {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return 0
	}
	total := 0
	for _, card := range e.getAllFieldCards(e.State.Players[playerID]) {
		if card != nil && (card.Card.IsHero() || card.Card.IsCompanion()) {
			total += card.DamageTakenThisTurn
		}
	}
	return total
}

type Card3521109Regroup struct{ AlwaysActive }

func (Card3521109Regroup) ID() string   { return "3521109" }
func (Card3521109Regroup) Name() string { return "重整旗鼓" }
func (Card3521109Regroup) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil {
		return nil
	}
	attacker, ok := ctx.ExtraData["attacker"].(int)
	if !ok || attacker == ctx.PlayerID || !canTriggeredRegroupBeUsed(ctx.Source) {
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion()
	})
	if len(candidates) == 0 {
		return nil
	}
	source := ctx.Source
	ctx.Engine.SetPendingAction(ctx.PlayerID, "regroup_buff_companion",
		"重整旗鼓:选择1个友方伙伴获得+1血和负载+1光", candidates, 1, 1,
		func(selected []string) {
			if !canTriggeredRegroupBeUsed(source) || ctx.Engine.findSkill(ctx.Engine.State.Players[ctx.PlayerID], source.InstanceID) != source {
				return
			}
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target == nil || zone != "unit" || target.Card == nil || !target.Card.IsCompanion() {
				return
			}
			source.IsHorizontal = true
			ctx.Engine.ApplyKeywordOnSkillUse(source)
			ctx.Engine.applySkillUseCooldownModifiers(ctx.Engine.State.Players[ctx.PlayerID], source)
			target.Statuses["max_life_bonus"]++
			target.CurrentLife++
			ctx.Engine.addElementsGainBonus(target, ctx.PlayerID, model.ElementLight, 1, source)
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source": cardToInfo(source),
				"target": cardToInfo(target),
				"effect": "regroup_buff_companion",
			}})
		})
	return nil
}

func canTriggeredRegroupBeUsed(skill *CardInstance) bool {
	return skill != nil && skill.Card != nil && skill.Card.Number == "3521109" && !skill.IsHorizontal && skill.Statuses[StatusCooldown] <= 0
}

type Card1021113MagicMoth struct{ AlwaysActive }

func (Card1021113MagicMoth) ID() string   { return "1021113" }
func (Card1021113MagicMoth) Name() string { return "魔法飞蛾" }

func (e *Engine) triggerMagicMothAfterFocusSpellCast(playerID int, skill *CardInstance) {
	if e == nil || skill == nil || skill.Card == nil || !hasCardTag(skill.Card, "聚能") || playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	ps := e.State.Players[playerID]
	for _, card := range ps.Deck {
		if card == nil || card.Card == nil || card.Card.Number != "1021113" {
			continue
		}
		moth := card
		e.SetPendingAction(playerID, "magic_moth_draw",
			"魔法飞蛾:是否从卡组抽取本卡", []map[string]any{candidateInfo(moth, "deck", "own")}, 0, 1,
			func(selected []string) {
				if len(selected) == 0 {
					return
				}
				for i, current := range ps.Deck {
					if current != moth {
						continue
					}
					ps.Deck = append(ps.Deck[:i], ps.Deck[i+1:]...)
					e.appendCardsToHand(playerID, []*CardInstance{moth})
					e.notifyCardDrawn(playerID, moth)
					e.shuffleDeck(playerID)
					e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
						"effect": "magic_moth_draw",
						"source": cardToInfo(skill),
						"card":   cardToInfo(moth),
					}})
					e.enforceImmediateHandLimitAfterHandGain(playerID)
					return
				}
			})
		return
	}
}

const sinTargetTagStatusPrefix = "罪责目标种类:"

type Card3521104Sin struct{ AlwaysActive }

func (Card3521104Sin) ID() string   { return "3521104" }
func (Card3521104Sin) Name() string { return "罪责" }
func (Card3521104Sin) OnEnter(ctx *EffectContext) error {
	choices := companionTagChoices()
	ctx.Engine.SetPendingAction(ctx.PlayerID, "sin_choose_companion_kind",
		"罪责:选择1个伙伴种类", choices, 1, 1,
		func(selected []string) {
			tag := firstSelected(selected)
			if !validCompanionTagChoice(tag) {
				return
			}
			clearStatusPrefix(ctx.Source, sinTargetTagStatusPrefix)
			ctx.Source.Statuses[sinTargetTagStatusPrefix+tag] = 1
		})
	return nil
}

func (Card3521104Sin) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Source == nil || ctx.ExtraData == nil || ctx.ExtraData["purpose"] != string(skillPurposeAttack) {
		return
	}
	target, _ := ctx.ExtraData["spell_target_unit"].(*CardInstance)
	if !sinMatchesTargetTag(ctx.Source, target) {
		return
	}
	stats.PowerBonus += 2
	stats.Pierce = true
}

func companionTagChoices() []map[string]any {
	tags := []string{"人类", "巫师", "野兽", "精灵", "恶魔", "造物", "植物", "灵体", "龙"}
	choices := make([]map[string]any, 0, len(tags))
	for _, tag := range tags {
		choices = append(choices, map[string]any{
			"instance_id": tag,
			"name":        tag,
			"zone":        "choice",
			"side":        "own",
		})
	}
	return choices
}

func validCompanionTagChoice(tag string) bool {
	for _, choice := range companionTagChoices() {
		if choice["instance_id"] == tag {
			return true
		}
	}
	return false
}

func sinMatchesTargetTag(skill *CardInstance, target *CardInstance) bool {
	if skill == nil || target == nil || target.Card == nil || !target.Card.IsCompanion() {
		return false
	}
	for status, amount := range skill.Statuses {
		if amount <= 0 || !strings.HasPrefix(status, sinTargetTagStatusPrefix) {
			continue
		}
		tag := strings.TrimPrefix(status, sinTargetTagStatusPrefix)
		return hasCardTag(target.Card, tag)
	}
	return false
}

type Card1121109DivineFireRider struct{ AlwaysActive }

func (Card1121109DivineFireRider) ID() string   { return "1121109" }
func (Card1121109DivineFireRider) Name() string { return "神火兽骑手" }
func (Card1121109DivineFireRider) OnUltimate(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != ctx.Source && isFireCompanion(card) && ctx.Engine.canConsumeCard(card)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "divine_fire_rider_consume_companion",
		"神火兽骑手:消耗1个其他友方火焰伙伴", candidates, 1, 1,
		func(selected []string) {
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target == nil || target == ctx.Source || zone != "unit" || !isFireCompanion(target) || !ctx.Engine.canConsumeCard(target) {
				return
			}
			bonus := totalElementCost(target.Card.ElementsCost)
			if bonus <= 0 {
				return
			}
			ctx.Engine.consumeCardForEffectWithTriggers(ctx.PlayerID, target, ctx.Engine.effectiveElementsGain(target), "1121109")
			ctx.Engine.addNextElementSpellPowerBonus(ctx.PlayerID, model.ElementFire, bonus)
		})
	return nil
}

func isFireBeastOrMonsterCompanion(card *CardInstance) bool {
	if card == nil || card.Card == nil || !card.Card.IsCompanion() || card.Card.Category != model.ElementFire {
		return false
	}
	return strings.Contains(card.Card.Tag, "野兽") || strings.Contains(card.Card.Tag, "异兽")
}

type Card1121101VolcanoSalamander struct{ AlwaysActive }

func (Card1121101VolcanoSalamander) ID() string      { return "1121101" }
func (Card1121101VolcanoSalamander) Name() string    { return "火山蝾螈" }
func (Card1121101VolcanoSalamander) MasteryMax() int { return 2 }
func (Card1121101VolcanoSalamander) OnMastery(ctx *EffectContext, level int) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || level != 2 || !ctx.Engine.cardStillOnField(ctx.Source) {
		return nil
	}
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, isFireCompanionWithEntryCostLessThanEight)
	if len(candidates) == 0 || len(friendlyPositionsAfterRemovingSource(ctx.Engine, ctx.PlayerID, ctx.Source)) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "volcano_salamander_summon_card",
		"火山蝾螈:选择1个入场费用小于8的火焰伙伴免费召唤", candidates, 1, 1,
		func(selected []string) {
			cardID := firstSelected(selected)
			positions := friendlyPositionsAfterRemovingSource(ctx.Engine, ctx.PlayerID, ctx.Source)
			if cardID == "" || len(positions) == 0 {
				return
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "volcano_salamander_summon_position",
				"火山蝾螈:选择召唤位置", positions, 1, 1,
				func(posSelected []string) {
					pos, ok := positionFromSelectionID(firstSelected(posSelected))
					if !ok || !ctx.Engine.cardStillOnField(ctx.Source) {
						return
					}
					if card, _ := ctx.Engine.State.Players[ctx.PlayerID].FindHandCard(cardID); !isFireCompanionWithEntryCostLessThanEight(card) {
						return
					}
					ctx.Engine.destroyUnitWithCause(ctx.Source, ctx.PlayerID, DeathCauseSacrifice)
					summonCardFreeFromHandOrDeckAtPosition(ctx, cardID, pos)
				})
		})
	return nil
}

func isFireCompanionWithEntryCostLessThanEight(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.IsCompanion() &&
		card.Card.Category == model.ElementFire &&
		totalElementCost(card.Card.ElementsCost) < 8
}

func friendlyPositionsAfterRemovingSource(e *Engine, playerID int, source *CardInstance) []map[string]any {
	positions := e.friendlyEmptyUnitPositions(playerID)
	if source == nil || source.Position == nil {
		return positions
	}
	pos := *source.Position
	if !pos.Valid() {
		return positions
	}
	return append(positions, map[string]any{
		"instance_id": positionSelectionID(pos),
		"name":        fmt.Sprintf("位置 (%d,%d)", pos.Col, pos.Row),
		"zone":        "field_position",
		"side":        "own",
		"col":         pos.Col,
		"row":         pos.Row,
	})
}

type Card1121114LegionGeneral struct{ AlwaysActive }

func (Card1121114LegionGeneral) ID() string            { return "1121114" }
func (Card1121114LegionGeneral) Name() string          { return "军团将星" }
func (Card1121114LegionGeneral) IsPrayerAbility() bool { return true }

func (Card1121114LegionGeneral) OnPerTurn(ctx *EffectContext) error {
	choices := []map[string]any{
		{"instance_id": "power", "name": "火焰法术+2威", "zone": "choice", "side": "own"},
		{"instance_id": "attack", "name": "火焰法术+1攻", "zone": "choice", "side": "own"},
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "legion_general_prayer",
		"军团将星:选择你的火焰法术直到下个回合结束获得的效果", choices, 1, 1,
		func(selected []string) {
			switch firstSelected(selected) {
			case "power":
				ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
					Type:        TempModSkillPowerBonus,
					Element:     model.ElementFire,
					Amount:      2,
					ExpiresTurn: ctx.Engine.State.TurnNumber + 2,
				})
			case "attack":
				ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
					Type:        TempModCurrentTurnElementDamage,
					Element:     model.ElementFire,
					Amount:      1,
					ExpiresTurn: ctx.Engine.State.TurnNumber + 2,
				})
			}
		})
	return nil
}

type Card1121110LavafortArchivist struct{ AlwaysActive }

func (Card1121110LavafortArchivist) ID() string   { return "1121110" }
func (Card1121110LavafortArchivist) Name() string { return "熔岩堡档案员" }
func (Card1121110LavafortArchivist) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.Card == nil || ctx.Source.UltimateUsed {
		return nil
	}
	if !isFriendlySpellCast(ctx) || !hasCardTag(ctx.Target.Card, "创造") {
		return nil
	}
	drawn := ctx.Engine.flipDeckMatchesToHand(ctx.PlayerID, 1, 0, isRuneOrScroll)
	if len(drawn) > 0 {
		ctx.Source.UltimateUsed = true
	}
	return nil
}

type Card1121115LegionStaffOfficer struct{ AlwaysActive }

func (Card1121115LegionStaffOfficer) ID() string   { return "1121115" }
func (Card1121115LegionStaffOfficer) Name() string { return "军团参谋" }
func (Card1121115LegionStaffOfficer) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.Card == nil {
		return nil
	}
	if ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) || !isFriendlySpellCast(ctx) || !hasCardTag(ctx.Target.Card, "创造") {
		return nil
	}
	if !deckHasMatch(ctx.Engine.State.Players[ctx.PlayerID], isFireConsumableItem) {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "legion_staff_officer_flip_fire_consumable",
		"军团参谋:是否翻取1个火焰消耗品道具", []map[string]any{candidateInfo(ctx.Source, "unit", "own")}, 0, 1,
		func(selected []string) {
			if len(selected) == 0 || ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) {
				return
			}
			drawn := ctx.Engine.flipDeckMatchesToHand(ctx.PlayerID, 1, 0, isFireConsumableItem)
			if len(drawn) == 0 {
				return
			}
			ps := ctx.Engine.State.Players[ctx.PlayerID]
			if ps.DiscardAtTurnEnd == nil {
				ps.DiscardAtTurnEnd = make(map[string]bool)
			}
			ps.DiscardAtTurnEnd[drawn[0].InstanceID] = true
			ctx.Source.UsedThisTurn++
		})
	return nil
}

func isFireConsumableItem(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.Category == model.ElementFire && isConsumableCardInstance(card)
}

func deckHasMatch(ps *PlayerState, predicate func(*CardInstance) bool) bool {
	if ps == nil {
		return false
	}
	for _, card := range ps.Deck {
		if card != nil && (predicate == nil || predicate(card)) {
			return true
		}
	}
	return false
}

type Card1211103SeaHeroineCoralWendy struct{ AlwaysActive }

func (Card1211103SeaHeroineCoralWendy) ID() string   { return "1211103" }
func (Card1211103SeaHeroineCoralWendy) Name() string { return "海上巾帼 珊瑚 雯迪" }
func (Card1211103SeaHeroineCoralWendy) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.Card == nil {
		return nil
	}
	if ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) || !isFriendlySpellCast(ctx) || totalElementCost(ctx.Target.Card.ElementsExpense) >= 3 {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	cost := map[string]int{model.ElementWater: 2}
	if !ctx.Engine.canPayCost(ps, cost) {
		return nil
	}
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "coral_wendy_reset_spell",
		"海上巾帼 珊瑚 雯迪:是否支付2水重置刚使用的法术", []map[string]any{candidateInfo(ctx.Target, "skill", "own")}, 0, 1, cost, false,
		func(selected []string, data map[string]any) error {
			if len(selected) == 0 {
				return nil
			}
			if selected[0] != ctx.Target.InstanceID || ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) {
				return fmt.Errorf("invalid coral wendy reset")
			}
			target := ctx.Engine.findSkill(ctx.Engine.State.Players[ctx.PlayerID], selected[0])
			if target == nil || target.Card == nil || !target.IsHorizontal || totalElementCost(target.Card.ElementsExpense) >= 3 {
				return fmt.Errorf("invalid coral wendy reset")
			}
			if !ctx.Engine.payCostForAction(ps, cost, ActionMessage{Data: data}) {
				return fmt.Errorf("invalid coral wendy payment")
			}
			target.IsHorizontal = false
			ctx.Source.UsedThisTurn++
			return nil
		})
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

type Card1221108HeartLotusMirrorMage struct{ AlwaysActive }

func (Card1221108HeartLotusMirrorMage) ID() string   { return "1221108" }
func (Card1221108HeartLotusMirrorMage) Name() string { return "心莲镜魔师" }
func (Card1221108HeartLotusMirrorMage) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.Card == nil {
		return nil
	}
	if ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) || !isFriendlySpellCast(ctx) || !hasCardTag(ctx.Target.Card, "创造") {
		return nil
	}
	if !deckHasMatch(ctx.Engine.State.Players[ctx.PlayerID], isWaterItem) {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "heart_lotus_mirror_mage_flip_water_item",
		"心莲镜魔师:是否翻取1张水纹道具", []map[string]any{candidateInfo(ctx.Source, "unit", "own")}, 0, 1,
		func(selected []string) {
			if len(selected) == 0 || ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) {
				return
			}
			drawn := ctx.Engine.flipDeckMatchesToHandThen(ctx.PlayerID, 1, 0, isWaterItem, func(drawn []*CardInstance) {
				if len(drawn) == 0 {
					return
				}
				card := drawn[0]
				if isCounterTrapCard(card.Card.Number) && ctx.Engine.freeEquipmentSlots(ctx.PlayerID) > 0 {
					makeEntryCostZero(card)
					ctx.Engine.setHandCounterTrapFree(ctx.PlayerID, card)
				}
			})
			if len(drawn) == 0 {
				return
			}
			ctx.Source.UsedThisTurn++
		})
	return nil
}

func isWaterItem(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.IsItem() && card.Card.Category == model.ElementWater
}

type Card1411103KingOfBeasts struct{ AlwaysActive }

func (Card1411103KingOfBeasts) ID() string   { return "1411103" }
func (Card1411103KingOfBeasts) Name() string { return "百兽之王 莱恩克塞斯" }
func (Card1411103KingOfBeasts) OnEnter(ctx *EffectContext) error {
	drawn := ctx.Engine.flipDeckMatchesToHandThen(ctx.PlayerID, 1, 0, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion()
	}, func(drawn []*CardInstance) {
		if len(drawn) == 0 || drawn[0].Card.Category != model.ElementEarth {
			return
		}
		positions := ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)
		if len(positions) == 0 {
			return
		}
		cardID := drawn[0].InstanceID
		ctx.Engine.SetPendingAction(ctx.PlayerID, "king_of_beasts_summon_earth_companion",
			"百兽之王 莱恩克塞斯:选择位置免费召唤翻取的地脉伙伴", positions, 1, 1,
			func(selected []string) {
				pos, ok := positionFromSelectionID(firstSelected(selected))
				if !ok {
					return
				}
				card, _ := ctx.Engine.State.Players[ctx.PlayerID].FindHandCard(cardID)
				if card == nil || card.Card == nil || !card.Card.IsCompanion() || card.Card.Category != model.ElementEarth {
					return
				}
				summonCardFreeFromHandOrDeckAtPosition(ctx, cardID, pos)
			})
	})
	_ = drawn
	return nil
}

type Card1411101AgedFrankenBaililan struct{ AlwaysActive }

func (Card1411101AgedFrankenBaililan) ID() string            { return "1411101" }
func (Card1411101AgedFrankenBaililan) Name() string          { return "苍老者 弗兰肯 拜利兰" }
func (Card1411101AgedFrankenBaililan) IsPrayerAbility() bool { return true }
func (Card1411101AgedFrankenBaililan) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	load := dragonBloodTreantReducibleLoad(ctx.Source)
	candidates := make([]map[string]any, 0)
	for _, elem := range model.AllElements {
		if load[elem] <= 0 {
			continue
		}
		candidates = append(candidates, map[string]any{
			"instance_id": elem,
			"name":        elem,
			"zone":        "choice",
			"side":        "own",
		})
	}
	removeLoad := func(elem string) {
		if elem == "" || reducibleElementLoad(ctx.Source, elem) <= 0 {
			return
		}
		ctx.Engine.reduceCardElementLoadWithTriggers(ctx.PlayerID, ctx.Source, elem, 1, ctx.Source)
	}
	if len(candidates) == 0 {
		return nil
	}
	if len(candidates) == 1 {
		elem, _ := candidates[0]["instance_id"].(string)
		removeLoad(elem)
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "aged_franken_prayer_load",
		"苍老者 弗兰肯 拜利兰:选择失去1点负载", candidates, 1, 1,
		func(selected []string) {
			removeLoad(firstSelected(selected))
		})
	return nil
}
func (Card1411101AgedFrankenBaililan) OnLoadLoss(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil {
		return nil
	}
	lossTarget, _ := ctx.ExtraData["load_loss_target"].(*CardInstance)
	if lossTarget != ctx.Source || !ctx.Engine.cardStillOnField(ctx.Source) {
		return nil
	}
	candidates := ctx.Engine.enemySkills(ctx.PlayerID, nil)
	if len(candidates) == 0 {
		return nil
	}
	weakenSkill := func(selection string) {
		target, _ := ctx.Engine.findFriendlyCandidate(ctx.OpponentID, selection)
		if target == nil || target.Card == nil || !target.Card.IsSkill() {
			return
		}
		target.PowerBonus -= 2
		ctx.Engine.emit(GameEvent{
			Type:   "aged_franken_weaken_spell",
			Player: -1,
			Data: map[string]any{
				"player": ctx.PlayerID,
				"source": cardToInfo(ctx.Source),
				"target": cardToInfo(target),
				"power":  -2,
			},
		})
	}
	if len(candidates) == 1 {
		id, _ := candidates[0]["instance_id"].(string)
		weakenSkill(id)
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "aged_franken_weaken_spell",
		"苍老者 弗兰肯 拜利兰:选择1个敌方法术永久-2威", candidates, 1, 1,
		func(selected []string) {
			weakenSkill(firstSelected(selected))
		})
	return nil
}

var _ PrayerAbility = Card1411101AgedFrankenBaililan{}
var _ PerTurnAbility = Card1411101AgedFrankenBaililan{}
var _ OnLoadLossBehavior = Card1411101AgedFrankenBaililan{}

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
	for _, skill := range ps.SkillPool {
		if skill == nil || skill.Card == nil || !skill.Card.IsSkill() || totalElementCost(skill.Card.ElementsCost) >= 4 {
			continue
		}
		hasEmptySlot := false
		for i := 0; i < skillSlotCapacity(ps); i++ {
			if ps.Skills[i] == nil && skillAllowedInSlot(ps, skill, i) {
				hasEmptySlot = true
				break
			}
		}
		if hasEmptySlot {
			candidates = append(candidates, candidateInfo(skill, "skill_pool", "own"))
			allowed[skill.InstanceID] = true
			continue
		}
		for slotIdx, learned := range ps.Skills {
			if learned == nil || learned.IsHorizontal {
				continue
			}
			if !skillAllowedInSlot(ps, skill, slotIdx) {
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

type Card1221115WinterfellAntiMage struct{ AlwaysActive }

func (Card1221115WinterfellAntiMage) ID() string   { return "1221115" }
func (Card1221115WinterfellAntiMage) Name() string { return "凛冬城御魔师" }
func (Card1221115WinterfellAntiMage) OnPrayer(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	for _, skill := range ps.Skills {
		if skill == nil || skill.Card == nil || !skill.Card.IsSkill() {
			continue
		}
		ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
			Type:             TempModNextSkillUseCostMinus,
			SourceCardNumber: ctx.Source.Card.Number,
			SourceName:       ctx.Source.Card.Name,
			TargetInstanceID: skill.InstanceID,
			Element:          model.ElementWater,
			Amount:           1,
			RemainingUses:    1,
		})
	}
	return nil
}

type Card4211102WinterfellWarlockSophia struct{ AlwaysActive }

func (Card4211102WinterfellWarlockSophia) ID() string   { return "4211102" }
func (Card4211102WinterfellWarlockSophia) Name() string { return "凛冰魔巫 索菲娅" }
func (Card4211102WinterfellWarlockSophia) HasNegativeStatusImmunity(status string) bool {
	return status == StatusFreeze
}

func (Card4211102WinterfellWarlockSophia) OnUltimate(ctx *EffectContext) error {
	candidates := make([]map[string]any, 0)
	for playerID, ps := range ctx.Engine.State.Players {
		if ps == nil {
			continue
		}
		side := "enemy"
		if playerID == ctx.PlayerID {
			side = "own"
		}
		for _, unit := range ps.Units {
			for _, card := range unit {
				if card != nil && card.Card != nil && (card.Card.IsHero() || card.Card.IsCompanion()) && card.Statuses[StatusFreeze] > 0 {
					candidates = append(candidates, candidateInfo(card, "unit", side))
				}
			}
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	allowed := make(map[string]bool, len(candidates))
	for _, candidate := range candidates {
		if id, _ := candidate["instance_id"].(string); id != "" {
			allowed[id] = true
		}
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "sophia_thaw_strike",
		"凛冰魔巫 索菲娅:选择1个冻结单位移除1层冻结并造成2点伤害", candidates, 1, 1,
		func(selected []string) {
			id := firstSelected(selected)
			if !allowed[id] {
				return
			}
			target := ctx.Engine.findUnitByInstanceID(id)
			if target == nil || target.Card == nil || (!target.Card.IsHero() && !target.Card.IsCompanion()) || target.Statuses[StatusFreeze] <= 0 {
				return
			}
			target.Statuses[StatusFreeze]--
			ctx.Engine.dealDamageWithExtra(target, 2, target.OwnerID, map[string]any{
				"damage_source": "effect",
				"attacker":      ctx.PlayerID,
			})
		})
	return nil
}

type Card1421112SandDustDemon struct{ AlwaysActive }

func (Card1421112SandDustDemon) ID() string            { return "1421112" }
func (Card1421112SandDustDemon) Name() string          { return "沙尘恶魔" }
func (Card1421112SandDustDemon) IsPrayerAbility() bool { return true }
func (Card1421112SandDustDemon) OnPerTurn(ctx *EffectContext) error {
	opponent := ctx.Engine.State.Players[ctx.OpponentID]
	frontRow := opponent.GetFrontRow()
	if frontRow < 0 || frontRow >= 3 {
		return nil
	}
	for col := 0; col < 3; col++ {
		if unit := opponent.Units[col][frontRow]; unit != nil {
			ctx.Engine.addStatus(unit, StatusPetrify, 1)
		}
	}
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
		ctx.Engine.reduceCardElementLoadWithTriggers(ctx.PlayerID, target, elem, 1, ctx.Source)
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

func reducibleElementLoad(card *CardInstance, elem string) int {
	return dragonBloodTreantReducibleLoad(card)[elem]
}

func reduceCardElementLoad(card *CardInstance, elem string, amount int) int {
	removed := 0
	for i := 0; i < amount; i++ {
		if reducibleElementLoad(card, elem) <= 0 {
			return removed
		}
		dragonBloodTreantRemoveLoad(card, elem)
		removed++
	}
	return removed
}

func (e *Engine) reduceCardElementLoadWithTriggers(playerID int, card *CardInstance, elem string, amount int, cause *CardInstance) int {
	removed := reduceCardElementLoad(card, elem, amount)
	if e == nil || removed <= 0 || card == nil || card.Card == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return removed
	}
	data := map[string]any{
		"load_loss_player": playerID,
		"load_loss_target": card,
		"load_loss_source": cause,
		"element":          elem,
		"amount":           removed,
	}
	e.emit(GameEvent{
		Type:   "load_loss",
		Player: -1,
		Data: map[string]any{
			"player":  playerID,
			"target":  cardToInfo(card),
			"source":  cardToInfo(cause),
			"element": elem,
			"amount":  removed,
		},
	})
	e.triggerEffects(TriggerOnLoadLoss, card, cause, data)
	e.triggerFieldEffectsWithData(TriggerOnLoadLoss, playerID, card, data)
	e.triggerFieldEffectsWithData(TriggerOnLoadLoss, 1-playerID, card, data)
	return removed
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

type Card1321101SkycarrierE2 struct{ AlwaysActive }

func (Card1321101SkycarrierE2) ID() string            { return "1321101" }
func (Card1321101SkycarrierE2) Name() string          { return "翱翔者E2型运输舰" }
func (Card1321101SkycarrierE2) IsPrayerAbility() bool { return true }

func (Card1321101SkycarrierE2) OnPerTurn(ctx *EffectContext) error {
	choices := []map[string]any{
		{"instance_id": "draw", "name": "抽2张牌", "zone": "choice", "side": "own"},
	}
	if len(airGraveyardCandidates(ctx.Engine.State.Players[ctx.PlayerID])) >= 2 {
		choices = append(choices, map[string]any{"instance_id": "recycle", "name": "将2张大气弃牌洗回牌组", "zone": "choice", "side": "own"})
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "skycarrier_e2_prayer",
		"翱翔者E2型运输舰:选择祈咒效果", choices, 1, 1,
		func(selected []string) {
			switch firstSelected(selected) {
			case "draw":
				ctx.Engine.drawCards(ctx.PlayerID, 2)
			case "recycle":
				openSkycarrierRecyclePrompt(ctx)
			}
		})
	return nil
}

func openSkycarrierRecyclePrompt(ctx *EffectContext) {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	candidates := airGraveyardCandidates(ps)
	if len(candidates) < 2 {
		return
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "skycarrier_e2_recycle",
		"翱翔者E2型运输舰:选择2张大气弃牌洗回牌组", candidates, 2, 2,
		func(selected []string) {
			moveSelectedAirGraveyardCardsToDeck(ps, selected, 2)
			ctx.Engine.shuffleDeck(ctx.PlayerID)
		})
}

func airGraveyardCandidates(ps *PlayerState) []map[string]any {
	candidates := make([]map[string]any, 0)
	if ps == nil {
		return candidates
	}
	for _, card := range ps.Graveyard {
		if card != nil && card.Card != nil && card.Card.Category == model.ElementAir {
			candidates = append(candidates, candidateInfo(card, "graveyard", "own"))
		}
	}
	return candidates
}

func moveSelectedAirGraveyardCardsToDeck(ps *PlayerState, selected []string, maxCount int) {
	if ps == nil || maxCount <= 0 {
		return
	}
	selectedSet := make(map[string]bool, len(selected))
	for _, id := range selected {
		selectedSet[id] = true
	}
	moved := 0
	for i := 0; i < len(ps.Graveyard) && moved < maxCount; {
		card := ps.Graveyard[i]
		if card != nil && selectedSet[card.InstanceID] && card.Card != nil && card.Card.Category == model.ElementAir {
			ps.Graveyard = append(ps.Graveyard[:i], ps.Graveyard[i+1:]...)
			ps.Deck = append(ps.Deck, card)
			moved++
			continue
		}
		i++
	}
}

type Card3421101ForestInsight struct{ AlwaysActive }

func (Card3421101ForestInsight) ID() string   { return "3421101" }
func (Card3421101ForestInsight) Name() string { return "森之洞察" }
func (Card3421101ForestInsight) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || !isFriendlySpellCast(ctx) || !isSpellBeingCast(ctx) {
		return nil
	}
	if ctx.Source.Card.Number != "3421101" {
		return nil
	}
	drawCount := min(5, countFriendlyEarthCompanions(ctx.Engine, ctx.PlayerID))
	if drawCount <= 0 {
		return nil
	}
	drawn := ctx.Engine.drawCards(ctx.PlayerID, drawCount)
	shuffleBackCount := len(drawn)
	if shuffleBackCount <= 0 {
		return nil
	}
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, nil)
	if len(candidates) < shuffleBackCount {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "forest_insight_shuffle_hand",
		fmt.Sprintf("森之洞察:选择%d张手牌洗回卡组", shuffleBackCount), candidates, shuffleBackCount, shuffleBackCount,
		func(selected []string) {
			if moveSelectedHandCardsToDeck(ctx.Engine, ctx.PlayerID, selected, shuffleBackCount) > 0 {
				ctx.Engine.shuffleDeck(ctx.PlayerID)
			}
		})
	return nil
}

type Card1321105Illusionist struct{ AlwaysActive }

func (Card1321105Illusionist) ID() string   { return "1321105" }
func (Card1321105Illusionist) Name() string { return "幻术师" }
func (Card1321105Illusionist) OnUltimate(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() &&
			totalElementCost(card.Card.ElementsCost) < 6
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "illusionist_return_companion",
		"幻术师:选择1个入场花费小于6的友方伙伴移回手牌", candidates, 1, 1,
		func(selected []string) {
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target == nil || zone != "unit" || target.Card == nil || !target.Card.IsCompanion() || totalElementCost(target.Card.ElementsCost) >= 6 {
				return
			}
			gain := copyElementAmounts(effectiveElementsGain(target))
			ctx.Engine.returnUnitToHand(target, ctx.PlayerID)
			resetCardForHiddenZone(target)
			ctx.Engine.State.Players[ctx.PlayerID].GainElements(gain)
		})
	return nil
}

type Card1621115SoulDevourer struct{ AlwaysActive }

func (Card1621115SoulDevourer) ID() string   { return "1621115" }
func (Card1621115SoulDevourer) Name() string { return "灵魂吸食者" }
func (Card1621115SoulDevourer) PerTurnLabel(*CardInstance) string {
	return "主动"
}
func (Card1621115SoulDevourer) OnPerTurn(ctx *EffectContext) error {
	candidates := soulMarkedFriendlyFieldCandidates(ctx.Engine, ctx.PlayerID)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "soul_devourer_remove_marker",
		"灵魂吸食者:移除你场上的1个灵魂标记物，抽2张并获得2暗", candidates, 1, 1,
		func(selected []string) {
			target := findFriendlyFieldCardIncludingBoundSkill(ctx.Engine, ctx.PlayerID, firstSelected(selected))
			if target == nil || target.Statuses[soulMarkerStatus] <= 0 {
				return
			}
			removeSoulMarkerFromCard(target)
			ctx.Engine.drawCards(ctx.PlayerID, 2)
			ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementShadow: 2})
		})
	return nil
}

func copyElementAmounts(src map[string]int) map[string]int {
	copied := make(map[string]int, len(src))
	for elem, amount := range src {
		if amount > 0 {
			copied[elem] = amount
		}
	}
	return copied
}

func countFriendlyEarthCompanions(e *Engine, playerID int) int {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return 0
	}
	count := 0
	for _, card := range e.getAllFieldCards(e.State.Players[playerID]) {
		if card != nil && card.Card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementEarth {
			count++
		}
	}
	return count
}

func moveSelectedHandCardsToDeck(e *Engine, playerID int, selected []string, maxCount int) int {
	if e == nil || maxCount <= 0 || playerID < 0 || playerID >= len(e.State.Players) {
		return 0
	}
	ps := e.State.Players[playerID]
	selectedSet := make(map[string]bool, len(selected))
	for _, id := range selected {
		selectedSet[id] = true
	}
	moved := 0
	for i := 0; i < len(ps.Hand) && moved < maxCount; {
		card := ps.Hand[i]
		if card != nil && selectedSet[card.InstanceID] {
			ps.Hand = append(ps.Hand[:i], ps.Hand[i+1:]...)
			resetCardForHiddenZone(card)
			ps.Deck = append(ps.Deck, card)
			moved++
			continue
		}
		i++
	}
	return moved
}

func soulMarkedFriendlyFieldCandidates(e *Engine, playerID int) []map[string]any {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	ps := e.State.Players[playerID]
	candidates := make([]map[string]any, 0)
	for _, card := range e.getAllFieldCards(ps) {
		if card == nil {
			continue
		}
		if card.Statuses[soulMarkerStatus] > 0 {
			candidates = append(candidates, candidateInfo(card, "field", "own"))
		}
		for _, skill := range card.BoundSkills {
			if skill != nil && skill.Statuses[soulMarkerStatus] > 0 {
				candidates = append(candidates, candidateInfo(skill, "bound_skill", "own"))
			}
		}
	}
	for _, skill := range ps.Skills {
		if skill != nil && skill.Statuses[soulMarkerStatus] > 0 {
			candidates = append(candidates, candidateInfo(skill, "skill", "own"))
		}
	}
	return candidates
}

func findFriendlyFieldCardIncludingBoundSkill(e *Engine, playerID int, instanceID string) *CardInstance {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) || instanceID == "" {
		return nil
	}
	ps := e.State.Players[playerID]
	if card, zone := e.findFriendlyCandidate(playerID, instanceID); card != nil && (zone == "unit" || zone == "equipment" || zone == "skill") {
		return card
	}
	for _, card := range e.getAllFieldCards(ps) {
		if card == nil {
			continue
		}
		for _, skill := range card.BoundSkills {
			if skill != nil && skill.InstanceID == instanceID {
				return skill
			}
		}
	}
	return nil
}

func shadowCompanionGraveyardCandidates(ps *PlayerState) []map[string]any {
	candidates := make([]map[string]any, 0)
	if ps == nil {
		return candidates
	}
	for _, card := range ps.Graveyard {
		if card != nil && card.Card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementShadow {
			candidates = append(candidates, candidateInfo(card, "graveyard", "own"))
		}
	}
	return candidates
}

func moveSelectedShadowCompanionsFromGraveyardToExile(e *Engine, playerID int, selected []string, maxCount int) int {
	if e == nil || maxCount <= 0 || playerID < 0 || playerID >= len(e.State.Players) {
		return 0
	}
	ps := e.State.Players[playerID]
	selectedSet := make(map[string]bool, len(selected))
	for _, id := range selected {
		selectedSet[id] = true
	}
	moved := 0
	for _, card := range append([]*CardInstance(nil), ps.Graveyard...) {
		if moved >= maxCount {
			break
		}
		if card == nil || !selectedSet[card.InstanceID] || card.Card == nil || !card.Card.IsCompanion() || card.Card.Category != model.ElementShadow {
			continue
		}
		if e.exileCard(playerID, card) {
			moved++
		}
	}
	return moved
}

func isShadowSpellInstance(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.Category == model.ElementShadow && isSpellLikeCard(card.Card)
}

func addSoulMarkerToSpell(skill *CardInstance) {
	if skill == nil || skill.Card == nil || !isSpellLikeCard(skill.Card) {
		return
	}
	skill.Statuses[soulMarkerStatus]++
	skill.PowerBonus += 2
}

func removeSoulMarkerFromCard(card *CardInstance) {
	if card == nil || card.Statuses[soulMarkerStatus] <= 0 {
		return
	}
	card.Statuses[soulMarkerStatus]--
	if card.Card != nil && isSpellLikeCard(card.Card) {
		card.PowerBonus -= 2
	}
	if card.Statuses[soulMarkerStatus] <= 0 {
		delete(card.Statuses, soulMarkerStatus)
	}
}

const prayerFlameMarkerStatus = "祈祷之焰标记物"

type Card3121103PrayerFlame struct{ AlwaysActive }

func (Card3121103PrayerFlame) ID() string   { return "3121103" }
func (Card3121103PrayerFlame) Name() string { return "祈祷之焰" }
func (Card3121103PrayerFlame) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || !isFriendlySpellCast(ctx) || !isSpellBeingCast(ctx) {
		return nil
	}
	if ctx.Source.Card.Number != "3121103" {
		return nil
	}
	choices := []map[string]any{
		{"instance_id": "add_markers", "name": "放置3个标记物", "zone": "choice", "side": "own"},
	}
	if prayerFlameHasSummonTarget(ctx.Engine, ctx.PlayerID, ctx.Source.Statuses[prayerFlameMarkerStatus]) {
		choices = append(choices, map[string]any{"instance_id": "summon", "name": "取除标记物免费召唤火焰伙伴", "zone": "choice", "side": "own"})
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "prayer_flame_choice",
		"祈祷之焰:选择放置标记物或取除标记物召唤火焰伙伴", choices, 1, 1,
		func(selected []string) {
			if firstSelected(selected) != "summon" {
				ctx.Source.Statuses[prayerFlameMarkerStatus] += 3
				return
			}
			markers := ctx.Source.Statuses[prayerFlameMarkerStatus]
			if markers <= 0 {
				return
			}
			openPrayerFlameSummonPrompt(ctx, markers)
		})
	return nil
}

func prayerFlameHasSummonTarget(e *Engine, playerID int, markers int) bool {
	return markers > 0 && len(e.friendlyHandCards(playerID, func(card *CardInstance) bool {
		return isFireCompanionWithEntryCostAtMost(card, markers)
	})) > 0 && len(e.friendlyEmptyUnitPositions(playerID)) > 0
}

func openPrayerFlameSummonPrompt(ctx *EffectContext, markers int) {
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, func(card *CardInstance) bool {
		return isFireCompanionWithEntryCostAtMost(card, markers)
	})
	if len(candidates) == 0 {
		return
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "prayer_flame_summon_card",
		"祈祷之焰:选择1个火焰伙伴免费召唤", candidates, 1, 1,
		func(selected []string) {
			cardID := firstSelected(selected)
			positions := ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)
			if len(positions) == 0 {
				return
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "prayer_flame_summon_position",
				"祈祷之焰:选择召唤位置", positions, 1, 1,
				func(posSelected []string) {
					pos, ok := positionFromSelectionID(firstSelected(posSelected))
					if !ok {
						return
					}
					card := summonCardFreeFromHandOrDeckAtPosition(ctx, cardID, pos)
					if card != nil {
						delete(ctx.Source.Statuses, prayerFlameMarkerStatus)
					}
				})
		})
}

func isFireCompanionWithEntryCostAtMost(card *CardInstance, maxCost int) bool {
	return card != nil && card.Card != nil && card.Card.IsCompanion() &&
		card.Card.Category == model.ElementFire &&
		totalElementCost(card.Card.ElementsCost) <= maxCost
}

type Card1421106PhantomLizard struct{ AlwaysActive }

func (Card1421106PhantomLizard) ID() string   { return "1421106" }
func (Card1421106PhantomLizard) Name() string { return "幻影蜥蜴" }
func (Card1421106PhantomLizard) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.Card == nil || ctx.Source.UltimateUsed {
		return nil
	}
	if !isFriendlySpellCast(ctx) || !hasCardTag(ctx.Target.Card, "灵媒") || !ctx.Engine.canConsumeCard(ctx.Source) {
		return nil
	}
	if len(ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)) < 1 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "phantom_lizard_split",
		"幻影蜥蜴:消耗此卡并分裂为两个普通蜥蜴", []map[string]any{candidateInfo(ctx.Source, "unit", "own")}, 1, 1,
		func(selected []string) {
			if !ctx.Engine.canConsumeCard(ctx.Source) || !ctx.Engine.cardStillOnField(ctx.Source) {
				return
			}
			ctx.Source.UltimateUsed = true
			ctx.Engine.consumeCardForEffectWithTriggers(ctx.PlayerID, ctx.Source, ctx.Engine.effectiveElementsGain(ctx.Source), "")
			moveUnitToGraveyardWithoutDeath(ctx.Engine, ctx.PlayerID, ctx.Source)
			positions := ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)
			if len(positions) < 2 {
				return
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "phantom_lizard_first_position",
				"幻影蜥蜴:选择第1个普通蜥蜴的位置", positions, 1, 1,
				func(firstSelectedPos []string) {
					firstPos, ok := positionFromSelectionID(firstSelected(firstSelectedPos))
					if !ok || ctx.Engine.summonFreshCardAtPosition(ctx.PlayerID, "1401101", firstPos, true) == nil {
						return
					}
					secondPositions := ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)
					if len(secondPositions) == 0 {
						return
					}
					ctx.Engine.SetPendingAction(ctx.PlayerID, "phantom_lizard_second_position",
						"幻影蜥蜴:选择第2个普通蜥蜴的位置", secondPositions, 1, 1,
						func(secondSelectedPos []string) {
							secondPos, ok := positionFromSelectionID(firstSelected(secondSelectedPos))
							if ok {
								ctx.Engine.summonFreshCardAtPosition(ctx.PlayerID, "1401101", secondPos, true)
							}
						})
				})
		})
	return nil
}

func (e *Engine) consumeCardForEffectWithTriggers(playerID int, card *CardInstance, gains map[string]int, sourceNumber string) {
	if e == nil || card == nil || playerID < 0 || playerID >= len(e.State.Players) || !e.canConsumeCard(card) {
		return
	}
	gains = copyElementCost(gains)
	card.IsHorizontal = true
	ps := e.State.Players[playerID]
	ps.GainElements(gains)
	e.emit(GameEvent{
		Type:   "consume",
		Player: -1,
		Data: map[string]any{
			"player":      playerID,
			"instance_id": card.InstanceID,
			"elements":    ps.Elements,
			"gained":      gains,
		},
	})
	consumeData := map[string]any{
		"consumed_player": playerID,
		"gained":          gains,
	}
	if sourceNumber != "" {
		consumeData["consume_source"] = sourceNumber
	}
	e.triggerEffects(TriggerOnConsume, card, nil, consumeData)
	e.triggerFieldEffectsWithData(TriggerOnConsume, playerID, card, consumeData)
	e.triggerFieldEffectsWithData(TriggerOnConsume, 1-playerID, card, consumeData)
	e.advanceMastery(card, playerID, 1)
	e.destroyFuyeDoomedCardAfterExert(card)
}

func moveUnitToGraveyardWithoutDeath(e *Engine, playerID int, unit *CardInstance) {
	if e == nil || unit == nil || unit.Position == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	ps := e.State.Players[playerID]
	if ps.Units[unit.Position.Col][unit.Position.Row] != unit {
		return
	}
	ps.Units[unit.Position.Col][unit.Position.Row] = nil
	unit.Position = nil
	e.releaseUnderCardsToGraveyard(playerID, unit)
	e.exileTransferredBoundSkills(playerID, unit)
	unit.BoundSkills = nil
	e.addToGraveyard(playerID, unit)
	e.emit(GameEvent{Type: "unit_transformed", Player: -1, Data: map[string]any{
		"player": playerID,
		"card":   cardToInfo(unit),
	}})
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
					ctx.Engine.triggerHolyChildAfterLifeGain(ctx.PlayerID, ctx.Source)
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
			ctx.Engine.consumeCardForEffectWithTriggers(ctx.PlayerID, target, target.Card.ElementsCost, "2121108")
		})
	return nil
}

const beastTamingCollarTargetPrefix = "驯兽项圈目标:"

type Card2121106BeastTamingCollar struct{ AlwaysActive }

func (Card2121106BeastTamingCollar) ID() string   { return "2121106" }
func (Card2121106BeastTamingCollar) Name() string { return "驯兽项圈" }
func (Card2121106BeastTamingCollar) OnEnter(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, isCollarEligibleFireCompanion)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "beast_taming_collar_target",
		"驯兽项圈:选择1个巫师以外的火焰伙伴", candidates, 1, 1,
		func(selected []string) {
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target == nil || zone != "unit" || !isCollarEligibleFireCompanion(target) {
				return
			}
			clearStatusPrefix(ctx.Source, beastTamingCollarTargetPrefix)
			ctx.Source.Statuses[beastTamingCollarTargetPrefix+target.InstanceID] = 1
		})
	return nil
}
func (Card2121106BeastTamingCollar) PerTurnLabel(*CardInstance) string {
	return "消耗目标"
}
func (Card2121106BeastTamingCollar) OnPerTurn(ctx *EffectContext) error {
	if !ctx.Engine.equipmentInOwnerSlot(ctx.PlayerID, ctx.Source) {
		return nil
	}
	target := collarTarget(ctx.Engine, ctx.PlayerID, ctx.Source)
	if target == nil || !ctx.Engine.canConsumeCard(target) {
		return nil
	}
	ctx.Engine.consumeCardForEffectWithTriggers(ctx.PlayerID, target, target.Card.ElementsCost, "2121106")
	return nil
}

func isCollarEligibleFireCompanion(card *CardInstance) bool {
	return isFireCompanion(card) && !hasCardTag(card.Card, "巫师")
}

func collarTarget(e *Engine, playerID int, collar *CardInstance) *CardInstance {
	if e == nil || collar == nil {
		return nil
	}
	for status, amount := range collar.Statuses {
		if amount <= 0 || !strings.HasPrefix(status, beastTamingCollarTargetPrefix) {
			continue
		}
		target, zone := e.findFriendlyCandidate(playerID, strings.TrimPrefix(status, beastTamingCollarTargetPrefix))
		if target != nil && zone == "unit" && isCollarEligibleFireCompanion(target) {
			return target
		}
	}
	return nil
}

type Card2421102RoseWhip struct{ AlwaysActive }

func (Card2421102RoseWhip) ID() string   { return "2421102" }
func (Card2421102RoseWhip) Name() string { return "蔷薇之鞭" }
func (Card2421102RoseWhip) OnLoadLoss(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.ExtraData == nil {
		return nil
	}
	lossPlayer, _ := ctx.ExtraData["load_loss_player"].(int)
	if lossPlayer != ctx.PlayerID || ctx.Target.OwnerID != ctx.PlayerID || ctx.Target == ctx.Source || !ctx.Engine.cardStillOnField(ctx.Source) {
		return nil
	}
	currentBonus := ctx.Source.ElementsGainBonus[model.ElementShadow]
	if currentBonus >= 2 {
		return nil
	}
	ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementShadow, 1, ctx.Source)
	ctx.Engine.emit(GameEvent{
		Type:   "rose_whip_load_gain",
		Player: -1,
		Data: map[string]any{
			"player": ctx.PlayerID,
			"source": cardToInfo(ctx.Source),
			"target": cardToInfo(ctx.Target),
		},
	})
	return nil
}

var _ OnLoadLossBehavior = Card2421102RoseWhip{}

type Card2421104BloodRoseContract struct{ AlwaysActive }

func (Card2421104BloodRoseContract) ID() string   { return "2421104" }
func (Card2421104BloodRoseContract) Name() string { return "血蔷薇契约" }
func (Card2421104BloodRoseContract) OnUseItem(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil {
		return nil
	}
	spellCandidates := ctx.Engine.friendlySkills(ctx.PlayerID, func(skill *CardInstance) bool {
		return skill != nil && skill.Card != nil && isSpellLikeCard(skill.Card)
	})
	if len(spellCandidates) == 0 {
		return nil
	}
	hostCandidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() &&
			(card.Card.Category == model.ElementEarth || card.Card.Category == model.ElementShadow)
	})
	if len(hostCandidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "blood_rose_contract_spell",
		"血蔷薇契约:选择要绑定的己方法术", spellCandidates, 1, 1,
		func(selected []string) {
			skillID := firstSelected(selected)
			skill, skillIndex := findSkillSlotByInstance(ctx.Engine.State.Players[ctx.PlayerID], skillID)
			if skill == nil || skill.Card == nil || !isSpellLikeCard(skill.Card) {
				return
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "blood_rose_contract_host",
				"血蔷薇契约:选择地脉或暗影伙伴作为绑定宿主", hostCandidates, 1, 1,
				func(hostSelected []string) {
					host, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(hostSelected))
					if host == nil || zone != "unit" || host.Card == nil || !host.Card.IsCompanion() ||
						(host.Card.Category != model.ElementEarth && host.Card.Category != model.ElementShadow) {
						return
					}
					currentSkill, currentIndex := findSkillSlotByInstance(ctx.Engine.State.Players[ctx.PlayerID], skillID)
					if currentSkill != skill || currentIndex != skillIndex {
						return
					}
					bonus := ctx.Engine.totalLoad(host)
					ctx.Engine.State.Players[ctx.PlayerID].Skills[skillIndex] = nil
					skill.SlotIndex = -1
					markTransferredBoundSkill(skill)
					skill.PowerBonus += bonus
					host.BoundSkills = append(host.BoundSkills, skill)
					ctx.Engine.emit(GameEvent{
						Type:   "blood_rose_contract_bind",
						Player: -1,
						Data: map[string]any{
							"player":      ctx.PlayerID,
							"skill":       cardToInfo(skill),
							"host":        cardToInfo(host),
							"power_bonus": bonus,
						},
					})
				})
		})
	return nil
}

func findSkillSlotByInstance(ps *PlayerState, instanceID string) (*CardInstance, int) {
	if ps == nil || instanceID == "" {
		return nil, -1
	}
	for i, skill := range ps.Skills {
		if skill != nil && skill.InstanceID == instanceID {
			return skill, i
		}
	}
	return nil, -1
}

var _ OnUseItemBehavior = Card2421104BloodRoseContract{}

type Card2421105NaturalCommunion struct{ AlwaysActive }

func (Card2421105NaturalCommunion) ID() string   { return "2421105" }
func (Card2421105NaturalCommunion) Name() string { return "自然交感" }
func (Card2421105NaturalCommunion) OnUseItem(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil {
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementEarth
	})
	if len(candidates) < 2 {
		return nil
	}
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "natural_communion_distribute",
		"自然交感:选择2个地脉伙伴并提交新的负载分配", candidates, 2, 2, nil, false,
		func(selected []string, data map[string]any) error {
			a, _ := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, selected[0])
			b, _ := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, selected[1])
			if !isEarthCompanionUnit(a) || !isEarthCompanionUnit(b) || a == b {
				return fmt.Errorf("自然交感需要2个不同的地脉伙伴")
			}
			distribution, err := parseLoadDistribution(data["load_distribution"])
			if err != nil {
				return err
			}
			aLoad, okA := distribution[a.InstanceID]
			bLoad, okB := distribution[b.InstanceID]
			if !okA || !okB || len(distribution) != 2 {
				return fmt.Errorf("自然交感需要为选择的2个伙伴提交负载分配")
			}
			totalBefore := mergeElementCosts(ctx.Engine.effectiveElementsGain(a), ctx.Engine.effectiveElementsGain(b))
			totalAfter := mergeElementCosts(aLoad, bLoad)
			if !sameElementCost(totalBefore, totalAfter) {
				return fmt.Errorf("自然交感的新负载分配必须保持总负载不变")
			}
			applyRedistributedLoad(a, aLoad)
			applyRedistributedLoad(b, bLoad)
			ctx.Engine.emit(GameEvent{
				Type:   "natural_communion_distribute",
				Player: -1,
				Data: map[string]any{
					"player": ctx.PlayerID,
					"cards":  []any{cardToInfo(a), cardToInfo(b)},
				},
			})
			return nil
		})
	return nil
}

func isEarthCompanionUnit(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementEarth
}

func parseLoadDistribution(raw any) (map[string]map[string]int, error) {
	input, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("missing load_distribution")
	}
	result := make(map[string]map[string]int, len(input))
	for instanceID, value := range input {
		elemMap, ok := value.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid load distribution")
		}
		load := make(map[string]int)
		for elem, amountRaw := range elemMap {
			if !isNonArcaneElement(elem) && elem != model.ElementArcane {
				return nil, fmt.Errorf("invalid load element")
			}
			amount, ok := intFromAny(amountRaw)
			if !ok || amount < 0 {
				return nil, fmt.Errorf("invalid load amount")
			}
			if amount > 0 {
				load[elem] = amount
			}
		}
		result[instanceID] = load
	}
	return result, nil
}

func sameElementCost(a, b map[string]int) bool {
	for _, elem := range model.AllElements {
		if a[elem] != b[elem] {
			return false
		}
	}
	return true
}

func intFromAny(raw any) (int, bool) {
	switch v := raw.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		i := int(v)
		return i, float64(i) == v
	case float32:
		i := int(v)
		return i, float32(i) == v
	default:
		return 0, false
	}
}

func applyRedistributedLoad(card *CardInstance, load map[string]int) {
	if card == nil {
		return
	}
	card.ElementsGainBonus = make(map[string]int)
	setElementsGain(card, load)
}

var _ OnUseItemBehavior = Card2421105NaturalCommunion{}

type Card2421106AgingPotion struct{ AlwaysActive }

func (Card2421106AgingPotion) ID() string   { return "2421106" }
func (Card2421106AgingPotion) Name() string { return "苍老药剂" }
func (Card2421106AgingPotion) OnUseItem(ctx *EffectContext) error {
	candidates := friendlyFieldCardsIncludingBound(ctx.Engine, ctx.PlayerID, func(card *CardInstance) bool {
		behavior, ok := masteryBehavior(card)
		return ok && card.Statuses[StatusMastery] < behavior.MasteryMax() && reducibleElementLoad(card, model.ElementEarth) > 0
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "aging_potion_mastery",
		"苍老药剂:移除1点地负载并立刻达到下一次精通", candidates, 1, 1,
		func(selected []string) {
			target := ctx.Engine.findFriendlyCardIncludingBound(ctx.PlayerID, firstSelected(selected))
			behavior, ok := masteryBehavior(target)
			if target == nil || !ok || target.Statuses[StatusMastery] >= behavior.MasteryMax() || reducibleElementLoad(target, model.ElementEarth) <= 0 {
				return
			}
			ctx.Engine.reduceCardElementLoadWithTriggers(ctx.PlayerID, target, model.ElementEarth, 1, ctx.Source)
			ctx.Engine.advanceMastery(target, ctx.PlayerID, 1)
		})
	return nil
}

func clearStatusPrefix(card *CardInstance, prefix string) {
	if card == nil || prefix == "" {
		return
	}
	for status := range card.Statuses {
		if strings.HasPrefix(status, prefix) {
			delete(card.Statuses, status)
		}
	}
}

func friendlyFieldCardsIncludingBound(e *Engine, playerID int, predicate func(*CardInstance) bool) []map[string]any {
	ps := e.State.Players[playerID]
	candidates := make([]map[string]any, 0)
	for _, card := range e.getAllFieldCards(ps) {
		if card == nil || predicate != nil && !predicate(card) {
			continue
		}
		candidates = append(candidates, candidateInfo(card, "field", "own"))
		for _, skill := range card.BoundSkills {
			if skill == nil || predicate != nil && !predicate(skill) {
				continue
			}
			candidates = append(candidates, candidateInfo(skill, "bound_skill", "own"))
		}
	}
	return candidates
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

func (e *Engine) firstFreeEquipmentSlot(playerID int) int {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return -1
	}
	ps := e.State.Players[playerID]
	for i := 0; i < equipmentSlotCapacity(ps); i++ {
		if ps.Equipment[i] == nil {
			return i
		}
	}
	return -1
}

func (e *Engine) freeEquipmentSlots(playerID int) int {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return 0
	}
	ps := e.State.Players[playerID]
	free := 0
	for i := 0; i < equipmentSlotCapacity(ps); i++ {
		if ps.Equipment[i] == nil {
			free++
		}
	}
	return free
}

func (e *Engine) equipGeneratedCard(playerID int, number string) *CardInstance {
	if e == nil || number == "" || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	cardDef, ok := cards.PlayableCardDB[number]
	if !ok || !cardDef.IsItem() || !isEquipmentCard(cardDef) {
		return nil
	}
	slot := e.firstFreeEquipmentSlot(playerID)
	if slot < 0 {
		return nil
	}
	card := NewCardInstance(cardDef, playerID, e.State.TurnNumber)
	card.IsHorizontal = true
	card.SlotIndex = slot
	ps := e.State.Players[playerID]
	ps.Equipment[slot] = card
	e.emit(GameEvent{
		Type:   "equip",
		Player: -1,
		Data: map[string]any{
			"player": playerID,
			"card":   cardToInfo(card),
			"slot":   slot,
		},
	})
	e.triggerEffects(TriggerOnEquip, card, nil, nil)
	e.triggerEffects(TriggerOnEnter, card, nil, nil)
	e.notifyCardEntered(playerID, card, map[string]any{"entered_player": playerID, "equipped": true})
	return card
}

func (e *Engine) equipCardFromHandOrDeckFree(playerID int, number string) *CardInstance {
	if e == nil || number == "" || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	slot := e.firstFreeEquipmentSlot(playerID)
	if slot < 0 {
		return nil
	}
	ps := e.State.Players[playerID]
	var card *CardInstance
	for i, handCard := range ps.Hand {
		if handCard != nil && handCard.Card != nil && handCard.Card.Number == number {
			card = handCard
			ps.RemoveFromHand(i)
			break
		}
	}
	if card == nil {
		for i, deckCard := range ps.Deck {
			if deckCard != nil && deckCard.Card != nil && deckCard.Card.Number == number {
				card = deckCard
				ps.Deck = append(ps.Deck[:i], ps.Deck[i+1:]...)
				break
			}
		}
	}
	if card == nil || card.Card == nil || !card.Card.IsItem() || !isEquipmentCard(card.Card) {
		return nil
	}
	card.OwnerID = playerID
	card.Position = nil
	card.IsSetCounter = false
	card.IsHorizontal = true
	card.SlotIndex = slot
	card.EnterTurn = e.State.TurnNumber
	ps.Equipment[slot] = card
	e.emit(GameEvent{
		Type:   "equip",
		Player: -1,
		Data: map[string]any{
			"player": playerID,
			"card":   cardToInfo(card),
			"slot":   slot,
			"free":   true,
		},
	})
	e.triggerEffects(TriggerOnEquip, card, nil, nil)
	e.triggerEffects(TriggerOnEnter, card, nil, nil)
	e.notifyCardEntered(playerID, card, map[string]any{"entered_player": playerID, "equipped": true, "free": true})
	return card
}

func (e *Engine) setHandCounterTrapFree(playerID int, card *CardInstance) bool {
	if e == nil || card == nil || card.Card == nil || playerID < 0 || playerID >= len(e.State.Players) || !isCounterTrapCard(card.Card.Number) {
		return false
	}
	ps := e.State.Players[playerID]
	_, handIdx := ps.FindHandCard(card.InstanceID)
	if handIdx < 0 || e.firstFreeEquipmentSlot(playerID) < 0 {
		return false
	}
	return e.placeCounterTrap(playerID, card, handIdx) == nil
}

const entryCostZeroStatus = "入场费用变为0"

func makeEntryCostZero(card *CardInstance) {
	if card == nil || card.Card == nil {
		return
	}
	if card.Statuses == nil {
		card.Statuses = make(map[string]int)
	}
	card.Statuses[entryCostZeroStatus] = 1
	for elem, amount := range card.Card.ElementsCost {
		if amount > 0 {
			card.Statuses["入场费用"+elem+"-"+fmt.Sprint(amount)]++
		}
	}
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

type Card2621110AndisGift struct{ AlwaysActive }

func (Card2621110AndisGift) ID() string   { return "2621110" }
func (Card2621110AndisGift) Name() string { return "安迪斯的赠与" }

const andisGiftDoomedStatus = "安迪斯的赠与回合结束死亡"

func (Card2621110AndisGift) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "andis_gift_target",
		"安迪斯的赠与:选择1个友方单位获得负载+2暗,回合结束时死亡", candidates, 1, 1,
		func(selected []string) {
			target := ctx.Engine.findUnitByInstanceID(firstSelected(selected))
			if target == nil || target.OwnerID != ctx.PlayerID || target.Position == nil {
				return
			}
			ctx.Engine.addElementsGainBonus(target, ctx.PlayerID, model.ElementShadow, 2, ctx.Source)
			target.Statuses[andisGiftDoomedStatus] = ctx.Engine.State.TurnNumber
		})
	return nil
}

func (e *Engine) destroyAndisGiftDoomedUnits(ps *PlayerState) {
	if e == nil || ps == nil {
		return
	}
	for _, card := range append([]*CardInstance(nil), e.getAllFieldCards(ps)...) {
		if card == nil || card.Card == nil || card.Statuses[andisGiftDoomedStatus] <= 0 {
			continue
		}
		delete(card.Statuses, andisGiftDoomedStatus)
		if card.Card.IsHero() {
			card.CurrentLife = 0
			e.emit(GameEvent{Type: "unit_destroyed", Player: -1, Data: map[string]any{
				"player": ps.PlayerID,
				"card":   cardToInfo(card),
				"reason": "andis_gift",
			}})
			e.checkWinCondition()
			continue
		}
		e.destroyUnitWithCause(card, ps.PlayerID, "andis_gift")
	}
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

type Card2621101BlackPineWand struct{ AlwaysActive }

func (Card2621101BlackPineWand) ID() string   { return "2621101" }
func (Card2621101BlackPineWand) Name() string { return "黑松木魔杖" }
func (Card2621101BlackPineWand) ModifySkillUseCost(ctx *EffectContext, cost map[string]int) {
	if ctx == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.ExtraData == nil {
		return
	}
	target, _ := ctx.ExtraData["spell_target"].(SpellTarget)
	targetUnit, _ := ctx.ExtraData["spell_target_unit"].(*CardInstance)
	if target.Type != "unit" || targetUnit == nil || targetUnit.OwnerID != ctx.PlayerID || ctx.Source.Card.Category == "" {
		return
	}
	reduceCost(cost, ctx.Source.Card.Category, 1)
}

var _ SkillUseCostModifier = Card2621101BlackPineWand{}

type Card2621102BloodRoseCurse struct{ AlwaysActive }

func (Card2621102BloodRoseCurse) ID() string   { return "2621102" }
func (Card2621102BloodRoseCurse) Name() string { return "血蔷薇诅咒" }
func (Card2621102BloodRoseCurse) OnUseItem(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil {
		return nil
	}
	opponentID := ctx.OpponentID
	spellCandidates := ctx.Engine.enemySkills(ctx.PlayerID, func(skill *CardInstance) bool {
		return skill != nil && skill.Card != nil && isSpellLikeCard(skill.Card)
	})
	if len(spellCandidates) == 0 {
		return nil
	}
	hostCandidates := ctx.Engine.enemyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion()
	})
	if len(hostCandidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "blood_rose_curse_spell",
		"血蔷薇诅咒:选择敌方1个法术", spellCandidates, 1, 1,
		func(selected []string) {
			skillID := firstSelected(selected)
			skill, skillIndex := findSkillSlotByInstance(ctx.Engine.State.Players[opponentID], skillID)
			if skill == nil || skill.Card == nil || !isSpellLikeCard(skill.Card) {
				return
			}
			ctx.Engine.SetPendingAction(opponentID, "blood_rose_curse_host",
				"血蔷薇诅咒:选择你的1个伙伴作为绑定宿主", hostCandidates, 1, 1,
				func(hostSelected []string) {
					host, zone := ctx.Engine.findFriendlyCandidate(opponentID, firstSelected(hostSelected))
					if host == nil || zone != "unit" || host.Card == nil || !host.Card.IsCompanion() {
						return
					}
					currentSkill, currentIndex := findSkillSlotByInstance(ctx.Engine.State.Players[opponentID], skillID)
					if currentSkill != skill || currentIndex != skillIndex {
						return
					}
					ctx.Engine.State.Players[opponentID].Skills[skillIndex] = nil
					skill.SlotIndex = -1
					markTransferredBoundSkill(skill)
					host.BoundSkills = append(host.BoundSkills, skill)
					ctx.Engine.emit(GameEvent{
						Type:   "blood_rose_curse_bind",
						Player: -1,
						Data: map[string]any{
							"player": ctx.PlayerID,
							"skill":  cardToInfo(skill),
							"host":   cardToInfo(host),
						},
					})
				})
		})
	return nil
}

var _ OnUseItemBehavior = Card2621102BloodRoseCurse{}

type Card2221108WesternChart struct{ AlwaysActive }

func (Card2221108WesternChart) ID() string   { return "2221108" }
func (Card2221108WesternChart) Name() string { return "西境航海图" }

const westernChartPierceTargetPrefix = "西境航海图穿透目标:"

func (Card2221108WesternChart) OnEnter(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlySkillsIncludingBound(ctx.PlayerID, func(skill *CardInstance) bool {
		return skill != nil && skill.Card != nil && isSpellLikeCard(skill.Card) && skill.Card.Category == model.ElementWater
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
	ctx.Engine.SetPendingAction(ctx.PlayerID, "western_chart_pierce_target",
		"西境航海图:选择你的1个水纹法术获得穿透", candidates, 1, 1,
		func(selected []string) {
			id := firstSelected(selected)
			if !allowed[id] {
				return
			}
			skill := ctx.Engine.findSkill(ctx.Engine.State.Players[ctx.PlayerID], id)
			if skill == nil || skill.Card == nil || !isSpellLikeCard(skill.Card) || skill.Card.Category != model.ElementWater {
				return
			}
			for status := range ctx.Source.Statuses {
				if strings.HasPrefix(status, westernChartPierceTargetPrefix) {
					delete(ctx.Source.Statuses, status)
				}
			}
			ctx.Source.Statuses[westernChartPierceTargetPrefix+skill.InstanceID] = 1
		})
	return nil
}

type Card4211101CoralBelly struct{ AlwaysActive }

func (Card4211101CoralBelly) ID() string   { return "4211101" }
func (Card4211101CoralBelly) Name() string { return "海神之使 珊瑚 贝莉" }

const coralBellyFirstSpellAttackUsedStatus = "海神之使首次法术攻击已触发"

func (e *Engine) applyCoralBellyFirstSpellAttackBonus(playerID int, skill *CardInstance) {
	if e == nil || skill == nil || skill.Card == nil || !isSpellLikeCard(skill.Card) || isSorcerySkill(skill.Card) {
		return
	}
	ps := e.State.Players[playerID]
	if ps == nil {
		return
	}
	for _, card := range e.getAllFieldCards(ps) {
		if card == nil || card.Card == nil || card.Card.Number != "4211101" || e.hasEffectiveStatus(card, StatusPetrify) {
			continue
		}
		if card.Statuses[coralBellyFirstSpellAttackUsedStatus] > 0 {
			continue
		}
		card.Statuses[coralBellyFirstSpellAttackUsedStatus] = 1
		skill.PowerBonus += 3
		e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
			"source": cardToInfo(card),
			"target": cardToInfo(skill),
			"effect": "first_spell_attack_power_bonus",
			"amount": 3,
		}})
		return
	}
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
			ctx.Engine.exileTransferredBoundSkills(ctx.PlayerID, ctx.Source)
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

type Card2621112SoulStaff struct{ AlwaysActive }

func (Card2621112SoulStaff) ID() string   { return "2621112" }
func (Card2621112SoulStaff) Name() string { return "灵魂法杖" }
func (Card2621112SoulStaff) PerTurnLabel(*CardInstance) string {
	return "主动"
}
func (Card2621112SoulStaff) OnPerTurn(ctx *EffectContext) error {
	graveyardCandidates := shadowCompanionGraveyardCandidates(ctx.Engine.State.Players[ctx.PlayerID])
	if len(graveyardCandidates) < 2 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "soul_staff_exile_companions",
		"灵魂法杖:选择2张暗影伙伴移出游戏", graveyardCandidates, 2, 2,
		func(selected []string) {
			if moveSelectedShadowCompanionsFromGraveyardToExile(ctx.Engine, ctx.PlayerID, selected, 2) < 2 {
				return
			}
			spellCandidates := ctx.Engine.friendlySkillsIncludingBound(ctx.PlayerID, isShadowSpellInstance)
			if len(spellCandidates) == 0 {
				return
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "soul_staff_mark_spell",
				"灵魂法杖:选择1个暗影法术放置1个灵魂标记物", spellCandidates, 1, 1,
				func(spellSelected []string) {
					skill := findFriendlySkillIncludingBound(ctx.Engine, ctx.PlayerID, firstSelected(spellSelected))
					if !isShadowSpellInstance(skill) {
						return
					}
					addSoulMarkerToSpell(skill)
				})
		})
	return nil
}

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
	addSoulMarkerToSpell(skill)
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
				addSoulMarkerToSpell(skill)
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

type Card3521108Grace struct{ AlwaysActive }

func (Card3521108Grace) ID() string   { return "3521108" }
func (Card3521108Grace) Name() string { return "恩典" }
func (Card3521108Grace) OnSpellCast(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() && card.CurrentLife < maxLife(card)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "grace_heal_companion",
		"恩典:选择1个受伤友方伙伴回复2血", candidates, 1, 1,
		func(selected []string) {
			target := ctx.Engine.findUnitByInstanceID(firstSelected(selected))
			if target == nil || target.OwnerID != ctx.PlayerID || target.Card == nil || !target.Card.IsCompanion() || target.Card.IsHero() || target.CurrentLife >= maxLife(target) {
				return
			}
			healUnit(target, 2)
			if target.CurrentLife >= maxLife(target) {
				target.Statuses["max_life_bonus"]++
				target.CurrentLife++
				ctx.Engine.addElementsGainBonus(target, ctx.PlayerID, model.ElementLight, 1, ctx.Source)
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

type Card2201103DreamRipple struct{ AlwaysActive }

func (Card2201103DreamRipple) ID() string   { return "2201103" }
func (Card2201103DreamRipple) Name() string { return "幻创之梦-波纹" }
func (Card2201103DreamRipple) OnUseItem(ctx *EffectContext) error {
	candidates := frontRowEnemyCandidates(ctx)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "dream_ripple_damage",
		"幻创之梦-波纹:选择前排敌人分配共计3点伤害", candidates, 1, min(3, len(candidates)),
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			allocations := map[string]int{}
			order := make([]string, 0, len(selected))
			for _, id := range selected {
				if allocations[id] == 0 {
					order = append(order, id)
				}
				allocations[id]++
			}
			for remaining := 3 - len(selected); remaining > 0 && len(order) > 0; remaining-- {
				allocations[order[0]]++
			}
			for _, id := range order {
				target := ctx.Engine.findUnitByInstanceID(id)
				if target == nil || target.OwnerID != ctx.OpponentID || target.Position == nil || !isCurrentFrontRowUnit(ctx.Engine.State.Players[ctx.OpponentID], target) {
					continue
				}
				ctx.Engine.dealDamageWithExtra(target, allocations[id], target.OwnerID, map[string]any{
					"damage_source": "effect",
					"attacker":      ctx.PlayerID,
				})
			}
		})
	return nil
}

type Card1621108DemonChild struct{ AlwaysActive }

func (Card1621108DemonChild) ID() string   { return "1621108" }
func (Card1621108DemonChild) Name() string { return "恶魔之子" }
func (Card1621108DemonChild) DevourCardRequirement() DevourCardRequirement {
	return DevourCardRequirement{Count: 1, Category: model.ElementShadow, CompanionOnly: true}
}

type Card3221108SixPetalSnowflake struct{ AlwaysActive }

func (Card3221108SixPetalSnowflake) ID() string   { return "3221108" }
func (Card3221108SixPetalSnowflake) Name() string { return "六瓣雪花" }
func (Card3221108SixPetalSnowflake) SpellHitStatuses(ctx *EffectContext) map[string]int {
	if ctx == nil || ctx.Target == nil || ctx.Target.Card == nil || ctx.Target.Card.IsHero() {
		return nil
	}
	return map[string]int{StatusFreeze: 1}
}

type Card3321105SweepingWind struct{ AlwaysActive }

func (Card3321105SweepingWind) ID() string   { return "3321105" }
func (Card3321105SweepingWind) Name() string { return "风卷残云" }
func (Card3321105SweepingWind) OnDamaged(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Target == nil || ctx.Target.Position == nil || ctx.Target.CurrentLife != 1 {
		return nil
	}
	ctx.Engine.destroyUnit(ctx.Target, ctx.Target.OwnerID)
	return nil
}

type Card3121107WarTrample struct{ AlwaysActive }

func (Card3121107WarTrample) ID() string   { return "3121107" }
func (Card3121107WarTrample) Name() string { return "战争践踏" }
func (Card3121107WarTrample) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Source == nil || ctx.Target != ctx.Source || ctx.ExtraData["purpose"] != string(skillPurposeAttack) || ctx.ExtraData["stat"] != "damage" {
		return
	}
	stats.DamageBonus -= spellTargetUnitCount(ctx.ExtraData)
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

func frontRowEnemyCandidates(ctx *EffectContext) []map[string]any {
	if ctx == nil || ctx.Engine == nil || ctx.OpponentID < 0 || ctx.OpponentID >= len(ctx.Engine.State.Players) {
		return nil
	}
	opponent := ctx.Engine.State.Players[ctx.OpponentID]
	frontRow := opponent.GetFrontRow()
	if frontRow < 0 || frontRow >= 3 {
		return nil
	}
	candidates := make([]map[string]any, 0, 3)
	for col := 0; col < 3; col++ {
		if unit := opponent.Units[col][frontRow]; unit != nil {
			candidates = append(candidates, candidateInfo(unit, "unit", "enemy"))
		}
	}
	return candidates
}

func isCurrentFrontRowUnit(ps *PlayerState, card *CardInstance) bool {
	if ps == nil || card == nil || card.Position == nil {
		return false
	}
	frontRow := ps.GetFrontRow()
	return frontRow >= 0 && card.Position.Row == frontRow
}

func spellTargetUnitCount(data map[string]any) int {
	if units, ok := data["affected_units"].([]*CardInstance); ok {
		return len(units)
	}
	if targets, ok := data["spell_targets"].([]SpellTarget); ok {
		count := 0
		for _, target := range targets {
			if target.Type == "unit" {
				count++
			}
		}
		return count
	}
	if target, ok := data["spell_target"].(SpellTarget); ok && target.Type == "unit" {
		return 1
	}
	return 0
}

func addGeneratedCardToPlayerHand(ctx *EffectContext, playerID int, cardNumber string) *CardInstance {
	card := getCardDB()[cardNumber]
	if card == nil {
		return nil
	}
	instance := NewCardInstance(card, playerID, ctx.Engine.State.TurnNumber)
	ctx.Engine.addCardToHand(playerID, instance)
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
	return len(e.discardSelectedHandCardsMatching(playerID, selected, limit, nil))
}

func (e *Engine) discardSelectedHandCardsMatching(playerID int, selected []string, limit int, predicate func(*CardInstance) bool) []*CardInstance {
	if playerID < 0 || playerID >= len(e.State.Players) || limit <= 0 {
		return nil
	}
	ps := e.State.Players[playerID]
	selectedSet := map[string]bool{}
	for _, id := range selected {
		if id != "" {
			selectedSet[id] = true
		}
	}
	discarded := make([]*CardInstance, 0, limit)
	for i := len(ps.Hand) - 1; i >= 0 && len(discarded) < limit; i-- {
		card := ps.Hand[i]
		if card == nil || !selectedSet[card.InstanceID] {
			continue
		}
		if predicate != nil && !predicate(card) {
			continue
		}
		if discardedCard := e.discardHandCardAt(playerID, i); discardedCard != nil {
			discarded = append(discarded, discardedCard)
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
	if card.Card.Number == "1321102" {
		e.offerDiscardedSpeckledSparrowSummon(playerID, card)
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
		e.notifyCardEntered(playerID, card, map[string]any{"entered_player": playerID})
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

const dragonSnowfieldTriggerStatus = "dragon_snowfield_trigger_count"

type Card3211102DragonSnowfield struct{ AlwaysActive }

func (Card3211102DragonSnowfield) ID() string   { return "3211102" }
func (Card3211102DragonSnowfield) Name() string { return "龙吟雪域" }

func (Card3211102DragonSnowfield) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	casterID, ok := spellCasterFromData(ctx)
	if !ok || casterID < 0 || casterID >= len(ctx.Engine.State.Players) {
		return nil
	}
	candidates := ctx.Engine.dragonSnowfieldFreezeCandidates(casterID)
	if len(candidates) == 0 {
		return nil
	}
	sourceID := ctx.Source.InstanceID
	ownerID := ctx.PlayerID
	ctx.Engine.SetPendingActionWithError(casterID, "dragon_snowfield_freeze",
		"龙吟雪域:选择法力范围内1个单位冻结1", candidates, 1, 1, nil, false,
		func(selected []string, _ map[string]any) error {
			target := ctx.Engine.findUnitByInstanceID(firstSelected(selected))
			if target == nil || target.OwnerID != 1-casterID || target.Position == nil || !ctx.Engine.IsInSpellRange(casterID, target.Position.Col, target.Position.Row, false) {
				return fmt.Errorf("invalid dragon snowfield freeze target")
			}
			target.Statuses[StatusFreeze]++
			source := ctx.Engine.findSkill(ctx.Engine.State.Players[ownerID], sourceID)
			if source == nil {
				return nil
			}
			source.Statuses[dragonSnowfieldTriggerStatus]++
			ctx.Engine.promptDragonSnowfieldSummon(ownerID, source)
			return nil
		})
	return nil
}

func (e *Engine) dragonSnowfieldFreezeCandidates(casterID int) []map[string]any {
	if e == nil || casterID < 0 || casterID >= len(e.State.Players) {
		return nil
	}
	candidates := []map[string]any{}
	ownerID := 1 - casterID
	ps := e.State.Players[ownerID]
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := ps.Units[col][row]
			if unit == nil || unit.CurrentLife <= 0 || !e.IsInSpellRange(casterID, col, row, false) {
				continue
			}
			info := cardToInfo(unit)
			info["side"] = "enemy"
			info["zone"] = "unit"
			info["position"] = Position{Col: col, Row: row}
			candidates = append(candidates, info)
		}
	}
	return candidates
}

func (e *Engine) promptDragonSnowfieldSummon(ownerID int, source *CardInstance) {
	if e == nil || source == nil || source.Statuses[dragonSnowfieldTriggerStatus] < 5 || ownerID < 0 || ownerID >= len(e.State.Players) {
		return
	}
	positions := e.friendlyEmptyUnitPositions(ownerID)
	if len(positions) == 0 {
		return
	}
	sourceID := source.InstanceID
	e.SetPendingActionWithError(ownerID, "dragon_snowfield_summon_frost_dragon",
		"龙吟雪域:是否召唤凛冰之龙", positions, 0, 1, nil, false,
		func(selected []string, _ map[string]any) error {
			if len(selected) == 0 {
				return nil
			}
			pos, ok := positionFromSelectionID(firstSelected(selected))
			if !ok || !pos.Valid() || e.State.Players[ownerID].Units[pos.Col][pos.Row] != nil {
				return fmt.Errorf("invalid dragon snowfield summon position")
			}
			source := e.findSkill(e.State.Players[ownerID], sourceID)
			if source == nil || source.Statuses[dragonSnowfieldTriggerStatus] < 5 {
				return nil
			}
			if e.summonFreshCardAtPosition(ownerID, "1201101", pos, true) == nil {
				return fmt.Errorf("failed to summon frost dragon")
			}
			source.Statuses[dragonSnowfieldTriggerStatus] -= 5
			return nil
		})
}

type Card3121105Embers struct{ AlwaysActive }

func (Card3121105Embers) ID() string   { return "3121105" }
func (Card3121105Embers) Name() string { return "余火" }
func (Card3121105Embers) OnTurnEnd(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.PlayerID < 0 || ctx.PlayerID >= len(ctx.Engine.State.Players) {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ps.Elements[model.ElementFire] != 0 || ctx.Engine.findSkill(ps, ctx.Source.InstanceID) != ctx.Source {
		return nil
	}
	if err := ctx.Engine.validateSkillForPurpose(ctx.Source, skillPurposeAttack); err != nil {
		return nil
	}
	candidates := ctx.Engine.spellTargetCandidates(ctx.PlayerID, ctx.Source)
	if len(candidates) == 0 {
		return nil
	}
	sourceID := ctx.Source.InstanceID
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "embers_free_cast_target",
		"余火:选择目标免费使用此卡", candidates, 1, 1, nil, false,
		func(selected []string, _ map[string]any) error {
			source := ctx.Engine.findSkill(ctx.Engine.State.Players[ctx.PlayerID], sourceID)
			if source == nil {
				return nil
			}
			target := selectedSpellTargetFromCandidates(ctx.Engine, ctx.PlayerID, source, firstSelected(selected), candidates)
			if target == nil {
				return fmt.Errorf("invalid embers target")
			}
			return ctx.Engine.startFreeSpellCastNoBoost(ctx.PlayerID, source, *target, map[string]any{"triggered_by": "3121105"})
		})
	return nil
}

func selectedSpellTargetFromCandidates(e *Engine, playerID int, skill *CardInstance, instanceID string, candidates []map[string]any) *SpellTarget {
	if e == nil || skill == nil || instanceID == "" {
		return nil
	}
	for _, candidate := range candidates {
		id, _ := candidate["instance_id"].(string)
		if id != instanceID {
			continue
		}
		target := SpellTarget{Type: "unit"}
		if pos, ok := candidate["position"].(Position); ok {
			target.Position = pos
		} else {
			unit := e.findUnitByInstanceID(instanceID)
			if unit == nil || unit.Position == nil {
				return nil
			}
			target.Position = *unit.Position
		}
		if owner, ok := candidate["target_owner"].(int); ok {
			target.OwnerID = &owner
		}
		if err := e.validateSpellTarget(playerID, skill, target); err != nil {
			return nil
		}
		return &target
	}
	return nil
}

type Card3111101FlameInferno struct{ AlwaysActive }

func (Card3111101FlameInferno) ID() string   { return "3111101" }
func (Card3111101FlameInferno) Name() string { return "火炎炼狱" }
func (Card3111101FlameInferno) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.ExtraData == nil {
		return nil
	}
	power := intFromData(ctx.ExtraData, "power", 0)
	burn := max(power, 0) / 8
	if burn <= 0 {
		return nil
	}
	for _, unit := range spellHitAffectedUnitsFromData(ctx) {
		if unit == nil || unit.CurrentLife <= 0 {
			continue
		}
		unit.Statuses[StatusBurn] += burn
	}
	return nil
}

type Card3501101FiveRainbowBeam struct{ AlwaysActive }

func (Card3501101FiveRainbowBeam) ID() string   { return "3501101" }
func (Card3501101FiveRainbowBeam) Name() string { return "五虹之束" }
func (Card3501101FiveRainbowBeam) HasActivePierce(card *CardInstance) bool {
	return card != nil && card.Statuses[fiveRainbowBeamSelectedStatus(model.ElementAir)] > 0
}
func (Card3501101FiveRainbowBeam) HasPierce() bool { return true }
func (Card3501101FiveRainbowBeam) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Source == nil || stats == nil {
		return
	}
	fire := ctx.Source.Statuses[fiveRainbowBeamSelectedStatus(model.ElementFire)]
	earth := ctx.Source.Statuses[fiveRainbowBeamSelectedStatus(model.ElementEarth)]
	stats.DamageBonus += fire * 2
	stats.PowerBonus += earth * 3
	if ctx.Source.Statuses[fiveRainbowBeamAllStatus] > 0 {
		stats.PowerBonus *= 2
	}
}

func fiveRainbowBeamSelectedStatus(elem string) string {
	return "五虹之束消耗:" + elem
}

const fiveRainbowBeamAllStatus = "五虹之束五色齐发"

type Card3311101SkyPhantasm struct{ AlwaysActive }

func (Card3311101SkyPhantasm) ID() string   { return "3311101" }
func (Card3311101SkyPhantasm) Name() string { return "苍穹幻韵" }

func (e *Engine) promptSkyPhantasmSpellChoice(playerID int, source *CardInstance) error {
	if e == nil || source == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	candidates := e.friendlySkills(playerID, func(skill *CardInstance) bool {
		return skill != nil && skill.Card != nil &&
			skill.InstanceID != source.InstanceID &&
			canUseSkillForPurpose(skill.Card, skillPurposeAttack) &&
			(hasCardTag(skill.Card, "驱动") || hasCardTag(skill.Card, "聚能"))
	})
	if len(candidates) == 0 {
		return nil
	}
	sourceID := source.InstanceID
	e.SetPendingActionWithError(playerID, "sky_phantasm_spell_choice",
		"苍穹幻韵:选择另一个驱动或聚能法术", candidates, 1, 1, nil, false,
		func(selected []string, _ map[string]any) error {
			if e.findSkill(e.State.Players[playerID], sourceID) == nil {
				return nil
			}
			baseSkill := e.findSkill(e.State.Players[playerID], firstSelected(selected))
			if baseSkill == nil || baseSkill.Card == nil || baseSkill.InstanceID == sourceID ||
				!canUseSkillForPurpose(baseSkill.Card, skillPurposeAttack) ||
				(!hasCardTag(baseSkill.Card, "驱动") && !hasCardTag(baseSkill.Card, "聚能")) {
				return fmt.Errorf("invalid sky phantasm spell")
			}
			virtual := cloneVirtualSpell(baseSkill, playerID, e.State.TurnNumber)
			targets := e.spellTargetCandidates(playerID, virtual)
			if skillNeedsTargetInstance(virtual) {
				if len(targets) == 0 {
					return nil
				}
				e.SetPendingActionWithError(playerID, "sky_phantasm_target",
					fmt.Sprintf("苍穹幻韵:选择%s的目标", virtual.Card.Name), targets, 1, 1, nil, false,
					func(targetSelected []string, _ map[string]any) error {
						target := selectedSpellTargetFromCandidates(e, playerID, virtual, firstSelected(targetSelected), targets)
						if target == nil {
							return fmt.Errorf("invalid sky phantasm target")
						}
						return e.startVirtualSpellCastNoBoost(playerID, virtual, *target, map[string]any{
							"triggered_by":        "3311101",
							"source_skill":        cardToInfo(baseSkill),
							"source_skill_hidden": false,
						})
					})
				return nil
			}
			return e.startVirtualSpellCastNoBoost(playerID, virtual, SpellTarget{Type: "none"}, map[string]any{
				"triggered_by":        "3311101",
				"source_skill":        cardToInfo(baseSkill),
				"source_skill_hidden": false,
			})
		})
	return nil
}

func cloneVirtualSpell(source *CardInstance, ownerID int, turn int) *CardInstance {
	virtual := NewCardInstance(source.Card, ownerID, turn)
	virtual.AttackBonus = source.AttackBonus
	virtual.PowerBonus = source.PowerBonus
	virtual.ElementsGainBonus = copyElements(source.ElementsGainBonus)
	virtual.Statuses = copyStatuses(source.Statuses)
	virtual.BoundSkills = append([]*CardInstance{}, source.BoundSkills...)
	return virtual
}

func copyElements(values map[string]int) map[string]int {
	if len(values) == 0 {
		return nil
	}
	copied := make(map[string]int, len(values))
	for key, value := range values {
		copied[key] = value
	}
	return copied
}

func copyStatuses(values map[string]int) map[string]int {
	if len(values) == 0 {
		return map[string]int{}
	}
	copied := make(map[string]int, len(values))
	for key, value := range values {
		copied[key] = value
	}
	return copied
}

func (e *Engine) startVirtualSpellCastNoBoost(playerID int, skill *CardInstance, target SpellTarget, extraData map[string]any) error {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) || skill == nil || skill.Card == nil {
		return fmt.Errorf("invalid virtual spell cast")
	}
	if !canUseSkillForPurpose(skill.Card, skillPurposeAttack) {
		return fmt.Errorf("virtual spell cannot attack")
	}
	if skillNeedsTargetInstance(skill) {
		if err := e.validateSpellTargetWithPierce(playerID, skill, target, e.skillHasPierce(playerID, skill)); err != nil {
			return err
		}
	}
	e.applyCoralBellyFirstSpellAttackBonus(playerID, skill)
	totalPower := e.effectiveSpellPower(playerID, skill, nil, target)
	powerSources := e.spellPowerSources(playerID, skill, nil, totalPower, target)
	e.consumeNextSpellPowerBonuses(e.State.Players[playerID], skill)

	isSorcery := isSorcerySkill(skill.Card)
	e.recordSpellCast(playerID, skill)
	e.triggerMagicMothAfterFocusSpellCast(playerID, skill)
	spellCastData := map[string]any{
		"cast_player":  playerID,
		"attacker":     playerID,
		"skill":        cardToInfo(skill),
		"target":       target,
		"power":        totalPower,
		"boost_count":  0,
		"is_sorcery":   isSorcery,
		"virtual_cast": true,
	}
	for key, value := range extraData {
		spellCastData[key] = value
	}
	e.emit(GameEvent{Type: "spell_cast", Player: -1, Data: spellCastData})
	e.triggerEffects(TriggerOnSpellCast, skill, nil, spellCastData)

	if isSorcery {
		resolveSorcery := func() {
			if e.shouldResolveSorceryHit(skill) {
				e.resolveSpellHit(playerID, skill, target, nil, nil)
			}
		}
		if e.triggerSpellCastFieldEffectsWithContinuation(playerID, skill, spellCastData, resolveSorcery) {
			return nil
		}
		resolveSorcery()
		return nil
	}

	e.State.PendingSpell = &SpellCast{
		AttackerID:   playerID,
		Skill:        skill,
		Target:       target,
		TotalPower:   totalPower,
		PowerSources: powerSources,
	}
	resolveWithoutDefense := func() {
		e.resolvePendingSpellHit()
	}
	openDefenseWindow := func() {
		if e.State.PendingSpell == nil {
			return
		}
		e.State.ResumePhase = PhaseDefenseWindow
		e.State.Phase = PhaseDefenseWindow
		e.emit(GameEvent{Type: "defense_window", Player: 1 - playerID, Data: map[string]any{"timeout": 30}})
	}
	continueSpell := openDefenseWindow
	if !e.spellAllowsDefense(playerID, skill, target) {
		continueSpell = resolveWithoutDefense
	}
	if e.triggerSpellCastFieldEffectsWithContinuation(playerID, skill, spellCastData, continueSpell) {
		if e.spellAllowsDefense(playerID, skill, target) {
			e.State.ResumePhase = PhaseDefenseWindow
		}
		return nil
	}
	continueSpell()
	return nil
}

func (e *Engine) prepareFiveRainbowBeamMarkers(playerID int, skill *CardInstance, action ActionMessage) (map[string]int, *CardInstance, error) {
	markers := map[string]int{}
	if skill == nil || skill.Card == nil || skill.Card.Number != "3501101" {
		return markers, nil, nil
	}
	for _, elem := range fiveRainbowElements() {
		markers[elem] = 0
	}
	raw, ok := action.Data["rainbow_markers"]
	if !ok {
		return markers, nil, nil
	}
	values, ok := raw.(map[string]any)
	if !ok {
		return nil, nil, fmt.Errorf("invalid rainbow markers")
	}
	ring := e.findFiveRainbowRingForSkill(playerID, skill)
	if ring == nil {
		return nil, nil, fmt.Errorf("five rainbow beam requires five rainbow ring")
	}
	for _, elem := range fiveRainbowElements() {
		amount, ok := intFromAny(values[elem])
		if !ok {
			amount = 0
		}
		if amount < 0 {
			return nil, nil, fmt.Errorf("invalid rainbow marker amount")
		}
		if amount > ring.Statuses[fiveRainbowMarkerStatus(elem)] {
			return nil, nil, fmt.Errorf("not enough rainbow markers")
		}
		markers[elem] = amount
	}
	return markers, ring, nil
}

func (e *Engine) applyFiveRainbowBeamMarkers(skill *CardInstance, ring *CardInstance, markers map[string]int) {
	if skill == nil || skill.Card == nil || skill.Card.Number != "3501101" || len(markers) == 0 {
		return
	}
	clearFiveRainbowBeamSelection(skill)
	fullSet := true
	for _, elem := range fiveRainbowElements() {
		amount := markers[elem]
		if amount <= 0 {
			fullSet = false
			continue
		}
		skill.Statuses[fiveRainbowBeamSelectedStatus(elem)] = amount
		if ring != nil {
			ring.Statuses[fiveRainbowMarkerStatus(elem)] -= amount
			if ring.Statuses[fiveRainbowMarkerStatus(elem)] <= 0 {
				delete(ring.Statuses, fiveRainbowMarkerStatus(elem))
			}
		}
	}
	if fullSet {
		skill.Statuses[fiveRainbowBeamAllStatus] = 1
	}
}

func clearFiveRainbowBeamSelection(skill *CardInstance) {
	if skill == nil {
		return
	}
	for _, elem := range fiveRainbowElements() {
		delete(skill.Statuses, fiveRainbowBeamSelectedStatus(elem))
	}
	delete(skill.Statuses, fiveRainbowBeamAllStatus)
}

func (e *Engine) findFiveRainbowRingForSkill(playerID int, skill *CardInstance) *CardInstance {
	if e == nil || skill == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	ps := e.State.Players[playerID]
	for _, card := range e.getAllFieldCards(ps) {
		if card == nil || card.Card == nil || card.Card.Number != "2511102" {
			continue
		}
		for _, bound := range card.BoundSkills {
			if bound == skill {
				return card
			}
		}
	}
	return nil
}

func fiveRainbowElements() []string {
	return []string{model.ElementFire, model.ElementWater, model.ElementEarth, model.ElementAir, model.ElementLight}
}

func (e *Engine) fiveRainbowBeamExtraTargetsFromAction(playerID int, skill *CardInstance, mainTarget SpellTarget, action ActionMessage, maxTargets int, hasPierce bool) ([]SpellTarget, error) {
	if maxTargets <= 0 {
		return nil, nil
	}
	raw, ok := action.Data["extra_targets"].([]any)
	if !ok || len(raw) == 0 {
		return nil, nil
	}
	if len(raw) > maxTargets {
		return nil, fmt.Errorf("too many five rainbow beam extra targets")
	}
	result := make([]SpellTarget, 0, len(raw))
	for _, value := range raw {
		data, ok := value.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid five rainbow beam extra target")
		}
		colF, colOK := data["col"].(float64)
		rowF, rowOK := data["row"].(float64)
		if !colOK || !rowOK {
			return nil, fmt.Errorf("invalid five rainbow beam extra target")
		}
		extra := SpellTarget{Type: "unit", Position: Position{Col: int(colF), Row: int(rowF)}}
		if ownerF, ok := data["owner"].(float64); ok {
			owner := int(ownerF)
			extra.OwnerID = &owner
		}
		if extra.Position == mainTarget.Position && !e.allowsSameSpellExtraTarget(e.State.Players[playerID], skill) {
			continue
		}
		if err := e.validateSpellTargetWithPierce(playerID, skill, extra, hasPierce); err != nil {
			return nil, err
		}
		result = append(result, extra)
	}
	return result, nil
}

type Card2221104WaterMirrorScroll struct{ AlwaysActive }

func (Card2221104WaterMirrorScroll) ID() string   { return "2221104" }
func (Card2221104WaterMirrorScroll) Name() string { return "水镜卷轴" }
func (Card2221104WaterMirrorScroll) OnUseItem(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.PlayerID < 0 || ctx.PlayerID >= len(ctx.Engine.State.Players) {
		return nil
	}
	recorded := ctx.Engine.State.Players[ctx.PlayerID].LastLowCostWaterSpell
	if recorded == nil || recorded.Card == nil {
		return nil
	}
	virtual := cloneVirtualSpell(recorded, ctx.PlayerID, ctx.Engine.State.TurnNumber)
	targets := ctx.Engine.spellTargetCandidates(ctx.PlayerID, virtual)
	if skillNeedsTargetInstance(virtual) {
		if len(targets) == 0 {
			return nil
		}
		ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "water_mirror_scroll_target",
			fmt.Sprintf("水镜卷轴:选择%s的目标", virtual.Card.Name), targets, 1, 1, nil, false,
			func(selected []string, _ map[string]any) error {
				target := selectedSpellTargetFromCandidates(ctx.Engine, ctx.PlayerID, virtual, firstSelected(selected), targets)
				if target == nil {
					return fmt.Errorf("invalid water mirror scroll target")
				}
				return ctx.Engine.startVirtualSpellCastNoBoost(ctx.PlayerID, virtual, *target, map[string]any{
					"triggered_by":        "2221104",
					"source_item":         cardToInfo(ctx.Source),
					"source_skill_hidden": false,
				})
			})
		return nil
	}
	return ctx.Engine.startVirtualSpellCastNoBoost(ctx.PlayerID, virtual, SpellTarget{Type: "none"}, map[string]any{
		"triggered_by":        "2221104",
		"source_item":         cardToInfo(ctx.Source),
		"source_skill_hidden": false,
	})
}

func (e *Engine) promptCounterWindHoleScroll(counter *CardInstance, original *CardInstance, extraData map[string]any) {
	if e == nil || counter == nil || original == nil || counter.OwnerID < 0 || counter.OwnerID >= len(e.State.Players) {
		return
	}
	ownerID := counter.OwnerID
	candidates := e.enemyUnits(ownerID, false, func(unit *CardInstance) bool {
		return unit != nil && unit.Card != nil && unit.Card.IsCompanion()
	})
	if len(candidates) == 0 {
		return
	}
	counterID := counter.InstanceID
	e.SetPendingActionWithError(ownerID, "counter_wind_hole_scroll_target",
		"反击风洞卷轴:选择敌方单位反击该法术", candidates, 1, 1, nil, false,
		func(selected []string, _ map[string]any) error {
			targetUnit := selectedUnitFromCandidates(e, selected, candidates)
			if targetUnit == nil || targetUnit.Position == nil {
				return fmt.Errorf("invalid counter wind hole target")
			}
			virtual := cloneVirtualSpell(original, ownerID, e.State.TurnNumber)
			boosts := cloneSpellInstances(spellInstancesFromData(extraData, "boost_skills"), ownerID, e.State.TurnNumber)
			target := SpellTarget{Type: "unit", Position: *targetUnit.Position}
			if targetUnit.OwnerID == ownerID {
				target.OwnerID = &ownerID
			}
			if e.State.PendingSpell != nil {
				clearFiveRainbowBeamSelection(e.State.PendingSpell.Skill)
				e.State.PendingSpell = nil
			}
			return e.startVirtualSpellCastWithBoosts(ownerID, virtual, target, boosts, map[string]any{
				"triggered_by": "2321111",
				"source_item":  counterID,
			})
		})
}

func cloneSpellInstances(skills []*CardInstance, ownerID int, turn int) []*CardInstance {
	clones := make([]*CardInstance, 0, len(skills))
	for _, skill := range skills {
		if skill == nil || skill.Card == nil {
			continue
		}
		clones = append(clones, cloneVirtualSpell(skill, ownerID, turn))
	}
	return clones
}

func (e *Engine) startVirtualSpellCastWithBoosts(playerID int, skill *CardInstance, target SpellTarget, boostSkills []*CardInstance, extraData map[string]any) error {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) || skill == nil || skill.Card == nil {
		return fmt.Errorf("invalid virtual boosted spell cast")
	}
	if !canUseSkillForPurpose(skill.Card, skillPurposeAttack) {
		return fmt.Errorf("virtual spell cannot attack")
	}
	if target.Type == "unit" && !target.Position.Valid() {
		return fmt.Errorf("invalid virtual spell target")
	}
	e.applyCoralBellyFirstSpellAttackBonus(playerID, skill)
	totalPower := e.effectiveSpellPower(playerID, skill, boostSkills, target)
	powerSources := e.spellPowerSources(playerID, skill, boostSkills, totalPower, target)
	isSorcery := isSorcerySkill(skill.Card)
	e.recordSpellCast(playerID, skill)
	e.triggerMagicMothAfterFocusSpellCast(playerID, skill)
	spellCastData := map[string]any{
		"cast_player":  playerID,
		"attacker":     playerID,
		"skill":        cardToInfo(skill),
		"target":       target,
		"power":        totalPower,
		"boost_count":  len(boostSkills),
		"boost_skills": boostSkills,
		"is_sorcery":   isSorcery,
		"virtual_cast": true,
	}
	for key, value := range extraData {
		spellCastData[key] = value
	}
	e.emit(GameEvent{Type: "spell_cast", Player: -1, Data: spellCastData})
	e.triggerEffects(TriggerOnSpellCast, skill, nil, spellCastData)
	if isSorcery {
		resolveSorcery := func() {
			if e.shouldResolveSorceryHit(skill) {
				e.resolveSpellHit(playerID, skill, target, boostSkills, nil)
			}
		}
		if e.triggerSpellCastFieldEffectsWithContinuation(playerID, skill, spellCastData, resolveSorcery) {
			return nil
		}
		resolveSorcery()
		return nil
	}
	e.State.PendingSpell = &SpellCast{
		AttackerID:   playerID,
		Skill:        skill,
		Target:       target,
		TotalPower:   totalPower,
		PowerSources: powerSources,
		BoostSkills:  boostSkills,
	}
	resolveWithoutDefense := func() {
		e.resolvePendingSpellHit()
	}
	openDefenseWindow := func() {
		if e.State.PendingSpell == nil {
			return
		}
		e.State.ResumePhase = PhaseDefenseWindow
		e.State.Phase = PhaseDefenseWindow
		e.emit(GameEvent{Type: "defense_window", Player: 1 - playerID, Data: map[string]any{"timeout": 30}})
	}
	continueSpell := openDefenseWindow
	if !e.spellAllowsDefense(playerID, skill, target) {
		continueSpell = resolveWithoutDefense
	}
	continueAfterMainCounters := func() {
		if e.promptAttackBoostSpellCastCounters(playerID, boostSkills, continueSpell) {
			if e.spellAllowsDefense(playerID, skill, target) {
				e.State.ResumePhase = PhaseDefenseWindow
			}
			return
		}
		continueSpell()
	}
	if e.triggerSpellCastFieldEffectsWithContinuation(playerID, skill, spellCastData, continueAfterMainCounters) {
		if e.spellAllowsDefense(playerID, skill, target) {
			e.State.ResumePhase = PhaseDefenseWindow
		}
		return nil
	}
	continueAfterMainCounters()
	return nil
}

const endlessWindTideAirCostBonusStatus = "endless_wind_tide_air_cost_bonus"

type Card2321106EndlessWindTide struct{ AlwaysActive }

func (Card2321106EndlessWindTide) ID() string   { return "2321106" }
func (Card2321106EndlessWindTide) Name() string { return "无尽风潮" }

func (Card2321106EndlessWindTide) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || !isFriendlySpellHit(ctx) {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	removed := false
	for i, skill := range ps.Skills {
		if skill == ctx.Source {
			ps.Skills[i] = nil
			removed = true
			break
		}
	}
	if !removed {
		removed = ctx.Engine.removeCardFromGraveyard(ctx.PlayerID, ctx.Source)
	}
	if !removed {
		return nil
	}
	ctx.Source.SlotIndex = -1
	ctx.Source.IsHorizontal = true
	ctx.Source.PowerBonus += 2
	ctx.Source.Statuses[endlessWindTideAirCostBonusStatus]++
	ctx.Source.Statuses[StatusCannotUseSkillUntilTurn] = ctx.Engine.State.TurnNumber
	ps.Hand = append(ps.Hand, ctx.Source)
	ctx.Engine.emit(GameEvent{
		Type:   "endless_wind_tide_return",
		Player: -1,
		Data: map[string]any{
			"player": ctx.PlayerID,
			"card":   cardToInfo(ctx.Source),
		},
	})
	return nil
}

func (Card2321106EndlessWindTide) ModifySelfCardPlayCost(ctx *EffectContext, cost map[string]int) {
	if ctx == nil || ctx.Source == nil {
		return
	}
	cost[model.ElementAir] += ctx.Source.Statuses[endlessWindTideAirCostBonusStatus]
}

type Card2221103IceLockRune struct{ AlwaysActive }

func (Card2221103IceLockRune) ID() string   { return "2221103" }
func (Card2221103IceLockRune) Name() string { return "冰锁符文" }

func (Card2221103IceLockRune) OnCardEnter(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Target == nil || ctx.Target.Card == nil || ctx.ExtraData == nil {
		return nil
	}
	learned, _ := ctx.ExtraData["learned_skill"].(bool)
	enteredPlayer, _ := ctx.ExtraData["entered_player"].(int)
	if !learned || enteredPlayer == ctx.PlayerID || !ctx.Target.Card.IsSkill() {
		return nil
	}
	ctx.Target.Statuses[StatusCannotUseSkillUntilTurn] = ctx.Engine.State.TurnNumber + 1
	ctx.Engine.emit(GameEvent{
		Type:   "skill_locked",
		Player: -1,
		Data: map[string]any{
			"player": ctx.Target.OwnerID,
			"source": cardToInfo(ctx.Source),
			"target": cardToInfo(ctx.Target),
			"until":  ctx.Target.Statuses[StatusCannotUseSkillUntilTurn],
		},
	})
	return nil
}

type Card2021114GuardianRune struct{ AlwaysActive }

func (Card2021114GuardianRune) ID() string   { return "2021114" }
func (Card2021114GuardianRune) Name() string { return "神护符文" }

func (Card2021114GuardianRune) OnDamaged(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.ExtraData == nil {
		return nil
	}
	prevented, _ := ctx.ExtraData["prevent_damage"].(*bool)
	if prevented == nil || !lethalDamageFromData(ctx.ExtraData, ctx.Target) {
		return nil
	}
	*prevented = true
	ctx.Engine.emit(GameEvent{
		Type:   "damage_prevented",
		Player: -1,
		Data: map[string]any{
			"source": cardToInfo(ctx.Source),
			"target": cardToInfo(ctx.Target),
			"amount": damageFromData(ctx.ExtraData),
			"reason": "guardian_rune",
		},
	})
	return nil
}

type Card2221111IceSoulSealForge struct{ AlwaysActive }

func (Card2221111IceSoulSealForge) ID() string   { return "2221111" }
func (Card2221111IceSoulSealForge) Name() string { return "冰魄印 淬" }

func (Card2221111IceSoulSealForge) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Target == nil || ctx.Target.OwnerID == ctx.PlayerID {
		return nil
	}
	original := spellPowerFromData(ctx.ExtraData)
	if original <= 10 {
		return nil
	}
	reduced := (original + 1) / 2
	if ctx.Engine.State.PendingSpell != nil && ctx.Engine.State.PendingSpell.Skill == ctx.Target {
		ctx.Engine.State.PendingSpell.TotalPower = reduced
	}
	if ctx.ExtraData != nil {
		ctx.ExtraData["power"] = reduced
	}
	ctx.Engine.emit(GameEvent{
		Type:   "spell_power_reduced",
		Player: -1,
		Data: map[string]any{
			"player":   ctx.PlayerID,
			"source":   cardToInfo(ctx.Source),
			"spell":    cardToInfo(ctx.Target),
			"original": original,
			"power":    reduced,
		},
	})
	return nil
}

type Card2521109PunishmentRune struct{ AlwaysActive }

func (Card2521109PunishmentRune) ID() string   { return "2521109" }
func (Card2521109PunishmentRune) Name() string { return "惩戒符文" }

func (Card2521109PunishmentRune) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.OwnerID == ctx.PlayerID {
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.Target.OwnerID, false, func(unit *CardInstance) bool {
		return unit != nil && unit.Card != nil && unit.Card.IsCompanion()
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "punishment_rune_damage",
		"惩戒符文:选择1个敌方伙伴造成2点伤害", candidates, 1, 1,
		nil, false, func(selected []string, _ map[string]any) error {
			target := selectedUnitFromCandidates(ctx.Engine, selected, candidates)
			if target == nil || target.Card == nil || target.OwnerID != ctx.Target.OwnerID || !target.Card.IsCompanion() {
				return fmt.Errorf("invalid punishment rune target")
			}
			ctx.Engine.dealDamageWithExtra(target, 2, ctx.PlayerID, map[string]any{
				"damage_source":  "punishment_rune",
				"damage_element": model.ElementLight,
				"source_card":    ctx.Source,
			})
			ctx.Engine.emit(GameEvent{
				Type:   "punishment_rune_damage",
				Player: -1,
				Data: map[string]any{
					"player": ctx.PlayerID,
					"source": cardToInfo(ctx.Source),
					"target": cardToInfo(target),
					"damage": 2,
				},
			})
			return nil
		})
	return nil
}

const bloodRoseSealExpireTurnStatus = "blood_rose_seal_expire_turn"

type Card3621104BloodRoseSeal struct{ AlwaysActive }

func (Card3621104BloodRoseSeal) ID() string   { return "3621104" }
func (Card3621104BloodRoseSeal) Name() string { return "血蔷薇咒印" }

func (Card3621104BloodRoseSeal) OnEnter(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, false, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "blood_rose_seal_mark",
		"血蔷薇咒印:选择1个敌方单位标记", candidates, 1, 1,
		nil, false, func(selected []string, _ map[string]any) error {
			target := selectedUnitFromCandidates(ctx.Engine, selected, candidates)
			if target == nil || target.OwnerID != ctx.OpponentID {
				return fmt.Errorf("invalid blood rose seal target")
			}
			target.Statuses[bloodRoseSealMarkerStatus(ctx.Source)] = 1
			ctx.Source.Statuses[bloodRoseSealExpireTurnStatus] = ctx.Engine.State.TurnNumber + 2
			ctx.Engine.emit(GameEvent{
				Type:   "blood_rose_seal_mark",
				Player: -1,
				Data: map[string]any{
					"player": ctx.PlayerID,
					"source": cardToInfo(ctx.Source),
					"target": cardToInfo(target),
				},
			})
			return nil
		})
	return nil
}

func (Card3621104BloodRoseSeal) OnEnemyDeath(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Source.Card == nil {
		return nil
	}
	if ctx.Source.Statuses[bloodRoseSealExpireTurnStatus] > 0 && ctx.Engine.State.TurnNumber > ctx.Source.Statuses[bloodRoseSealExpireTurnStatus] {
		return nil
	}
	if ctx.Target.Statuses[bloodRoseSealMarkerStatus(ctx.Source)] <= 0 {
		return nil
	}
	if ctx.Engine.findSkill(ctx.Engine.State.Players[ctx.PlayerID], ctx.Source.InstanceID) != ctx.Source {
		return nil
	}
	if !ctx.Engine.bindSkillToHero(ctx.PlayerID, ctx.Source) {
		return nil
	}
	ctx.Source.Statuses["使用费用"+model.ElementShadow+"-1"]++
	delete(ctx.Source.Statuses, bloodRoseSealExpireTurnStatus)
	ctx.Engine.clearBloodRoseSealMarkers(ctx.Source)
	ctx.Engine.emit(GameEvent{
		Type:   "blood_rose_seal_bound",
		Player: -1,
		Data: map[string]any{
			"player": ctx.PlayerID,
			"source": cardToInfo(ctx.Source),
			"target": cardToInfo(ctx.Target),
		},
	})
	return nil
}

func (Card3621104BloodRoseSeal) OnTurnEnd(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	endedPlayer := ctx.PlayerID
	if ctx.ExtraData != nil {
		endedPlayer, _ = ctx.ExtraData["ended_player"].(int)
	}
	if endedPlayer != ctx.PlayerID {
		return nil
	}
	expires := ctx.Source.Statuses[bloodRoseSealExpireTurnStatus]
	if expires <= 0 || ctx.Engine.State.TurnNumber < expires {
		return nil
	}
	delete(ctx.Source.Statuses, bloodRoseSealExpireTurnStatus)
	ctx.Engine.clearBloodRoseSealMarkers(ctx.Source)
	return nil
}

type Card1611102BloodThornGarden struct{ AlwaysActive }

func (Card1611102BloodThornGarden) ID() string   { return "1611102" }
func (Card1611102BloodThornGarden) Name() string { return "蔷薇花园的血荆棘" }
func (Card1611102BloodThornGarden) OnDeath(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Position == nil || !bloodThornKilledByFriendlyCard(ctx) {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ps == nil || ps.Elements[model.ElementShadow] < 1 {
		return nil
	}
	pos := *ctx.Source.Position
	if !pos.Valid() || ps.Units[pos.Col][pos.Row] != nil {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "blood_thorn_resummon",
		"蔷薇花园的血荆棘:是否支付1暗重新召唤", []map[string]any{candidateInfo(ctx.Source, "graveyard", "own")}, 0, 1,
		func(selected []string) {
			if len(selected) == 0 || ps.Elements[model.ElementShadow] < 1 || !pos.Valid() || ps.Units[pos.Col][pos.Row] != nil {
				return
			}
			if !ctx.Engine.removeCardFromGraveyard(ctx.PlayerID, ctx.Source) {
				return
			}
			ps.Elements[model.ElementShadow]--
			resetCardForResummon(ctx.Source)
			if !ctx.Engine.placeExistingCompanionAtPosition(ctx.PlayerID, ctx.Source, pos, true) {
				ctx.Engine.addToGraveyard(ctx.PlayerID, ctx.Source)
				return
			}
			ctx.Engine.emit(GameEvent{
				Type:   "blood_thorn_resummon",
				Player: -1,
				Data: map[string]any{
					"player": ctx.PlayerID,
					"card":   cardToInfo(ctx.Source),
				},
			})
		})
	return nil
}

func bloodThornKilledByFriendlyCard(ctx *EffectContext) bool {
	if ctx == nil || ctx.ExtraData == nil {
		return false
	}
	if attacker, ok := ctx.ExtraData["attacker"].(int); ok && attacker == ctx.PlayerID {
		return true
	}
	if source, ok := ctx.ExtraData["source_card"].(*CardInstance); ok && source != nil && source.OwnerID == ctx.PlayerID {
		return true
	}
	return false
}

func resetCardForResummon(card *CardInstance) {
	if card == nil || card.Card == nil {
		return
	}
	card.CurrentLife = card.Card.Life
	card.CurrentAttack = card.Card.Attack
	card.DamageTakenThisTurn = 0
	card.IsHorizontal = true
	card.Position = nil
	card.Statuses = make(map[string]int)
	card.ElementsGainBonus = make(map[string]int)
	card.ElementsGainSet = nil
	card.PowerBonus = 0
	card.AttackBonus = 0
	card.UsedThisTurn = 0
	card.UltimateUsed = false
	card.BoundSkills = nil
	card.UnderCards = nil
	card.AttachedBehaviors = nil
}

const robertBlackPineMarkerStatus = "robert_black_pine_markers"

type Card1611103RobertBlackPine struct{ AlwaysActive }

func (Card1611103RobertBlackPine) ID() string   { return "1611103" }
func (Card1611103RobertBlackPine) Name() string { return "鲜血贵公子 罗伯特 黑松" }
func (Card1611103RobertBlackPine) OnDamaged(ctx *EffectContext) error {
	if ctx == nil || ctx.Source == nil || ctx.Target == nil || ctx.ExtraData == nil {
		return nil
	}
	if ctx.Target.OwnerID != ctx.PlayerID {
		return nil
	}
	if attacker, ok := ctx.ExtraData["attacker"].(int); !ok || attacker != ctx.PlayerID {
		return nil
	}
	if damage, _ := ctx.ExtraData["damage"].(int); damage <= 0 {
		return nil
	}
	ctx.Source.Statuses[robertBlackPineMarkerStatus]++
	return nil
}
func (Card1611103RobertBlackPine) OnFriendlyDeath(ctx *EffectContext) error {
	if ctx == nil || ctx.Source == nil || ctx.Target == nil || !bloodThornKilledByFriendlyCard(ctx) {
		return nil
	}
	ctx.Source.Statuses[robertBlackPineMarkerStatus] += 2
	return nil
}
func (Card1611103RobertBlackPine) HasActivePerTurn(card *CardInstance) bool {
	return card != nil && card.Statuses[robertBlackPineMarkerStatus] >= 3
}
func (Card1611103RobertBlackPine) PerTurnLabel(*CardInstance) string {
	return "移除标记"
}
func (Card1611103RobertBlackPine) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Statuses[robertBlackPineMarkerStatus] < 3 {
		return fmt.Errorf("罗伯特需要移除3个标记物")
	}
	choices := []map[string]any{
		{"instance_id": "life", "name": "+1血", "zone": "choice", "side": "own"},
		{"instance_id": "load", "name": "负载+1暗", "zone": "choice", "side": "own"},
		{"instance_id": "attack", "name": "+1攻", "zone": "choice", "side": "own"},
	}
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "robert_black_pine_reward",
		"鲜血贵公子 罗伯特 黑松:选择奖励", choices, 1, 1,
		nil, false, func(selected []string, _ map[string]any) error {
			if ctx.Source.Statuses[robertBlackPineMarkerStatus] < 3 || !ctx.Engine.cardStillOnField(ctx.Source) {
				return fmt.Errorf("invalid Robert reward")
			}
			ctx.Source.Statuses[robertBlackPineMarkerStatus] -= 3
			switch firstSelected(selected) {
			case "life":
				ctx.Source.CurrentLife++
			case "load":
				ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementShadow, 1, ctx.Source)
			case "attack":
				ctx.Source.AttackBonus++
			default:
				return fmt.Errorf("invalid Robert reward")
			}
			ctx.Engine.emit(GameEvent{
				Type:   "robert_black_pine_reward",
				Player: -1,
				Data: map[string]any{
					"player": ctx.PlayerID,
					"source": cardToInfo(ctx.Source),
					"choice": firstSelected(selected),
				},
			})
			return nil
		})
	return nil
}

type Card1411102WhisperElfKingSindariel struct{ AlwaysActive }

func (Card1411102WhisperElfKingSindariel) ID() string   { return "1411102" }
func (Card1411102WhisperElfKingSindariel) Name() string { return "谧语精灵王 辛达瑞尔" }
func (Card1411102WhisperElfKingSindariel) OnEnter(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.PlayerID < 0 || ctx.PlayerID >= len(ctx.Engine.State.Players) {
		return nil
	}
	opponentID := 1 - ctx.PlayerID
	opponent := ctx.Engine.State.Players[opponentID]
	if opponent == nil {
		return nil
	}
	maxTargets := 0
	if opponent.SpellHitsLastTurn >= 3 {
		maxTargets++
	}
	if opponent.SpellHitTargetsLastTurn >= 3 {
		maxTargets++
	}
	if opponent.SpellDamageLastTurn >= 3 {
		maxTargets++
	}
	if maxTargets == 0 {
		return nil
	}
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, true, nil)
	if len(candidates) == 0 {
		return nil
	}
	if maxTargets > len(candidates) {
		maxTargets = len(candidates)
	}
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "sindariel_entry_damage",
		"谧语精灵王 辛达瑞尔:选择敌人造成2点伤害", candidates, 0, maxTargets,
		nil, false, func(selected []string, _ map[string]any) error {
			allowed := make(map[string]bool, len(candidates))
			for _, candidate := range candidates {
				if id, _ := candidate["instance_id"].(string); id != "" {
					allowed[id] = true
				}
			}
			for _, id := range selected {
				if !allowed[id] {
					return fmt.Errorf("invalid Sindariel target")
				}
				target := ctx.Engine.findUnitByInstanceID(id)
				if target == nil || target.OwnerID != opponentID {
					continue
				}
				ctx.Engine.dealDamageWithExtra(target, 2, opponentID, map[string]any{
					"damage_source":  "effect",
					"damage_element": model.ElementEarth,
					"source_card":    ctx.Source,
					"attacker":       ctx.PlayerID,
				})
				if target.CurrentLife <= 0 && !target.Card.IsHero() && ctx.Engine.unitInOwnerGrid(target, opponentID) {
					ctx.Engine.destroyUnitWithData(target, opponentID, map[string]any{
						"death_cause": "sindariel",
						"source_card": ctx.Source,
						"attacker":    ctx.PlayerID,
					})
				}
			}
			return nil
		})
	return nil
}

func bloodRoseSealMarkerStatus(source *CardInstance) string {
	if source == nil {
		return "blood_rose_seal_mark:"
	}
	return "blood_rose_seal_mark:" + source.InstanceID
}

func (e *Engine) bindSkillToHero(playerID int, skill *CardInstance) bool {
	if e == nil || skill == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return false
	}
	ps := e.State.Players[playerID]
	if ps == nil || ps.Hero == nil {
		return false
	}
	for _, bound := range ps.Hero.BoundSkills {
		if bound == skill {
			return true
		}
	}
	for i := range ps.Skills {
		if ps.Skills[i] == skill {
			ps.Skills[i] = nil
			break
		}
	}
	skill.SlotIndex = -1
	markTransferredBoundSkill(skill)
	ps.Hero.BoundSkills = append(ps.Hero.BoundSkills, skill)
	return true
}

func (e *Engine) clearBloodRoseSealMarkers(source *CardInstance) {
	if e == nil || source == nil {
		return
	}
	marker := bloodRoseSealMarkerStatus(source)
	for _, ps := range e.State.Players {
		if ps == nil {
			continue
		}
		for col := 0; col < 3; col++ {
			for row := 0; row < 3; row++ {
				if unit := ps.Units[col][row]; unit != nil {
					delete(unit.Statuses, marker)
				}
			}
		}
	}
}

type Card2611102SpiritCandle struct{ AlwaysActive }

func (Card2611102SpiritCandle) ID() string   { return "2611102" }
func (Card2611102SpiritCandle) Name() string { return "渡灵之烛" }
func (Card2611102SpiritCandle) OnEquip(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil {
		return nil
	}
	ctx.Engine.enforceSlotCapacities(ctx.Engine.State.Players[ctx.PlayerID])
	return nil
}

var _ OnSpellHitBehavior = Card2321106EndlessWindTide{}
var _ SelfCardPlayCostModifier = Card2321106EndlessWindTide{}

type Card2421111DesertLeggings struct{ AlwaysActive }

func (Card2421111DesertLeggings) ID() string   { return "2421111" }
func (Card2421111DesertLeggings) Name() string { return "沙漠护腿" }

func (Card2421111DesertLeggings) ModifyFieldDamageAmount(ctx *EffectContext, amount int) int {
	if ctx == nil || ctx.Source == nil || ctx.Source.UltimateUsed || ctx.Target == nil || ctx.Target.Card == nil || amount < 2 {
		return amount
	}
	if !ctx.Target.Card.IsCompanion() {
		return amount
	}
	ctx.Source.UltimateUsed = true
	return max(amount-2, 0)
}

var _ FieldDamageAmountModifier = Card2421111DesertLeggings{}

type Card3111102PrimalDivineFlameLopsius struct{ AlwaysActive }

func (Card3111102PrimalDivineFlameLopsius) ID() string   { return "3111102" }
func (Card3111102PrimalDivineFlameLopsius) Name() string { return "原初神炎 洛普修斯" }
func (Card3111102PrimalDivineFlameLopsius) PerTurnLabel(*CardInstance) string {
	return "献祭火焰技能"
}
func (Card3111102PrimalDivineFlameLopsius) ValidateSkillUse(ctx *EffectContext, skill *CardInstance, purpose skillPurpose) error {
	if purpose == skillPurposeBoost || purpose == skillPurposeAttackBoost || purpose == skillPurposeDefenseBoost {
		return fmt.Errorf("原初神炎 洛普修斯不能用于强化")
	}
	return nil
}
func (Card3111102PrimalDivineFlameLopsius) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) {
		return nil
	}
	candidates := ctx.Engine.friendlySkillsIncludingBound(ctx.PlayerID, func(skill *CardInstance) bool {
		return isFireSpellInstance(skill) && skill != ctx.Source
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "primal_divine_flame_exile",
		"原初神炎 洛普修斯:选择1个火焰技能移出游戏", candidates, 1, 1,
		func(selected []string) {
			if ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) {
				return
			}
			target := findFriendlySkillIncludingBound(ctx.Engine, ctx.PlayerID, firstSelected(selected))
			if !isFireSpellInstance(target) || target == ctx.Source {
				return
			}
			if !ctx.Engine.exileCard(ctx.PlayerID, target) {
				return
			}
			ctx.Source.AttackBonus++
			ctx.Source.PowerBonus += 2
			ctx.Source.UsedThisTurn++
			ctx.Engine.emit(GameEvent{
				Type:   "primal_divine_flame_growth",
				Player: -1,
				Data: map[string]any{
					"player": ctx.PlayerID,
					"source": cardToInfo(ctx.Source),
					"exiled": firstSelected(selected),
				},
			})
		})
	return nil
}

var _ PerTurnAbility = Card3111102PrimalDivineFlameLopsius{}
var _ SkillUsePermissionModifier = Card3111102PrimalDivineFlameLopsius{}

const mindSeaMazeAnyRangeUntilStatus = "mind_sea_maze_any_range_until"

type Card3211101MindSeaMaze struct{ AlwaysActive }

func (Card3211101MindSeaMaze) ID() string   { return "3211101" }
func (Card3211101MindSeaMaze) Name() string { return "心海迷离" }
func (Card3211101MindSeaMaze) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || !isFriendlySpellCast(ctx) {
		return nil
	}
	ctx.Source.Statuses[mindSeaMazeAnyRangeUntilStatus] = ctx.Engine.State.TurnNumber
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModSkillPowerBonus,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		TargetInstanceID: ctx.Source.InstanceID,
		Amount:           1,
		ExpiresTurn:      ctx.Engine.State.TurnNumber + 1,
	})
	return nil
}
func (Card3211101MindSeaMaze) ModifySpellArea(ctx *EffectContext, area *SpellArea) {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target != ctx.Source || area == nil {
		return
	}
	if ctx.Source.Statuses[mindSeaMazeAnyRangeUntilStatus] >= ctx.Engine.State.TurnNumber {
		*area = SpellAreaAll
	}
}

var _ OnSpellCastBehavior = Card3211101MindSeaMaze{}
var _ SpellAreaModifier = Card3211101MindSeaMaze{}

const (
	treadingWaveTriggerTurnStatus  = "treading_wave_trigger_turn"
	treadingWaveTriggerCountStatus = "treading_wave_trigger_count"
)

type Card3221101TreadingWave struct{ AlwaysActive }

func (Card3221101TreadingWave) ID() string   { return "3221101" }
func (Card3221101TreadingWave) Name() string { return "踏浪术" }
func (Card3221101TreadingWave) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || !isFriendlySpellCast(ctx) {
		return nil
	}
	castSkill := spellCastCardForTrigger(ctx)
	if castSkill == nil || castSkill.Card == nil || castSkill.Card.Category != model.ElementWater {
		return nil
	}
	candidates := ctx.Engine.friendlySkillsIncludingBound(ctx.PlayerID, func(skill *CardInstance) bool {
		return skill != nil && skill != castSkill && skill.Card != nil &&
			skill.Card.Category == model.ElementWater && skill.IsHorizontal
	})
	if len(candidates) == 0 {
		return nil
	}
	resetSkill := func(selection string) {
		target := ctx.Engine.findFriendlyCardIncludingBound(ctx.PlayerID, selection)
		if target == nil || target == castSkill || target.Card == nil || target.Card.Category != model.ElementWater || !target.IsHorizontal {
			return
		}
		if ctx.Source.Statuses[treadingWaveTriggerTurnStatus] != ctx.Engine.State.TurnNumber {
			ctx.Source.Statuses[treadingWaveTriggerTurnStatus] = ctx.Engine.State.TurnNumber
			ctx.Source.Statuses[treadingWaveTriggerCountStatus] = 0
		}
		ctx.Source.Statuses[treadingWaveTriggerCountStatus]++
		bonus := ctx.Source.Statuses[treadingWaveTriggerCountStatus] + 1
		target.IsHorizontal = false
		target.Statuses[skillUseExtraCostStatus(model.ElementWater, bonus)]++
		ctx.Engine.emit(GameEvent{
			Type:   "treading_wave_reset",
			Player: -1,
			Data: map[string]any{
				"player":      ctx.PlayerID,
				"source":      cardToInfo(ctx.Source),
				"cast_skill":  cardToInfo(castSkill),
				"reset_skill": cardToInfo(target),
				"cost_bonus":  bonus,
			},
		})
	}
	if len(candidates) == 1 {
		id, _ := candidates[0]["instance_id"].(string)
		resetSkill(id)
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "treading_wave_reset",
		"踏浪术:选择另一个横置的水纹法术重置", candidates, 1, 1,
		func(selected []string) {
			resetSkill(firstSelected(selected))
		})
	return nil
}

func spellCastCardForTrigger(ctx *EffectContext) *CardInstance {
	if ctx == nil {
		return nil
	}
	if ctx.Target != nil {
		return ctx.Target
	}
	return ctx.Source
}

var _ OnSpellCastBehavior = Card3221101TreadingWave{}

func (e *Engine) validateSpellPowerSacrifice(playerID int, skill *CardInstance, action ActionMessage) (*CardInstance, int, error) {
	if e == nil || skill == nil || skill.Card == nil || skill.Card.Number != "3121104" {
		return nil, 0, nil
	}
	sacrificeID, _ := action.Data["sacrifice_id"].(string)
	if sacrificeID == "" {
		return nil, 0, nil
	}
	target, zone := e.findFriendlyCandidate(playerID, sacrificeID)
	if zone != "unit" || target == nil || target.Card == nil || !target.Card.IsCompanion() || target.Card.Category != model.ElementFire {
		return nil, 0, fmt.Errorf("3121104 requires a friendly fire companion sacrifice")
	}
	return target, totalElementCost(target.Card.ElementsCost), nil
}

func (e *Engine) validateSpellPowerSacrificeForSources(playerID int, sources []*CardInstance, action ActionMessage) (*CardInstance, *CardInstance, int, error) {
	for _, source := range sources {
		sacrifice, bonus, err := e.validateSpellPowerSacrifice(playerID, source, action)
		if err != nil {
			return nil, nil, 0, err
		}
		if sacrifice != nil && bonus > 0 {
			return sacrifice, source, bonus, nil
		}
	}
	return nil, nil, 0, nil
}

func (e *Engine) validateOracleGlorySupport(playerID int, scroll *CardInstance, action ActionMessage) (int, error) {
	if e == nil || scroll == nil || scroll.Card == nil || scroll.Card.Number != "2521111" {
		return 0, nil
	}
	supportID, _ := action.Data["support_id"].(string)
	if supportID == "" {
		return 0, fmt.Errorf("2521111 requires a friendly companion support")
	}
	target, zone := e.findFriendlyCandidate(playerID, supportID)
	if zone != "unit" || target == nil || target.Card == nil || !target.Card.IsCompanion() {
		return 0, fmt.Errorf("2521111 requires a friendly companion support")
	}
	bonus := max(target.CurrentLife, 0) + e.totalLoad(target)
	if bonus <= 5 {
		return 0, fmt.Errorf("2521111 support companion life plus load must be greater than 5")
	}
	return bonus, nil
}

func (e *Engine) validateFlameArrayScrollSacrifice(playerID int, scroll *CardInstance, action ActionMessage) (*CardInstance, int, error) {
	if e == nil || scroll == nil || scroll.Card == nil || scroll.Card.Number != "2121105" {
		return nil, 0, nil
	}
	sacrificeID, _ := action.Data["sacrifice_id"].(string)
	if sacrificeID == "" {
		return nil, 0, fmt.Errorf("2121105 requires a friendly fire companion sacrifice")
	}
	target, zone := e.findFriendlyCandidate(playerID, sacrificeID)
	if zone != "unit" || target == nil || target.Card == nil || !target.Card.IsCompanion() || target.Card.Category != model.ElementFire {
		return nil, 0, fmt.Errorf("2121105 requires a friendly fire companion sacrifice")
	}
	bonus := totalElementCost(target.Card.ElementsCost)
	if bonus <= 4 {
		return nil, 0, fmt.Errorf("2121105 sacrifice entry cost must be greater than 4")
	}
	return target, bonus, nil
}

type Card2011102ForesightOrb struct{ AlwaysActive }

func (Card2011102ForesightOrb) ID() string   { return "2011102" }
func (Card2011102ForesightOrb) Name() string { return "预知宝珠" }

func (Card2011102ForesightOrb) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	if ctx.Source.IsHorizontal {
		return fmt.Errorf("预知宝珠需要竖置才能消耗")
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	lookCount := min(3, len(ps.Deck))
	if lookCount == 0 {
		return fmt.Errorf("预知宝珠需要牌库中有牌")
	}
	ctx.Source.IsHorizontal = true
	looked := append([]*CardInstance(nil), ps.Deck[:lookCount]...)
	candidates := make([]map[string]any, 0, len(looked))
	for i, card := range looked {
		info := candidateInfo(card, "deck", "own")
		info["deck_index"] = i
		info["can_select"] = false
		candidates = append(candidates, info)
	}
	ctx.Engine.SetPendingActionWithData(ctx.PlayerID, "foresight_orb_reorder",
		"预知宝珠:将牌库顶3张以任意顺序放回牌库顶或牌库底", candidates, 0, 0,
		func(selected []string, data map[string]any) {
			resolveTopDeckReorder(ctx.Engine, ctx.PlayerID, ctx.PlayerID, looked, data, "foresight_orb_reorder")
		})
	return nil
}

func resolveTopDeckReorder(e *Engine, playerID int, deckOwnerID int, looked []*CardInstance, data map[string]any, effect string) {
	if e == nil {
		return
	}
	ps := e.State.Players[deckOwnerID]
	lookCount := min(len(looked), len(ps.Deck))
	rest := append([]*CardInstance(nil), ps.Deck[lookCount:]...)
	pool := make(map[string]*CardInstance, len(looked))
	for _, card := range looked {
		if card != nil {
			pool[card.InstanceID] = card
		}
	}
	used := make(map[string]bool, len(pool))
	top := orderedWaterDivinationCards(stringsFromActionData(data, "top_order"), pool, used)
	bottom := orderedWaterDivinationCards(stringsFromActionData(data, "bottom_order"), pool, used)
	for _, card := range looked {
		if card == nil || used[card.InstanceID] {
			continue
		}
		top = append(top, card)
		used[card.InstanceID] = true
	}
	ps.Deck = append(append(top, rest...), bottom...)
	e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
		"effect":       effect,
		"deck_owner":   deckOwnerID,
		"top_order":    cardsToInfo(top),
		"bottom_order": cardsToInfo(bottom),
	}})
}

func (e *Engine) hasForesightOrbActive(playerID int) bool {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return false
	}
	ps := e.State.Players[playerID]
	for _, card := range ps.Equipment {
		if card != nil && card.Card != nil && card.Card.Number == "2011102" && !e.hasEffectiveStatus(card, StatusPetrify) {
			return true
		}
	}
	return false
}

var _ PerTurnAbility = Card2011102ForesightOrb{}

type Card1521112CouncilGuard struct{ AlwaysActive }

func (Card1521112CouncilGuard) ID() string   { return "1521112" }
func (Card1521112CouncilGuard) Name() string { return "议庭护法" }

func (Card1521112CouncilGuard) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil {
		return nil
	}
	opponent := ctx.Engine.State.Players[ctx.OpponentID]
	lookCount := min(5, len(opponent.Deck))
	if lookCount == 0 {
		return fmt.Errorf("议庭护法需要对手牌库中有牌")
	}
	looked := append([]*CardInstance(nil), opponent.Deck[:lookCount]...)
	hasMark := false
	for _, card := range looked {
		if card != nil && card.Card != nil && card.Card.Number == "2001102" {
			hasMark = true
			break
		}
	}
	if !hasMark {
		ctx.Engine.shuffleDeck(ctx.OpponentID)
		ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
			"effect":     "council_guard_shuffle",
			"deck_owner": ctx.OpponentID,
		}})
		return nil
	}
	candidates := make([]map[string]any, 0, len(looked))
	for i, card := range looked {
		info := candidateInfo(card, "deck", "enemy")
		info["deck_index"] = i
		info["can_select"] = false
		candidates = append(candidates, info)
	}
	ctx.Engine.SetPendingActionWithData(ctx.PlayerID, "council_guard_reorder",
		"议庭护法:调整对手牌库顶5张的顺序并放回牌库顶或牌库底", candidates, 0, 0,
		func(selected []string, data map[string]any) {
			resolveTopDeckReorder(ctx.Engine, ctx.PlayerID, ctx.OpponentID, looked, data, "council_guard_reorder")
		})
	return nil
}

var _ PerTurnAbility = Card1521112CouncilGuard{}

type Card1511103RoseProphetLori struct{ AlwaysActive }

func (Card1511103RoseProphetLori) ID() string   { return "1511103" }
func (Card1511103RoseProphetLori) Name() string { return "\"玫瑰先知\" 洛莉" }

func (e *Engine) triggerRoseProphetAfterOpponentShuffle(shuffledPlayerID int) {
	if e == nil || shuffledPlayerID < 0 || shuffledPlayerID >= len(e.State.Players) {
		return
	}
	viewerID := 1 - shuffledPlayerID
	shuffled := e.State.Players[shuffledPlayerID]
	if len(shuffled.Deck) == 0 {
		return
	}
	for _, card := range e.getAllFieldCards(e.State.Players[viewerID]) {
		if card == nil || card.Card == nil || card.Card.Number != "1511103" || e.hasEffectiveStatus(card, StatusPetrify) {
			continue
		}
		lookCount := min(3, len(shuffled.Deck))
		looked := append([]*CardInstance(nil), shuffled.Deck[:lookCount]...)
		candidates := make([]map[string]any, 0, len(looked))
		for i, deckCard := range looked {
			info := candidateInfo(deckCard, "deck", "enemy")
			info["deck_index"] = i
			info["can_select"] = false
			candidates = append(candidates, info)
		}
		e.SetPendingActionWithData(viewerID, "rose_prophet_reorder",
			"玫瑰先知:调整对手牌库顶3张的顺序并放回牌库顶或牌库底", candidates, 0, 0,
			func(selected []string, data map[string]any) {
				resolveTopDeckReorder(e, viewerID, shuffledPlayerID, looked, data, "rose_prophet_reorder")
			})
	}
}

const holyChildResolvingStatus = "holy_child_resolving"

type Card1521102HolyChild struct{ AlwaysActive }

func (Card1521102HolyChild) ID() string   { return "1521102" }
func (Card1521102HolyChild) Name() string { return "神圣之子" }
func (Card1521102HolyChild) OnLoadGain(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target != ctx.Source {
		return nil
	}
	ctx.Engine.triggerHolyChildBonus(ctx.PlayerID, ctx.Source)
	return nil
}

func (e *Engine) triggerHolyChildBonus(playerID int, child *CardInstance) {
	if e == nil || child == nil || child.Card == nil || child.Card.Number != "1521102" || child.UltimateUsed || child.Statuses[holyChildResolvingStatus] > 0 {
		return
	}
	choices := []map[string]any{
		{"instance_id": "gain_light_load", "name": "额外获得负载+1光", "zone": "choice", "side": "own"},
		{"instance_id": "gain_life", "name": "额外获得+1血", "zone": "choice", "side": "own"},
	}
	e.SetPendingAction(playerID, "holy_child_bonus",
		"神圣之子:选择额外获得负载或生命", choices, 1, 1,
		func(selected []string) {
			if child == nil || child.CurrentLife <= 0 || child.UltimateUsed || child.Statuses[holyChildResolvingStatus] > 0 {
				return
			}
			child.UltimateUsed = true
			child.Statuses[holyChildResolvingStatus] = 1
			defer delete(child.Statuses, holyChildResolvingStatus)
			switch firstSelected(selected) {
			case "gain_life":
				child.CurrentLife++
				e.emit(GameEvent{
					Type:   "life_gain",
					Player: -1,
					Data: map[string]any{
						"player": playerID,
						"target": cardToInfo(child),
						"amount": 1,
					},
				})
			default:
				e.addElementsGainBonus(child, playerID, model.ElementLight, 1, child)
			}
		})
}

func (e *Engine) triggerHolyChildAfterLifeGain(playerID int, target *CardInstance) {
	if target == nil || target.OwnerID != playerID || target.Card == nil || target.Card.Number != "1521102" {
		return
	}
	e.triggerHolyChildBonus(playerID, target)
}

var _ OnLoadGainBehavior = Card1521102HolyChild{}

type Card2321101ThunderChain struct{ AlwaysActive }

func (Card2321101ThunderChain) ID() string   { return "2321101" }
func (Card2321101ThunderChain) Name() string { return "雷之链" }
func (Card2321101ThunderChain) PerTurnLabel(*CardInstance) string {
	return "充能"
}
func (Card2321101ThunderChain) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	if ctx.Source.IsHorizontal {
		return fmt.Errorf("雷之链需要竖置才能消耗")
	}
	ctx.Source.IsHorizontal = true
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModNextDriveSpellExtraTarget,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		RemainingUses:    1,
	})
	return nil
}

var _ PerTurnAbility = Card2321101ThunderChain{}

type Card3621109SoulRendingScream struct{ AlwaysActive }

func (Card3621109SoulRendingScream) ID() string   { return "3621109" }
func (Card3621109SoulRendingScream) Name() string { return "裂魂尖啸" }
func (Card3621109SoulRendingScream) OnDefend(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil {
		return nil
	}
	attackSkill, _ := ctx.ExtraData["attack_skill"].(*CardInstance)
	if attackSkill != ctx.Source {
		return nil
	}
	defenseSkills, _ := ctx.ExtraData["defense_skills"].([]*CardInstance)
	defenseBoosts, _ := ctx.ExtraData["defense_boosts"].([]*CardInstance)
	weakened := 0
	for _, skill := range append(defenseSkills, defenseBoosts...) {
		if skill == nil || skill.OwnerID == ctx.PlayerID {
			continue
		}
		if ctx.Engine.addStatus(skill, StatusWeaken, 1) {
			weakened++
		}
	}
	if weakened > 0 {
		ctx.Engine.emit(GameEvent{
			Type:   "soul_rending_scream_weaken",
			Player: -1,
			Data: map[string]any{
				"player": ctx.PlayerID,
				"source": cardToInfo(ctx.Source),
				"count":  weakened,
			},
		})
	}
	return nil
}

var _ OnDefendBehavior = Card3621109SoulRendingScream{}

type Card3421104NaturalEcho struct{ AlwaysActive }

func (Card3421104NaturalEcho) ID() string   { return "3421104" }
func (Card3421104NaturalEcho) Name() string { return "自然回响" }
func (Card3421104NaturalEcho) PerTurnLabel(*CardInstance) string {
	return "回响"
}

func (Card3421104NaturalEcho) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	if ctx.Source.UsedThisTurn > 0 {
		return fmt.Errorf("自然回响本回合已经发动")
	}
	candidates := make([]map[string]any, 0)
	for _, card := range ctx.Engine.getAllFieldCards(ctx.Engine.State.Players[ctx.PlayerID]) {
		if card == nil || card.Card == nil || reducibleElementLoad(card, model.ElementEarth) <= 0 {
			continue
		}
		info := candidateInfo(card, "field", "own")
		info["load_element"] = model.ElementEarth
		candidates = append(candidates, info)
	}
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "natural_echo_remove_load",
		"自然回响:移除1点友方卡牌地脉负载", candidates, 1, 1,
		nil, false, func(selected []string, _ map[string]any) error {
			target := ctx.Engine.findFieldCardByInstance(ctx.Engine.State.Players[ctx.PlayerID], firstSelected(selected))
			if target == nil || reducibleElementLoad(target, model.ElementEarth) <= 0 {
				return fmt.Errorf("invalid natural echo target")
			}
			reduceCardElementLoad(target, model.ElementEarth, 1)
			ctx.Source.UsedThisTurn++
			resetCard(ctx.Source)
			ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
				Type:             TempModSkillPowerBonus,
				SourceCardNumber: ctx.Source.Card.Number,
				SourceName:       ctx.Source.Card.Name,
				TargetInstanceID: ctx.Source.InstanceID,
				Amount:           2,
				RemainingUses:    1,
			})
			ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
				Type:             TempModNextSpellExtraTarget,
				SourceCardNumber: ctx.Source.Card.Number,
				SourceName:       ctx.Source.Card.Name,
				TargetInstanceID: ctx.Source.InstanceID,
				RemainingUses:    1,
				AllowSameTarget:  true,
			})
			ctx.Engine.emit(GameEvent{
				Type:   "natural_echo",
				Player: -1,
				Data: map[string]any{
					"player": ctx.PlayerID,
					"source": cardToInfo(ctx.Source),
					"target": cardToInfo(target),
				},
			})
			return nil
		})
	return nil
}

var _ PerTurnAbility = Card3421104NaturalEcho{}

type Card3621106RedMoonDevour struct{ AlwaysActive }

func (Card3621106RedMoonDevour) ID() string   { return "3621106" }
func (Card3621106RedMoonDevour) Name() string { return "红月吞噬" }
func (Card3621106RedMoonDevour) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || !isOwnSpellHit(ctx) {
		return nil
	}
	if ctx.Target.Card == nil || !ctx.Target.Card.IsCompanion() || ctx.Target.OwnerID == ctx.PlayerID {
		return nil
	}
	remainingLife := max(ctx.Target.CurrentLife, 0)
	targetOwner := ctx.Target.OwnerID
	ctx.Engine.destroyUnitWithData(ctx.Target, targetOwner, map[string]any{
		"destroyed_by": ctx.Source.InstanceID,
		"attacker":     ctx.PlayerID,
	})
	if remainingLife <= 0 || !ctx.Engine.redMoonActive(ctx.PlayerID) {
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.Category == model.ElementShadow
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "red_moon_devour_life",
		fmt.Sprintf("红月吞噬:选择1个友方暗影单位获得+%d血", remainingLife), candidates, 1, 1,
		func(selected []string) {
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target == nil || zone != "unit" || target.Card == nil || target.Card.Category != model.ElementShadow {
				return
			}
			target.CurrentLife += remainingLife
			ctx.Engine.triggerHolyChildAfterLifeGain(ctx.PlayerID, target)
			ctx.Engine.emit(GameEvent{
				Type:   "effect_trigger",
				Player: ctx.PlayerID,
				Data: map[string]any{
					"source": cardToInfo(ctx.Source),
					"target": cardToInfo(target),
					"effect": "modify_life",
					"amount": remainingLife,
				},
			})
		})
	return nil
}

var _ OnSpellHitBehavior = Card3621106RedMoonDevour{}

type Card3621108Moonshadow struct{ AlwaysActive }

func (Card3621108Moonshadow) ID() string   { return "3621108" }
func (Card3621108Moonshadow) Name() string { return "月影" }
func (Card3621108Moonshadow) OnDefend(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil {
		return nil
	}
	attackSkill, _ := ctx.ExtraData["attack_skill"].(*CardInstance)
	if attackSkill != ctx.Source {
		return nil
	}
	defenseSkills, _ := ctx.ExtraData["defense_skills"].([]*CardInstance)
	defenseBoosts, _ := ctx.ExtraData["defense_boosts"].([]*CardInstance)
	for _, skill := range append(defenseSkills, defenseBoosts...) {
		if skill != nil && skill.OwnerID != ctx.PlayerID && skill.Statuses[StatusWeaken] > 0 && ctx.Engine.hasEffectiveStatus(skill, StatusWeaken) {
			source := ctx.Source
			reason := skill
			ctx.Engine.SetPendingAction(ctx.PlayerID, "moonshadow_reset",
				"月影:是否重置此卡", []map[string]any{candidateInfo(source, "skill", "own")}, 0, 1,
				func(selected []string) {
					if len(selected) == 0 || source == nil || source.Card == nil {
						return
					}
					ctx.Engine.resetCard(source)
					ctx.Engine.emit(GameEvent{
						Type:   "moonshadow_reset",
						Player: -1,
						Data: map[string]any{
							"player": ctx.PlayerID,
							"source": cardToInfo(source),
							"reason": cardToInfo(reason),
						},
					})
				})
			return nil
		}
	}
	return nil
}

var _ OnDefendBehavior = Card3621108Moonshadow{}

type Card4511102RedeemerEveAutumnMaple struct{ AlwaysActive }

func (Card4511102RedeemerEveAutumnMaple) ID() string   { return "4511102" }
func (Card4511102RedeemerEveAutumnMaple) Name() string { return "救赎者 伊芙 秋枫" }
func (Card4511102RedeemerEveAutumnMaple) OnUltimate(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.UltimateUsed {
		return nil
	}
	ownUnits := ctx.Engine.friendlyUnits(ctx.PlayerID, true, nil)
	enemyUnits := ctx.Engine.enemyUnits(ctx.PlayerID, true, nil)
	if len(enemyUnits) <= len(ownUnits) {
		return fmt.Errorf("敌方场上单位数量必须比我方多")
	}
	woundedCount := 0
	forEachPlayerUnit(ctx.Engine.State.Players[ctx.PlayerID], true, func(unit *CardInstance) {
		if unit != nil && unit.CurrentLife < maxLife(unit) {
			woundedCount++
		}
	})
	if woundedCount <= 0 {
		return fmt.Errorf("救赎者 伊芙 秋枫需要受伤的友方单位")
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ps.Elements[model.ElementLight] < woundedCount {
		return fmt.Errorf("救赎者 伊芙 秋枫需要%d点光元素", woundedCount)
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() && !card.Card.IsHero() && card.CurrentLife < maxLife(card)
	})
	if len(candidates) == 0 {
		return fmt.Errorf("救赎者 伊芙 秋枫需要受伤的友方伙伴")
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "redeemer_eve_autumn_maple_target",
		fmt.Sprintf("救赎者 伊芙 秋枫:支付%d光,选择1个受伤友方伙伴获得+%d血和负载+%d光", woundedCount, woundedCount, woundedCount),
		candidates, 1, 1,
		func(selected []string) {
			if ctx.Source == nil || !ctx.Engine.cardStillOnField(ctx.Source) {
				return
			}
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target == nil || zone != "unit" || target.Card == nil || !target.Card.IsCompanion() || target.Card.IsHero() || target.CurrentLife >= maxLife(target) {
				return
			}
			ownNow := ctx.Engine.friendlyUnits(ctx.PlayerID, true, nil)
			enemyNow := ctx.Engine.enemyUnits(ctx.PlayerID, true, nil)
			if len(enemyNow) <= len(ownNow) {
				return
			}
			x := 0
			forEachPlayerUnit(ctx.Engine.State.Players[ctx.PlayerID], true, func(unit *CardInstance) {
				if unit != nil && unit.CurrentLife < maxLife(unit) {
					x++
				}
			})
			if x <= 0 || ps.Elements[model.ElementLight] < x {
				return
			}
			ps.Elements[model.ElementLight] -= x
			target.Statuses["max_life_bonus"] += x
			target.CurrentLife += x
			ctx.Engine.triggerHolyChildAfterLifeGain(ctx.PlayerID, target)
			ctx.Engine.addElementsGainBonus(target, ctx.PlayerID, model.ElementLight, x, ctx.Source)
			ctx.Source.UltimateUsed = true
			ctx.Engine.emit(GameEvent{
				Type:   "redeemer_eve_autumn_maple_blessing",
				Player: -1,
				Data: map[string]any{
					"player": ctx.PlayerID,
					"source": cardToInfo(ctx.Source),
					"target": cardToInfo(target),
					"amount": x,
				},
			})
		})
	return nil
}

var _ UltimateAbility = Card4511102RedeemerEveAutumnMaple{}

func forEachPlayerUnit(ps *PlayerState, includeHero bool, fn func(*CardInstance)) {
	if ps == nil || fn == nil {
		return
	}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := ps.Units[col][row]
			if unit == nil {
				continue
			}
			if !includeHero && unit.Card != nil && unit.Card.IsHero() {
				continue
			}
			fn(unit)
		}
	}
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

type Card3421107Burrow struct{ AlwaysActive }

func (Card3421107Burrow) ID() string   { return "3421107" }
func (Card3421107Burrow) Name() string { return "破土而出" }

func (Card3421107Burrow) MasteryMax() int { return 2 }

func (Card3421107Burrow) OnMastery(*EffectContext, int) error { return nil }

var _ MasteryBehavior = Card3421107Burrow{}

type Card3121110CursedFire struct{ AlwaysActive }

func (Card3121110CursedFire) ID() string   { return "3121110" }
func (Card3121110CursedFire) Name() string { return "咒火" }

func (Card3121110CursedFire) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.Source.Card.Number != "3121110" || ctx.PlayerID < 0 {
		return nil
	}
	drawn := ctx.Engine.flipDeckMatchesToHandThen(ctx.PlayerID, 1, 0, isCheapFireSpellScroll, func(drawn []*CardInstance) {
		for _, card := range drawn {
			makeEntryCostZero(card)
		}
	})
	if len(drawn) > 0 {
		ctx.Engine.emit(GameEvent{
			Type:   "cursed_fire_flip_scroll",
			Player: ctx.PlayerID,
			Data: map[string]any{
				"player": ctx.PlayerID,
				"source": cardToInfo(ctx.Source),
				"card":   cardToInfo(drawn[0]),
			},
		})
		ctx.Engine.promptCursedFireImmediateScrollUse(ctx.PlayerID, ctx.Source, drawn[0])
	}
	return nil
}

func (e *Engine) promptCursedFireImmediateScrollUse(playerID int, source *CardInstance, scroll *CardInstance) {
	if e == nil || source == nil || scroll == nil || scroll.Card == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	if isDefenseOnlySkill(scroll.Card) || !canUseSkillForPurpose(scroll.Card, skillPurposeAttack) {
		return
	}
	if !skillNeedsTargetInstance(scroll) {
		e.SetPendingAction(playerID, "cursed_fire_use_scroll",
			"咒火:是否立刻使用翻取的卷轴", []map[string]any{candidateInfo(scroll, "hand", "own")}, 0, 1,
			func(selected []string) {
				if len(selected) == 0 {
					return
				}
				e.useCursedFireScrollFromHand(playerID, scroll, SpellTarget{Type: "none"})
			})
		return
	}
	candidates := e.enemyUnits(playerID, true, func(target *CardInstance) bool {
		if target == nil || target.Position == nil {
			return false
		}
		return e.validateSpellTarget(playerID, scroll, SpellTarget{Type: "unit", Position: *target.Position}) == nil
	})
	if len(candidates) == 0 {
		return
	}
	e.SetPendingAction(playerID, "cursed_fire_use_scroll_target",
		"咒火:选择目标并立刻使用翻取的卷轴", candidates, 0, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			target := e.findUnitByInstanceID(firstSelected(selected))
			if target == nil || target.Position == nil {
				return
			}
			targetOwner := target.OwnerID
			e.useCursedFireScrollFromHand(playerID, scroll, SpellTarget{Type: "unit", Position: *target.Position, OwnerID: &targetOwner})
		})
}

func (e *Engine) useCursedFireScrollFromHand(playerID int, scroll *CardInstance, target SpellTarget) {
	ps := e.State.Players[playerID]
	current, handIdx := ps.FindHandCard(scroll.InstanceID)
	if current != scroll || handIdx < 0 {
		return
	}
	data := map[string]any{"instance_id": scroll.InstanceID, "target_type": target.Type}
	if target.Type == "unit" {
		data["target_col"] = float64(target.Position.Col)
		data["target_row"] = float64(target.Position.Row)
		if target.OwnerID != nil {
			data["target_owner"] = float64(*target.OwnerID)
		}
	}
	_ = e.handleUseSpellScrollItem(playerID, ActionMessage{Data: data}, scroll, handIdx)
}

func isCheapFireSpellScroll(card *CardInstance) bool {
	return card != nil &&
		card.Card != nil &&
		card.Card.Category == model.ElementFire &&
		isSpellScrollCard(card.Card) &&
		totalElementCost(card.Card.ElementsCost) < 4
}

var _ OnSpellHitBehavior = Card3121110CursedFire{}

type Card1221114JadeFacedSnowFox struct{ AlwaysActive }

func (Card1221114JadeFacedSnowFox) ID() string   { return "1221114" }
func (Card1221114JadeFacedSnowFox) Name() string { return "玉面雪狐" }

func (Card1221114JadeFacedSnowFox) HasActiveSpellReaction(card *CardInstance) bool {
	return card != nil && card.Position != nil && !card.UltimateUsed
}

func (Card1221114JadeFacedSnowFox) CanReactToSpell(ctx *EffectContext, spell *SpellCast) bool {
	return ctx != nil &&
		ctx.Engine != nil &&
		ctx.Source != nil &&
		ctx.Source.Position != nil &&
		!ctx.Source.UltimateUsed &&
		spell != nil &&
		spell.AttackerID != ctx.PlayerID &&
		spell.Skill != nil &&
		spell.Skill.Card != nil &&
		canUseSkillForPurpose(spell.Skill.Card, skillPurposeAttack)
}

func (Card1221114JadeFacedSnowFox) OnSpellReaction(ctx *EffectContext, spell *SpellCast) error {
	if !(Card1221114JadeFacedSnowFox{}).CanReactToSpell(ctx, spell) {
		return nil
	}
	sourceID := ctx.Source.InstanceID
	positions := ctx.Engine.allUnitPositionsForPlayer(ctx.PlayerID, ctx.PlayerID)
	ctx.Engine.SetPendingAction(ctx.PlayerID, "jade_faced_snow_fox_move",
		"玉面雪狐:移动此卡并获得2水", positions, 1, 1,
		func(selected []string) {
			if ctx.Engine.State.PendingSpell == nil {
				return
			}
			source, _ := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, sourceID)
			if source == nil || source.Position == nil || source.UltimateUsed {
				return
			}
			pos, ok := positionFromSelectionID(firstSelected(selected))
			if !ok {
				return
			}
			ctx.Engine.moveOrSwapUnitToPosition(ctx.PlayerID, sourceID, pos)
			ctx.Engine.State.Players[ctx.PlayerID].Elements[model.ElementWater] += 2
			source.UltimateUsed = true
			ctx.Engine.emit(GameEvent{
				Type:   "jade_faced_snow_fox_move",
				Player: -1,
				Data: map[string]any{
					"player": ctx.PlayerID,
					"source": cardToInfo(source),
					"water":  2,
				},
			})
			ctx.Engine.promptIllusionScrollRetarget(ctx.PlayerID, source, ctx.ExtraData)
		})
	return nil
}

var _ SpellReactionBehavior = Card1221114JadeFacedSnowFox{}

type Card2321108ScatterAway struct{ AlwaysActive }

func (Card2321108ScatterAway) ID() string   { return "2321108" }
func (Card2321108ScatterAway) Name() string { return "散去" }

func (Card2321108ScatterAway) OnDamaged(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Target == nil || ctx.Target.Card == nil {
		return nil
	}
	if ctx.Target.OwnerID != ctx.PlayerID || !ctx.Target.Card.IsCompanion() || ctx.Target.Card.Category != model.ElementAir || damageFromData(ctx.ExtraData) <= 0 {
		return nil
	}
	ctx.Target.Statuses[temporaryDamageAndNegativeImmunityUntilStatus] = ctx.Engine.State.TurnNumber + 1
	ctx.Engine.emit(GameEvent{
		Type:   "scatter_away_immunity",
		Player: -1,
		Data: map[string]any{
			"player": ctx.PlayerID,
			"source": cardToInfo(ctx.Source),
			"target": cardToInfo(ctx.Target),
		},
	})
	return nil
}

var _ OnDamagedBehavior = Card2321108ScatterAway{}

type Card2021115InfusionRuneE struct{ AlwaysActive }

func (Card2021115InfusionRuneE) ID() string   { return "2021115" }
func (Card2021115InfusionRuneE) Name() string { return "注能符文E型" }

func (Card2021115InfusionRuneE) OnDefend(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.ExtraData == nil {
		return nil
	}
	reset := 0
	for _, skill := range append(spellInstancesFromData(ctx.ExtraData, "defense_skills"), spellInstancesFromData(ctx.ExtraData, "defense_boosts")...) {
		if skill == nil || skill.Card == nil || !skill.Card.IsSkill() || skill.OwnerID != ctx.PlayerID {
			continue
		}
		skill.IsHorizontal = false
		reset++
	}
	if reset > 0 {
		ctx.Engine.emit(GameEvent{
			Type:   "infusion_rune_e_reset",
			Player: -1,
			Data: map[string]any{
				"player": ctx.PlayerID,
				"source": cardToInfo(ctx.Source),
				"count":  reset,
			},
		})
	}
	return nil
}

var _ OnDefendBehavior = Card2021115InfusionRuneE{}

type Card2121104FireRebirthScroll struct{ AlwaysActive }

func (Card2121104FireRebirthScroll) ID() string   { return "2121104" }
func (Card2121104FireRebirthScroll) Name() string { return "浴火重生卷轴" }

func (Card2121104FireRebirthScroll) OnTurnEnd(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	candidates := ctx.Engine.rebirthScrollReviveCandidates(ctx.PlayerID)
	if len(candidates) == 0 {
		return nil
	}
	if len(ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)) == 0 {
		return nil
	}
	if len(candidates) == 1 {
		ctx.Engine.promptFireRebirthScrollPosition(ctx.PlayerID, ctx.Source, candidates[0].InstanceID)
		return nil
	}
	choices := make([]map[string]any, 0, len(candidates))
	for _, card := range candidates {
		info := cardToInfo(card)
		info["zone"] = "graveyard"
		info["side"] = "own"
		choices = append(choices, info)
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "fire_rebirth_scroll",
		"浴火重生卷轴:选择1个本回合死亡的友方火焰伙伴复活", choices, 1, 1,
		func(selected []string) {
			ctx.Engine.promptFireRebirthScrollPosition(ctx.PlayerID, ctx.Source, firstSelected(selected))
		})
	return nil
}

func (e *Engine) promptFireRebirthScrollPosition(playerID int, source *CardInstance, instanceID string) {
	positions := e.friendlyEmptyUnitPositions(playerID)
	if len(positions) == 0 || e.findRecentFireRebirthCandidate(playerID, instanceID) == nil {
		return
	}
	e.SetPendingAction(playerID, "fire_rebirth_scroll_position",
		"浴火重生卷轴:选择复活位置", positions, 1, 1,
		func(selected []string) {
			pos, ok := positionFromSelectionID(firstSelected(selected))
			if !ok {
				return
			}
			if e.reviveRecentFireCompanionAtPosition(playerID, instanceID, pos) {
				e.emit(GameEvent{
					Type:   "fire_rebirth_scroll_revive",
					Player: -1,
					Data: map[string]any{
						"player":   playerID,
						"source":   cardToInfo(source),
						"revived":  instanceID,
						"position": pos,
						"count":    1,
					},
				})
			}
		})
}

func (e *Engine) findRecentFireRebirthCandidate(playerID int, instanceID string) *CardInstance {
	for _, card := range e.rebirthScrollReviveCandidates(playerID) {
		if card != nil && card.InstanceID == instanceID {
			return card
		}
	}
	return nil
}

func (e *Engine) reviveRecentFireCompanionAtPosition(playerID int, instanceID string, pos Position) bool {
	card := e.findRecentFireRebirthCandidate(playerID, instanceID)
	if card == nil || !pos.Valid() || e.State.Players[playerID].Units[pos.Col][pos.Row] != nil {
		return false
	}
	if e.reviveCompanionFromGraveyardWithLifeAtPosition(playerID, card.InstanceID, 0, false, pos) {
		card.IsHorizontal = false
		if card.Statuses != nil {
			delete(card.Statuses, enteredGraveyardTurnStatus)
		}
		return true
	}
	return false
}

func (e *Engine) rebirthScrollReviveCandidates(playerID int) []*CardInstance {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	ps := e.State.Players[playerID]
	candidates := make([]*CardInstance, 0)
	for _, card := range ps.Graveyard {
		if card == nil || card.Card == nil || !card.Card.IsCompanion() || card.Card.Category != model.ElementFire {
			continue
		}
		if card.Statuses[enteredGraveyardTurnStatus] == e.State.TurnNumber {
			candidates = append(candidates, card)
		}
	}
	return candidates
}

var _ OnTurnEndBehavior = Card2121104FireRebirthScroll{}
