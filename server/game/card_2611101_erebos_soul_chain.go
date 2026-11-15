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
		unit.Statuses[erebosSoulChainMarkedUnitStatus]++
		unit.Statuses[soulChainUnitKey(ctx.Source)] = 1
		markedUnits++
	}
	if markedUnits == 0 {
		return nil
	}
	markedSpells := 0
	for _, spell := range append([]*CardInstance{ctx.Target}, spellInstancesFromData(ctx.ExtraData, "boost_skills")...) {
		if spell == nil || spell.Card == nil || spell.OwnerID != castPlayer || !isCombatSpellCard(spell.Card) {
			continue
		}
		if spell.Statuses == nil {
			spell.Statuses = make(map[string]int)
		}
		spell.Statuses[erebosSoulChainMarkedSpellStatus]++
		spell.Statuses[soulChainSpellKey(ctx.Source)] = 1
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
	ctx.Engine.weakenSoulChainGroup(ctx.Source, ctx.Target)
	return nil
}

var _ OnSpellCastBehavior = Card2611101ErebosSoulChain{}

var _ OnConsumeBehavior = Card2611101ErebosSoulChain{}

func soulChainUnitKey(source *CardInstance) string  { return "soul_chain_unit:" + source.InstanceID }
func soulChainSpellKey(source *CardInstance) string { return "soul_chain_spell:" + source.InstanceID }

func (e *Engine) weakenErebosSoulChainMarkedSpellsForUnit(unit *CardInstance) int {
	if e == nil || unit == nil {
		return 0
	}
	total := 0
	for _, ps := range e.State.Players {
		for _, source := range e.ruleSources(ps) {
			if source.Card != nil && source.Card.Number == "2611101" && !e.hasEffectiveStatus(source, StatusPetrify) {
				total += e.weakenSoulChainGroup(source, unit)
			}
		}
	}
	return total
}
func (e *Engine) weakenSoulChainGroup(source, unit *CardInstance) int {
	if source == nil || unit == nil || unit.Statuses[soulChainUnitKey(source)] <= 0 {
		return 0
	}
	count := 0
	for _, ps := range e.State.Players {
		for _, spell := range e.ruleSources(ps) {
			if spell.Statuses[soulChainSpellKey(source)] > 0 && canInstanceBeWeakened(spell) && e.addStatus(spell, StatusWeaken, 1) {
				count++
			}
		}
	}
	return count
}
