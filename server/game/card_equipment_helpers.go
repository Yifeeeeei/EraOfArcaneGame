package game

func (e *Engine) equipmentInOwnerSlot(playerID int, target *CardInstance) bool {
	if e == nil || target == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return false
	}
	for _, card := range e.State.Players[playerID].Equipment {
		if card == target {
			return true
		}
	}
	return false
}

func (e *Engine) firstFreeEquipmentSlot(playerID int) int {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return -1
	}
	ps := e.State.Players[playerID]
	for i := 0; i < equipmentSlotCapacity(ps); i++ {
		if ps.Equipment[i] == nil {
			return i
		}
	}
	return -1
}
