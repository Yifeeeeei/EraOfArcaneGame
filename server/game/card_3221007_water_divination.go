package game

import "eraofarcane/model"

type Card3221007WaterDivination struct{ AlwaysActive }

func (Card3221007WaterDivination) ID() string   { return "3221007" }
func (Card3221007WaterDivination) Name() string { return "水占术" }

func (Card3221007WaterDivination) OnSpellCast(ctx *EffectContext) error {
	if !isSpellBeingCast(ctx) {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	lookCount := min(4, len(ps.Deck))
	if lookCount == 0 {
		return nil
	}
	looked := append([]*CardInstance(nil), ps.Deck[:lookCount]...)
	candidates := make([]map[string]any, 0, len(looked))
	hasWater := false
	for i, card := range looked {
		info := candidateInfo(card, "deck", "own")
		info["deck_index"] = i
		info["can_select"] = card.Card.Category == model.ElementWater
		if card.Card.Category == model.ElementWater {
			hasWater = true
		}
		candidates = append(candidates, info)
	}
	minSelect := 0
	if hasWater {
		minSelect = 1
	}
	ctx.Engine.SetPendingActionWithData(ctx.PlayerID, "water_divination_search",
		"查看牌堆顶4张，检索其中1张水纹卡，其余按任意顺序置于牌堆顶或牌堆底", candidates, minSelect, 1,
		func(selected []string, data map[string]any) {
			resolveWaterDivination(ctx.Engine, ctx.PlayerID, looked, selected, data)
		})
	return nil
}

func resolveWaterDivination(e *Engine, playerID int, looked []*CardInstance, selected []string, data map[string]any) {
	if e == nil {
		return
	}
	ps := e.State.Players[playerID]
	lookedByID := make(map[string]*CardInstance, len(looked))
	for _, card := range looked {
		if card != nil {
			lookedByID[card.InstanceID] = card
		}
	}

	var searched *CardInstance
	if len(selected) > 0 {
		if card := lookedByID[selected[0]]; card != nil && card.Card.Category == model.ElementWater {
			searched = card
		}
	}

	lookCount := min(len(looked), len(ps.Deck))
	rest := append([]*CardInstance(nil), ps.Deck[lookCount:]...)
	pool := make(map[string]*CardInstance, len(looked))
	for _, card := range looked {
		if card == nil || card == searched {
			continue
		}
		pool[card.InstanceID] = card
	}

	topIDs := stringsFromActionData(data, "top_order")
	bottomIDs := stringsFromActionData(data, "bottom_order")
	used := make(map[string]bool, len(pool))
	top := orderedWaterDivinationCards(topIDs, pool, used)
	bottom := orderedWaterDivinationCards(bottomIDs, pool, used)
	for _, card := range looked {
		if card == nil || card == searched || used[card.InstanceID] {
			continue
		}
		top = append(top, card)
		used[card.InstanceID] = true
	}

	ps.Deck = append(append(top, rest...), bottom...)
	if searched != nil {
		ps.Hand = append(ps.Hand, searched)
		e.emit(GameEvent{Type: "search_card", Player: playerID, Data: map[string]any{"card": cardToInfo(searched)}})
	}
	e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
		"effect":       "water_divination_reorder",
		"searched":     cardToInfo(searched),
		"top_order":    cardsToInfo(top),
		"bottom_order": cardsToInfo(bottom),
	}})
}

func stringsFromActionData(data map[string]any, key string) []string {
	if data == nil {
		return nil
	}
	values, _ := data[key].([]any)
	return stringsFromAnySlice(values)
}

func orderedWaterDivinationCards(ids []string, pool map[string]*CardInstance, used map[string]bool) []*CardInstance {
	cards := make([]*CardInstance, 0, len(ids))
	for _, id := range ids {
		if used[id] {
			continue
		}
		card := pool[id]
		if card == nil {
			continue
		}
		cards = append(cards, card)
		used[id] = true
	}
	return cards
}
