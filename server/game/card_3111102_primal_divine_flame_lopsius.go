package game

import (
	"fmt"
)

type Card3111102PrimalDivineFlameLopsius struct{ AlwaysActive }

func (Card3111102PrimalDivineFlameLopsius) ID() string { return "3111102" }

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
	usedAtActivation := ctx.Source.UsedThisTurn
	ctx.Engine.SetPendingAction(ctx.PlayerID, "primal_divine_flame_exile",
		"原初神炎 洛普修斯:选择1个火焰技能移出游戏", candidates, 1, 1,
		func(selected []string) {
			if ctx.Source.UsedThisTurn > usedAtActivation+1 {
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
			if ctx.Source.UsedThisTurn == usedAtActivation {
				ctx.Source.UsedThisTurn++
			}
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
