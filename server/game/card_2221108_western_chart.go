package game

import (
	"eraofarcane/model"
	"strings"
)

type Card2221108WesternChart struct{ AlwaysActive }

func (Card2221108WesternChart) ID() string { return "2221108" }

func (Card2221108WesternChart) Name() string { return "西境航海图" }

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

func (Card2221108WesternChart) SpellTargetGrant(ctx *EffectContext, skill *CardInstance, _ SpellTarget) SpellTargetGrant {
	return SpellTargetGrant{Pierce: skill != nil && skill.InstanceID != "" && ctx.Source.Statuses[westernChartPierceTargetPrefix+skill.InstanceID] > 0}
}

const westernChartPierceTargetPrefix = "西境航海图穿透目标:"
