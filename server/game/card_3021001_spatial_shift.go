package game

type Card3021001SpatialShift struct{ AlwaysActive }

func (Card3021001SpatialShift) ID() string   { return "3021001" }
func (Card3021001SpatialShift) Name() string { return "移形换影" }

func (Card3021001SpatialShift) OnSpellCast(ctx *EffectContext) error {
	if !isSpellBeingCast(ctx) {
		return nil
	}
	units := ctx.Engine.friendlyUnits(ctx.PlayerID, true, nil)
	if len(units) == 0 || len(ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "spatial_shift_unit",
		"选择1个友方单位移动", units, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			unitID := selected[0]
			positions := ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)
			if len(positions) == 0 {
				return
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "spatial_shift_position",
				"选择移动后的空位", positions, 1, 1,
				func(selected []string) {
					if len(selected) == 0 {
						return
					}
					pos, ok := positionFromSelectionID(selected[0])
					if !ok {
						return
					}
					ctx.Engine.moveFriendlyUnitToPosition(ctx.PlayerID, unitID, pos)
				})
		})
	return nil
}

func (e *Engine) moveFriendlyUnitToPosition(playerID int, instanceID string, pos Position) bool {
	if !pos.Valid() {
		return false
	}
	ps := e.State.Players[playerID]
	if ps.Units[pos.Col][pos.Row] != nil {
		return false
	}
	unit := e.findFieldCardByInstance(ps, instanceID)
	if unit == nil || unit.Position == nil {
		return false
	}
	oldPos := *unit.Position
	ps.Units[oldPos.Col][oldPos.Row] = nil
	unit.Position = &Position{Col: pos.Col, Row: pos.Row}
	ps.Units[pos.Col][pos.Row] = unit
	e.emit(GameEvent{Type: "unit_moved", Player: -1, Data: map[string]any{
		"player": playerID,
		"card":   cardToInfo(unit),
		"from":   oldPos,
		"to":     pos,
	}})
	return true
}
