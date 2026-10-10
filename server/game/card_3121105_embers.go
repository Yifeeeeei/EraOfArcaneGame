package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card3121105Embers struct{ AlwaysActive }

func (Card3121105Embers) ID() string { return "3121105" }

func (Card3121105Embers) Name() string { return "余火" }

func (Card3121105Embers) OnTurnEnd(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.PlayerID < 0 || ctx.PlayerID >= len(ctx.Engine.State.Players) {
		return nil
	}
	if ctx.Engine.State.CurrentTurn != ctx.PlayerID {
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
		"余火:可以选择目标免费使用此卡", candidates, 0, 1, nil, false,
		func(selected []string, _ map[string]any) error {
			if len(selected) == 0 {
				return nil
			}
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
