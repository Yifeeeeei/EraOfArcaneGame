package game

import (
	"eraofarcane/model"
)

type Card1621113WhisperElfPriest struct{ AlwaysActive }

func (Card1621113WhisperElfPriest) ID() string { return "1621113" }

func (Card1621113WhisperElfPriest) Name() string { return "谧语精灵祭司" }

func (Card1621113WhisperElfPriest) OnDeath(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion()
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "whisper_elf_priest_load",
		"谧语精灵祭司:选择1个友方伙伴负载+1暗", candidates, 1, 1,
		func(selected []string) {
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(selected))
			if target != nil && zone == "unit" && target.Card != nil && target.Card.IsCompanion() {
				ctx.Engine.addElementsGainBonus(target, ctx.PlayerID, model.ElementShadow, 1, ctx.Source)
			}
		})
	return nil
}
