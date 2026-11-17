package game

type Card4211003CrystalHeart struct{ AlwaysActive }

func (Card4211003CrystalHeart) ID() string { return "4211003" }

func (Card4211003CrystalHeart) Name() string { return "凛冬城主 水晶心" }

func (Card4211003CrystalHeart) OnUltimate(ctx *EffectContext) error {
	ctx.Source.Statuses["水晶心冻结法术"] = 1
	return nil
}

func (Card4211003CrystalHeart) OnSpellHit(ctx *EffectContext) error {
	if ctx.Source.Statuses["水晶心冻结法术"] <= 0 {
		return nil
	}
	if !isFriendlySpellHit(ctx) {
		return nil
	}
	if targets, ok := ctx.ExtraData["affected_units"].([]*CardInstance); ok {
		for _, target := range targets {
			if target != nil {
				ctx.Engine.addStatus(target, StatusFreeze, 1)
			}
		}
		return nil
	}
	if ctx.Target != nil && !ctx.Target.Card.IsSkill() {
		ctx.Engine.addStatus(ctx.Target, StatusFreeze, 1)
	}
	return nil
}

func (Card4211003CrystalHeart) OnTurnEnd(ctx *EffectContext) error {
	delete(ctx.Source.Statuses, "水晶心冻结法术")
	return nil
}
