package game

import (
	"eraofarcane/model"
)

type Card2321106EndlessWindTide struct{ AlwaysActive }

func (Card2321106EndlessWindTide) ID() string { return "2321106" }

func (Card2321106EndlessWindTide) Name() string { return "无尽风潮" }

func (Card2321106EndlessWindTide) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || !isFriendlySpellHit(ctx) {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	source := ctx.Source
	removed := false
	for i, skill := range ps.Skills {
		if skill == ctx.Source {
			ps.Skills[i] = nil
			source = ctx.Engine.leaveBattlefield(source)
			removed = true
			break
		}
	}
	if !removed {
		removed = ctx.Engine.removeCardFromGraveyard(ctx.PlayerID, ctx.Source)
	}
	if !removed {
		return nil
	}
	source.SlotIndex = -1
	source.IsHorizontal = true
	source.PowerBonus += 2
	source.Statuses[endlessWindTideAirCostBonusStatus]++
	source.Statuses[StatusCannotUseSkillUntilTurn] = ctx.Engine.State.TurnNumber
	ps.Hand = append(ps.Hand, source)
	ctx.Engine.emit(GameEvent{
		Type:   "endless_wind_tide_return",
		Player: -1,
		Data: map[string]any{
			"player": ctx.PlayerID,
			"card":   cardToInfo(source),
		},
	})
	return nil
}

func (Card2321106EndlessWindTide) ModifySelfCardPlayCost(ctx *EffectContext, cost map[string]int) {
	if ctx == nil || ctx.Source == nil {
		return
	}
	cost[model.ElementAir] += ctx.Source.Statuses[endlessWindTideAirCostBonusStatus]
}

var _ OnSpellHitBehavior = Card2321106EndlessWindTide{}

var _ SelfCardPlayCostModifier = Card2321106EndlessWindTide{}

const endlessWindTideAirCostBonusStatus = "endless_wind_tide_air_cost_bonus"
