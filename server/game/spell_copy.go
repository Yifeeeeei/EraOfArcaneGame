package game

import (
	"eraofarcane/model"
	"fmt"
)

// Copying changes the spell being resolved, while its physical source still
// owns payment, orientation and zone movement. It never consumes the copied card.
type SpellCopyBehavior interface {
	HasActiveSpellCopy(*CardInstance) bool
	SpellCopies(*EffectContext) []*CardInstance
	AutomaticSpellCopy() bool
}

func (AlwaysActive) HasActiveSpellCopy(*CardInstance) bool { return true }

type spellCopyUse struct {
	id      string
	purpose skillPurpose
}

func (e *Engine) validateCopiedOrLearnedSpell(spell, handSource *CardInstance, purpose skillPurpose) error {
	if handSource != nil {
		return e.validateHandSpellScrollForPurpose(spell, purpose)
	}
	return e.validateSkillForPurpose(spell, purpose)
}

func physicalSpellSource(card *CardInstance) *CardInstance {
	if card != nil && card.copySource != nil {
		return card.copySource
	}
	return card
}

func spellForm(card *CardInstance, forms map[string]*CardInstance) *CardInstance {
	if card != nil && forms[card.InstanceID] != nil {
		return forms[card.InstanceID]
	}
	return card
}

func (e *Engine) spellCopies(source *CardInstance, purpose skillPurpose) []*CardInstance {
	if source == nil || e.hasEffectiveStatus(source, StatusPetrify) {
		return nil
	}
	b, ok := cardBehavior(source).(SpellCopyBehavior)
	if !ok || !b.HasActiveSpellCopy(source) {
		return nil
	}
	var result []*CardInstance
	for _, card := range b.SpellCopies(e.skillContext(source.OwnerID, source)) {
		if card != nil && card.Card != nil && canUseSkillForPurpose(card.Card, purpose) {
			result = append(result, card)
		}
	}
	return result
}

func (e *Engine) copiedSpellForm(source, base *CardInstance) *CardInstance {
	form := snapshotFieldCard(base)
	form.InstanceID, form.OwnerID = source.InstanceID, source.OwnerID
	form.copySource, form.fieldLife = source, nil
	form.IsHorizontal = source.IsHorizontal
	form.usePayment, form.arcaneOnenessReveal = nil, nil
	return form
}

func copyActionData(action ActionMessage) map[string]any {
	data := map[string]any{}
	for k, v := range action.Data {
		data[k] = v
	}
	return data
}

func (e *Engine) prepareSpellCopies(playerID int, action ActionMessage, uses []spellCopyUse, retry func(ActionMessage) error) (map[string]*CardInstance, bool, error) {
	forms := map[string]*CardInstance{}
	selected, _ := action.Data["spell_copies"].(map[string]any)
	ps := e.State.Players[playerID]
	// Copy choices cannot change the physical cards' printed use costs.
	copyPresent := false
	plan := newActionCostPlan(ps)
	totalCost := map[string]int{}
	var paymentParts []spellPaymentPart
	for _, use := range uses {
		source := e.findSkill(ps, use.id)
		var cost map[string]int
		if source != nil {
			paymentParts = append(paymentParts, spellPaymentPart{source, cost})
			cost = plan.skillUseCost(e, source, use.purpose, nil)
		} else {
			source, _ = ps.FindHandCard(use.id)
			if source != nil {
				cost = plan.cardPlayCost(e, source)
			}
		}
		if source != nil {
			if _, ok := cardBehavior(source).(SpellCopyBehavior); ok {
				copyPresent = true
			}
		}
		totalCost = mergeElementCosts(totalCost, cost)
	}
	if copyPresent {
		overexertIDs, _ := action.Data["overexert_ids"].([]any)
		units, err := e.collectOverexertUnits(ps, stringsFromAnySlice(overexertIDs))
		if err != nil {
			return nil, false, err
		}
		if !e.canPayCostWithOverexertOptions(ps, totalCost, units, e.playerHasLightWildcard(ps)) {
			return nil, false, fmt.Errorf("not enough elements for copied spell")
		}
		available := e.availableElementsWithOverexert(ps, units)
		available[model.ElementArcane] += e.availableStrictArcaneWithOverexert(ps, units)
		if _, ok := e.constrainedSpellPayments(ps, paymentParts, available, nil, true); !ok {
			return nil, false, fmt.Errorf("not enough eligible elements for copied spell costs")
		}
	}
	for _, use := range uses {
		source := e.findSkill(ps, use.id)
		if source == nil {
			source, _ = ps.FindHandCard(use.id)
		}
		if source == nil {
			continue
		} // The ordinary collector reports bad IDs.
		b, ok := cardBehavior(source).(SpellCopyBehavior)
		if !ok {
			continue
		}
		if source.Card.IsSkill() {
			if err := e.validateReadySkill(source); err != nil {
				return nil, false, err
			}
		}
		options := e.spellCopies(source, use.purpose)
		boosts, _ := action.Data["boost_ids"].([]any)
		if len(boosts) > 0 && (use.purpose == skillPurposeAttack || use.purpose == skillPurposeDefend) {
			compatible := options[:0]
			for _, option := range options {
				if canSkillBeBoosted(option) {
					compatible = append(compatible, option)
				}
			}
			options = compatible
		}
		if len(options) == 0 {
			return nil, false, fmt.Errorf("no eligible spell to copy")
		}
		id, _ := selected[use.id].(string)
		if id == "" && len(options) == 1 && b.AutomaticSpellCopy() {
			id = options[0].InstanceID
		}
		if id == "" {
			candidates := make([]map[string]any, 0, len(options))
			for _, card := range options {
				candidates = append(candidates, candidateInfo(card, "skill", "own"))
			}
			e.SetPendingActionWithError(playerID, "spell_copy_choice", source.Card.Name+":选择本次复制的法术", candidates, 1, 1, nil, false, func(ids []string, _ map[string]any) error {
				data := copyActionData(action)
				choices := map[string]any{}
				for key, value := range selected {
					choices[key] = value
				}
				choices[use.id] = firstSelected(ids)
				data["spell_copies"] = choices
				return retry(ActionMessage{Action: action.Action, Data: data})
			})
			return nil, true, nil
		}
		var base *CardInstance
		for _, candidate := range options {
			if candidate.InstanceID == id {
				base = candidate
				break
			}
		}
		if base == nil {
			return nil, false, fmt.Errorf("invalid copied spell")
		}
		forms[use.id] = e.copiedSpellForm(source, base)
	}
	return forms, false, nil
}

