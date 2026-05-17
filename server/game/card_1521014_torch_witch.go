package game

type Card1521014TorchWitch struct{ AlwaysActive }

func (Card1521014TorchWitch) ID() string   { return "1521014" }
func (Card1521014TorchWitch) Name() string { return "炬之女巫" }
func (Card1521014TorchWitch) OnEnter(ctx *EffectContext) error {
	return ApplyStatusToSelf(StatusBurn, 2)(ctx)
}
