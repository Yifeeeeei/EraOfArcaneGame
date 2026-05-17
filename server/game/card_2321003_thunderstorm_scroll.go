package game

type Card2321003ThunderstormScroll struct{ AlwaysActive }

func (Card2321003ThunderstormScroll) ID() string   { return "2321003" }
func (Card2321003ThunderstormScroll) Name() string { return "雷暴卷轴" }

func (Card2321003ThunderstormScroll) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, true, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "thunderstorm_scroll",
		"选择1个敌方单位,对敌方方阵造成1伤并使命中的伙伴眩晕1", candidates, 1, 1,
		func(selected []string) {
			target := findEnemyByID(ctx, selected)
			if target == nil {
				return
			}
			for _, unit := range ctx.Engine.spellAffectedUnits(ctx.OpponentID, ctx.Source, SpellTarget{Type: "unit", Position: *target.Position}) {
				ctx.Engine.dealDamage(unit, 1, ctx.OpponentID)
				if unit.Card.IsCompanion() {
					unit.Statuses[StatusStun]++
				}
			}
		})
	return nil
}
