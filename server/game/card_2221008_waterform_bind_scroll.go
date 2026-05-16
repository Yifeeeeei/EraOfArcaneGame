package game

type Card2221008WaterformBindScroll struct{}

func (Card2221008WaterformBindScroll) ID() string   { return "2221008" }
func (Card2221008WaterformBindScroll) Name() string { return "水形之束卷轴" }

func (Card2221008WaterformBindScroll) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, false, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "waterform_bind_consume",
		"选择1个敌方伙伴并消耗它", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			target := ctx.Engine.findCardOnField(ctx.Engine.State.Players[ctx.OpponentID], selected[0])
			if target == nil || !target.Card.IsCompanion() {
				return
			}
			target.IsHorizontal = true
			ctx.Engine.dealDamage(target, 1, ctx.OpponentID)
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source": cardToInfo(ctx.Source),
				"target": cardToInfo(target),
				"effect": "consume_target",
			}})
		})
	return nil
}
