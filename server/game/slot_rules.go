package game

// SlotGrant declares equipment-provided capacity. Group is an explicit stacking
// key: repeated grants in a nonempty group apply once. Restricted skill slots
// follow unrestricted slots and accept any tag listed in their grant.
type SlotGrant struct {
	Group                      string
	SkillSlots                 int
	SkillTags                  []string
	EquipmentSlots             int
	DuplicateEquipmentSubtypes bool
}

type SlotGrantBehavior interface {
	HasActiveSlotGrant(*CardInstance) bool
	SlotGrant(*CardInstance) SlotGrant
}

func (AlwaysActive) HasActiveSlotGrant(*CardInstance) bool { return true }

func equipmentSlotGrants(ps *PlayerState) []SlotGrant {
	if ps == nil {
		return nil
	}
	var grants []SlotGrant
	seen := map[string]bool{}
	for _, source := range ps.Equipment {
		if source == nil || source.Statuses[StatusPetrify] > 0 {
			continue
		}
		b, ok := cardBehavior(source).(SlotGrantBehavior)
		if !ok || !b.HasActiveSlotGrant(source) {
			continue
		}
		grant := b.SlotGrant(source)
		if grant.Group != "" {
			if seen[grant.Group] {
				continue
			}
			seen[grant.Group] = true
		}
		grants = append(grants, grant)
	}
	return grants
}

func skillSlotCapacity(ps *PlayerState) int {
	capacity := BaseSkillSlots
	for _, grant := range equipmentSlotGrants(ps) {
		capacity += grant.SkillSlots
	}
	return min(capacity, MaxSkillSlots)
}

func baseSkillSlotCapacity(ps *PlayerState) int {
	capacity := BaseSkillSlots
	for _, grant := range equipmentSlotGrants(ps) {
		if len(grant.SkillTags) == 0 {
			capacity += grant.SkillSlots
		}
	}
	return min(capacity, MaxSkillSlots)
}

func skillAllowedInSlot(ps *PlayerState, skill *CardInstance, slotIdx int) bool {
	if slotIdx < 0 || slotIdx >= skillSlotCapacity(ps) {
		return false
	}
	cursor := baseSkillSlotCapacity(ps)
	if slotIdx < cursor {
		return true
	}
	for _, grant := range equipmentSlotGrants(ps) {
		if len(grant.SkillTags) == 0 {
			continue
		}
		if slotIdx < cursor+grant.SkillSlots {
			if skill == nil || skill.Card == nil {
				return false
			}
			for _, tag := range grant.SkillTags {
				if hasCardTag(skill.Card, tag) {
					return true
				}
			}
			return false
		}
		cursor += grant.SkillSlots
	}
	return false
}

func equipmentSlotCapacity(ps *PlayerState) int {
	capacity := BaseEquipmentSlots
	for _, grant := range equipmentSlotGrants(ps) {
		capacity += grant.EquipmentSlots
	}
	return min(capacity, MaxEquipmentSlots)
}

func playerCanEquipDuplicateSubtypes(ps *PlayerState) bool {
	for _, grant := range equipmentSlotGrants(ps) {
		if grant.DuplicateEquipmentSubtypes {
			return true
		}
	}
	return false
}
