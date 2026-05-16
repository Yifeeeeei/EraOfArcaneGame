package game

func findEnemyByID(ctx *EffectContext, selected []string) *CardInstance {
	if ctx == nil || len(selected) == 0 {
		return nil
	}
	return ctx.Engine.findCardOnField(ctx.Engine.State.Players[ctx.OpponentID], selected[0])
}
