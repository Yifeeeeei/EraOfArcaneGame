package game

type Card2621005SacrificeRune struct{}

func (Card2621005SacrificeRune) ID() string   { return "2621005" }
func (Card2621005SacrificeRune) Name() string { return "献祭符文" }

func (Card2621005SacrificeRune) OnFriendlyDeath(ctx *EffectContext) error {
	return DrawCards(2)(ctx)
}
