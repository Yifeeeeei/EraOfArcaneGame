package game

import (
	"eraofarcane/model"
)

type Card1611102BloodThornGarden struct{ AlwaysActive }

func (Card1611102BloodThornGarden) ID() string { return "1611102" }

func (Card1611102BloodThornGarden) Name() string { return "蔷薇花园的血荆棘" }

func (Card1611102BloodThornGarden) OnDeath(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || !bloodThornKilledByFriendlyCard(ctx) {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ps == nil || ps.Elements[model.ElementShadow] < 1 {
		return nil
	}
	sourceID := ctx.Source.InstanceID
	ctx.Engine.SetPendingAction(ctx.PlayerID, "blood_thorn_resummon",
		"蔷薇花园的血荆棘:是否支付1暗重新召唤", []map[string]any{candidateInfo(ctx.Source, "graveyard", "own")}, 0, 1,
		func(selected []string) {
			if len(selected) == 0 || firstSelected(selected) != sourceID || ps.Elements[model.ElementShadow] < 1 {
				return
			}
			ctx.Engine.promptBloodThornResummonPosition(ctx.PlayerID, sourceID)
		})
	return nil
}

func (e *Engine) promptBloodThornResummonPosition(playerID int, instanceID string) {
	positions := e.friendlyEmptyUnitPositions(playerID)
	if len(positions) == 0 {
		return
	}
	e.SetPendingAction(playerID, "blood_thorn_resummon_position",
		"蔷薇花园的血荆棘:选择重新召唤位置", positions, 1, 1,
		func(selected []string) {
			pos, ok := positionFromSelectionID(firstSelected(selected))
			if !ok {
				return
			}
			ps := e.State.Players[playerID]
			var source *CardInstance
			for _, card := range ps.Graveyard {
				if card != nil && card.InstanceID == instanceID {
					source = card
					break
				}
			}
			if source == nil || !pos.Valid() || ps.Units[pos.Col][pos.Row] != nil || ps.Elements[model.ElementShadow] < 1 {
				return
			}
			if !e.removeCardFromGraveyard(playerID, source) {
				return
			}
			ps.Elements[model.ElementShadow]--
			resetCardForResummon(source)
			if !e.placeExistingCompanionAtPosition(playerID, source, pos, true) {
				e.addToGraveyard(playerID, source)
				return
			}
			e.emit(GameEvent{
				Type:   "blood_thorn_resummon",
				Player: -1,
				Data: map[string]any{
					"player": playerID,
					"card":   cardToInfo(source),
				},
			})
		})
}
