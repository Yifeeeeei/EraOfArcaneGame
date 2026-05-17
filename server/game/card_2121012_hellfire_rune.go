package game

type Card2121012HellfireRune struct{ AlwaysActive }

func (Card2121012HellfireRune) ID() string   { return "2121012" }
func (Card2121012HellfireRune) Name() string { return "狱火符文" }

func (Card2121012HellfireRune) OnUnitEnter(ctx *EffectContext) error {
	if ctx.Target == nil || ctx.Target.OwnerID == ctx.PlayerID {
		return nil
	}
	ctx.Target.Statuses[StatusStun] += 2
	ctx.Target.Statuses[StatusPetrify] += 2
	ctx.Target.Statuses[StatusBurn] += 2
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(ctx.Source),
		"target": cardToInfo(ctx.Target),
		"effect": "apply_status",
		"statuses": map[string]int{
			StatusStun:    2,
			StatusPetrify: 2,
			StatusBurn:    2,
		},
	}})
	return nil
}
