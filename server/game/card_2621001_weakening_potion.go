package game

type Card2621001WeakeningPotion struct{ AlwaysActive }

func (Card2621001WeakeningPotion) ID() string   { return "2621001" }
func (Card2621001WeakeningPotion) Name() string { return "虚弱药剂" }

func (Card2621001WeakeningPotion) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.enemySkills(ctx.PlayerID, canInstanceBeWeakened)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "weakening_potion",
		"选择最多2个敌方法术虚弱2", candidates, 1, 2,
		func(selected []string) {
			for _, id := range selected {
				for _, skill := range ctx.Engine.State.Players[ctx.OpponentID].Skills {
					if skill != nil && skill.InstanceID == id {
						ctx.Engine.addStatus(skill, StatusWeaken, 2)
					}
				}
			}
		})
	return nil
}
