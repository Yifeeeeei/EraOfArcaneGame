package game

import (
	"eraofarcane/model"
)

type Card3011101AbsolutePurityArcaneOneness struct{ AlwaysActive }

func (Card3011101AbsolutePurityArcaneOneness) ID() string { return "3011101" }

func (Card3011101AbsolutePurityArcaneOneness) Name() string { return "绝对纯净 奥能一心" }

func (Card3011101AbsolutePurityArcaneOneness) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.Source.Card.Number != "3011101" {
		return
	}
	stats.PowerBonus += countTopConsecutiveArcaneCards(ctx.Engine.State.Players[ctx.PlayerID])
}

func (Card3011101AbsolutePurityArcaneOneness) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.Source.Card.Number != "3011101" || !isFriendlySpellCast(ctx) {
		return nil
	}
	count := countTopConsecutiveArcaneCards(ctx.Engine.State.Players[ctx.PlayerID])
	if count > 0 || len(ctx.Engine.State.Players[ctx.PlayerID].Deck) > 0 {
		ctx.Engine.shuffleDeck(ctx.PlayerID)
	}
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source":       cardToInfo(ctx.Source),
		"effect":       "arcane_oneness_reveal",
		"arcane_count": count,
	}})
	return nil
}

func (Card3011101AbsolutePurityArcaneOneness) PaymentConstraint(_ *CardInstance, purpose paymentPurpose, cost map[string]int) PaymentConstraint {
	if purpose == paymentPurposeLearn || purpose == paymentPurposeUse {
		return PaymentConstraint{StrictElements: map[string]int{model.ElementArcane: totalElementCost(cost)}}
	}
	return PaymentConstraint{}
}

func countTopConsecutiveArcaneCards(ps *PlayerState) int {
	if ps == nil {
		return 0
	}
	count := 0
	for _, card := range ps.Deck {
		if card == nil || card.Card == nil || card.Card.Category != model.ElementArcane {
			break
		}
		count++
	}
	return count
}
