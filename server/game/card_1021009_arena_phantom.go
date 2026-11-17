package game

type Card1021009ArenaPhantom struct{ AlwaysActive }

func (Card1021009ArenaPhantom) ID() string { return "1021009" }

func (Card1021009ArenaPhantom) Name() string { return "竞技场虚像" }

func (Card1021009ArenaPhantom) PreventsDamage(ctx *EffectContext) bool {
	return ctx.ExtraData["damage_source"] != "spell"
}
