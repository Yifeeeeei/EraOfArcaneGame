package game

import (
	"fmt"

	"eraofarcane/model"
)

func emitBatchEffect(ctx *EffectContext, effect string) {
	if ctx == nil || ctx.Engine == nil {
		return
	}
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(ctx.Source),
		"effect": effect,
	}})
}

func selectedUnitFromCandidates(e *Engine, selected []string, candidates []map[string]any) *CardInstance {
	if e == nil || len(selected) == 0 {
		return nil
	}
	allowed := make(map[string]bool, len(candidates))
	for _, candidate := range candidates {
		if id, _ := candidate["instance_id"].(string); id != "" {
			allowed[id] = true
		}
	}
	for _, id := range selected {
		if allowed[id] {
			return e.findUnitByInstanceID(id)
		}
	}
	return nil
}

func healUnit(card *CardInstance, amount int) {
	if card == nil || amount <= 0 {
		return
	}
	wasFull := maxLife(card) > 0 && card.CurrentLife >= maxLife(card)
	card.CurrentLife += amount
	if life := maxLife(card); life > 0 && card.CurrentLife > life {
		card.CurrentLife = life
	}
	if wasFull && card.Card != nil && card.Card.Number == "1521016" {
		card.Statuses["max_life_bonus"]++
		card.CurrentLife++
	}
}

func maxLife(card *CardInstance) int {
	if card == nil || card.Card == nil {
		return 0
	}
	return card.Card.Life + card.Statuses["max_life_bonus"]
}

func addGeneratedCardToHand(ctx *EffectContext, cardNumber string) {
	card := getCardDB()[cardNumber]
	if card == nil {
		return
	}
	ctx.Engine.State.Players[ctx.PlayerID].Hand = append(ctx.Engine.State.Players[ctx.PlayerID].Hand, NewCardInstance(card, ctx.PlayerID, ctx.Engine.State.TurnNumber))
	emitBatchEffect(ctx, "add_generated_card_to_hand")
}

func resetInstance(card *CardInstance) {
	if card == nil {
		return
	}
	card.IsHorizontal = false
	card.Statuses[StatusCooldown] = 0
	card.UsedThisTurn = 0
}

func reduceCost(cost map[string]int, elem string, amount int) {
	if amount <= 0 {
		return
	}
	if cost[elem] > 0 {
		cost[elem] -= amount
		if cost[elem] < 0 {
			cost[elem] = 0
		}
	}
}

func summonHandCompanionFree(ctx *EffectContext, predicate func(*CardInstance) bool) *CardInstance {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	pos := ps.FindEmptyPosition()
	if pos == nil {
		return nil
	}
	for i, card := range ps.Hand {
		if card == nil || !card.Card.IsCompanion() || (predicate != nil && !predicate(card)) {
			continue
		}
		ps.Hand = append(ps.Hand[:i], ps.Hand[i+1:]...)
		card.Position = pos
		card.IsHorizontal = true
		card.EnterTurn = ctx.Engine.State.TurnNumber
		ps.Units[pos.Col][pos.Row] = card
		ctx.Engine.triggerEffects(TriggerOnEnter, card, nil, nil)
		ctx.Engine.triggerFieldEffectsWithData(TriggerOnUnitEnter, ctx.PlayerID, card, map[string]any{"entered_player": ctx.PlayerID})
		return card
	}
	return nil
}

func totalFieldLoad(ps *PlayerState) int {
	total := 0
	for _, card := range ps.Units {
		for _, unit := range card {
			for _, amount := range effectiveElementsGain(unit) {
				total += amount
			}
		}
	}
	for _, card := range ps.Equipment {
		for _, amount := range effectiveElementsGain(card) {
			total += amount
		}
	}
	return total
}

type Card1011002WizardTower struct{ AlwaysActive }

func (Card1011002WizardTower) ID() string   { return "1011002" }
func (Card1011002WizardTower) Name() string { return "巫师之塔 通天阁" }
func (Card1011002WizardTower) HasGlobalSpellRange() bool {
	return true
}
func (Card1011002WizardTower) OnEnter(ctx *EffectContext) error {
	ctx.Source.Statuses["全场法力范围"] = 1
	return nil
}

type Card1021008ForesightProphet struct{ AlwaysActive }

func (Card1021008ForesightProphet) ID() string                            { return "1021008" }
func (Card1021008ForesightProphet) Name() string                          { return "预见先知" }
func (Card1021008ForesightProphet) HasActiveTurnStart(*CardInstance) bool { return false }
func (Card1021008ForesightProphet) OnTurnStart(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if len(ps.Deck) == 0 || ctx.Engine.State.PendingAction != nil {
		return nil
	}
	card := ps.Deck[0]
	candidates := []map[string]any{candidateInfo(card, "deck", "own")}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "foresight_prophet_top_card",
		"预见先知:查看牌堆顶1张牌,选择它则置于牌堆底,不选择则放回牌堆顶",
		candidates, 0, 1, func(selected []string) {
			if len(selected) == 0 || selected[0] != card.InstanceID || len(ps.Deck) == 0 || ps.Deck[0] != card {
				return
			}
			ps.Deck = append(ps.Deck[1:], card)
			emitBatchEffect(ctx, "peek_top_to_bottom")
		})
	return nil
}

type Card1021014ImpatientJunior struct{ AlwaysActive }

func (Card1021014ImpatientJunior) ID() string   { return "1021014" }
func (Card1021014ImpatientJunior) Name() string { return "急不可耐的小师弟" }
func (Card1021014ImpatientJunior) OnEnter(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{Type: TempModNextLearnedSkillHaste, RemainingUses: 1, ExpiresTurn: ctx.Engine.State.TurnNumber + 1})
	return nil
}

type Card1111003Bifang struct{ AlwaysActive }

func (Card1111003Bifang) ID() string   { return "1111003" }
func (Card1111003Bifang) Name() string { return "毕方" }
func (Card1111003Bifang) OnDamaged(ctx *EffectContext) error {
	if ctx.ExtraData != nil {
		if damagedPlayer, ok := ctx.ExtraData["damaged_player"].(int); ok && damagedPlayer != ctx.PlayerID && ctx.Target != nil && ctx.ExtraData["status_damage"] == StatusBurn {
			ctx.Engine.dealDamage(ctx.Target, 1, damagedPlayer)
		}
	}
	return nil
}

type Card1121012FireInsight struct{ AlwaysActive }

func (Card1121012FireInsight) ID() string   { return "1121012" }
func (Card1121012FireInsight) Name() string { return "火焰洞察者" }
func (Card1121012FireInsight) OnDamaged(ctx *EffectContext) error {
	if ctx.ExtraData != nil {
		if ctx.ExtraData["damage_element"] == model.ElementFire || ctx.ExtraData["status_damage"] == StatusBurn {
			ctx.Engine.drawCards(ctx.PlayerID, 1)
		}
	}
	return nil
}

