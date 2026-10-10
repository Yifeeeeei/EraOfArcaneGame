package game

type Card1621114SoulSymbiote struct{ AlwaysActive }

func (Card1621114SoulSymbiote) ID() string { return "1621114" }

func (Card1621114SoulSymbiote) Name() string { return "灵魂共生体" }

func (Card1621114SoulSymbiote) OnDeath(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlySkillsIncludingBound(ctx.PlayerID, func(skill *CardInstance) bool {
		return skill != nil && skill.Card != nil && skill.Card.IsSkill()
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "soul_symbiote_mark_skills",
		"灵魂共生体:选择最多2个法术放置灵魂标记物", candidates, 0, 2,
		func(selected []string) {
			seen := map[string]bool{}
			for _, id := range selected {
				if seen[id] {
					continue
				}
				seen[id] = true
				skill := ctx.Engine.findSkill(ctx.Engine.State.Players[ctx.PlayerID], id)
				if skill == nil || skill.Card == nil || !skill.Card.IsSkill() {
					continue
				}
				addSoulMarkerToSpell(skill)
			}
		})
	return nil
}
