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
	case RuleSummon, RuleEquip, RuleUseItem, RulePlaceTerrain, RuleAttack:
		return true
	case RuleUseSkill, RuleLearnSkill:
		return card != nil && isCombatSpellCard(card.Card)
	default:
		return false
	}
}

func (Card3411101TimeCycle) HasActiveActionRestriction(card *CardInstance) bool {
	return abilityDurationActive(card)
}
