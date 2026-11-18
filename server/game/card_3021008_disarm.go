package game

type Card3021008Disarm struct{ AlwaysActive }

func (Card3021008Disarm) ID() string { return "3021008" }

func (Card3021008Disarm) Name() string { return "缴械" }

func (Card3021008Disarm) OnSpellHit(ctx *EffectContext) error {
	if !isOwnSpellHit(ctx) {
		return nil
	}
	owner := ctx.OpponentID
	if ctx.Target != nil {
		owner = ctx.Target.OwnerID
	}
	candidates := ctx.Engine.friendlyEquipment(owner, nil)
	for _, candidate := range candidates {
		if owner != ctx.PlayerID {
			candidate["side"] = "enemy"
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "disarm_destroy_equipment",
		"选择命中目标控制者的1个装备摧毁", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			for i, item := range ctx.Engine.State.Players[owner].Equipment {
				if item != nil && item.InstanceID == selected[0] {
					ctx.Engine.moveEquipmentToGraveyard(owner, i, item)
					break
				}
			}
		})
	return nil
}
