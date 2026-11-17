package game

func (Card1221001DolphinPartner) CanPreventLethalBySacrifice(ctx *EffectContext, event DamageEvent) bool {
	return event.Target != nil && event.Target != ctx.Source && event.Target.OwnerID == ctx.PlayerID &&
		event.Amount > 0 && event.Target.CurrentLife-event.Amount <= 0
}

type Card1221001DolphinPartner struct{ AlwaysActive }

func (Card1221001DolphinPartner) ID() string { return "1221001" }

func (Card1221001DolphinPartner) Name() string { return "海豚伙伴" }
