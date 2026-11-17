package game

type Card1321013TeleportMage struct{ AlwaysActive }

func (Card1321013TeleportMage) ID() string { return "1321013" }

func (Card1321013TeleportMage) Name() string { return "传送法师" }

func (Card1321013TeleportMage) OnPerTurn(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	var moving *CardInstance
	for col := 0; col < 3 && moving == nil; col++ {
		for row := 0; row < 3; row++ {
			if u := ps.Units[col][row]; u != nil && !u.Card.IsHero() {
				moving = u
				break
			}
		}
	}
	pos := ps.FindEmptyPosition()
	if moving != nil && pos != nil {
		ps.Units[moving.Position.Col][moving.Position.Row] = nil
		moving.Position = pos
		ps.Units[pos.Col][pos.Row] = moving
	}
	return nil
}
