package game

import (
	"strings"

	"eraofarcane/model"
)

func addCardToDeck(ctx *EffectContext, cardNumber string, count int) {
	card := getCardDB()[cardNumber]
	if card == nil || count <= 0 {
		return
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	for i := 0; i < count; i++ {
		ps.Deck = append(ps.Deck, NewCardInstance(card, ctx.PlayerID, ctx.Engine.State.TurnNumber))
	}
	ctx.Engine.shuffleDeck(ctx.PlayerID)
	emitBatchEffect(ctx, "add_card_to_deck")
}

func addCardToOpponentDeck(ctx *EffectContext, cardNumber string, count int) {
	card := getCardDB()[cardNumber]
	if card == nil || count <= 0 {
		return
	}
	ps := ctx.Engine.State.Players[ctx.OpponentID]
	for i := 0; i < count; i++ {
		ps.Deck = append(ps.Deck, NewCardInstance(card, ctx.OpponentID, ctx.Engine.State.TurnNumber))
	}
	ctx.Engine.shuffleDeck(ctx.OpponentID)
	emitBatchEffect(ctx, "add_card_to_opponent_deck")
}

func allFriendlyUnits(ctx *EffectContext) []*CardInstance {
	return ctx.Engine.getAllFieldCards(ctx.Engine.State.Players[ctx.PlayerID])
}

type Card2411001AncientTreeHeart struct{ AlwaysActive }

func (Card2411001AncientTreeHeart) ID() string   { return "2411001" }
func (Card2411001AncientTreeHeart) Name() string { return "古树之心" }
func (Card2411001AncientTreeHeart) OnPerTurn(ctx *EffectContext) error {
	for _, unit := range allFriendlyUnits(ctx) {
		if unit != nil {
			healUnit(unit, 1)
			ctx.Engine.addElementsGainBonus(unit, ctx.PlayerID, model.ElementEarth, 1, ctx.Source)
		}
	}
	return nil
}

type Card2411002EarthsplitterSword struct{ AlwaysActive }

func (Card2411002EarthsplitterSword) ID() string   { return "2411002" }
func (Card2411002EarthsplitterSword) Name() string { return "裂地巨剑 阿托比斯" }
func (Card2411002EarthsplitterSword) OnConsume(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{Type: TempModSkillPowerBonus, Amount: 4, RemainingUses: 1, ExpiresTurn: ctx.Engine.State.TurnNumber + 1})
	return nil
}

type Card2421001KnowledgeTreeCare struct{ AlwaysActive }

func (Card2421001KnowledgeTreeCare) ID() string   { return "2421001" }
func (Card2421001KnowledgeTreeCare) Name() string { return "知识古树的关怀" }
func (Card2421001KnowledgeTreeCare) OnMasteryAchieved(ctx *EffectContext, level int) error {
	if ctx.Source == nil || ctx.Source.IsHorizontal {
		return nil
	}
	candidates := []map[string]any{candidateInfo(ctx.Source, "equipment", "own")}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "knowledge_tree_care",
		"你的卡牌达到精通，是否消耗知识古树的关怀抽1张牌并获得1地？", candidates, 0, 1,
		func(selected []string) {
			if len(selected) == 0 || selected[0] != ctx.Source.InstanceID || ctx.Source.IsHorizontal {
				return
			}
			ctx.Source.IsHorizontal = true
			ctx.Engine.drawCards(ctx.PlayerID, 1)
			ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementEarth: 1})
		})
	return nil
}

type Card2421007ParasiticTouch struct{ AlwaysActive }

func (Card2421007ParasiticTouch) ID() string      { return "2421007" }
func (Card2421007ParasiticTouch) Name() string    { return "寄生之触" }
func (Card2421007ParasiticTouch) MasteryMax() int { return 1 }
func (Card2421007ParasiticTouch) OnMastery(ctx *EffectContext, level int) error {
	if level == 1 {
		ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementEarth, 1, ctx.Source)
	}
	return nil
}

type Card2421013GeographyPrimer struct{ AlwaysActive }

func (Card2421013GeographyPrimer) ID() string   { return "2421013" }
func (Card2421013GeographyPrimer) Name() string { return "《地理学入门》" }
func (Card2421013GeographyPrimer) ModifyCardPlayCost(ctx *EffectContext, card *CardInstance, cost map[string]int) {
	if card == nil || card.Card.TotalCost() <= 5 {
		return
	}
	reduceCost(cost, model.ElementEarth, 2)
}

