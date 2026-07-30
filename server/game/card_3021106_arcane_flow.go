package game

import "eraofarcane/model"

type Card3021106ArcaneFlow struct{ AlwaysActive }

func (Card3021106ArcaneFlow) ID() string   { return "3021106" }
func (Card3021106ArcaneFlow) Name() string { return "奥能流贯" }

func (Card3021106ArcaneFlow) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if ctx == nil || ctx.Source == nil || ctx.Target == nil || ctx.Source != ctx.Target {
		return
	}
	if !friendlyFieldHasOnlyArcaneCards(ctx.Engine, ctx.PlayerID) {
		return
	}
	stats.PowerBonus += max(ctx.Source.Card.Power+ctx.Source.PowerBonus, 0)
}

func friendlyFieldHasOnlyArcaneCards(engine *Engine, playerID int) bool {
	if engine == nil || playerID < 0 || playerID >= len(engine.State.Players) {
		return false
	}
	for _, card := range engine.getAllFieldCards(engine.State.Players[playerID]) {
		if card == nil || card.Card == nil {
			continue
		}
		if card.Card.Category != model.ElementArcane {
			return false
		}
	}
	return true
}

var _ SpellStatModifier = Card3021106ArcaneFlow{}
