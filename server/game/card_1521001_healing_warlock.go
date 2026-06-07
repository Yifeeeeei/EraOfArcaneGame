package game

type Card1521001HealingWarlock struct{ AlwaysActive }

func (Card1521001HealingWarlock) ID() string            { return "1521001" }
func (Card1521001HealingWarlock) Name() string          { return "治疗术士" }
func (Card1521001HealingWarlock) IsPrayerAbility() bool { return true }
func (Card1521001HealingWarlock) OnPerTurn(ctx *EffectContext) error {
	if ctx.ExtraData != nil && ctx.ExtraData["prayer"] == true {
		candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, nil)
		if len(candidates) == 0 {
			return nil
		}
		ctx.Engine.SetPendingAction(ctx.PlayerID, "healing_warlock_prayer",
			"治疗术士:选择1个友方单位回复1生命", candidates, 1, 1,
			func(selected []string) {
				if len(selected) == 0 {
					return
				}
				target := ctx.Engine.findFieldCardByInstance(ctx.Engine.State.Players[ctx.PlayerID], selected[0])
				if target != nil {
					healUnit(target, 1)
				}
			})
		return nil
	}
	if ctx.Target == nil || ctx.Target.OwnerID != ctx.PlayerID {
		return nil
	}
	healUnit(ctx.Target, 1)
	return nil
}
