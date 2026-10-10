package game

import (
	"fmt"
)

type Card1321109StormHorn struct{ AlwaysActive }

func (Card1321109StormHorn) ID() string { return "1321109" }

func (Card1321109StormHorn) Name() string { return "风暴之角" }

func (Card1321109StormHorn) OnUltimate(ctx *EffectContext) error {
	handCandidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, nil)
	if len(handCandidates) == 0 || !ctx.Engine.hasAirEquipmentInDeck(ctx.PlayerID) {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "storm_horn_discard",
		"风暴之角:丢弃1张手牌", handCandidates, 1, 1,
		func(selected []string) {
			if ctx.Engine.discardSelectedHandCards(ctx.PlayerID, selected, 1) != 1 {
				return
			}
			searchDeckToHandByPredicate(ctx, "storm_horn_search_air_equipment", "风暴之角:翻取1张大气装备", isAirEquipment)
		})
	return nil
}

func (Card1321109StormHorn) ValidateAbility(ctx *EffectContext, trigger EffectTrigger) error {
	if trigger != TriggerUltimate {
		return nil
	}
	if len(ctx.Engine.State.Players[ctx.PlayerID].Hand) == 0 {
		return fmt.Errorf("风暴之角需要丢弃1张手牌")
	}
	if !ctx.Engine.hasAirEquipmentInDeck(ctx.PlayerID) {
		return fmt.Errorf("风暴之角需要卡组中有可翻取的大气装备")
	}
	return nil
}

func (e *Engine) hasAirEquipmentInDeck(playerID int) bool {
	if playerID < 0 || playerID >= len(e.State.Players) {
		return false
	}
	for _, card := range e.State.Players[playerID].Deck {
		if isAirEquipment(card) && canFlipOrSearchCard(card) {
			return true
		}
	}
	return false
}
