package game

type Card1221016IceSpikeFortress struct{ AlwaysActive }

func (Card1221016IceSpikeFortress) ID() string { return "1221016" }

func (Card1221016IceSpikeFortress) Name() string { return "冰刺堡垒" }

func (Card1221016IceSpikeFortress) DamageScope() DamageScope { return DamageSelf }

func (Card1221016IceSpikeFortress) OnDamaged(ctx *EffectContext, event DamageEvent) error {
	attacker, hasAttacker := event.SourcePlayer, event.SourcePlayer >= 0
	if !hasAttacker || attacker == ctx.PlayerID {
		return nil
	}
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card.Position != nil && ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "ice_spike_fortress",
		"冰刺堡垒:选择法力范围内1个敌人冻结1，若已冻结则造成1点伤害", candidates, 1, 1,
		func(selected []string) {
			target := selectedUnitFromCandidates(ctx.Engine, selected, candidates)
			if target == nil {
				return
			}
			if target.Statuses[StatusFreeze] > 0 {
				ctx.Engine.ApplyDamage(DamageRequest{Target: target, Amount: 1, Kind: "effect", SourcePlayer: ctx.PlayerID, SourceKnown: true, Source: ctx.Source})
				return
			}
			ctx.Engine.addStatus(target, StatusFreeze, 1)
		})
	return nil
}
