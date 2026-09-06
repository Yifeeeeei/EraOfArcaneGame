package game

type Card1121103BeaconGuard struct{ AlwaysActive }

func (Card1121103BeaconGuard) ID() string { return "1121103" }

func (Card1121103BeaconGuard) Name() string { return "烽火台守卫" }

func (Card1121103BeaconGuard) OnEnter(ctx *EffectContext) error {
	if royalCompanionCount(ctx.Engine.State.Players[ctx.PlayerID]) < royalCompanionCount(ctx.Engine.State.Players[ctx.OpponentID]) {
		ctx.Engine.gainPlayerShield(ctx.PlayerID, 3)
	}
	return nil
}

func royalCompanionCount(ps *PlayerState) int {
	if ps == nil {
		return 0
	}
	count := 0
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := ps.Units[col][row]
			if unit != nil && unit.Card != nil && unit.Card.IsCompanion() {
				count++
			}
		}
	}
	return count
}
