package game

import (
	"eraofarcane/model"
)

type Card1401001LifeSeed struct{ AlwaysActive }

func (Card1401001LifeSeed) ID() string { return "1401001" }

func (Card1401001LifeSeed) Name() string { return "生命种子" }

func (Card1401001LifeSeed) MasteryMax() int { return 2 }

func (Card1401001LifeSeed) OnMastery(ctx *EffectContext, level int) error {
	if level != 2 {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ps.FindEmptyPosition() == nil {
		return nil
	}
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, func(card *CardInstance) bool {
		return card.Card.IsCompanion() && card.Card.Category == model.ElementEarth
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "life_seed_summon", "生命种子:可以召唤1个地属性伙伴并继承生命种子的加成", candidates, 0, 1, func(selected []string) {
		if len(selected) == 0 {
			return
		}
		cardID := selected[0]
		positions := ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)
		if len(positions) == 0 {
			return
		}
		ctx.Engine.SetPendingAction(ctx.PlayerID, "life_seed_summon_position", "生命种子:选择召唤位置", positions, 1, 1, func(posSelected []string) {
			if len(posSelected) == 0 {
				return
			}
			pos, ok := positionFromSelectionID(posSelected[0])
			if !ok {
				return
			}
			summoned := summonCardFreeFromHandOrDeckAtPosition(ctx, cardID, pos)
			if summoned == nil {
				return
			}
			inheritLifeSeedBonuses(ctx.Engine, ctx.Source, summoned, ctx.PlayerID)
			ctx.Engine.destroyUnit(ctx.Source, ctx.PlayerID)
		})
	})
	return nil
}
