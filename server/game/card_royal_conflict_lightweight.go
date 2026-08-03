package game

import (
	"fmt"

	"eraofarcane/model"
)

type Card4311101SkyWitchSoland struct{ AlwaysActive }

func (Card4311101SkyWitchSoland) ID() string   { return "4311101" }
func (Card4311101SkyWitchSoland) Name() string { return "司天魔巫 索兰德" }
func (Card4311101SkyWitchSoland) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Target == nil || ctx.Target.Card == nil || !isSpellLikeCard(ctx.Target.Card) {
		return
	}
	if hasCardTag(ctx.Target.Card, "驱动") || hasCardTag(ctx.Target.Card, "聚能") {
		stats.PowerBonus++
	}
}
func (Card4311101SkyWitchSoland) ValidateSkillLearn(ctx *EffectContext, skill *CardInstance) error {
	if skill == nil || skill.Card == nil {
		return nil
	}
	if hasCardTag(skill.Card, "驱动") || hasCardTag(skill.Card, "聚能") {
		return nil
	}
	return fmt.Errorf("司天魔巫 索兰德只能让你学习驱动或聚能法术")
}

type Card2221110LingeringFrostScroll struct{ AlwaysActive }

func (Card2221110LingeringFrostScroll) ID() string   { return "2221110" }
func (Card2221110LingeringFrostScroll) Name() string { return "残霜飞雪卷轴" }
func (Card2221110LingeringFrostScroll) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.Source.Card.Number != "2221110" {
		return
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ps == nil {
		return
	}
	stats.PowerBonus += ps.SpellsCastThisTurn[model.ElementWater] * 3
}

type Card2521112OracleScrollUnity struct{ AlwaysActive }

func (Card2521112OracleScrollUnity) ID() string   { return "2521112" }
func (Card2521112OracleScrollUnity) Name() string { return "神谕卷轴 团结" }
func (Card2521112OracleScrollUnity) ModifySelfCardPlayCost(ctx *EffectContext, cost map[string]int) {
	if ctx == nil || ctx.Engine == nil {
		return
	}
	reduceCost(cost, model.ElementLight, countFriendlyLightUnits(ctx.Engine, ctx.PlayerID))
}

type Card3521101Gospel struct{ AlwaysActive }

func (Card3521101Gospel) ID() string   { return "3521101" }
func (Card3521101Gospel) Name() string { return "福音" }
func (Card3521101Gospel) OnConsume(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.Card == nil {
		return nil
	}
	consumedPlayer, _ := ctx.ExtraData["consumed_player"].(int)
	if consumedPlayer != ctx.PlayerID || !ctx.Target.Card.IsCompanion() || ctx.Target.Card.Category != model.ElementLight {
		return nil
	}
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModCurrentTurnSkillUseCostMinus,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		TargetInstanceID: ctx.Source.InstanceID,
		Element:          model.ElementLight,
		Amount:           1,
		ExpiresTurn:      ctx.Engine.State.TurnNumber,
	})
	return nil
}

type Card3421106RottingErosion struct{ AlwaysActive }

func (Card3421106RottingErosion) ID() string      { return "3421106" }
func (Card3421106RottingErosion) Name() string    { return "腐朽侵蚀" }
func (Card3421106RottingErosion) MasteryMax() int { return 6 }
func (Card3421106RottingErosion) OnMastery(ctx *EffectContext, level int) error {
	if ctx == nil || ctx.Source == nil {
		return nil
	}
	if level == 3 || level == 6 {
		ctx.Source.PowerBonus++
		ctx.Source.AttackBonus++
	}
	return nil
}
func (Card3421106RottingErosion) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || !isOwnSpellHit(ctx) {
		return nil
	}
	amount := max(ctx.Source.Card.Attack+ctx.Source.AttackBonus, 0)
	if amount > 0 {
		for _, skill := range enemySpellInstancesIncludingBound(ctx.Engine, ctx.PlayerID) {
			if canInstanceBeWeakened(skill) {
				ctx.Engine.addStatus(skill, StatusWeaken, amount)
			}
		}
	}
	ctx.Engine.advanceMastery(ctx.Source, ctx.PlayerID, 1)
	return nil
}