type Card1121013Arsonist struct{ AlwaysActive }

func (Card1121013Arsonist) ID() string   { return "1121013" }
func (Card1121013Arsonist) Name() string { return "纵火者" }
func (Card1121013Arsonist) OnSpellCast(ctx *EffectContext) error {
	if ctx.Source.UsedThisTurn > 0 || !isFriendlySpellCast(ctx) || spellCastSourceElement(ctx) != model.ElementFire || ctx.Target == nil || isSorcerySkill(ctx.Target.Card) {
		return nil
	}
	candidates := append(ctx.Engine.friendlyUnits(ctx.PlayerID, true, nil), ctx.Engine.enemyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card.Position != nil && ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, cardHasPierce(ctx.Target))
	})...)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "arsonist_burn",
		"纵火者:可以选择法力范围内1个单位点燃1", candidates, 0, 1,
		func(selected []string) {
			target := selectedUnitFromCandidates(ctx.Engine, selected, candidates)
			if target != nil {
				ctx.Engine.addStatus(target, StatusBurn, 1)
				ctx.Source.UsedThisTurn++
			}
		})
	return nil
}

const leviathanCooldownStatus = "leviathan_cooldown"

type Card1211002Leviathan struct{ AlwaysActive }

func (Card1211002Leviathan) ID() string   { return "1211002" }
func (Card1211002Leviathan) Name() string { return "深渊巨口 利维坦" }

func (Card1211002Leviathan) PerTurnLabel(*CardInstance) string {
	return "消耗"
}

func (Card1211002Leviathan) OnPerTurn(ctx *EffectContext) error {
	if ctx.Source.IsHorizontal {
		return fmt.Errorf("深渊巨口 利维坦已经横置")
	}
	if ctx.Source.Statuses[leviathanCooldownStatus] > 0 {
		return fmt.Errorf("深渊巨口 利维坦的主动效果正在冷却")
	}
	targets := ctx.Engine.enemyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card.Card.IsCompanion() && card.Position != nil && ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false)
	})
	if len(targets) == 0 {
		return fmt.Errorf("没有可消灭的敌方伙伴")
	}
	ctx.Source.IsHorizontal = true
	ctx.Engine.SetPendingAction(ctx.PlayerID, "leviathan_destroy",
		"利维坦:选择法力范围内1个敌方伙伴消灭", targets, 1, 1,
		func(selected []string) {
			target := selectedUnitFromCandidates(ctx.Engine, selected, targets)
			if target != nil {
				ctx.Engine.destroyUnit(target, ctx.OpponentID)
				ctx.Source.Statuses[leviathanCooldownStatus] = 2
			}
		})
	return nil
}

func (Card1211002Leviathan) OnTurnEnd(ctx *EffectContext) error {
	if endedPlayer, ok := ctx.ExtraData["ended_player"].(int); ok && endedPlayer != ctx.PlayerID {
		return nil
	}
	if ctx.Source.Statuses[leviathanCooldownStatus] > 0 {
		ctx.Source.Statuses[leviathanCooldownStatus]--
	}
	return nil
}

type Card1211003SnowWoman struct{ AlwaysActive }

func (Card1211003SnowWoman) ID() string   { return "1211003" }
func (Card1211003SnowWoman) Name() string { return "\"雪女\" 天户凌" }
func (Card1211003SnowWoman) OnEnter(ctx *EffectContext) error {
	ctx.Source.Statuses["引魔"] = 1
	return nil
}
func (Card1211003SnowWoman) HasActivePerTurn(*CardInstance) bool { return false }
func (Card1211003SnowWoman) OnPerTurn(ctx *EffectContext) error {
	frontRow := ctx.Engine.State.Players[ctx.OpponentID].GetFrontRow()
	if frontRow < 0 {
		return nil
	}
	targets := ctx.Engine.enemyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card.Position != nil && card.Position.Row == frontRow
	})
	if len(targets) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "snow_woman_freeze",
		"雪女:选择1个前排敌人冻结1", targets, 1, 1,
		func(selected []string) {
			target := selectedUnitFromCandidates(ctx.Engine, selected, targets)
			if target != nil {
				ctx.Engine.addStatus(target, StatusFreeze, 1)
			}
		})
	return nil
}

type Card1221001DolphinPartner struct{ AlwaysActive }

func (Card1221001DolphinPartner) ID() string   { return "1221001" }
func (Card1221001DolphinPartner) Name() string { return "海豚伙伴" }

type Card1221010WallKeeper struct{ AlwaysActive }

func (Card1221010WallKeeper) ID() string   { return "1221010" }
func (Card1221010WallKeeper) Name() string { return "护壁者" }
func (Card1221010WallKeeper) OnEnter(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{Type: TempModAllSpellDamageZero, RemainingUses: 1, ExpiresTurn: ctx.Engine.State.TurnNumber + 1})
	return nil
}

type Card1221012DragonPrinceDescendant struct{ AlwaysActive }

func (Card1221012DragonPrinceDescendant) ID() string      { return "1221012" }
func (Card1221012DragonPrinceDescendant) Name() string    { return "龙王子裔" }
func (Card1221012DragonPrinceDescendant) MasteryMax() int { return 2 }
func (Card1221012DragonPrinceDescendant) OnMastery(ctx *EffectContext, level int) error {
	if level != 2 {
		return nil
	}
	searchDeckToHandByPredicateWithResult(ctx, "dragon_prince_search", "检索1个水纹伙伴并使其入场费用-1水", isWaterCompanion, func(card *CardInstance) {
		card.Statuses["入场费用水-1"]++
	})
	return nil
}

type Card1311003WindBladeKarina struct{ AlwaysActive }

func (Card1311003WindBladeKarina) ID() string   { return "1311003" }
func (Card1311003WindBladeKarina) Name() string { return "\"风刃\" 卡琳娜" }
func (Card1311003WindBladeKarina) ModifySkillUseCost(ctx *EffectContext, cost map[string]int) {
	if ctx.Source != nil && ctx.Source.Card.Category == model.ElementAir && skillNeedsTargetInstance(ctx.Source) && !cardHasPierce(ctx.Source) {
		cost[model.ElementAir]++
	}
}

type Card1321013TeleportMage struct{ AlwaysActive }

func (Card1321013TeleportMage) ID() string   { return "1321013" }
func (Card1321013TeleportMage) Name() string { return "传送法师" }
func (Card1321013TeleportMage) OnPerTurn(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	var moving *CardInstance
	for col := 0; col < 3 && moving == nil; col++ {
		for row := 0; row < 3; row++ {
			if u := ps.Units[col][row]; u != nil && !u.Card.IsHero() {
				moving = u
				break
			}
		}
	}
	pos := ps.FindEmptyPosition()
	if moving != nil && pos != nil {
		ps.Units[moving.Position.Col][moving.Position.Row] = nil
		moving.Position = pos
		ps.Units[pos.Col][pos.Row] = moving
	}
	return nil
}

