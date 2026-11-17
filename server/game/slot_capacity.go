package game

func (e *Engine) enforceAllSlotCapacities() {
	for _, ps := range e.State.Players {
		e.enforceSlotCapacities(ps)
	}
}

func (e *Engine) enforceSlotCapacities(ps *PlayerState) {
	if ps == nil {
		return
	}

	equipmentCap := equipmentSlotCapacity(ps)
	for i := equipmentCap; i < len(ps.Equipment); i++ {
		equipment := ps.Equipment[i]
		if equipment == nil {
			continue
		}
		e.moveEquipmentToGraveyard(ps.PlayerID, i, equipment)
	}

	for i := 0; i < len(ps.Skills); i++ {
		skill := ps.Skills[i]
		if skill == nil {
			continue
		}
		if skillAllowedInSlot(ps, skill, i) {
			continue
		}
		ps.Skills[i] = nil
		skill = e.leaveBattlefield(skill)
		ps.SkillPool = append(ps.SkillPool, skill)
		e.emit(GameEvent{Type: "skill_returned_to_pool", Player: ps.PlayerID, Data: map[string]any{"card": cardToInfo(skill)}})
	}
}
