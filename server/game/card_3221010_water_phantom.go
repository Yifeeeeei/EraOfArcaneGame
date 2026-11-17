package game

import (
	"eraofarcane/model"
)

type Card3221010WaterPhantom struct{ AlwaysActive }

func (Card3221010WaterPhantom) ID() string { return "3221010" }

func (Card3221010WaterPhantom) Name() string { return "水幻影" }

func (Card3221010WaterPhantom) OnSpellCast(ctx *EffectContext) error {
	if !isFriendlySpellCast(ctx) || ctx.Target != nil {
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card.Card.IsCompanion() && card.Card.Category == model.ElementWater && card.EnterTurn == ctx.Engine.State.TurnNumber
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "water_phantom_copy_target", "水幻影:选择本回合你召唤的1个水纹伙伴", candidates, 1, 1, func(selected []string) {
		target := selectedUnitFromCandidates(ctx.Engine, selected, candidates)
		if target == nil {
			return
		}
		positions := ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)
		if len(positions) == 0 {
			return
		}
		targetID := target.InstanceID
		ctx.Engine.SetPendingAction(ctx.PlayerID, "water_phantom_copy_position", "水幻影:选择复制体的入场位置", positions, 1, 1, func(posSelected []string) {
			pos, ok := positionFromSelectionID(firstSelected(posSelected))
			if !ok {
				return
			}
			ctx.Engine.summonWaterPhantomCopy(ctx.PlayerID, targetID, pos)
		})
	})
	return nil
}
