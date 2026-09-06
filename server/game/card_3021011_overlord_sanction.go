package game

import "fmt"

type Card3021011OverlordSanction struct{ AlwaysActive }

func (Card3021011OverlordSanction) ID() string { return "3021011" }

func (Card3021011OverlordSanction) Name() string { return "统御者的制裁" }

func (Card3021011OverlordSanction) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	stats.DamageBonus += 1
}

func (Card3021011OverlordSanction) ValidateOwnCost(ctx *EffectContext, cost map[string]int, action ActionMessage) error {
	if !validateSingleElementPayment(ctx.Engine.State.Players[ctx.PlayerID].Elements, cost, action) {
		return fmt.Errorf("overlord sanction cost must be paid with one element")
	}
	return nil
}
