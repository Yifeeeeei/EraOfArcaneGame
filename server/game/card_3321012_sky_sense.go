package game

type Card3321012SkySense struct{ AlwaysActive }

func (Card3321012SkySense) ID() string   { return "3321012" }
func (Card3321012SkySense) Name() string { return "空天感应" }

func (Card3321012SkySense) HasActiveSpellStatModifier(card *CardInstance) bool {
	return abilityDurationActive(card)
}

func (Card3321012SkySense) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if ctx.ExtraData["purpose"] != string(skillPurposeAttack) {
		return
	}
	target, ok := ctx.ExtraData["spell_target"].(SpellTarget)
	if !ok || target.Type != "unit" {
		return
	}
	defenderID := ctx.OpponentID
	frontRow := ctx.Engine.State.Players[defenderID].GetFrontRow()
	if frontRow < 0 {
		return
	}
	for _, unit := range ctx.Engine.spellAffectedUnits(defenderID, ctx.Target, target) {
		if unit.Position != nil && unit.Position.Row > frontRow {
			stats.PowerBonus += 2
			return
		}
	}
}
