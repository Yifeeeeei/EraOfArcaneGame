package game

type Card1121014Firethorn struct{ AlwaysActive }

func (Card1121014Firethorn) ID() string   { return "1121014" }
func (Card1121014Firethorn) Name() string { return "火荆" }
func (Card1121014Firethorn) OnDeath(ctx *EffectContext) error {
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, true, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "firethorn_death_burn",
		"火荆遗言:选择1个敌方单位施加点燃1", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			target := ctx.Engine.findCardOnField(ctx.Engine.State.Players[ctx.OpponentID], selected[0])
			if target != nil {
				target.Statuses[StatusBurn]++
			}
		})
	return nil
}
