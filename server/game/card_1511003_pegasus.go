package game

type Card1511003Pegasus struct{ AlwaysActive }

func (Card1511003Pegasus) ID() string { return "1511003" }

func (Card1511003Pegasus) Name() string { return "天枢圣兽 珀伽索斯" }

func (Card1511003Pegasus) PreventsFieldDamage(ctx *EffectContext) bool {
	if ctx.Source == nil || ctx.Target == nil || ctx.Source == ctx.Target {
		return false
	}
	if ctx.ExtraData["damage_source"] != "spell" {
		return false
	}
	attacker, ok := ctx.ExtraData["attacker"].(int)
	return ok && attacker != ctx.PlayerID
}
