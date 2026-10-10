package game

type Card2421103Dreamcatcher struct{ AlwaysActive }

func (Card2421103Dreamcatcher) ID() string { return "2421103" }

func (Card2421103Dreamcatcher) Name() string { return "捕梦网" }

func (Card2421103Dreamcatcher) OnEnter(ctx *EffectContext) error {
	for _, skill := range ctx.Engine.State.Players[ctx.PlayerID].Skills {
		if skill != nil && skill.Card != nil && isSpellLikeCard(skill.Card) && hasCardTag(skill.Card, "灵媒") {
			skill.PowerBonus += 2
			ctx.Engine.refreshPendingSpellPowerForModifiedSkill(ctx.PlayerID, skill)
		}
	}
	return nil
}
