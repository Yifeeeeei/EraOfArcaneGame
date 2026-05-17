package game

type Card3621007AndisPunishment struct{ AlwaysActive }

func (Card3621007AndisPunishment) ID() string   { return "3621007" }
func (Card3621007AndisPunishment) Name() string { return "安迪斯的惩罚" }

func (Card3621007AndisPunishment) OnDamaged(ctx *EffectContext) error {
	damagedPlayer, _ := ctx.ExtraData["damaged_player"].(int)
	if damagedPlayer == ctx.PlayerID {
		amount, _ := ctx.ExtraData["damage"].(int)
		ctx.Source.PowerBonus += max(amount, 0)
	}
	return nil
}
