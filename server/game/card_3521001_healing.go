package game

type Card3521001Healing struct{ AlwaysActive }

func (Card3521001Healing) ID() string   { return "3521001" }
func (Card3521001Healing) Name() string { return "治疗术" }

func (Card3521001Healing) OnSpellCast(ctx *EffectContext) error {
	if !isFriendlySpellCast(ctx) || !isSpellBeingCast(ctx) {
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card.CurrentLife < card.Card.Life
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "healing_spell",
		"选择1个友方单位回复2血", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			target, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, selected[0])
			if target == nil || zone != "unit" {
				return
			}
			target.CurrentLife += 2
			if target.CurrentLife > target.Card.Life {
				target.CurrentLife = target.Card.Life
			}
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source": cardToInfo(ctx.Source),
				"target": cardToInfo(target),
				"effect": "heal",
				"amount": 2,
			}})
		})
	return nil
}
