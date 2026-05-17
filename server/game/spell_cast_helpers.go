package game

func isFriendlySpellCast(ctx *EffectContext) bool {
	if ctx == nil || ctx.ExtraData == nil {
		return true
	}
	castPlayer, ok := ctx.ExtraData["cast_player"].(int)
	return !ok || castPlayer == ctx.PlayerID
}

func isEnemySpellCast(ctx *EffectContext) bool {
	if ctx == nil || ctx.ExtraData == nil {
		return false
	}
	castPlayer, ok := ctx.ExtraData["cast_player"].(int)
	return ok && castPlayer != ctx.PlayerID
}

func isSpellBeingCast(ctx *EffectContext) bool {
	return ctx != nil && ctx.Target == nil
}

func isFriendlySpellHit(ctx *EffectContext) bool {
	if ctx == nil || ctx.ExtraData == nil {
		return true
	}
	attacker, ok := ctx.ExtraData["attacker"].(int)
	return !ok || attacker == ctx.PlayerID
}
