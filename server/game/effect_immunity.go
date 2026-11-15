package game

// CardEffectImmunityBehavior protects printed rules from other card effects.
// Costs, orientation reset, and the card's own printed ability are still rules.
type CardEffectImmunityBehavior interface{ ImmuneToOtherCardEffects(*CardInstance) bool }

func cardImmuneToOtherEffects(card *CardInstance) bool {
	if card == nil || card.Card == nil {
		return false
	}
	b, ok := cardBehavior(card).(CardEffectImmunityBehavior)
	return ok && b.ImmuneToOtherCardEffects(card)
}
func (e *Engine) acceptsCardEffect(source, target *CardInstance) bool {
	return source == target || !cardImmuneToOtherEffects(target)
}
func (e *Engine) filterEffectCandidates(candidates []map[string]any) []map[string]any {
	result := make([]map[string]any, 0, len(candidates))
	var source *CardInstance
	if e.effectSource != nil {
		source = e.effectSource.state()
	}
	for _, info := range candidates {
		id, _ := info["instance_id"].(string)
		var card *CardInstance
		for _, ps := range e.State.Players {
			if card = e.findSkill(ps, id); card != nil {
				break
			}
		}
		if e.acceptsCardEffect(source, card) {
			result = append(result, info)
		}
	}
	return result
}
