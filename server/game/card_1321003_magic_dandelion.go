package game

type Card1321003MagicDandelion struct{}

func (Card1321003MagicDandelion) ID() string   { return "1321003" }
func (Card1321003MagicDandelion) Name() string { return "魔法蒲公英" }

func (Card1321003MagicDandelion) RevealsOnDraw() bool { return true }

func (Card1321003MagicDandelion) OnEnter(ctx *EffectContext) error {
	// Runtime still lacks per-instance draw-turn tracking. Keep the explicit
	// playable behavior that used to come from text parsing: entering draws 1.
	return DrawCards(1)(ctx)
}
