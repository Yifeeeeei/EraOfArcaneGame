package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card2321102WindCycle struct{ AlwaysActive }

func (Card2321102WindCycle) ID() string { return "2321102" }

func (Card2321102WindCycle) Name() string { return "风之轮回" }

func (Card2321102WindCycle) PerTurnLabel(*CardInstance) string {
	return "主动"
}

func (Card2321102WindCycle) OnPerTurn(ctx *EffectContext) error {
	if !ctx.Engine.canConsumeCard(ctx.Source) {
		return fmt.Errorf("风之轮回不能被消耗")
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	inEquipment := false
	for _, card := range ps.Equipment {
		if card == ctx.Source {
			inEquipment = true
			break
		}
	}
	if !inEquipment {
		return fmt.Errorf("风之轮回必须从装备区献祭")
	}
	candidates := make([]map[string]any, 0)
	allowed := make(map[string]bool)
	for _, card := range ps.Graveyard {
		if card != nil && card.Card != nil && card.Card.Category == model.ElementAir {
			candidates = append(candidates, candidateInfo(card, "graveyard", "own"))
			allowed[card.InstanceID] = true
		}
	}
	ctx.Source.IsHorizontal = true
	if !ctx.Engine.sacrificeEquipment(ctx.PlayerID, ctx.Source.InstanceID) {
		return fmt.Errorf("风之轮回必须从装备区献祭")
	}
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "wind_cycle_shuffle_air",
		"风之轮回:选择任意数量的大气弃牌洗回卡组", candidates, 0, len(candidates),
		func(selected []string) {
			selectedSet := make(map[string]bool, len(selected))
			for _, id := range selected {
				if allowed[id] {
					selectedSet[id] = true
				}
			}
			if len(selectedSet) == 0 {
				return
			}
			for i := 0; i < len(ps.Graveyard); {
				card := ps.Graveyard[i]
				if card != nil && selectedSet[card.InstanceID] && card.Card != nil && card.Card.Category == model.ElementAir {
					ps.Graveyard = append(ps.Graveyard[:i], ps.Graveyard[i+1:]...)
					ps.Deck = append(ps.Deck, card)
					continue
				}
				i++
			}
			ctx.Engine.shuffleDeck(ctx.PlayerID)
		})
	return nil
}
