package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card2421105NaturalCommunion struct{ AlwaysActive }

func (Card2421105NaturalCommunion) ID() string { return "2421105" }

func (Card2421105NaturalCommunion) Name() string { return "自然交感" }

func (Card2421105NaturalCommunion) OnUseItem(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil {
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementEarth
	})
	if len(candidates) < 2 {
		return nil
	}
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "natural_communion_distribute",
		"自然交感:选择2个地脉伙伴并提交新的负载分配", candidates, 2, 2, nil, false,
		func(selected []string, data map[string]any) error {
			a, _ := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, selected[0])
			b, _ := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, selected[1])
			if !isEarthCompanionUnit(a) || !isEarthCompanionUnit(b) || a == b {
				return fmt.Errorf("自然交感需要2个不同的地脉伙伴")
			}
			distribution, err := parseLoadDistribution(data["load_distribution"])
			if err != nil {
				return err
			}
			aLoad, okA := distribution[a.InstanceID]
			bLoad, okB := distribution[b.InstanceID]
			if !okA || !okB || len(distribution) != 2 {
				return fmt.Errorf("自然交感需要为选择的2个伙伴提交负载分配")
			}
			totalBefore := mergeElementCosts(ctx.Engine.effectiveElementsGain(a), ctx.Engine.effectiveElementsGain(b))
			totalAfter := mergeElementCosts(aLoad, bLoad)
			if !sameElementCost(totalBefore, totalAfter) {
				return fmt.Errorf("自然交感的新负载分配必须保持总负载不变")
			}
			applyRedistributedLoad(a, aLoad)
			applyRedistributedLoad(b, bLoad)
			ctx.Engine.emit(GameEvent{
				Type:   "natural_communion_distribute",
				Player: -1,
				Data: map[string]any{
					"player": ctx.PlayerID,
					"cards":  []any{cardToInfo(a), cardToInfo(b)},
				},
			})
			return nil
		})
	return nil
}

var _ OnUseItemBehavior = Card2421105NaturalCommunion{}

func (Card2421105NaturalCommunion) ValidateItemUse(ctx *EffectContext) error {
	e, playerID := ctx.Engine, ctx.PlayerID
	if len(e.friendlyUnits(playerID, false, isEarthCompanionUnit)) < 2 {
		return fmt.Errorf("Natural Communion requires two friendly earth companions")
	}
	return nil
}

func isEarthCompanionUnit(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementEarth
}

func parseLoadDistribution(raw any) (map[string]map[string]int, error) {
	input, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("missing load_distribution")
	}
	result := make(map[string]map[string]int, len(input))
	for instanceID, value := range input {
		elemMap, ok := value.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid load distribution")
		}
		load := make(map[string]int)
		for elem, amountRaw := range elemMap {
			if !isNonArcaneElement(elem) && elem != model.ElementArcane {
				return nil, fmt.Errorf("invalid load element")
			}
			amount, ok := intFromAny(amountRaw)
			if !ok || amount < 0 {
				return nil, fmt.Errorf("invalid load amount")
			}
			if amount > 0 {
				load[elem] = amount
			}
		}
		result[instanceID] = load
	}
	return result, nil
}

func sameElementCost(a, b map[string]int) bool {
	for _, elem := range model.AllElements {
		if a[elem] != b[elem] {
			return false
		}
	}
	return true
}

func applyRedistributedLoad(card *CardInstance, load map[string]int) {
	if card == nil {
		return
	}
	card.ElementsGainBonus = make(map[string]int)
	setElementsGain(card, load)
}
