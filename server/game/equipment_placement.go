package game

import "fmt"

func cardCanBeReplaced(card *CardInstance) bool {
	return card != nil && !card.IsHorizontal && card.Statuses[StatusCooldown] <= 0
}

// Validate the resulting field, including any capacity/duplicate grant removed
// by replacement. A weapon may replace any legal slot if the result is legal.
func prepareEquipmentPlacement(ps *PlayerState, card *CardInstance, replaceID string) (int, *CardInstance, error) {
	slot := -1
	var replaced *CardInstance
	for i := 0; i < equipmentSlotCapacity(ps); i++ {
		old := ps.Equipment[i]
		if replaceID != "" {
			if old != nil && old.InstanceID == replaceID {
				if !cardCanBeReplaced(old) {
					return -1, nil, fmt.Errorf("equipment cannot be replaced while horizontal or cooling down")
				}
				slot, replaced = i, old
				break
			}
		} else if old == nil {
			slot = i
			break
		}
	}
	if slot < 0 {
		return -1, nil, fmt.Errorf("equipment area is full or replacement not found")
	}
	projected := *ps
	projected.Equipment[slot] = card
	if slot >= equipmentSlotCapacity(&projected) {
		return -1, nil, fmt.Errorf("replacement removes the required equipment slot")
	}
	if !playerCanEquipDuplicateSubtypes(&projected) {
		seen := map[string]bool{}
		for _, item := range projected.Equipment {
			if item == nil || item.IsSetCounter {
				continue
			}
			subtype := restrictedEquipmentSubtype(item.Card)
			if subtype == "" {
				continue
			}
			if seen[subtype] {
				return -1, nil, fmt.Errorf("same subtype equipment must be replaced")
			}
			seen[subtype] = true
		}
	}
	return slot, replaced, nil
}
