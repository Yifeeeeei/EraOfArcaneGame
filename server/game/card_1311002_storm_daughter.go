package game

import (
	"eraofarcane/model"
)

type Card1311002StormDaughter struct{ AlwaysActive }

func (Card1311002StormDaughter) ID() string { return "1311002" }

func (Card1311002StormDaughter) Name() string { return `"风暴之女" 艾拉雅` }

func (Card1311002StormDaughter) OnEnter(ctx *EffectContext) error {
	bindSkillToHost(ctx, "3301001")
	return nil
}

func (Card1311002StormDaughter) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if ctx.Target == nil || ctx.Target.Card == nil || ctx.Target.Card.Category != model.ElementAir {
		return
	}
	if len(ctx.Engine.State.Players[ctx.PlayerID].Hand) >= ctx.Engine.State.HandLimit {
		stats.PowerBonus += len(ctx.Engine.State.Players[ctx.PlayerID].Hand)
	}
}
