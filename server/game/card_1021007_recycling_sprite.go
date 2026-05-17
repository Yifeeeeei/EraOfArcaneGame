package game

type Card1021007RecyclingSprite struct{ AlwaysActive }

func (Card1021007RecyclingSprite) ID() string   { return "1021007" }
func (Card1021007RecyclingSprite) Name() string { return "回收小精灵" }

func (Card1021007RecyclingSprite) OnEnter(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if len(ps.Graveyard) == 0 {
		return nil
	}
	candidates := make([]map[string]any, 0, len(ps.Graveyard))
	for _, card := range ps.Graveyard {
		if card == nil {
			continue
		}
		candidates = append(candidates, candidateInfo(card, "graveyard", "own"))
	}
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "graveyard_to_deck_top",
		"选择1张弃牌堆中的牌放到卡组顶", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			ctx.Engine.moveGraveyardCardToDeckTop(ctx.PlayerID, selected[0])
		})
	return nil
}
