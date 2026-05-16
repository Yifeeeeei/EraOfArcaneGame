package game

type Card2221013DeepFrostCurseScroll struct{}

func (Card2221013DeepFrostCurseScroll) ID() string   { return "2221013" }
func (Card2221013DeepFrostCurseScroll) Name() string { return "深寒诅咒卷轴" }

func (Card2221013DeepFrostCurseScroll) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, false, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "deep_frost_curse",
		"选择1个敌方伙伴永久冻结", candidates, 1, 1,
		func(selected []string) {
			target := findEnemyByID(ctx, selected)
			if target == nil || !target.Card.IsCompanion() {
				return
			}
			target.Statuses[StatusFreeze] += 99
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source": cardToInfo(ctx.Source),
				"target": cardToInfo(target),
				"effect": "apply_status",
				"status": StatusFreeze,
				"amount": 99,
			}})
		})
	return nil
}
