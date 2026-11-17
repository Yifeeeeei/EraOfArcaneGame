package game

type Card4511101ProtectorSival struct{ AlwaysActive }

func (Card4511101ProtectorSival) ID() string { return "4511101" }

func (Card4511101ProtectorSival) Name() string { return "庇护者 西瓦尔" }

func (Card4511101ProtectorSival) DamageScope() DamageScope { return DamageFriendly }

func (Card4511101ProtectorSival) OnDamaged(ctx *EffectContext, event DamageEvent) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || event.Target == nil || ctx.Source.UltimateUsed {
		return nil
	}
	damagedPlayer := event.Target.OwnerID
	if damagedPlayer != ctx.PlayerID || friendlyDamageTakenThisTurn(ctx.Engine, ctx.PlayerID) < 3 {
		return nil
	}
	source := ctx.Source
	if !ctx.Engine.cardStillOnField(source) {
		return nil
	}
	source.UltimateUsed = true
	source.Statuses[protectorSivalPreventionUntilStatus] = ctx.Engine.turnOrdinal() + 1
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(source),
		"effect": "prevent_friendly_damage",
	}})
	return nil
}

func (Card4511101ProtectorSival) PreventsFieldDamage(ctx *EffectContext) bool {
	if ctx == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.OwnerID != ctx.PlayerID {
		return false
	}
	return ctx.Source.Statuses[protectorSivalPreventionUntilStatus] >= ctx.Engine.turnOrdinal()
}

func (Card4511101ProtectorSival) PreventsDamage(ctx *EffectContext) bool {
	if ctx == nil || ctx.Source == nil || ctx.Engine == nil {
		return false
	}
	return ctx.Source.Statuses[protectorSivalPreventionUntilStatus] >= ctx.Engine.turnOrdinal()
}

const protectorSivalPreventionUntilStatus = "西瓦尔防止伤害至回合"

func friendlyDamageTakenThisTurn(e *Engine, playerID int) int {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return 0
	}
	return e.State.Players[playerID].friendlyUnitDamageThisTurn
}
