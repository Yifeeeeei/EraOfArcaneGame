package game

import "fmt"

func (Card2121110OfferingTorch) ValidateItemUse(ctx *EffectContext) error {
	e, playerID := ctx.Engine, ctx.PlayerID
	if len(e.friendlySkillsIncludingBound(playerID, isFireSpellInstance)) < 2 {
		return fmt.Errorf("Offering Torch requires at least two friendly fire spells")
	}
	return nil
}
