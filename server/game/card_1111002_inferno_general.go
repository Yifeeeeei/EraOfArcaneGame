package game

type Card1111002InfernoGeneral struct{ AlwaysActive }

func (Card1111002InfernoGeneral) ID() string { return "1111002" }

func (Card1111002InfernoGeneral) Name() string { return "炎狱大将军 狄斯托德" }

func (Card1111002InfernoGeneral) OnUnitEnter(ctx *EffectContext) error {
	if ctx.Target == nil || ctx.Target.OwnerID == ctx.PlayerID {
		return nil
	}
	statuses := map[string]int{}
	if ctx.Engine.addStatus(ctx.Target, StatusBurn, 2) {
		statuses[StatusBurn] = 2
	}
	if ctx.Engine.addStatus(ctx.Target, StatusPetrify, 2) {
		statuses[StatusPetrify] = 2
	}
	if len(statuses) == 0 {
		return nil
	}
	ctx.Engine.emit(GameEvent{
		Type:   "effect_trigger",
		Player: ctx.PlayerID,
		Data: map[string]any{
			"source":   cardToInfo(ctx.Source),
			"target":   cardToInfo(ctx.Target),
			"effect":   "apply_status",
			"statuses": statuses,
		},
	})
	return nil
}
