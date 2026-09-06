package game

type Card2611101ErebosSoulChain struct{ AlwaysActive }

func (Card2611101ErebosSoulChain) ID() string { return "2611101" }

func (Card2611101ErebosSoulChain) Name() string { return "厄瑞波斯的魂链" }

func (Card2611101ErebosSoulChain) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.UltimateUsed || ctx.Target == nil || ctx.Target.Card == nil || ctx.ExtraData == nil {
		return nil
	}
	castPlayer, _ := ctx.ExtraData["cast_player"].(int)
	if castPlayer == ctx.PlayerID {
		return nil
	}
	overexertUnits := spellInstancesFromData(ctx.ExtraData, "overexert_units")
	if len(overexertUnits) == 0 {
		return nil
	}
	markedUnits := 0
	for _, unit := range overexertUnits {
		if unit == nil || unit.Card == nil || unit.OwnerID != castPlayer || !unit.Card.IsCompanion() {
			continue
		}
		if unit.Statuses == nil {
			unit.Statuses = make(map[string]int)
		}
		unit.Statuses[erebosSoulChainMarkedUnitStatus] = 1
		markedUnits++
	}
	if markedUnits == 0 {
		return nil
	}
	markedSpells := 0
	for _, spell := range append([]*CardInstance{ctx.Target}, spellInstancesFromData(ctx.ExtraData, "boost_skills")...) {
		if spell == nil || spell.Card == nil || spell.OwnerID != castPlayer || !isSpellLikeCard(spell.Card) {
			continue
		}
		if spell.Statuses == nil {
			spell.Statuses = make(map[string]int)
		}
		spell.Statuses[erebosSoulChainMarkedSpellStatus] = 1
		markedSpells++
	}
	if markedSpells > 0 {
		ctx.Source.UltimateUsed = true
		ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
			"source":        cardToInfo(ctx.Source),
			"effect":        "erebos_soul_chain_mark",
			"cast_player":   castPlayer,
			"marked_units":  markedUnits,
			"marked_spells": markedSpells,
		}})
	}
	return nil
}

func (Card2611101ErebosSoulChain) OnConsume(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Target == nil {
		return nil
	}
	ctx.Engine.weakenErebosSoulChainMarkedSpellsForUnit(ctx.Target)
	return nil
}

var _ OnSpellCastBehavior = Card2611101ErebosSoulChain{}

var _ OnConsumeBehavior = Card2611101ErebosSoulChain{}

func (e *Engine) weakenErebosSoulChainMarkedSpellsForUnit(unit *CardInstance) int {
	if e == nil || unit == nil || unit.Statuses[erebosSoulChainMarkedUnitStatus] <= 0 || unit.OwnerID < 0 || unit.OwnerID >= len(e.State.Players) {
		return 0
	}
	weakened := 0
	for _, card := range e.getAllFieldCards(e.State.Players[unit.OwnerID]) {
		if card == nil || card.Statuses[erebosSoulChainMarkedSpellStatus] <= 0 || !canInstanceBeWeakened(card) {
			continue
		}
		if e.addStatus(card, StatusWeaken, 1) {
			weakened++
		}
	}
	if weakened > 0 {
		e.emit(GameEvent{Type: "effect_trigger", Player: unit.OwnerID, Data: map[string]any{
			"effect": "erebos_soul_chain_weaken",
			"unit":   cardToInfo(unit),
			"count":  weakened,
		}})
	}
	return weakened
}
