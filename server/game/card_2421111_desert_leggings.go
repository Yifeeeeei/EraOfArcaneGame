package game

type Card2421111DesertLeggings struct{ AlwaysActive }

func (Card2421111DesertLeggings) ID() string { return "2421111" }

func (Card2421111DesertLeggings) Name() string { return "沙漠护腿" }

func (Card2421111DesertLeggings) ModifyFieldDamageAmount(ctx *EffectContext, amount int) int {
	if ctx == nil || ctx.Source == nil || ctx.Source.UltimateUsed || ctx.Target == nil || ctx.Target.Card == nil || amount < 2 {
		return amount
	}
	if !ctx.Target.Card.IsCompanion() {
		return amount
	}
	return max(amount-2, 0)
}

var _ FieldDamageAmountModifier = Card2421111DesertLeggings{}

func (Card2421111DesertLeggings) ConsumeDamageModifier(ctx *EffectContext, before, after int) {
	if after < before {
		ctx.Source.UltimateUsed = true
	}
}
