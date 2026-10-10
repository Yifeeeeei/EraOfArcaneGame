package game

type Card1611003HeartPiercer struct{ AlwaysActive }

func (Card1611003HeartPiercer) ID() string { return "1611003" }

func (Card1611003HeartPiercer) Name() string { return "\"穿心人\"" }

func (Card1611003HeartPiercer) OnEnter(ctx *EffectContext) error {
	addGeneratedCardToHand(ctx, "2601001")
	return nil
}
