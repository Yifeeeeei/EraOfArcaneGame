package game

type Card2321110PigeonRaidOrder struct{ AlwaysActive }

func (Card2321110PigeonRaidOrder) ID() string { return "2321110" }

func (Card2321110PigeonRaidOrder) Name() string { return "飞鸽急袭令" }

func (Card2321110PigeonRaidOrder) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlySkills(ctx.PlayerID, func(skill *CardInstance) bool {
		return isLearnedRushSkillThisTurn(ctx.Engine, skill)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "pigeon_raid_order_skill",
		"飞鸽急袭令:选择本回合学习的速攻法术", candidates, 1, 1,
		func(selected []string) {
			skill := findSkillSlotCard(ctx.Engine.State.Players[ctx.PlayerID], firstSelected(selected))
			if !isLearnedRushSkillThisTurn(ctx.Engine, skill) {
				return
			}
			ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
				Type:             TempModSkillPowerBonus,
				TargetInstanceID: skill.InstanceID,
				Amount:           1,
				RemainingUses:    1,
				ExpiresTurn:      ctx.Engine.State.TurnNumber + 2,
			})
			ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
				Type:             TempModNextSkillUseAttackBonus,
				TargetInstanceID: skill.InstanceID,
				Amount:           1,
				RemainingUses:    1,
				ExpiresTurn:      ctx.Engine.State.TurnNumber + 2,
			})
		})
	return nil
}

func isLearnedRushSkillThisTurn(e *Engine, skill *CardInstance) bool {
	return e != nil &&
		skill != nil &&
		skill.Card != nil &&
		skill.Card.IsSkill() &&
		skill.EnterTurn == e.State.TurnNumber &&
		cardHasRush(skill)
}
