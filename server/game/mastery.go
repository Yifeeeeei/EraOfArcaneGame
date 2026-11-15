package game

// Thresholds describe printed mastery milestones; Max alone cannot express
// sparse progressions such as 1,3,5 or 3,6.
type MasteryThresholdBehavior interface {
	MasteryThresholds() []int
}

func nextMasteryThreshold(card *CardInstance) int {
	b, ok := masteryBehavior(card)
	if !ok {
		return 0
	}
	thresholds := []int{b.MasteryMax()}
	if progression, ok := b.(MasteryThresholdBehavior); ok {
		thresholds = progression.MasteryThresholds()
	}
	next := 0
	for _, level := range thresholds {
		if level > card.Statuses[StatusMastery] && (next == 0 || level < next) {
			next = level
		}
	}
	return next
}

func (e *Engine) advanceToNextMastery(card *CardInstance, playerID int) {
	if next := nextMasteryThreshold(card); next > 0 {
		e.advanceMastery(card, playerID, next-card.Statuses[StatusMastery])
	}
}

func masteryBehavior(card *CardInstance) (MasteryBehavior, bool) {
	if card == nil || card.Card == nil {
		return nil, false
	}
	behavior, ok := globalRegistry.GetBehavior(card.Card.Number).(MasteryBehavior)
	if !ok || !behavior.HasActiveMastery(card) {
		return nil, false
	}
	return behavior, true
}

func (e *Engine) advanceMastery(card *CardInstance, playerID int, amount int) {
	card = physicalSpellSource(card)
	if amount <= 0 {
		return
	}
	behavior, ok := masteryBehavior(card)
	if !ok {
		return
	}
	maxLevel := behavior.MasteryMax()
	if maxLevel <= 0 {
		return
	}
	oldLevel := card.Statuses[StatusMastery]
	newLevel := oldLevel + amount
	if newLevel > maxLevel {
		newLevel = maxLevel
	}
	if newLevel <= oldLevel {
		return
	}
	card.Statuses[StatusMastery] = newLevel
	for level := oldLevel + 1; level <= newLevel; level++ {
		ctx := &EffectContext{
			Engine:     e,
			Source:     card,
			PlayerID:   playerID,
			OpponentID: 1 - playerID,
			ExtraData:  map[string]any{"mastery": level},
		}
		_ = behavior.OnMastery(ctx, level)
		e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
			"source":  cardToInfo(card),
			"effect":  "mastery",
			"mastery": level,
		}})
		e.triggerFieldEffectsWithData(TriggerOnMastery, playerID, card, map[string]any{
			"mastery":         level,
			"mastered_card":   card,
			"mastered_player": playerID,
		})
	}
}

func (e *Engine) advanceMasteryForUsedSkills(playerID int, skills ...*CardInstance) {
	for _, skill := range skills {
		if skill == nil || skill.Card == nil || !skill.Card.IsSkill() {
			continue
		}
		e.advanceMastery(skill, playerID, 1)
	}
}

func (e *Engine) advanceAllMasteryToMax(ps *PlayerState) {
	if ps == nil {
		return
	}
	for _, card := range e.getAllFieldCards(ps) {
		behavior, ok := masteryBehavior(card)
		if !ok {
			continue
		}
		e.advanceMastery(card, ps.PlayerID, behavior.MasteryMax()-card.Statuses[StatusMastery])
	}
}

func (e *Engine) settleMastery(ps *PlayerState) {
	// 精通 does not advance during mark settlement. It advances when the card
	// itself is consumed, then threshold effects resolve immediately.
}
