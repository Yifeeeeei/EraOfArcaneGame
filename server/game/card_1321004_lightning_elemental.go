package game

type Card1321004LightningElemental struct{ AlwaysActive }

func (Card1321004LightningElemental) ID() string   { return "1321004" }
func (Card1321004LightningElemental) Name() string { return "雷电元素" }
func (Card1321004LightningElemental) OnEnter(ctx *EffectContext) error {
	return ApplyStatusAuto(StatusStun, 1)(ctx)
}
