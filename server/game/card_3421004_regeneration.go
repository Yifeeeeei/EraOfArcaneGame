package game

import "eraofarcane/model"

type Card3421004Regeneration struct{}

func (Card3421004Regeneration) ID() string   { return "3421004" }
func (Card3421004Regeneration) Name() string { return "再生之力" }
func (Card3421004Regeneration) NeedsSpellTarget() bool {
	return false
}
func (Card3421004Regeneration) OnSpellHit(ctx *EffectContext) error {
	targets := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card.Card.IsCompanion() && card.Card.Category == model.ElementEarth
	})
	if len(targets) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "regeneration_reset",
		"选择1张地脉伙伴重置", targets, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			card, _ := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, selected[0])
			resetInstance(card)
		})
	return nil
}
