package game

type Card1521002LightforgedTitan struct{}

func (Card1521002LightforgedTitan) ID() string   { return "1521002" }
func (Card1521002LightforgedTitan) Name() string { return "光铸泰坦" }
func (Card1521002LightforgedTitan) OnEnter(ctx *EffectContext) error {
	return DrawCards(2)(ctx)
}
