package game

import (
	"fmt"
)

type Card1211002Leviathan struct{ AlwaysActive }

func (Card1211002Leviathan) ID() string { return "1211002" }

func (Card1211002Leviathan) Name() string { return "深渊巨口 利维坦" }

func (Card1211002Leviathan) PerTurnLabel(*CardInstance) string {
	return "消耗"
}

func (Card1211002Leviathan) OnPerTurn(ctx *EffectContext) error {
	if ctx.Source.IsHorizontal {
		return fmt.Errorf("深渊巨口 利维坦已经横置")
	}
	if ctx.Source.Statuses[leviathanCooldownStatus] > 0 {
		return fmt.Errorf("深渊巨口 利维坦的主动效果正在冷却")
	}
	targets := ctx.Engine.enemyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card.Card.IsCompanion() && card.Position != nil && ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false)
	})
	if len(targets) == 0 {
		return fmt.Errorf("没有可消灭的敌方伙伴")
	}
	ctx.Source.IsHorizontal = true
	ctx.Engine.SetPendingAction(ctx.PlayerID, "leviathan_destroy",
		"利维坦:选择法力范围内1个敌方伙伴消灭", targets, 1, 1,
		func(selected []string) {
			target := selectedUnitFromCandidates(ctx.Engine, selected, targets)
			if target != nil {
				ctx.Engine.destroyUnit(target, ctx.OpponentID)
				ctx.Source.Statuses[leviathanCooldownStatus] = 2
			}
		})
	return nil
}

func (Card1211002Leviathan) OnTurnEnd(ctx *EffectContext) error {
	if endedPlayer, ok := ctx.ExtraData["ended_player"].(int); ok && endedPlayer != ctx.PlayerID {
		return nil
	}
	if ctx.Source.Statuses[leviathanCooldownStatus] > 0 {
		ctx.Source.Statuses[leviathanCooldownStatus]--
	}
	return nil
}
