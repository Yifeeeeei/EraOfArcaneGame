package game

type Card3421110CrushingStone struct{ AlwaysActive }

func (Card3421110CrushingStone) ID() string   { return "3421110" }
func (Card3421110CrushingStone) Name() string { return "粉碎石破" }

func (Card3421110CrushingStone) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.ExtraData == nil || ctx.ExtraData["purpose"] != string(skillPurposeAttack) {
		return
	}
	target, _ := ctx.ExtraData["spell_target_unit"].(*CardInstance)
	if target != nil && target.CurrentLife > 2 {
		stats.PowerBonus++
	}
}
