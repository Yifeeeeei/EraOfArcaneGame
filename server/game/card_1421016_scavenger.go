package game

import (
	"eraofarcane/model"
)

type Card1421016Scavenger struct{ AlwaysActive }

func (Card1421016Scavenger) ID() string { return "1421016" }

func (Card1421016Scavenger) Name() string { return "食腐者" }

func (Card1421016Scavenger) DamageScope() DamageScope { return DamageOtherFriendly }

func (Card1421016Scavenger) OnDamaged(ctx *EffectContext, event DamageEvent) error {
	if event.Target == nil || event.Target == ctx.Source {
		return nil
	}
	damagedPlayer := event.Target.OwnerID
	attacker, hasAttacker := event.SourcePlayer, event.SourcePlayer >= 0
	if damagedPlayer == ctx.PlayerID && hasAttacker && attacker != ctx.PlayerID {
		ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementEarth: 2})
	}
	return nil
}
