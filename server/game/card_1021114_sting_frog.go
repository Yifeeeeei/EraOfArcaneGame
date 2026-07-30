package game

type Card1021114StingFrog struct{ AlwaysActive }

func (Card1021114StingFrog) ID() string   { return "1021114" }
func (Card1021114StingFrog) Name() string { return "蛰蛙" }

func (Card1021114StingFrog) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Target == nil || ctx.Target.Card == nil || !isFriendlySpellCast(ctx) {
		return nil
	}
	if !hasCardTag(ctx.Target.Card, "驱动") && !hasCardTag(ctx.Target.Card, "聚能") {
		return nil
	}
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:        TempModSkillPowerBonus,
		Amount:      1,
		ExpiresTurn: ctx.Engine.State.TurnNumber,
	})
	return nil
}

var _ OnSpellCastBehavior = Card1021114StingFrog{}
