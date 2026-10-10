package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card2421112AutumnMapleGem struct{ AlwaysActive }

func (Card2421112AutumnMapleGem) ID() string { return "2421112" }

func (Card2421112AutumnMapleGem) Name() string { return "秋枫宝钻" }

func (Card2421112AutumnMapleGem) OnEnter(ctx *EffectContext) error {
	ctx.Source.Statuses[autumnMapleGemCounter] += 2
	return nil
}

func (Card2421112AutumnMapleGem) PerTurnLabel(*CardInstance) string {
	return "回合技"
}

func (Card2421112AutumnMapleGem) OnPerTurn(ctx *EffectContext) error {
	if ctx.Source.Statuses[autumnMapleGemCounter] <= 0 {
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementEarth && card.IsHorizontal
	})
	if len(candidates) == 0 {
		return nil
	}
	allowed := make(map[string]bool, len(candidates))
	for _, candidate := range candidates {
		if id, _ := candidate["instance_id"].(string); id != "" {
			allowed[id] = true
		}
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "autumn_maple_gem_reset",
		"秋枫宝钻:选择1个地脉伙伴重置", candidates, 1, 1,
		func(selected []string) {
			id := firstSelected(selected)
			if !allowed[id] || ctx.Source.Statuses[autumnMapleGemCounter] <= 0 {
				return
			}
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, id)
			if target == nil || zone != "unit" || target.Card == nil || !target.Card.IsCompanion() || target.Card.Category != model.ElementEarth {
				return
			}
			ctx.Source.Statuses[autumnMapleGemCounter]--
			target.IsHorizontal = false
		})
	return nil
}

func (Card2421112AutumnMapleGem) ValidateAbility(ctx *EffectContext, trigger EffectTrigger) error {
	if trigger != TriggerPerTurn {
		return nil
	}
	if ctx.Source.Statuses[autumnMapleGemCounter] <= 0 {
		return fmt.Errorf("秋枫宝钻没有标记物")
	}
	if !ctx.Engine.hasResettableEarthCompanion(ctx.PlayerID) {
		return fmt.Errorf("秋枫宝钻需要1个已横置的地脉伙伴")
	}
	return nil
}

const autumnMapleGemCounter = "秋枫宝钻标记物"

func (e *Engine) hasResettableEarthCompanion(playerID int) bool {
	return len(e.friendlyUnits(playerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() &&
			card.Card.Category == model.ElementEarth &&
			card.IsHorizontal
	})) > 0
}
