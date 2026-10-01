package game

import (
	"eraofarcane/model"
)

type Card2111001FireDragonHeart struct{ AlwaysActive }

func (Card2111001FireDragonHeart) ID() string { return "2111001" }

func (Card2111001FireDragonHeart) Name() string { return "火龙之心" }

func (Card2111001FireDragonHeart) OnPerTurn(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	spend := min(ps.Elements[model.ElementFire], 3)
	if spend > 0 {
		ps.Elements[model.ElementFire] -= spend
		ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{Type: TempModSkillPowerBonus, Amount: spend * 3, RemainingUses: 1, ExpiresTurn: ctx.Engine.State.TurnNumber + 1})
	}
	return nil
}
