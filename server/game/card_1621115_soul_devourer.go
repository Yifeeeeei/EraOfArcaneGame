package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card1621115SoulDevourer struct{ AlwaysActive }

func (Card1621115SoulDevourer) ID() string { return "1621115" }

func (Card1621115SoulDevourer) Name() string { return "灵魂吸食者" }

func (Card1621115SoulDevourer) PerTurnLabel(*CardInstance) string {
	return "主动"
}

func (Card1621115SoulDevourer) OnPerTurn(ctx *EffectContext) error {
	candidates := soulMarkedFriendlyFieldCandidates(ctx.Engine, ctx.PlayerID)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "soul_devourer_remove_marker",
		"灵魂吸食者:移除你场上的1个灵魂标记物，抽2张并获得2暗", candidates, 1, 1,
		func(selected []string) {
			target := findFriendlyFieldCardIncludingBoundSkill(ctx.Engine, ctx.PlayerID, firstSelected(selected))
			if target == nil || target.Statuses[soulMarkerStatus] <= 0 {
				return
			}
			removeSoulMarkerFromCard(target)
			ctx.Engine.drawCards(ctx.PlayerID, 2)
			ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementShadow: 2})
		})
	return nil
}

func (Card1621115SoulDevourer) ValidateAbility(ctx *EffectContext, trigger EffectTrigger) error {
	if trigger != TriggerPerTurn {
		return nil
	}
	if len(soulMarkedFriendlyFieldCandidates(ctx.Engine, ctx.PlayerID)) == 0 {
		return fmt.Errorf("灵魂吸食者需要你场上有灵魂标记物")
	}
	return nil
}

func soulMarkedFriendlyFieldCandidates(e *Engine, playerID int) []map[string]any {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	ps := e.State.Players[playerID]
	candidates := make([]map[string]any, 0)
	for _, card := range e.getAllFieldCards(ps) {
		if card == nil {
			continue
		}
		if card.Statuses[soulMarkerStatus] > 0 {
			candidates = append(candidates, candidateInfo(card, "field", "own"))
		}
		for _, skill := range card.BoundSkills {
			if skill != nil && skill.Statuses[soulMarkerStatus] > 0 {
				candidates = append(candidates, candidateInfo(skill, "bound_skill", "own"))
			}
		}
	}
	for _, skill := range ps.Skills {
		if skill != nil && skill.Statuses[soulMarkerStatus] > 0 {
			candidates = append(candidates, candidateInfo(skill, "skill", "own"))
		}
	}
	return candidates
}

func findFriendlyFieldCardIncludingBoundSkill(e *Engine, playerID int, instanceID string) *CardInstance {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) || instanceID == "" {
		return nil
	}
	ps := e.State.Players[playerID]
	if card, zone := e.findFriendlyCandidate(playerID, instanceID); card != nil && (zone == "unit" || zone == "equipment" || zone == "skill") {
		return card
	}
	for _, card := range e.getAllFieldCards(ps) {
		if card == nil {
			continue
		}
		for _, skill := range card.BoundSkills {
			if skill != nil && skill.InstanceID == instanceID {
				return skill
			}
		}
	}
	return nil
}

func removeSoulMarkerFromCard(card *CardInstance) {
	if card == nil || card.Statuses[soulMarkerStatus] <= 0 {
		return
	}
	card.Statuses[soulMarkerStatus]--
	if card.Card != nil && isSpellLikeCard(card.Card) {
		card.PowerBonus -= 2
	}
	if card.Statuses[soulMarkerStatus] <= 0 {
		delete(card.Statuses, soulMarkerStatus)
	}
}