func countFriendlyLightUnits(e *Engine, playerID int) int {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return 0
	}
	count := 0
	for _, card := range e.getAllFieldCards(e.State.Players[playerID]) {
		if card != nil && card.Card != nil && (card.Card.IsHero() || card.Card.IsCompanion()) && card.Card.Category == model.ElementLight && !e.hasEffectiveStatus(card, StatusPetrify) {
			count++
		}
	}
	return count
}

func enemySpellInstancesIncludingBound(e *Engine, playerID int) []*CardInstance {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	opponent := e.State.Players[1-playerID]
	if opponent == nil {
		return nil
	}
	result := make([]*CardInstance, 0)
	for _, skill := range opponent.Skills {
		if skill != nil {
			result = append(result, skill)
		}
	}
	for _, card := range e.getAllFieldCards(opponent) {
		if card == nil {
			continue
		}
		result = append(result, card.BoundSkills...)
	}
	return result
}

type Card3321107HeldBreath struct{ AlwaysActive }

func (Card3321107HeldBreath) ID() string   { return "3321107" }
func (Card3321107HeldBreath) Name() string { return "屏息凝神" }
func (Card3321107HeldBreath) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Engine == nil || ctx.Target == nil || ctx.Target.Card == nil {
		return
	}
	if ctx.Target.Card.Category != model.ElementAir || !isSpellLikeCard(ctx.Target.Card) {
		return
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ps == nil || ps.DrawCountThisTurn > 1 {
		return
	}
	stats.PowerBonus++
}

type Card2621104DevotionContract struct{ AlwaysActive }

func (Card2621104DevotionContract) ID() string   { return "2621104" }
func (Card2621104DevotionContract) Name() string { return "献身契约" }
func (Card2621104DevotionContract) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.Card == nil {
		return nil
	}
	if ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) || !isFriendlySpellCast(ctx) || !hasCardTag(ctx.Target.Card, "代赎") {
		return nil
	}
	hero := ctx.Engine.State.Players[ctx.PlayerID].Hero
	if hero != nil {
		ctx.Engine.dealDamageWithExtra(hero, 1, ctx.PlayerID, map[string]any{
			"damage_source": "devotion_contract",
			"attacker":      ctx.PlayerID,
		})
	}
	ctx.Engine.drawCards(ctx.PlayerID, 1)
	ctx.Source.UsedThisTurn++
	return nil
}

const skyCityZenithStoneMarkerStatus = "云霄城的天顶石标记物"

type Card2311101SkyCityZenithStone struct{ AlwaysActive }

func (Card2311101SkyCityZenithStone) ID() string   { return "2311101" }
func (Card2311101SkyCityZenithStone) Name() string { return "云霄城的天顶石" }
func (Card2311101SkyCityZenithStone) OnDraw(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil {
		return nil
	}
	drawnPlayer, ok := ctx.ExtraData["drawn_player"].(int)
	if !ok || drawnPlayer < 0 || drawnPlayer >= len(ctx.Engine.State.Players) {
		return nil
	}
	ctx.Source.Statuses[skyCityZenithStoneMarkerStatus]++
	if ctx.Source.Statuses[skyCityZenithStoneMarkerStatus] < 5 {
		return nil
	}
	delete(ctx.Source.Statuses, skyCityZenithStoneMarkerStatus)

	ps := ctx.Engine.State.Players[drawnPlayer]
	if ps == nil {
		return nil
	}
	frontRow := ps.GetFrontRow()
	if frontRow < 0 || frontRow >= 3 {
		return nil
	}
	targets := make([]*CardInstance, 0, 3)
	for col := 0; col < 3; col++ {
		if unit := ps.Units[col][frontRow]; unit != nil {
			targets = append(targets, unit)
		}
	}
	for _, target := range targets {
		ctx.Engine.dealDamageWithExtra(target, 1, drawnPlayer, map[string]any{
			"damage_source": "sky_city_zenith_stone",
			"attacker":      ctx.PlayerID,
		})
		ctx.Engine.addStatus(target, StatusStun, 1)
	}
	return nil
}

const bloodGuMarkerStatus = "血蛊标记物"

type Card2621103BloodGu struct{ AlwaysActive }

