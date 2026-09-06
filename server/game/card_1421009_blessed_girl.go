package game

import (
	"eraofarcane/model"
)

type Card1421009BlessedGirl struct{ AlwaysActive }

func (Card1421009BlessedGirl) ID() string { return "1421009" }

func (Card1421009BlessedGirl) Name() string { return "被祝福的少女" }

func (Card1421009BlessedGirl) IsPrayerAbility() bool { return true }

func (Card1421009BlessedGirl) OnPerTurn(ctx *EffectContext) error {
	if ctx.Source == nil || ctx.Source.Position == nil {
		return nil
	}
	candidates := make([]map[string]any, 0)
	for _, unit := range adjacentUnits(ctx.Engine.State.Players[ctx.PlayerID], ctx.Source.Position) {
		if unit.Card.IsCompanion() && unit.Card.Category == model.ElementEarth {
			candidates = append(candidates, candidateInfo(unit, "unit", "own"))
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "blessed_girl_load",
		"选择相邻地脉伙伴获得负载+1地", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			target := ctx.Engine.findFieldCardByInstance(ctx.Engine.State.Players[ctx.PlayerID], selected[0])
			if target != nil {
				ctx.Engine.addElementsGainBonus(target, ctx.PlayerID, model.ElementEarth, 1, ctx.Source)
			}
		})
	return nil
}
