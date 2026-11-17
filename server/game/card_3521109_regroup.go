package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card3521109Regroup struct{ AlwaysActive }

func (Card3521109Regroup) ID() string { return "3521109" }

func (Card3521109Regroup) Name() string { return "重整旗鼓" }

func (Card3521109Regroup) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil {
		return nil
	}
	attacker, ok := ctx.ExtraData["attacker"].(int)
	if !ok || attacker == ctx.PlayerID || !canTriggeredRegroupBeUsed(ctx.Source) || !ctx.Engine.canUseTriggeredRegroup(ctx.Source) {
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion()
	})
	if len(candidates) == 0 {
		return nil
	}
	source := ctx.Source
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	cost := ctx.Engine.effectiveSkillUseCostForPurpose(ps, source, skillPurposeReaction)
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "regroup_buff_companion",
		"重整旗鼓:可以支付费用选择1个友方伙伴获得+1血和负载+1光", candidates, 0, 1, cost, true,
		func(selected []string, data map[string]any) error {
			if len(selected) == 0 {
				return nil
			}
			if !canTriggeredRegroupBeUsed(source) || ctx.Engine.findSkill(ctx.Engine.State.Players[ctx.PlayerID], source.InstanceID) != source {
				return fmt.Errorf("regroup is no longer available")
			}
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target == nil || zone != "unit" || target.Card == nil || !target.Card.IsCompanion() {
				return fmt.Errorf("invalid regroup target")
			}
			if !ctx.Engine.canUseTriggeredRegroup(source) {
				return fmt.Errorf("regroup is not ready")
			}
			if err := ctx.Engine.validateSkillUsePermissionModifiers(source, skillPurposeReaction); err != nil {
				return err
			}
			units, err := ctx.Engine.collectOverexertUnits(ps, stringsFromActionData(data, "overexert_ids"))
			if err != nil {
				return err
			}
			if !ctx.Engine.payDefenseCostWithOptions(ps, cost, ActionMessage{Data: data}, units, ctx.Engine.playerHasLightWildcard(ps)) {
				return fmt.Errorf("invalid regroup payment")
			}
			source.IsHorizontal = true
			ctx.Engine.ApplyKeywordOnSkillUse(source)
			ctx.Engine.applySkillUseCooldownModifiers(ctx.Engine.State.Players[ctx.PlayerID], source)
			ctx.Engine.consumeNextSkillUseModifiersForPurpose(ps, source, skillPurposeReaction)
			ctx.Engine.notifySkillUseCommitted(ctx.PlayerID, source)
			ctx.Engine.advanceMasteryForUsedSkills(ctx.PlayerID, source)
			ctx.Engine.destroyFuyeDoomedAfterExert(units)
			ctx.Engine.recordSpellCast(ctx.PlayerID, source)
			castData := map[string]any{"cast_player": ctx.PlayerID, "attacker": ctx.PlayerID, "purpose": string(skillPurposeReaction), "skill": cardToInfo(source), "spell_source": source, "overexert_units": units}
			finish := func() {
				if !ctx.Engine.unitStillOnField(target) {
					return
				}
				ctx.Engine.gainLife(target, 1, source)
				ctx.Engine.addElementsGainBonus(target, ctx.PlayerID, model.ElementLight, 1, source)
				ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{"source": cardToInfo(source), "target": cardToInfo(target), "effect": "regroup_buff_companion"}})
			}
			ctx.Engine.emit(GameEvent{Type: "spell_cast", Player: -1, Data: castData})
			ctx.Engine.triggerEffects(TriggerOnSpellCast, source, nil, castData)
			if !ctx.Engine.triggerSpellCastFieldEffectsWithContinuation(ctx.PlayerID, source, castData, finish) {
				finish()
			}

			return nil
		})
	return nil
}

func canTriggeredRegroupBeUsed(skill *CardInstance) bool {
	return skill != nil && skill.Card != nil && skill.Card.Number == "3521109" && !skill.IsHorizontal && skill.Statuses[StatusCooldown] <= 0
}

func (e *Engine) canUseTriggeredRegroup(skill *CardInstance) bool {
	return canTriggeredRegroupBeUsed(skill) && e.validateReadySkill(skill) == nil && !e.hasEffectiveStatus(skill, StatusStun) && skill.Statuses[StatusCannotUseSkillUntilTurn] < e.turnOrdinal() && !e.actionRestricted(RuleUseSkill, skill)
}
