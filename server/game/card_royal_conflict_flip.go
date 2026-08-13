package game

import "eraofarcane/model"

type Card2121107SacredFireRally struct{ AlwaysActive }

func (Card2121107SacredFireRally) ID() string   { return "2121107" }
func (Card2121107SacredFireRally) Name() string { return "神火集结号" }
func (Card2121107SacredFireRally) OnEnter(ctx *EffectContext) error {
	ctx.Engine.flipDeckMatchesToHand(ctx.PlayerID, 2, 0, isFireCompanion)
	return nil
}

type Card2421110SandwormBait struct{ AlwaysActive }

func (Card2421110SandwormBait) ID() string   { return "2421110" }
func (Card2421110SandwormBait) Name() string { return "沙虫之饵" }
func (Card2421110SandwormBait) OnUseItem(ctx *EffectContext) error {
	drawn := ctx.Engine.flipDeckMatchesToHand(ctx.PlayerID, 1, 0, isEarthCompanionWithCostAboveFive)
	if len(drawn) > 0 && drawn[0].Card.Number == "1421114" {
		drawn[0].Statuses["入场费用"+model.ElementEarth+"-2"]++
	}
	return nil
}

type Card2521110AngelPrayer struct{ AlwaysActive }

func (Card2521110AngelPrayer) ID() string   { return "2521110" }
func (Card2521110AngelPrayer) Name() string { return "天使之祈祷" }
func (Card2521110AngelPrayer) OnUseItem(ctx *EffectContext) error {
	drawn := ctx.Engine.flipDeckMatchesToHand(ctx.PlayerID, 1, 0, isLightSpirit)
	if len(drawn) > 0 {
		drawn[0].Statuses["入场费用"+model.ElementLight+"-1"]++
	}
	return nil
}
