package game

import "eraofarcane/model"

type Card3221007WaterDivination struct{ AlwaysActive }

func (Card3221007WaterDivination) ID() string   { return "3221007" }
func (Card3221007WaterDivination) Name() string { return "水占术" }

func (Card3221007WaterDivination) OnSpellCast(ctx *EffectContext) error {
	if !isSpellBeingCast(ctx) {
		return nil
	}
	candidates := ctx.Engine.friendlyTopDeckCards(ctx.PlayerID, 4, func(card *CardInstance) bool {
		return card.Card.Category == model.ElementWater
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "water_divination_search",
		"查看牌堆顶4张，检索其中1张水纹卡", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			ctx.Engine.searchDeckToHand(ctx.PlayerID, selected[0])
		})
	return nil
}
