package game

type Card1011003ArchleaderFarolank struct{ AlwaysActive }

func (Card1011003ArchleaderFarolank) ID() string   { return "1011003" }
func (Card1011003ArchleaderFarolank) Name() string { return "盟主 法罗兰克" }

func (Card1011003ArchleaderFarolank) OnEnter(ctx *EffectContext) error {
	if ctx.Source == nil || ctx.Source.Position == nil {
		return nil
	}
	bindSkillToHost(ctx, "3001002")
	gains := make(map[string]int)
	for _, unit := range adjacentUnits(ctx.Engine.State.Players[ctx.PlayerID], ctx.Source.Position) {
		if unit.Card.IsCompanion() {
			for elem, amount := range effectiveElementsGain(unit) {
				gains[elem] += amount
			}
		}
	}
	for elem, amount := range gains {
		ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, elem, amount, ctx.Source)
	}
	return nil
}
