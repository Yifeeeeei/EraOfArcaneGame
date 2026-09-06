package game

type Card1121113LavaFortHellhound struct{ AlwaysActive }

func (Card1121113LavaFortHellhound) ID() string { return "1121113" }

func (Card1121113LavaFortHellhound) Name() string { return "熔岩堡地狱犬" }

func (Card1121113LavaFortHellhound) OnConsume(ctx *EffectContext) error {
	if ctx == nil || ctx.Source == nil || ctx.Engine == nil || !triggeredTurnAvailable(ctx.Source) || ctx.ExtraData == nil {
		return nil
	}
	if ctx.Target != nil && ctx.Target != ctx.Source {
		return nil
	}
	if source, _ := ctx.ExtraData["consume_source"].(string); source == "" || source == ctx.Source.Card.Number {
		return nil
	}
	candidates := companionSpellRangeCandidates(ctx, false)
	if len(candidates) < 2 {
		return nil
	}
	if !useTriggeredTurn(ctx.Source) {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "lava_fort_hellhound_damage",
		"熔岩堡地狱犬:选择法力范围内2个不同单位各造成1点伤害", candidates, 2, 2,
		func(selected []string) {
			seen := map[string]bool{}
			for _, id := range selected {
				if seen[id] {
					continue
				}
				seen[id] = true
				target := ctx.Engine.findUnitByInstanceID(id)
				if target == nil || target.Card == nil || !target.Card.IsCompanion() || target.Position == nil {
					continue
				}
				if target.OwnerID != ctx.PlayerID && !ctx.Engine.IsInSpellRange(ctx.PlayerID, target.Position.Col, target.Position.Row, false) {
					continue
				}
				ctx.Engine.ApplyDamage(DamageRequest{Target: target, Amount: 1, Kind: "effect", SourcePlayer: ctx.PlayerID, SourceKnown: true, Source: ctx.Source})
			}
		})
	return nil
}
