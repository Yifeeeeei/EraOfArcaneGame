package game

import (
	"eraofarcane/model"
)

type Card2621110AndisGift struct{ AlwaysActive }

func (Card2621110AndisGift) ID() string { return "2621110" }

func (Card2621110AndisGift) Name() string { return "安迪斯的赠与" }

func (Card2621110AndisGift) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "andis_gift_target",
		"安迪斯的赠与:选择1个友方单位获得负载+2暗,回合结束时死亡", candidates, 1, 1,
		func(selected []string) {
			target := ctx.Engine.findUnitByInstanceID(firstSelected(selected))
			if target == nil || target.OwnerID != ctx.PlayerID || target.Position == nil {
				return
			}
			ctx.Engine.addElementsGainBonus(target, ctx.PlayerID, model.ElementShadow, 2, ctx.Source)
			target.Statuses[andisGiftDoomedStatus] = ctx.Engine.State.TurnNumber
		})
	return nil
}

const andisGiftDoomedStatus = "安迪斯的赠与回合结束死亡"
