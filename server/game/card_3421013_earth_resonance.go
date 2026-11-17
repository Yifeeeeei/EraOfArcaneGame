package game

type Card3421013EarthResonance struct{ AlwaysActive }

func (Card3421013EarthResonance) ID() string { return "3421013" }

func (Card3421013EarthResonance) Name() string { return "大地共鸣" }

func (Card3421013EarthResonance) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx.ExtraData["purpose"] != string(skillPurposeAttack) {
		return
	}
	player := ctx.Engine.State.Players[ctx.PlayerID]
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := player.Units[col][row]
			if unit == nil || unit.Card == nil {
				continue
			}
			if totalLoad(unit) > 3 || unit.CurrentLife > 3 {
				stats.DamageBonus++
			}
		}
	}
}
