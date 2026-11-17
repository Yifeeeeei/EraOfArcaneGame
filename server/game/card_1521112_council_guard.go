package game

import (
	"fmt"
)

type Card1521112CouncilGuard struct{ AlwaysActive }

func (Card1521112CouncilGuard) ID() string { return "1521112" }

func (Card1521112CouncilGuard) Name() string { return "议庭护法" }

func (Card1521112CouncilGuard) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil {
		return nil
	}
	opponent := ctx.Engine.State.Players[ctx.OpponentID]
	lookCount := min(5, len(opponent.Deck))
	if lookCount == 0 {
		return fmt.Errorf("议庭护法需要对手牌库中有牌")
	}
	looked := append([]*CardInstance(nil), opponent.Deck[:lookCount]...)
	hasMark := false
	for _, card := range looked {
		if card != nil && card.Card != nil && card.Card.Number == "2001102" {
			hasMark = true
			break
		}
	}
	candidates := make([]map[string]any, 0, len(looked))
	for i, card := range looked {
		info := candidateInfo(card, "deck", "enemy")
		info["deck_index"] = i
		info["can_select"] = false
		candidates = append(candidates, info)
	}
	ctx.Engine.SetPendingActionWithData(ctx.PlayerID, "council_guard_reorder",
		"议庭护法:查看牌库顶5张；有印记时调整顺序并放回牌库顶，否则洗混牌库", candidates, 0, 0,
		func(selected []string, data map[string]any) {
			if !hasMark {
				ctx.Engine.shuffleDeck(ctx.OpponentID)
				return
			}
			top := append(stringsFromActionData(data, "top_order"), stringsFromActionData(data, "bottom_order")...)
			order := make([]any, len(top))
			for i, id := range top {
				order[i] = id
			}
			resolveTopDeckReorder(ctx.Engine, ctx.PlayerID, ctx.OpponentID, looked, map[string]any{"top_order": order}, "council_guard_reorder")
		})
	return nil
}

var _ PerTurnAbility = Card1521112CouncilGuard{}
