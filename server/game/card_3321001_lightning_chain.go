package game

type Card3321001LightningChain struct{ AlwaysActive }

func (Card3321001LightningChain) ID() string { return "3321001" }

func (Card3321001LightningChain) Name() string { return "闪电链" }

func (Card3321001LightningChain) SpellDamage(ctx *EffectContext) int {
	return 1
}

func (Card3321001LightningChain) PrepareSpellCast(*EffectContext, SpellTarget, ActionMessage) (SpellCastOptions, error) {
	return SpellCastOptions{AllowExtraTarget: true}, nil
}
