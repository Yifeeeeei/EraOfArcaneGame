package game

type Card3221108SixPetalSnowflake struct{ AlwaysActive }

func (Card3221108SixPetalSnowflake) ID() string { return "3221108" }

func (Card3221108SixPetalSnowflake) Name() string { return "六瓣雪花" }

func (Card3221108SixPetalSnowflake) SpellHitStatuses(ctx *EffectContext) map[string]int {
	if ctx == nil || ctx.Target == nil || ctx.Target.Card == nil || ctx.Target.Card.IsHero() {
		return nil
	}
	return map[string]int{StatusFreeze: 1}
}
