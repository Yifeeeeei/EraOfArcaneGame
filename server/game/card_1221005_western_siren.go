package game

type Card1221005WesternSiren struct{ AlwaysActive }

func (Card1221005WesternSiren) ID() string   { return "1221005" }
func (Card1221005WesternSiren) Name() string { return "西境海妖" }

func (Card1221005WesternSiren) OnPerTurn(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != ctx.Source && card.Card.IsCompanion()
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "western_siren_consume",
		"选择法力范围内1个伙伴横置", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			target := ctx.Engine.findFieldCardByInstance(ctx.Engine.State.Players[ctx.PlayerID], selected[0])
			if target != nil {
				target.IsHorizontal = true
			}
		})
	return nil
}
