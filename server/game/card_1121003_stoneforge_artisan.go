package game

type Card1121003StoneforgeArtisan struct{ AlwaysActive }

func (Card1121003StoneforgeArtisan) ID() string   { return "1121003" }
func (Card1121003StoneforgeArtisan) Name() string { return "锻石工匠" }

func (Card1121003StoneforgeArtisan) OnConsume(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlySkills(ctx.PlayerID, func(card *CardInstance) bool {
		return card.Card.IsSkill() && canUseSkillForPurpose(card.Card, skillPurposeAttack)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "spell_power_bonus",
		"选择你的1个法术，本回合获得+2威",
		candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
				Type:             TempModSkillPowerBonus,
				SourceCardNumber: ctx.Source.Card.Number,
				SourceName:       ctx.Source.Card.Name,
				TargetInstanceID: selected[0],
				Amount:           2,
				ExpiresTurn:      ctx.Engine.State.TurnNumber,
			})
		})
	return nil
}
