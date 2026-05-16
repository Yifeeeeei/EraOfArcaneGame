package game

type Card2521011FlashRune struct{}

func (Card2521011FlashRune) ID() string   { return "2521011" }
func (Card2521011FlashRune) Name() string { return "闪光符文" }

func (Card2521011FlashRune) OnSpellCast(ctx *EffectContext) error {
	if !isEnemySpellCast(ctx) || ctx.Target == nil || !ctx.Target.Card.IsSkill() {
		return nil
	}
	opponent := ctx.Engine.State.Players[ctx.OpponentID]
	for col := 0; col < 3; col++ {
		unit := opponent.Units[col][0]
		if unit == nil {
			continue
		}
		unit.Statuses[StatusStun]++
		ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
			"source": cardToInfo(ctx.Source),
			"target": cardToInfo(unit),
			"effect": "apply_status",
			"status": StatusStun,
			"amount": 1,
		}})
	}
	return nil
}
