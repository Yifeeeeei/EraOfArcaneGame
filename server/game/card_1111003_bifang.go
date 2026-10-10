package game

type Card1111003Bifang struct{ AlwaysActive }

func (Card1111003Bifang) ID() string { return "1111003" }

func (Card1111003Bifang) Name() string { return "毕方" }

func (Card1111003Bifang) DamageAdjustmentScope() DamageScope { return DamageEnemy }

func (Card1111003Bifang) AdjustDamage(_ *EffectContext, event DamageEvent, amount int) int {
	if event.Status == StatusBurn {
		return amount + 1
	}
	return amount
}
