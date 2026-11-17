package game

type Card2621013WitchcraftRing struct{ AlwaysActive }

func (Card2621013WitchcraftRing) ID() string { return "2621013" }

func (Card2621013WitchcraftRing) Name() string { return "巫术指环" }

func (Card2621013WitchcraftRing) OnStatusGain(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.Card == nil || ctx.ExtraData == nil {
		return nil
	}
	if ctx.ExtraData["status"] != StatusWeaken || ctx.Target.OwnerID == ctx.PlayerID || !canInstanceBeWeakened(ctx.Target) || !useTriggeredTurn(ctx.Source) {
		return nil
	}
	ctx.Target.Statuses[StatusWeaken]++
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(ctx.Source),
		"target": cardToInfo(ctx.Target),
		"effect": "increase_weaken",
		"amount": 1,
	}})
	return nil
}
