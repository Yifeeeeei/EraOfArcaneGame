package game

type Card2321012WindCloak struct{ AlwaysActive }

func (Card2321012WindCloak) ID() string { return "2321012" }

func (Card2321012WindCloak) Name() string { return "随风斗篷" }

func (Card2321012WindCloak) OnUltimate(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	pos := ps.FindEmptyPosition()
	if pos != nil && ps.Hero != nil && ps.Hero.Position != nil {
		ps.Units[ps.Hero.Position.Col][ps.Hero.Position.Row] = nil
		ps.Hero.Position = pos
		ps.Units[pos.Col][pos.Row] = ps.Hero
	}
	return nil
}
