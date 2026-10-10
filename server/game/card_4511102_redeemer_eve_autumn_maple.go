package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card4511102RedeemerEveAutumnMaple struct{ AlwaysActive }

func (Card4511102RedeemerEveAutumnMaple) ID() string { return "4511102" }

func (Card4511102RedeemerEveAutumnMaple) Name() string { return "救赎者 伊芙 秋枫" }

func (Card4511102RedeemerEveAutumnMaple) OnUltimate(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.UltimateUsed {
		return nil
	}
	ownUnits := ctx.Engine.friendlyUnits(ctx.PlayerID, true, nil)
	enemyUnits := ctx.Engine.enemyUnits(ctx.PlayerID, true, nil)
	if len(enemyUnits) <= len(ownUnits) {
		return fmt.Errorf("敌方场上单位数量必须比我方多")
	}
	woundedCount := 0
	forEachPlayerUnit(ctx.Engine.State.Players[ctx.PlayerID], true, func(unit *CardInstance) {
		if unit != nil && unit.CurrentLife < maxLife(unit) {
			woundedCount++
		}
	})
	if woundedCount <= 0 {
		return fmt.Errorf("救赎者 伊芙 秋枫需要受伤的友方单位")
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ps.Elements[model.ElementLight] < woundedCount {
		return fmt.Errorf("救赎者 伊芙 秋枫需要%d点光元素", woundedCount)
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() && !card.Card.IsHero() && card.CurrentLife < maxLife(card)
	})
	if len(candidates) == 0 {
		return fmt.Errorf("救赎者 伊芙 秋枫需要受伤的友方伙伴")
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "redeemer_eve_autumn_maple_target",
		fmt.Sprintf("救赎者 伊芙 秋枫:支付%d光,选择1个受伤友方伙伴获得+%d血和负载+%d光", woundedCount, woundedCount, woundedCount),
		candidates, 1, 1,
		func(selected []string) {
			if ctx.Source == nil || !ctx.Engine.cardStillOnField(ctx.Source) {
				return
			}
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target == nil || zone != "unit" || target.Card == nil || !target.Card.IsCompanion() || target.Card.IsHero() || target.CurrentLife >= maxLife(target) {
				return
			}
			ownNow := ctx.Engine.friendlyUnits(ctx.PlayerID, true, nil)
			enemyNow := ctx.Engine.enemyUnits(ctx.PlayerID, true, nil)
			if len(enemyNow) <= len(ownNow) {
				return
			}
			x := 0
			forEachPlayerUnit(ctx.Engine.State.Players[ctx.PlayerID], true, func(unit *CardInstance) {
				if unit != nil && unit.CurrentLife < maxLife(unit) {
					x++
				}
			})
			if x <= 0 || ps.Elements[model.ElementLight] < x {
				return
			}
			ps.Elements[model.ElementLight] -= x
			target.Statuses["max_life_bonus"] += x
			ctx.Engine.gainLife(target, x, ctx.Source)
			ctx.Engine.addElementsGainBonus(target, ctx.PlayerID, model.ElementLight, x, ctx.Source)
			ctx.Source.UltimateUsed = true
			ctx.Engine.emit(GameEvent{
				Type:   "redeemer_eve_autumn_maple_blessing",
				Player: -1,
				Data: map[string]any{
					"player": ctx.PlayerID,
					"source": cardToInfo(ctx.Source),
					"target": cardToInfo(target),
					"amount": x,
				},
			})
		})
	return nil
}

var _ UltimateAbility = Card4511102RedeemerEveAutumnMaple{}

func forEachPlayerUnit(ps *PlayerState, includeHero bool, fn func(*CardInstance)) {
	if ps == nil || fn == nil {
		return
	}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := ps.Units[col][row]
			if unit == nil {
				continue
			}
			if !includeHero && unit.Card != nil && unit.Card.IsHero() {
				continue
			}
			fn(unit)
		}
	}
}