type Card1321015WindSpeaker struct{ AlwaysActive }

func (Card1321015WindSpeaker) ID() string   { return "1321015" }
func (Card1321015WindSpeaker) Name() string { return "风语者" }
func (Card1321015WindSpeaker) OnPerTurn(ctx *EffectContext) error {
	ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementAir: 1})
	return nil
}

type Card1401001LifeSeed struct{ AlwaysActive }

func (Card1401001LifeSeed) ID() string      { return "1401001" }
func (Card1401001LifeSeed) Name() string    { return "生命种子" }
func (Card1401001LifeSeed) MasteryMax() int { return 2 }
func (Card1401001LifeSeed) OnMastery(ctx *EffectContext, level int) error {
	if level != 2 {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ps.FindEmptyPosition() == nil {
		return nil
	}
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, func(card *CardInstance) bool {
		return card.Card.IsCompanion() && card.Card.Category == model.ElementEarth
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "life_seed_summon", "生命种子:可以召唤1个地属性伙伴并继承生命种子的加成", candidates, 0, 1, func(selected []string) {
		if len(selected) == 0 {
			return
		}
		cardID := selected[0]
		positions := ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)
		if len(positions) == 0 {
			return
		}
		ctx.Engine.SetPendingAction(ctx.PlayerID, "life_seed_summon_position", "生命种子:选择召唤位置", positions, 1, 1, func(posSelected []string) {
			if len(posSelected) == 0 {
				return
			}
			pos, ok := positionFromSelectionID(posSelected[0])
			if !ok {
				return
			}
			summoned := summonCardFreeFromHandOrDeckAtPosition(ctx, cardID, pos)
			if summoned == nil {
				return
			}
			inheritLifeSeedBonuses(ctx.Engine, ctx.Source, summoned, ctx.PlayerID)
			ctx.Engine.destroyUnit(ctx.Source, ctx.PlayerID)
		})
	})
	return nil
}

func inheritLifeSeedBonuses(e *Engine, source *CardInstance, target *CardInstance, playerID int) {
	if source == nil || target == nil {
		return
	}
	if bonusLife := source.CurrentLife - source.Card.Life; bonusLife > 0 {
		target.CurrentLife += bonusLife
	}
	target.CurrentAttack += max(source.CurrentAttack-source.Card.Attack, 0)
	target.PowerBonus += source.PowerBonus
	target.AttackBonus += source.AttackBonus
	for elem, amount := range source.ElementsGainBonus {
		if amount != 0 {
			e.addElementsGainBonus(target, playerID, elem, amount, source)
		}
	}
	if len(source.ElementsGainSet) > 0 {
		target.ElementsGainSet = copyElementCost(source.ElementsGainSet)
	}
	for status, amount := range source.Statuses {
		if amount <= 0 || status == StatusMastery {
			continue
		}
		target.Statuses[status] += amount
	}
}

type Card1401002SpiritBeastXinke struct{ AlwaysActive }

func (Card1401002SpiritBeastXinke) ID() string   { return "1401002" }
func (Card1401002SpiritBeastXinke) Name() string { return "灵兽 辛柯" }
func (Card1401002SpiritBeastXinke) OnFriendlyDamagedFromHidden(ctx *EffectContext) error {
	if ctx.ExtraData == nil {
		return nil
	}
	attacker, hasAttacker := ctx.ExtraData["attacker"].(int)
	if damagedPlayer, ok := ctx.ExtraData["damaged_player"].(int); ok && damagedPlayer == ctx.PlayerID && hasAttacker && attacker != ctx.PlayerID {
		if ctx.Engine.State.PendingAction != nil {
			return nil
		}
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		candidates := make([]map[string]any, 0)
		for _, card := range ps.Hand {
			if card != nil && card.Card.Number == "1401002" {
				candidates = append(candidates, candidateInfo(card, "hand", "own"))
			}
		}
		candidates = append(candidates, ctx.Engine.friendlyDeckCards(ctx.PlayerID, func(card *CardInstance) bool { return card.Card.Number == "1401002" })...)
		ctx.Engine.SetPendingAction(ctx.PlayerID, "xinke_summon", "免费召唤灵兽 辛柯", candidates, 0, 1,
			func(selected []string) {
				if len(selected) > 0 {
					cardID := selected[0]
					positions := ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)
					if len(positions) == 0 {
						return
					}
					ctx.Engine.SetPendingAction(ctx.PlayerID, "xinke_summon_position", "选择灵兽 辛柯的入场位置", positions, 1, 1,
						func(posSelected []string) {
							if len(posSelected) == 0 {
								return
							}
							pos, ok := positionFromSelectionID(posSelected[0])
							if !ok {
								return
							}
							summonCardFreeFromHandOrDeckAtPosition(ctx, cardID, pos)
						})
				}
			})
	}
	return nil
}

type Card1411001GreatDruidCycle struct{ AlwaysActive }

func (Card1411001GreatDruidCycle) ID() string                           { return "1411001" }
func (Card1411001GreatDruidCycle) Name() string                         { return "\"轮回不息\" 大德鲁伊 烟尘" }
func (Card1411001GreatDruidCycle) HasActiveUltimate(*CardInstance) bool { return false }
func (Card1411001GreatDruidCycle) OnUltimate(ctx *EffectContext) error {
	ctx.Source.Statuses["轮回不息"] = 1
	return nil
}
func (Card1411001GreatDruidCycle) OnFriendlyDeath(ctx *EffectContext) error {
	if ctx.Source.UltimateUsed {
		return nil
	}
	ctx.Source.Statuses["轮回不息"] = 1
	if ctx.Source.Statuses["轮回不息"] == 0 || ctx.Target == nil || !ctx.Target.Card.IsCompanion() {
		return nil
	}
	if ctx.Engine.State.Players[ctx.PlayerID].FindEmptyPosition() == nil || getCardDB()["1401001"] == nil {
		return nil
	}
	ctx.Source.UltimateUsed = true
	candidates := []map[string]any{candidateInfo(ctx.Source, "companion", "own")}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "great_druid_life_seed", "\"轮回不息\" 大德鲁伊 烟尘:是否召唤1个生命种子", candidates, 0, 1, func(selected []string) {
		if len(selected) == 0 {
			return
		}
		positions := ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)
		if len(positions) == 0 {
			return
		}
		ctx.Engine.SetPendingAction(ctx.PlayerID, "great_druid_life_seed_position", "\"轮回不息\" 大德鲁伊 烟尘:选择生命种子位置", positions, 1, 1, func(posSelected []string) {
			if len(posSelected) == 0 {
				return
			}
			pos, ok := positionFromSelectionID(posSelected[0])
			if !ok {
				return
			}
			summonGreatDruidLifeSeedAtPosition(ctx, pos)
		})
	})
	return nil
}

