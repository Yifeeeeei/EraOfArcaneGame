package game

type royalShieldOnEnter struct {
	AlwaysActive
	id     string
	name   string
	amount int
}

func (c royalShieldOnEnter) ID() string   { return c.id }
func (c royalShieldOnEnter) Name() string { return c.name }
func (c royalShieldOnEnter) OnEnter(ctx *EffectContext) error {
	ctx.Engine.gainPlayerShield(ctx.PlayerID, c.amount)
	return nil
}

type Card1421102EmeraldGuard struct{ AlwaysActive }

func (Card1421102EmeraldGuard) ID() string   { return "1421102" }
func (Card1421102EmeraldGuard) Name() string { return "翡翠守卫" }
func (Card1421102EmeraldGuard) OnEnter(ctx *EffectContext) error {
	if ctx.Engine.State.Players[ctx.PlayerID].Shield == 0 {
		ctx.Engine.gainPlayerShield(ctx.PlayerID, 2)
	}
	return nil
}

type Card1021110RockWallGuard struct{ AlwaysActive }

func (Card1021110RockWallGuard) ID() string   { return "1021110" }
func (Card1021110RockWallGuard) Name() string { return "岩壁护卫军" }
func (Card1021110RockWallGuard) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Source == nil || isFriendlySpellHit(ctx) {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ps == nil || ps.Shield > 0 {
		return nil
	}
	ctx.Engine.gainPlayerShield(ctx.PlayerID, 2)
	return nil
}

type Card2011101ArcaneArmorSky struct{ AlwaysActive }

func (Card2011101ArcaneArmorSky) ID() string   { return "2011101" }
func (Card2011101ArcaneArmorSky) Name() string { return "奥术铠甲 天穹" }
func (Card2011101ArcaneArmorSky) OnEnter(ctx *EffectContext) error {
	ctx.Engine.gainPlayerShield(ctx.PlayerID, 2)
	ctx.Engine.State.Players[ctx.PlayerID].CannotGainShield = true
	return nil
}

type Card2021102DemonBreakingBlade struct{ AlwaysActive }

func (Card2021102DemonBreakingBlade) ID() string   { return "2021102" }
func (Card2021102DemonBreakingBlade) Name() string { return "破魔之刃" }
func (Card2021102DemonBreakingBlade) OnEnter(ctx *EffectContext) error {
	ctx.Engine.losePlayerShield(ctx.OpponentID, 3)
	return nil
}

type Card2021113ArcaneBarrierScroll struct{ AlwaysActive }

func (Card2021113ArcaneBarrierScroll) ID() string   { return "2021113" }
func (Card2021113ArcaneBarrierScroll) Name() string { return "奥术结界卷轴" }
func (Card2021113ArcaneBarrierScroll) OnSpellHitBeforeDamage(ctx *EffectContext) error {
	ctx.Engine.gainPlayerShield(ctx.PlayerID, 2)
	return nil
}

type Card2221102OceanShieldScroll struct{ AlwaysActive }

func (Card2221102OceanShieldScroll) ID() string   { return "2221102" }
func (Card2221102OceanShieldScroll) Name() string { return "海洋之盾卷轴" }
func (Card2221102OceanShieldScroll) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.gainPlayerShield(ctx.PlayerID, 2)
	return nil
}

type Card2421107EmeraldBarrierScroll struct{ AlwaysActive }

func (Card2421107EmeraldBarrierScroll) ID() string   { return "2421107" }
func (Card2421107EmeraldBarrierScroll) Name() string { return "翡翠结界卷轴" }
func (Card2421107EmeraldBarrierScroll) OnUseItem(ctx *EffectContext) error {
	own := learnedSkillCount(ctx.Engine.State.Players[ctx.PlayerID])
	enemy := learnedSkillCount(ctx.Engine.State.Players[ctx.OpponentID])
	if diff := enemy - own; diff > 0 {
		ctx.Engine.gainPlayerShield(ctx.PlayerID, diff)
	}
	return nil
}

type Card2411101EmeraldImmortality struct{ AlwaysActive }

func (Card2411101EmeraldImmortality) ID() string   { return "2411101" }
func (Card2411101EmeraldImmortality) Name() string { return "翡翠永生" }
func (Card2411101EmeraldImmortality) OnEnter(ctx *EffectContext) error {
	ctx.Engine.gainPlayerShield(ctx.PlayerID, 2)
	return nil
}
func (Card2411101EmeraldImmortality) PreventsFieldDamage(ctx *EffectContext) bool {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	return ps != nil && ps.Shield > 0
}

func learnedSkillCount(ps *PlayerState) int {
	if ps == nil {
		return 0
	}
	count := 0
	for _, skill := range ps.Skills {
		if skill != nil {
			count++
		}
	}
	return count
}
