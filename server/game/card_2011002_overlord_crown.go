package game

type Card2011002OverlordCrown struct{ AlwaysActive }

func (Card2011002OverlordCrown) ID() string { return "2011002" }

func (Card2011002OverlordCrown) Name() string { return "统御者之冠" }

func (Card2011002OverlordCrown) OnUnitEnter(ctx *EffectContext) error {
	if ctx.ExtraData == nil || ctx.Target == nil || !ctx.Target.Card.IsCompanion() {
		return nil
	}
	if enteredPlayer, ok := ctx.ExtraData["entered_player"].(int); ok && enteredPlayer == ctx.PlayerID && ctx.Target != ctx.Source {
		setElementsGain(ctx.Target, map[string]int{})
	}
	return nil
}