func firstActiveCardByNumber(e *Engine, ps *PlayerState, number string) *CardInstance {
	for _, fieldCard := range e.getAllFieldCards(ps) {
		if fieldCard != nil && fieldCard.Card != nil && fieldCard.Card.Number == number && !e.hasEffectiveStatus(fieldCard, StatusPetrify) {
			return fieldCard
		}
	}
	return nil
}

type Card2501001Shackle struct{ AlwaysActive }

func (Card2501001Shackle) ID() string   { return "2501001" }
func (Card2501001Shackle) Name() string { return "桎梏" }
func (Card2501001Shackle) RevealsOnDraw() bool {
	return true
}
func (Card2501001Shackle) OnSelfDraw(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	if ctx.ExtraData != nil {
		if initial, _ := ctx.ExtraData["initial_hand"].(bool); initial {
			return nil
		}
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	for i, card := range ps.Hand {
		if card == nil || card.InstanceID != ctx.Source.InstanceID {
			continue
		}
		ps.Hand = append(ps.Hand[:i], ps.Hand[i+1:]...)
		ps.Graveyard = append(ps.Graveyard, card)
		delete(ps.RevealedHand, card.InstanceID)
		ctx.Engine.emit(GameEvent{
			Type:   "discard",
			Player: ctx.PlayerID,
			Data:   map[string]any{"card": cardToInfo(card)},
		})
		ctx.Engine.drawCards(ctx.PlayerID, 1)
		return nil
	}
	return nil
}
func (Card2501001Shackle) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.drawCards(ctx.PlayerID, 1)
	return nil
}

type Card2511002ShiningShield struct{ AlwaysActive }

func (Card2511002ShiningShield) ID() string   { return "2511002" }
func (Card2511002ShiningShield) Name() string { return "辉之盾 闪耀" }
func (Card2511002ShiningShield) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if ctx.ExtraData != nil && ctx.ExtraData["purpose"] == string(skillPurposeDefend) {
		stats.PowerBonus += 2
	}
}
func (Card2511002ShiningShield) OnDefend(ctx *EffectContext) error {
	success, _ := ctx.ExtraData["defense_success"].(bool)
	defender, _ := ctx.ExtraData["defender"].(int)
	if !success || defender != ctx.PlayerID {
		return nil
	}
	for _, target := range ctx.Engine.getAllFieldCards(ctx.Engine.State.Players[ctx.OpponentID]) {
		if target == nil || target.Position == nil || !ctx.Engine.IsInSpellRange(ctx.PlayerID, target.Position.Col, target.Position.Row, false) {
			continue
		}
		ctx.Engine.addStatus(target, StatusStun, 1)
	}
	return nil
}

type Card2521002ShelterRune struct{ AlwaysActive }

func (Card2521002ShelterRune) ID() string   { return "2521002" }
func (Card2521002ShelterRune) Name() string { return "庇护符文" }
func (Card2521002ShelterRune) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{Type: "shelter_rune", Amount: 10, RemainingUses: 1, ExpiresTurn: ctx.Engine.State.TurnNumber + 2})
	return nil
}

type Card2521004HolySanctionScroll struct{ AlwaysActive }

func (Card2521004HolySanctionScroll) ID() string   { return "2521004" }
func (Card2521004HolySanctionScroll) Name() string { return "神圣制裁卷轴" }
func (Card2521004HolySanctionScroll) OnUseItem(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{Type: "holy_sanction", RemainingUses: 1, ExpiresTurn: ctx.Engine.State.TurnNumber + 2})
	return nil
}

type Card2601001PhantomPain struct{ AlwaysActive }

func (Card2601001PhantomPain) ID() string   { return "2601001" }
func (Card2601001PhantomPain) Name() string { return "幻痛" }
func (Card2601001PhantomPain) OnDefend(ctx *EffectContext) error {
	success, _ := ctx.ExtraData["defense_success"].(bool)
	defender, _ := ctx.ExtraData["defender"].(int)
	if !success || defender == ctx.PlayerID {
		return nil
	}
	weakenDefenseCards := func(key string) {
		skills, _ := ctx.ExtraData[key].([]*CardInstance)
		for _, skill := range skills {
			if skill != nil {
				ctx.Engine.addStatus(skill, StatusWeaken, 2)
			}
		}
	}
	weakenDefenseCards("defense_skills")
	weakenDefenseCards("defense_boosts")
	ctx.Engine.promptHeartPiercerAfterPhantomPain(ctx)
	return nil
}

