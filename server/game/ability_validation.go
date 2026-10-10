package game

// AbilityValidationBehavior checks card-specific requirements without changing
// state. The engine separately checks timing, active state and use limits.
type AbilityValidationBehavior interface {
	ValidateAbility(*EffectContext, EffectTrigger) error
}

func (e *Engine) validateAbility(card *CardInstance, trigger EffectTrigger) error {
	behavior, ok := cardBehavior(card).(AbilityValidationBehavior)
	if !ok {
		return nil
	}
	return behavior.ValidateAbility(&EffectContext{Engine: e, Source: card, PlayerID: card.OwnerID, OpponentID: 1 - card.OwnerID}, trigger)
}

// ItemUseValidationBehavior validates targets and additional costs before an
// item leaves the hand or its entry cost is paid.
type ItemUseValidationBehavior interface{ ValidateItemUse(*EffectContext) error }
