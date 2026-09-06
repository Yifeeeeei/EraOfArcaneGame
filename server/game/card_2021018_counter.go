package game

func (Card2021018ArcaneRune) CounterTriggers() []EffectTrigger {
	return []EffectTrigger{TriggerOnSpellCast}
}

func (Card2021018ArcaneRune) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Trigger == TriggerOnSpellCast && ctx.Event.PlayerID != ctx.Source.OwnerID && len(ctx.Engine.friendlySkillsIncludingBound(ctx.Source.OwnerID, nil)) > 0
}

type Card2021018ArcaneRune struct{ AlwaysActive }

func (Card2021018ArcaneRune) ID() string { return "2021018" }

func (Card2021018ArcaneRune) Name() string { return "奥术符文" }

func (Card2021018ArcaneRune) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlySkillsIncludingBound(ctx.PlayerID, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "arcane_rune_skill", "奥术符文:选择己方1个法术获得+3威", candidates, 1, 1, func(selected []string) {
		skill := ctx.Engine.findSkill(ctx.Engine.State.Players[ctx.PlayerID], firstSelected(selected))
		if skill != nil {
			ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
				Type:              TempModSkillPowerBonus,
				SourceCardNumber:  ctx.Source.Card.Number,
				SourceName:        ctx.Source.Card.Name,
				TargetInstanceID:  skill.InstanceID,
				Amount:            3,
				ExpiresAtTurnEnd:  true,
				ExpiresOnPlayerID: ctx.Engine.State.CurrentTurn,
			})
			ctx.Engine.refreshPendingSpellPowerForModifiedSkill(ctx.PlayerID, skill)
		}
	})
	return nil
}