func (Card2621103BloodGu) ID() string   { return "2621103" }
func (Card2621103BloodGu) Name() string { return "血蛊" }
func (Card2621103BloodGu) OnDamaged(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.ExtraData == nil {
		return nil
	}
	damagedPlayer, ok := ctx.ExtraData["damaged_player"].(int)
	if !ok || damagedPlayer != ctx.PlayerID || ctx.Target != ctx.Engine.playerHeroCard(ctx.PlayerID) {
		return nil
	}
	damage, _ := ctx.ExtraData["damage"].(int)
	if damage <= 0 {
		damage = 1
	}
	ctx.Source.Statuses[bloodGuMarkerStatus] = min(6, ctx.Source.Statuses[bloodGuMarkerStatus]+damage)
	return nil
}
func (Card2621103BloodGu) PerTurnLabel(*CardInstance) string {
	return "主动"
}
func (Card2621103BloodGu) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	amount := ctx.Source.Statuses[bloodGuMarkerStatus] / 2
	if amount <= 0 {
		return fmt.Errorf("血蛊没有足够标记物")
	}
	if !ctx.Engine.sacrificeEquipment(ctx.PlayerID, ctx.Source.InstanceID) {
		return fmt.Errorf("血蛊必须从装备区献祭")
	}
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModSkillPowerBonus,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		Amount:           amount,
		ExpiresTurn:      ctx.Engine.State.TurnNumber,
	})
	return nil
}

type Card2521108CouncilJudgmentHammer struct{ AlwaysActive }

func (Card2521108CouncilJudgmentHammer) ID() string   { return "2521108" }
func (Card2521108CouncilJudgmentHammer) Name() string { return "议庭审判锤" }
func (Card2521108CouncilJudgmentHammer) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.Card == nil {
		return nil
	}
	if ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) || !isEnemySpellCast(ctx) || isSorcerySkill(ctx.Target.Card) {
		return nil
	}
	addGeneratedCardsToPlayerDeck(ctx, ctx.OpponentID, "2001102", 3)
	ctx.Source.UsedThisTurn++
	return nil
}

type Card2521103RedAgateChalice struct{ AlwaysActive }

func (Card2521103RedAgateChalice) ID() string   { return "2521103" }
func (Card2521103RedAgateChalice) Name() string { return "红玛瑙圣杯" }
func (Card2521103RedAgateChalice) ModifyElementsGain(ctx *EffectContext, target *CardInstance, gains map[string]int) {
	if ctx == nil || ctx.Engine == nil || target == nil || target.Card == nil || !redAgateChaliceSetComplete(ctx.Engine, ctx.PlayerID) {
		return
	}
	switch target.Card.Number {
	case "2521006", "2521007", "2521103":
		gains[model.ElementLight]++
	}
}

func redAgateChaliceSetComplete(e *Engine, playerID int) bool {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return false
	}
	seen := map[string]bool{}
	for _, card := range e.State.Players[playerID].Equipment {
		if card == nil || card.Card == nil {
			continue
		}
		switch card.Card.Number {
		case "2521006", "2521007", "2521103":
			seen[card.Card.Number] = true
		}
	}
	return seen["2521006"] && seen["2521007"] && seen["2521103"]
}

type Card2221109QuickIceBullet struct{ AlwaysActive }

func (Card2221109QuickIceBullet) ID() string   { return "2221109" }
func (Card2221109QuickIceBullet) Name() string { return "速射冰弹" }
func (Card2221109QuickIceBullet) OnUseItem(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModNextItemOrSkillCostMinus,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		Element:          model.ElementWater,
		Amount:           3,
		RemainingUses:    1,
		ExpiresTurn:      ctx.Engine.State.TurnNumber,
	})
	return nil
}

type Card3511102LastStandLight struct{ AlwaysActive }

func (Card3511102LastStandLight) ID() string   { return "3511102" }
func (Card3511102LastStandLight) Name() string { return "绝境之光 孤星闪耀" }
func (Card3511102LastStandLight) ValidateSkillUse(ctx *EffectContext, skill *CardInstance, purpose skillPurpose) error {
	if ctx == nil || ctx.Engine == nil || skill == nil || skill.Card == nil || skill.Card.Number != "3511102" {
		return nil
	}
	if royalCompanionCount(ctx.Engine.State.Players[ctx.PlayerID]) < royalCompanionCount(ctx.Engine.State.Players[ctx.OpponentID]) {
		return nil
	}
	return fmt.Errorf("绝境之光 孤星闪耀只能在你场上单位比对方少时使用")
}
func (Card3511102LastStandLight) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.Source.Card.Number != "3511102" {
		return
	}
	stats.PowerBonus += highestFriendlyLightCompanionLifeAndLoad(ctx.Engine, ctx.PlayerID)
}

