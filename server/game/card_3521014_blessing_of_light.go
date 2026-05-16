package game

import "eraofarcane/model"

type Card3521014BlessingOfLight struct{}

func (Card3521014BlessingOfLight) ID() string   { return "3521014" }
func (Card3521014BlessingOfLight) Name() string { return "光之祝福" }

func (Card3521014BlessingOfLight) OnSpellCast(ctx *EffectContext) error {
	if !isFriendlySpellCast(ctx) || !isSpellBeingCast(ctx) {
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "blessing_of_light",
		"选择1个友方伙伴获得+1血和负载+1光", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, selected[0])
			if target == nil || zone != "unit" {
				return
			}
			target.CurrentLife++
			addElementsGainBonus(target, model.ElementLight, 1)
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source":  cardToInfo(ctx.Source),
				"target":  cardToInfo(target),
				"effect":  "life_and_load",
				"life":    1,
				"element": model.ElementLight,
				"amount":  1,
			}})
		})
	return nil
}
