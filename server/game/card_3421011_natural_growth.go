package game

import "eraofarcane/model"

type Card3421011NaturalGrowth struct{}

func (Card3421011NaturalGrowth) ID() string   { return "3421011" }
func (Card3421011NaturalGrowth) Name() string { return "自然生长" }

func (Card3421011NaturalGrowth) OnSpellCast(ctx *EffectContext) error {
	if !isSpellBeingCast(ctx) {
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card.Card.IsCompanion() && card.Card.Category == model.ElementEarth && totalLoad(card) < 4
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "natural_growth",
		"选择1个地脉伙伴，回合结束时负载+1地", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			scheduleLoadGainAtTurnEnd(ctx.Engine.State.Players[ctx.PlayerID], selected[0], model.ElementEarth, 1)
		})
	return nil
}
