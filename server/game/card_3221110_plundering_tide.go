package game

type Card3221110PlunderingTide struct{ AlwaysActive }

func (Card3221110PlunderingTide) ID() string { return "3221110" }

func (Card3221110PlunderingTide) Name() string { return "劫掠之潮" }

func (Card3221110PlunderingTide) OnSpellHitBeforeDamage(ctx *EffectContext) error {
	if ctx.ExtraData == nil {
		return nil
	}
	affected, _ := ctx.ExtraData["affected_units"].([]*CardInstance)
	hitUnits := 0
	for _, unit := range affected {
		if unit != nil {
			hitUnits++
		}
	}
	if hitUnits == 0 {
		return nil
	}
	for i := 0; i < hitUnits; i++ {
		ctx.Engine.discardRandomHandCard(ctx.OpponentID)
	}
	ctx.Engine.drawCards(ctx.PlayerID, hitUnits)
	return nil
}
