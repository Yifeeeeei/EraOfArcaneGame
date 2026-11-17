package game

import "fmt"

func (Card2121101LavafortAshes) ValidateItemUse(ctx *EffectContext) error {
	e, playerID := ctx.Engine, ctx.PlayerID
	if len(lavaFortAshSourceCandidates(e, playerID)) == 0 {
		return fmt.Errorf("Lavafort Ashes requires a fire skill and a higher-cost fire card in deck")
	}
	return nil
}
