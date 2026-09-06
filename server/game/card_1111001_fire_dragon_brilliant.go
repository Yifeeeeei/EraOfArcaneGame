package game

import (
	"eraofarcane/model"
)

type Card1111001FireDragonBrilliant struct{ AlwaysActive }

func (Card1111001FireDragonBrilliant) ID() string { return "1111001" }

func (Card1111001FireDragonBrilliant) Name() string { return `火龙 "辉煌"` }

func (Card1111001FireDragonBrilliant) DevourRequirement() map[string]int {
	return map[string]int{model.ElementFire: 3}
}

func (Card1111001FireDragonBrilliant) OnEnter(ctx *EffectContext) error {
	bindSkillToHost(ctx, "3101001")
	return nil
}
