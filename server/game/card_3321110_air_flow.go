package game

import (
	"eraofarcane/model"
)

type Card3321110AirFlow struct{ AlwaysActive }

func (Card3321110AirFlow) ID() string { return "3321110" }

func (Card3321110AirFlow) Name() string { return "气蕴成流" }

func (Card3321110AirFlow) OnEnter(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModNextLearnedSkillHaste,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		Element:          model.ElementAir,
		RemainingUses:    1,
	})
	return nil
}
