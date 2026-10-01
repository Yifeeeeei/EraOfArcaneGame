package game

import (
	"eraofarcane/model"
)

type Card3521109Regroup struct{ AlwaysActive }

func (Card3521109Regroup) ID() string { return "3521109" }

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
			ctx.Engine.gainLife(target, 1, source)
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
