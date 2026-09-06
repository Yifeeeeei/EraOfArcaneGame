package game

import (
	"fmt"
)

type Card2621107CurseBox struct{ AlwaysActive }

func (Card2621107CurseBox) ID() string { return "2621107" }

func (Card2621107CurseBox) Name() string { return "诅咒魔盒" }

func (Card2621107CurseBox) OnFriendlyDeath(ctx *EffectContext) error {
	return addCurseBoxMarker(ctx)
}

func (Card2621107CurseBox) OnEnemyDeath(ctx *EffectContext) error {
	return addCurseBoxMarker(ctx)
}

func (Card2621107CurseBox) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Source == nil || ctx.Engine == nil {
		return nil
	}
	markers := ctx.Source.Statuses[curseBoxMarkerStatus]
	if markers <= 0 {
		return fmt.Errorf("诅咒魔盒没有标记物")
	}
	candidates := ctx.Engine.enemySkills(ctx.PlayerID, canInstanceBeWeakened)
	maxSelect := min(3, min(markers, len(candidates)))
	if maxSelect <= 0 {
		return fmt.Errorf("没有可虚弱的敌方法术")
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "curse_box_weaken",
		"诅咒魔盒:移除最多3个标记物，使等量敌方法术虚弱1", candidates, 1, maxSelect,
		func(selected []string) {
			removed := 0
			seen := map[string]bool{}
			for _, id := range selected {
				if seen[id] || ctx.Source.Statuses[curseBoxMarkerStatus] <= 0 {
					continue
				}
				seen[id] = true
				for _, skill := range ctx.Engine.State.Players[ctx.OpponentID].Skills {
					if skill != nil && skill.InstanceID == id && canInstanceBeWeakened(skill) {
						ctx.Engine.addStatus(skill, StatusWeaken, 1)
						ctx.Source.Statuses[curseBoxMarkerStatus]--
						removed++
						break
					}
				}
			}
			if ctx.Source.Statuses[curseBoxMarkerStatus] <= 0 {
				delete(ctx.Source.Statuses, curseBoxMarkerStatus)
			}
			if removed == 0 && ctx.Source.Statuses[curseBoxMarkerStatus] <= 0 {
				delete(ctx.Source.Statuses, curseBoxMarkerStatus)
			}
		})
	return nil
}

func addCurseBoxMarker(ctx *EffectContext) error {
	if ctx == nil || ctx.Source == nil {
		return nil
	}
	ctx.Source.Statuses[curseBoxMarkerStatus]++
	return nil
}
