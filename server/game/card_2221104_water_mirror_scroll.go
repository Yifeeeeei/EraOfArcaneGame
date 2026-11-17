package game

import (
	"fmt"
)

type Card2221104WaterMirrorScroll struct{ AlwaysActive }

func (Card2221104WaterMirrorScroll) ID() string { return "2221104" }

func (Card2221104WaterMirrorScroll) Name() string { return "水镜卷轴" }

func (Card2221104WaterMirrorScroll) AutomaticSpellCopy() bool { return true }
func (Card2221104WaterMirrorScroll) SpellCopies(ctx *EffectContext) []*CardInstance {
	recorded := ctx.Engine.State.Players[ctx.PlayerID].LastLowCostWaterSpell
	if recorded == nil {
		return nil
	}
	return []*CardInstance{recorded}
}

func (Card2221104WaterMirrorScroll) ValidateItemUse(ctx *EffectContext) error {
	e, playerID := ctx.Engine, ctx.PlayerID
	recorded := e.State.Players[playerID].LastLowCostWaterSpell
	if recorded == nil || recorded.Card == nil {
		return fmt.Errorf("Water Mirror Scroll requires a previous low-cost water spell")
	}
	if skillNeedsTargetInstance(recorded) && len(e.spellTargetCandidates(playerID, recorded)) == 0 {
		return fmt.Errorf("Water Mirror Scroll requires a legal target")
	}
	return nil
}
