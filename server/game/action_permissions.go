package game

// RuleAction describes a gameplay operation, including internal rule checks
// which do not originate from a WebSocket action.
type RuleAction string

const (
	RuleSummon       RuleAction = "summon"
	RuleEquip        RuleAction = "equip"
	RuleUseItem      RuleAction = "use_item"
	RulePlaceTerrain RuleAction = "place_terrain"
	RuleLearnSkill   RuleAction = "learn_skill"
	RuleUseSkill     RuleAction = "use_skill"
	RuleAttack       RuleAction = "attack"
)

type ActionRestrictionBehavior interface {
	HasActiveActionRestriction(*CardInstance) bool
	RestrictsAction(*CardInstance, RuleAction, *CardInstance) bool
}

func (AlwaysActive) HasActiveActionRestriction(*CardInstance) bool { return true }
func (e *Engine) actionRestricted(action RuleAction, card *CardInstance) bool {
	if e == nil {
		return false
	}
	for _, ps := range e.State.Players {
		if ps == nil {
			continue
		}
		for _, source := range e.getAllFieldCards(ps) {
			if source == nil || e.hasEffectiveStatus(source, StatusPetrify) {
				continue
			}
			b, ok := cardBehavior(source).(ActionRestrictionBehavior)
			if ok && b.HasActiveActionRestriction(source) && b.RestrictsAction(source, action, card) {
				return true
			}
		}
	}
	return false
}
