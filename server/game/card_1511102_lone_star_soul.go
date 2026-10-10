package game

type Card1511102LoneStarSoul struct{ AlwaysActive }

func (Card1511102LoneStarSoul) ID() string { return "1511102" }

func (Card1511102LoneStarSoul) Name() string { return "孤星之魂 凯拉莫将军" }

func (Card1511102LoneStarSoul) DamageScope() DamageScope { return DamageSelf }

func (Card1511102LoneStarSoul) OnDamaged(ctx *EffectContext, event DamageEvent) error {
	if ctx == nil || ctx.Source == nil {
		return nil
	}
	attacker, hasAttacker := event.SourcePlayer, event.SourcePlayer >= 0
	if !hasAttacker || attacker == ctx.PlayerID {
		return nil
	}
	if len(adjacentFriendlyCompanions(ctx)) > 0 {
		return nil
	}
	ctx.Engine.gainPlayerShield(ctx.PlayerID, 1)
	ctx.Source.CurrentAttack++
	return nil
}
