package game

import (
	"eraofarcane/model"
)

type Card1521113RadiantWatchdog struct{ AlwaysActive }

func (Card1521113RadiantWatchdog) ID() string { return "1521113" }

func (Card1521113RadiantWatchdog) Name() string { return "辉之都戒卫犬" }

func (Card1521113RadiantWatchdog) OnDeath(ctx *EffectContext) error {
	if ctx.ExtraData == nil {
		return nil
	}
	attacker, ok := ctx.ExtraData["attacker"].(int)
	if !ok || attacker == ctx.PlayerID {
		return nil
	}
	candidates := ctx.Engine.friendlyDeckCards(ctx.PlayerID, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion()
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "radiant_watchdog_search",
		"辉之都戒卫犬:翻取1个伙伴牌并使其入场花费-1光", candidates, 1, 1,
		func(selected []string) {
			card := ctx.Engine.searchDeckCardToHand(ctx.PlayerID, firstSelected(selected))
			if card != nil {
				card.Statuses["入场费用"+model.ElementLight+"-1"]++
			}
		})
	return nil
}
