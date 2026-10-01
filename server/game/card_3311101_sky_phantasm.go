package game

import (
	"fmt"
)

type Card3311101SkyPhantasm struct{ AlwaysActive }

func (Card3311101SkyPhantasm) ID() string { return "3311101" }

func (Card3311101SkyPhantasm) Name() string { return "苍穹幻韵" }

func (Card3311101SkyPhantasm) PrepareSpellCast(ctx *EffectContext, _ SpellTarget, action ActionMessage) (SpellCastOptions, error) {
	if boosts, _ := action.Data["boost_ids"].([]any); len(boosts) > 0 {
		return SpellCastOptions{}, fmt.Errorf("sky phantasm cannot be boosted directly")
	}
	return SpellCastOptions{ResolveInstead: func() { _ = ctx.Engine.promptSkyPhantasmSpellChoice(ctx.PlayerID, ctx.Source) }}, nil
}

func (e *Engine) promptSkyPhantasmSpellChoice(playerID int, source *CardInstance) error {
	if e == nil || source == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	candidates := e.friendlySkills(playerID, func(skill *CardInstance) bool {
		return skill != nil && skill.Card != nil &&
			skill.InstanceID != source.InstanceID &&
			canUseSkillForPurpose(skill.Card, skillPurposeAttack) &&
			(hasCardTag(skill.Card, "驱动") || hasCardTag(skill.Card, "聚能"))
	})
	if len(candidates) == 0 {
		return nil
	}
	sourceID := source.InstanceID
	e.SetPendingActionWithError(playerID, "sky_phantasm_spell_choice",
		"苍穹幻韵:选择另一个驱动或聚能法术", candidates, 1, 1, nil, false,
		func(selected []string, _ map[string]any) error {
			if e.findSkill(e.State.Players[playerID], sourceID) == nil {
				return nil
			}
			baseSkill := e.findSkill(e.State.Players[playerID], firstSelected(selected))
			if baseSkill == nil || baseSkill.Card == nil || baseSkill.InstanceID == sourceID ||
				!canUseSkillForPurpose(baseSkill.Card, skillPurposeAttack) ||
				(!hasCardTag(baseSkill.Card, "驱动") && !hasCardTag(baseSkill.Card, "聚能")) {
				return fmt.Errorf("invalid sky phantasm spell")
			}
			virtual := e.cloneVirtualSpell(baseSkill, playerID, e.State.TurnNumber)
			targets := e.spellTargetCandidates(playerID, virtual)
			if skillNeedsTargetInstance(virtual) {
				if len(targets) == 0 {
					return nil
				}
				e.SetPendingActionWithError(playerID, "sky_phantasm_target",
					fmt.Sprintf("苍穹幻韵:选择%s的目标", virtual.Card.Name), targets, 1, 1, nil, false,
					func(targetSelected []string, _ map[string]any) error {
						target := selectedSpellTargetFromCandidates(e, playerID, virtual, firstSelected(targetSelected), targets)
						if target == nil {
							return fmt.Errorf("invalid sky phantasm target")
						}
						return e.startVirtualSpellCastNoBoost(playerID, virtual, *target, map[string]any{
							"triggered_by":        "3311101",
							"source_skill":        cardToInfo(baseSkill),
							"source_skill_hidden": false,
						})
					})
				return nil
			}
			return e.startVirtualSpellCastNoBoost(playerID, virtual, SpellTarget{Type: "none"}, map[string]any{
				"triggered_by":        "3311101",
				"source_skill":        cardToInfo(baseSkill),
				"source_skill_hidden": false,
			})
		})
	return nil
}
