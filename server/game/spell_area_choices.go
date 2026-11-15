package game

import "fmt"

// Area selection belongs to the cast, not to a mutable global card definition.
type SpellAreaChoiceBehavior interface {
	SpellAreaChoices(*EffectContext) []SpellArea
}

func (e *Engine) spellAreaChoices(skill *CardInstance) []SpellArea {
	if skill == nil || skill.Card == nil {
		return nil
	}
	if b, ok := cardBehavior(skill).(SpellAreaChoiceBehavior); ok {
		return b.SpellAreaChoices(e.skillContext(skill.OwnerID, skill))
	}
	return nil
}
func (e *Engine) spellAreaForTarget(skill *CardInstance, target SpellTarget) SpellArea {
	if target.Area != "" {
		return target.Area
	}
	return e.effectiveSpellArea(skill)
}
func (e *Engine) chooseCastSpellArea(playerID int, skill *CardInstance, action ActionMessage) bool {
	choices := e.spellAreaChoices(skill)
	if len(choices) < 2 || action.Data["area"] != nil {
		return false
	}
	labels := map[SpellArea]string{SpellAreaFrontRow: "前排", SpellAreaColumn: "纵列", SpellAreaSquare: "方阵", SpellAreaSplashCross: "溅射"}
	candidates := []map[string]any{}
	for _, area := range choices {
		candidates = append(candidates, map[string]any{"instance_id": string(area), "name": labels[area], "zone": "choice"})
	}
	e.SetPendingActionWithError(playerID, "spell_area_choice", "选择"+skill.Card.Name+"的范围", candidates, 1, 1, nil, false, func(selected []string, _ map[string]any) error {
		data := make(map[string]any, len(action.Data)+1)
		for k, v := range action.Data {
			data[k] = v
		}
		data["area"] = firstSelected(selected)
		return e.handleCastSpell(playerID, ActionMessage{Data: data})
	})
	return true
}
func (e *Engine) parseSpellAreaChoice(skill *CardInstance, target SpellTarget, action ActionMessage) (SpellTarget, error) {
	raw, ok := action.Data["area"].(string)
	if !ok {
		return target, nil
	}
	for _, choice := range e.spellAreaChoices(skill) {
		if raw == string(choice) {
			target.Area = choice
			return target, nil
		}
	}
	return target, fmt.Errorf("invalid spell area choice")
}
func (Card3511101DivineRadianceSkyward) SpellAreaChoices(*EffectContext) []SpellArea {
	return []SpellArea{SpellAreaFrontRow, SpellAreaColumn}
}