func (e *Engine) chooseCopiedAttackTarget(playerID int, spell *CardInstance, action ActionMessage, retry func(ActionMessage) error) (bool, error) {
	if spell.copySource == nil || !skillNeedsTargetInstance(spell) || action.Data["target_type"] == "unit" {
		return false, nil
	}
	candidates := e.spellTargetCandidates(playerID, spell)
	if len(candidates) == 0 {
		return false, fmt.Errorf("copied spell has no legal target")
	}
	e.SetPendingActionWithError(playerID, "copied_spell_target", "选择复制法术「"+spell.Card.Name+"」的目标", candidates, 1, 1, nil, false, func(ids []string, _ map[string]any) error {
		target := selectedSpellTargetFromCandidates(e, playerID, spell, firstSelected(ids), candidates)
		if target == nil {
			return fmt.Errorf("invalid copied spell target")
		}
		data := copyActionData(action)
		data["target_type"], data["target_col"], data["target_row"] = target.Type, float64(target.Position.Col), float64(target.Position.Row)
		return retry(ActionMessage{Action: action.Action, Data: data})
	})
	return true, nil
}

func (e *Engine) chooseCopiedItemBoosts(playerID int, spell *CardInstance, action ActionMessage) bool {
	if _, supplied := action.Data["boost_ids"]; supplied || !canSkillBeBoosted(spell) {
		return false
	}
	ps := e.State.Players[playerID]
	candidates := []map[string]any{}
	for _, card := range e.ruleSources(ps) {
		if card == nil || card.InstanceID == spell.InstanceID || !card.Card.IsSkill() {
			continue
		}
		if e.validateSkillForPurpose(card, skillPurposeAttackBoost) == nil || len(e.spellCopies(card, skillPurposeAttackBoost)) > 0 {
			candidates = append(candidates, candidateInfo(card, "skill", "own"))
		}
	}
	for _, card := range ps.Hand {
		if card.InstanceID == spell.InstanceID || !isSpellScrollCard(card.Card) {
			continue
		}
		if e.validateHandSpellScrollForPurpose(card, skillPurposeAttackBoost) == nil || len(e.spellCopies(card, skillPurposeAttackBoost)) > 0 {
			candidates = append(candidates, candidateInfo(card, "hand", "own"))
		}
	}
	if len(candidates) == 0 {
		return false
	}
	e.SetPendingActionWithError(playerID, "copied_spell_boosts", "为复制的法术选择强化，也可跳过", candidates, 0, len(candidates), nil, false, func(ids []string, _ map[string]any) error {
		data := copyActionData(action)
		values := []any{}
		for _, id := range ids {
			values = append(values, id)
		}
		data["boost_ids"] = values
		delete(data, "payment")
		return e.handleCastSpell(playerID, ActionMessage{Data: data})
	})
	return true
}

func hasCopiedSpellUse(skill *CardInstance, boosts []*CardInstance) bool {
	if skill != nil && skill.copySource != nil {
		return true
	}
	for _, boost := range boosts {
		if boost.copySource != nil {
			return true
		}
	}
	return false
}

func (e *Engine) retryCopiedSpellPayment(playerID int, action ActionMessage, cost, available map[string]int, retry func(ActionMessage) error) {
	e.SetPendingActionWithErrorAndContext(playerID, "copied_spell_payment", "实际支付不满足卡牌约束，请重新选择元素", []map[string]any{{"instance_id": "pay", "name": "重新支付", "zone": "choice"}}, 1, 1, cost, false, map[string]any{"payment_available": available}, func(_ []string, choice map[string]any) error {
		data := copyActionData(action)
		data["payment"] = choice["payment"]
		data["copy_payment_done"] = true
		return retry(ActionMessage{Action: action.Action, Data: data})
	})
}
