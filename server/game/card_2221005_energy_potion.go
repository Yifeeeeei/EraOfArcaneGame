package game

type Card2221005EnergyPotion struct{ AlwaysActive }

func (Card2221005EnergyPotion) ID() string   { return "2221005" }
func (Card2221005EnergyPotion) Name() string { return "精力药剂" }

func (Card2221005EnergyPotion) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModResetSkillsOnOpponentTurnEnd,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		Amount:           1,
	})
	return nil
}
