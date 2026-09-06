package game

import (
	"fmt"
)

type Card3211102DragonSnowfield struct{ AlwaysActive }

func (Card3211102DragonSnowfield) ID() string { return "3211102" }

func (Card3211102DragonSnowfield) Name() string { return "龙吟雪域" }

func (Card3211102DragonSnowfield) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	casterID, ok := spellCasterFromData(ctx)
	if !ok || casterID < 0 || casterID >= len(ctx.Engine.State.Players) {
		return nil
	}
	candidates := ctx.Engine.dragonSnowfieldFreezeCandidates(casterID)
	if len(candidates) == 0 {
		return nil
	}
	sourceID := ctx.Source.InstanceID
	ownerID := ctx.PlayerID
	ctx.Engine.SetPendingActionWithError(casterID, "dragon_snowfield_freeze",
		"龙吟雪域:选择法力范围内1个单位冻结1", candidates, 1, 1, nil, false,
		func(selected []string, _ map[string]any) error {
			target := ctx.Engine.findUnitByInstanceID(firstSelected(selected))
			if target == nil || target.OwnerID != 1-casterID || target.Position == nil || !ctx.Engine.IsInSpellRange(casterID, target.Position.Col, target.Position.Row, false) {
				return fmt.Errorf("invalid dragon snowfield freeze target")
			}
			target.Statuses[StatusFreeze]++
			source := ctx.Engine.findSkill(ctx.Engine.State.Players[ownerID], sourceID)
			if source == nil {
				return nil
			}
			source.Statuses[dragonSnowfieldTriggerStatus]++
			ctx.Engine.promptDragonSnowfieldSummon(ownerID, source)
			return nil
		})
	return nil
}

const dragonSnowfieldTriggerStatus = "dragon_snowfield_trigger_count"

func (e *Engine) dragonSnowfieldFreezeCandidates(casterID int) []map[string]any {
	if e == nil || casterID < 0 || casterID >= len(e.State.Players) {
		return nil
	}
	candidates := []map[string]any{}
	ownerID := 1 - casterID
	ps := e.State.Players[ownerID]
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := ps.Units[col][row]
			if unit == nil || unit.CurrentLife <= 0 || !e.IsInSpellRange(casterID, col, row, false) {
				continue
			}
			info := cardToInfo(unit)
			info["side"] = "enemy"
			info["zone"] = "unit"
			info["position"] = Position{Col: col, Row: row}
			candidates = append(candidates, info)
		}
	}
	return candidates
}

func (e *Engine) promptDragonSnowfieldSummon(ownerID int, source *CardInstance) {
	if e == nil || source == nil || source.Statuses[dragonSnowfieldTriggerStatus] < 5 || ownerID < 0 || ownerID >= len(e.State.Players) {
		return
	}
	positions := e.friendlyEmptyUnitPositions(ownerID)
	if len(positions) == 0 {
		return
	}
	sourceID := source.InstanceID
	e.SetPendingActionWithError(ownerID, "dragon_snowfield_summon_frost_dragon",
		"龙吟雪域:是否召唤凛冰之龙", positions, 0, 1, nil, false,
		func(selected []string, _ map[string]any) error {
			if len(selected) == 0 {
				return nil
			}
			pos, ok := positionFromSelectionID(firstSelected(selected))
			if !ok || !pos.Valid() || e.State.Players[ownerID].Units[pos.Col][pos.Row] != nil {
				return fmt.Errorf("invalid dragon snowfield summon position")
			}
			source := e.findSkill(e.State.Players[ownerID], sourceID)
			if source == nil || source.Statuses[dragonSnowfieldTriggerStatus] < 5 {
				return nil
			}
			if e.summonFreshCardAtPosition(ownerID, "1201101", pos, true) == nil {
				return fmt.Errorf("failed to summon frost dragon")
			}
			source.Statuses[dragonSnowfieldTriggerStatus] -= 5
			return nil
		})
}
