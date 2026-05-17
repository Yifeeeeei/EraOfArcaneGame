package game

type Card1221004FrostPuppet struct{ AlwaysActive }

func (Card1221004FrostPuppet) ID() string   { return "1221004" }
func (Card1221004FrostPuppet) Name() string { return "寒霜傀儡" }
func (Card1221004FrostPuppet) OnEnter(ctx *EffectContext) error {
	if ctx.Target != nil {
		ctx.Target.Statuses[StatusFreeze]++
		return nil
	}
	opponent := ctx.Engine.State.Players[ctx.OpponentID]
	frontRow := opponent.GetFrontRow()
	if frontRow < 0 {
		return nil
	}
	for col := 0; col < 3; col++ {
		if opponent.Units[col][frontRow] != nil {
			opponent.Units[col][frontRow].Statuses[StatusFreeze]++
			return nil
		}
	}
	return nil
}
