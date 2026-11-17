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
	ctx.Engine.flipDeckMatchesToHandThen(ctx.PlayerID, 1, 0, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion()
	}, func(drawn []*CardInstance) {
		for _, card := range drawn {
			card.Statuses["入场费用"+model.ElementLight+"-1"]++
		}
	})
	return nil
}
