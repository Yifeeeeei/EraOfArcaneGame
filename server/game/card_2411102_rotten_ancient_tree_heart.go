package game

import (
	"fmt"
	"strings"
)

type Card2411102RottenAncientTreeHeart struct{ AlwaysActive }

func (Card2411102RottenAncientTreeHeart) ID() string { return "2411102" }

func (Card2411102RottenAncientTreeHeart) Name() string { return "腐朽的古树之心" }

func (Card2411102RottenAncientTreeHeart) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil {
		return nil
	}
	castPlayer, ok := ctx.ExtraData["cast_player"].(int)
	if !ok || castPlayer < 0 || castPlayer >= len(ctx.Engine.State.Players) {
		return nil
	}
	key := fmt.Sprintf("%s%d", rottenAncientTreeHeartSpellCountPrefix, castPlayer)
	ctx.Source.Statuses[key]++
	if ctx.Source.Statuses[key] < 2 {
		return nil
	}
	candidates := loadRemovalCandidates(ctx.Engine, castPlayer, "own")
	if len(candidates) == 0 {
		return nil
	}
	removeLoad := func(selection string) {
		targetID, elem, ok := strings.Cut(selection, "|")
		if !ok || elem == "" {
			return
		}
		target := ctx.Engine.findCardOnField(ctx.Engine.State.Players[castPlayer], targetID)
		if target == nil || reducibleElementLoad(target, elem) <= 0 {
			return
		}
		ctx.Engine.reduceCardElementLoadWithTriggers(castPlayer, target, elem, 1, ctx.Source)
		ctx.Source.Statuses[key] -= 2
		if ctx.Source.Statuses[key] < 0 {
			ctx.Source.Statuses[key] = 0
		}
		ctx.Engine.emit(GameEvent{
			Type:   "rotten_ancient_tree_heart_remove_load",
			Player: -1,
			Data: map[string]any{
				"player":  castPlayer,
				"source":  cardToInfo(ctx.Source),
				"target":  cardToInfo(target),
				"element": elem,
			},
		})
		if target == ctx.Source && ctx.Engine.totalLoad(ctx.Source) <= 0 {
			ctx.Engine.sacrificeEquipment(ctx.Source.OwnerID, ctx.Source.InstanceID)
		}
	}
	if len(candidates) == 1 {
		id, _ := candidates[0]["instance_id"].(string)
		removeLoad(id)
		return nil
	}
	ctx.Engine.SetPendingAction(castPlayer, "rotten_ancient_tree_heart_remove_load",
		"腐朽的古树之心:选择自己场上1点负载移除", candidates, 1, 1,
		func(selected []string) {
			removeLoad(firstSelected(selected))
		})
	return nil
}

var _ OnSpellCastBehavior = Card2411102RottenAncientTreeHeart{}

const rottenAncientTreeHeartSpellCountPrefix = "rotten_ancient_tree_heart_spell_count_p"

func loadRemovalCandidates(e *Engine, playerID int, side string) []map[string]any {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	candidates := make([]map[string]any, 0)
	for _, card := range e.getAllFieldCards(e.State.Players[playerID]) {
		if card == nil || card.Card == nil || e.hasEffectiveStatus(card, StatusPetrify) {
			continue
		}
		for elem, amount := range dragonBloodTreantReducibleLoad(card) {
			if amount <= 0 {
				continue
			}
			info := candidateInfo(card, "field", side)
			info["instance_id"] = card.InstanceID + "|" + elem
			info["load_element"] = elem
			info["name"] = fmt.Sprintf("%s - 移除%s负载", card.Card.Name, elem)
			candidates = append(candidates, info)
		}
	}
	return candidates
}
