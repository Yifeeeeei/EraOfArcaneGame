package game

import (
	"eraofarcane/model"
)

type Card1121114LegionGeneral struct{ AlwaysActive }

func (Card1121114LegionGeneral) ID() string { return "1121114" }

func (Card1121114LegionGeneral) Name() string { return "军团将星" }

func (Card1121114LegionGeneral) IsPrayerAbility() bool { return true }

func (Card1121114LegionGeneral) OnPerTurn(ctx *EffectContext) error {
	choices := []map[string]any{
		{"instance_id": "power", "name": "火焰法术+2威", "zone": "choice", "side": "own"},
		{"instance_id": "attack", "name": "火焰法术+1攻", "zone": "choice", "side": "own"},
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "legion_general_prayer",
		"军团将星:选择你的火焰法术直到下个回合结束获得的效果", choices, 1, 1,
		func(selected []string) {
			switch firstSelected(selected) {
			case "power":
				ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
					Type:        TempModSkillPowerBonus,
					Element:     model.ElementFire,
					Amount:      2,
					ExpiresTurn: ctx.Engine.State.TurnNumber + 2,
				})
			case "attack":
				ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
					Type:        TempModCurrentTurnElementDamage,
					Element:     model.ElementFire,
					Amount:      1,
					ExpiresTurn: ctx.Engine.State.TurnNumber + 2,
				})
			}
		})
	return nil
}
