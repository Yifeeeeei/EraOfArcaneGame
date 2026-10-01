package game

type Card1411001GreatDruidCycle struct{ AlwaysActive }

func (Card1411001GreatDruidCycle) ID() string { return "1411001" }

func (Card1411001GreatDruidCycle) Name() string { return "\"轮回不息\" 大德鲁伊 烟尘" }

func (Card1411001GreatDruidCycle) HasActiveUltimate(*CardInstance) bool { return false }

func (Card1411001GreatDruidCycle) OnUltimate(ctx *EffectContext) error {
	ctx.Source.Statuses["轮回不息"] = 1
	return nil
}

func (Card1411001GreatDruidCycle) OnFriendlyDeath(ctx *EffectContext) error {
	if ctx.Source.UltimateUsed {
		return nil
	}
	ctx.Source.Statuses["轮回不息"] = 1
	if ctx.Source.Statuses["轮回不息"] == 0 || ctx.Target == nil || !ctx.Target.Card.IsCompanion() {
		return nil
	}
	if ctx.Engine.State.Players[ctx.PlayerID].FindEmptyPosition() == nil || getCardDB()["1401001"] == nil {
		return nil
	}
	ctx.Source.UltimateUsed = true
	candidates := []map[string]any{candidateInfo(ctx.Source, "companion", "own")}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "great_druid_life_seed", "\"轮回不息\" 大德鲁伊 烟尘:是否召唤1个生命种子", candidates, 0, 1, func(selected []string) {
		if len(selected) == 0 {
			return
		}
		positions := ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)
		if len(positions) == 0 {
			return
		}
		ctx.Engine.SetPendingAction(ctx.PlayerID, "great_druid_life_seed_position", "\"轮回不息\" 大德鲁伊 烟尘:选择生命种子位置", positions, 1, 1, func(posSelected []string) {
			if len(posSelected) == 0 {
				return
			}
			pos, ok := positionFromSelectionID(posSelected[0])
			if !ok {
				return
			}
			summonGreatDruidLifeSeedAtPosition(ctx, pos)
		})
	})
	return nil
}
