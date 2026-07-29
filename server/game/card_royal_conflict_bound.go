package game

type Card1011103Gambler struct{ AlwaysActive }

func (Card1011103Gambler) ID() string   { return "1011103" }
func (Card1011103Gambler) Name() string { return "\"弈者\"" }
func (Card1011103Gambler) OnEnter(ctx *EffectContext) error {
	bindSkillToHost(ctx, "3001101")
	return nil
}

type Card2511102FiveRainbowRing struct{ AlwaysActive }

func (Card2511102FiveRainbowRing) ID() string   { return "2511102" }
func (Card2511102FiveRainbowRing) Name() string { return "五虹之环" }
func (Card2511102FiveRainbowRing) OnEnter(ctx *EffectContext) error {
	bindSkillToHost(ctx, "3501101")
	return nil
}
