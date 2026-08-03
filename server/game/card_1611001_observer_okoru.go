package game

type Card1611001ObserverOkoru struct{ AlwaysActive }

func (Card1611001ObserverOkoru) ID() string   { return "1611001" }
func (Card1611001ObserverOkoru) Name() string { return "\"观察者\" 欧柯茹" }
func (Card1611001ObserverOkoru) OnEnter(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	lookCount := min(5, len(ps.Deck))
	if lookCount == 0 {
		return nil
	}
	looked := append([]*CardInstance(nil), ps.Deck[:lookCount]...)
	candidates := make([]map[string]any, 0, len(looked))
	for i, card := range looked {
		info := candidateInfo(card, "deck", "own")
		info["deck_index"] = i
		candidates = append(candidates, info)
	}
	ctx.Engine.SetPendingActionWithData(ctx.PlayerID, "okoru_observe",
		"欧柯茹:查看牌堆顶5张，可抽取任意张，其余按任意顺序置于牌堆顶或牌堆底；每抽1张受到1点伤害",
		candidates, 0, lookCount,
		func(selected []string, data map[string]any) {
			resolveOkoruObserve(ctx.Engine, ctx.PlayerID, looked, selected, data)
		})
	return nil
}

func resolveOkoruObserve(e *Engine, playerID int, looked []*CardInstance, selected []string, data map[string]any) {
	if e == nil {
		return
	}
	ps := e.State.Players[playerID]
	selectedSet := make(map[string]bool, len(selected))
	for _, id := range selected {
		selectedSet[id] = true
	}
	lookCount := min(len(looked), len(ps.Deck))
	rest := append([]*CardInstance(nil), ps.Deck[lookCount:]...)
	pool := make(map[string]*CardInstance, len(looked))
	drawn := make([]*CardInstance, 0, len(selected))
	for _, card := range looked {
		if card == nil {
			continue
		}
		if selectedSet[card.InstanceID] {
			drawn = append(drawn, card)
			continue
		}
		pool[card.InstanceID] = card
	}
	used := make(map[string]bool, len(pool))
	top := orderedWaterDivinationCards(stringsFromActionData(data, "top_order"), pool, used)
	bottom := orderedWaterDivinationCards(stringsFromActionData(data, "bottom_order"), pool, used)
	for _, card := range looked {
		if card == nil || selectedSet[card.InstanceID] || used[card.InstanceID] {
			continue
		}
		top = append(top, card)
		used[card.InstanceID] = true
	}
	ps.Deck = append(append(top, rest...), bottom...)
	e.appendCardsToHand(playerID, drawn)
	for _, card := range drawn {
		ps.DrawCountThisTurn++
		e.emit(GameEvent{Type: "draw_card", Player: playerID, Data: map[string]any{"card": cardToInfo(card)}})
		e.triggerFieldEffectsWithData(TriggerOnDraw, playerID, card, map[string]any{
			"drawn_card":           card,
			"drawn_player":         playerID,
			"draw_count_this_turn": ps.DrawCountThisTurn,
		})
	}
	e.enforceImmediateHandLimitAfterHandGain(playerID)
	if len(drawn) > 0 && ps.Hero != nil {
		e.dealDamage(ps.Hero, len(drawn), playerID)
	}
	e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
		"effect":       "okoru_observe",
		"drawn":        cardsToInfo(drawn),
		"top_order":    cardsToInfo(top),
		"bottom_order": cardsToInfo(bottom),
	}})
}
