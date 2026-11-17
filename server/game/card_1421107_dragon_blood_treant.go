package game

import (
	"eraofarcane/model"
	"strings"
)

type Card1421107DragonBloodTreant struct{ AlwaysActive }

func (Card1421107DragonBloodTreant) ID() string { return "1421107" }

func (Card1421107DragonBloodTreant) Name() string { return "龙血树精" }

func (Card1421107DragonBloodTreant) OnEnter(ctx *EffectContext) error {
	candidates := make([]map[string]any, 0)
	for _, card := range ctx.Engine.getAllFieldCards(ctx.Engine.State.Players[ctx.PlayerID]) {
		if card == nil || card.Card == nil {
			continue
		}
		load := dragonBloodTreantReducibleLoad(card)
		for _, elem := range model.AllElements {
			if load[elem] <= 0 {
				continue
			}
			candidate := candidateInfo(card, "field", "own")
			candidate["instance_id"] = card.InstanceID + "|" + elem
			candidate["name"] = card.Card.Name + " - " + elem
			candidate["element"] = elem
			candidates = append(candidates, candidate)
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	removeLoad := func(selection string) {
		instanceID, elem, ok := strings.Cut(selection, "|")
		if !ok || elem == "" {
			return
		}
		target, _ := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, instanceID)
		if target == nil || dragonBloodTreantReducibleLoad(target)[elem] <= 0 {
			return
		}
		ctx.Engine.reduceCardElementLoadWithTriggers(ctx.PlayerID, target, elem, 1, ctx.Source)
		ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementShadow, 1, ctx.Source)
	}
	if len(candidates) == 1 {
		id, _ := candidates[0]["instance_id"].(string)
		removeLoad(id)
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "dragon_blood_treant_remove_load",
		"龙血树精:选择1个友方卡牌失去1点负载", candidates, 1, 1,
		func(selected []string) {
			removeLoad(firstSelected(selected))
		})
	return nil
}
