package game

type Card2011003KingRobe struct{ AlwaysActive }

func (Card2011003KingRobe) ID() string { return "2011003" }

func (Card2011003KingRobe) Name() string { return "君王法袍 至贤" }

func (Card2011003KingRobe) OnUltimate(ctx *EffectContext) error {
	diff := totalFieldLoad(ctx.Engine.State.Players[ctx.PlayerID]) - totalFieldLoad(ctx.Engine.State.Players[ctx.OpponentID])
	amount := diff / 2
	if amount <= 0 {
		return nil
	}
	ctx.Source.Statuses[kingRobeReductionStatus] = amount
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(ctx.Source),
		"effect": "enemy_spell_damage_modifier",
		"amount": -amount,
	}})
	return nil
}

func (Card2011003KingRobe) ModifyEnemySpellStats(ctx *EffectContext, stats *SpellStats) {
	amount := ctx.Source.Statuses[kingRobeReductionStatus]
	if amount > 0 {
		stats.DamageBonus -= amount
	}
}

func (Card2011003KingRobe) OnTurnEnd(ctx *EffectContext) error {
	delete(ctx.Source.Statuses, kingRobeReductionStatus)
	return nil
}
