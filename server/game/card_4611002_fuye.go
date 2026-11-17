package game

import (
	"fmt"
)

const fuyeDeathAfterExertStatus = "fuye_death_after_exert"

type Card4611002Fuye struct{ AlwaysActive }

func (Card4611002Fuye) ID() string { return "4611002" }

func (Card4611002Fuye) Name() string { return "芙雅夫人" }

func (Card4611002Fuye) OnUltimate(ctx *EffectContext) error {
	if ctx.Target != nil {
		if !isValidFuyeUltimateTarget(ctx, ctx.Target) {
			return fmt.Errorf("invalid Fuye target")
		}
		applyFuyeUltimate(ctx.Target)
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() && !card.IsHorizontal
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "fuye_ultimate_target",
		"芙雅夫人:选择1个友方伙伴，其攻击力和负载翻倍，并在消耗或透支后死亡", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, selected[0])
			if zone != "unit" || !isValidFuyeUltimateTarget(ctx, target) {
				return
			}
			applyFuyeUltimate(target)
		})
	return nil
}

func isValidFuyeUltimateTarget(ctx *EffectContext, target *CardInstance) bool {
	if ctx == nil || ctx.Engine == nil || target == nil || target.Card == nil {
		return false
	}
	if target.OwnerID != ctx.PlayerID || !target.Card.IsCompanion() || target.IsHorizontal {
		return false
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	return ps != nil && ctx.Engine.findUnitOnGrid(ps, target.InstanceID) == target
}

func applyFuyeUltimate(target *CardInstance) {
	target.CurrentAttack *= 2
	gain := effectiveElementsGain(target)
	doubled := make(map[string]int, len(gain))
	for elem, amount := range gain {
		doubled[elem] = amount * 2
	}
	setElementsGain(target, doubled)
	target.Statuses[fuyeDeathAfterExertStatus] = 1
}

func (e *Engine) destroyFuyeDoomedAfterExert(cards []*CardInstance) {
	for _, card := range cards {
		e.destroyFuyeDoomedCardAfterExert(card)
	}
}

func (e *Engine) destroyFuyeDoomedCardAfterExert(card *CardInstance) {
	if card == nil || card.Statuses[fuyeDeathAfterExertStatus] <= 0 {
		return
	}
	ps := e.State.Players[card.OwnerID]
	if ps == nil || e.findUnitOnGrid(ps, card.InstanceID) == nil {
		return
	}
	delete(card.Statuses, fuyeDeathAfterExertStatus)
	e.destroyUnitWithCause(card, card.OwnerID, "fuye_exert")
}
