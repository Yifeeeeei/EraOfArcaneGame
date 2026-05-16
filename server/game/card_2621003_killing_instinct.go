package game

type Card2621003KillingInstinct struct{}

func (Card2621003KillingInstinct) ID() string   { return "2621003" }
func (Card2621003KillingInstinct) Name() string { return "杀戮本能" }

func (Card2621003KillingInstinct) OnUnitEnter(ctx *EffectContext) error {
	if ctx.Target == nil || ctx.Target.OwnerID == ctx.PlayerID {
		return nil
	}
	ctx.Engine.dealDamage(ctx.Target, 2, ctx.Target.OwnerID)
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(ctx.Source),
		"target": cardToInfo(ctx.Target),
		"effect": "counter_damage",
		"amount": 2,
	}})
	return nil
}
