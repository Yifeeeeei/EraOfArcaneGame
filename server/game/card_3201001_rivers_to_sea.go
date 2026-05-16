package game

import "eraofarcane/model"

type Card3201001RiversToSea struct{}

func (Card3201001RiversToSea) ID() string   { return "3201001" }
func (Card3201001RiversToSea) Name() string { return "百川归海" }

func (Card3201001RiversToSea) OnDefend(ctx *EffectContext) error {
	success, _ := ctx.ExtraData["defense_success"].(bool)
	if !success {
		return nil
	}
	total := 0
	if skill, _ := ctx.ExtraData["attack_skill"].(*CardInstance); skill != nil && skill.Card != nil {
		total += max(skill.Card.Attack+skill.AttackBonus, 0)
	}
	if boosts, _ := ctx.ExtraData["boost_skills"].([]*CardInstance); len(boosts) > 0 {
		for _, skill := range boosts {
			if skill != nil && skill.Card != nil {
				total += max(skill.Card.Attack+skill.AttackBonus, 0)
			}
		}
	}
	if total <= 0 {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	ps.GainElements(map[string]int{model.ElementWater: total})
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source":  cardToInfo(ctx.Source),
		"effect":  "gain_element",
		"element": model.ElementWater,
		"amount":  total,
	}})
	return nil
}
