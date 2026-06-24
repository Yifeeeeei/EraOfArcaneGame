package game

type Card2221008WaterformBindScroll struct{ AlwaysActive }

func (Card2221008WaterformBindScroll) ID() string   { return "2221008" }
func (Card2221008WaterformBindScroll) Name() string { return "水形之束卷轴" }

func (Card2221008WaterformBindScroll) OnSpellHit(ctx *EffectContext) error {
	if !isOwnSpellHit(ctx) {
		return nil
	}
	target := ctx.Target
	if target == nil || !target.Card.IsCompanion() {
		return nil
	}
	target.IsHorizontal = true
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(ctx.Source),
		"target": cardToInfo(target),
		"effect": "consume_target",
	}})
	return nil
}