func highestFriendlyLightCompanionLifeAndLoad(e *Engine, playerID int) int {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return 0
	}
	best := 0
	ps := e.State.Players[playerID]
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := ps.Units[col][row]
			if unit == nil || unit.Card == nil || !unit.Card.IsCompanion() || unit.Card.Category != model.ElementLight {
				continue
			}
			best = max(best, unit.CurrentLife+totalElementCost(e.effectiveElementsGain(unit)))
		}
	}
	return best
}

const (
	collectorEquipTriggeredTurnStatus = "收藏家装备触发回合"
	collectorItemTriggeredTurnStatus  = "收藏家消耗品触发回合"
)

type Card1011101CollectorCoralFenlo struct{ AlwaysActive }

func (Card1011101CollectorCoralFenlo) ID() string   { return "1011101" }
func (Card1011101CollectorCoralFenlo) Name() string { return "收藏家 珊瑚 芬洛" }
func (Card1011101CollectorCoralFenlo) OnCardEnter(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.Card == nil || ctx.ExtraData == nil {
		return nil
	}
	enteredPlayer, _ := ctx.ExtraData["entered_player"].(int)
	equipped, _ := ctx.ExtraData["equipped"].(bool)
	if enteredPlayer != ctx.PlayerID || !equipped || !isEquipmentCard(ctx.Target.Card) || ctx.Source.Statuses[collectorEquipTriggeredTurnStatus] == ctx.Engine.State.TurnNumber {
		return nil
	}
	ctx.Engine.drawCards(ctx.PlayerID, 1)
	ctx.Source.Statuses[collectorEquipTriggeredTurnStatus] = ctx.Engine.State.TurnNumber
	return nil
}
func (Card1011101CollectorCoralFenlo) OnUseItem(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.Card == nil || ctx.ExtraData == nil {
		return nil
	}
	usedPlayer, _ := ctx.ExtraData["used_player"].(int)
	if usedPlayer != ctx.PlayerID || !isConsumableCardInstance(ctx.Target) || ctx.Source.Statuses[collectorItemTriggeredTurnStatus] == ctx.Engine.State.TurnNumber {
		return nil
	}
	ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementArcane: 1})
	ctx.Source.Statuses[collectorItemTriggeredTurnStatus] = ctx.Engine.State.TurnNumber
	return nil
}

type Card1521111CouncilConsul struct{ AlwaysActive }

func (Card1521111CouncilConsul) ID() string   { return "1521111" }
func (Card1521111CouncilConsul) Name() string { return "议庭执政官" }
func (Card1521111CouncilConsul) OnUnitEnter(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Target == nil || ctx.Target.Card == nil || ctx.ExtraData == nil {
		return nil
	}
	enteredPlayer, _ := ctx.ExtraData["entered_player"].(int)
	if enteredPlayer != ctx.OpponentID || !ctx.Target.Card.IsCompanion() {
		return nil
	}
	addGeneratedCardsToPlayerDeck(ctx, ctx.OpponentID, "2001102", 3)
	return nil
}

type Card4011101PureSpiritOshis struct{ AlwaysActive }

func (Card4011101PureSpiritOshis) ID() string   { return "4011101" }
func (Card4011101PureSpiritOshis) Name() string { return "纯净灵体 奥希斯" }
func (Card4011101PureSpiritOshis) OnCardEnter(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Target == nil || ctx.Target.Card == nil || ctx.ExtraData == nil {
		return nil
	}
	enteredPlayer, _ := ctx.ExtraData["entered_player"].(int)
	if enteredPlayer != ctx.PlayerID || ctx.Target.Card.Category == model.ElementArcane {
		return nil
	}
	for _, skill := range friendlySpellInstancesIncludingBound(ctx.Engine, ctx.PlayerID) {
		if canInstanceBeWeakened(skill) {
			ctx.Engine.addStatus(skill, StatusWeaken, 2)
		}
	}
	return nil
}

