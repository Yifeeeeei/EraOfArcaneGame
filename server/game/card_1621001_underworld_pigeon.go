package game

type Card1621001UnderworldPigeon struct{ AlwaysActive }

func (Card1621001UnderworldPigeon) ID() string { return "1621001" }

func (Card1621001UnderworldPigeon) Name() string { return "冥界信鸽" }

func (Card1621001UnderworldPigeon) OnDeath(ctx *EffectContext) error {
	return DrawCards(1)(ctx)
}
