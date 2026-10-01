package game

type Card1421113RockWallMonk struct{ AlwaysActive }

func (Card1421113RockWallMonk) ID() string { return "1421113" }

func (Card1421113RockWallMonk) Name() string { return "岩壁修道士" }

func (Card1421113RockWallMonk) OnSpellHitBeforeDamage(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil {
		return nil
	}
	if ctx.Source.UsedThisTurn >= perTurnLimit(ctx.Source) {
		return nil
	}
	attacker, ok := ctx.ExtraData["attacker"].(int)
	if !ok || attacker == ctx.PlayerID {
		return nil
	}
	if playerHasLearnedSkill(ctx.Engine.State.Players[ctx.PlayerID]) {
		return nil
	}
	damagePtr, ok := ctx.ExtraData["damage_ptr"].(*int)
	if !ok || damagePtr == nil {
		return nil
	}
	*damagePtr = 0
	ctx.ExtraData["damage"] = 0
	ctx.Source.UsedThisTurn++
	return nil
}

func playerHasLearnedSkill(ps *PlayerState) bool {
	if ps == nil {
		return false
	}
	for _, skill := range ps.Skills {
		if skill != nil {
			return true
		}
	}
	return false
}

var _ OnSpellHitBeforeDamageBehavior = Card1421113RockWallMonk{}
