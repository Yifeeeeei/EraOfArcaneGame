package game

type Card1321004LightningElemental struct{ AlwaysActive }

func (Card1321004LightningElemental) ID() string   { return "1321004" }
func (Card1321004LightningElemental) Name() string { return "雷电元素" }
func (Card1321004LightningElemental) OnEnter(ctx *EffectContext) error {
	candidates := append(ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != ctx.Source && card.Card.IsCompanion()
	}), ctx.Engine.enemyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card.Card.IsCompanion() && card.Position != nil && ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false)
	})...)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "lightning_elemental_stun",
		"雷电元素:选择法力范围内1个伙伴晕眩1", candidates, 1, 1,
		func(selected []string) {
			target := selectedUnitFromCandidates(ctx.Engine, selected, candidates)
			if target != nil {
				target.Statuses[StatusStun]++
			}
		})
	return nil
}
