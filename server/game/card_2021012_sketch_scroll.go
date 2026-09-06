package game

import "fmt"

type Card2021012SketchScroll struct{ AlwaysActive }

func (Card2021012SketchScroll) ID() string { return "2021012" }

func (Card2021012SketchScroll) Name() string { return "速写卷轴" }

func (Card2021012SketchScroll) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.sketchScrollSkillCandidates(ctx.PlayerID)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "sketch_scroll_skill",
		"选择1个已学习法术释放，本次无需消耗该技能", candidates, 1, 1,
		nil, false, func(selected []string, _ map[string]any) error {
			if len(selected) == 0 {
				return nil
			}
			return ctx.Engine.resolveSketchScrollSkill(ctx.PlayerID, selected[0])
		})
	return nil
}

func (Card2021012SketchScroll) ValidateItemUse(ctx *EffectContext) error {
	e, playerID := ctx.Engine, ctx.PlayerID
	if len(e.sketchScrollSkillCandidates(playerID)) == 0 {
		return fmt.Errorf("Sketch Scroll requires a payable learned attack spell")
	}
	return nil
}
