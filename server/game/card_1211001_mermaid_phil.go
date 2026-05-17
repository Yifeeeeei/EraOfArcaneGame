package game

type Card1211001MermaidPhil struct{ AlwaysActive }

func (Card1211001MermaidPhil) ID() string   { return "1211001" }
func (Card1211001MermaidPhil) Name() string { return "人鱼 菲尔" }

func (Card1211001MermaidPhil) OnPerTurn(ctx *EffectContext) error {
	if ctx.Source == nil || ctx.Source.Position == nil {
		return nil
	}
	for _, unit := range adjacentUnits(ctx.Engine.State.Players[ctx.PlayerID], ctx.Source.Position) {
		if unit.Card.IsCompanion() {
			return nil
		}
	}
	searchDeckToHandByPredicate(ctx, "mermaid_search_water_companion", "检索1张水纹伙伴", isWaterCompanion)
	return nil
}
