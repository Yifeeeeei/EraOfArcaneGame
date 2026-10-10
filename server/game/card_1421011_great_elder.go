package game

type Card1421011GreatElder struct{ AlwaysActive }

func (Card1421011GreatElder) ID() string { return "1421011" }

func (Card1421011GreatElder) Name() string { return "大长老" }

func (Card1421011GreatElder) MasteryMax() int { return 3 }

func (Card1421011GreatElder) OnMastery(ctx *EffectContext, level int) error {
	if level == 1 || level == 3 {
		ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
			Type:          TempModNextEarthSkillLearnCostMinus,
			Amount:        2,
			RemainingUses: 1,
		})
	}
	return nil
}
