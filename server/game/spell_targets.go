package game

func (e *Engine) spellAffectedUnits(defenderID int, skill *CardInstance, target SpellTarget) []*CardInstance {
	if target.Type != "unit" {
		return nil
	}
	defender := e.State.Players[defenderID]
	units := make([]*CardInstance, 0, 9)

	if affected, ok := e.customSpellAffectedUnits(defenderID, skill, target); ok {
		return affected
	}

	switch e.effectiveSpellArea(skill) {
	case SpellAreaSquare:
		startCol := min(max(target.Position.Col, 0), 1)
		startRow := min(max(target.Position.Row, 0), 1)
		for col := startCol; col < startCol+2; col++ {
			for row := startRow; row < startRow+2; row++ {
				if defender.Units[col][row] != nil {
					units = append(units, defender.Units[col][row])
				}
			}
		}
		return units
	case SpellAreaAll:
		for col := 0; col < 3; col++ {
			for row := 0; row < 3; row++ {
				if defender.Units[col][row] != nil {
					units = append(units, defender.Units[col][row])
				}
			}
		}
		return units
	case SpellAreaColumn:
		for row := 0; row < 3; row++ {
			if defender.Units[target.Position.Col][row] != nil {
				units = append(units, defender.Units[target.Position.Col][row])
			}
		}
		return units
	case SpellAreaFrontRow:
		frontRow := defender.GetFrontRow()
		if frontRow >= 0 {
			for col := 0; col < 3; col++ {
				if defender.Units[col][frontRow] != nil {
					units = append(units, defender.Units[col][frontRow])
				}
			}
		}
		return units
	case SpellAreaSplashCross:
		for _, delta := range []struct{ col, row int }{{0, 0}, {-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			col := target.Position.Col + delta.col
			row := target.Position.Row + delta.row
			if col < 0 || col >= 3 || row < 0 || row >= 3 {
				continue
			}
			if defender.Units[col][row] != nil {
				units = append(units, defender.Units[col][row])
			}
		}
		return units
	}

	if target.Position.Valid() && defender.Units[target.Position.Col][target.Position.Row] != nil {
		return []*CardInstance{defender.Units[target.Position.Col][target.Position.Row]}
	}
	return nil
}

func (e *Engine) effectiveSpellArea(skill *CardInstance) SpellArea {
	area := spellArea(skill)
	if skill == nil || skill.Card == nil || skill.OwnerID < 0 || skill.OwnerID >= len(e.State.Players) {
		return area
	}
	ps := e.State.Players[skill.OwnerID]
	ctx := &EffectContext{
		Engine:     e,
		Source:     skill,
		PlayerID:   skill.OwnerID,
		OpponentID: 1 - skill.OwnerID,
	}
	for _, fieldCard := range e.getAllFieldCards(ps) {
		if fieldCard == nil || fieldCard.Card == nil || e.hasEffectiveStatus(fieldCard, StatusPetrify) {
			continue
		}
		behavior := globalRegistry.GetBehavior(fieldCard.Card.Number)
		if modifier, ok := behavior.(SpellAreaModifier); ok && modifier.HasActiveSpellAreaModifier(fieldCard) {
			ctx.Target = fieldCard
			modifier.ModifySpellArea(ctx, &area)
		}
	}
	return area
}

func (e *Engine) applyGenericSpellEffects(attackerID int, defenderID int, skill *CardInstance, targets []*CardInstance, target SpellTarget) {
	e.applyGenericElementGain(attackerID, skill)
	for _, unit := range targets {
		e.applyExplicitSpellHitStatuses(skill, unit)
	}
}
