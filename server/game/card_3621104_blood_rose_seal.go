package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card3621104BloodRoseSeal struct{ AlwaysActive }

func (Card3621104BloodRoseSeal) ID() string { return "3621104" }

func (Card3621104BloodRoseSeal) Name() string { return "血蔷薇咒印" }

func (Card3621104BloodRoseSeal) OnEnter(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, false, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "blood_rose_seal_mark",
		"血蔷薇咒印:选择1个敌方单位标记", candidates, 1, 1,
		nil, false, func(selected []string, _ map[string]any) error {
			target := selectedUnitFromCandidates(ctx.Engine, selected, candidates)
			if target == nil || target.OwnerID != ctx.OpponentID {
				return fmt.Errorf("invalid blood rose seal target")
			}
			target.Statuses[bloodRoseSealMarkerStatus(ctx.Source)] = 1
			ctx.Source.Statuses[bloodRoseSealExpireTurnStatus] = ctx.Engine.State.TurnNumber + 2
			ctx.Engine.emit(GameEvent{
				Type:   "blood_rose_seal_mark",
				Player: -1,
				Data: map[string]any{
					"player": ctx.PlayerID,
					"source": cardToInfo(ctx.Source),
					"target": cardToInfo(target),
				},
			})
			return nil
		})
	return nil
}

func (Card3621104BloodRoseSeal) OnEnemyDeath(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Source.Card == nil {
		return nil
	}
	if ctx.Source.Statuses[bloodRoseSealExpireTurnStatus] > 0 && ctx.Engine.State.TurnNumber > ctx.Source.Statuses[bloodRoseSealExpireTurnStatus] {
		return nil
	}
	if ctx.Target.Statuses[bloodRoseSealMarkerStatus(ctx.Source)] <= 0 {
		return nil
	}
	if ctx.Engine.findSkill(ctx.Engine.State.Players[ctx.PlayerID], ctx.Source.InstanceID) != ctx.Source {
		return nil
	}
	if !ctx.Engine.bindSkillToHero(ctx.PlayerID, ctx.Source) {
		return nil
	}
	ctx.Source.Statuses["使用费用"+model.ElementShadow+"-1"]++
	delete(ctx.Source.Statuses, bloodRoseSealExpireTurnStatus)
	ctx.Engine.clearBloodRoseSealMarkers(ctx.Source)
	ctx.Engine.emit(GameEvent{
		Type:   "blood_rose_seal_bound",
		Player: -1,
		Data: map[string]any{
			"player": ctx.PlayerID,
			"source": cardToInfo(ctx.Source),
			"target": cardToInfo(ctx.Target),
		},
	})
	return nil
}

func (Card3621104BloodRoseSeal) OnTurnEnd(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	endedPlayer := ctx.PlayerID
	if ctx.ExtraData != nil {
		endedPlayer, _ = ctx.ExtraData["ended_player"].(int)
	}
	if endedPlayer != ctx.PlayerID {
		return nil
	}
	expires := ctx.Source.Statuses[bloodRoseSealExpireTurnStatus]
	if expires <= 0 || ctx.Engine.State.TurnNumber < expires {
		return nil
	}
	delete(ctx.Source.Statuses, bloodRoseSealExpireTurnStatus)
	ctx.Engine.clearBloodRoseSealMarkers(ctx.Source)
	return nil
}

const bloodRoseSealExpireTurnStatus = "blood_rose_seal_expire_turn"

func bloodRoseSealMarkerStatus(source *CardInstance) string {
	if source == nil {
		return "blood_rose_seal_mark:"
	}
	return "blood_rose_seal_mark:" + source.InstanceID
}

func (e *Engine) bindSkillToHero(playerID int, skill *CardInstance) bool {
	if e == nil || skill == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return false
	}
	ps := e.State.Players[playerID]
	if ps == nil || ps.Hero == nil {
		return false
	}
	for _, bound := range ps.Hero.BoundSkills {
		if bound == skill {
			return true
		}
	}
	for i := range ps.Skills {
		if ps.Skills[i] == skill {
			ps.Skills[i] = nil
			break
		}
	}
	skill.SlotIndex = -1
	markTransferredBoundSkill(skill)
	ps.Hero.BoundSkills = append(ps.Hero.BoundSkills, skill)
	return true
}

func (e *Engine) clearBloodRoseSealMarkers(source *CardInstance) {
	if e == nil || source == nil {
		return
	}
	marker := bloodRoseSealMarkerStatus(source)
	for _, ps := range e.State.Players {
		if ps == nil {
			continue
		}
		for col := 0; col < 3; col++ {
			for row := 0; row < 3; row++ {
				if unit := ps.Units[col][row]; unit != nil {
					delete(unit.Statuses, marker)
				}
			}
		}
	}
}
