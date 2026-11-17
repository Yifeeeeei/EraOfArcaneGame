package game

func (Card3621107WillErosion) PrepareSpellCast(*EffectContext, SpellTarget, ActionMessage) (SpellCastOptions, error) {
	return SpellCastOptions{AllowExtraTarget: true}, nil
}

func (Card3621107WillErosion) SpellTargetGrant(ctx *EffectContext, skill *CardInstance, _ SpellTarget) SpellTargetGrant {
	if ctx.Source != skill {
		return SpellTargetGrant{}
	}
	return SpellTargetGrant{AllowSameExtraTarget: true, Pierce: ctx.Engine.redMoonActive(ctx.PlayerID)}
}
