package game

type Card3421105AgingTouch struct{ AlwaysActive }

func (Card3421105AgingTouch) ID() string { return "3421105" }

func (Card3421105AgingTouch) Name() string { return "苍老之触" }

func (Card3421105AgingTouch) OnSpellHit(ctx *EffectContext) error {
	if ctx.Target == nil || ctx.Target.Card == nil || !ctx.Target.Card.IsCompanion() {
		return nil
	}
	setElementsGain(ctx.Target, map[string]int{})
	ctx.Target.ElementsGainBonus = make(map[string]int)
	return nil
}
