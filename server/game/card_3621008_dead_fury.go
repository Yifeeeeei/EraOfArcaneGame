package game

type Card3621008DeadFury struct{ AlwaysActive }

func (Card3621008DeadFury) ID() string   { return "3621008" }
func (Card3621008DeadFury) Name() string { return "亡者之怒" }

func (Card3621008DeadFury) OnFriendlyDeath(ctx *EffectContext) error {
	addDeadFuryPower(ctx)
	return nil
}

func (Card3621008DeadFury) OnEnemyDeath(ctx *EffectContext) error {
	addDeadFuryPower(ctx)
	return nil
}

func addDeadFuryPower(ctx *EffectContext) {
	if ctx.Source == nil || ctx.Source.Card == nil || !ctx.Source.Card.IsSkill() {
		return
	}
	ctx.Source.PowerBonus++
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(ctx.Source),
		"effect": "skill_power_up",
		"amount": 1,
	}})
}
