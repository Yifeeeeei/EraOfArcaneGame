package game

type Card3621108Moonshadow struct{ AlwaysActive }

func (Card3621108Moonshadow) ID() string { return "3621108" }

func (Card3621108Moonshadow) Name() string { return "月影" }

func (Card3621108Moonshadow) ObservesSpellBattles() bool { return true }

func (Card3621108Moonshadow) OnDefend(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil {
		return nil
	}
	if !triggeredTurnAvailable(ctx.Source) {
		return nil
	}
	for _, skill := range opposingBattleSpells(ctx) {
		if skill != nil && skill.OwnerID != ctx.PlayerID && skill.Statuses[StatusWeaken] > 0 && ctx.Engine.hasEffectiveStatus(skill, StatusWeaken) {
			source := ctx.Source
			reason := skill
			ctx.Engine.SetTriggeredTurnAction(source, ctx.PlayerID, "moonshadow_reset",
				"月影:是否重置此卡", []map[string]any{candidateInfo(source, "skill", "own")}, 0, 1,
				func(selected []string) {
					if len(selected) == 0 || source == nil || source.Card == nil {
						return
					}
					if !useTriggeredTurn(source) {
						return
					}
					ctx.Engine.resetCard(source)
					ctx.Engine.emit(GameEvent{
						Type:   "moonshadow_reset",
						Player: -1,
						Data: map[string]any{
							"player": ctx.PlayerID,
							"source": cardToInfo(source),
							"reason": cardToInfo(reason),
						},
					})
				})
			return nil
		}
	}
	return nil
}

var _ OnDefendBehavior = Card3621108Moonshadow{}
