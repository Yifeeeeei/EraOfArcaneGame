package game

type Card4111002WitchVerland struct{}

func (Card4111002WitchVerland) ID() string   { return "4111002" }
func (Card4111002WitchVerland) Name() string { return "女巫 维兰德" }
func (Card4111002WitchVerland) OnPerTurn(ctx *EffectContext) error {
	ctx.Source.Statuses[StatusBurn]++
	return nil
}
