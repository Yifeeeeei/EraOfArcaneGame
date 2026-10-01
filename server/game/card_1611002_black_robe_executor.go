package game

type Card1611002BlackRobeExecutor struct{ AlwaysActive }

func (Card1611002BlackRobeExecutor) ID() string { return "1611002" }

func (Card1611002BlackRobeExecutor) Name() string { return "黑袍执行官 无心" }

func (Card1611002BlackRobeExecutor) OnFriendlyDeath(ctx *EffectContext) error {
	cause, _ := ctx.ExtraData["death_cause"].(string)
	if ctx.Target != nil && ctx.Target.Card.IsCompanion() && (cause == DeathCauseSacrifice || cause == DeathCauseDevour) {
		ctx.Source.Statuses["暗影标记"] += max(ctx.Target.Card.Life, 1)
	}
	return nil
}

func (Card1611002BlackRobeExecutor) OnUltimate(ctx *EffectContext) error {
	targets := ctx.Engine.enemyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card.Card.IsCompanion() && card.Position != nil &&
			ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, cardHasPierce(ctx.Source)) &&
			ctx.Source.Statuses["暗影标记"] >= max(card.CurrentLife, 1)
	})
	if len(targets) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "black_robe_executor_destroy",
		"黑袍执行官:选择法力范围内1个可支付暗影标记的敌方伙伴消灭", targets, 1, 1,
		func(selected []string) {
			target := selectedUnitFromCandidates(ctx.Engine, selected, targets)
			if target == nil || target.Position == nil || !target.Card.IsCompanion() {
				return
			}
			cost := max(target.CurrentLife, 1)
			if ctx.Source.Statuses["暗影标记"] < cost {
				return
			}
			ctx.Source.Statuses["暗影标记"] -= cost
			ctx.Engine.destroyUnit(target, target.OwnerID)
		})
	return nil
}
