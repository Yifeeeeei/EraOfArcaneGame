package game

type Card2421011ElfArmor struct{ AlwaysActive }

func (Card2421011ElfArmor) ID() string   { return "2421011" }
func (Card2421011ElfArmor) Name() string { return "精灵铠" }
func (Card2421011ElfArmor) OnPerTurn(ctx *EffectContext) error {
	healUnit(ctx.Engine.State.Players[ctx.PlayerID].Hero, 1)
	return nil
}
