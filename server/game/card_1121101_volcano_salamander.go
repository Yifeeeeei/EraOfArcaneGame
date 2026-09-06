package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card1121101VolcanoSalamander struct{ AlwaysActive }

func (Card1121101VolcanoSalamander) ID() string { return "1121101" }

func (Card1121101VolcanoSalamander) Name() string { return "火山蝾螈" }

func (Card1121101VolcanoSalamander) MasteryMax() int { return 2 }

func (Card1121101VolcanoSalamander) OnMastery(ctx *EffectContext, level int) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || level != 2 || !ctx.Engine.cardStillOnField(ctx.Source) {
		return nil
	}
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, isFireCompanionWithEntryCostLessThanEight)
	if len(candidates) == 0 || len(friendlyPositionsAfterRemovingSource(ctx.Engine, ctx.PlayerID, ctx.Source)) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "volcano_salamander_summon_card",
		"火山蝾螈:选择1个入场费用小于8的火焰伙伴免费召唤", candidates, 1, 1,
		func(selected []string) {
			cardID := firstSelected(selected)
			positions := friendlyPositionsAfterRemovingSource(ctx.Engine, ctx.PlayerID, ctx.Source)
			if cardID == "" || len(positions) == 0 {
				return
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "volcano_salamander_summon_position",
				"火山蝾螈:选择召唤位置", positions, 1, 1,
				func(posSelected []string) {
					pos, ok := positionFromSelectionID(firstSelected(posSelected))
					if !ok || !ctx.Engine.cardStillOnField(ctx.Source) {
						return
					}
					if card, _ := ctx.Engine.State.Players[ctx.PlayerID].FindHandCard(cardID); !isFireCompanionWithEntryCostLessThanEight(card) {
						return
					}
					ctx.Engine.destroyUnitWithCause(ctx.Source, ctx.PlayerID, DeathCauseSacrifice)
					summonCardFreeFromHandOrDeckAtPosition(ctx, cardID, pos)
				})
		})
	return nil
}

func isFireCompanionWithEntryCostLessThanEight(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.IsCompanion() &&
		card.Card.Category == model.ElementFire &&
		totalElementCost(card.Card.ElementsCost) < 8
}

func friendlyPositionsAfterRemovingSource(e *Engine, playerID int, source *CardInstance) []map[string]any {
	positions := e.friendlyEmptyUnitPositions(playerID)
	if source == nil || source.Position == nil {
		return positions
	}
	pos := *source.Position
	if !pos.Valid() {
		return positions
	}
	return append(positions, map[string]any{
		"instance_id": positionSelectionID(pos),
		"name":        fmt.Sprintf("位置 (%d,%d)", pos.Col, pos.Row),
		"zone":        "field_position",
		"side":        "own",
		"col":         pos.Col,
		"row":         pos.Row,
	})
}
