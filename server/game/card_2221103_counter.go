package game

func (Card2221103IceLockRune) CounterTriggers() []EffectTrigger {
	return []EffectTrigger{TriggerOnCardEnter}
}

func (Card2221103IceLockRune) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Trigger == TriggerOnCardEnter && ctx.Event.PlayerID != ctx.Source.OwnerID && ctx.Event.LearnedSkill && ctx.Event.Card != nil && ctx.Event.Card.Card.IsSkill()
}

type Card2221103IceLockRune struct{ AlwaysActive }

func (Card2221103IceLockRune) ID() string { return "2221103" }

func (Card2221103IceLockRune) Name() string { return "冰锁符文" }

func (Card2221103IceLockRune) OnCardEnter(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Target == nil || ctx.Target.Card == nil || ctx.ExtraData == nil {
		return nil
	}
	learned, _ := ctx.ExtraData["learned_skill"].(bool)
	enteredPlayer, _ := ctx.ExtraData["entered_player"].(int)
	if !learned || enteredPlayer == ctx.PlayerID || !ctx.Target.Card.IsSkill() {
		return nil
	}
	ctx.Target.Statuses[StatusCannotUseSkillUntilTurn] = ctx.Engine.State.TurnNumber + 1
	ctx.Engine.emit(GameEvent{
		Type:   "skill_locked",
		Player: -1,
		Data: map[string]any{
			"player": ctx.Target.OwnerID,
			"source": cardToInfo(ctx.Source),
			"target": cardToInfo(ctx.Target),
			"until":  ctx.Target.Statuses[StatusCannotUseSkillUntilTurn],
		},
	})
	return nil
}
