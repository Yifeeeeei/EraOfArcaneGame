package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card2621112SoulStaff struct{ AlwaysActive }

func (Card2621112SoulStaff) ID() string { return "2621112" }

func (Card2621112SoulStaff) Name() string { return "灵魂法杖" }

func (Card2621112SoulStaff) PerTurnLabel(*CardInstance) string {
	return "主动"
}

func (Card2621112SoulStaff) OnPerTurn(ctx *EffectContext) error {
	graveyardCandidates := shadowCompanionGraveyardCandidates(ctx.Engine.State.Players[ctx.PlayerID])
	if len(graveyardCandidates) < 2 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "soul_staff_exile_companions",
		"灵魂法杖:选择2张暗影伙伴移出游戏", graveyardCandidates, 2, 2,
		func(selected []string) {
			if moveSelectedShadowCompanionsFromGraveyardToExile(ctx.Engine, ctx.PlayerID, selected, 2) < 2 {
				return
			}
			spellCandidates := ctx.Engine.friendlySkillsIncludingBound(ctx.PlayerID, isShadowSpellInstance)
			if len(spellCandidates) == 0 {
				return
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "soul_staff_mark_spell",
				"灵魂法杖:选择1个暗影法术放置1个灵魂标记物", spellCandidates, 1, 1,
				func(spellSelected []string) {
					skill := findFriendlySkillIncludingBound(ctx.Engine, ctx.PlayerID, firstSelected(spellSelected))
					if !isShadowSpellInstance(skill) {
						return
					}
					addSoulMarkerToSpell(skill)
				})
		})
	return nil
}

func (Card2621112SoulStaff) ValidateAbility(ctx *EffectContext, trigger EffectTrigger) error {
	if trigger != TriggerPerTurn {
		return nil
	}
	if len(shadowCompanionGraveyardCandidates(ctx.Engine.State.Players[ctx.PlayerID])) < 2 {
		return fmt.Errorf("灵魂法杖需要弃牌堆中至少2张暗影伙伴")
	}
	if len(ctx.Engine.friendlySkillsIncludingBound(ctx.PlayerID, isShadowSpellInstance)) == 0 {
		return fmt.Errorf("灵魂法杖需要1个暗影法术")
	}
	return nil
}

func shadowCompanionGraveyardCandidates(ps *PlayerState) []map[string]any {
	candidates := make([]map[string]any, 0)
	if ps == nil {
		return candidates
	}
	for _, card := range ps.Graveyard {
		if card != nil && card.Card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementShadow {
			candidates = append(candidates, candidateInfo(card, "graveyard", "own"))
		}
	}
	return candidates
}

func moveSelectedShadowCompanionsFromGraveyardToExile(e *Engine, playerID int, selected []string, maxCount int) int {
	if e == nil || maxCount <= 0 || playerID < 0 || playerID >= len(e.State.Players) {
		return 0
	}
	ps := e.State.Players[playerID]
	selectedSet := make(map[string]bool, len(selected))
	for _, id := range selected {
		selectedSet[id] = true
	}
	moved := 0
	for _, card := range append([]*CardInstance(nil), ps.Graveyard...) {
		if moved >= maxCount {
			break
		}
		if card == nil || !selectedSet[card.InstanceID] || card.Card == nil || !card.Card.IsCompanion() || card.Card.Category != model.ElementShadow {
			continue
		}
		if e.exileCard(playerID, card) {
			moved++
		}
	}
	return moved
}

func isShadowSpellInstance(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.Category == model.ElementShadow && isSpellLikeCard(card.Card)
}
