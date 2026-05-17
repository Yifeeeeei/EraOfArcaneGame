package game

type Card1621005CursedGolem struct{ AlwaysActive }

func (Card1621005CursedGolem) ID() string   { return "1621005" }
func (Card1621005CursedGolem) Name() string { return "诅咒魔像" }
func (Card1621005CursedGolem) OnEnter(ctx *EffectContext) error {
	opponent := ctx.Engine.State.Players[ctx.OpponentID]
	for _, skill := range opponent.Skills {
		if skill != nil {
			skill.Statuses[StatusWeaken] += 2
			return nil
		}
	}
	return nil
}
