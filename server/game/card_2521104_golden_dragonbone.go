package game

import (
	"fmt"
)

type Card2521104GoldenDragonbone struct{ AlwaysActive }

func (Card2521104GoldenDragonbone) ID() string { return "2521104" }

func (Card2521104GoldenDragonbone) Name() string { return "黄金龙骨" }

func (Card2521104GoldenDragonbone) PerTurnLabel(*CardInstance) string {
	return "主动"
}

func (Card2521104GoldenDragonbone) OnPerTurn(ctx *EffectContext) error {
	if !ctx.Engine.sacrificeEquipment(ctx.PlayerID, ctx.Source.InstanceID) {
		return fmt.Errorf("golden dragonbone must be sacrificed from equipment")
	}
	ctx.Engine.drawCards(ctx.PlayerID, 2)
	return nil
}
