package game

import (
	"eraofarcane/model"
)

type Card3121110CursedFire struct{ AlwaysActive }

func (Card3121110CursedFire) ID() string { return "3121110" }

func (Card3121110CursedFire) Name() string { return "咒火" }

func (Card3121110CursedFire) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.Source.Card.Number != "3121110" || ctx.PlayerID < 0 {
		return nil
	}
	drawn := ctx.Engine.flipDeckMatchesToHandThen(ctx.PlayerID, 1, 0, isCheapFireSpellScroll, func(drawn []*CardInstance) {
		for _, card := range drawn {
			makeEntryCostZero(card)
		}
	})
	if len(drawn) > 0 {
		ctx.Engine.emit(GameEvent{
			Type:   "cursed_fire_flip_scroll",
			Player: ctx.PlayerID,
			Data: map[string]any{
				"player": ctx.PlayerID,
				"source": cardToInfo(ctx.Source),
				"card":   cardToInfo(drawn[0]),
			},
		})
		ctx.Engine.promptCursedFireImmediateScrollUse(ctx.PlayerID, ctx.Source, drawn[0])
	}
	return nil
}

var _ OnSpellHitBehavior = Card3121110CursedFire{}

func (e *Engine) promptCursedFireImmediateScrollUse(playerID int, source *CardInstance, scroll *CardInstance) {
	if e == nil || source == nil || scroll == nil || scroll.Card == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	if isDefenseOnlySkill(scroll.Card) || !canUseSkillForPurpose(scroll.Card, skillPurposeAttack) {
		return
	}
	if !skillNeedsTargetInstance(scroll) {
		e.SetPendingAction(playerID, "cursed_fire_use_scroll",
			"咒火:是否立刻使用翻取的卷轴", []map[string]any{candidateInfo(scroll, "hand", "own")}, 0, 1,
			func(selected []string) {
				if len(selected) == 0 {
					return
				}
				e.useCursedFireScrollFromHand(playerID, scroll, SpellTarget{Type: "none"})
			})
		return
	}
	candidates := e.enemyUnits(playerID, true, func(target *CardInstance) bool {
		if target == nil || target.Position == nil {
			return false
		}
		return e.validateSpellTarget(playerID, scroll, SpellTarget{Type: "unit", Position: *target.Position}) == nil
	})
	if len(candidates) == 0 {
		return
	}
	e.SetPendingAction(playerID, "cursed_fire_use_scroll_target",
		"咒火:选择目标并立刻使用翻取的卷轴", candidates, 0, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			target := e.findUnitByInstanceID(firstSelected(selected))
			if target == nil || target.Position == nil {
				return
			}
			targetOwner := target.OwnerID
			e.useCursedFireScrollFromHand(playerID, scroll, SpellTarget{Type: "unit", Position: *target.Position, OwnerID: &targetOwner})
		})
}

func (e *Engine) useCursedFireScrollFromHand(playerID int, scroll *CardInstance, target SpellTarget) {
	ps := e.State.Players[playerID]
	current, handIdx := ps.FindHandCard(scroll.InstanceID)
	if current != scroll || handIdx < 0 {
		return
	}
	data := map[string]any{"instance_id": scroll.InstanceID, "target_type": target.Type}
	if target.Type == "unit" {
		data["target_col"] = float64(target.Position.Col)
		data["target_row"] = float64(target.Position.Row)
		if target.OwnerID != nil {
			data["target_owner"] = float64(*target.OwnerID)
		}
	}
	_ = e.handleUseSpellScrollItem(playerID, ActionMessage{Data: data}, scroll, handIdx)
}

func isCheapFireSpellScroll(card *CardInstance) bool {
	return card != nil &&
		card.Card != nil &&
		card.Card.Category == model.ElementFire &&
		isSpellScrollCard(card.Card) &&
		totalElementCost(card.Card.ElementsCost) < 4
}
