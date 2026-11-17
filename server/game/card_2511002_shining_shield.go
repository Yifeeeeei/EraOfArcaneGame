package game

type Card2511002ShiningShield struct{ AlwaysActive }

func (Card2511002ShiningShield) ID() string { return "2511002" }

func (Card2511002ShiningShield) Name() string { return "辉之盾 闪耀" }

func (Card2511002ShiningShield) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if ctx.ExtraData != nil && ctx.ExtraData["purpose"] == string(skillPurposeDefend) {
		stats.PowerBonus += 2
	}
}

func (Card2511002ShiningShield) OnDefend(ctx *EffectContext) error {
	success, _ := ctx.ExtraData["defense_success"].(bool)
	defender, _ := ctx.ExtraData["defender"].(int)
	if !success || defender != ctx.PlayerID || !useTriggeredTurn(ctx.Source) {
		return nil
	}
	for _, target := range ctx.Engine.getAllFieldCards(ctx.Engine.State.Players[ctx.OpponentID]) {
		if target == nil || target.Position == nil || !ctx.Engine.IsInSpellRange(ctx.PlayerID, target.Position.Col, target.Position.Row, false) {
			continue
		}
		ctx.Engine.addStatus(target, StatusStun, 1)
	}
	return nil
}