func summonGreatDruidLifeSeedAtPosition(ctx *EffectContext, pos Position) {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	seedCard := getCardDB()["1401001"]
	if !pos.Valid() || ps.Units[pos.Col][pos.Row] != nil || seedCard == nil || ctx.Target == nil {
		return
	}
	seed := NewCardInstance(seedCard, ctx.PlayerID, ctx.Engine.State.TurnNumber)
	seed.Position = &Position{Col: pos.Col, Row: pos.Row}
	seed.EnterTurn = ctx.Engine.State.TurnNumber
	seed.IsHorizontal = true
	if bonusLife := ctx.Target.CurrentLife - ctx.Target.Card.Life; bonusLife > 0 {
		seed.CurrentLife += bonusLife
	}
	for elem, amount := range ctx.Target.ElementsGainBonus {
		ctx.Engine.addElementsGainBonus(seed, ctx.PlayerID, elem, amount, ctx.Source)
	}
	if len(ctx.Target.ElementsGainSet) > 0 {
		seed.ElementsGainSet = copyElementCost(ctx.Target.ElementsGainSet)
	}
	ps.Units[pos.Col][pos.Row] = seed
	ctx.Engine.triggerEffects(TriggerOnEnter, seed, nil, nil)
	ctx.Engine.triggerFieldEffectsWithData(TriggerOnUnitEnter, ctx.PlayerID, seed, map[string]any{"entered_player": ctx.PlayerID})
	ctx.Engine.triggerFieldEffectsWithData(TriggerOnUnitEnter, ctx.OpponentID, seed, map[string]any{"entered_player": ctx.PlayerID})
}

type Card1411002KnowledgeTreeDeepRoot struct{ AlwaysActive }

func (Card1411002KnowledgeTreeDeepRoot) ID() string   { return "1411002" }
func (Card1411002KnowledgeTreeDeepRoot) Name() string { return "\"知识古树\" 深耕" }
func (Card1411002KnowledgeTreeDeepRoot) OnEnter(ctx *EffectContext) error {
	ctx.Engine.advanceAllMasteryToMax(ctx.Engine.State.Players[ctx.PlayerID])
	return nil
}

type Card1411003SandWitchSommer struct{ AlwaysActive }

func (Card1411003SandWitchSommer) ID() string   { return "1411003" }
func (Card1411003SandWitchSommer) Name() string { return "沙之魔巫 梭默" }
func (Card1411003SandWitchSommer) ModifySpellArea(ctx *EffectContext, area *SpellArea) {
	if ctx.Source != nil && ctx.Source.Card.Category == model.ElementEarth && !isSorcerySkill(ctx.Source.Card) && *area == SpellAreaSingle {
		*area = SpellAreaSquare
	}
}

type Card1421003GrowingTreant struct{ AlwaysActive }

func (Card1421003GrowingTreant) ID() string      { return "1421003" }
func (Card1421003GrowingTreant) Name() string    { return "成长的树人" }
func (Card1421003GrowingTreant) MasteryMax() int { return 4 }
func (Card1421003GrowingTreant) OnMastery(ctx *EffectContext, level int) error {
	if level != 2 && level != 4 {
		return nil
	}
	loadID := fmt.Sprintf("%s:load:%d", ctx.Source.InstanceID, level)
	lifeID := fmt.Sprintf("%s:life:%d", ctx.Source.InstanceID, level)
	choices := []map[string]any{
		{"instance_id": loadID, "name": "+1地负载", "zone": "choice", "side": "own"},
		{"instance_id": lifeID, "name": "+1生命", "zone": "choice", "side": "own"},
	}
	source := ctx.Source
	ctx.Engine.SetPendingAction(ctx.PlayerID, "growing_treant_mastery_choice", "成长的树人:选择精通奖励", choices, 1, 1, func(selected []string) {
		if len(selected) == 0 || !ctx.Engine.cardStillOnField(source) {
			return
		}
		if selected[0] == lifeID {
			source.CurrentLife++
			return
		}
		if selected[0] == loadID {
			ctx.Engine.addElementsGainBonus(source, ctx.PlayerID, model.ElementEarth, 1, source)
		}
	})
	return nil
}

type Card1421004ForestGuard struct{ AlwaysActive }

func (Card1421004ForestGuard) ID() string      { return "1421004" }
func (Card1421004ForestGuard) Name() string    { return "森林守卫" }
func (Card1421004ForestGuard) MasteryMax() int { return 5 }
func (Card1421004ForestGuard) OnMastery(ctx *EffectContext, level int) error {
	switch level {
	case 1:
		ctx.Source.CurrentLife++
	case 3:
		ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementEarth, 1, ctx.Source)
	case 5:
		ctx.Source.AttackBonus += 2
	}
	return nil
}

type Card1421007HighlandTitan struct{ AlwaysActive }

func (Card1421007HighlandTitan) ID() string   { return "1421007" }
func (Card1421007HighlandTitan) Name() string { return "高地泰坦" }
func (Card1421007HighlandTitan) OnDamaged(ctx *EffectContext) error {
	if (ctx.Target == nil || ctx.Target == ctx.Source) &&
		ctx.ExtraData != nil &&
		ctx.ExtraData["damage_source"] == "spell" &&
		ctx.ExtraData["boost_count"] == 0 {
		ctx.Source.CurrentLife--
	}
	return nil
}

type Card1421010PlantationGardener struct{ AlwaysActive }

func (Card1421010PlantationGardener) ID() string   { return "1421010" }
func (Card1421010PlantationGardener) Name() string { return "种植园丁" }
func (Card1421010PlantationGardener) OnTurnStart(ctx *EffectContext) error {
	ctx.Source.Statuses["地脉标记"]++
	return nil
}
func (Card1421010PlantationGardener) OnLoadGain(ctx *EffectContext) error {
	if ctx.ExtraData == nil {
		return nil
	}
	if player, ok := ctx.ExtraData["load_gain_player"].(int); ok && player == ctx.PlayerID {
		ctx.Source.Statuses["地脉标记"]++
	}
	return nil
}
func (Card1421010PlantationGardener) OnPerTurn(ctx *EffectContext) error {
	if ctx.Source.Statuses["地脉标记"] >= 2 {
		ctx.Source.Statuses["地脉标记"] -= 2
		ctx.Engine.drawCards(ctx.PlayerID, 1)
	}
	return nil
}

type Card1421011GreatElder struct{ AlwaysActive }

func (Card1421011GreatElder) ID() string      { return "1421011" }
func (Card1421011GreatElder) Name() string    { return "大长老" }
func (Card1421011GreatElder) MasteryMax() int { return 3 }
func (Card1421011GreatElder) OnMastery(ctx *EffectContext, level int) error {
	if level == 1 || level == 3 {
		ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
			Type:          TempModNextEarthSkillLearnCostMinus,
			Amount:        2,
			RemainingUses: 1,
		})
	}
	return nil
}

type Card1511001WhiteRobeSage struct{ AlwaysActive }

