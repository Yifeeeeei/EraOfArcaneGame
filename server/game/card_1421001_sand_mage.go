package game

type Card1421001SandMage struct{}

func (Card1421001SandMage) ID() string   { return "1421001" }
func (Card1421001SandMage) Name() string { return "流沙法师" }
func (Card1421001SandMage) OnEnter(ctx *EffectContext) error {
	return ApplyStatusAuto(StatusPetrify, 1)(ctx)
}
