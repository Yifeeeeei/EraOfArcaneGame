package game

import (
	"eraofarcane/model"
)

type Card3221101TreadingWave struct{ AlwaysActive }

func (Card3221101TreadingWave) ID() string { return "3221101" }

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

var _ OnSpellCastBehavior = Card3221101TreadingWave{}

func spellCastCardForTrigger(ctx *EffectContext) *CardInstance {
	if ctx == nil {
		return nil
	}
	if ctx.Target != nil {
		return ctx.Target
	}
	return ctx.Source
}
