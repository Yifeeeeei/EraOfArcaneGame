package game

type Card3121013FireBacklash struct{ AlwaysActive }

func (Card3121013FireBacklash) ID() string   { return "3121013" }
func (Card3121013FireBacklash) Name() string { return "烈焰反噬" }

func (Card3121013FireBacklash) OnDefend(ctx *EffectContext) error {
	success, _ := ctx.ExtraData["defense_success"].(bool)
	if !success {
		return nil
	}
	hero := ctx.Engine.State.Players[ctx.OpponentID].Hero
	if hero == nil {
		return nil
	}
	if !ctx.Engine.addStatus(hero, StatusBurn, 1) {
		return nil
	}
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(ctx.Source),
		"target": cardToInfo(hero),
		"effect": "apply_status",
		"status": StatusBurn,
		"amount": 1,
	}})
	return nil
}

type Card3221014IceField struct{ AlwaysActive }

func (Card3221014IceField) ID() string   { return "3221014" }
func (Card3221014IceField) Name() string { return "坚冰领域" }

func (Card3221014IceField) OnDefend(ctx *EffectContext) error {
	success, _ := ctx.ExtraData["defense_success"].(bool)
	if !success {
		return nil
	}
	opponent := ctx.Engine.State.Players[ctx.OpponentID]
	row := opponent.GetFrontRow()
	if row < 0 {
		return nil
	}
	for col := 0; col < 3; col++ {
		unit := opponent.Units[col][row]
		if unit == nil {
			continue
		}
		if !ctx.Engine.addStatus(unit, StatusFreeze, 1) {
			continue
		}
		ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
			"source": cardToInfo(ctx.Source),
			"target": cardToInfo(unit),
			"effect": "apply_status",
			"status": StatusFreeze,
			"amount": 1,
		}})
	}
	return nil
}

type Card3621014Karma struct{ AlwaysActive }

func (Card3621014Karma) ID() string   { return "3621014" }
func (Card3621014Karma) Name() string { return "业障" }

func (Card3621014Karma) OnDefend(ctx *EffectContext) error {
	success, _ := ctx.ExtraData["defense_success"].(bool)
	if !success {
		return nil
	}
	attackSkill, _ := ctx.ExtraData["attack_skill"].(*CardInstance)
	if attackSkill == nil {
		return nil
	}
	if !ctx.Engine.addStatus(attackSkill, StatusWeaken, 2) {
		return nil
	}
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(ctx.Source),
		"target": cardToInfo(attackSkill),
		"effect": "apply_status",
		"status": StatusWeaken,
		"amount": 2,
	}})
	return nil
}
