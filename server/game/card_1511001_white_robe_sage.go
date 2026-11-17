package game

import (
	"fmt"
)

type Card1511001WhiteRobeSage struct{ AlwaysActive }

func (Card1511001WhiteRobeSage) ID() string { return "1511001" }

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
