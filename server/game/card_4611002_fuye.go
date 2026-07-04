package game

type Card4611002Fuye struct{ AlwaysActive }

func (Card4611002Fuye) ID() string   { return "4611002" }
func (Card4611002Fuye) Name() string { return "芙雅夫人" }
func (Card4611002Fuye) OnUltimate(ctx *EffectContext) error {
	if ctx.Target != nil {
		applyFuyeUltimate(ctx.Target)
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion()
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "fuye_ultimate_target",
		"芙雅夫人:选择1个友方伙伴，其攻击力和负载翻倍，并获得临时", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, selected[0])
			if target == nil || zone != "unit" || target.Card == nil || !target.Card.IsCompanion() {
				return
			}
			applyFuyeUltimate(target)
		})
	return nil
}

func applyFuyeUltimate(target *CardInstance) {
	target.CurrentAttack *= 2
	gain := effectiveElementsGain(target)
	doubled := make(map[string]int, len(gain))
	for elem, amount := range gain {
		doubled[elem] = amount * 2
	}
	setElementsGain(target, doubled)
	target.Statuses["临时"] = 1
}
