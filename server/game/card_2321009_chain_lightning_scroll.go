package game

type Card2321009ChainLightningScroll struct{ AlwaysActive }

func (Card2321009ChainLightningScroll) ID() string   { return "2321009" }
func (Card2321009ChainLightningScroll) Name() string { return "连锁闪电卷轴" }

func (Card2321009ChainLightningScroll) OnSpellHit(ctx *EffectContext) error {
	if !isOwnSpellHit(ctx) {
		return nil
	}
	return DrawCards(1)(ctx)
}
