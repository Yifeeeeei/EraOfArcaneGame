package game

const boneKnightRebornStatus = "失去遗言"

type Card1621011BoneKnight struct{ AlwaysActive }

func (Card1621011BoneKnight) ID() string { return "1621011" }

func (Card1621011BoneKnight) Name() string { return "白骨骑士" }

func (Card1621011BoneKnight) HasActiveDeathrattle(card *CardInstance) bool {
	return card != nil && card.Statuses[boneKnightRebornStatus] <= 0
}

func (Card1621011BoneKnight) OnDeath(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	sourceID := ctx.Source.InstanceID
	ctx.Engine.SetPendingAction(ctx.PlayerID, "bone_knight_reborn",
		"白骨骑士:选择此卡重新召唤", []map[string]any{candidateInfo(ctx.Source, "graveyard", "own")}, 1, 1,
		func(selected []string) {
			if firstSelected(selected) != sourceID {
				return
			}
			ctx.Engine.promptBoneKnightRebornPosition(ctx.PlayerID, sourceID)
		})
	return nil
}

func (e *Engine) promptBoneKnightRebornPosition(playerID int, instanceID string) {
	positions := e.friendlyEmptyUnitPositions(playerID)
	if len(positions) == 0 {
		return
	}
	e.SetPendingAction(playerID, "bone_knight_reborn_position",
		"白骨骑士:选择重新召唤位置", positions, 1, 1,
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
			if source == nil || !pos.Valid() || ps.Units[pos.Col][pos.Row] != nil {
				return
			}
			if !e.removeCardFromGraveyard(playerID, source) {
				return
			}
			resetCardForResummon(source)
			source.Statuses[boneKnightRebornStatus] = 1
			if !e.placeExistingCompanionAtPosition(playerID, source, pos, false) {
				e.addToGraveyard(playerID, source)
				return
			}
			e.emit(GameEvent{Type: "summon", Player: playerID, Data: map[string]any{
				"player":   playerID,
				"card":     cardToInfo(source),
				"position": pos,
			}})
		})
}
