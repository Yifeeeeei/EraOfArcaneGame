package game

import "eraofarcane/model"

type Card4411001Whitebeard struct{ AlwaysActive }

func (Card4411001Whitebeard) ID() string   { return "4411001" }
func (Card4411001Whitebeard) Name() string { return "森林隐士 白须" }

const whitebeardFirstTurnChecked = "whitebeard_first_turn_checked"

func (Card4411001Whitebeard) OnTurnStart(ctx *EffectContext) error {
	if ctx.Source.Statuses[whitebeardFirstTurnChecked] > 0 {
		return nil
	}
	ctx.Source.Statuses[whitebeardFirstTurnChecked] = 1
	candidates := ctx.Engine.friendlyDeckCards(ctx.PlayerID, func(card *CardInstance) bool {
		return card.Card.Category == model.ElementEarth && isBeastPlantOrSpirit(card)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "whitebeard_first_turn_search", "森林隐士 白须:可以检索1张地属性野兽、植物或精灵", candidates, 0, 1, func(selected []string) {
		if len(selected) == 0 {
			return
		}
		ctx.Engine.searchDeckCardToHand(ctx.PlayerID, selected[0])
	})
	return nil
}
