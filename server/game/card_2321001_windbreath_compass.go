package game

import (
	"eraofarcane/model"
)

type Card2321001WindbreathCompass struct{ AlwaysActive }

func (Card2321001WindbreathCompass) ID() string { return "2321001" }

func (Card2321001WindbreathCompass) Name() string { return "风息罗盘" }

func (Card2321001WindbreathCompass) OnDraw(ctx *EffectContext) error {
	if ctx.ExtraData == nil {
		return nil
	}
	if player, ok := ctx.ExtraData["drawn_player"].(int); !ok || player != ctx.PlayerID {
		return nil
	}
	drawn, _ := ctx.ExtraData["drawn_card"].(*CardInstance)
	if drawn == nil || drawn.Card == nil || drawn.Card.Category != model.ElementAir {
		return nil
	}
	ctx.Engine.SetTriggeredTurnAction(ctx.Source, ctx.PlayerID, "windbreath_compass",
		"风息罗盘:是否展示抽到的大气卡牌并临时获得负载+1气", []map[string]any{candidateInfo(drawn, "hand", "own")}, 0, 1,
		func(selected []string) {
			accepted := len(selected) > 0 && selected[0] == drawn.InstanceID && ctx.Engine.findFriendlyHandCard(ctx.PlayerID, drawn.InstanceID) == drawn
			if !accepted || !useTriggeredTurn(ctx.Source) {
				return
			}
			ps := ctx.Engine.State.Players[ctx.PlayerID]
			if ps.RevealedHand == nil {
				ps.RevealedHand = make(map[string]bool)
			}
			ps.RevealedHand[drawn.InstanceID] = true
			ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementAir, 1, ctx.Source)
			ctx.Source.Statuses[windbreathCompassTemporaryAirStatus]++
		})
	return nil
}

func (Card2321001WindbreathCompass) OnTurnEnd(ctx *EffectContext) error {
	count := ctx.Source.Statuses[windbreathCompassTemporaryAirStatus]
	if count > 0 {
		addElementsGainBonus(ctx.Source, model.ElementAir, -count)
		delete(ctx.Source.Statuses, windbreathCompassTemporaryAirStatus)
	}
	return nil
}
