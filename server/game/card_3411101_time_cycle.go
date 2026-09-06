package game

import "eraofarcane/model"

type Card3411101TimeCycle struct{ AlwaysActive }

func (Card3411101TimeCycle) ID() string   { return "3411101" }
func (Card3411101TimeCycle) Name() string { return "时岁轮转" }
func (Card3411101TimeCycle) PaymentConstraint(_ *CardInstance, purpose paymentPurpose, cost map[string]int) PaymentConstraint {
	if purpose != paymentPurposeLearn && purpose != paymentPurposeUse {
		return PaymentConstraint{}
	}
	earth := cost[model.ElementEarth]
	arcane := totalElementCost(cost) - earth
	required := map[string]int{}
	if earth > 0 {
		required[model.ElementEarth] = earth
	}
	if arcane > 0 {
		required[model.ElementArcane] = arcane
	}
	return PaymentConstraint{StrictElements: required}
}

func (Card3411101TimeCycle) RestrictsAction(_ *CardInstance, action RuleAction, card *CardInstance) bool {
	switch action {
	case RuleSummon, RuleEquip, RuleUseItem, RulePlaceTerrain, RuleLearnSkill, RuleAttack:
		return true
	case RuleUseSkill:
		// The ongoing sorcery must remain usable to spend its remaining duration.
		return card == nil || card.Card == nil || card.Card.Number != "3411101"
	default:
		return false
	}
}
