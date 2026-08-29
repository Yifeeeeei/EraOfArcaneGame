package game

type Card1621009DemonSummoner struct{ AlwaysActive }

func (Card1621009DemonSummoner) ID() string   { return "1621009" }
func (Card1621009DemonSummoner) Name() string { return "唤魔邪术士" }

func (Card1621009DemonSummoner) OnFriendlyDeath(ctx *EffectContext) error {
	if ctx.Source == nil || ctx.Target == nil || ctx.Target.Card == nil || !ctx.Target.Card.IsCompanion() {
		return nil
	}
	if !triggeredTurnAvailable(ctx.Source) {
		return nil
	}
	candidates := ctx.Engine.friendlyDeckCards(ctx.PlayerID, isShadowConstructOrDemon)
	if len(candidates) == 0 {
		return nil
	}
	if !useTriggeredTurn(ctx.Source) {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "demon_summoner_search",
		"唤魔邪术士:检索1个暗影造物或恶魔", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			ps := ctx.Engine.State.Players[ctx.PlayerID]
			for _, card := range ps.Deck {
				if card != nil && card.InstanceID == selected[0] && isShadowConstructOrDemon(card) {
					ctx.Engine.searchDeckCardToHand(ctx.PlayerID, selected[0])
					return
				}
			}
		})
	return nil
}