func friendlySpellInstancesIncludingBound(e *Engine, playerID int) []*CardInstance {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	ps := e.State.Players[playerID]
	if ps == nil {
		return nil
	}
	result := make([]*CardInstance, 0)
	for _, skill := range ps.Skills {
		if skill != nil {
			result = append(result, skill)
		}
	}
	for _, card := range e.getAllFieldCards(ps) {
		if card == nil {
			continue
		}
		for _, skill := range card.BoundSkills {
			if skill != nil {
				result = append(result, skill)
			}
		}
	}
	return result
}

type Card3611102ClawOfErebos struct{ AlwaysActive }

func (Card3611102ClawOfErebos) ID() string   { return "3611102" }
func (Card3611102ClawOfErebos) Name() string { return "厄瑞波斯之爪" }
func (Card3611102ClawOfErebos) ValidateSkillLearn(ctx *EffectContext, skill *CardInstance) error {
	if ctx == nil || ctx.Engine == nil || enemySkillSlotWeakenedSpellLayers(ctx.Engine, ctx.PlayerID) >= 3 {
		return nil
	}
	return fmt.Errorf("敌方需要有三个及以上虚弱法术才能学习厄瑞波斯之爪")
}
func (Card3611102ClawOfErebos) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.Source.Card.Number != "3611102" {
		return
	}
	stats.PowerBonus += enemySkillSlotWeakenedSpellLayers(ctx.Engine, ctx.PlayerID)
}
func (Card3611102ClawOfErebos) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.Source.Card.Number != "3611102" || !isFriendlySpellCast(ctx) {
		return nil
	}
	candidates := ctx.Engine.enemySkills(ctx.PlayerID, canInstanceBeWeakened)
	if len(candidates) == 0 {
		return nil
	}
	maxSelect := min(3, len(candidates))
	ctx.Engine.SetPendingAction(ctx.PlayerID, "claw_of_erebos_weaken",
		"厄瑞波斯之爪:使最多3个不同的敌方法术虚弱1", candidates, 0, maxSelect,
		func(selected []string) {
			seen := map[string]bool{}
			for _, id := range selected {
				if seen[id] {
					continue
				}
				seen[id] = true
				for _, skill := range ctx.Engine.State.Players[ctx.OpponentID].Skills {
					if skill != nil && skill.InstanceID == id && canInstanceBeWeakened(skill) {
						ctx.Engine.addStatus(skill, StatusWeaken, 1)
						break
					}
				}
			}
		})
	return nil
}

func enemySkillSlotWeakenedSpellLayers(e *Engine, playerID int) int {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return 0
	}
	total := 0
	opponent := e.State.Players[1-playerID]
	if opponent == nil {
		return 0
	}
	for _, skill := range opponent.Skills {
		if skill != nil && skill.Statuses[StatusWeaken] > 0 && e.hasEffectiveStatus(skill, StatusWeaken) {
			total += skill.Statuses[StatusWeaken]
		}
	}
	return total
}

type Card3621102Retribution struct{ AlwaysActive }

func (Card3621102Retribution) ID() string   { return "3621102" }
func (Card3621102Retribution) Name() string { return "报应" }
func (Card3621102Retribution) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.Source.Card.Number != "3621102" {
		return
	}
	if ctx.ExtraData["purpose"] != string(skillPurposeAttack) {
		return
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ps == nil {
		return
	}
	stats.DamageBonus += ps.HeroDamageTakenThisTurn + ps.HeroDamageTakenLastTurn
}

type Card2521107PanaceaP struct{ AlwaysActive }

func (Card2521107PanaceaP) ID() string   { return "2521107" }
func (Card2521107PanaceaP) Name() string { return "百灵药P型" }
func (Card2521107PanaceaP) OnUseItem(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil {
		return nil
	}
	resolveHealAndDraw := func() {
		targets := ctx.Engine.friendlyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
			return card != nil && card.CurrentLife < maxLife(card)
		})
		if len(targets) == 0 {
			ctx.Engine.drawCards(ctx.PlayerID, 1)
			return
		}
		ctx.Engine.SetPendingAction(ctx.PlayerID, "panacea_p_heal",
			"百灵药P型:选择1个友方单位回复1点生命", targets, 0, 1,
			func(selected []string) {
				target := selectedUnitFromCandidates(ctx.Engine, selected, targets)
				if target != nil {
					healUnit(target, 1)
					ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
						"source": cardToInfo(ctx.Source),
						"target": cardToInfo(target),
						"effect": "heal",
						"amount": 1,
					}})
				}
				ctx.Engine.drawCards(ctx.PlayerID, 1)
			})
	}

	targets := append(ctx.Engine.friendlyUnits(ctx.PlayerID, true, nil), ctx.Engine.enemyUnits(ctx.PlayerID, true, nil)...)
	if len(targets) == 0 {
		resolveHealAndDraw()
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "panacea_p_damage",
		"百灵药P型:选择1个单位造成1点伤害", targets, 1, 1,
		func(selected []string) {
			target := selectedUnitFromCandidates(ctx.Engine, selected, targets)
			if target != nil {
				ctx.Engine.dealDamageWithExtra(target, 1, target.OwnerID, map[string]any{
					"damage_source": "panacea_p",
					"attacker":      ctx.PlayerID,
				})
			}
			resolveHealAndDraw()
		})
	return nil
}

