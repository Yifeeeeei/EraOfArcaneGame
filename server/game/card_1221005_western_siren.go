package game

type Card1221005WesternSiren struct{ AlwaysActive }

func (Card1221005WesternSiren) ID() string            { return "1221005" }
func (Card1221005WesternSiren) Name() string          { return "西境海妖" }
func (Card1221005WesternSiren) IsPrayerAbility() bool { return true }

func (Card1221005WesternSiren) OnPerTurn(ctx *EffectContext) error {
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card.Card.IsCompanion() && card.Position != nil && ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "western_siren_consume",
		"西境海妖:选择法力范围内1个敌方伙伴横置", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			target := ctx.Engine.findFieldCardByInstance(ctx.Engine.State.Players[ctx.OpponentID], selected[0])
			if target != nil {
				target.IsHorizontal = true
			}
		})
	return nil
}
