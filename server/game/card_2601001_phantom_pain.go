package game

type Card2601001PhantomPain struct{ AlwaysActive }

func (Card2601001PhantomPain) ID() string { return "2601001" }

func (Card2601001PhantomPain) Name() string { return "幻痛" }

func (Card2601001PhantomPain) OnDefend(ctx *EffectContext) error {
	success, _ := ctx.ExtraData["defense_success"].(bool)
	defender, _ := ctx.ExtraData["defender"].(int)
	if !success || defender == ctx.PlayerID || !useTriggeredTurn(ctx.Source) {
		return nil
	}
	weakenDefenseCards := func(key string) {
		skills, _ := ctx.ExtraData[key].([]*CardInstance)
		for _, skill := range skills {
			if skill != nil {
				ctx.Engine.addStatus(skill, StatusWeaken, 2)
			}
		}
	}
	weakenDefenseCards("defense_skills")
	weakenDefenseCards("defense_boosts")
	ctx.Engine.promptHeartPiercerAfterPhantomPain(ctx)
	return nil
}
