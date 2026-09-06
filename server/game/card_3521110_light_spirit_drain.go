package game

import (
	"eraofarcane/model"
)

type Card3521110LightSpiritDrain struct{ AlwaysActive }

func (Card3521110LightSpiritDrain) ID() string { return "3521110" }

func (Card3521110LightSpiritDrain) Name() string { return "光灵汲取" }

func (Card3521110LightSpiritDrain) OnSpellHit(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() && card.Card.Category == model.ElementLight
	})
	if len(candidates) == 0 {
		return nil
	}
	applyLoad := func(instanceID string) {
		target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, instanceID)
		if target == nil || zone != "unit" || target.Card == nil || !target.Card.IsCompanion() || target.Card.Category != model.ElementLight {
			return
		}
		ctx.Engine.addElementsGainBonus(target, ctx.PlayerID, model.ElementLight, 1, ctx.Source)
	}
	if len(candidates) == 1 {
		id, _ := candidates[0]["instance_id"].(string)
		applyLoad(id)
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "light_spirit_drain_load",
		"光灵汲取:选择1个友方光辉伙伴获得负载+1光", candidates, 1, 1,
		func(selected []string) {
			applyLoad(firstSelected(selected))
		})
	return nil
}
