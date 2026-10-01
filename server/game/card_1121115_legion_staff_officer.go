package game

import (
	"eraofarcane/model"
)

type Card1121115LegionStaffOfficer struct{ AlwaysActive }

func (Card1121115LegionStaffOfficer) ID() string { return "1121115" }

func (Card1121115LegionStaffOfficer) Name() string { return "军团参谋" }

func (Card1121115LegionStaffOfficer) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.Card == nil {
		return nil
	}
	if !triggeredTurnAvailable(ctx.Source) || !isFriendlySpellCast(ctx) || !hasCardTag(ctx.Target.Card, "创造") {
		return nil
	}
	if !deckHasMatch(ctx.Engine.State.Players[ctx.PlayerID], isFireConsumableItem) {
		return nil
	}
	ctx.Engine.SetTriggeredTurnAction(ctx.Source, ctx.PlayerID, "legion_staff_officer_flip_fire_consumable",
		"军团参谋:是否翻取1个火焰消耗品道具", []map[string]any{candidateInfo(ctx.Source, "unit", "own")}, 0, 1,
		func(selected []string) {
			if len(selected) == 0 || !deckHasMatch(ctx.Engine.State.Players[ctx.PlayerID], isFireConsumableItem) || !useTriggeredTurn(ctx.Source) {
				return
			}
			drawn := ctx.Engine.flipDeckMatchesToHand(ctx.PlayerID, 1, 0, isFireConsumableItem)
			if len(drawn) == 0 {
				return
			}
			ps := ctx.Engine.State.Players[ctx.PlayerID]
			if ps.DiscardAtTurnEnd == nil {
				ps.DiscardAtTurnEnd = make(map[string]bool)
			}
			ps.DiscardAtTurnEnd[drawn[0].InstanceID] = true
		})
	return nil
}

func isFireConsumableItem(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.Category == model.ElementFire && isConsumableCardInstance(card)
}