func (e *Engine) promptHeartPiercerAfterPhantomPain(ctx *EffectContext) {
	if ctx == nil || ctx.Engine == nil {
		return
	}
	owner := ctx.PlayerID
	if firstActiveCardByNumber(e, e.State.Players[owner], "1611003") == nil {
		return
	}
	candidates := e.enemySpellCardCandidates(owner)
	if len(candidates) == 0 {
		return
	}
	opponent := e.State.Players[1-owner]
	e.SetPendingAction(owner, "heart_piercer_phantom_pain_extra", "\"穿心人\":幻痛触发,可以额外选择1个敌方法术虚弱2", candidates, 0, 1, func(selected []string) {
		if len(selected) == 0 {
			return
		}
		target := e.findSkill(opponent, selected[0])
		if canInstanceBeWeakened(target) {
			e.addStatus(target, StatusWeaken, 2)
		}
	})
}

func (e *Engine) enemySpellCardCandidates(playerID int) []map[string]any {
	ps := e.State.Players[1-playerID]
	candidates := make([]map[string]any, 0)
	for _, skill := range ps.Skills {
		if canInstanceBeWeakened(skill) {
			candidates = append(candidates, candidateInfo(skill, "skill", "enemy"))
		}
	}
	for _, card := range e.getAllFieldCards(ps) {
		if card == nil {
			continue
		}
		for _, skill := range card.BoundSkills {
			if canInstanceBeWeakened(skill) {
				candidates = append(candidates, candidateInfo(skill, "bound_skill", "enemy"))
			}
		}
	}
	return candidates
}

type Card2621002VoodooDoll struct{ AlwaysActive }

func (Card2621002VoodooDoll) ID() string   { return "2621002" }
func (Card2621002VoodooDoll) Name() string { return "巫毒娃娃" }
func (Card2621002VoodooDoll) OnEquip(ctx *EffectContext) error {
	ctx.Source.Statuses["暗影标记"] = 3
	candidates := voodooDollLinkCandidates(ctx)
	if len(candidates) < 2 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "voodoo_doll_link", "巫毒娃娃:选择法力范围内2个伙伴建立连结", candidates, 2, 2, func(selected []string) {
		clearVoodooDollLinks(ctx.Source)
		for _, id := range selected {
			ctx.Source.Statuses["巫毒连结:"+id] = 1
		}
	})
	return nil
}
func (Card2621002VoodooDoll) OnDamaged(ctx *EffectContext) error {
	if ctx.ExtraData != nil && ctx.ExtraData["damage_source"] == "voodoo_doll" {
		return nil
	}
	damage, _ := ctx.ExtraData["damage"].(int)
	if damage <= 0 || ctx.Source.Statuses["暗影标记"] < damage || ctx.Target == nil || !voodooDollIsLinked(ctx.Source, ctx.Target) {
		return nil
	}
	linked := voodooDollOtherLinkedUnit(ctx.Engine, ctx.Source, ctx.Target)
	if linked == nil {
		return nil
	}
	candidates := []map[string]any{candidateInfo(linked, "unit", voodooDollSide(ctx.PlayerID, linked))}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "voodoo_doll_damage", "巫毒娃娃:是否让连结的另一伙伴受到同等伤害", candidates, 0, 1, func(selected []string) {
		if len(selected) == 0 || selected[0] != linked.InstanceID || ctx.Source.Statuses["暗影标记"] < damage {
			return
		}
		ctx.Source.Statuses["暗影标记"] -= damage
		ctx.Engine.dealDamageWithExtra(linked, damage, linked.OwnerID, map[string]any{
			"damage_source": "voodoo_doll",
			"attacker":      ctx.PlayerID,
		})
	})
	return nil
}

func voodooDollLinkCandidates(ctx *EffectContext) []map[string]any {
	candidates := make([]map[string]any, 0)
	for ownerID, ps := range ctx.Engine.State.Players {
		side := "enemy"
		if ownerID == ctx.PlayerID {
			side = "own"
		}
		for col := 0; col < 3; col++ {
			for row := 0; row < 3; row++ {
				unit := ps.Units[col][row]
				if unit == nil || !unit.Card.IsCompanion() {
					continue
				}
				if ownerID != ctx.PlayerID && !ctx.Engine.IsInSpellRange(ctx.PlayerID, col, row, false) {
					continue
				}
				candidates = append(candidates, candidateInfo(unit, "unit", side))
			}
		}
	}
	return candidates
}

func clearVoodooDollLinks(doll *CardInstance) {
	for status := range doll.Statuses {
		if strings.HasPrefix(status, "巫毒连结:") {
			delete(doll.Statuses, status)
		}
	}
}

func voodooDollIsLinked(doll *CardInstance, unit *CardInstance) bool {
	return doll != nil && unit != nil && doll.Statuses["巫毒连结:"+unit.InstanceID] > 0
}

