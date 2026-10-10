package game

import (
	"eraofarcane/model"
)

type Card1411103KingOfBeasts struct{ AlwaysActive }

func (Card1411103KingOfBeasts) ID() string { return "1411103" }

func (Card1411103KingOfBeasts) Name() string { return "百兽之王 莱恩克塞斯" }

func (Card1411103KingOfBeasts) OnEnter(ctx *EffectContext) error {
	drawn := ctx.Engine.flipDeckMatchesToHandThen(ctx.PlayerID, 1, 0, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion()
	}, func(drawn []*CardInstance) {
		if len(drawn) == 0 || drawn[0].Card.Category != model.ElementEarth {
			return
		}
		positions := ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)
		if len(positions) == 0 {
			return
		}
		cardID := drawn[0].InstanceID
		ctx.Engine.SetPendingAction(ctx.PlayerID, "king_of_beasts_summon_earth_companion",
			"百兽之王 莱恩克塞斯:选择位置免费召唤翻取的地脉伙伴", positions, 1, 1,
			func(selected []string) {
				pos, ok := positionFromSelectionID(firstSelected(selected))
				if !ok {
					return
				}
				card, _ := ctx.Engine.State.Players[ctx.PlayerID].FindHandCard(cardID)
				if card == nil || card.Card == nil || !card.Card.IsCompanion() || card.Card.Category != model.ElementEarth {
					return
				}
				summonCardFreeFromHandOrDeckAtPosition(ctx, cardID, pos)
			})
	})
	_ = drawn
	return nil
}
