package game

type Card4311001Su struct{}

func (Card4311001Su) ID() string   { return "4311001" }
func (Card4311001Su) Name() string { return "雷术士 肃" }
func (Card4311001Su) OnUltimate(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	discarded := 0
	for i := len(ps.Hand) - 1; i >= 0 && discarded < 2; i-- {
		if ps.Hand[i].Card.Category == "气" {
			ps.Graveyard = append(ps.Graveyard, ps.Hand[i])
			ps.Hand = append(ps.Hand[:i], ps.Hand[i+1:]...)
			discarded++
		}
	}
	if discarded < 2 {
		return nil
	}
	target := findAnyUnit(ctx.Engine.State.Players[ctx.OpponentID])
	if target != nil {
		ctx.Engine.dealDamage(target, 1, ctx.OpponentID)
	}
	return nil
}
