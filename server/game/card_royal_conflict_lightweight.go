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
var _ SkillLearnPermissionModifier = Card3611102ClawOfErebos{}
var _ SkillContributionModifier = Card3611102ClawOfErebos{}
var _ OnSpellCastBehavior = Card3611102ClawOfErebos{}
