package game

import (
	"fmt"
)

type Card2521102MoonlightDust struct{ AlwaysActive }

func (Card2521102MoonlightDust) ID() string { return "2521102" }

func (Card2521102MoonlightDust) Name() string { return "月霞之尘" }

func (Card2521102MoonlightDust) OnUseItem(ctx *EffectContext) error {
	choices := make([]map[string]any, 0, 2)
	if ctx.Engine.hasEnemySetCounter(ctx.PlayerID) {
		choices = append(choices, map[string]any{"instance_id": "destroy_counters", "name": "摧毁敌方盖放的所有卡牌", "zone": "choice", "side": "own"})
	}
	if ctx.Engine.hasEnemyFrontStealth(ctx.PlayerID) {
		choices = append(choices, map[string]any{"instance_id": "remove_front_stealth", "name": "使前排敌人失去隐蔽", "zone": "choice", "side": "own"})
	}
	if len(choices) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "moonlight_dust_mode",
		"月霞之尘:选择1项效果", choices, 1, 1,
		func(selected []string) {
			switch firstSelected(selected) {
			case "destroy_counters":
				ctx.Engine.destroyEnemySetCounters(ctx.PlayerID)
			case "remove_front_stealth":
				ctx.Engine.removeEnemyFrontStealth(ctx.PlayerID)
			}
		})
	return nil
}

func (Card2521102MoonlightDust) ValidateItemUse(ctx *EffectContext) error {
	e, playerID := ctx.Engine, ctx.PlayerID
	if !e.hasEnemySetCounter(playerID) && !e.hasEnemyFrontStealth(playerID) {
		return fmt.Errorf("Moonlight Dust requires enemy set counters or stealthy front enemies")
	}
	return nil
}

func (e *Engine) hasEnemySetCounter(playerID int) bool {
	opponent := e.State.Players[1-playerID]
	if opponent == nil {
		return false
	}
	for _, card := range opponent.Equipment {
		if card != nil && card.IsSetCounter {
			return true
		}
	}
	return false
}

func (e *Engine) destroyEnemySetCounters(playerID int) {
	opponentID := 1 - playerID
	opponent := e.State.Players[opponentID]
	if opponent == nil {
		return
	}
	for i := range opponent.Equipment {
		card := opponent.Equipment[i]
		if card != nil && card.IsSetCounter {
			e.moveEquipmentToGraveyard(opponentID, i, card)
		}
	}
}

func (e *Engine) hasEnemyFrontStealth(playerID int) bool {
	opponent := e.State.Players[1-playerID]
	if opponent == nil {
		return false
	}
	frontRow := opponent.GetFrontRow()
	if frontRow < 0 || frontRow >= 3 {
		return false
	}
	for col := 0; col < 3; col++ {
		unit := opponent.Units[col][frontRow]
		if unit != nil && unit.Statuses[StatusStealth] > 0 {
			return true
		}
	}
	return false
}

func (e *Engine) removeEnemyFrontStealth(playerID int) {
	opponent := e.State.Players[1-playerID]
	if opponent == nil {
		return
	}
	frontRow := opponent.GetFrontRow()
	if frontRow < 0 || frontRow >= 3 {
		return
	}
	for col := 0; col < 3; col++ {
		unit := opponent.Units[col][frontRow]
		if unit != nil {
			delete(unit.Statuses, StatusStealth)
		}
	}
}
