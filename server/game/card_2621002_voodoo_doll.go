package game

type Card2621002VoodooDoll struct{ AlwaysActive }

func (Card2621002VoodooDoll) ID() string { return "2621002" }

func (Card2621002VoodooDoll) Name() string { return "巫毒娃娃" }

func (Card2621002VoodooDoll) OnEquip(ctx *EffectContext) error {
	ctx.Source.Statuses["暗影标记"] = 3
	candidates := voodooDollLinkCandidates(ctx)
	if len(candidates) < 2 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "voodoo_doll_link", "巫毒娃娃:选择法力范围内2个伙伴建立连结", candidates, 2, 2, func(selected []string) {
		clearVoodooDollLinks(ctx.Source)
		for _, id := range selected {
			ctx.Source.Statuses["巫毒连结:"+id] = 1
		}
	})
	return nil
}

func (Card2621002VoodooDoll) DamageScope() DamageScope { return DamageAny }

func (Card2621002VoodooDoll) OnDamaged(ctx *EffectContext, event DamageEvent) error {
	if event.Kind == "voodoo_doll" {
		return nil
	}
	damage := event.Amount
	if damage <= 0 || ctx.Source.Statuses["暗影标记"] < damage || event.Target == nil || !voodooDollIsLinked(ctx.Source, event.Target) {
		return nil
	}
	linked := voodooDollOtherLinkedUnit(ctx.Engine, ctx.Source, event.Target)
	if linked == nil {
		return nil
	}
	candidates := []map[string]any{candidateInfo(linked, "unit", voodooDollSide(ctx.PlayerID, linked))}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "voodoo_doll_damage", "巫毒娃娃:是否让连结的另一伙伴受到同等伤害", candidates, 0, 1, func(selected []string) {
		if len(selected) == 0 || selected[0] != linked.InstanceID || ctx.Source.Statuses["暗影标记"] < damage {
			return
		}
		ctx.Source.Statuses["暗影标记"] -= damage
		ctx.Engine.ApplyDamage(DamageRequest{Target: linked, Amount: damage, Kind: "voodoo_doll", SourcePlayer: ctx.PlayerID, SourceKnown: true, Source: ctx.Source})
	})
	return nil
}