func (Card1511001WhiteRobeSage) ID() string   { return "1511001" }
func (Card1511001WhiteRobeSage) Name() string { return "白袍大贤者 掌号使" }
func (Card1511001WhiteRobeSage) OnUltimate(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ps.FindEmptyPosition() == nil {
		return nil
	}
	targets := ctx.Engine.enemyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		if card == nil || card.Position == nil || !card.Card.IsCompanion() {
			return false
		}
		if !ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, cardHasPierce(ctx.Source)) {
			return false
		}
		return ctx.Engine.canPayCost(ps, ctx.Engine.effectiveCardPlayCost(ps, card))
	})
	if len(targets) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "white_robe_sage_control",
		"白袍大贤者:选择法力范围内1个可支付费用的敌方伙伴获得控制权", targets, 1, 1,
		func(selected []string) {
			target := selectedUnitFromCandidates(ctx.Engine, selected, targets)
			if target == nil || target.Position == nil || !target.Card.IsCompanion() {
				return
			}
			if !ctx.Engine.IsInSpellRange(ctx.PlayerID, target.Position.Col, target.Position.Row, cardHasPierce(ctx.Source)) {
				return
			}
			cost := ctx.Engine.effectiveCardPlayCost(ps, target)
			positions := ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)
			if len(positions) == 0 {
				return
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "white_robe_sage_position", "白袍大贤者:选择获得控制权后的入场位置", positions, 1, 1,
				func(posSelected []string) {
					if len(posSelected) == 0 {
						return
					}
					pos, ok := positionFromSelectionID(posSelected[0])
					if !ok {
						return
					}
					candidate := candidateInfo(target, "unit", "enemy")
					ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "white_robe_sage_payment",
						"白袍大贤者:支付目标入场费用以获得控制权", []map[string]any{candidate}, 1, 1, cost, false,
						func(selected []string, data map[string]any) error {
							if len(selected) == 0 || selected[0] != target.InstanceID {
								return fmt.Errorf("invalid control target")
							}
							return resolveWhiteRobeSageControl(ctx, target, pos, cost, data)
						})
				})
		})
	return nil
}

func resolveWhiteRobeSageControl(ctx *EffectContext, target *CardInstance, pos Position, cost map[string]int, data map[string]any) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if target == nil || target.Position == nil || !target.Card.IsCompanion() || target.OwnerID != ctx.OpponentID {
		return fmt.Errorf("invalid control target")
	}
	if !ctx.Engine.IsInSpellRange(ctx.PlayerID, target.Position.Col, target.Position.Row, cardHasPierce(ctx.Source)) {
		return fmt.Errorf("target out of range")
	}
	if !pos.Valid() || ps.Units[pos.Col][pos.Row] != nil {
		return fmt.Errorf("no empty position")
	}
	if !ctx.Engine.payCostForAction(ps, cost, ActionMessage{Data: data}) {
		return fmt.Errorf("invalid payment")
	}
	op := ctx.Engine.State.Players[ctx.OpponentID]
	op.Units[target.Position.Col][target.Position.Row] = nil
	target.OwnerID = ctx.PlayerID
	target.Position = &Position{Col: pos.Col, Row: pos.Row}
	ps.Units[pos.Col][pos.Row] = target
	return nil
}

type Card1521007RainbowAngel struct{ AlwaysActive }

func (Card1521007RainbowAngel) ID() string   { return "1521007" }
func (Card1521007RainbowAngel) Name() string { return "虹之天使" }

type Card1611002BlackRobeExecutor struct{ AlwaysActive }

func (Card1611002BlackRobeExecutor) ID() string   { return "1611002" }
func (Card1611002BlackRobeExecutor) Name() string { return "黑袍执行官 无心" }
func (Card1611002BlackRobeExecutor) OnFriendlyDeath(ctx *EffectContext) error {
	cause, _ := ctx.ExtraData["death_cause"].(string)
	if ctx.Target != nil && ctx.Target.Card.IsCompanion() && (cause == DeathCauseSacrifice || cause == DeathCauseDevour) {
		ctx.Source.Statuses["暗影标记"] += max(ctx.Target.Card.Life, 1)
	}
	return nil
}
func (Card1611002BlackRobeExecutor) OnUltimate(ctx *EffectContext) error {
	targets := ctx.Engine.enemyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card.Card.IsCompanion() && card.Position != nil &&
			ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, cardHasPierce(ctx.Source)) &&
			ctx.Source.Statuses["暗影标记"] >= max(card.CurrentLife, 1)
	})
	if len(targets) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "black_robe_executor_destroy",
		"黑袍执行官:选择法力范围内1个可支付暗影标记的敌方伙伴消灭", targets, 1, 1,
		func(selected []string) {
			target := selectedUnitFromCandidates(ctx.Engine, selected, targets)
			if target == nil || target.Position == nil || !target.Card.IsCompanion() {
				return
			}
			cost := max(target.CurrentLife, 1)
			if ctx.Source.Statuses["暗影标记"] < cost {
				return
			}
			ctx.Source.Statuses["暗影标记"] -= cost
			ctx.Engine.destroyUnit(target, target.OwnerID)
		})
	return nil
}

type Card1611003HeartPiercer struct{ AlwaysActive }

func (Card1611003HeartPiercer) ID() string   { return "1611003" }
func (Card1611003HeartPiercer) Name() string { return "\"穿心人\"" }
func (Card1611003HeartPiercer) OnEnter(ctx *EffectContext) error {
	addGeneratedCardToHand(ctx, "2601001")
	return nil
}

type Card1621013WordSpirit struct{ AlwaysActive }

func (Card1621013WordSpirit) ID() string   { return "1621013" }
func (Card1621013WordSpirit) Name() string { return "言灵" }
func (Card1621013WordSpirit) OnSpellCast(ctx *EffectContext) error {
	if !isEnemySpellCast(ctx) {
		return nil
	}
	for _, skill := range ctx.Engine.State.Players[ctx.OpponentID].Skills {
		if skill != nil && skill.IsHorizontal {
			ctx.Engine.addStatus(skill, StatusWeaken, 1)
		}
	}
	return nil
}

type Card2011001ArchmageStaff struct{ AlwaysActive }

func (Card2011001ArchmageStaff) ID() string   { return "2011001" }
func (Card2011001ArchmageStaff) Name() string { return "大法师之杖" }

const archmageStaffStoredSkillStatus = "archmage_staff_stored_skill"

