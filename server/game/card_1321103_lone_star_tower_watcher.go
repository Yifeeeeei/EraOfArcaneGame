package game

type Card1321103LoneStarTowerWatcher struct{ AlwaysActive }

func (Card1321103LoneStarTowerWatcher) ID() string { return "1321103" }

func (Card1321103LoneStarTowerWatcher) Name() string { return "孤星塔守望者" }

func (Card1321103LoneStarTowerWatcher) OnUltimate(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "lone_star_tower_watcher_discard",
		"孤星塔守望者:丢弃至多3张手牌并获得等量护盾", candidates, 0, min(3, len(candidates)),
		func(selected []string) {
			discarded := ctx.Engine.discardSelectedHandCards(ctx.PlayerID, selected, 3)
			if discarded > 0 {
				ctx.Engine.gainPlayerShield(ctx.PlayerID, discarded)
			}
		})
	return nil
}
