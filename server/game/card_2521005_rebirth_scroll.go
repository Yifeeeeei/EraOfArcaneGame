package game

import "eraofarcane/model"

type Card2521005RebirthScroll struct{ AlwaysActive }

func (Card2521005RebirthScroll) ID() string   { return "2521005" }
func (Card2521005RebirthScroll) Name() string { return "新生卷轴" }

func (Card2521005RebirthScroll) OnUseItem(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	candidates := make([]map[string]any, 0)
	for _, card := range ps.Graveyard {
		if card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementLight && ps.CanPayCost(card.Card.ElementsCost) {
			candidates = append(candidates, candidateInfo(card, "graveyard", "own"))
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "rebirth_scroll",
		"选择1个死亡的光辉伙伴，支付入场费用并复活", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			ctx.Engine.reviveCompanionFromGraveyard(ctx.PlayerID, selected[0])
		})
	return nil
}
