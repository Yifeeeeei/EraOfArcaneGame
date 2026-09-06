package game

type Card1421108CelticDeer struct{ AlwaysActive }

func (Card1421108CelticDeer) ID() string { return "1421108" }

func (Card1421108CelticDeer) Name() string { return "凯尔特灵鹿" }

func (Card1421108CelticDeer) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Source == nil || !triggeredTurnAvailable(ctx.Source) || ctx.Target == nil || ctx.Target.Card == nil {
		return nil
	}
	if !ctx.Target.Card.IsSkill() || !hasCardTag(ctx.Target.Card, "灵媒") {
		return nil
	}
	if !useTriggeredTurn(ctx.Source) {
		return nil
	}
	resetCard(ctx.Source)
	return nil
}
