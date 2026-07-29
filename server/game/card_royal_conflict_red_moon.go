package game

import "eraofarcane/model"

type Card3611101RedMoon struct{ AlwaysActive }

func (Card3611101RedMoon) ID() string   { return "3611101" }
func (Card3611101RedMoon) Name() string { return "红月" }
func (Card3611101RedMoon) HasActiveSpellStatModifier(card *CardInstance) bool {
	return abilityDurationActive(card)
}
func (Card3611101RedMoon) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if !isAttackPurpose(ctx) || ctx.Target == nil || ctx.Target.Card == nil || ctx.Target.Card.Category != model.ElementShadow {
		return
	}
	stats.PowerBonus += 2
}

type Card1621110ScarletBeast struct{ AlwaysActive }

func (Card1621110ScarletBeast) ID() string   { return "1621110" }
func (Card1621110ScarletBeast) Name() string { return "猩红魔兽" }
func (Card1621110ScarletBeast) HasActiveSpellStatModifier(card *CardInstance) bool {
	return card != nil && card.Statuses[StatusPetrify] <= 0
}
func (Card1621110ScarletBeast) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if !isAttackPurpose(ctx) || !ctx.Engine.redMoonActive(ctx.PlayerID) || ctx.Target == nil || ctx.Target.Card == nil || ctx.Target.Card.Category != model.ElementShadow {
		return
	}
	stats.PowerBonus += 2
}

type Card3621107WillErosion struct{ AlwaysActive }

func (Card3621107WillErosion) ID() string   { return "3621107" }
func (Card3621107WillErosion) Name() string { return "意志侵蚀" }
func (Card3621107WillErosion) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if !isAttackPurpose(ctx) || !ctx.Engine.redMoonActive(ctx.PlayerID) {
		return
	}
	stats.PowerBonus++
	stats.Pierce = true
}

func isAttackPurpose(ctx *EffectContext) bool {
	if ctx == nil || ctx.ExtraData == nil {
		return true
	}
	purpose, _ := ctx.ExtraData["purpose"].(string)
	return purpose == "" || purpose == string(skillPurposeAttack)
}
