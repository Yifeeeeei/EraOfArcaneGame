package game

type Card1421106PhantomLizard struct{ AlwaysActive }

func (Card1421106PhantomLizard) ID() string { return "1421106" }

func (Card1421106PhantomLizard) Name() string { return "幻影蜥蜴" }

func (Card1421106PhantomLizard) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.Card == nil || ctx.Source.UltimateUsed {
		return nil
	}
	if !isFriendlySpellCast(ctx) || !hasCardTag(ctx.Target.Card, "灵媒") || !ctx.Engine.canConsumeCard(ctx.Source) {
		return nil
	}
	if len(ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)) < 1 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "phantom_lizard_split",
		"幻影蜥蜴:消耗此卡并分裂为两个普通蜥蜴", []map[string]any{candidateInfo(ctx.Source, "unit", "own")}, 1, 1,
		func(selected []string) {
			if !ctx.Engine.canConsumeCard(ctx.Source) || !ctx.Engine.cardStillOnField(ctx.Source) {
				return
			}
			ctx.Source.UltimateUsed = true
			ctx.Engine.consumeCardForEffectWithTriggers(ctx.PlayerID, ctx.Source, ctx.Engine.effectiveElementsGain(ctx.Source), "")
			moveUnitToGraveyardWithoutDeath(ctx.Engine, ctx.PlayerID, ctx.Source)
			positions := ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)
			if len(positions) < 2 {
				return
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "phantom_lizard_first_position",
				"幻影蜥蜴:选择第1个普通蜥蜴的位置", positions, 1, 1,
				func(firstSelectedPos []string) {
					firstPos, ok := positionFromSelectionID(firstSelected(firstSelectedPos))
					if !ok || ctx.Engine.summonFreshCardAtPosition(ctx.PlayerID, "1401101", firstPos, true) == nil {
						return
					}
					secondPositions := ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)
					if len(secondPositions) == 0 {
						return
					}
					ctx.Engine.SetPendingAction(ctx.PlayerID, "phantom_lizard_second_position",
						"幻影蜥蜴:选择第2个普通蜥蜴的位置", secondPositions, 1, 1,
						func(secondSelectedPos []string) {
							secondPos, ok := positionFromSelectionID(firstSelected(secondSelectedPos))
							if ok {
								ctx.Engine.summonFreshCardAtPosition(ctx.PlayerID, "1401101", secondPos, true)
							}
						})
				})
		})
	return nil
}

func moveUnitToGraveyardWithoutDeath(e *Engine, playerID int, unit *CardInstance) {
	if e == nil || unit == nil || unit.Position == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	ps := e.State.Players[playerID]
	if ps.Units[unit.Position.Col][unit.Position.Row] != unit {
		return
	}
	ps.Units[unit.Position.Col][unit.Position.Row] = nil
	unit.Position = nil
	e.releaseUnderCardsToGraveyard(playerID, unit)
	e.exileTransferredBoundSkills(playerID, unit)
	unit.BoundSkills = nil
	e.addToGraveyard(playerID, unit)
	e.emit(GameEvent{Type: "unit_transformed", Player: -1, Data: map[string]any{
		"player": playerID,
		"card":   cardToInfo(unit),
	}})
}