func (Card2011001ArchmageStaff) OnEnter(ctx *EffectContext) error {
	staffPlayer := ctx.Engine.State.Players[ctx.PlayerID]
	candidates := make([]map[string]any, 0, len(staffPlayer.SkillPool))
	for _, skill := range staffPlayer.SkillPool {
		if skill == nil || !isSpellLikeCard(skill.Card) {
			continue
		}
		candidates = append(candidates, candidateInfo(skill, "skill_pool", "own"))
	}
	if len(candidates) == 0 {
		return nil
	}
	ctx.Source.Statuses["存储技能"] = 1
	ctx.Engine.SetPendingAction(ctx.PlayerID, "archmage_staff_store_skill", "大法师之杖:选择技能池中的1个法术置于此卡上", candidates, 1, 1, func(selected []string) {
		selectedID := firstSelected(selected)
		for i, skill := range staffPlayer.SkillPool {
			if skill == nil || skill.InstanceID != selectedID || !isSpellLikeCard(skill.Card) {
				continue
			}
			staffPlayer.SkillPool = append(staffPlayer.SkillPool[:i], staffPlayer.SkillPool[i+1:]...)
			skill.SlotIndex = -1
			skill.IsHorizontal = true
			skill.Statuses[archmageStaffStoredSkillStatus] = 1
			ctx.Source.BoundSkills = append(ctx.Source.BoundSkills, skill)
			ctx.Source.Statuses["存储技能"] = 1
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source": cardToInfo(ctx.Source),
				"effect": "store_skill",
				"card":   cardToInfo(skill),
			}})
			return
		}
	})
	return nil
}

type Card2011002OverlordCrown struct{ AlwaysActive }

func (Card2011002OverlordCrown) ID() string   { return "2011002" }
func (Card2011002OverlordCrown) Name() string { return "统御者之冠" }
func (Card2011002OverlordCrown) OnUnitEnter(ctx *EffectContext) error {
	if ctx.ExtraData == nil || ctx.Target == nil || !ctx.Target.Card.IsCompanion() {
		return nil
	}
	if enteredPlayer, ok := ctx.ExtraData["entered_player"].(int); ok && enteredPlayer == ctx.PlayerID && ctx.Target != ctx.Source {
		setElementsGain(ctx.Target, map[string]int{})
	}
	return nil
}

type Card2011003KingRobe struct{ AlwaysActive }

func (Card2011003KingRobe) ID() string   { return "2011003" }
func (Card2011003KingRobe) Name() string { return "君王法袍 至贤" }

const kingRobeReductionStatus = "君王法袍绝技减攻"

func (Card2011003KingRobe) OnUltimate(ctx *EffectContext) error {
	diff := totalFieldLoad(ctx.Engine.State.Players[ctx.PlayerID]) - totalFieldLoad(ctx.Engine.State.Players[ctx.OpponentID])
	amount := diff / 2
	if amount <= 0 {
		return nil
	}
	ctx.Source.Statuses[kingRobeReductionStatus] = amount
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(ctx.Source),
		"effect": "enemy_spell_damage_modifier",
		"amount": -amount,
	}})
	return nil
}

func (Card2011003KingRobe) ModifyEnemySpellStats(ctx *EffectContext, stats *SpellStats) {
	amount := ctx.Source.Statuses[kingRobeReductionStatus]
	if amount > 0 {
		stats.DamageBonus -= amount
	}
}

func (Card2011003KingRobe) OnTurnEnd(ctx *EffectContext) error {
	delete(ctx.Source.Statuses, kingRobeReductionStatus)
	return nil
}

type Card2021002MemoryNecklace struct{ AlwaysActive }

func (Card2021002MemoryNecklace) ID() string   { return "2021002" }
func (Card2021002MemoryNecklace) Name() string { return "记忆项链" }
func (Card2021002MemoryNecklace) OnEquip(ctx *EffectContext) error {
	ctx.Source.Statuses["技能槽位+1"] = 1
	return nil
}

type Card2021012SketchScroll struct{ AlwaysActive }

func (Card2021012SketchScroll) ID() string   { return "2021012" }
func (Card2021012SketchScroll) Name() string { return "速写卷轴" }
func (Card2021012SketchScroll) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlySkills(ctx.PlayerID, func(skill *CardInstance) bool {
		return canUseSkillForPurpose(skill.Card, skillPurposeAttack) && !isSorcerySkill(skill.Card)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "sketch_scroll_skill",
		"选择1个已学习法术释放，本次无需消耗该技能", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			ctx.Engine.resolveSketchScrollSkill(ctx.PlayerID, selected[0])
		})
	return nil
}

type Card2021015ManaBoosterC struct{ AlwaysActive }

func (Card2021015ManaBoosterC) ID() string   { return "2021015" }
func (Card2021015ManaBoosterC) Name() string { return "法力增强剂C型" }
func (Card2021015ManaBoosterC) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModCurrentTurnSkillCostZero,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		RemainingUses:    99,
		ExpiresTurn:      ctx.Engine.State.TurnNumber + 1,
	})
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModSkillUseCooldownAdd,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		Amount:           2,
		ExpiresTurn:      ctx.Engine.State.TurnNumber + 1,
	})
	return nil
}

type Card2021017TravelPack struct{ AlwaysActive }

func (Card2021017TravelPack) ID() string   { return "2021017" }
func (Card2021017TravelPack) Name() string { return "旅行行囊" }
func (Card2021017TravelPack) OnEquip(ctx *EffectContext) error {
	ctx.Source.Statuses["道具槽位+3"] = 1
	return nil
}

type Card2021018ArcaneRune struct{ AlwaysActive }

func (Card2021018ArcaneRune) ID() string   { return "2021018" }
func (Card2021018ArcaneRune) Name() string { return "奥术符文" }
func (Card2021018ArcaneRune) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlySkillsIncludingBound(ctx.PlayerID, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "arcane_rune_skill", "奥术符文:选择己方1个法术获得+3威", candidates, 1, 1, func(selected []string) {
		skill := ctx.Engine.findSkill(ctx.Engine.State.Players[ctx.PlayerID], firstSelected(selected))
		if skill != nil {
			skill.PowerBonus += 3
			ctx.Engine.refreshPendingSpellPowerForModifiedSkill(ctx.PlayerID, skill)
		}
	})
	return nil
}

type Card2021022CounterRune struct{ AlwaysActive }

func (Card2021022CounterRune) ID() string   { return "2021022" }
func (Card2021022CounterRune) Name() string { return "反制符文" }
func (Card2021022CounterRune) OnUseItem(ctx *EffectContext) error {
	emitBatchEffect(ctx, "counter_rune_ready")
	return nil
}

type Card2111001FireDragonHeart struct{ AlwaysActive }

func (Card2111001FireDragonHeart) ID() string   { return "2111001" }
func (Card2111001FireDragonHeart) Name() string { return "火龙之心" }
func (Card2111001FireDragonHeart) OnPerTurn(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	spend := min(ps.Elements[model.ElementFire], 3)
	if spend > 0 {
		ps.Elements[model.ElementFire] -= spend
		ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{Type: TempModSkillPowerBonus, Amount: spend * 3, RemainingUses: 1, ExpiresTurn: ctx.Engine.State.TurnNumber + 1})
	}
	return nil
}

type Card2111002NurEye struct{ AlwaysActive }

func (Card2111002NurEye) ID() string            { return "2111002" }
func (Card2111002NurEye) Name() string          { return "努尔之眼" }
func (Card2111002NurEye) IsPrayerAbility() bool { return true }

