package game

import (
	"eraofarcane/model"
)

type Card2521101BlessedLoneStar struct{ AlwaysActive }

func (Card2521101BlessedLoneStar) ID() string { return "2521101" }

func (Card2521101BlessedLoneStar) Name() string { return "赐福之孤星" }

func (Card2521101BlessedLoneStar) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion()
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "blessed_lone_star_target",
		"赐福之孤星:选择1个友方伙伴获得负载+1光和+1血", candidates, 1, 1,
		func(selected []string) {
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target == nil || zone != "unit" || target.Card == nil || !target.Card.IsCompanion() {
				return
			}
			ctx.Engine.gainLife(target, 1, ctx.Source)
			ctx.Engine.addElementsGainBonus(target, ctx.PlayerID, model.ElementLight, 1, ctx.Source)
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
