package game

type Card1421001SandMage struct{ AlwaysActive }

func (Card1421001SandMage) ID() string { return "1421001" }

func (Card1421001SandMage) Name() string { return "流沙法师" }

func (Card1421001SandMage) OnEnter(ctx *EffectContext) error {
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, true, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "sand_mage_petrify",
		"流沙法师:选择1个敌人石化1", candidates, 1, 1,
		func(selected []string) {
			target := selectedUnitFromCandidates(ctx.Engine, selected, candidates)
			if target != nil {
				ctx.Engine.addStatus(target, StatusPetrify, 1)
			}
		})
	return nil
}
