package game

type Card2321109MistMask struct{ AlwaysActive }

func (Card2321109MistMask) ID() string { return "2321109" }

func (Card2321109MistMask) Name() string { return "幻雾面罩" }

func (Card2321109MistMask) OnSpellHitBeforeDamage(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil || ctx.Source.UltimateUsed {
		return nil
	}
	attacker, ok := ctx.ExtraData["attacker"].(int)
	damagePtr, hasDamage := ctx.ExtraData["damage_ptr"].(*int)
	if !ok || attacker == ctx.PlayerID || !hasDamage || damagePtr == nil || *damagePtr <= 0 {
		return nil
	}
	candidates := ctx.Engine.friendlyHandCards(ctx.PlayerID, nil)
	if len(candidates) == 0 {
		return nil
	}
	mask := ctx.Source
	ctx.Engine.SetPendingAction(ctx.PlayerID, "mist_mask_discard_reduce_spell_attack",
		"幻雾面罩:丢弃最多3张手牌降低该法术伤害", candidates, 0, min(3, len(candidates)),
		func(selected []string) {
			if len(selected) == 0 || mask.UltimateUsed || !ctx.Engine.cardStillOnField(mask) {
				return
			}
			discarded := ctx.Engine.discardSelectedHandCardsMatching(ctx.PlayerID, selected, 3, nil)
			reduction := min(len(discarded), *damagePtr)
			if reduction <= 0 {
				return
			}
			*damagePtr -= reduction
			ctx.ExtraData["damage"] = *damagePtr
			mask.UltimateUsed = true
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source":    cardToInfo(mask),
				"effect":    "mist_mask_reduce_spell_attack",
				"discarded": len(discarded),
				"reduction": reduction,
				"damage":    *damagePtr,
			}})
		})
	return nil
}
