package game

type Card4611003Jieying struct{ AlwaysActive }

func (Card4611003Jieying) ID() string { return "4611003" }

func (Card4611003Jieying) Name() string { return "咒言师 结影" }

func (Card4611003Jieying) OnTurnStart(ctx *EffectContext) error {
	if ctx.Source.Statuses["咒言书"] == 0 {
		addCardToDeck(ctx, "2601002", 3)
		ctx.Source.Statuses["咒言书"] = 1
	}
	return nil
}
