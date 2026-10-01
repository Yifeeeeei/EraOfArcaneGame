package game

import (
	"eraofarcane/model"
)

type Card2421102RoseWhip struct{ AlwaysActive }

func (Card2421102RoseWhip) ID() string { return "2421102" }

func (Card2421102RoseWhip) Name() string { return "蔷薇之鞭" }

func (Card2421102RoseWhip) OnLoadLoss(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.ExtraData == nil {
		return nil
	}
	lossPlayer, _ := ctx.ExtraData["load_loss_player"].(int)
	if lossPlayer != ctx.PlayerID || ctx.Target.OwnerID != ctx.PlayerID || ctx.Target == ctx.Source || !ctx.Engine.cardStillOnField(ctx.Source) {
		return nil
	}
	currentBonus := ctx.Source.ElementsGainBonus[model.ElementShadow]
	if currentBonus >= 2 {
		return nil
	}
	ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementShadow, 1, ctx.Source)
	ctx.Engine.emit(GameEvent{
		Type:   "rose_whip_load_gain",
		Player: -1,
		Data: map[string]any{
			"player": ctx.PlayerID,
			"source": cardToInfo(ctx.Source),
			"target": cardToInfo(ctx.Target),
		},
	})
	return nil
}

var _ OnLoadLossBehavior = Card2421102RoseWhip{}
