package game

import "eraofarcane/model"

type Card1021010SpecialistMage struct{}

func (Card1021010SpecialistMage) ID() string   { return "1021010" }
func (Card1021010SpecialistMage) Name() string { return "专精法师" }

func (Card1021010SpecialistMage) OnEnter(ctx *EffectContext) error {
	amount := totalLoad(ctx.Source)
	if amount <= 0 {
		amount = 1
	}
	candidates := make([]map[string]any, 0, len(model.AllElements))
	for _, elem := range model.AllElements {
		candidates = append(candidates, map[string]any{
			"instance_id": elem,
			"number":      "1021010",
			"name":        elem,
			"type":        "属性",
			"zone":        "choice",
			"side":        "own",
		})
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "specialist_mage_element",
		"选择专精法师的负载属性", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			setElementsGain(ctx.Source, map[string]int{selected[0]: amount})
		})
	return nil
}
