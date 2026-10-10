package game

import (
	"eraofarcane/model"
)

type Card2421001KnowledgeTreeCare struct{ AlwaysActive }

func (Card2421001KnowledgeTreeCare) ID() string { return "2421001" }

func (Card2421001KnowledgeTreeCare) Name() string { return "知识古树的关怀" }

func (Card2421001KnowledgeTreeCare) OnMasteryAchieved(ctx *EffectContext, level int) error {
	if ctx.Source == nil || ctx.Source.IsHorizontal {
		return nil
	}
	candidates := []map[string]any{candidateInfo(ctx.Source, "equipment", "own")}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "knowledge_tree_care",
		"你的卡牌达到精通，是否消耗知识古树的关怀抽1张牌并获得1地？", candidates, 0, 1,
		func(selected []string) {
			if len(selected) == 0 || selected[0] != ctx.Source.InstanceID || ctx.Source.IsHorizontal {
				return
			}
			ctx.Source.IsHorizontal = true
			ctx.Engine.drawCards(ctx.PlayerID, 1)
			ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementEarth: 1})
		})
	return nil
}
