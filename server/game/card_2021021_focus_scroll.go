package game

import "eraofarcane/model"

type Card2021021FocusScroll struct{}

func (Card2021021FocusScroll) ID() string   { return "2021021" }
func (Card2021021FocusScroll) Name() string { return "聚能卷轴" }

func (Card2021021FocusScroll) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModDelayedElementGain,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		Status:           model.ElementArcane,
		Amount:           3,
	})
	return nil
}
