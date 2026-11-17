package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card3411102BloodSandArray struct{ AlwaysActive }

func (Card3411102BloodSandArray) ID() string { return "3411102" }

func (Card3411102BloodSandArray) Name() string { return "蔽天阵 血沙" }

func (Card3411102BloodSandArray) PerTurnLabel(*CardInstance) string {
	return "主动"
}

func (Card3411102BloodSandArray) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	e, source := ctx.Engine, ctx.Source
	e.promptBloodSandPayment(ctx.PlayerID, source, "blood_sand_array_pay", func(first int) {
		e.promptBloodSandPayment(ctx.OpponentID, source, "blood_sand_array_pay_opponent", func(second int) {
			if !e.cardStillOnField(source) {
				return
			}
			diff := bloodSandAbsDiff(first - second)
			source.Statuses[bloodSandArrayMarkerStatus] += diff
			e.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{"source": cardToInfo(source), "effect": "blood_sand_array_markers", "paid_owner": first, "paid_enemy": second, "markers_add": diff, "markers": source.Statuses[bloodSandArrayMarkerStatus]}})
		})
	})
	return nil
}

func (e *Engine) promptBloodSandPayment(playerID int, source *CardInstance, kind string, after func(int)) {
	candidates := e.bloodSandPaymentCandidates(playerID)
	e.SetPendingActionWithError(playerID, kind, "蔽天阵 血沙:选择最多3点己方单位负载或生命支付", candidates, 0, min(3, len(candidates)), nil, false, func(selected []string, data map[string]any) error {
		if !e.cardStillOnField(source) {
			return nil
		}
		paid, err := e.applyBloodSandPayments(playerID, selected, data, source)
		if err != nil {
			return err
		}
		e.runResolution("blood sand payment", e.settleTurnDeaths, func() { after(paid) })
		return nil
	})
}

func (Card3411102BloodSandArray) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Source == nil {
		return
	}
	markers := ctx.Source.Statuses[bloodSandArrayMarkerStatus]
	if markers <= 0 {
		return
	}
	stats.PowerBonus += markers * 3
	stats.DamageBonus += markers
}

var _ PerTurnAbility = Card3411102BloodSandArray{}

var _ SkillContributionModifier = Card3411102BloodSandArray{}

const bloodSandArrayMarkerStatus = "蔽天阵 血沙标记物"

func (e *Engine) bloodSandPaymentCandidates(playerID int) []map[string]any {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	candidates := make([]map[string]any, 0)
	for _, unit := range e.getAllFieldCards(e.State.Players[playerID]) {
		if unit == nil || unit.Card == nil || unit.Position == nil {
			continue
		}
		if e.bloodSandPayablePoints(unit) <= 0 {
			continue
		}
		info := candidateInfo(unit, "unit", "own")
		info["payable_load"] = dragonBloodTreantReducibleLoad(unit)
		candidates = append(candidates, info)
	}
	return candidates
}

