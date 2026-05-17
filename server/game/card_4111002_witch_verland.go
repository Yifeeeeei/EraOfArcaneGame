package game

import "eraofarcane/model"

type Card4111002WitchVerland struct{}

func (Card4111002WitchVerland) ID() string   { return "4111002" }
func (Card4111002WitchVerland) Name() string { return "女巫 维兰德" }
func (Card4111002WitchVerland) OnPerTurn(ctx *EffectContext) error {
	ctx.Source.Statuses[StatusBurn]++
	if effectiveElementsGain(ctx.Source)[model.ElementFire] > 0 {
		addElementsGainBonus(ctx.Source, model.ElementFire, -1)
		addElementsGainBonus(ctx.Source, model.ElementArcane, 1)
		ctx.Source.Statuses["维兰德负载转换"]++
	}
	return nil
}

func (Card4111002WitchVerland) OnTurnEnd(ctx *EffectContext) error {
	count := ctx.Source.Statuses["维兰德负载转换"]
	if count <= 0 {
		return nil
	}
	addElementsGainBonus(ctx.Source, model.ElementFire, count)
	addElementsGainBonus(ctx.Source, model.ElementArcane, -count)
	delete(ctx.Source.Statuses, "维兰德负载转换")
	return nil
}
