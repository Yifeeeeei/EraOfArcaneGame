package game

import "eraofarcane/model"

type Card1521014TorchWitch struct{ AlwaysActive }

func (Card1521014TorchWitch) ID() string            { return "1521014" }
func (Card1521014TorchWitch) Name() string          { return "炬之女巫" }
func (Card1521014TorchWitch) IsPrayerAbility() bool { return true }
func (Card1521014TorchWitch) OnEnter(ctx *EffectContext) error {
	return ApplyStatusToSelf(StatusBurn, 2)(ctx)
}

func (Card1521014TorchWitch) OnPerTurn(ctx *EffectContext) error {
	if ctx.Source == nil || ctx.Source.Position == nil {
		return nil
	}
	candidates := make([]map[string]any, 0)
	for _, unit := range adjacentUnits(ctx.Engine.State.Players[ctx.PlayerID], ctx.Source.Position) {
		if unit.Card.IsCompanion() {
			candidates = append(candidates, candidateInfo(unit, "unit", "own"))
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "torch_witch_prayer",
		"炬之女巫:选择相邻伙伴获得负载+1光", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			target := ctx.Engine.findFieldCardByInstance(ctx.Engine.State.Players[ctx.PlayerID], selected[0])
			if target != nil {
				ctx.Engine.addElementsGainBonus(target, ctx.PlayerID, model.ElementLight, 1, ctx.Source)
			}
		})
	return nil
}
