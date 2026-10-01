package game

import (
	"eraofarcane/model"
)

type Card2311001ThunderSource struct{ AlwaysActive }

func (Card2311001ThunderSource) ID() string { return "2311001" }

func (Card2311001ThunderSource) Name() string { return "雷之源" }

func (Card2311001ThunderSource) ModifyCardPlayCost(ctx *EffectContext, card *CardInstance, cost map[string]int) {
	reduceCost(cost, model.ElementAir, 1)
}

func (Card2311001ThunderSource) ModifySkillUseCost(ctx *EffectContext, cost map[string]int) {
	reduceCost(cost, model.ElementAir, 1)
}
