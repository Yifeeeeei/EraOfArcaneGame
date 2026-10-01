package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card1411102WhisperElfKingSindariel struct{ AlwaysActive }

func (Card1411102WhisperElfKingSindariel) ID() string { return "1411102" }

func (Card1411102WhisperElfKingSindariel) Name() string { return "谧语精灵王 辛达瑞尔" }

func (Card1411102WhisperElfKingSindariel) OnEnter(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.PlayerID < 0 || ctx.PlayerID >= len(ctx.Engine.State.Players) {
		return nil
	}
	opponentID := 1 - ctx.PlayerID
	opponent := ctx.Engine.State.Players[opponentID]
	if opponent == nil {
		return nil
	}
	maxTargets := 0
	if opponent.SpellHitsLastTurn >= 3 {
		maxTargets++
	}
	if opponent.SpellHitTargetsLastTurn >= 3 {
		maxTargets++
	}
	if opponent.SpellDamageLastTurn >= 3 {
		maxTargets++
	}
	if maxTargets == 0 {
		return nil
	}
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, true, nil)
	if len(candidates) == 0 {
		return nil
	}
	if maxTargets > len(candidates) {
		maxTargets = len(candidates)
	}
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "sindariel_entry_damage",
		"谧语精灵王 辛达瑞尔:选择敌人造成2点伤害", candidates, 0, maxTargets,
		nil, false, func(selected []string, _ map[string]any) error {
			allowed := make(map[string]bool, len(candidates))
			for _, candidate := range candidates {
				if id, _ := candidate["instance_id"].(string); id != "" {
					allowed[id] = true
				}
			}
			for _, id := range selected {
				if !allowed[id] {
					return fmt.Errorf("invalid Sindariel target")
				}
				target := ctx.Engine.findUnitByInstanceID(id)
				if target == nil || target.OwnerID != opponentID {
					continue
				}
				ctx.Engine.ApplyDamage(DamageRequest{Target: target, Amount: 2, Kind: "effect", Element: model.ElementEarth, Source: ctx.Source, SourcePlayer: ctx.PlayerID, SourceKnown: true})
				if target.CurrentLife <= 0 && !target.Card.IsHero() && ctx.Engine.unitInOwnerGrid(target, opponentID) {
					ctx.Engine.destroyUnitWithData(target, opponentID, map[string]any{
						"death_cause": "sindariel",
						"source_card": ctx.Source,
						"attacker":    ctx.PlayerID,
					})
				}
			}
			return nil
		})
	return nil
}
