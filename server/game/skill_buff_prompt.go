package game

func promptPermanentSkillBuff(ctx *EffectContext, prompt string) {
	targets := ctx.Engine.friendlySkills(ctx.PlayerID, func(skill *CardInstance) bool {
		return ctx.Source == nil || skill.InstanceID != ctx.Source.InstanceID
	})
	if len(targets) == 0 {
		return
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "permanent_skill_buff_target",
		prompt, targets, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			targetID := selected[0]
			sourceNumber := ""
			if ctx.Source != nil && ctx.Source.Card != nil {
				sourceNumber = ctx.Source.Card.Number
			}
			choices := []map[string]any{
				{"instance_id": "power", "number": sourceNumber, "name": "+3威", "type": "选择", "zone": "choice", "side": "own"},
				{"instance_id": "attack", "number": sourceNumber, "name": "+1攻", "type": "选择", "zone": "choice", "side": "own"},
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "permanent_skill_buff_mode",
				"选择永久强化", choices, 1, 1,
				func(selected []string) {
					if len(selected) == 0 {
						return
					}
					for _, skill := range ctx.Engine.State.Players[ctx.PlayerID].Skills {
						if skill == nil || skill.InstanceID != targetID {
							continue
						}
						if selected[0] == "attack" {
							skill.AttackBonus++
						} else {
							skill.PowerBonus += 3
						}
						return
					}
				})
		})
}
