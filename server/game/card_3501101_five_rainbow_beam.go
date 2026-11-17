package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card3501101FiveRainbowBeam struct{ AlwaysActive }

func (Card3501101FiveRainbowBeam) ID() string { return "3501101" }

func (Card3501101FiveRainbowBeam) Name() string { return "五虹之束" }

func (Card3501101FiveRainbowBeam) HasActivePierce(card *CardInstance) bool {
	return card != nil && card.Statuses[fiveRainbowBeamSelectedStatus(model.ElementAir)] > 0
}

func (Card3501101FiveRainbowBeam) HasPierce() bool { return true }

func (Card3501101FiveRainbowBeam) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Source == nil || stats == nil {
		return
	}
	fire := ctx.Source.Statuses[fiveRainbowBeamSelectedStatus(model.ElementFire)]
	earth := ctx.Source.Statuses[fiveRainbowBeamSelectedStatus(model.ElementEarth)]
	stats.DamageBonus += fire * 2
	stats.PowerBonus += earth * 3
	if ctx.Source.Statuses[fiveRainbowBeamAllStatus] > 0 {
		stats.PowerBonus *= 2
	}
}

func (Card3501101FiveRainbowBeam) PrepareSpellCast(ctx *EffectContext, target SpellTarget, action ActionMessage) (SpellCastOptions, error) {
	markers, ring, err := ctx.Engine.prepareFiveRainbowBeamMarkers(ctx.PlayerID, ctx.Source, action)
	if err != nil {
		return SpellCastOptions{}, err
	}
	return SpellCastOptions{
		Pierce: markers[model.ElementAir] > 0,
		ExtraTargets: func(pierce bool) ([]SpellTarget, error) {
			return ctx.Engine.fiveRainbowBeamExtraTargetsFromAction(ctx.PlayerID, ctx.Source, target, action, markers[model.ElementWater], pierce)
		},
		ModifyCost: func(cost map[string]int) {
			if markers[model.ElementLight] > 0 {
				reduceGenericCost(cost, model.ElementLight, markers[model.ElementLight]*2)
			}
		},
		Commit: func() { ctx.Engine.applyFiveRainbowBeamMarkers(ctx.Source, ring, markers) },
	}, nil
}

func fiveRainbowBeamSelectedStatus(elem string) string {
	return "五虹之束消耗:" + elem
}

const fiveRainbowBeamAllStatus = "五虹之束五色齐发"

func (e *Engine) prepareFiveRainbowBeamMarkers(playerID int, skill *CardInstance, action ActionMessage) (map[string]int, *CardInstance, error) {
	markers := map[string]int{}
	if skill == nil || skill.Card == nil || skill.Card.Number != "3501101" {
		return markers, nil, nil
	}
	for _, elem := range fiveRainbowElements() {
		markers[elem] = 0
	}
	raw, ok := action.Data["rainbow_markers"]
	if !ok {
		return markers, nil, nil
	}
	values, ok := raw.(map[string]any)
	if !ok {
		return nil, nil, fmt.Errorf("invalid rainbow markers")
	}
	ring := e.findFiveRainbowRingForSkill(playerID, skill)
	if ring == nil {
		return nil, nil, fmt.Errorf("five rainbow beam requires five rainbow ring")
	}
	for _, elem := range fiveRainbowElements() {
		amount, ok := intFromAny(values[elem])
		if !ok {
			amount = 0
		}
		if amount < 0 {
			return nil, nil, fmt.Errorf("invalid rainbow marker amount")
		}
		if amount > ring.Statuses[fiveRainbowMarkerStatus(elem)] {
			return nil, nil, fmt.Errorf("not enough rainbow markers")
		}
		markers[elem] = amount
	}
	return markers, ring, nil
}

func (e *Engine) applyFiveRainbowBeamMarkers(skill *CardInstance, ring *CardInstance, markers map[string]int) {
	if skill == nil || skill.Card == nil || skill.Card.Number != "3501101" || len(markers) == 0 {
		return
	}
	clearFiveRainbowBeamSelection(skill)
	fullSet := true
	for _, elem := range fiveRainbowElements() {
		amount := markers[elem]
		if amount <= 0 {
			fullSet = false
			continue
		}
		skill.Statuses[fiveRainbowBeamSelectedStatus(elem)] = amount
		if ring != nil {
			ring.Statuses[fiveRainbowMarkerStatus(elem)] -= amount
			if ring.Statuses[fiveRainbowMarkerStatus(elem)] <= 0 {
				delete(ring.Statuses, fiveRainbowMarkerStatus(elem))
			}
		}
	}
	if fullSet {
		skill.Statuses[fiveRainbowBeamAllStatus] = 1
	}
}

func clearFiveRainbowBeamSelection(skill *CardInstance) {
	if skill == nil {
		return
	}
	for _, elem := range fiveRainbowElements() {
		delete(skill.Statuses, fiveRainbowBeamSelectedStatus(elem))
	}
	delete(skill.Statuses, fiveRainbowBeamAllStatus)
}

func (e *Engine) findFiveRainbowRingForSkill(playerID int, skill *CardInstance) *CardInstance {
	if e == nil || skill == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	ps := e.State.Players[playerID]
	for _, card := range e.getAllFieldCards(ps) {
		if card == nil || card.Card == nil || card.Card.Number != "2511102" {
			continue
		}
		for _, bound := range card.BoundSkills {
			if bound == skill {
				return card
			}
		}
	}
	return nil
}

func fiveRainbowElements() []string {
	return []string{model.ElementFire, model.ElementWater, model.ElementEarth, model.ElementAir, model.ElementLight}
}

func (e *Engine) fiveRainbowBeamExtraTargetsFromAction(playerID int, skill *CardInstance, mainTarget SpellTarget, action ActionMessage, maxTargets int, hasPierce bool) ([]SpellTarget, error) {
	if maxTargets <= 0 {
		return nil, nil
	}
	raw, ok := action.Data["extra_targets"].([]any)
	if !ok || len(raw) == 0 {
		return nil, nil
	}
	if len(raw) > maxTargets {
		return nil, fmt.Errorf("too many five rainbow beam extra targets")
	}
	result := make([]SpellTarget, 0, len(raw))
	for _, value := range raw {
		data, ok := value.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid five rainbow beam extra target")
		}
		pos, err := requiredBoardPosition(data, "col", "row")
		if err != nil {
			return nil, fmt.Errorf("invalid five rainbow beam extra target")
		}
		extra := SpellTarget{Type: "unit", Position: pos}
		if ownerF, ok := data["owner"].(float64); ok {
			owner := int(ownerF)
			extra.OwnerID = &owner
		}
		if extra.Position == mainTarget.Position && !e.allowsSameSpellExtraTarget(e.State.Players[playerID], skill) {
			continue
		}
		if err := e.validateSpellTargetWithPierce(playerID, skill, extra, hasPierce); err != nil {
			return nil, err
		}
		result = append(result, extra)
	}
	return result, nil
}
