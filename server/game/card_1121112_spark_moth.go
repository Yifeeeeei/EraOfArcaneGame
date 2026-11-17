package game

import (
	"eraofarcane/model"
)

type Card1121112SparkMoth struct{ AlwaysActive }

func (Card1121112SparkMoth) ID() string { return "1121112" }

func (Card1121112SparkMoth) Name() string { return "火花飞蛾" }

func (e *Engine) triggerSparkMothAfterSpellHit(skill *CardInstance) {
	if skill == nil || skill.Card == nil || skill.Card.Category != model.ElementFire || !isSpellLikeCard(skill.Card) {
		return
	}
	for playerID := range e.State.Players {
		e.promptSparkMothReveal(playerID)
	}
}

func (e *Engine) promptSparkMothReveal(playerID int) {
	ps := e.State.Players[playerID]
	if ps == nil {
		return
	}
	candidates := e.friendlyHandCards(playerID, isSparkMoth)
	if len(candidates) == 0 {
		return
	}
	e.SetPendingAction(playerID, "spark_moth_reveal",
		"火花飞蛾:可以展示手牌中的此卡并-1入场花费", candidates, 0, len(candidates),
		func(selected []string) {
			for _, id := range selected {
				card, _ := ps.FindHandCard(id)
				if !isSparkMoth(card) {
					continue
				}
				if ps.RevealedHand == nil {
					ps.RevealedHand = make(map[string]bool)
				}
				ps.RevealedHand[card.InstanceID] = true
				card.Statuses["入场费用"+model.ElementFire+"-1"]++
			}
		})
}

func isSparkMoth(card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Card.Number == "1121112"
}
