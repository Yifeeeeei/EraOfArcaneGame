package game

type Card1521107RadiantGuard struct{ AlwaysActive }

func (Card1521107RadiantGuard) ID() string   { return "1521107" }
func (Card1521107RadiantGuard) Name() string { return "辉之圣防军" }

func (Card1521107RadiantGuard) ModifySelfCardPlayCost(ctx *EffectContext, cost map[string]int) {
	if ctx == nil || ctx.Engine == nil {
		return
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ps == nil || !ps.FriendlyUnitDamagedLastTurn {
		return
	}
	for elem := range cost {
		cost[elem] = 0
	}
}

var _ SelfCardPlayCostModifier = Card1521107RadiantGuard{}
