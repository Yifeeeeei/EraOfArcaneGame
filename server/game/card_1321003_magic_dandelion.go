package game

type Card1321003MagicDandelion struct{ AlwaysActive }

func (Card1321003MagicDandelion) ID() string { return "1321003" }

func (Card1321003MagicDandelion) Name() string { return "魔法蒲公英" }

func (Card1321003MagicDandelion) RevealsOnDraw() bool { return true }

func (Card1321003MagicDandelion) OnEnter(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ps == nil || ps.DrawnTurn == nil || ps.DrawnTurn[ctx.Source.InstanceID] != ctx.Engine.State.TurnNumber {
		return nil
	}
	return DrawCards(1)(ctx)
}
