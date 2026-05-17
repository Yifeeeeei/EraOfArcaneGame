package game

type Card3001002PureArcane struct{ AlwaysActive }

func (Card3001002PureArcane) ID() string   { return "3001002" }
func (Card3001002PureArcane) Name() string { return "纯净奥术" }
func (Card3001002PureArcane) NeedsSpellTarget() bool {
	return false
}
func (Card3001002PureArcane) OnSpellHit(ctx *EffectContext) error {
	choices := ctx.Engine.pureArcaneChoices(ctx.PlayerID)
	if len(choices) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "pure_arcane_spend",
		"选择同种元素与数量,使下一次该属性法术+等量威力", choices, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			elem, amount, ok := parsePureArcaneChoice(selected[0])
			if !ok {
				return
			}
			ps := ctx.Engine.State.Players[ctx.PlayerID]
			if ps.Elements[elem] < amount {
				return
			}
			ps.Elements[elem] -= amount
			ctx.Engine.addNextElementSpellPowerBonus(ctx.PlayerID, elem, amount)
		})
	return nil
}