type Card2121110OfferingTorch struct{ AlwaysActive }

func (Card2121110OfferingTorch) ID() string   { return "2121110" }
func (Card2121110OfferingTorch) Name() string { return "供奉之炬" }
func (Card2121110OfferingTorch) OnUseItem(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil {
		return nil
	}
	sources := ctx.Engine.friendlySkillsIncludingBound(ctx.PlayerID, isFireSpellInstance)
	if len(sources) < 2 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "offering_torch_exile",
		"供奉之炬:选择要移出游戏的火焰法术", sources, 1, 1,
		func(selected []string) {
			sourceID := firstSelected(selected)
			exiledSkill := findFriendlySkillIncludingBound(ctx.Engine, ctx.PlayerID, sourceID)
			if !isFireSpellInstance(exiledSkill) {
				return
			}
			powerBonus := max(exiledSkill.Card.Power+exiledSkill.PowerBonus, 0)
			attackBonus := max(exiledSkill.Card.Attack+exiledSkill.AttackBonus, 0)
			if !ctx.Engine.exileCard(ctx.PlayerID, exiledSkill) {
				return
			}
			targets := ctx.Engine.friendlySkillsIncludingBound(ctx.PlayerID, func(skill *CardInstance) bool {
				return isFireSpellInstance(skill) && skill.InstanceID != sourceID
			})
			if len(targets) == 0 {
				return
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "offering_torch_buff",
				"供奉之炬:选择另一个火焰法术永久增加威和攻", targets, 1, 1,
				func(buffSelected []string) {
					targetID := firstSelected(buffSelected)
					target := findFriendlySkillIncludingBound(ctx.Engine, ctx.PlayerID, targetID)
					if !isFireSpellInstance(target) || target.InstanceID == sourceID {
						return
					}
					target.PowerBonus += powerBonus
					target.AttackBonus += attackBonus
					ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
						"source":       cardToInfo(ctx.Source),
						"exiled_skill": sourceID,
						"target":       cardToInfo(target),
						"effect":       "permanent_spell_buff",
						"power":        powerBonus,
						"attack":       attackBonus,
					}})
				})
		})
	return nil
}

type Card2121101LavafortAshes struct{ AlwaysActive }

func (Card2121101LavafortAshes) ID() string   { return "2121101" }
func (Card2121101LavafortAshes) Name() string { return "熔岩堡的灰烬" }
func (Card2121101LavafortAshes) OnUseItem(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil {
		return nil
	}
	sources := lavaFortAshSourceCandidates(ctx.Engine, ctx.PlayerID)
	if len(sources) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "lavafort_ashes_exile_fire_skill",
		"熔岩堡的灰烬:选择场上或技能池1个火焰技能移出游戏", sources, 1, 1,
		func(selected []string) {
			source := findFriendlyFieldOrPoolSkill(ctx.Engine, ctx.PlayerID, firstSelected(selected))
			if !isFireLearnableSkillInstance(source) {
				return
			}
			sourceCost := totalElementCost(source.Card.ElementsCost)
			targets := lavaFortAshDeckTargets(ctx.Engine, ctx.PlayerID, sourceCost)
			if len(targets) == 0 || !ctx.Engine.exileCard(ctx.PlayerID, source) {
				return
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "lavafort_ashes_search_fire_card",
				"熔岩堡的灰烬:翻取1张更高入场花费的火焰卡牌", targets, 1, 1,
				func(searchSelected []string) {
					card := ctx.Engine.searchDeckCardToHand(ctx.PlayerID, firstSelected(searchSelected))
					if card == nil {
						return
					}
					card.Statuses["入场费用"+model.ElementFire+"-1"]++
				})
		})
	return nil
}

