package game

import (
	"fmt"
)

type Card2421109CaveElfPickaxe struct{ AlwaysActive }

func (Card2421109CaveElfPickaxe) ID() string { return "2421109" }

func (Card2421109CaveElfPickaxe) Name() string { return "地穴精灵矿镐" }

func (Card2421109CaveElfPickaxe) PerTurnLabel(*CardInstance) string {
	return "消耗"
}

func (Card2421109CaveElfPickaxe) OnPerTurn(ctx *EffectContext) error {
	if !ctx.Engine.canConsumeCard(ctx.Source) {
		return fmt.Errorf("地穴精灵矿镐不能被消耗")
	}
	if !ctx.Engine.equipmentInOwnerSlot(ctx.PlayerID, ctx.Source) {
		return fmt.Errorf("地穴精灵矿镐必须从装备区发动")
	}
	choices := []map[string]any{
		{"instance_id": "companion", "number": "2421109", "name": "伙伴", "type": "选择", "zone": "choice", "side": "own"},
		{"instance_id": "item", "number": "2421109", "name": "道具", "type": "选择", "zone": "choice", "side": "own"},
	}
	ctx.Source.IsHorizontal = true
	ctx.Engine.SetPendingAction(ctx.PlayerID, "cave_elf_pickaxe_kind",
		"地穴精灵矿镐:选择翻取伙伴或道具", choices, 1, 1,
		func(selected []string) {
			switch firstSelected(selected) {
			case "companion":
				ctx.Engine.flipDeckMatchesToHand(ctx.PlayerID, 1, 5, func(card *CardInstance) bool {
					return card != nil && card.Card != nil && card.Card.IsCompanion()
				})
			case "item":
				ctx.Engine.flipDeckMatchesToHand(ctx.PlayerID, 1, 5, func(card *CardInstance) bool {
					return card != nil && card.Card != nil && card.Card.IsItem()
				})
			}
		})
	return nil
}