func voodooDollOtherLinkedUnit(e *Engine, doll *CardInstance, damaged *CardInstance) *CardInstance {
	for status := range doll.Statuses {
		if !strings.HasPrefix(status, "巫毒连结:") {
			continue
		}
		id := strings.TrimPrefix(status, "巫毒连结:")
		if id == "" || damaged != nil && id == damaged.InstanceID {
			continue
		}
		if unit := e.findUnitByInstanceID(id); unit != nil {
			return unit
		}
	}
	return nil
}

func voodooDollSide(playerID int, unit *CardInstance) string {
	if unit != nil && unit.OwnerID == playerID {
		return "own"
	}
	return "enemy"
}

type Card2621004ShadowVeil struct{ AlwaysActive }

func (Card2621004ShadowVeil) ID() string   { return "2621004" }
func (Card2621004ShadowVeil) Name() string { return "暗影帷幕" }
func (Card2621004ShadowVeil) OnSpellHit(ctx *EffectContext) error {
	if isEnemySpellCast(ctx) {
		ctx.Engine.State.Players[ctx.PlayerID].Hero.Statuses["引魔"] = 1
	}
	return nil
}

type Card2621010DragIntoAbyss struct{ AlwaysActive }

func (Card2621010DragIntoAbyss) ID() string   { return "2621010" }
func (Card2621010DragIntoAbyss) Name() string { return "拖入深渊" }
func (Card2621010DragIntoAbyss) OnUseItem(ctx *EffectContext) error {
	if ctx.Target == nil || ctx.Target.DamageTakenThisTurn <= 0 {
		return nil
	}
	damage := ctx.Target.DamageTakenThisTurn
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card.Position != nil && ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false)
	})
	ctx.Engine.SetPendingAction(ctx.PlayerID, "drag_into_abyss_target", "Drag Into Abyss: choose an enemy unit", candidates, 1, 1, func(selected []string) {
		target := selectedUnitFromCandidates(ctx.Engine, selected, candidates)
		if target != nil {
			ctx.Engine.dealDamage(target, damage, ctx.OpponentID)
		}
	})
	return nil
}

type Card2621011FrenzyRune struct{ AlwaysActive }

func (Card2621011FrenzyRune) ID() string   { return "2621011" }
func (Card2621011FrenzyRune) Name() string { return "狂乱符文" }
func (Card2621011FrenzyRune) OnUseItem(ctx *EffectContext) error {
	if ctx.Target == nil || ctx.Target.CurrentAttack <= 0 {
		return nil
	}
	attacker := ctx.Target
	candidates := ctx.Engine.adjacentUnitCandidatesForCounter(ctx.PlayerID, attacker)
	ctx.Engine.SetPendingAction(ctx.PlayerID, "frenzy_rune_target", "Frenzy Rune: choose an adjacent unit", candidates, 1, 1, func(selected []string) {
		if len(selected) == 0 {
			return
		}
		target := ctx.Engine.findFieldCardByInstance(ctx.Engine.State.Players[attacker.OwnerID], selected[0])
		ctx.Engine.resolveForcedUnitAttack(attacker.OwnerID, attacker, target, "frenzy_rune")
	})
	return nil
}

type Card2621012ShadowCloak struct{ AlwaysActive }

func (Card2621012ShadowCloak) ID() string   { return "2621012" }
func (Card2621012ShadowCloak) Name() string { return "暗影披风" }
func (Card2621012ShadowCloak) HasActiveSpellHitBeforeDamage(card *CardInstance) bool {
	return card != nil && card.Statuses[shadowCloakUsedStatus] == 0
}
func (Card2621012ShadowCloak) OnSpellHitBeforeDamage(ctx *EffectContext) error {
	if ctx == nil || ctx.ExtraData == nil || ctx.Source == nil {
		return nil
	}
	attacker, ok := ctx.ExtraData["attacker"].(int)
	if !ok || attacker == ctx.PlayerID {
		return nil
	}
	damagePtr, ok := ctx.ExtraData["damage_ptr"].(*int)
	if !ok || damagePtr == nil {
		return nil
	}
	*damagePtr = 0
	ctx.ExtraData["damage"] = 0
	ctx.Source.Statuses[shadowCloakUsedStatus] = 1
	return nil
}

const shadowCloakUsedStatus = "shadow_cloak_used"

type Card2621013WitchcraftRing struct{ AlwaysActive }

func (Card2621013WitchcraftRing) ID() string   { return "2621013" }
func (Card2621013WitchcraftRing) Name() string { return "巫术指环" }
func (Card2621013WitchcraftRing) OnPerTurn(ctx *EffectContext) error {
	for _, skill := range ctx.Engine.State.Players[ctx.OpponentID].Skills {
		if skill != nil && skill.Statuses[StatusWeaken] > 0 {
			ctx.Engine.addStatus(skill, StatusWeaken, 1)
		}
	}
	return nil
}

