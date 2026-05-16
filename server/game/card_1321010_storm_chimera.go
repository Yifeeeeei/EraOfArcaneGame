package game

import "eraofarcane/model"

type Card1321010StormChimera struct{}

func (Card1321010StormChimera) ID() string   { return "1321010" }
func (Card1321010StormChimera) Name() string { return "风暴奇美拉" }

func (Card1321010StormChimera) DevourRequirement() map[string]int {
	return map[string]int{model.ElementAir: 3}
}

func (Card1321010StormChimera) ModifySkillUseCost(ctx *EffectContext, cost map[string]int) {
	if ctx.Source == nil || ctx.Source.Card == nil {
		return
	}
	if ctx.Source.Card.IsSkill() && ctx.Source.Card.Category == model.ElementAir && cost[model.ElementAir] > 0 {
		cost[model.ElementAir]--
		if cost[model.ElementAir] <= 0 {
			delete(cost, model.ElementAir)
		}
	}
}
