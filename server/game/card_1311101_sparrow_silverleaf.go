package game

import ()

type Card1311101SparrowSilverleaf struct{ AlwaysActive }

func (Card1311101SparrowSilverleaf) ID() string { return "1311101" }

func (Card1311101SparrowSilverleaf) Name() string { return "斯帕罗 银叶" }

func (Card1311101SparrowSilverleaf) OnEnter(ctx *EffectContext) error {
	damage := min(ctx.Engine.State.Players[ctx.PlayerID].DiscardedHandCountThisTurn, 3)
	if damage <= 0 {
		return nil
	}
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card.Position != nil && ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "sparrow_silverleaf_entry_damage",
		"斯帕罗 银叶:选择法力范围内1名敌人造成弃牌数量伤害", candidates, 1, 1,
		func(selected []string) {
			target := findEnemyCardCandidate(ctx.Engine, ctx.PlayerID, firstSelected(selected), candidates)
			if target == nil || target.Position == nil || !ctx.Engine.IsInSpellRange(ctx.PlayerID, target.Position.Col, target.Position.Row, false) {
				return
			}
			ctx.DealDamage(target, damage)
		})
	return nil
}
