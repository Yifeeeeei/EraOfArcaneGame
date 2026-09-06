package game

import (
	"eraofarcane/model"
)

type Card1521114HuiPrayer struct{ AlwaysActive }

func (Card1521114HuiPrayer) ID() string { return "1521114" }

func (Card1521114HuiPrayer) Name() string { return "辉之都祈祷者" }

func (Card1521114HuiPrayer) OnEnter(ctx *EffectContext) error {
	wounded := 0
	for _, unit := range royalFriendlyUnits(ctx) {
		if unit != nil && unit.CurrentLife < maxLife(unit) {
			wounded++
		}
	}
	if wounded > 0 {
		ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementLight: wounded})
	}
	return nil
}
