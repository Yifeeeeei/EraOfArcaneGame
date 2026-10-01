package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card4211001Bartel struct{ AlwaysActive }

func (Card4211001Bartel) ID() string { return "4211001" }

func (Card4211001Bartel) Name() string { return "\"浪之人\" 巴特尔" }

func (Card4211001Bartel) OnUltimate(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "bartel_convert_hand_to_water",
		"巴特尔:展示1张手牌，其属性、入场花费和负载永久变为等量水", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			ps := ctx.Engine.State.Players[ctx.PlayerID]
			card, _ := ps.FindHandCard(selected[0])
			if card == nil || card.Card == nil {
				return
			}
			convertCardInstanceToWater(card)
			if ps.RevealedHand == nil {
				ps.RevealedHand = make(map[string]bool)
			}
			ps.RevealedHand[card.InstanceID] = true
			ctx.Engine.emit(GameEvent{
				Type:   "effect_trigger",
				Player: -1,
				Data: map[string]any{
					"source": cardToInfo(ctx.Source),
					"target": cardToInfo(card),
					"effect": "bartel_convert_hand_to_water",
				},
			})
		})
	return nil
}

func convertCardInstanceToWater(card *CardInstance) {
	if card == nil || card.Card == nil {
		return
	}
	cardCopy := *card.Card
	cardCopy.Category = model.ElementWater
	cardCopy.ElementsCost = convertElementMapToWater(card.Card.ElementsCost)
	cardCopy.ElementsGain = convertElementMapToWater(effectiveElementsGain(card))
	card.Card = &cardCopy
	card.ElementsGainSet = nil
	card.ElementsGainBonus = make(map[string]int)
}

func convertElementMapToWater(elements map[string]int) map[string]int {
	total := 0
	for _, amount := range elements {
		if amount > 0 {
			total += amount
		}
	}
	if total <= 0 {
		return map[string]int{}
	}
	return map[string]int{model.ElementWater: total}
}

func (Card4211001Bartel) ValidateAbility(ctx *EffectContext, trigger EffectTrigger) error {
	if trigger != TriggerUltimate {
		return nil
	}
	if len(ctx.Engine.State.Players[ctx.PlayerID].Hand) == 0 {
		return fmt.Errorf("Bartel ultimate requires a hand card")
	}
	return nil
}
