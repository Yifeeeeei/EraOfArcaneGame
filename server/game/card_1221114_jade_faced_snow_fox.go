package game

import (
	"eraofarcane/model"
)

type Card1221114JadeFacedSnowFox struct{ AlwaysActive }

func (Card1221114JadeFacedSnowFox) ID() string { return "1221114" }

func (Card1221114JadeFacedSnowFox) Name() string { return "玉面雪狐" }

func (Card1221114JadeFacedSnowFox) HasActiveSpellReaction(card *CardInstance) bool {
	return card != nil && card.Position != nil && !card.UltimateUsed
}

func (Card1221114JadeFacedSnowFox) CanReactToSpell(ctx *EffectContext, spell *SpellCast) bool {
	return ctx != nil &&
		ctx.Engine != nil &&
		ctx.Source != nil &&
		ctx.Source.Position != nil &&
		!ctx.Source.UltimateUsed &&
		spell != nil &&
		spell.AttackerID != ctx.PlayerID &&
		spell.Skill != nil &&
		spell.Skill.Card != nil &&
		canUseSkillForPurpose(spell.Skill.Card, skillPurposeAttack)
}

func (Card1221114JadeFacedSnowFox) OnSpellReaction(ctx *EffectContext, spell *SpellCast) error {
	if !(Card1221114JadeFacedSnowFox{}).CanReactToSpell(ctx, spell) {
		return nil
	}
	sourceID := ctx.Source.InstanceID
	positions := ctx.Engine.allUnitPositionsForPlayer(ctx.PlayerID, ctx.PlayerID)
	ctx.Engine.SetPendingAction(ctx.PlayerID, "jade_faced_snow_fox_move",
		"玉面雪狐:移动此卡并获得2水", positions, 1, 1,
		func(selected []string) {
			if ctx.Engine.State.PendingSpell == nil {
				return
			}
			source, _ := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, sourceID)
			if source == nil || source.Position == nil || source.UltimateUsed {
				return
			}
			pos, ok := positionFromSelectionID(firstSelected(selected))
			if !ok {
				return
			}
			ctx.Engine.moveOrSwapUnitToPosition(ctx.PlayerID, sourceID, pos)
			ctx.Engine.State.Players[ctx.PlayerID].Elements[model.ElementWater] += 2
			source.UltimateUsed = true
			ctx.Engine.emit(GameEvent{
				Type:   "jade_faced_snow_fox_move",
				Player: -1,
				Data: map[string]any{
					"player": ctx.PlayerID,
					"source": cardToInfo(source),
					"water":  2,
				},
			})
			ctx.Engine.promptSpellRetarget(ctx.PlayerID, source, ctx.ExtraData,
				"jade_faced_snow_fox_retarget",
				"玉面雪狐:重新选择法术攻击目标",
				"jade_faced_snow_fox")
		})
	return nil
}

var _ SpellReactionBehavior = Card1221114JadeFacedSnowFox{}
