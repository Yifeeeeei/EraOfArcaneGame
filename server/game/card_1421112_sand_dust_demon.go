package game

type Card1421112SandDustDemon struct{ AlwaysActive }

func (Card1421112SandDustDemon) ID() string { return "1421112" }

func (Card1421112SandDustDemon) Name() string { return "沙尘恶魔" }

func (Card1421112SandDustDemon) IsPrayerAbility() bool { return true }

func (Card1421112SandDustDemon) OnPerTurn(ctx *EffectContext) error {
	opponent := ctx.Engine.State.Players[ctx.OpponentID]
	frontRow := opponent.GetFrontRow()
	if frontRow < 0 || frontRow >= 3 {
		return nil
	}
	for col := 0; col < 3; col++ {
		if unit := opponent.Units[col][frontRow]; unit != nil {
			ctx.Engine.addStatus(unit, StatusPetrify, 1)
		}
	}
	return nil
}
