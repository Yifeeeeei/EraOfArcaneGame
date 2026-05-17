package game

type Card2521001LifePotion struct{ AlwaysActive }

func (Card2521001LifePotion) ID() string   { return "2521001" }
func (Card2521001LifePotion) Name() string { return "生命药剂" }

func (Card2521001LifePotion) OnUseItem(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card.CurrentLife < card.Card.Life
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "heal_unit",
		"选择1个友方单位回复2点生命",
		candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			card, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, selected[0])
			if card == nil || zone != "unit" {
				return
			}
			card.CurrentLife += 2
			if card.CurrentLife > card.Card.Life {
				card.CurrentLife = card.Card.Life
			}
			ctx.Engine.emit(GameEvent{
				Type:   "effect_trigger",
				Player: -1,
				Data: map[string]any{
					"source": cardToInfo(ctx.Source),
					"effect": "heal",
					"amount": 2,
					"target": cardToInfo(card),
				},
			})
		})
	return nil
}