func isFireSpellInstance(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.Category == model.ElementFire && isSpellLikeCard(card.Card)
}

func isFireLearnableSkillInstance(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.IsSkill() && card.Card.Category == model.ElementFire
}

func findFriendlySkillIncludingBound(e *Engine, playerID int, instanceID string) *CardInstance {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) || instanceID == "" {
		return nil
	}
	ps := e.State.Players[playerID]
	if ps == nil {
		return nil
	}
	for _, skill := range ps.Skills {
		if skill != nil && skill.InstanceID == instanceID {
			return skill
		}
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

func lavaFortAshSourceCandidates(e *Engine, playerID int) []map[string]any {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	ps := e.State.Players[playerID]
	if ps == nil {
		return nil
	}
	candidates := make([]map[string]any, 0)
	addIfValid := func(skill *CardInstance, zone string) {
		if !isFireLearnableSkillInstance(skill) {
			return
		}
		if len(lavaFortAshDeckTargets(e, playerID, totalElementCost(skill.Card.ElementsCost))) == 0 {
			return
		}
		candidates = append(candidates, candidateInfo(skill, zone, "own"))
	}
	for _, skill := range ps.Skills {
		addIfValid(skill, "skill")
	}
	for _, skill := range ps.SkillPool {
		addIfValid(skill, "skill_pool")
	}
	return candidates
}

func lavaFortAshDeckTargets(e *Engine, playerID int, sourceCost int) []map[string]any {
	return e.friendlyDeckCards(playerID, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.Category == model.ElementFire && totalElementCost(card.Card.ElementsCost) > sourceCost
	})
}

func findFriendlyFieldOrPoolSkill(e *Engine, playerID int, instanceID string) *CardInstance {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) || instanceID == "" {
		return nil
	}
	ps := e.State.Players[playerID]
	if ps == nil {
		return nil
	}
	for _, skill := range ps.Skills {
		if skill != nil && skill.InstanceID == instanceID {
			return skill
		}
	}
	for _, skill := range ps.SkillPool {
		if skill != nil && skill.InstanceID == instanceID {
			return skill
		}
	}
	return nil
}

var _ SpellStatModifier = Card4311101SkyWitchSoland{}
var _ SkillLearnPermissionModifier = Card4311101SkyWitchSoland{}
var _ SkillContributionModifier = Card2221110LingeringFrostScroll{}
var _ SelfCardPlayCostModifier = Card2521112OracleScrollUnity{}
var _ SpellStatModifier = Card3321107HeldBreath{}
var _ OnConsumeBehavior = Card3521101Gospel{}
var _ MasteryBehavior = Card3421106RottingErosion{}
var _ OnSpellHitBehavior = Card3421106RottingErosion{}
var _ OnSpellCastBehavior = Card2621104DevotionContract{}
var _ OnDrawBehavior = Card2311101SkyCityZenithStone{}
var _ OnDamagedBehavior = Card2621103BloodGu{}
var _ PerTurnAbility = Card2621103BloodGu{}
var _ OnSpellCastBehavior = Card2521108CouncilJudgmentHammer{}
var _ ElementsGainModifier = Card2521103RedAgateChalice{}
var _ OnUseItemBehavior = Card2221109QuickIceBullet{}
var _ SkillUsePermissionModifier = Card3511102LastStandLight{}
var _ SkillContributionModifier = Card3511102LastStandLight{}
var _ OnCardEnterBehavior = Card1011101CollectorCoralFenlo{}
var _ OnUseItemBehavior = Card1011101CollectorCoralFenlo{}
var _ OnUnitEnterBehavior = Card1521111CouncilConsul{}
var _ OnCardEnterBehavior = Card4011101PureSpiritOshis{}
var _ SkillLearnPermissionModifier = Card3611102ClawOfErebos{}
var _ SkillContributionModifier = Card3611102ClawOfErebos{}
var _ OnSpellCastBehavior = Card3611102ClawOfErebos{}
var _ SkillContributionModifier = Card3621102Retribution{}
var _ OnUseItemBehavior = Card2521107PanaceaP{}
var _ OnUseItemBehavior = Card2121110OfferingTorch{}
var _ OnUseItemBehavior = Card2121101LavafortAshes{}
