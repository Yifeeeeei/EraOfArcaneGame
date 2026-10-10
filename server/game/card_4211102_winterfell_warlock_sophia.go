package game

type Card4211102WinterfellWarlockSophia struct{ AlwaysActive }

func (Card4211102WinterfellWarlockSophia) ID() string { return "4211102" }

func (Card4211102WinterfellWarlockSophia) Name() string { return "凛冰魔巫 索菲娅" }

func (Card4211102WinterfellWarlockSophia) HasNegativeStatusImmunity(status string) bool {
	return status == StatusFreeze
}

func (Card4211102WinterfellWarlockSophia) OnUltimate(ctx *EffectContext) error {
	candidates := make([]map[string]any, 0)
	for playerID, ps := range ctx.Engine.State.Players {
		if ps == nil {
			continue
		}
		side := "enemy"
		if playerID == ctx.PlayerID {
			side = "own"
		}
		for _, unit := range ps.Units {
			for _, card := range unit {
				if card != nil && card.Card != nil && (card.Card.IsHero() || card.Card.IsCompanion()) && card.Statuses[StatusFreeze] > 0 {
					candidates = append(candidates, candidateInfo(card, "unit", side))
				}
			}
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	allowed := make(map[string]bool, len(candidates))
	for _, candidate := range candidates {
		if id, _ := candidate["instance_id"].(string); id != "" {
			allowed[id] = true
		}
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "sophia_thaw_strike",
		"凛冰魔巫 索菲娅:选择1个冻结单位移除1层冻结并造成2点伤害", candidates, 1, 1,
		func(selected []string) {
			id := firstSelected(selected)
			if !allowed[id] {
				return
			}
			target := ctx.Engine.findUnitByInstanceID(id)
			if target == nil || target.Card == nil || (!target.Card.IsHero() && !target.Card.IsCompanion()) || target.Statuses[StatusFreeze] <= 0 {
				return
			}
			target.Statuses[StatusFreeze]--
			ctx.Engine.ApplyDamage(DamageRequest{Target: target, Amount: 2, Kind: "effect", SourcePlayer: ctx.PlayerID, SourceKnown: true, Source: ctx.Source})
		})
	return nil
}
