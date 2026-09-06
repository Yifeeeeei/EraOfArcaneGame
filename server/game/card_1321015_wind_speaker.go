package game

import (
	"eraofarcane/model"
)

type Card1321015WindSpeaker struct{ AlwaysActive }

func (Card1321015WindSpeaker) ID() string { return "1321015" }

func (Card1321015WindSpeaker) Name() string { return "风语者" }

func (Card1321015WindSpeaker) OnDiscard(ctx *EffectContext) error {
	if ctx.ExtraData == nil {
		return nil
	}
	discardedPlayer, _ := ctx.ExtraData["discarded_player"].(int)
	if discardedPlayer == ctx.PlayerID && useTriggeredTurn(ctx.Source) {
		ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementAir: 1})
	}
	return nil
}
