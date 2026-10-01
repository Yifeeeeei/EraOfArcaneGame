package game

import (
	"fmt"
)

type Card2011102ForesightOrb struct{ AlwaysActive }

func (Card2011102ForesightOrb) ID() string { return "2011102" }

func (Card2011102ForesightOrb) Name() string { return "预知宝珠" }

func (Card2011102ForesightOrb) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	if ctx.Source.IsHorizontal {
		return fmt.Errorf("预知宝珠需要竖置才能消耗")
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	lookCount := min(3, len(ps.Deck))
	if lookCount == 0 {
		return fmt.Errorf("预知宝珠需要牌库中有牌")
	}
	ctx.Source.IsHorizontal = true
	looked := append([]*CardInstance(nil), ps.Deck[:lookCount]...)
	candidates := make([]map[string]any, 0, len(looked))
	for i, card := range looked {
		info := candidateInfo(card, "deck", "own")
		info["deck_index"] = i
		info["can_select"] = false
		candidates = append(candidates, info)
	}
	ctx.Engine.SetPendingActionWithData(ctx.PlayerID, "foresight_orb_reorder",
		"预知宝珠:将牌库顶3张以任意顺序放回牌库顶或牌库底", candidates, 0, 0,
		func(selected []string, data map[string]any) {
			resolveTopDeckReorder(ctx.Engine, ctx.PlayerID, ctx.PlayerID, looked, data, "foresight_orb_reorder")
		})
	return nil
}

var _ PerTurnAbility = Card2011102ForesightOrb{}
