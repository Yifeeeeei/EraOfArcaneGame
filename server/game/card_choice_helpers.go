package game

import (
	"eraofarcane/model"
	"fmt"
)

const (
	erebosSoulChainMarkedUnitStatus  = "erebos_soul_chain_marked_unit"
	erebosSoulChainMarkedSpellStatus = "erebos_soul_chain_marked_spell"
)

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

const (
	fireButterflyTemporaryLoadStatus     = "火蝴蝶临时负载"
	fireButterflyPreviousLoadSetStatus   = "火蝴蝶原负载覆盖"
	fireButterflyPreviousLoadValuePrefix = "火蝴蝶原负载:"
)

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

func isThunderlightItem(card *CardInstance) bool {
	if card == nil || card.Card == nil || !card.Card.IsItem() {
		return false
	}
	gain := effectiveElementsGain(card)
	return gain[model.ElementAir] > 0 && gain[model.ElementLight] > 0
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

type royalInfusionRune struct {
	AlwaysActive
	id          string
	name        string
	powerBonus  int
	attackBonus int
}

func (r royalInfusionRune) ID() string { return r.id }

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
			virtual := e.cloneVirtualSpell(original, ownerID, e.State.TurnNumber)
			boosts := e.cloneSpellInstances(spellInstancesFromData(extraData, "boost_skills"), ownerID, e.State.TurnNumber)
			target := SpellTarget{Type: "unit", Position: *targetUnit.Position}
			if targetUnit.OwnerID == ownerID {
				target.OwnerID = &ownerID
			}
			return e.replacePendingSpell(func() error {
				return e.startVirtualSpellCastWithBoosts(ownerID, virtual, target, boosts, map[string]any{
					"triggered_by": "2321111",
					"source_item":  counterID,
				})
			})
		})
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

const (
	treadingWaveTriggerTurnStatus  = "treading_wave_trigger_turn"
	treadingWaveTriggerCountStatus = "treading_wave_trigger_count"
)

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
