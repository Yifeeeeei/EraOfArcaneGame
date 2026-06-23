package game

type Card2121012HellfireRune struct{ AlwaysActive }

func (Card2121012HellfireRune) ID() string   { return "2121012" }
func (Card2121012HellfireRune) Name() string { return "狱火符文" }

func (Card2121012HellfireRune) OnUnitEnter(ctx *EffectContext) error {
	if ctx.Target == nil || ctx.Target.OwnerID == ctx.PlayerID {
		return nil
	}
	statuses := map[string]int{}
	if ctx.Engine.addStatus(ctx.Target, StatusStun, 2) {
		statuses[StatusStun] = 2
	}
	if ctx.Engine.addStatus(ctx.Target, StatusPetrify, 2) {
		statuses[StatusPetrify] = 2
	}
	if ctx.Engine.addStatus(ctx.Target, StatusBurn, 2) {
		statuses[StatusBurn] = 2
	}
	if len(statuses) == 0 {
		return nil
	}
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source":   cardToInfo(ctx.Source),
		"target":   cardToInfo(ctx.Target),
		"effect":   "apply_status",
		"statuses": statuses,
	}})
	return nil
}
