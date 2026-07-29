package game

type Card3321014CallLightning struct{ AlwaysActive }

func (Card3321014CallLightning) ID() string   { return "3321014" }
func (Card3321014CallLightning) Name() string { return "引雷" }

func (Card3321014CallLightning) OnSpellCast(ctx *EffectContext) error {
	if !isSpellBeingCast(ctx) {
		return nil
	}
	hand := ctx.Engine.friendlyHandCards(ctx.PlayerID, nil)
	targets := ctx.Engine.enemyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() &&
			card.Position.Valid() &&
			ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false)
	})
	if len(hand) == 0 || len(targets) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "call_lightning_discard",
		"丢弃1张手牌", hand, 1, 1,
		func(selected []string) {
			if len(selected) == 0 || !ctx.Engine.discardFriendlyCandidate(ctx.PlayerID, selected[0]) {
				return
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "call_lightning_stun",
				"选择1个敌方伙伴晕眩1", targets, 1, 1,
				func(selected []string) {
					for _, id := range selected {
						card := ctx.Engine.findFieldCardByInstance(ctx.Engine.State.Players[ctx.OpponentID], id)
						if card != nil && card.Card != nil && card.Card.IsCompanion() {
							ctx.Engine.addStatus(card, StatusStun, 1)
						}
						return
					}
				})
		})
	return nil
}
