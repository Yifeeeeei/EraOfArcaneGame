package game

import (
	"eraofarcane/model"
)

type Card2001101FinaleViolin struct{ AlwaysActive }

func (Card2001101FinaleViolin) ID() string { return "2001101" }

func (Card2001101FinaleViolin) Name() string { return "落幕提琴" }

func (Card2001101FinaleViolin) AttackCost(*EffectContext) map[string]int {
	return map[string]int{model.ElementArcane: 2}
}
