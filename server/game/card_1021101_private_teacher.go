package game

import (
	"fmt"
	"strings"
)

type Card1021101PrivateTeacher struct{ AlwaysActive }

func (Card1021101PrivateTeacher) ID() string { return "1021101" }

func (Card1021101PrivateTeacher) Name() string { return "私家教师" }

func (Card1021101PrivateTeacher) OnEnter(ctx *EffectContext) error {
	candidates := make([]map[string]any, 0)
	allowed := make(map[string]bool)
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	for _, skill := range ps.SkillPool {
		if skill == nil || skill.Card == nil || !skill.Card.IsSkill() || totalElementCost(skill.Card.ElementsCost) >= 4 {
			continue
		}
		hasEmptySlot := false
		for i := 0; i < skillSlotCapacity(ps); i++ {
			if ps.Skills[i] == nil && skillAllowedInSlot(ps, skill, i) {
				hasEmptySlot = true
				break
			}
		}
		if hasEmptySlot {
			candidates = append(candidates, candidateInfo(skill, "skill_pool", "own"))
			allowed[skill.InstanceID] = true
			continue
		}
		for slotIdx, learned := range ps.Skills {
			if learned == nil || learned.IsHorizontal {
				continue
			}
			if !skillAllowedInSlot(ps, skill, slotIdx) {
				continue
			}
			id := skill.InstanceID + "|" + learned.InstanceID
			candidate := candidateInfo(skill, "skill_pool", "own")
			candidate["instance_id"] = id
			candidate["name"] = fmt.Sprintf("学习%s，替换%s", skill.Card.Name, learned.Card.Name)
			candidate["replace_id"] = learned.InstanceID
			allowed[id] = true
			candidates = append(candidates, candidate)
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "private_teacher_learn_skill",
		"私家教师:选择1个学习花费小于4的法术免费学习", candidates, 1, 1,
		func(selected []string) {
			id := firstSelected(selected)
			if !allowed[id] {
				return
			}
			skillID := id
			replaceID := ""
			if before, after, ok := strings.Cut(id, "|"); ok {
				skillID = before
				replaceID = after
			}
			ctx.Engine.learnSkillFromPoolWithoutCost(ctx.PlayerID, skillID, replaceID)
		})
	return nil
}
