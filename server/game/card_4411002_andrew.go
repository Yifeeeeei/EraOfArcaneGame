package game

type Card4411002Andrew struct{ AlwaysActive }

func (Card4411002Andrew) ID() string { return "4411002" }

func (Card4411002Andrew) Name() string { return "大法师 安德鲁" }

func (Card4411002Andrew) OnEnter(ctx *EffectContext) error {
	addCardToDeck(ctx, "1401002", 1)
	return nil
}
