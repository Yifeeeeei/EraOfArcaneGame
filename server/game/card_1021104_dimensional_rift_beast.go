package game

type Card1021104DimensionalRiftBeast struct{ AlwaysActive }

func (Card1021104DimensionalRiftBeast) ID() string { return "1021104" }

func (Card1021104DimensionalRiftBeast) Name() string { return "次元撕裂兽" }

func (Card1021104DimensionalRiftBeast) OnEnter(ctx *EffectContext) error {
	candidates := companionSpellRangeCandidates(ctx, false)
	filtered := candidates[:0]
	for _, candidate := range candidates {
		if candidate["side"] == "enemy" {
			filtered = append(filtered, candidate)
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "dimensional_rift_beast_exile",
		"次元撕裂兽:选择法力范围内1个敌方伙伴移出游戏", filtered, 1, 1,
		func(selected []string) {
			target := selectedUnitFromCandidates(ctx.Engine, selected, filtered)
			if target != nil && target.Card != nil && target.Card.IsCompanion() {
				ctx.Engine.exileCard(target.OwnerID, target)
			}
		})
	return nil
}
