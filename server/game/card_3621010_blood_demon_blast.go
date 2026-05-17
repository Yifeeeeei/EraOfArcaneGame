package game

type Card3621010BloodDemonBlast struct{ AlwaysActive }

func (Card3621010BloodDemonBlast) ID() string   { return "3621010" }
func (Card3621010BloodDemonBlast) Name() string { return "血魔爆" }

func (Card3621010BloodDemonBlast) OnSpellCast(ctx *EffectContext) error {
	if !isSpellBeingCast(ctx) {
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card.Card.IsCompanion()
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "blood_demon_blast",
		"献祭1个伙伴，对敌方前排造成其生命值的伤害", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			player := ctx.Engine.State.Players[ctx.PlayerID]
			unit := ctx.Engine.findFieldCardByInstance(player, selected[0])
			if unit == nil || unit.Position == nil {
				return
			}
			damage := max(unit.CurrentLife, 0)
			ctx.Engine.destroyUnit(unit, ctx.PlayerID)
			row := ctx.Engine.State.Players[ctx.OpponentID].GetFrontRow()
			if row < 0 || damage <= 0 {
				return
			}
			for col := 0; col < 3; col++ {
				target := ctx.Engine.State.Players[ctx.OpponentID].Units[col][row]
				if target != nil {
					ctx.Engine.dealDamage(target, damage, ctx.OpponentID)
				}
			}
		})
	return nil
}
