package game

type Card1121002LivelyHearth struct{ AlwaysActive }

func (Card1121002LivelyHearth) ID() string   { return "1121002" }
func (Card1121002LivelyHearth) Name() string { return "活泼的炉火" }
func (Card1121002LivelyHearth) OnEnter(ctx *EffectContext) error {
	return DrawCards(1)(ctx)
}
