package game

type Card1221004FrostPuppet struct{ AlwaysActive }

func (Card1221004FrostPuppet) ID() string   { return "1221004" }
func (Card1221004FrostPuppet) Name() string { return "寒霜傀儡" }
func (Card1221004FrostPuppet) OnEnter(ctx *EffectContext) error {
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card.Card.IsCompanion() && card.Position != nil && ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "frost_puppet_freeze",
		"寒霜傀儡:选择法力范围内1个敌方伙伴冻结1", candidates, 1, 1,
		func(selected []string) {
			target := selectedUnitFromCandidates(ctx.Engine, selected, candidates)
			if target != nil {
				target.Statuses[StatusFreeze]++
			}
		})
	return nil
}
