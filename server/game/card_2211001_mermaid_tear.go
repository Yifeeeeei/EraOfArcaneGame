package game

type Card2211001MermaidTear struct{ AlwaysActive }

func (Card2211001MermaidTear) ID() string   { return "2211001" }
func (Card2211001MermaidTear) Name() string { return "人鱼之泪" }
func (Card2211001MermaidTear) OnUltimate(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	candidates := make([]map[string]any, 0)
	for _, card := range ps.Graveyard {
		if card != nil && card.Card.IsCompanion() {
			candidates = append(candidates, candidateInfo(card, "graveyard", "own"))
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "mermaid_tear_revive",
		"选择1个死亡伙伴复活并令其只有1血", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			reviveID := selected[0]
			positions := ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)
			if len(positions) == 0 {
				return
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "mermaid_tear_position",
				"选择复活位置", positions, 1, 1,
				func(posSelected []string) {
					if len(posSelected) == 0 {
						return
					}
					pos, ok := positionFromSelectionID(posSelected[0])
					if !ok {
						return
					}
					ctx.Engine.removeEquipmentFromGame(ctx.PlayerID, ctx.Source.InstanceID)
					ctx.Engine.reviveCompanionFromGraveyardWithLifeAtPosition(ctx.PlayerID, reviveID, 1, false, pos)
				})
		})
	return nil
}
