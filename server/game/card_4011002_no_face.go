package game

type Card4011002NoFace struct{}

func (Card4011002NoFace) ID() string   { return "4011002" }
func (Card4011002NoFace) Name() string { return "\"无面\"" }
func (Card4011002NoFace) OnUnitEnter(ctx *EffectContext) error {
	if ctx.Target == nil || ctx.Source == nil {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	for _, c := range ctx.Engine.getAllFieldCards(ps) {
		if c == ctx.Target {
			continue
		}
		if c.Card.Category == ctx.Target.Card.Category && ps.Hero != nil {
			ctx.Engine.dealDamage(ps.Hero, 1, ctx.PlayerID)
			ctx.Engine.emit(GameEvent{
				Type:   "effect_trigger",
				Player: -1,
				Data: map[string]any{
					"source": cardToInfo(ctx.Source),
					"effect": "same_element_penalty",
					"damage": 1,
				},
			})
			return nil
		}
	}
	return nil
}
