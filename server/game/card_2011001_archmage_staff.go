package game

type Card2011001ArchmageStaff struct{ AlwaysActive }

func (Card2011001ArchmageStaff) ID() string { return "2011001" }

func (Card2011001ArchmageStaff) Name() string { return "大法师之杖" }

func (Card2011001ArchmageStaff) OnEnter(ctx *EffectContext) error {
	staffPlayer := ctx.Engine.State.Players[ctx.PlayerID]
	candidates := make([]map[string]any, 0, len(staffPlayer.SkillPool))
	for _, skill := range staffPlayer.SkillPool {
		if skill == nil || !isSpellLikeCard(skill.Card) {
			continue
		}
		candidates = append(candidates, candidateInfo(skill, "skill_pool", "own"))
	}
	if len(candidates) == 0 {
		return nil
	}
	ctx.Source.Statuses["存储技能"] = 1
	ctx.Engine.SetPendingAction(ctx.PlayerID, "archmage_staff_store_skill", "大法师之杖:选择技能池中的1个法术置于此卡上", candidates, 1, 1, func(selected []string) {
		selectedID := firstSelected(selected)
		for i, skill := range staffPlayer.SkillPool {
			if skill == nil || skill.InstanceID != selectedID || !isSpellLikeCard(skill.Card) {
				continue
			}
			staffPlayer.SkillPool = append(staffPlayer.SkillPool[:i], staffPlayer.SkillPool[i+1:]...)
			skill.SlotIndex = -1
			skill.IsHorizontal = true
			skill.Statuses[archmageStaffStoredSkillStatus] = 1
			ctx.Source.BoundSkills = append(ctx.Source.BoundSkills, skill)
			ctx.Source.Statuses["存储技能"] = 1
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source": cardToInfo(ctx.Source),
				"effect": "store_skill",
				"card":   cardToInfo(skill),
			}})
			return
		}
	})
	return nil
}
