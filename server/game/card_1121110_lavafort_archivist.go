package game

type Card1121110LavafortArchivist struct{ AlwaysActive }

func (Card1121110LavafortArchivist) ID() string { return "1121110" }

func (Card1121110LavafortArchivist) Name() string { return "熔岩堡档案员" }

func (Card1121110LavafortArchivist) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.Card == nil || ctx.Source.UltimateUsed {
		return nil
	}
	if !isFriendlySpellCast(ctx) || !hasCardTag(ctx.Target.Card, "创造") {
		return nil
	}
	drawn := ctx.Engine.flipDeckMatchesToHand(ctx.PlayerID, 1, 0, isRuneOrScroll)
	if len(drawn) > 0 {
		ctx.Source.UltimateUsed = true
	}
	return nil
}
