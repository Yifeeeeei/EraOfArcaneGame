package game

type Card2621012ShadowCloak struct{ AlwaysActive }

func (Card2621012ShadowCloak) ID() string { return "2621012" }

func (Card2621012ShadowCloak) Name() string { return "暗影披风" }

func (Card2621012ShadowCloak) HasActiveSpellHitBeforeDamage(card *CardInstance) bool {
	return card != nil && card.Statuses[shadowCloakUsedStatus] == 0
}

func (Card2621012ShadowCloak) OnSpellHitBeforeDamage(ctx *EffectContext) error {
	if ctx == nil || ctx.ExtraData == nil || ctx.Source == nil {
		return nil
	}
	attacker, ok := ctx.ExtraData["attacker"].(int)
	if !ok || attacker == ctx.PlayerID {
		return nil
	}
	damagePtr, ok := ctx.ExtraData["damage_ptr"].(*int)
	if !ok || damagePtr == nil {
		return nil
	}
	*damagePtr = 0
	ctx.ExtraData["damage"] = 0
	ctx.Source.Statuses[shadowCloakUsedStatus] = 1
	return nil
}
