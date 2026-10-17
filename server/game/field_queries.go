package game

import ()

func (e *Engine) findCardOnField(ps *PlayerState, instanceID string) *CardInstance {
	if ps.Hero != nil && ps.Hero.InstanceID == instanceID {
		return ps.Hero
	}
	// Check units
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if ps.Units[col][row] != nil && ps.Units[col][row].InstanceID == instanceID {
				return ps.Units[col][row]
			}
		}
	}
	// Check terrain
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if ps.Terrain[col][row] != nil && ps.Terrain[col][row].InstanceID == instanceID {
				return ps.Terrain[col][row]
			}
		}
	}
	// Check equipment
	for i := range ps.Equipment {
		if ps.Equipment[i] != nil && ps.Equipment[i].InstanceID == instanceID {
			return ps.Equipment[i]
		}
	}
	return nil
}

func (e *Engine) findUnitOnGrid(ps *PlayerState, instanceID string) *CardInstance {
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if ps.Units[col][row] != nil && ps.Units[col][row].InstanceID == instanceID {
				return ps.Units[col][row]
			}
		}
	}
	return nil
}

func (e *Engine) findEquipment(ps *PlayerState, instanceID string) *CardInstance {
	for i := range ps.Equipment {
		if ps.Equipment[i] != nil && ps.Equipment[i].InstanceID == instanceID {
			return ps.Equipment[i]
		}
	}
	return nil
}

func (e *Engine) findSkill(ps *PlayerState, instanceID string) *CardInstance {
	for i := range ps.Skills {
		if ps.Skills[i] != nil && ps.Skills[i].InstanceID == instanceID {
			return ps.Skills[i]
		}
	}
	for _, card := range e.getAllFieldCards(ps) {
		if card == nil {
			continue
		}
		for _, skill := range card.BoundSkills {
			if skill != nil && skill.InstanceID == instanceID {
				return skill
			}
		}
	}
	return nil
}

func (e *Engine) findReactionCard(ps *PlayerState, instanceID string) *CardInstance {
	if skill := e.findSkill(ps, instanceID); skill != nil {
		return skill
	}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := ps.Units[col][row]
			if unit != nil && unit.InstanceID == instanceID {
				return unit
			}
		}
	}
	for _, equipment := range ps.Equipment {
		if equipment != nil && equipment.InstanceID == instanceID {
			return equipment
		}
	}
	return nil
}

func (e *Engine) getAllFieldCards(ps *PlayerState) []*CardInstance {
	var cards []*CardInstance
	if ps.Hero != nil {
		cards = append(cards, ps.Hero)
	}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if ps.Units[col][row] != nil && ps.Units[col][row] != ps.Hero {
				cards = append(cards, ps.Units[col][row])
			}
		}
	}
	for i := range ps.Skills {
		if ps.Skills[i] != nil {
			cards = append(cards, ps.Skills[i])
		}
	}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if ps.Terrain[col][row] != nil {
				cards = append(cards, ps.Terrain[col][row])
			}
		}
	}
	for i := range ps.Equipment {
		if ps.Equipment[i] != nil {
			cards = append(cards, ps.Equipment[i])
		}
	}
	return cards
}

// spellModifierSources never inserts a virtual source into game state. Only
// explicitly opted-in self modifiers join this one evaluation, once by identity.
func (e *Engine) spellModifierSources(ps *PlayerState, skill *CardInstance) []*CardInstance {
	sources := e.getAllFieldCards(ps)
	if skill == nil || skill.Card == nil || skill.OwnerID != ps.PlayerID {
		return sources
	}
	if _, ok := behaviorForNumber(skill.Card.Number).(SpellSourceModifier); !ok {
		return sources
	}
	for _, source := range sources {
		if source == skill {
			return sources
		}
	}
	return append(sources, skill)
}
