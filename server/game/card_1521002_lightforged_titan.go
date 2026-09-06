package game

import (
	"strings"
)

type Card1521002LightforgedTitan struct{ AlwaysActive }

func (Card1521002LightforgedTitan) ID() string { return "1521002" }

func (Card1521002LightforgedTitan) Name() string { return "光铸泰坦" }

func (Card1521002LightforgedTitan) OnEnter(ctx *EffectContext) error {
	return DrawCards(2)(ctx)
}

// AdjustDamage is a pure calculation. Prevention, receipts and damage observers
// all see the resulting amount; this aura never performs a second life mutation.
func (Card1521002LightforgedTitan) DamageAdjustmentScope() DamageScope { return DamageSelf }

func (Card1521002LightforgedTitan) AdjustDamage(_ *EffectContext, event DamageEvent, amount int) int {
	if event.Kind != "spell" || event.Spell == nil {
		return amount
	}
	tag := event.Spell.Tag
	if strings.Contains(tag, "驱动") || strings.Contains(tag, "神秘") || strings.Contains(tag, "聚能") {
		return amount + 1
	}
	return amount
}
