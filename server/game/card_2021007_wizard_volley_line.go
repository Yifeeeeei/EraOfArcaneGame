package game

type Card2021007WizardVolleyLine struct{ AlwaysActive }

func (Card2021007WizardVolleyLine) ID() string   { return "2021007" }
func (Card2021007WizardVolleyLine) Name() string { return "巫师齐射线列" }
func (Card2021007WizardVolleyLine) OnUseItem(ctx *EffectContext) error {
	companionCount := 0
	for _, card := range ctx.Engine.getAllFieldCards(ctx.Engine.State.Players[ctx.PlayerID]) {
		if card != nil && card.Card.IsCompanion() {
			companionCount++
		}
	}
	if companionCount < 7 {
		return nil
	}
	targets := ctx.Engine.friendlySkills(ctx.PlayerID, func(card *CardInstance) bool {
		return card.Card.IsSkill()
	})
	if len(targets) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "wizard_volley_line",
		"选择1个法术重置，其下一次范围变为前排", targets, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			card, _ := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, selected[0])
			if card == nil {
				return
			}
			resetInstance(card)
			card.Statuses["下一次范围前排"] = 1
		})
	return nil
}
