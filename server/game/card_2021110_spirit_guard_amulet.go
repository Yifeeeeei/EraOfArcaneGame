package game

import "eraofarcane/model"

type Card2021110SpiritGuardAmulet struct{ AlwaysActive }

func (Card2021110SpiritGuardAmulet) ID() string   { return "2021110" }
func (Card2021110SpiritGuardAmulet) Name() string { return "灵守护符" }

func (Card2021110SpiritGuardAmulet) ModifyElementsGain(ctx *EffectContext, target *CardInstance, gains map[string]int) {
	if ctx == nil || ctx.Engine == nil || ctx.Target == nil || target == nil || target != ctx.Target {
		return
	}
	if !ctx.Engine.isOnlyFriendlyEquipment(ctx.PlayerID, ctx.Target) {
		return
	}
	gains[model.ElementArcane]++
}

func (e *Engine) isOnlyFriendlyEquipment(playerID int, equipment *CardInstance) bool {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) || equipment == nil {
		return false
	}
	ps := e.State.Players[playerID]
	if ps == nil {
		return false
	}
	count := 0
	found := false
	for _, card := range ps.Equipment {
		if card == nil {
			continue
		}
		count++
		if card == equipment {
			found = true
		}
	}
	return found && count == 1
}