type Card3021006InsightEye struct{ AlwaysActive }

func (Card3021006InsightEye) ID() string   { return "3021006" }
func (Card3021006InsightEye) Name() string { return "洞察之眼" }
func (Card3021006InsightEye) OnSpellCast(ctx *EffectContext) error {
	if !isFriendlySpellCast(ctx) || !isSpellBeingCast(ctx) {
		return nil
	}
	candidates := ctx.Engine.enemyEquipment(ctx.PlayerID, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "insight_eye_destroy_equipment",
		"选择1张敌方盖放的卡牌摧毁", candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			ctx.Engine.destroyEnemyEquipment(ctx.PlayerID, selected[0])
		})
	return nil
}

type Card3021010Dispel struct{ AlwaysActive }

func (Card3021010Dispel) ID() string   { return "3021010" }
func (Card3021010Dispel) Name() string { return "解咒" }

type Card3021011OverlordSanction struct{ AlwaysActive }

func (Card3021011OverlordSanction) ID() string   { return "3021011" }
func (Card3021011OverlordSanction) Name() string { return "统御者的制裁" }
func (Card3021011OverlordSanction) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	stats.DamageBonus += 1
}

type Card3121015BurningWind struct{ AlwaysActive }

func (Card3121015BurningWind) ID() string   { return "3121015" }
func (Card3121015BurningWind) Name() string { return "焚风" }
func (Card3121015BurningWind) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if ctx.ExtraData != nil {
		purpose, _ := ctx.ExtraData["purpose"].(string)
		if !isBoostPurpose(skillPurpose(purpose)) {
			return
		}
		stats.Pierce = true
	}
}

type Card3221008IceDissolve struct{ AlwaysActive }

func (Card3221008IceDissolve) ID() string   { return "3221008" }
func (Card3221008IceDissolve) Name() string { return "冰封消解" }

func (Card3221008IceDissolve) CanReactToSpell(ctx *EffectContext, spell *SpellCast) bool {
	return ctx != nil && spell != nil && spell.AttackerID != ctx.PlayerID && spell.TotalPower > 0
}

func (Card3221008IceDissolve) OnSpellReaction(ctx *EffectContext, spell *SpellCast) error {
	sources := positiveSpellPowerSources(spell)
	if len(sources) == 0 {
		return nil
	}
	if len(sources) == 1 {
		applyIceDissolveToSource(ctx, spell, sources[0].InstanceID)
		return nil
	}
	candidates := make([]map[string]any, 0, len(sources))
	for _, source := range sources {
		candidates = append(candidates, map[string]any{
			"instance_id": source.InstanceID,
			"name":        source.CardName,
			"power":       source.Power,
			"is_main":     source.IsMain,
		})
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "ice_dissolve_power_source", "冰封消解:选择1个法术威力变为0", candidates, 1, 1, func(selected []string) {
		if len(selected) == 0 {
			return
		}
		applyIceDissolveToSource(ctx, spell, selected[0])
	})
	return nil
}

func positiveSpellPowerSources(spell *SpellCast) []SpellPowerSource {
	if spell == nil {
		return nil
	}
	sources := make([]SpellPowerSource, 0, len(spell.PowerSources))
	for _, source := range spell.PowerSources {
		if source.Power > 0 && source.InstanceID != "" {
			sources = append(sources, source)
		}
	}
	if len(sources) == 0 && spell.TotalPower > 0 {
		source := spellPowerSourceForCard(spell.Skill, spell.TotalPower, true)
		if source.InstanceID == "" {
			source.InstanceID = "__total_power__"
		}
		sources = append(sources, source)
	}
	return sources
}

func applyIceDissolveToSource(ctx *EffectContext, spell *SpellCast, instanceID string) {
	if ctx == nil || spell == nil || instanceID == "" {
		return
	}
	removed := 0
	for i := range spell.PowerSources {
		if spell.PowerSources[i].InstanceID != instanceID {
			continue
		}
		removed = spell.PowerSources[i].Power
		spell.PowerSources[i].Power = 0
		break
	}
	if removed == 0 && ((spell.Skill != nil && spell.Skill.InstanceID == instanceID) || instanceID == "__total_power__") {
		removed = spell.TotalPower
	}
	spell.TotalPower = max(spell.TotalPower-removed, 0)
	ctx.Engine.emit(GameEvent{
		Type:   "spell_reaction",
		Player: -1,
		Data: map[string]any{
			"player":             ctx.PlayerID,
			"card":               cardToInfo(ctx.Source),
			"effect":             "power_zero",
			"power":              spell.TotalPower,
			"power_source_id":    instanceID,
			"power_source_power": removed,
		},
	})
}

