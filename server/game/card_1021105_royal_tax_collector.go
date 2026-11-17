package game

type Card1021105RoyalTaxCollector struct{ AlwaysActive }

func (Card1021105RoyalTaxCollector) ID() string   { return "1021105" }
func (Card1021105RoyalTaxCollector) Name() string { return "皇城征税员" }
func (Card1021105RoyalTaxCollector) OnEnter(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{Type: TempModRoyalTax, SourceCardNumber: "1021105", SourceName: "皇城征税员", Amount: 1, ExpiresAtTurnEnd: true, ExpiresOnPlayerID: ctx.OpponentID})
	return nil
}
