package game

import (
	"eraofarcane/model"
)

type Card1121013Arsonist struct{ AlwaysActive }

func (Card1121013Arsonist) ID() string { return "1121013" }

func (Card1121013Arsonist) Name() string { return "纵火者" }

func (Card1121013Arsonist) OnSpellCast(ctx *EffectContext) error {
	if !triggeredTurnAvailable(ctx.Source) || !isFriendlySpellCast(ctx) || spellCastSourceElement(ctx) != model.ElementFire || ctx.Target == nil || isSorcerySkill(ctx.Target.Card) {
		return nil
	}
	candidates := append(ctx.Engine.friendlyUnits(ctx.PlayerID, true, nil), ctx.Engine.enemyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card.Position != nil && ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, cardHasPierce(ctx.Target))
	})...)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetTriggeredTurnAction(ctx.Source, ctx.PlayerID, "arsonist_burn",
		"纵火者:可以选择法力范围内1个单位点燃1", candidates, 0, 1,
		func(selected []string) {
			target := selectedUnitFromCandidates(ctx.Engine, selected, candidates)
			if target == nil || !useTriggeredTurn(ctx.Source) {
				return
			}
			ctx.Engine.addStatus(target, StatusBurn, 1)
		})
	return nil
}
