package game

type Card1021006Grocer struct{ AlwaysActive }

func (Card1021006Grocer) ID() string   { return "1021006" }
func (Card1021006Grocer) Name() string { return "杂货商贩" }
func (Card1021006Grocer) OnEnter(ctx *EffectContext) error {
	return DrawCards(2)(ctx)
}
