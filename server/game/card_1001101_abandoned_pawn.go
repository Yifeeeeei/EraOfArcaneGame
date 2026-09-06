package game

type Card1001101AbandonedPawn struct{ AlwaysActive }

func (Card1001101AbandonedPawn) ID() string { return "1001101" }

func (Card1001101AbandonedPawn) Name() string { return "弃子" }

func (Card1001101AbandonedPawn) OnDeath(ctx *EffectContext) error {
	if ctx.Source == nil || ctx.Source.Position == nil {
		return nil
	}
	pos := *ctx.Source.Position
	damaged := adjacentUnits(ctx.Engine.State.Players[ctx.PlayerID], &pos)
	damaged = append(damaged, adjacentUnits(ctx.Engine.State.Players[ctx.OpponentID], &pos)...)
	for _, target := range damaged {
		if target == nil || target.CurrentLife <= 0 {
			continue
		}
		targetPos := Position{}
		if target.Position != nil {
			targetPos = *target.Position
		}
		ctx.Engine.ApplyDamage(DamageRequest{Target: target, Amount: 1, Kind: "effect", SourcePlayer: ctx.PlayerID, SourceKnown: true, Source: ctx.Source})
		if target.CurrentLife <= 0 && !target.Card.IsHero() {
			ownerID := target.OwnerID
			if ctx.Engine.unitInOwnerGrid(target, ownerID) {
				ctx.Engine.destroyUnitWithData(target, ownerID, map[string]any{
					"death_cause": "abandoned_pawn",
					"attacker":    ctx.PlayerID,
				})
			}
			if ctx.Engine.State.Players[ownerID].Units[targetPos.Col][targetPos.Row] == nil {
				ctx.Engine.summonFreshCardAtPosition(ownerID, "1001101", targetPos, true)
			}
		}
	}
	return nil
}
