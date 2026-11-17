package game

import (
	"eraofarcane/model"
	"strings"
)

func addCardToDeck(ctx *EffectContext, cardNumber string, count int) {
	card := getCardDB()[cardNumber]
	if card == nil || count <= 0 {
		return
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	for i := 0; i < count; i++ {
		ps.Deck = append(ps.Deck, ctx.Engine.newCardInstance(card, ctx.PlayerID, ctx.Engine.State.TurnNumber))
	}
	ctx.Engine.shuffleDeck(ctx.PlayerID)
	emitBatchEffect(ctx, "add_card_to_deck")
}

func addCardToOpponentDeck(ctx *EffectContext, cardNumber string, count int) {
	card := getCardDB()[cardNumber]
	if card == nil || count <= 0 {
		return
	}
	ps := ctx.Engine.State.Players[ctx.OpponentID]
	for i := 0; i < count; i++ {
		ps.Deck = append(ps.Deck, ctx.Engine.newCardInstance(card, ctx.OpponentID, ctx.Engine.State.TurnNumber))
	}
	ctx.Engine.shuffleDeck(ctx.OpponentID)
	emitBatchEffect(ctx, "add_card_to_opponent_deck")
}

func firstActiveCardByNumber(e *Engine, ps *PlayerState, number string) *CardInstance {
	for _, fieldCard := range e.getAllFieldCards(ps) {
		if fieldCard != nil && fieldCard.Card != nil && fieldCard.Card.Number == number && !e.hasEffectiveStatus(fieldCard, StatusPetrify) {
			return fieldCard
		}
	}
	return nil
}

func (e *Engine) promptHeartPiercerAfterPhantomPain(ctx *EffectContext) {
	if ctx == nil || ctx.Engine == nil {
		return
	}
	owner := ctx.PlayerID
	if firstActiveCardByNumber(e, e.State.Players[owner], "1611003") == nil {
		return
	}
	candidates := e.enemySpellCardCandidates(owner)
	if len(candidates) == 0 {
		return
	}
	opponent := e.State.Players[1-owner]
	e.SetPendingAction(owner, "heart_piercer_phantom_pain_extra", "\"穿心人\":幻痛触发,可以额外选择1个敌方法术虚弱2", candidates, 0, 1, func(selected []string) {
		if len(selected) == 0 {
			return
		}
		target := e.findSkill(opponent, selected[0])
		if canInstanceBeWeakened(target) {
			e.addStatus(target, StatusWeaken, 2)
		}
	})
}

func (e *Engine) enemySpellCardCandidates(playerID int) []map[string]any {
	ps := e.State.Players[1-playerID]
	candidates := make([]map[string]any, 0)
	for _, skill := range ps.Skills {
		if canInstanceBeWeakened(skill) {
			candidates = append(candidates, candidateInfo(skill, "skill", "enemy"))
		}
	}
	for _, card := range e.getAllFieldCards(ps) {
		if card == nil {
			continue
		}
		for _, skill := range card.BoundSkills {
			if canInstanceBeWeakened(skill) {
				candidates = append(candidates, candidateInfo(skill, "bound_skill", "enemy"))
			}
		}
	}
	return candidates
}

func voodooDollLinkCandidates(ctx *EffectContext) []map[string]any {
	candidates := make([]map[string]any, 0)
	for ownerID, ps := range ctx.Engine.State.Players {
		side := "enemy"
		if ownerID == ctx.PlayerID {
			side = "own"
		}
		for col := 0; col < 3; col++ {
			for row := 0; row < 3; row++ {
				unit := ps.Units[col][row]
				if unit == nil || !unit.Card.IsCompanion() {
					continue
				}
				if ownerID != ctx.PlayerID && !ctx.Engine.IsInSpellRange(ctx.PlayerID, col, row, false) {
					continue
				}
				candidates = append(candidates, candidateInfo(unit, "unit", side))
			}
		}
	}
	return candidates
}

func clearVoodooDollLinks(doll *CardInstance) {
	for status := range doll.Statuses {
		if strings.HasPrefix(status, "巫毒连结:") {
			delete(doll.Statuses, status)
		}
	}
}

func voodooDollIsLinked(doll *CardInstance, unit *CardInstance) bool {
	return doll != nil && unit != nil && doll.Statuses["巫毒连结:"+unit.InstanceID] > 0
}

func voodooDollOtherLinkedUnit(e *Engine, doll *CardInstance, damaged *CardInstance) *CardInstance {
	for status := range doll.Statuses {
		if !strings.HasPrefix(status, "巫毒连结:") {
			continue
		}
		id := strings.TrimPrefix(status, "巫毒连结:")
		if id == "" || damaged != nil && id == damaged.InstanceID {
			continue
		}
		if unit := e.findUnitByInstanceID(id); unit != nil {
			return unit
		}
	}
	return nil
}

func voodooDollSide(playerID int, unit *CardInstance) string {
	if unit != nil && unit.OwnerID == playerID {
		return "own"
	}
	return "enemy"
}

const shadowCloakUsedStatus = "shadow_cloak_used"

func positiveSpellPowerSources(spell *SpellCast) []SpellPowerSource {
	if spell == nil {
		return nil
	}
	sources := make([]SpellPowerSource, 0, len(spell.PowerSources))
	for _, source := range spell.PowerSources {
		if source.Power > 0 && source.InstanceID != "" {
			sources = append(sources, source)
		}
	}
	if len(sources) == 0 && spell.TotalPower > 0 {
		source := spellPowerSourceForCard(spell.Skill, spell.TotalPower, true)
		if source.InstanceID == "" {
			source.InstanceID = "__total_power__"
		}
		sources = append(sources, source)
	}
	return sources
}

func applyIceDissolveToSource(ctx *EffectContext, spell *SpellCast, instanceID string) {
	if ctx == nil || spell == nil || instanceID == "" {
		return
	}
	removed := 0
	for i := range spell.PowerSources {
		if spell.PowerSources[i].InstanceID != instanceID {
			continue
		}
		removed = spell.PowerSources[i].Power
		spell.PowerSources[i].Power = 0
		break
	}
	if removed == 0 && ((spell.Skill != nil && spell.Skill.InstanceID == instanceID) || instanceID == "__total_power__") {
		removed = spell.TotalPower
	}
	spell.TotalPower = max(spell.TotalPower-removed, 0)
	ctx.Engine.emit(GameEvent{
		Type:   "spell_reaction",
		Player: -1,
		Data: map[string]any{
			"player":             ctx.PlayerID,
			"card":               cardToInfo(ctx.Source),
			"effect":             "power_zero",
			"power":              spell.TotalPower,
			"power_source_id":    instanceID,
			"power_source_power": removed,
		},
	})
}

func (e *Engine) summonWaterPhantomCopy(playerID int, targetID string, pos Position) *CardInstance {
	if e == nil || !pos.Valid() || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	ps := e.State.Players[playerID]
	if ps.Units[pos.Col][pos.Row] != nil {
		return nil
	}
	source := e.findUnitByInstanceID(targetID)
	if source == nil || source.OwnerID != playerID || source.Card == nil || !source.Card.IsCompanion() || source.Card.Category != model.ElementWater || source.EnterTurn != e.State.TurnNumber {
		return nil
	}
	cardCopy := *source.Card
	cardCopy.Life = 1
	copyUnit := e.newCardInstance(&cardCopy, playerID, e.State.TurnNumber)
	copyUnit.Position = &Position{Col: pos.Col, Row: pos.Row}
	copyUnit.IsHorizontal = true
	copyUnit.CurrentLife = 1
	copyUnit.Statuses["水幻影复制"] = 1
	ps.Units[pos.Col][pos.Row] = copyUnit
	e.ApplySummonModifiersOnEnter(copyUnit)
	e.emit(GameEvent{
		Type:   "summon",
		Player: -1,
		Data: map[string]any{
			"player":   playerID,
			"card":     cardToInfo(copyUnit),
			"position": pos,
			"via":      "water_phantom",
		},
	})
	e.triggerEffects(TriggerOnEnter, copyUnit, nil, nil)
	enterData := map[string]any{"entered_player": playerID}
	e.notifyCardEntered(playerID, copyUnit, enterData)
	e.triggerFieldEffectsWithData(TriggerOnUnitEnter, playerID, copyUnit, enterData)
	e.triggerFieldEffectsWithData(TriggerOnUnitEnter, 1-playerID, copyUnit, enterData)
	return copyUnit
}

func (e *Engine) spellAffectedUnitsWithExtraTargets(defenderID int, skill *CardInstance, target SpellTarget, extraTargets []SpellTarget) []*CardInstance {
	affected := e.spellAffectedUnits(defenderID, skill, target)
	for _, extraTarget := range extraTargets {
		if extraTarget.Type != "unit" || !extraTarget.Position.Valid() {
			continue
		}
		targetOwnerID := defenderID
		if extraTarget.OwnerID != nil {
			targetOwnerID = *extraTarget.OwnerID
		}
		if targetOwnerID < 0 || targetOwnerID >= len(e.State.Players) {
			continue
		}
		unit := e.State.Players[targetOwnerID].Units[extraTarget.Position.Col][extraTarget.Position.Row]
		if unit == nil {
			continue
		}
		alreadyIncluded := false
		for _, existing := range affected {
			if existing == unit {
				alreadyIncluded = true
				break
			}
		}
		if !alreadyIncluded {
			affected = append(affected, unit)
		}
	}
	return affected
}

func ailimerShuffleShacklesOnce(ctx *EffectContext) {
	if ctx == nil || ctx.Source == nil || ctx.Source.Statuses["桎梏已洗入"] > 0 {
		return
	}
	addCardToOpponentDeck(ctx, "2501001", 5)
	ctx.Source.Statuses["桎梏已洗入"] = 1
}

func ailimerUnlockIfShacklesCleared(ctx *EffectContext) {
	if ctx == nil || ctx.Source == nil || ctx.Source.Statuses["爱里默解放"] > 0 {
		return
	}
	if countCardInstancesByNumber(ctx.Engine.State.Players[ctx.OpponentID].Graveyard, "2501001") < 5 {
		return
	}
	ctx.Source.Statuses["爱里默解放"] = 1
	ctx.Engine.emit(GameEvent{
		Type:   "effect_trigger",
		Player: -1,
		Data: map[string]any{
			"player": ctx.PlayerID,
			"card":   cardToInfo(ctx.Source),
			"effect": "ailimer_unlocked",
		},
	})
}

func countCardInstancesByNumber(cards []*CardInstance, number string) int {
	count := 0
	for _, card := range cards {
		if card != nil && card.Card != nil && card.Card.Number == number {
			count++
		}
	}
	return count
}
