package game

type Card1621106SoulHunter struct{ AlwaysActive }

func (Card1621106SoulHunter) ID() string { return "1621106" }

func (Card1621106SoulHunter) Name() string { return "猎魂者" }

func (Card1621106SoulHunter) OnSpellHit(ctx *EffectContext) error {
	if !isFriendlySpellHit(ctx) || !triggeredTurnAvailable(ctx.Source) {
		return nil
	}
	skill := ctx.Target
	if ctx.ExtraData != nil {
		if source, ok := ctx.ExtraData["spell_source"].(*CardInstance); ok && source != nil {
			skill = source
		}
	}
	if skill == nil || skill.Card == nil || !isSpellLikeCard(skill.Card) {
		return nil
	}
	if useTriggeredTurn(ctx.Source) {
		addSoulMarkerToSpell(skill)
	}
	return nil
}
