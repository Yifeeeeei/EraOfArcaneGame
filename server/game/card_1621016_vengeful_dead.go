package game

type Card1621016VengefulDead struct{ AlwaysActive }

func (Card1621016VengefulDead) ID() string { return "1621016" }

func (Card1621016VengefulDead) Name() string { return "复仇死者" }

func (Card1621016VengefulDead) OnDeath(ctx *EffectContext) error {
	sourcePlayer := ctx.Source.Statuses["lethal_source_player"] - 1
	if sourcePlayer < 0 || sourcePlayer >= len(ctx.Engine.State.Players) {
		return nil
	}
	hero := ctx.Engine.State.Players[sourcePlayer].Hero
	if hero != nil {
		ctx.Engine.ApplyDamage(DamageRequest{Target: hero, Amount: 2, Kind: "effect", SourcePlayer: ctx.PlayerID, SourceKnown: true, Source: ctx.Source})
	}
	return nil
}
