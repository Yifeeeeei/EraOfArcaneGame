package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card3601101BloodFeast struct{ AlwaysActive }

func (Card3601101BloodFeast) ID() string { return "3601101" }

func (Card3601101BloodFeast) Name() string { return "鲜血盛宴" }

func (Card3601101BloodFeast) AllowsFriendlySpellTarget() bool {
	return true
}

func (Card3601101BloodFeast) ValidateSpellTarget(ctx *EffectContext, target SpellTarget, targetUnit *CardInstance) error {
	if ctx == nil || target.Type != "unit" || targetUnit == nil || targetUnit.OwnerID != ctx.PlayerID {
		return fmt.Errorf("鲜血盛宴只能攻击友方单位")
	}
	return nil
}

func (Card3601101BloodFeast) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil || ctx.Source.Card.Number != "3601101" || !isOwnSpellHit(ctx) {
		return nil
	}
	choices := []map[string]any{
		{"instance_id": "gain_shadow", "name": "获得2暗", "zone": "choice", "side": "own"},
		{"instance_id": "heal_hero", "name": "人物回复1血", "zone": "choice", "side": "own"},
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "blood_feast_reward",
		"鲜血盛宴:选择命中奖励", choices, 1, 1,
		func(selected []string) {
			switch firstSelected(selected) {
			case "gain_shadow":
				ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementShadow: 2})
			case "heal_hero":
				hero := ctx.Engine.playerHeroCard(ctx.PlayerID)
				if hero != nil && hero.CurrentLife < maxLife(hero) {
					ctx.Engine.healUnit(hero, 1, ctx.Source)
				}
			}
		})
	return nil
}

func (Card3601101BloodFeast) PerTurnLabel(*CardInstance) string { return "绑定" }

func (Card3601101BloodFeast) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	cost := map[string]int{model.ElementShadow: 1}
	if !ctx.Engine.canPayCost(ps, cost) {
		return fmt.Errorf("鲜血盛宴需要支付1暗绑定到人物")
	}
	hero := ctx.Engine.playerHeroCard(ctx.PlayerID)
	if hero == nil {
		return fmt.Errorf("no hero")
	}
	slot := -1
	for i, skill := range ps.Skills {
		if skill == ctx.Source {
			slot = i
			break
		}
	}
	if slot < 0 {
		return fmt.Errorf("skill is already bound")
	}
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "blood_feast_bind", "鲜血盛宴:支付1暗绑定到人物", []map[string]any{candidateInfo(hero, "unit", "own")}, 1, 1, cost, false, func(selected []string, data map[string]any) error {
		if ps.Skills[slot] != ctx.Source || ctx.Engine.playerHeroCard(ctx.PlayerID) != hero {
			return fmt.Errorf("invalid binding target")
		}
		if !ctx.Engine.payCostForAction(ps, cost, ActionMessage{Data: data}) {
			return fmt.Errorf("invalid payment")
		}
		ps.Skills[slot] = nil
		ctx.Source.SlotIndex = -1
		markTransferredBoundSkill(ctx.Source)
		hero.BoundSkills = append(hero.BoundSkills, ctx.Source)
		ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{"source": cardToInfo(ctx.Source), "target": cardToInfo(hero), "effect": "bind_existing_skill"}})
		return nil
	})
	return nil
}

func (e *Engine) playerHeroCard(playerID int) *CardInstance {
	ps := e.State.Players[playerID]
	if ps == nil {
		return nil
	}
	if ps.Hero != nil {
		return ps.Hero
	}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			card := ps.Units[col][row]
			if card != nil && card.Card != nil && card.Card.IsHero() {
				return card
			}
		}
	}
	return nil
}
