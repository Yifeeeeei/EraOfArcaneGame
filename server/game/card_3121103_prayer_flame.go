package game

import (
	"eraofarcane/model"
)

type Card3121103PrayerFlame struct{ AlwaysActive }

func (Card3121103PrayerFlame) ID() string { return "3121103" }

func (Card3121103PrayerFlame) Name() string { return "祈祷之焰" }

func (Card3121103PrayerFlame) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || !isFriendlySpellCast(ctx) || !isSpellBeingCast(ctx) {
		return nil
	}
	if ctx.Source.Card.Number != "3121103" {
		return nil
	}
	choices := []map[string]any{
		{"instance_id": "add_markers", "name": "放置3个标记物", "zone": "choice", "side": "own"},
	}
	if prayerFlameHasSummonTarget(ctx.Engine, ctx.PlayerID, ctx.Source.Statuses[prayerFlameMarkerStatus]) {
		choices = append(choices, map[string]any{"instance_id": "summon", "name": "取除标记物免费召唤火焰伙伴", "zone": "choice", "side": "own"})
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "prayer_flame_choice",
		"祈祷之焰:选择放置标记物或取除标记物召唤火焰伙伴", choices, 1, 1,
		func(selected []string) {
			if firstSelected(selected) != "summon" {
				ctx.Source.Statuses[prayerFlameMarkerStatus] += 3
				return
			}
			markers := ctx.Source.Statuses[prayerFlameMarkerStatus]
			if markers <= 0 {
				return
			}
			openPrayerFlameSummonPrompt(ctx, markers)
		})
	return nil
}

const prayerFlameMarkerStatus = "祈祷之焰标记物"

func prayerFlameHasSummonTarget(e *Engine, playerID int, markers int) bool {
	return markers > 0 && len(e.friendlyHandCards(playerID, func(card *CardInstance) bool {
		return isFireCompanionWithEntryCostAtMost(card, markers)
	})) > 0 && len(e.friendlyEmptyUnitPositions(playerID)) > 0
}

func openPrayerFlameSummonPrompt(ctx *EffectContext, markers int) {
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, func(card *CardInstance) bool {
		return isFireCompanionWithEntryCostAtMost(card, markers)
	})
	if len(candidates) == 0 {
		return
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "prayer_flame_summon_card",
		"祈祷之焰:选择1个火焰伙伴免费召唤", candidates, 1, 1,
		func(selected []string) {
			cardID := firstSelected(selected)
			positions := ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)
			if len(positions) == 0 {
				return
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "prayer_flame_summon_position",
				"祈祷之焰:选择召唤位置", positions, 1, 1,
				func(posSelected []string) {
					pos, ok := positionFromSelectionID(firstSelected(posSelected))
					if !ok {
						return
					}
					card := summonCardFreeFromHandOrDeckAtPosition(ctx, cardID, pos)
					if card != nil {
						delete(ctx.Source.Statuses, prayerFlameMarkerStatus)
					}
				})
		})
}

func isFireCompanionWithEntryCostAtMost(card *CardInstance, maxCost int) bool {
	return card != nil && card.Card != nil && card.Card.IsCompanion() &&
		card.Card.Category == model.ElementFire &&
		totalElementCost(card.Card.ElementsCost) <= maxCost
}
