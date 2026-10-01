package game

type Card1521105RadiantCityPriest struct{ AlwaysActive }

func (Card1521105RadiantCityPriest) ID() string { return "1521105" }

func (Card1521105RadiantCityPriest) Name() string { return "辉之都祭司" }

func (Card1521105RadiantCityPriest) OnSpellHitBeforeDamage(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil || ctx.Source.UltimateUsed {
		return nil
	}
	attacker, ok := ctx.ExtraData["attacker"].(int)
	damagePtr, hasDamage := ctx.ExtraData["damage_ptr"].(*int)
	if !ok || attacker == ctx.PlayerID || !hasDamage || damagePtr == nil || *damagePtr <= 0 {
		return nil
	}
	source := ctx.Source
	ctx.Engine.SetPendingAction(ctx.PlayerID, "radiant_city_priest_convert_damage_to_burn",
		"辉之都祭司:是否将该法术伤害转为点燃", []map[string]any{candidateInfo(source, "unit", "own")}, 0, 1,
		func(selected []string) {
			if len(selected) == 0 || source.UltimateUsed || !ctx.Engine.cardStillOnField(source) {
				return
			}
			burn := *damagePtr
			if burn <= 0 {
				return
			}
			for _, target := range spellHitAffectedUnitsFromData(ctx) {
				if target != nil && ctx.Engine.unitStillOnField(target) {
					ctx.Engine.addStatus(target, StatusBurn, burn)
				}
			}
			*damagePtr = 0
			ctx.ExtraData["damage"] = 0
			source.UltimateUsed = true
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source": cardToInfo(source),
				"effect": "convert_spell_damage_to_burn",
				"burn":   burn,
			}})
		})
	return nil
}