const nurEyeFireMark = "nur_eye_fire_mark"

func (Card2111002NurEye) OnDamaged(ctx *EffectContext) error {
	if ctx.ExtraData != nil && (ctx.ExtraData["damage_element"] == model.ElementFire || ctx.ExtraData["status_damage"] == StatusBurn) {
		ctx.Source.Statuses[nurEyeFireMark]++
		ctx.Source.Statuses["火焰标记"]++
		return nil
	}
	if ctx.ExtraData == nil || (ctx.ExtraData["damage_element"] != model.ElementFire && ctx.ExtraData["status_damage"] != StatusBurn) {
		return nil
	}
	ctx.Source.Statuses["火焰标记"]++
	return nil
}
func (Card2111002NurEye) OnPerTurn(ctx *EffectContext) error {
	newMarkers := ctx.Source.Statuses[nurEyeFireMark]
	ctx.Source.Statuses[nurEyeFireMark] = 0
	ctx.Source.Statuses["火焰标记"] = 0
	switch newMarkers {
	case 0:
		ctx.Engine.discardFriendlyCandidate(ctx.PlayerID, ctx.Source.InstanceID)
		return nil
	case 1:
		ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementFire: 2})
	case 2:
		ctx.Engine.addNextElementSpellPowerBonus(ctx.PlayerID, model.ElementFire, 2)
	case 3:
		ctx.Engine.addNextElementSpellDamageBonus(ctx.PlayerID, model.ElementFire, 1)
	default:
		candidates := append(ctx.Engine.friendlyUnits(ctx.PlayerID, true, nil), ctx.Engine.enemyUnits(ctx.PlayerID, true, nil)...)
		ctx.Engine.SetPendingAction(ctx.PlayerID, "nur_eye_fire_damage", "努尔之眼:选择1个单位造成2点火焰伤害", candidates, 1, 1, func(selected []string) {
			for _, ps := range ctx.Engine.State.Players {
				target := ctx.Engine.findUnitOnGrid(ps, firstSelected(selected))
				if target != nil {
					ctx.Engine.dealDamageWithExtra(target, 2, target.OwnerID, map[string]any{
						"damage_source":  "effect",
						"damage_element": model.ElementFire,
					})
					return
				}
			}
		})
	}
	return nil
}

type Card2211002WinterBow struct{ AlwaysActive }

func (Card2211002WinterBow) ID() string   { return "2211002" }
func (Card2211002WinterBow) Name() string { return "嗜魔弓 凛冬" }

const winterBowWaterMark = "winter_bow_water_mark"

func (Card2211002WinterBow) OnEnter(ctx *EffectContext) error {
	bindSkillToHost(ctx, "3201002")
	return nil
}
func (Card2211002WinterBow) OnSpellCast(ctx *EffectContext) error {
	winterBowPlayer := ctx.Engine.State.Players[ctx.PlayerID]
	if ctx.ExtraData == nil {
		return nil
	}
	if _, ok := ctx.ExtraData["cast_player"].(int); !ok {
		return nil
	}
	if !ctx.Engine.canPayCost(winterBowPlayer, map[string]int{model.ElementWater: 1}) {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "winter_bow_water_mark", "嗜魔弓 凛冬:是否支付1水放置1个水纹标记物", []map[string]any{candidateInfo(ctx.Source, "equipment", "own")}, 0, 1, func(selected []string) {
		if len(selected) == 0 {
			return
		}
		if ctx.Engine.payCostForAction(winterBowPlayer, map[string]int{model.ElementWater: 1}, ActionMessage{}) {
			ctx.Source.Statuses[winterBowWaterMark]++
		}
	})
	return nil
}

type Card2221001FrostHeart struct{ AlwaysActive }

func (Card2221001FrostHeart) ID() string   { return "2221001" }
func (Card2221001FrostHeart) Name() string { return "冰霜之心" }
func (Card2221001FrostHeart) CanReactToSpell(ctx *EffectContext, spell *SpellCast) bool {
	return ctx != nil && spell != nil && spell.AttackerID != ctx.PlayerID
}
func (Card2221001FrostHeart) OnSpellReaction(ctx *EffectContext, spell *SpellCast) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:          TempModAllSpellDamageZero,
		RemainingUses: 1,
		ExpiresTurn:   ctx.Engine.State.TurnNumber + 1,
	})
	ctx.Engine.discardFriendlyCandidate(ctx.PlayerID, ctx.Source.InstanceID)
	ctx.Engine.emit(GameEvent{Type: "spell_reaction", Player: -1, Data: map[string]any{
		"player": ctx.PlayerID,
		"card":   cardToInfo(ctx.Source),
		"effect": "damage_zero",
	}})
	return nil
}

type Card2221010TideRune struct{ AlwaysActive }

func (Card2221010TideRune) ID() string   { return "2221010" }
func (Card2221010TideRune) Name() string { return "潮涌符文" }
func (Card2221010TideRune) OnUseItem(ctx *EffectContext) error {
	targets := ctx.Engine.friendlyUnits(ctx.PlayerID, false, isWaterCompanion)
	if len(targets) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "tide_rune_buff",
		"潮涌符文:选择你的1个水纹伙伴获得负载+2水", targets, 1, 1,
		func(selected []string) {
			target := selectedUnitFromCandidates(ctx.Engine, selected, targets)
			if target != nil {
				ctx.Engine.addElementsGainBonus(target, ctx.PlayerID, model.ElementWater, 2, ctx.Source)
			}
		})
	return nil
}

type Card2221011RainOfGrace struct{ AlwaysActive }

func (Card2221011RainOfGrace) ID() string   { return "2221011" }
func (Card2221011RainOfGrace) Name() string { return "恩惠之雨" }
func (Card2221011RainOfGrace) OnUseItem(ctx *EffectContext) error {
	for _, unit := range ctx.Engine.getAllFieldCards(ctx.Engine.State.Players[ctx.PlayerID]) {
		healUnit(unit, 2)
	}
	return nil
}

type Card2311001ThunderSource struct{ AlwaysActive }

func (Card2311001ThunderSource) ID() string   { return "2311001" }
func (Card2311001ThunderSource) Name() string { return "雷之源" }
func (Card2311001ThunderSource) ModifyCardPlayCost(ctx *EffectContext, card *CardInstance, cost map[string]int) {
	reduceCost(cost, model.ElementAir, 1)
}
func (Card2311001ThunderSource) ModifySkillUseCost(ctx *EffectContext, cost map[string]int) {
	reduceCost(cost, model.ElementAir, 1)
}

type Card2311002ThunderDrum struct{ AlwaysActive }

