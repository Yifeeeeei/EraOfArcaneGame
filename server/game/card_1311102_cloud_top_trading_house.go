package game

import "eraofarcane/model"

type Card1311102CloudTopTradingHouse struct{ AlwaysActive }

func (Card1311102CloudTopTradingHouse) ID() string   { return "1311102" }
func (Card1311102CloudTopTradingHouse) Name() string { return "云顶商行 克罗斯" }
func (Card1311102CloudTopTradingHouse) DrawFromEmptyDeck(ctx *EffectContext, n int) []*CardInstance {
	return ctx.Engine.drawCardsFromOpponentDeckAsNeutral(ctx.PlayerID, n)
}
func (e *Engine) drawCardsFromOpponentDeckAsNeutral(playerID int, n int) []*CardInstance {
	if e == nil || n <= 0 || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	ps := e.State.Players[playerID]
	opponent := e.State.Players[1-playerID]
	drawn := make([]*CardInstance, 0, n)
	for i := 0; i < n && len(opponent.Deck) > 0; i++ {
		card := opponent.Deck[0]
		opponent.Deck = opponent.Deck[1:]
		card.OwnerID = playerID
		card.Statuses[StatusEntryCostNeutralAmount] = totalElementCost(card.Card.ElementsCost)
		setElementsGain(card, map[string]int{model.ElementArcane: totalElementCost(card.Card.ElementsGain)})
		card.ElementsGainBonus = make(map[string]int)
		ps.Hand = append(ps.Hand, card)
		drawn = append(drawn, card)
	}
	if len(drawn) > 0 {
		e.emit(GameEvent{
			Type:   "cloud_top_trading_house_draw",
			Player: -1,
			Data: map[string]any{
				"player": playerID,
				"count":  len(drawn),
			},
		})
	}
	return drawn
}
