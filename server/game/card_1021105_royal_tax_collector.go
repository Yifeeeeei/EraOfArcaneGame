package game

import (
	"eraofarcane/model"
)

type Card1021105RoyalTaxCollector struct{ AlwaysActive }

func (Card1021105RoyalTaxCollector) ID() string { return "1021105" }

func (Card1021105RoyalTaxCollector) Name() string { return "皇城征税员" }

func (Card1021105RoyalTaxCollector) OnEnter(ctx *EffectContext) error {
	if ctx.Source == nil {
		return nil
	}
	ctx.Source.Statuses[royalTaxCollectorUntilOpponentTurnEndStatus] = ctx.Engine.State.TurnNumber
	return nil
}

func (Card1021105RoyalTaxCollector) OnDraw(ctx *EffectContext) error {
	if ctx.Source == nil || ctx.Source.Statuses[royalTaxCollectorUntilOpponentTurnEndStatus] <= 0 || ctx.ExtraData == nil {
		return nil
	}
	drawnPlayer, _ := ctx.ExtraData["drawn_player"].(int)
	if drawnPlayer != ctx.OpponentID {
		return nil
	}
	ctx.Engine.State.Players[ctx.PlayerID].Elements[model.ElementArcane]++
	return nil
}

func (Card1021105RoyalTaxCollector) OnTurnEnd(ctx *EffectContext) error {
	if ctx.Source == nil || ctx.Source.Statuses[royalTaxCollectorUntilOpponentTurnEndStatus] <= 0 || ctx.ExtraData == nil {
		return nil
	}
	endedPlayer, _ := ctx.ExtraData["ended_player"].(int)
	if endedPlayer == ctx.OpponentID && ctx.Engine.State.TurnNumber >= ctx.Source.Statuses[royalTaxCollectorUntilOpponentTurnEndStatus] {
		delete(ctx.Source.Statuses, royalTaxCollectorUntilOpponentTurnEndStatus)
	}
	return nil
}

const royalTaxCollectorUntilOpponentTurnEndStatus = "皇城征税员征税至对手回合结束"
