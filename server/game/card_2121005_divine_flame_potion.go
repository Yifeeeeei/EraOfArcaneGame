package game

type Card2121005DivineFlamePotion struct{ AlwaysActive }

func (Card2121005DivineFlamePotion) ID() string { return "2121005" }

func (Card2121005DivineFlamePotion) Name() string { return "神炎魔咒药剂" }

func (Card2121005DivineFlamePotion) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModSkillPowerBonus,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		Amount:           2,
		ExpiresTurn:      ctx.Engine.State.TurnNumber + 2,
	})
	hero := ctx.Engine.State.Players[ctx.PlayerID].Hero
	if hero != nil && ctx.Engine.addStatus(hero, StatusBurn, 1) {
		ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
			"source": cardToInfo(ctx.Source),
			"target": cardToInfo(hero),
			"effect": "apply_status",
			"status": StatusBurn,
			"amount": 1,
		}})
	}
	return nil
}
