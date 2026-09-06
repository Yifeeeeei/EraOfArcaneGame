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
	source := ctx.Source
	firstCandidates := ctx.Engine.bloodSandPaymentCandidates(ctx.PlayerID)
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "blood_sand_array_pay",
		"蔽天阵 血沙:选择最多3点己方单位负载或生命支付", firstCandidates, 0, min(3, len(firstCandidates)), nil, false,
		func(selected []string, data map[string]any) error {
			if !ctx.Engine.cardStillOnField(source) {
				return nil
			}
			firstPaid, err := ctx.Engine.applyBloodSandPayments(ctx.PlayerID, selected, data, source)
			if err != nil {
				return err
			}
			opponentID := 1 - ctx.PlayerID
			secondCandidates := ctx.Engine.bloodSandPaymentCandidates(opponentID)
			ctx.Engine.SetPendingActionWithError(opponentID, "blood_sand_array_pay_opponent",
				"蔽天阵 血沙:选择最多3点己方单位负载或生命支付", secondCandidates, 0, min(3, len(secondCandidates)), nil, false,
				func(opponentSelected []string, opponentData map[string]any) error {
					if !ctx.Engine.cardStillOnField(source) {
						return nil
					}
					secondPaid, err := ctx.Engine.applyBloodSandPayments(opponentID, opponentSelected, opponentData, source)
					if err != nil {
						return err
					}
					diff := bloodSandAbsDiff(firstPaid - secondPaid)
					if diff > 0 {
						source.Statuses[bloodSandArrayMarkerStatus] += diff
					}
					ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
						"source":      cardToInfo(source),
						"effect":      "blood_sand_array_markers",
						"paid_owner":  firstPaid,
						"paid_enemy":  secondPaid,
						"markers_add": diff,
						"markers":     source.Statuses[bloodSandArrayMarkerStatus],
					}})
					return nil
				})
			return nil
		})
	return nil
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
		candidates = append(candidates, candidateInfo(unit, "unit", "own"))
	}
	return candidates
}

func (e *Engine) applyBloodSandPayments(playerID int, selected []string, data map[string]any, source *CardInstance) (int, error) {
	requests := bloodSandPaymentRequests(selected, data)
	if len(requests) == 0 {
		return 0, nil
	}
	total := 0
	for _, request := range requests {
		if total >= 3 {
			break
		}
		unit := e.findOwnUnitByInstanceID(playerID, request.instanceID)
		if unit == nil {
			return total, fmt.Errorf("invalid blood sand payment target")
		}
		amount := min(request.amount, 3-total)
		for i := 0; i < amount; i++ {
			if request.mode == "life" {
				if !bloodSandPayLife(unit) {
					return total, fmt.Errorf("cannot pay life from selected unit")
				}
				total++
				continue
			}
			if e.payOneBloodSandLoad(playerID, unit, source) || bloodSandPayLife(unit) {
				total++
				continue
			}
			return total, fmt.Errorf("selected unit cannot pay load or life")
		}
	}
	return total, nil
}

type bloodSandPaymentRequest struct {
	instanceID string
	amount     int
	mode       string
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
				requests = append(requests, bloodSandPaymentRequest{instanceID: id, amount: amount, mode: mode})
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
	return min(3, totalLoad(unit)+max(unit.CurrentLife-1, 0))
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
	if unit == nil || unit.CurrentLife <= 1 {
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
