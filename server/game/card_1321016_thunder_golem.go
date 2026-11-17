package game

import ()

type Card1321016ThunderGolem struct{ AlwaysActive }

func (Card1321016ThunderGolem) ID() string { return "1321016" }

func (Card1321016ThunderGolem) Name() string { return "雷傀儡" }

func (Card1321016ThunderGolem) OnDeath(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.OpponentID]
	if len(ps.Hand) == 0 {
		return nil
	}
	idx := ctx.Engine.randomIntn(len(ps.Hand))
	ctx.Engine.discardHandCardAt(ctx.OpponentID, idx)
	return nil
}
