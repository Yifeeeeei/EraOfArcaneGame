package game

type Card2021011LifeAmulet struct{}

func (Card2021011LifeAmulet) ID() string   { return "2021011" }
func (Card2021011LifeAmulet) Name() string { return "生命护符" }

func (Card2021011LifeAmulet) OnEnter(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "life_amulet_buff",
		"选择1个友方角色获得+1血", candidates, 1, 1,
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
