package game

func (Card2021115InfusionRuneE) CounterTriggers() []EffectTrigger {
	return []EffectTrigger{TriggerOnDefend}
}

func (Card2021115InfusionRuneE) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Trigger == TriggerOnDefend && ctx.Source.OwnerID == ctx.Event.DefenderID &&
		len(ctx.Event.DefenseSkills)+len(ctx.Event.DefenseBoosts) > 0
}

type Card2021115InfusionRuneE struct{ AlwaysActive }

func (Card2021115InfusionRuneE) ID() string { return "2021115" }

func (Card2021115InfusionRuneE) Name() string { return "注能符文E型" }

func (Card2021115InfusionRuneE) OnDefend(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.ExtraData == nil {
		return nil
	}
	reset := 0
	for _, skill := range append(spellInstancesFromData(ctx.ExtraData, "defense_skills"), spellInstancesFromData(ctx.ExtraData, "defense_boosts")...) {
		if skill == nil || skill.Card == nil || !skill.Card.IsSkill() || skill.OwnerID != ctx.PlayerID {
			continue
		}
		skill.IsHorizontal = false
		reset++
	}
	if reset > 0 {
		ctx.Engine.emit(GameEvent{
			Type:   "infusion_rune_e_reset",
			Player: -1,
			Data: map[string]any{
				"player": ctx.PlayerID,
				"source": cardToInfo(ctx.Source),
				"count":  reset,
			},
		})
	}
	return nil
}

var _ OnDefendBehavior = Card2021115InfusionRuneE{}
