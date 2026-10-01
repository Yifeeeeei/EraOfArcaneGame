package game

type Card2521106MoonlightScroll struct{ AlwaysActive }

func (Card2521106MoonlightScroll) ID() string { return "2521106" }

func (Card2521106MoonlightScroll) Name() string { return "沐光卷轴" }

func (Card2521106MoonlightScroll) OnUseItem(ctx *EffectContext) error {
	for _, unit := range royalFriendlyUnits(ctx) {
		ctx.Engine.healUnit(unit, 2, ctx.Source)
	}
	return nil
}
