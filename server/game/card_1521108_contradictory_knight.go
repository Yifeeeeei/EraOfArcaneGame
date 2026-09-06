package game

type Card1521108ContradictoryKnight struct{ AlwaysActive }

func (Card1521108ContradictoryKnight) ID() string { return "1521108" }

func (Card1521108ContradictoryKnight) Name() string { return "矛盾的骑士" }

func (Card1521108ContradictoryKnight) OnDeath(ctx *EffectContext) error {
	opponentID := ctx.OpponentID
	candidates := ctx.Engine.friendlyEmptyUnitPositions(opponentID)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(opponentID, "contradictory_knight_summon",
		"矛盾的骑士:选择位置为你召唤此卡", candidates, 1, 1,
		func(selected []string) {
			pos, ok := positionFromSelectionID(firstSelected(selected))
			if !ok {
				return
			}
			opponent := ctx.Engine.State.Players[opponentID]
			if opponent.Units[pos.Col][pos.Row] != nil {
				return
			}
			if !ctx.Engine.removeCardFromGraveyard(ctx.PlayerID, ctx.Source) {
				return
			}
			if ctx.Source.Card.Life > 1 {
				cardCopy := *ctx.Source.Card
				cardCopy.Life--
				ctx.Source.Card = &cardCopy
			}
			ctx.Source.OwnerID = opponentID
			ctx.Source.CurrentLife = ctx.Source.Card.Life
			ctx.Source.CurrentAttack = ctx.Source.Card.Attack
			ctx.Source.DamageTakenThisTurn = 0
			ctx.Source.IsHorizontal = true
			ctx.Source.Position = nil
			ctx.Source.Statuses = make(map[string]int)
			ctx.Source.ElementsGainBonus = make(map[string]int)
			ctx.Source.ElementsGainSet = nil
			ctx.Source.PowerBonus = 0
			ctx.Source.AttackBonus = 0
			ctx.Source.UsedThisTurn = 0
			ctx.Source.UltimateUsed = false
			ctx.Engine.exileTransferredBoundSkills(ctx.PlayerID, ctx.Source)
			ctx.Source.BoundSkills = nil
			ctx.Source.UnderCards = nil
			ctx.Source.AttachedBehaviors = nil
			if !ctx.Engine.placeExistingCompanionAtPosition(opponentID, ctx.Source, pos, true) {
				ctx.Engine.addToGraveyard(ctx.PlayerID, ctx.Source)
			}
		})
	return nil
}
