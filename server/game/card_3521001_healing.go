package game

type Card3521001Healing struct{ AlwaysActive }

func (Card3521001Healing) ID() string   { return "3521001" }
func (Card3521001Healing) Name() string { return "治疗术" }

func (Card3521001Healing) AllowsFriendlySpellTarget() bool {
	return true
}

func (Card3521001Healing) OnSpellHit(ctx *EffectContext) error {
	target := ctx.Target
	if target == nil || target.OwnerID != ctx.PlayerID || target.CurrentLife >= target.Card.Life {
		return nil
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
	return nil
}
