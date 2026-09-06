package game

type Card1621013WordSpirit struct{ AlwaysActive }

func (Card1621013WordSpirit) ID() string { return "1621013" }

func (Card1621013WordSpirit) Name() string { return "言灵" }

func (Card1621013WordSpirit) OnSpellCast(ctx *EffectContext) error {
	if !isEnemySpellCast(ctx) || !triggeredTurnAvailable(ctx.Source) {
		return nil
	}
	targets := make([]*CardInstance, 0)
	for _, skill := range friendlySpellInstancesIncludingBound(ctx.Engine, ctx.OpponentID) {
		if skill != nil && skill.IsHorizontal && canInstanceBeWeakened(skill) {
			targets = append(targets, skill)
		}
	}
	if len(targets) == 0 {
		return nil
	}
	ctx.Engine.SetTriggeredTurnAction(ctx.Source, ctx.PlayerID, "word_spirit_weaken",
		"言灵:是否使敌方所有横置法术虚弱1", []map[string]any{candidateInfo(ctx.Source, "unit", "own")}, 0, 1,
		func(selected []string) {
			accepted := len(selected) > 0 && selected[0] == ctx.Source.InstanceID
			if !accepted || !useTriggeredTurn(ctx.Source) {
				return
			}
			for _, skill := range targets {
				if skill != nil && skill.IsHorizontal && canInstanceBeWeakened(skill) {
					ctx.Engine.addStatus(skill, StatusWeaken, 1)
				}
			}
		})
	return nil
}
