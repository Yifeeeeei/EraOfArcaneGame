package game

import "eraofarcane/model"

type Card2421002GrowthPotion struct{}

func (Card2421002GrowthPotion) ID() string   { return "2421002" }
func (Card2421002GrowthPotion) Name() string { return "生长药水" }
func (Card2421002GrowthPotion) OnUseItem(ctx *EffectContext) error {
	targets := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card.Card.IsCompanion() && card.Card.Category == model.ElementEarth
	})
	if target := firstUnitFromCandidates(ctx.Engine, ctx.PlayerID, targets); target != nil {
		resetInstance(target)
	}
	return nil
}
