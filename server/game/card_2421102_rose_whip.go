package game

import (
	"eraofarcane/model"
)

type Card2421102RoseWhip struct{ AlwaysActive }

func (Card2421102RoseWhip) ID() string { return "2421102" }

func (Card2421102RoseWhip) Name() string { return "蔷薇之鞭" }

func (Card2421102RoseWhip) OnLoadLoss(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil || !triggeredTurnAvailable(ctx.Source) {
		return nil
	}
	lossPlayer, _ := ctx.ExtraData["load_loss_player"].(int)
	target, _ := ctx.ExtraData["load_loss_target"].(*CardInstance)
	if lossPlayer != ctx.PlayerID || target == nil || target.OwnerID != ctx.PlayerID || !ctx.Engine.cardStillOnField(ctx.Source) {
		return nil
	}
	currentBonus := ctx.Source.ElementsGainBonus[model.ElementShadow]
	if currentBonus >= 2 {
		return nil
	}
	if !useTriggeredTurn(ctx.Source) {
		return nil
	}
	ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementShadow, 1, ctx.Source)
	ctx.Engine.emit(GameEvent{
		Type:   "rose_whip_load_gain",
		Player: -1,
		Data: map[string]any{
			"player": ctx.PlayerID,
			"source": cardToInfo(ctx.Source),
			"target": cardToInfo(target),
		},
	})
	return nil
}

var _ OnLoadLossBehavior = Card2421102RoseWhip{}
