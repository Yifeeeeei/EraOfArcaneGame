package game

import "eraofarcane/model"

type Card2521012IridescentPaint struct{}

func (Card2521012IridescentPaint) ID() string   { return "2521012" }
func (Card2521012IridescentPaint) Name() string { return "幻彩颜料" }
func (Card2521012IridescentPaint) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return effectiveElementsGain(card)[model.ElementLight] > 0
	})
	candidates = append(candidates, ctx.Engine.friendlyEquipment(ctx.PlayerID, func(card *CardInstance) bool {
		return effectiveElementsGain(card)[model.ElementLight] > 0
	})...)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "iridescent_paint",
		"选择要转换光辉负载的友方卡牌(最多转换4点)", candidates, 1, 4,
		func(selected []string) {
			remaining := 4
			for _, id := range selected {
				if remaining <= 0 {
					break
				}
				card, _ := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, id)
				if card == nil {
					continue
				}
				remaining -= convertLoad(card, model.ElementLight, model.ElementArcane, remaining)
			}
		})
	return nil
}
