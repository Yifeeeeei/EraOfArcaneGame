package game

type Card4511101ProtectorSival struct{ AlwaysActive }

func (Card4511101ProtectorSival) ID() string { return "4511101" }

func (Card4511101ProtectorSival) Name() string { return "庇护者 西瓦尔" }

func (Card4511101ProtectorSival) DamageScope() DamageScope { return DamageOtherFriendly }

func (Card4511101ProtectorSival) OnDamaged(ctx *EffectContext, event DamageEvent) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || event.Target == nil || ctx.Source.UltimateUsed {
		return nil
	}
	damagedPlayer := event.Target.OwnerID
	if damagedPlayer != ctx.PlayerID || friendlyDamageTakenThisTurn(ctx.Engine, ctx.PlayerID) < 3 {
		return nil
	}
	source := ctx.Source
	ctx.Engine.SetPendingAction(ctx.PlayerID, "protector_sival_prevent_all_damage",
		"庇护者 西瓦尔:是否发动绝技防止所有友方单位伤害", []map[string]any{candidateInfo(source, "hero", "own")}, 0, 1,
		func(selected []string) {
			if len(selected) == 0 || source.UltimateUsed || !ctx.Engine.cardStillOnField(source) {
				return
			}
			source.UltimateUsed = true
			source.Statuses[protectorSivalPreventionUntilStatus] = ctx.Engine.State.TurnNumber + 1
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source": cardToInfo(source),
				"effect": "prevent_friendly_damage",
			}})
		})
	return nil
}

func (Card4511101ProtectorSival) PreventsFieldDamage(ctx *EffectContext) bool {
	if ctx == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.OwnerID != ctx.PlayerID {
		return false
	}
	return ctx.Source.Statuses[protectorSivalPreventionUntilStatus] >= ctx.Engine.State.TurnNumber
}

func (Card4511101ProtectorSival) PreventsDamage(ctx *EffectContext) bool {
	if ctx == nil || ctx.Source == nil || ctx.Engine == nil {
		return false
	}
	return ctx.Source.Statuses[protectorSivalPreventionUntilStatus] >= ctx.Engine.State.TurnNumber
}

const protectorSivalPreventionUntilStatus = "西瓦尔防止伤害至回合"

func friendlyDamageTakenThisTurn(e *Engine, playerID int) int {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return 0
	}
	total := 0
	for _, card := range e.getAllFieldCards(e.State.Players[playerID]) {
		if card != nil && (card.Card.IsHero() || card.Card.IsCompanion()) {
			total += card.DamageTakenThisTurn
		}
	}
	return total
}
