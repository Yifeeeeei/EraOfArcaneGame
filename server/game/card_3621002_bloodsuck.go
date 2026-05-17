package game

type Card3621002Bloodsuck struct{ AlwaysActive }

func (Card3621002Bloodsuck) ID() string   { return "3621002" }
func (Card3621002Bloodsuck) Name() string { return "噬血" }

func (Card3621002Bloodsuck) OnSpellHit(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "bloodsuck_buff",
		"选择1个友方单位获得+2血", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, selected[0])
			if target == nil || zone != "unit" {
				return
			}
			target.CurrentLife += 2
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source": cardToInfo(ctx.Source),
				"target": cardToInfo(target),
				"effect": "modify_life",
				"amount": 2,
			}})
		})
	return nil
}
