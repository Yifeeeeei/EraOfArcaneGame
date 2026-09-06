package game

import "fmt"

type Card2221101MirrorseaSpring struct{ royalWaterUseCostReduction }

func (Card2221101MirrorseaSpring) ValidateItemUse(ctx *EffectContext) error {
	e, playerID := ctx.Engine, ctx.PlayerID
	if len(e.friendlySkillsIncludingBound(playerID, func(skill *CardInstance) bool {
		return skill != nil && skill.Card != nil && isSpellLikeCard(skill.Card)
	})) == 0 {
		return fmt.Errorf("Mirrorsea Spring requires a friendly spell")
	}
	return nil
}
