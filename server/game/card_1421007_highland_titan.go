package game

type Card1421007HighlandTitan struct{ AlwaysActive }

func (Card1421007HighlandTitan) ID() string { return "1421007" }

func (Card1421007HighlandTitan) Name() string { return "高地泰坦" }

func (Card1421007HighlandTitan) DamageAdjustmentScope() DamageScope { return DamageSelf }

func (Card1421007HighlandTitan) AdjustDamage(_ *EffectContext, event DamageEvent, amount int) int {
	if event.Kind == "spell" && event.BoostCount == 0 {
		return amount + 1
	}
	return amount
}
