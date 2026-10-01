package game

type Card1521015EmberWitch struct{ AlwaysActive }

func (Card1521015EmberWitch) ID() string { return "1521015" }

func (Card1521015EmberWitch) Name() string { return "烬之女巫" }

func (Card1521015EmberWitch) OnEnter(ctx *EffectContext) error {
	return ApplyStatusToSelf(StatusBurn, 3)(ctx)
}
