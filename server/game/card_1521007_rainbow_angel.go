package game

type Card1521007RainbowAngel struct{ AlwaysActive }

func (Card1521007RainbowAngel) ID() string { return "1521007" }

func (Card1521007RainbowAngel) Name() string { return "虹之天使" }

func (Card1521007RainbowAngel) PaymentRules(*CardInstance) PaymentRules {
	return PaymentRules{LightPaysAny: true}
}
