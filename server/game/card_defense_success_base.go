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

type Card3121101SummonFireSnake struct{ AlwaysActive }

func (Card3121101SummonFireSnake) ID() string   { return "3121101" }
func (Card3121101SummonFireSnake) Name() string { return "唤灵术 火蛇" }

func (Card3121101SummonFireSnake) OnDefend(ctx *EffectContext) error {
	return promptDefenseSuccessDamage(ctx, "summon_fire_snake_defense_damage", "唤灵术 火蛇:选择法力范围内1个敌人造成1点伤害", 1)
}

type Card2121109SummonBlazingHoundScroll struct{ AlwaysActive }

func (Card2121109SummonBlazingHoundScroll) ID() string   { return "2121109" }
func (Card2121109SummonBlazingHoundScroll) Name() string { return "唤灵术卷轴 烈焰犬" }

func (Card2121109SummonBlazingHoundScroll) OnDefend(ctx *EffectContext) error {
	return promptDefenseSuccessDamage(ctx, "blazing_hound_scroll_defense_damage", "唤灵术卷轴 烈焰犬:选择法力范围内1个敌人造成2点伤害", 2)
}

type Card3221102SummonFloodDragon struct{ AlwaysActive }

func (Card3221102SummonFloodDragon) ID() string   { return "3221102" }
func (Card3221102SummonFloodDragon) Name() string { return "唤灵术 蛟龙" }

func (Card3221102SummonFloodDragon) OnDefend(ctx *EffectContext) error {
	if !defenseSucceeded(ctx) {
		return nil
	}
	for _, candidate := range defenseSuccessDamageCandidates(ctx) {
		id, _ := candidate["instance_id"].(string)
		target := validDefenseSuccessDamageTarget(ctx, id)
		if target == nil {
			continue
		}
		ctx.Engine.ApplyDamage(DamageRequest{Target: target, Amount: 1, Kind: "effect", SourcePlayer: ctx.PlayerID, SourceKnown: true, Source: ctx.Source})
	}
	return nil
}

func promptDefenseSuccessDamage(ctx *EffectContext, actionType, prompt string, amount int) error {
	if !defenseSucceeded(ctx) {
		return nil
	}
	candidates := defenseSuccessDamageCandidates(ctx)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, actionType, prompt, candidates, 1, 1, func(selected []string) {
		target := validDefenseSuccessDamageTarget(ctx, firstSelected(selected))
		if target == nil {
			return
		}
		ctx.Engine.ApplyDamage(DamageRequest{Target: target, Amount: amount, Kind: "effect", SourcePlayer: ctx.PlayerID, SourceKnown: true, Source: ctx.Source})
	})
	return nil
}

func defenseSucceeded(ctx *EffectContext) bool {
	if ctx == nil || ctx.Engine == nil || ctx.ExtraData == nil {
		return false
	}
	success, _ := ctx.ExtraData["defense_success"].(bool)
	return success
}

func defenseSuccessDamageCandidates(ctx *EffectContext) []map[string]any {
	return ctx.Engine.enemyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return defenseSuccessDamageTargetInRange(ctx, card)
	})
}

func validDefenseSuccessDamageTarget(ctx *EffectContext, instanceID string) *CardInstance {
	target := ctx.Engine.findUnitByInstanceID(instanceID)
	if target == nil || target.OwnerID != ctx.OpponentID || !defenseSuccessDamageTargetInRange(ctx, target) {
		return nil
	}
	return target
}

func defenseSuccessDamageTargetInRange(ctx *EffectContext, card *CardInstance) bool {
	return card != nil && card.Card != nil && card.Position != nil && ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false)
}

type Card3121102LionGuardian struct{ AlwaysActive }

func (Card3121102LionGuardian) ID() string   { return "3121102" }
func (Card3121102LionGuardian) Name() string { return "雄狮之守护" }

func (Card3121102LionGuardian) OnDefend(ctx *EffectContext) error {
	if ctx.ExtraData == nil {
		return nil
	}
	success, _ := ctx.ExtraData["defense_success"].(bool)
	if !success {
		return nil
	}
	for _, skill := range ctx.Engine.State.Players[ctx.PlayerID].Skills {
		if skill == nil || skill == ctx.Source || skill.Card == nil || skill.Card.Category != "火" {
			continue
		}
		skill.PowerBonus++
	}
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

type Card3321104GatherMomentum struct{ AlwaysActive }

func (Card3321104GatherMomentum) ID() string   { return "3321104" }
func (Card3321104GatherMomentum) Name() string { return "收势" }

func (Card3321104GatherMomentum) OnDefend(ctx *EffectContext) error {
	if ctx.ExtraData == nil {
		return nil
	}
	success, _ := ctx.ExtraData["defense_success"].(bool)
	if !success {
		return nil
	}
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModNextAttackSpellPowerBonus,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		Amount:           3,
		RemainingUses:    1,
	})
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
