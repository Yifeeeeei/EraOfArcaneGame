package game

type Card1521006LifeFlower struct{ AlwaysActive }

func (Card1521006LifeFlower) ID() string   { return "1521006" }
func (Card1521006LifeFlower) Name() string { return "生命之花" }

func (Card1521006LifeFlower) OnEnter(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card != nil && card != ctx.Source
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "life_flower_buff",
		"选择1个其他友方单位获得+1血", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, selected[0])
			if target == nil || zone != "unit" {
				return
			}
			target.CurrentLife++
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source": cardToInfo(ctx.Source),
				"target": cardToInfo(target),
				"effect": "modify_life",
				"amount": 1,
			}})
		})
	return nil
}
