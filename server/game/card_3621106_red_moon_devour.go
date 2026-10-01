package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card3621106RedMoonDevour struct{ AlwaysActive }

func (Card3621106RedMoonDevour) ID() string { return "3621106" }

func (Card3621106RedMoonDevour) Name() string { return "红月吞噬" }

func (Card3621106RedMoonDevour) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || !isOwnSpellHit(ctx) {
		return nil
	}
	targets := spellHitAffectedUnitsFromData(ctx)
	if len(targets) == 0 && ctx.Target != nil {
		targets = []*CardInstance{ctx.Target}
	}
	destroyed := make([]*CardInstance, 0, len(targets))
	seen := make(map[string]bool, len(targets))
	totalRemainingLife := 0
	for _, target := range targets {
		if target == nil || target.Card == nil || !target.Card.IsCompanion() || target.OwnerID == ctx.PlayerID || seen[target.InstanceID] {
			continue
		}
		if !ctx.Engine.unitInOwnerGrid(target, target.OwnerID) {
			continue
		}
		seen[target.InstanceID] = true
		totalRemainingLife += max(target.CurrentLife, 0)
		destroyed = append(destroyed, target)
	}
	if len(destroyed) == 0 {
		return nil
	}
	for _, target := range destroyed {
		ctx.Engine.destroyUnitWithData(target, target.OwnerID, map[string]any{
			"destroyed_by": ctx.Source.InstanceID,
			"attacker":     ctx.PlayerID,
		})
	}
	if totalRemainingLife <= 0 || !ctx.Engine.redMoonActive(ctx.PlayerID) {
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.Category == model.ElementShadow
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "red_moon_devour_life",
		fmt.Sprintf("红月吞噬:选择1个友方暗影单位获得+%d血", totalRemainingLife), candidates, 1, 1,
		func(selected []string) {
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target == nil || zone != "unit" || target.Card == nil || target.Card.Category != model.ElementShadow {
				return
			}
			ctx.Engine.gainLife(target, totalRemainingLife, ctx.Source)
			ctx.Engine.emit(GameEvent{
				Type:   "effect_trigger",
				Player: ctx.PlayerID,
				Data: map[string]any{
					"source": cardToInfo(ctx.Source),
					"target": cardToInfo(target),
					"effect": "modify_life",
					"amount": totalRemainingLife,
				},
			})
		})
	return nil
}

var _ OnSpellHitBehavior = Card3621106RedMoonDevour{}
