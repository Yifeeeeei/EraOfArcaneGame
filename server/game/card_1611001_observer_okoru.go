package game

type Card1611001ObserverOkoru struct{ AlwaysActive }

func (Card1611001ObserverOkoru) ID() string   { return "1611001" }
func (Card1611001ObserverOkoru) Name() string { return "\"观察者\" 欧柯茹" }
func (Card1611001ObserverOkoru) OnEnter(ctx *EffectContext) error {
	if err := DrawCards(1)(ctx); err != nil {
		return err
	}
	return DealDamageToSelfHero(1)(ctx)
}
