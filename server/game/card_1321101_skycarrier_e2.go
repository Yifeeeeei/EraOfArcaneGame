package game

import (
	"eraofarcane/model"
)

type Card1321101SkycarrierE2 struct{ AlwaysActive }

func (Card1321101SkycarrierE2) ID() string { return "1321101" }

func (Card1321101SkycarrierE2) Name() string { return "翱翔者E2型运输舰" }

func (Card1321101SkycarrierE2) IsPrayerAbility() bool { return true }

func (Card1321101SkycarrierE2) OnPerTurn(ctx *EffectContext) error {
	choices := []map[string]any{
		{"instance_id": "draw", "name": "抽2张牌", "zone": "choice", "side": "own"},
	}
	if len(airGraveyardCandidates(ctx.Engine.State.Players[ctx.PlayerID])) >= 2 {
		choices = append(choices, map[string]any{"instance_id": "recycle", "name": "将2张大气弃牌洗回牌组", "zone": "choice", "side": "own"})
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "skycarrier_e2_prayer",
		"翱翔者E2型运输舰:选择祈咒效果", choices, 1, 1,
		func(selected []string) {
			switch firstSelected(selected) {
			case "draw":
				ctx.Engine.drawCards(ctx.PlayerID, 2)
			case "recycle":
				openSkycarrierRecyclePrompt(ctx)
			}
		})
	return nil
}

func openSkycarrierRecyclePrompt(ctx *EffectContext) {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	candidates := airGraveyardCandidates(ps)
	if len(candidates) < 2 {
		return
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "skycarrier_e2_recycle",
		"翱翔者E2型运输舰:选择2张大气弃牌洗回牌组", candidates, 2, 2,
		func(selected []string) {
			moveSelectedAirGraveyardCardsToDeck(ps, selected, 2)
			ctx.Engine.shuffleDeck(ctx.PlayerID)
		})
}

func airGraveyardCandidates(ps *PlayerState) []map[string]any {
	candidates := make([]map[string]any, 0)
	if ps == nil {
		return candidates
	}
	for _, card := range ps.Graveyard {
		if card != nil && card.Card != nil && card.Card.Category == model.ElementAir {
			candidates = append(candidates, candidateInfo(card, "graveyard", "own"))
		}
	}
	return candidates
}

func moveSelectedAirGraveyardCardsToDeck(ps *PlayerState, selected []string, maxCount int) {
	if ps == nil || maxCount <= 0 {
		return
	}
	selectedSet := make(map[string]bool, len(selected))
	for _, id := range selected {
		selectedSet[id] = true
	}
	moved := 0
	for i := 0; i < len(ps.Graveyard) && moved < maxCount; {
		card := ps.Graveyard[i]
		if card != nil && selectedSet[card.InstanceID] && card.Card != nil && card.Card.Category == model.ElementAir {
			ps.Graveyard = append(ps.Graveyard[:i], ps.Graveyard[i+1:]...)
			ps.Deck = append(ps.Deck, card)
			moved++
			continue
		}
		i++
	}
}
