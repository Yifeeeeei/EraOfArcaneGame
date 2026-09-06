package game

type Card2321007WindwhisperRing struct{ AlwaysActive }

func (Card2321007WindwhisperRing) ID() string { return "2321007" }

func (Card2321007WindwhisperRing) Name() string { return "风语之戒" }

func (Card2321007WindwhisperRing) OnEnter(ctx *EffectContext) error {
	return DrawCards(1)(ctx)
}
