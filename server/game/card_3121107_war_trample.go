package game

type Card3121107WarTrample struct{ AlwaysActive }

func (Card3121107WarTrample) ID() string { return "3121107" }

func (Card3121107WarTrample) Name() string { return "战争践踏" }

func (Card3121107WarTrample) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Source == nil || ctx.Target != ctx.Source || ctx.ExtraData["purpose"] != string(skillPurposeAttack) || ctx.ExtraData["stat"] != "damage" {
		return
	}
	stats.DamageBonus -= spellTargetUnitCount(ctx.ExtraData)
}

func spellTargetUnitCount(data map[string]any) int {
	if units, ok := data["affected_units"].([]*CardInstance); ok {
		return len(units)
	}
	if targets, ok := data["spell_targets"].([]SpellTarget); ok {
		count := 0
		for _, target := range targets {
			if target.Type == "unit" {
				count++
			}
		}
		return count
	}
	if target, ok := data["spell_target"].(SpellTarget); ok && target.Type == "unit" {
		return 1
	}
	return 0
}
