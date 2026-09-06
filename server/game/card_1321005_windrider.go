package game

import (
	"eraofarcane/model"
)

type Card1321005Windrider struct{ AlwaysActive }

func (Card1321005Windrider) ID() string { return "1321005" }

func (Card1321005Windrider) Name() string { return "驭风师" }

func (Card1321005Windrider) OnUltimate(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "discard_any",
		"选择任意数量手牌弃置，每张获得1点大气",
		candidates, 0, len(candidates),
		func(selected []string) {
			discarded := 0
			for _, id := range selected {
				if ctx.Engine.discardFriendlyCandidate(ctx.PlayerID, id) {
					discarded++
				}
			}
			ctx.Engine.State.Players[ctx.PlayerID].Elements[model.ElementAir] += discarded
		})
	return nil
}
