package game

type Card1621003FearDemon struct{}

func (Card1621003FearDemon) ID() string   { return "1621003" }
func (Card1621003FearDemon) Name() string { return "恐惧魔" }

func (Card1621003FearDemon) OnPerTurn(ctx *EffectContext) error {
	ctx.Source.CurrentLife -= 3
	if ctx.Source.CurrentLife <= 0 {
		ctx.Engine.destroyUnit(ctx.Source, ctx.PlayerID)
	}
	return nil
}