type Card3221010WaterPhantom struct{ AlwaysActive }

func (Card3221010WaterPhantom) ID() string   { return "3221010" }
func (Card3221010WaterPhantom) Name() string { return "水幻影" }
func (Card3221010WaterPhantom) OnSpellCast(ctx *EffectContext) error {
	if !isFriendlySpellCast(ctx) || ctx.Target != nil {
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card.Card.IsCompanion() && card.Card.Category == model.ElementWater && card.EnterTurn == ctx.Engine.State.TurnNumber
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "water_phantom_copy_target", "水幻影:选择本回合你召唤的1个水纹伙伴", candidates, 1, 1, func(selected []string) {
		target := selectedUnitFromCandidates(ctx.Engine, selected, candidates)
		if target == nil {
			return
		}
		positions := ctx.Engine.friendlyEmptyUnitPositions(ctx.PlayerID)
		if len(positions) == 0 {
			return
		}
		targetID := target.InstanceID
		ctx.Engine.SetPendingAction(ctx.PlayerID, "water_phantom_copy_position", "水幻影:选择复制体的入场位置", positions, 1, 1, func(posSelected []string) {
			pos, ok := positionFromSelectionID(firstSelected(posSelected))
			if !ok {
				return
			}
			ctx.Engine.summonWaterPhantomCopy(ctx.PlayerID, targetID, pos)
		})
	})
	return nil
}

func (e *Engine) summonWaterPhantomCopy(playerID int, targetID string, pos Position) *CardInstance {
	if e == nil || !pos.Valid() || playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	ps := e.State.Players[playerID]
	if ps.Units[pos.Col][pos.Row] != nil {
		return nil
	}
	source := e.findUnitByInstanceID(targetID)
	if source == nil || source.OwnerID != playerID || source.Card == nil || !source.Card.IsCompanion() || source.Card.Category != model.ElementWater || source.EnterTurn != e.State.TurnNumber {
		return nil
	}
	cardCopy := *source.Card
	cardCopy.Life = 1
	copyUnit := NewCardInstance(&cardCopy, playerID, e.State.TurnNumber)
	copyUnit.Position = &Position{Col: pos.Col, Row: pos.Row}
	copyUnit.IsHorizontal = true
	copyUnit.CurrentLife = 1
	copyUnit.Statuses["水幻影复制"] = 1
	ps.Units[pos.Col][pos.Row] = copyUnit
	e.emit(GameEvent{
		Type:   "summon",
		Player: -1,
		Data: map[string]any{
			"player":   playerID,
			"card":     cardToInfo(copyUnit),
			"position": pos,
			"via":      "water_phantom",
		},
	})
	e.triggerEffects(TriggerOnEnter, copyUnit, nil, nil)
	enterData := map[string]any{"entered_player": playerID}
	e.triggerFieldEffectsWithData(TriggerOnUnitEnter, playerID, copyUnit, enterData)
	e.triggerFieldEffectsWithData(TriggerOnUnitEnter, 1-playerID, copyUnit, enterData)
	return copyUnit
}

type Card3321001LightningChain struct{ AlwaysActive }

func (Card3321001LightningChain) ID() string   { return "3321001" }
func (Card3321001LightningChain) Name() string { return "闪电链" }
func (Card3321001LightningChain) SpellDamage(ctx *EffectContext) int {
	return 1
}

type Card3321008WindHole struct{ AlwaysActive }

func (Card3321008WindHole) ID() string   { return "3321008" }
func (Card3321008WindHole) Name() string { return "风洞" }
func (Card3321008WindHole) CanReactToSpell(ctx *EffectContext, spell *SpellCast) bool {
	return ctx != nil && spell != nil && spell.AttackerID != ctx.PlayerID && spellArea(spell.Skill) == SpellAreaSingle
}
func (Card3321008WindHole) OnSpellReaction(ctx *EffectContext, spell *SpellCast) error {
	ctx.Engine.cancelPendingSpell(ctx.PlayerID, ctx.Source, "wind_hole")
	return nil
}

type Card3421014ThousandMileQuicksand struct{ AlwaysActive }

func (Card3421014ThousandMileQuicksand) ID() string   { return "3421014" }
func (Card3421014ThousandMileQuicksand) Name() string { return "千里流沙" }
func (Card3421014ThousandMileQuicksand) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	stats.PowerBonus += 1
}

type Card3521011LightShelter struct{ AlwaysActive }

