package game

import (
	"eraofarcane/model"
	"fmt"
)

func spellHitAffectedUnitsFromData(ctx *EffectContext) []*CardInstance {
	if ctx == nil {
		return nil
	}
	if ctx.ExtraData != nil {
		if affected, ok := ctx.ExtraData["affected_units"].([]*CardInstance); ok && len(affected) > 0 {
			return affected
		}
	}
	if ctx.Target != nil {
		return []*CardInstance{ctx.Target}
	}
	return nil
}

func (e *Engine) enemyCompanionsInSpellRange(playerID int) []map[string]any {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	opponent := e.State.Players[1-playerID]
	candidates := make([]map[string]any, 0)
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := opponent.Units[col][row]
			if unit == nil || unit.Card == nil || !unit.Card.IsCompanion() {
				continue
			}
			if !e.IsInSpellRange(playerID, col, row, false) {
				continue
			}
			candidates = append(candidates, candidateInfo(unit, "unit", "enemy"))
		}
	}
	return candidates
}

func (e *Engine) triggerMagicMothAfterFocusSpellCast(playerID int, skill *CardInstance) {
	if e == nil || skill == nil || skill.Card == nil || !hasCardTag(skill.Card, "聚能") || playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	ps := e.State.Players[playerID]
	for _, card := range ps.Deck {
		if card == nil || card.Card == nil || card.Card.Number != "1021113" {
			continue
		}
		moth := card
		e.SetPendingAction(playerID, "magic_moth_draw",
			"魔法飞蛾:是否从卡组抽取本卡", []map[string]any{candidateInfo(moth, "deck", "own")}, 0, 1,
			func(selected []string) {
				if len(selected) == 0 {
					return
				}
				for i, current := range ps.Deck {
					if current != moth {
						continue
					}
					ps.Deck = append(ps.Deck[:i], ps.Deck[i+1:]...)
					e.appendCardsToHand(playerID, []*CardInstance{moth})
					e.notifyCardDrawn(playerID, moth)
					e.shuffleDeck(playerID)
					e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
						"effect": "magic_moth_draw",
						"source": cardToInfo(skill),
						"card":   cardToInfo(moth),
					}})
					e.enforceImmediateHandLimitAfterHandGain(playerID)
					return
				}
			})
		return
	}
}

func addSoulMarkerToSpell(skill *CardInstance) {
	if skill == nil || skill.Card == nil || !isSpellLikeCard(skill.Card) {
		return
	}
	skill.Statuses[soulMarkerStatus]++
	skill.PowerBonus += 2
}

func (e *Engine) applyCoralBellyFirstSpellAttackBonus(playerID int, skill *CardInstance) {
	if e == nil || skill == nil || skill.Card == nil || !isSpellLikeCard(skill.Card) || isSorcerySkill(skill.Card) {
		return
	}
	ps := e.State.Players[playerID]
	if ps == nil {
		return
	}
	for _, card := range e.getAllFieldCards(ps) {
		if card == nil || card.Card == nil || card.Card.Number != "4211101" || e.hasEffectiveStatus(card, StatusPetrify) {
			continue
		}
		if card.Statuses[coralBellyFirstSpellAttackUsedStatus] > 0 {
			continue
		}
		card.Statuses[coralBellyFirstSpellAttackUsedStatus] = 1
		skill.PowerBonus += 3
		e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
			"source": cardToInfo(card),
			"target": cardToInfo(skill),
			"effect": "first_spell_attack_power_bonus",
			"amount": 3,
		}})
		return
	}
}

func selectedSpellTargetFromCandidates(e *Engine, playerID int, skill *CardInstance, instanceID string, candidates []map[string]any) *SpellTarget {
	if e == nil || skill == nil || instanceID == "" {
		return nil
	}
	for _, candidate := range candidates {
		id, _ := candidate["instance_id"].(string)
		if id != instanceID {
			continue
		}
		target := SpellTarget{Type: "unit"}
		if pos, ok := candidate["position"].(Position); ok {
			target.Position = pos
		} else {
			unit := e.findUnitByInstanceID(instanceID)
			if unit == nil || unit.Position == nil {
				return nil
			}
			target.Position = *unit.Position
		}
		if owner, ok := candidate["target_owner"].(int); ok {
			target.OwnerID = &owner
		}
		if err := e.validateSpellTarget(playerID, skill, target); err != nil {
			return nil
		}
		return &target
	}
	return nil
}

func (e *Engine) validateSpellPowerSacrifice(playerID int, skill *CardInstance, action ActionMessage) (*CardInstance, int, error) {
	if e == nil || skill == nil || skill.Card == nil || skill.Card.Number != "3121104" {
		return nil, 0, nil
	}
	sacrificeID, _ := action.Data["sacrifice_id"].(string)
	if sacrificeID == "" {
		return nil, 0, nil
	}
	target, zone := e.findFriendlyCandidate(playerID, sacrificeID)
	if zone != "unit" || target == nil || target.Card == nil || !target.Card.IsCompanion() || target.Card.Category != model.ElementFire {
		return nil, 0, fmt.Errorf("3121104 requires a friendly fire companion sacrifice")
	}
	return target, totalElementCost(target.Card.ElementsCost), nil
}

func (e *Engine) validateSpellPowerSacrificeForSources(playerID int, sources []*CardInstance, action ActionMessage) (*CardInstance, *CardInstance, int, error) {
	for _, source := range sources {
		sacrifice, bonus, err := e.validateSpellPowerSacrifice(playerID, source, action)
		if err != nil {
			return nil, nil, 0, err
		}
		if sacrifice != nil && bonus > 0 {
			return sacrifice, source, bonus, nil
		}
	}
	return nil, nil, 0, nil
}
