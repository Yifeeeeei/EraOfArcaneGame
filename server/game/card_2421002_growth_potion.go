package game

import "eraofarcane/model"

type Card2421002GrowthPotion struct{ AlwaysActive }

func (Card2421002GrowthPotion) ID() string   { return "2421002" }
func (Card2421002GrowthPotion) Name() string { return "生长药水" }
func (Card2421002GrowthPotion) OnUseItem(ctx *EffectContext) error {
	targets := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card.Card.IsCompanion() && card.Card.Category == model.ElementEarth
	})
	if len(targets) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "growth_potion_reset",
		"生长药水:选择你的1个地脉伙伴重置", targets, 1, 1,
		func(selected []string) {
			target := selectedUnitFromCandidates(ctx.Engine, selected, targets)
			if target != nil {
				resetInstance(target)
			}
		})
	return nil
}
