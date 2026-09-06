package game

type Card1221010WallKeeper struct{ AlwaysActive }

func (Card1221010WallKeeper) ID() string { return "1221010" }

func (Card1221010WallKeeper) Name() string { return "护壁者" }

func (Card1221010WallKeeper) OnEnter(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{Type: TempModAllSpellDamageZero, RemainingUses: 1, ExpiresTurn: ctx.Engine.State.TurnNumber + 1})
	return nil
}
