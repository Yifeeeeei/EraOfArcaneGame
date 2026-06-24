package game

type Card2221013DeepFrostCurseScroll struct{ AlwaysActive }

func (Card2221013DeepFrostCurseScroll) ID() string   { return "2221013" }
func (Card2221013DeepFrostCurseScroll) Name() string { return "深寒诅咒卷轴" }

func (Card2221013DeepFrostCurseScroll) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card.Card.IsCompanion() && card.Position != nil && ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "deep_frost_curse",
		"选择1个敌方伙伴永久冻结", candidates, 1, 1,
		func(selected []string) {
			target := selectedUnitFromCandidates(ctx.Engine, selected, candidates)
			if target == nil {
				return
			}
			if !ctx.Engine.addStatus(target, StatusFreeze, 99) {
				return
			}
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
