package game

import (
	"eraofarcane/model"
	"fmt"
)

// Only cards whose result depends on the actual elements need attribution.
type UsePaymentAttributionBehavior interface {
	NeedsUsePaymentAttribution(*CardInstance) bool
}

func paymentAttributions(action ActionMessage) map[string]map[string]int {
	result := map[string]map[string]int{}
	if raw, ok := action.Data["payment_attributions"].(map[string]any); ok {
		for id, value := range raw {
			if p := paymentFromAction(ActionMessage{Data: map[string]any{"payment": value}}); p != nil {
				result[id] = p
			}
		}
	}
	return result
}

func (e *Engine) prepareUsePaymentAttributions(ps *PlayerState, parts []spellPaymentPart, payment map[string]int, action ActionMessage, retry func(ActionMessage) error) (bool, error) {
	assigned := paymentAttributions(action)
	if _, ok := e.constrainedSpellPayments(ps, parts, payment, assigned); !ok {
		return false, fmt.Errorf("实际支付无法满足各法术的费用约束")
	}
	remaining := copyElements(payment)
	for _, part := range parts {
		reserved, chosen := assigned[part.card.InstanceID]
		if !chosen {
			required := strictPaymentRequirement(physicalSpellSource(part.card), paymentPurposeUse, part.cost)
			if totalElementCost(required) == totalElementCost(part.cost) {
				reserved = required
			}
		}
		for elem, n := range reserved {
			remaining[elem] -= n
		}
	}
	for _, part := range parts {
		if _, ok := assigned[part.card.InstanceID]; ok {
			continue
		}
		b, ok := cardBehavior(part.card).(UsePaymentAttributionBehavior)
		if !ok || !b.NeedsUsePaymentAttribution(part.card) {
			continue
		}
		n := totalElementCost(part.cost)
		if n <= 0 || remaining[model.ElementArcane] < n {
			continue
		}
		legal := func(own map[string]int) bool {
			if !validateElementPaymentWithOptions(remaining, part.cost, own, e.playerHasLightWildcard(ps), e.playerHasLightCostWildcard(ps)) {
				return false
			}
			proposed := map[string]map[string]int{}
			for id, paid := range assigned {
				proposed[id] = paid
			}
			proposed[part.card.InstanceID] = own
			_, ok := e.constrainedSpellPayments(ps, parts, payment, proposed)
			return ok
		}
		if !legal(map[string]int{model.ElementArcane: n}) {
			continue
		}
		ambiguous := false
		for _, elem := range model.AllElements {
			if elem != model.ElementArcane && remaining[elem] > 0 && legal(map[string]int{model.ElementArcane: n - 1, elem: 1}) {
				ambiguous = true
				break
			}
		}
		if !ambiguous {
			continue
		}
		e.SetPendingActionWithErrorAndContext(ps.PlayerID, "spell_payment_attribution", "选择实际用于「"+part.card.Card.Name+"」的元素，其余用于其他法术", []map[string]any{{"instance_id": part.card.InstanceID, "name": part.card.Card.Name, "zone": "choice"}}, 1, 1, part.cost, false, map[string]any{"payment_available": remaining}, func(_ []string, choice map[string]any) error {
			own := paymentFromAction(ActionMessage{Data: choice})
			if own == nil {
				return fmt.Errorf("请选择实际支付给此卡的元素")
			}
			if !legal(own) {
				return fmt.Errorf("所选支付无法同时满足其余法术费用")
			}
			data := copyActionData(action)
			rows := map[string]any{}
			for id, paid := range assigned {
				rows[id] = paid
			}
			rows[part.card.InstanceID] = own
			data["payment_attributions"] = rows
			return retry(ActionMessage{Action: action.Action, Data: data})
		})
		return true, nil
	}
	return false, nil
}
