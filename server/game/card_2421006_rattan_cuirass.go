package game

type Card2421006RattanCuirass struct{ AlwaysActive }

func (Card2421006RattanCuirass) ID() string   { return "2421006" }
func (Card2421006RattanCuirass) Name() string { return "磐藤胸甲" }

func (Card2421006RattanCuirass) OnEnter(ctx *EffectContext) error {
	hero := ctx.Engine.State.Players[ctx.PlayerID].Hero
	if hero == nil {
		return nil
	}
	ctx.Engine.gainLife(hero, 2, ctx.Source)
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(ctx.Source),
		"target": cardToInfo(hero),
		"effect": "modify_life",
		"amount": 2,
	}})
	return nil
}
