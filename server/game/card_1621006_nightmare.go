package game

type Card1621006Nightmare struct{ AlwaysActive }

func (Card1621006Nightmare) ID() string   { return "1621006" }
func (Card1621006Nightmare) Name() string { return "梦魇" }

func (Card1621006Nightmare) OnFriendlyDeath(ctx *EffectContext) error {
	if ctx.Target == nil || ctx.Target == ctx.Source {
		return nil
	}
	ctx.Engine.gainLife(ctx.Source, 1, ctx.Source)
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(ctx.Source),
		"target": cardToInfo(ctx.Target),
		"effect": "modify_life_self",
		"amount": 1,
	}})
	return nil
}
