package game

import (
	"fmt"
)

type Card1321106SilverleafRanger struct{ AlwaysActive }

func (Card1321106SilverleafRanger) ID() string { return "1321106" }

func (Card1321106SilverleafRanger) Name() string { return "银叶游侠" }

func (Card1321106SilverleafRanger) PerTurnLabel(*CardInstance) string {
	return "主动"
}

func (Card1321106SilverleafRanger) OnPerTurn(ctx *EffectContext) error {
	if !ctx.Engine.canConsumeCard(ctx.Source) {
		return fmt.Errorf("银叶游侠不能被消耗")
	}
	ctx.Source.IsHorizontal = true
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:          TempModSkillAttackBonus,
		Amount:        1,
		RemainingUses: 1,
		ExpiresTurn:   ctx.Engine.State.TurnNumber + 2,
	})
	return nil
}