func (Card2311002ThunderDrum) ID() string   { return "2311002" }
func (Card2311002ThunderDrum) Name() string { return "唤雷震鼓" }
func (Card2311002ThunderDrum) OnDraw(ctx *EffectContext) error {
	if ctx.ExtraData == nil {
		return nil
	}
	drawnPlayer, _ := ctx.ExtraData["drawn_player"].(int)
	if drawnPlayer != ctx.PlayerID {
		return nil
	}
	candidates := []map[string]any{candidateInfo(ctx.Source, "equipment", "own")}
	if drawn, ok := ctx.ExtraData["drawn_card"].(*CardInstance); ok && drawn != nil {
		info := cardToInfo(drawn)
		info["zone"] = "hand"
		info["side"] = "own"
		candidates = append(candidates, info)
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "thunder_drum_mark",
		"唤雷震鼓:是否展示抽到的牌并放置1个标记?", candidates, 0, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			ctx.Source.Statuses["雷鼓标记"]++
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source": cardToInfo(ctx.Source),
				"effect": "thunder_drum_mark",
				"amount": 1,
			}})
		})
	return nil
}

func (Card2311002ThunderDrum) OnPerTurn(ctx *EffectContext) error {
	if thunderDrumMarks(ctx.Source) < 3 || ctx.Engine.State.PendingAction != nil {
		return nil
	}
	choices := []map[string]any{
		{"instance_id": "attack", "name": "本回合你的大气法术+1攻", "zone": "choice"},
		{"instance_id": "stun", "name": "本回合你的大气法术获得晕眩1", "zone": "choice"},
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "thunder_drum_bonus",
		"唤雷震鼓:移除3个标记,选择本回合大气法术获得的效果", choices, 1, 1,
		func(selected []string) {
			if len(selected) == 0 || thunderDrumMarks(ctx.Source) < 3 {
				return
			}
			spendThunderDrumMarks(ctx.Source, 3)
			switch selected[0] {
			case "attack":
				ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{Type: TempModCurrentTurnElementDamage, Element: model.ElementAir, Amount: 1, ExpiresTurn: ctx.Engine.State.TurnNumber})
			case "stun":
				ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{Type: TempModCurrentTurnElementHitStatus, Element: model.ElementAir, Status: StatusStun, Amount: 1, ExpiresTurn: ctx.Engine.State.TurnNumber})
			}
		})
	return nil
}

func thunderDrumMarks(source *CardInstance) int {
	if source == nil {
		return 0
	}
	return source.Statuses["雷鼓标记"]
}

func spendThunderDrumMarks(source *CardInstance, amount int) {
	if source == nil || amount <= 0 {
		return
	}
	spend := min(source.Statuses["雷鼓标记"], amount)
	source.Statuses["雷鼓标记"] -= spend
}

type Card2321001WindbreathCompass struct{ AlwaysActive }

func (Card2321001WindbreathCompass) ID() string   { return "2321001" }
func (Card2321001WindbreathCompass) Name() string { return "风息罗盘" }

const (
	windbreathCompassPendingStatus      = "风息罗盘待触发"
	windbreathCompassTemporaryAirStatus = "风息罗盘临时气负载"
)

func (Card2321001WindbreathCompass) OnDraw(ctx *EffectContext) error {
	if ctx.ExtraData == nil {
		return nil
	}
	if player, ok := ctx.ExtraData["drawn_player"].(int); !ok || player != ctx.PlayerID {
		return nil
	}
	ctx.Source.Statuses[windbreathCompassPendingStatus]++
	openWindbreathCompassPrompt(ctx.Engine, ctx.PlayerID, ctx.Source)
	return nil
}

func openWindbreathCompassPrompt(e *Engine, playerID int, source *CardInstance) {
	if e == nil || source == nil || source.Statuses[windbreathCompassPendingStatus] <= 0 {
		return
	}
	if e.State.PendingAction != nil && e.State.PendingAction.Type == "windbreath_compass" && e.State.PendingAction.PlayerID == playerID {
		return
	}
	candidates := []map[string]any{candidateInfo(source, "equipment", "own")}
	e.SetPendingAction(playerID, "windbreath_compass",
		"你抽牌了，是否触发风息罗盘获得临时负载+1气？", candidates, 0, 1,
		func(selected []string) {
			if source.Statuses[windbreathCompassPendingStatus] > 0 {
				source.Statuses[windbreathCompassPendingStatus]--
			}
			if len(selected) > 0 && selected[0] == source.InstanceID {
				e.addElementsGainBonus(source, playerID, model.ElementAir, 1, source)
				source.Statuses[windbreathCompassTemporaryAirStatus]++
			}
			openWindbreathCompassPrompt(e, playerID, source)
		})
}

func (Card2321001WindbreathCompass) OnTurnEnd(ctx *EffectContext) error {
	count := ctx.Source.Statuses[windbreathCompassTemporaryAirStatus]
	if count > 0 {
		addElementsGainBonus(ctx.Source, model.ElementAir, -count)
		delete(ctx.Source.Statuses, windbreathCompassTemporaryAirStatus)
	}
	delete(ctx.Source.Statuses, windbreathCompassPendingStatus)
	return nil
}

type Card2321010IllusionScroll struct{ AlwaysActive }

func (Card2321010IllusionScroll) ID() string   { return "2321010" }
func (Card2321010IllusionScroll) Name() string { return "幻术卷轴" }
func (Card2321010IllusionScroll) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.startIllusionScrollRearrange(ctx.PlayerID, ctx.Source, ctx.ExtraData)
	return nil
}

type Card2321011TeleportRune struct{ AlwaysActive }

func (Card2321011TeleportRune) ID() string   { return "2321011" }
func (Card2321011TeleportRune) Name() string { return "传送符文" }
func (Card2321011TeleportRune) OnUseItem(ctx *EffectContext) error {
	target := ctx.Target
	if target == nil || target.Card == nil || !target.Card.IsCompanion() || target.Position == nil {
		return nil
	}
	positions := ctx.Engine.emptyUnitPositionsForPlayer(target.OwnerID, ctx.PlayerID)
	if len(positions) == 0 {
		return nil
	}
	targetID := target.InstanceID
	targetOwner := target.OwnerID
	ctx.Engine.SetPendingAction(ctx.PlayerID, "teleport_rune_position",
		"Teleport Rune: choose another empty position", positions, 1, 1,
		func(selected []string) {
			pos, ok := positionFromSelectionID(firstSelected(selected))
			if !ok {
				return
			}
			ctx.Engine.moveUnitToPosition(targetOwner, targetID, pos)
		})
	return nil
}

type Card2321012WindCloak struct{ AlwaysActive }

func (Card2321012WindCloak) ID() string   { return "2321012" }
func (Card2321012WindCloak) Name() string { return "随风斗篷" }
func (Card2321012WindCloak) OnUltimate(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	pos := ps.FindEmptyPosition()
	if pos != nil && ps.Hero != nil && ps.Hero.Position != nil {
		ps.Units[ps.Hero.Position.Col][ps.Hero.Position.Row] = nil
		ps.Hero.Position = pos
		ps.Units[pos.Col][pos.Row] = ps.Hero
	}
	return nil
}
