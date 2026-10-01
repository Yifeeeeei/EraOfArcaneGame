package game

import (
	"fmt"
)

type Card2221104WaterMirrorScroll struct{ AlwaysActive }

func (Card2221104WaterMirrorScroll) ID() string { return "2221104" }

func (Card2221104WaterMirrorScroll) Name() string { return "水镜卷轴" }

func (Card2221104WaterMirrorScroll) OnUseItem(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.PlayerID < 0 || ctx.PlayerID >= len(ctx.Engine.State.Players) {
		return nil
	}
	recorded := ctx.Engine.State.Players[ctx.PlayerID].LastLowCostWaterSpell
	if recorded == nil || recorded.Card == nil {
		return nil
	}
	virtual := ctx.Engine.cloneVirtualSpell(recorded, ctx.PlayerID, ctx.Engine.State.TurnNumber)
	targets := ctx.Engine.spellTargetCandidates(ctx.PlayerID, virtual)
	if skillNeedsTargetInstance(virtual) {
		if len(targets) == 0 {
			return nil
		}
		ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "water_mirror_scroll_target",
			fmt.Sprintf("水镜卷轴:选择%s的目标", virtual.Card.Name), targets, 1, 1, nil, false,
			func(selected []string, _ map[string]any) error {
				target := selectedSpellTargetFromCandidates(ctx.Engine, ctx.PlayerID, virtual, firstSelected(selected), targets)
				if target == nil {
					return fmt.Errorf("invalid water mirror scroll target")
				}
				return ctx.Engine.startVirtualSpellCastNoBoost(ctx.PlayerID, virtual, *target, map[string]any{
					"triggered_by":        "2221104",
					"source_item":         cardToInfo(ctx.Source),
					"source_skill_hidden": false,
				})
			})
		return nil
	}
	return ctx.Engine.startVirtualSpellCastNoBoost(ctx.PlayerID, virtual, SpellTarget{Type: "none"}, map[string]any{
		"triggered_by":        "2221104",
		"source_item":         cardToInfo(ctx.Source),
		"source_skill_hidden": false,
	})
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
