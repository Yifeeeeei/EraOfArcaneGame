package game

type Card1521001HealingWarlock struct{}

func (Card1521001HealingWarlock) ID() string   { return "1521001" }
func (Card1521001HealingWarlock) Name() string { return "治疗术士" }
func (Card1521001HealingWarlock) OnPerTurn(ctx *EffectContext) error {
	if ctx.Target == nil || ctx.Target.OwnerID != ctx.PlayerID {
		return nil
	}
	healUnit(ctx.Target, 1)
	return nil
}
