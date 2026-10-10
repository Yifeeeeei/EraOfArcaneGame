package game

import (
	"eraofarcane/model"
	"fmt"
)

func dragonBloodTreantReducibleLoad(card *CardInstance) map[string]int {
	load := make(map[string]int)
	if card == nil || card.Card == nil {
		return load
	}
	base := card.Card.ElementsGain
	if card.ElementsGainSet != nil {
		base = card.ElementsGainSet
	}
	for elem, amount := range base {
		if amount > 0 {
			load[elem] += amount
		}
	}
	for elem, amount := range card.ElementsGainBonus {
		if amount > 0 {
			load[elem] += amount
		}
	}
	return load
}

func dragonBloodTreantRemoveLoad(card *CardInstance, elem string) {
	if card == nil || elem == "" {
		return
	}
	if card.ElementsGainBonus != nil && card.ElementsGainBonus[elem] > 0 {
		card.ElementsGainBonus[elem]--
		if card.ElementsGainBonus[elem] == 0 {
			delete(card.ElementsGainBonus, elem)
		}
		return
	}
	base := copyElementCost(card.Card.ElementsGain)
	if card.ElementsGainSet != nil {
		base = copyElementCost(card.ElementsGainSet)
	}
	if base[elem] <= 0 {
		return
	}
	base[elem]--
	setElementsGain(card, base)
}

func reducibleElementLoad(card *CardInstance, elem string) int {
	return dragonBloodTreantReducibleLoad(card)[elem]
}

func reduceCardElementLoad(card *CardInstance, elem string, amount int) int {
	removed := 0
	for i := 0; i < amount; i++ {
		if reducibleElementLoad(card, elem) <= 0 {
			return removed
		}
		dragonBloodTreantRemoveLoad(card, elem)
		removed++
	}
	return removed
}

func (e *Engine) reduceCardElementLoadWithTriggers(playerID int, card *CardInstance, elem string, amount int, cause *CardInstance) int {
	removed := reduceCardElementLoad(card, elem, amount)
	if e == nil || removed <= 0 || card == nil || card.Card == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return removed
	}
	data := map[string]any{
		"load_loss_player": playerID,
		"load_loss_target": card,
		"load_loss_source": cause,
		"element":          elem,
		"amount":           removed,
	}
	e.emit(GameEvent{
		Type:   "load_loss",
		Player: -1,
		Data: map[string]any{
			"player":  playerID,
			"target":  cardToInfo(card),
			"source":  cardToInfo(cause),
			"element": elem,
			"amount":  removed,
		},
	})
	e.triggerEffects(TriggerOnLoadLoss, card, cause, data)
	e.triggerFieldEffectsWithData(TriggerOnLoadLoss, playerID, card, data)
	e.triggerFieldEffectsWithData(TriggerOnLoadLoss, 1-playerID, card, data)
	return removed
}

const entryCostZeroStatus = "入场费用变为0"

func makeEntryCostZero(card *CardInstance) {
	if card == nil || card.Card == nil {
		return
	}
	if card.Statuses == nil {
		card.Statuses = make(map[string]int)
	}
	card.Statuses[entryCostZeroStatus] = 1
	for elem, amount := range card.Card.ElementsCost {
		if amount > 0 {
			card.Statuses["入场费用"+elem+"-"+fmt.Sprint(amount)]++
		}
	}
}

type royalWaterUseCostReduction struct {
	AlwaysActive
	id             string
	name           string
	requireWater   bool
	triggerOnEnter bool
}

func (r royalWaterUseCostReduction) ID() string { return r.id }

func (r royalWaterUseCostReduction) Name() string { return r.name }

func (r royalWaterUseCostReduction) HasActiveUseItem(*CardInstance) bool {
	return !r.triggerOnEnter
}

func (r royalWaterUseCostReduction) OnUseItem(ctx *EffectContext) error {
	return r.prompt(ctx)
}

func (r royalWaterUseCostReduction) OnEnter(ctx *EffectContext) error {
	if !r.triggerOnEnter {
		return nil
	}
	return r.prompt(ctx)
}

func (r royalWaterUseCostReduction) prompt(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlySkillsIncludingBound(ctx.PlayerID, func(skill *CardInstance) bool {
		return skill != nil && skill.Card != nil && isSpellLikeCard(skill.Card) && (!r.requireWater || skill.Card.Category == model.ElementWater)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "royal_water_use_cost_reduction",
		r.name+":选择你的1个法术使用花费-1水", candidates, 1, 1,
		func(selected []string) {
			skill := ctx.Engine.findSkill(ctx.Engine.State.Players[ctx.PlayerID], firstSelected(selected))
			if skill == nil || skill.Card == nil || !isSpellLikeCard(skill.Card) || (r.requireWater && skill.Card.Category != model.ElementWater) {
				return
			}
			skill.Statuses["使用费用"+model.ElementWater+"-1"]++
		})
	return nil
}

func elementChoiceCandidates(sourceNumber string, elements ...string) []map[string]any {
	candidates := make([]map[string]any, 0, len(elements))
	for _, elem := range elements {
		candidates = append(candidates, map[string]any{
			"instance_id": elem,
			"number":      sourceNumber,
			"name":        elem,
			"type":        "元素",
			"zone":        "choice",
			"side":        "own",
		})
	}
	return candidates
}

func isNonArcaneElement(elem string) bool {
	return elem == model.ElementFire || elem == model.ElementWater || elem == model.ElementEarth || elem == model.ElementAir || elem == model.ElementLight || elem == model.ElementShadow
}
