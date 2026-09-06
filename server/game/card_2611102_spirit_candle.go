package game

type Card2611102SpiritCandle struct{ AlwaysActive }

func (Card2611102SpiritCandle) ID() string { return "2611102" }

func (Card2611102SpiritCandle) Name() string { return "渡灵之烛" }

func (Card2611102SpiritCandle) OnEquip(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil {
		return nil
	}
	ctx.Engine.enforceSlotCapacities(ctx.Engine.State.Players[ctx.PlayerID])
	return nil
}

func (Card2611102SpiritCandle) SlotGrant(*CardInstance) SlotGrant {
	return SlotGrant{Group: "spirit_candle", SkillSlots: 2, SkillTags: []string{"灵媒", "神秘"}}
}

func (Card2611102SpiritCandle) SpellPowerScale(_ *CardInstance, skill *CardInstance) SpellPowerScale {
	if skill == nil || skill.Card == nil || hasCardTag(skill.Card, "灵媒") || hasCardTag(skill.Card, "神秘") {
		return SpellPowerScale{}
	}
	return SpellPowerScale{Group: "spirit_candle", Numerator: 1, Denominator: 2}
}
