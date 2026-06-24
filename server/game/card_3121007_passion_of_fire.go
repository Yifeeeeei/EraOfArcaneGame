package game

import "eraofarcane/model"

type Card3121007PassionOfFire struct{ AlwaysActive }

func (Card3121007PassionOfFire) ID() string   { return "3121007" }
func (Card3121007PassionOfFire) Name() string { return "激情之火" }

func (Card3121007PassionOfFire) HasActiveSpellHit(card *CardInstance) bool {
	return abilityDurationActive(card)
}

func (Card3121007PassionOfFire) OnSpellHit(ctx *EffectContext) error {
	if ctx.Target == nil || ctx.Target.Card == nil {
		return nil
	}
	if !ctx.Target.Card.IsSkill() || ctx.Target.Card.Category != model.ElementFire {
		return nil
	}
	return DrawCards(1)(ctx)
}