func (e *Engine) applyBloodSandPayments(playerID int, selected []string, data map[string]any, source *CardInstance) (int, error) {
	requests := bloodSandPaymentRequests(selected, data)
	selectedIDs := map[string]bool{}
	for _, id := range selected {
		if selectedIDs[id] {
			return 0, fmt.Errorf("duplicate payment target")
		}
		selectedIDs[id] = true
	}
	if raw, exists := data["payments"]; exists {
		entries, ok := raw.([]any)
		if !ok || len(entries) != len(requests) {
			return 0, fmt.Errorf("invalid blood sand payments")
		}
		for _, entry := range entries {
			fields, ok := entry.(map[string]any)
			if !ok {
				return 0, fmt.Errorf("invalid blood sand payment")
			}
			if rawAmount, supplied := fields["amount"]; supplied {
				if amount, valid := intFromAny(rawAmount); !valid || amount <= 0 || amount > 3 {
					return 0, fmt.Errorf("invalid blood sand payment amount")
				}
			}
		}
	}
	type payment struct {
		unit    *CardInstance
		element string
	}
	plan := []payment{}
	life := map[*CardInstance]int{}
	loads := map[*CardInstance]map[string]int{}
	usedIDs := map[string]bool{}
	for _, request := range requests {
		if !selectedIDs[request.instanceID] || request.amount <= 0 || len(plan)+request.amount > 3 || (request.mode != "" && request.mode != "auto" && request.mode != "load" && request.mode != "life") {
			return 0, fmt.Errorf("invalid blood sand payment")
		}
		unit := e.findOwnUnitByInstanceID(playerID, request.instanceID)
		if unit == nil {
			return 0, fmt.Errorf("invalid blood sand payment target")
		}
		usedIDs[request.instanceID] = true
		if loads[unit] == nil {
			loads[unit] = copyElements(dragonBloodTreantReducibleLoad(unit))
			life[unit] = unit.CurrentLife
		}
		for i := 0; i < request.amount; i++ {
			element := ""
			if request.element != "" {
				if request.mode == "life" || loads[unit][request.element] <= 0 {
					return 0, fmt.Errorf("selected load element is unavailable")
				}
				element = request.element
				loads[unit][element]--
			} else if request.mode != "life" {
				for _, elem := range model.AllElements {
					if loads[unit][elem] > 0 {
						element = elem
						loads[unit][elem]--
						break
					}
				}
			}
			if element == "" {
				if request.mode == "load" || life[unit] <= 0 {
					return 0, fmt.Errorf("selected unit cannot pay load or life")
				}
				life[unit]--
			}
			plan = append(plan, payment{unit, element})
		}
	}
	if len(usedIDs) != len(selectedIDs) {
		return 0, fmt.Errorf("payment targets do not match selection")
	}
	for _, p := range plan {
		if p.element == "" {
			p.unit.CurrentLife--
		} else {
			e.reduceCardElementLoadWithTriggers(playerID, p.unit, p.element, 1, source)
		}
	}
	for unit := range life {
		if unit.CurrentLife <= 0 && !unit.Card.IsHero() {
			e.queueDeathWithData(unit, playerID, map[string]any{"death_cause": "life_payment"})
		}
	}
	return len(plan), nil
}

type bloodSandPaymentRequest struct {
	instanceID string
	amount     int
	mode       string
	element    string
}

func bloodSandPaymentRequests(selected []string, data map[string]any) []bloodSandPaymentRequest {
	if data != nil {
		if raw, ok := data["payments"].([]any); ok {
			requests := make([]bloodSandPaymentRequest, 0, len(raw))
			for _, entry := range raw {
				m, ok := entry.(map[string]any)
				if !ok {
					continue
				}
				id, _ := m["instance_id"].(string)
				if id == "" {
					continue
				}
				amount := intFromData(m, "amount", 1)
				if amount <= 0 {
					continue
				}
				mode, _ := m["mode"].(string)
				element, _ := m["element"].(string)
				requests = append(requests, bloodSandPaymentRequest{instanceID: id, amount: amount, mode: mode, element: element})
			}
			return requests
		}
	}
	requests := make([]bloodSandPaymentRequest, 0, len(selected))
	for _, id := range selected {
		if id != "" {
			requests = append(requests, bloodSandPaymentRequest{instanceID: id, amount: 1})
		}
	}
	return requests
}

func (e *Engine) findOwnUnitByInstanceID(playerID int, instanceID string) *CardInstance {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) || instanceID == "" {
		return nil
	}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := e.State.Players[playerID].Units[col][row]
			if unit != nil && unit.InstanceID == instanceID {
				return unit
			}
		}
	}
	return nil
}

func (e *Engine) bloodSandPayablePoints(unit *CardInstance) int {
	if unit == nil {
		return 0
	}
	return min(3, totalElementCost(dragonBloodTreantReducibleLoad(unit))+max(unit.CurrentLife, 0))
}

func (e *Engine) payOneBloodSandLoad(playerID int, unit *CardInstance, source *CardInstance) bool {
	if unit == nil {
		return false
	}
	for _, elem := range model.AllElements {
		if e.effectiveElementsGain(unit)[elem] <= 0 {
			continue
		}
		return e.reduceCardElementLoadWithTriggers(playerID, unit, elem, 1, source) == 1
	}
	return false
}

func bloodSandPayLife(unit *CardInstance) bool {
	if unit == nil || unit.CurrentLife <= 0 {
		return false
	}
	unit.CurrentLife--
	return true
}

func bloodSandAbsDiff(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
