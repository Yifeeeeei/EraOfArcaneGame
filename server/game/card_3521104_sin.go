package game

import (
	"strings"
)

type Card3521104Sin struct{ AlwaysActive }

func (Card3521104Sin) ID() string { return "3521104" }

func (Card3521104Sin) Name() string { return "罪责" }

func (Card3521104Sin) OnEnter(ctx *EffectContext) error {
	choices := companionTagChoices()
	ctx.Engine.SetPendingAction(ctx.PlayerID, "sin_choose_companion_kind",
		"罪责:选择1个伙伴种类", choices, 1, 1,
		func(selected []string) {
			tag := firstSelected(selected)
			if !validCompanionTagChoice(tag) {
				return
			}
			clearStatusPrefix(ctx.Source, sinTargetTagStatusPrefix)
			ctx.Source.Statuses[sinTargetTagStatusPrefix+tag] = 1
		})
	return nil
}

func (Card3521104Sin) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Source == nil || ctx.ExtraData == nil || ctx.ExtraData["purpose"] != string(skillPurposeAttack) {
		return
	}
	target, _ := ctx.ExtraData["spell_target_unit"].(*CardInstance)
	if !sinMatchesTargetTag(ctx.Source, target) {
		return
	}
	stats.PowerBonus += 2
	stats.Pierce = true
}

const sinTargetTagStatusPrefix = "罪责目标种类:"

func companionTagChoices() []map[string]any {
	tags := []string{"人类", "巫师", "野兽", "精灵", "恶魔", "造物", "植物", "灵体", "龙"}
	choices := make([]map[string]any, 0, len(tags))
	for _, tag := range tags {
		choices = append(choices, map[string]any{
			"instance_id": tag,
			"name":        tag,
			"zone":        "choice",
			"side":        "own",
		})
	}
	return choices
}

func validCompanionTagChoice(tag string) bool {
	for _, choice := range companionTagChoices() {
		if choice["instance_id"] == tag {
			return true
		}
	}
	return false
}

func sinMatchesTargetTag(skill *CardInstance, target *CardInstance) bool {
	if skill == nil || target == nil || target.Card == nil || !target.Card.IsCompanion() {
		return false
	}
	for status, amount := range skill.Statuses {
		if amount <= 0 || !strings.HasPrefix(status, sinTargetTagStatusPrefix) {
			continue
		}
		tag := strings.TrimPrefix(status, sinTargetTagStatusPrefix)
		return hasCardTag(target.Card, tag)
	}
	return false
}
