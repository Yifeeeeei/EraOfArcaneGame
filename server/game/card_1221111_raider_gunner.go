package game

type Card1221111RaiderGunner struct{ AlwaysActive }

func (Card1221111RaiderGunner) ID() string   { return "1221111" }
func (Card1221111RaiderGunner) Name() string { return "掠夺者炮手" }

func (Card1221111RaiderGunner) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Source == nil || ctx.Source.UltimateUsed || !isFriendlySpellHit(ctx) {
		return nil
	}
	if ctx.Engine.discardRandomHandCard(ctx.OpponentID) != nil {
		ctx.Source.UltimateUsed = true
	}
	return nil
}
