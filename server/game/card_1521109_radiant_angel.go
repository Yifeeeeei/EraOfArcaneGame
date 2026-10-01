package game

type Card1521109RadiantAngel struct{ AlwaysActive }

func (Card1521109RadiantAngel) ID() string   { return "1521109" }
func (Card1521109RadiantAngel) Name() string { return "辉之天使" }
func (Card1521109RadiantAngel) PaymentRules(*CardInstance) PaymentRules {
	return PaymentRules{AnyPaysLight: true}
}
