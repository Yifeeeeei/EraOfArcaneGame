package game

import (
	"eraofarcane/model"
)

type Card4111003Brahma struct{ AlwaysActive }

func (Card4111003Brahma) ID() string { return "4111003" }

func (Card4111003Brahma) Name() string { return "大祭司 梵天" }

func (Card4111003Brahma) OnUltimate(ctx *EffectContext) error {
	ctx.Source.Statuses["梵天火焰命中"] = 1
	return nil
}

func (Card4111003Brahma) OnSpellHit(ctx *EffectContext) error {
	if ctx.Source.Statuses["梵天火焰命中"] <= 0 || ctx.Target == nil || ctx.Target.Card == nil {
		return nil
	}
	if !isFriendlySpellHit(ctx) || ctx.Target.Card.Category != model.ElementFire {
		return nil
	}
	ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementFire, 1, ctx.Source)
	return nil
}

func (Card4111003Brahma) OnTurnEnd(ctx *EffectContext) error {
	delete(ctx.Source.Statuses, "梵天火焰命中")
	return nil
}
