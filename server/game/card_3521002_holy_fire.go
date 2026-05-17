package game

type Card3521002HolyFire struct{}

func (Card3521002HolyFire) ID() string   { return "3521002" }
func (Card3521002HolyFire) Name() string { return "神圣之火" }
func (Card3521002HolyFire) AllowsFriendlySpellTarget() bool {
	return true
}
func (Card3521002HolyFire) SpellDamage(ctx *EffectContext) int {
	target, _ := ctx.ExtraData["target"].(SpellTarget)
	if target.Type == "unit" && target.Position.Valid() {
		if unit := ctx.Engine.State.Players[ctx.PlayerID].Units[target.Position.Col][target.Position.Row]; unit != nil {
			return 0
		}
	}
	return ctx.Source.Card.Attack
}
func (Card3521002HolyFire) OnSpellHit(ctx *EffectContext) error {
	if ctx.Target == nil || ctx.Target.OwnerID != ctx.PlayerID {
		return nil
	}
	for _, status := range negativeStatuses {
		ctx.Target.Statuses[status] = 0
	}
	return nil
}