func (Card3521011LightShelter) ID() string   { return "3521011" }
func (Card3521011LightShelter) Name() string { return "光之庇护" }
func (Card3521011LightShelter) AllowsFriendlySpellTarget() bool {
	return true
}
func (Card3521011LightShelter) OnSpellHit(ctx *EffectContext) error {
	if !isOwnSpellHit(ctx) {
		return nil
	}
	target := ctx.Target
	if target != nil && target.Card.IsCompanion() {
		target.Statuses["防止致命"] = 1
	}
	return nil
}

type Card3621015Siphon struct{ AlwaysActive }

func (Card3621015Siphon) ID() string   { return "3621015" }
func (Card3621015Siphon) Name() string { return "虹吸" }
func (Card3621015Siphon) CanReactToSpell(ctx *EffectContext, spell *SpellCast) bool {
	return ctx != nil && spell != nil && spell.AttackerID != ctx.PlayerID && spell.Skill != nil
}
func (Card3621015Siphon) OnSpellReaction(ctx *EffectContext, spell *SpellCast) error {
	defenderID := ctx.PlayerID
	affected := ctx.Engine.spellAffectedUnitsWithExtraTargets(defenderID, spell.Skill, spell.Target, spell.ExtraTargets)
	if len(affected) == 0 {
		return nil
	}
	baseDamage := max(spell.Skill.Card.Attack+spell.Skill.AttackBonus, 0)
	if override, ok := globalRegistry.SpellDamage(spell.Skill.Card.Number, &EffectContext{
		Engine:     ctx.Engine,
		Source:     spell.Skill,
		Target:     affected[0],
		PlayerID:   spell.AttackerID,
		OpponentID: defenderID,
		ExtraData:  map[string]any{"target": spell.Target},
	}); ok {
		baseDamage = max(override, 0)
	}
	damage := ctx.Engine.effectiveSpellDamage(spell.AttackerID, spell.Skill, baseDamage, spell.BoostSkills)
	for _, unit := range affected {
		healUnit(unit, damage)
	}
	ctx.Engine.cancelPendingSpell(ctx.PlayerID, ctx.Source, "siphon")
	return nil
}

func (e *Engine) spellAffectedUnitsWithExtraTargets(defenderID int, skill *CardInstance, target SpellTarget, extraTargets []SpellTarget) []*CardInstance {
	affected := e.spellAffectedUnits(defenderID, skill, target)
	for _, extraTarget := range extraTargets {
		if extraTarget.Type != "unit" || !extraTarget.Position.Valid() {
			continue
		}
		unit := e.State.Players[defenderID].Units[extraTarget.Position.Col][extraTarget.Position.Row]
		if unit == nil {
			continue
		}
		alreadyIncluded := false
		for _, existing := range affected {
			if existing == unit {
				alreadyIncluded = true
				break
			}
		}
		if !alreadyIncluded {
			affected = append(affected, unit)
		}
	}
	return affected
}

type Card4111001Longjuanhuo struct{ AlwaysActive }

func (Card4111001Longjuanhuo) ID() string   { return "4111001" }
func (Card4111001Longjuanhuo) Name() string { return "掌门 龙卷火" }
func (Card4111001Longjuanhuo) OnTurnStart(ctx *EffectContext) error {
	if ctx.Source.Statuses["开局技能"] == 0 {
		addSkillToPool(ctx, "3101002")
		ctx.Source.Statuses["开局技能"] = 1
	}
	return nil
}

type Card4211002Volport struct{ AlwaysActive }

func (Card4211002Volport) ID() string   { return "4211002" }
func (Card4211002Volport) Name() string { return "大贤者 沃尔波特" }
func (Card4211002Volport) OnTurnStart(ctx *EffectContext) error {
	if ctx.Source.Statuses["开局技能"] == 0 {
		addSkillToPool(ctx, "3201001")
		ctx.Source.Statuses["开局技能"] = 1
	}
	return nil
}

type Card4311002Raven struct{ AlwaysActive }

func (Card4311002Raven) ID() string   { return "4311002" }
func (Card4311002Raven) Name() string { return "\"渡鸦\" 睿文" }

type Card4411002Andrew struct{ AlwaysActive }

func (Card4411002Andrew) ID() string   { return "4411002" }
func (Card4411002Andrew) Name() string { return "大法师 安德鲁" }
func (Card4411002Andrew) OnEnter(ctx *EffectContext) error {
	addCardToDeck(ctx, "1401002", 1)
	return nil
}

type Card4411003ProfessorMaggie struct{ AlwaysActive }

