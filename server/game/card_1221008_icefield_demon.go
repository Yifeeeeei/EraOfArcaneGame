package game

type Card1221008IcefieldDemon struct{ AlwaysActive }

func (Card1221008IcefieldDemon) ID() string   { return "1221008" }
func (Card1221008IcefieldDemon) Name() string { return "冰域恶魔" }
func (Card1221008IcefieldDemon) OnEnter(ctx *EffectContext) error {
	opponent := ctx.Engine.State.Players[ctx.OpponentID]
	for _, target := range ctx.Engine.getAllFieldCards(opponent) {
		target.Statuses[StatusFreeze]++
	}
	return nil
}
