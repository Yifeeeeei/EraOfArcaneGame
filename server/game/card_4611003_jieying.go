package game

type Card4611003Jieying struct{ AlwaysActive }

func (Card4611003Jieying) ID() string { return "4611003" }

func (Card4611003Jieying) Name() string { return "咒言师 结影" }

func (Card4611003Jieying) OnEnter(ctx *EffectContext) error {
	addCardToDeck(ctx, "2601002", 3)
	return nil
}
