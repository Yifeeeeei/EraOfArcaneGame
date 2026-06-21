package game

type Card1111002InfernoGeneral struct{ AlwaysActive }

func (Card1111002InfernoGeneral) ID() string   { return "1111002" }
func (Card1111002InfernoGeneral) Name() string { return "炎狱大将军 狄斯托德" }

func (Card1111002InfernoGeneral) OnUnitEnter(ctx *EffectContext) error {
	if ctx.Target == nil || ctx.Target.OwnerID == ctx.PlayerID {
		return nil
	}
	ctx.Target.Statuses[StatusBurn]++
	ctx.Target.Statuses[StatusPetrify] += 2
	ctx.Engine.emit(GameEvent{
		Type:   "effect_trigger",
		Player: ctx.PlayerID,
		Data: map[string]any{
			"source":   cardToInfo(ctx.Source),
			"target":   cardToInfo(ctx.Target),
			"effect":   "apply_status",
			"statuses": map[string]int{StatusBurn: 1, StatusPetrify: 2},
		},
	})
	return nil
}