func (Card4411003ProfessorMaggie) ID() string   { return "4411003" }
func (Card4411003ProfessorMaggie) Name() string { return "麦吉教授" }
func (Card4411003ProfessorMaggie) ModifyCardPlayCost(ctx *EffectContext, card *CardInstance, cost map[string]int) {
	if ctx.Source.Statuses["麦吉折扣"] == 0 && card != nil && card.Card.TotalCost() > 5 {
		reduceCost(cost, model.ElementEarth, 2)
	}
}
func (Card4411003ProfessorMaggie) OnCardPlayCostPaid(ctx *EffectContext, card *CardInstance) {
	if ctx.Source.Statuses["麦吉折扣"] == 0 && card != nil && card.Card.TotalCost() > 5 {
		ctx.Source.Statuses["麦吉折扣"] = 1
	}
}

type Card4511002Ailimer struct{ AlwaysActive }

func (Card4511002Ailimer) ID() string   { return "4511002" }
func (Card4511002Ailimer) Name() string { return "神之眷子 爱里默" }
func (Card4511002Ailimer) OnEnter(ctx *EffectContext) error {
	ailimerShuffleShacklesOnce(ctx)
	ailimerUnlockIfShacklesCleared(ctx)
	return nil
}
func (Card4511002Ailimer) OnTurnStart(ctx *EffectContext) error {
	ailimerShuffleShacklesOnce(ctx)
	ailimerUnlockIfShacklesCleared(ctx)
	return nil
}
func (Card4511002Ailimer) OnUseItem(ctx *EffectContext) error {
	if ctx.Target != nil && ctx.Target.Card != nil && ctx.Target.Card.Number == "2501001" {
		ailimerUnlockIfShacklesCleared(ctx)
	}
	return nil
}
func (Card4511002Ailimer) HasActiveUltimate(card *CardInstance) bool {
	return card != nil && card.Statuses["爱里默解放"] > 0
}
func (Card4511002Ailimer) OnUltimate(ctx *EffectContext) error {
	candidates := ctx.Engine.nonHeroFieldCardCandidates(ctx.PlayerID)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "ailimer_remove_cards", "爱里默:移除最多3张非人物场上卡牌", candidates, 0, min(3, len(candidates)), func(selected []string) {
		for _, id := range selected {
			ctx.Engine.removeFieldCardFromGameByID(id)
		}
	})
	return nil
}

func ailimerShuffleShacklesOnce(ctx *EffectContext) {
	if ctx == nil || ctx.Source == nil || ctx.Source.Statuses["桎梏已洗入"] > 0 {
		return
	}
	addCardToOpponentDeck(ctx, "2501001", 5)
	ctx.Source.Statuses["桎梏已洗入"] = 1
}

func ailimerUnlockIfShacklesCleared(ctx *EffectContext) {
	if ctx == nil || ctx.Source == nil || ctx.Source.Statuses["爱里默解放"] > 0 {
		return
	}
	if countCardInstancesByNumber(ctx.Engine.State.Players[ctx.OpponentID].Graveyard, "2501001") < 5 {
		return
	}
	ctx.Source.Statuses["爱里默解放"] = 1
	ctx.Engine.emit(GameEvent{
		Type:   "effect_trigger",
		Player: -1,
		Data: map[string]any{
			"player": ctx.PlayerID,
			"card":   cardToInfo(ctx.Source),
			"effect": "ailimer_unlocked",
		},
	})
}

func countCardInstancesByNumber(cards []*CardInstance, number string) int {
	count := 0
	for _, card := range cards {
		if card != nil && card.Card != nil && card.Card.Number == number {
			count++
		}
	}
	return count
}

type Card4511003Lexia struct{ AlwaysActive }

func (Card4511003Lexia) ID() string   { return "4511003" }
func (Card4511003Lexia) Name() string { return "骑士团长 蕾曦娅" }
func (Card4511003Lexia) OnTurnStart(ctx *EffectContext) error {
	if ctx.Source.Statuses["团结希望"] == 0 {
		if !replaceSkillInPool(ctx, "3521007", "3501001") {
			addSkillToPool(ctx, "3501001")
		}
		ctx.Source.Statuses["团结希望"] = 1
	}
	return nil
}

type Card4611003Jieying struct{ AlwaysActive }

func (Card4611003Jieying) ID() string   { return "4611003" }
func (Card4611003Jieying) Name() string { return "咒言师 结影" }
func (Card4611003Jieying) OnTurnStart(ctx *EffectContext) error {
	if ctx.Source.Statuses["咒言书"] == 0 {
		addCardToDeck(ctx, "2601002", 3)
		ctx.Source.Statuses["咒言书"] = 1
	}
	return nil
}
