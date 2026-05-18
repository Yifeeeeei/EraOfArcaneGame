package game

import "eraofarcane/model"

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

func firstEnemyUnit(ctx *EffectContext) *CardInstance {
	return firstUnitFromCandidates(ctx.Engine, ctx.PlayerID, ctx.Engine.enemyUnits(ctx.PlayerID, true, nil))
}

func firstFriendlyCompanion(ctx *EffectContext) *CardInstance {
	return firstUnitFromCandidates(ctx.Engine, ctx.PlayerID, ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card.Card.IsCompanion()
	}))
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
	if card != nil && card.Card.TotalCost() > 5 {
		reduceCost(cost, model.ElementEarth, 2)
	}
}

type Card2501001Shackle struct{ AlwaysActive }

func (Card2501001Shackle) ID() string   { return "2501001" }
func (Card2501001Shackle) Name() string { return "桎梏" }
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
func (Card2511002ShiningShield) OnPerTurn(ctx *EffectContext) error {
	for _, target := range ctx.Engine.getAllFieldCards(ctx.Engine.State.Players[ctx.OpponentID]) {
		target.Statuses[StatusStun]++
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
func (Card2601001PhantomPain) OnPerTurn(ctx *EffectContext) error {
	for _, skill := range ctx.Engine.State.Players[ctx.OpponentID].Skills {
		if skill != nil {
			skill.Statuses[StatusWeaken] += 2
		}
	}
	return nil
}

type Card2621002VoodooDoll struct{ AlwaysActive }

func (Card2621002VoodooDoll) ID() string   { return "2621002" }
func (Card2621002VoodooDoll) Name() string { return "巫毒娃娃" }
func (Card2621002VoodooDoll) OnEquip(ctx *EffectContext) error {
	ctx.Source.Statuses["暗影标记"] = 3
	return nil
}
func (Card2621002VoodooDoll) OnDamaged(ctx *EffectContext) error {
	if ctx.Source.Statuses["暗影标记"] <= 0 || ctx.Target == nil || ctx.Target.OwnerID != ctx.PlayerID {
		return nil
	}
	if target := firstEnemyUnit(ctx); target != nil {
		ctx.Source.Statuses["暗影标记"]--
		ctx.Engine.dealDamage(target, 1, ctx.OpponentID)
	}
	return nil
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
	if target := firstEnemyUnit(ctx); target != nil {
		ctx.Engine.dealDamage(target, 2, ctx.OpponentID)
	}
	return nil
}

type Card2621011FrenzyRune struct{ AlwaysActive }

func (Card2621011FrenzyRune) ID() string   { return "2621011" }
func (Card2621011FrenzyRune) Name() string { return "狂乱符文" }
func (Card2621011FrenzyRune) OnUseItem(ctx *EffectContext) error {
	if target := firstEnemyUnit(ctx); target != nil {
		target.Statuses[StatusStun]++
	}
	return nil
}

type Card2621012ShadowCloak struct{ AlwaysActive }

func (Card2621012ShadowCloak) ID() string   { return "2621012" }
func (Card2621012ShadowCloak) Name() string { return "暗影披风" }
func (Card2621012ShadowCloak) ModifyEnemySpellStats(ctx *EffectContext, stats *SpellStats) {
	if isEnemySpellCast(ctx) && ctx.ExtraData["stat"] == "damage" && ctx.Source.Statuses["已防护"] == 0 {
		stats.DamageBonus -= 99
		ctx.Source.Statuses["已防护"] = 1
	}
}

type Card2621013WitchcraftRing struct{ AlwaysActive }

func (Card2621013WitchcraftRing) ID() string   { return "2621013" }
func (Card2621013WitchcraftRing) Name() string { return "巫术指环" }
func (Card2621013WitchcraftRing) OnPerTurn(ctx *EffectContext) error {
	for _, skill := range ctx.Engine.State.Players[ctx.OpponentID].Skills {
		if skill != nil && skill.Statuses[StatusWeaken] > 0 {
			skill.Statuses[StatusWeaken]++
		}
	}
	return nil
}

type Card3021006InsightEye struct{ AlwaysActive }

func (Card3021006InsightEye) ID() string   { return "3021006" }
func (Card3021006InsightEye) Name() string { return "洞察之眼" }
func (Card3021006InsightEye) OnSpellHit(ctx *EffectContext) error {
	for _, equipment := range ctx.Engine.State.Players[ctx.OpponentID].Equipment {
		if equipment != nil {
			ctx.Engine.destroyEnemyEquipment(ctx.PlayerID, equipment.InstanceID)
			break
		}
	}
	return nil
}

type Card3021010Dispel struct{ AlwaysActive }

func (Card3021010Dispel) ID() string   { return "3021010" }
func (Card3021010Dispel) Name() string { return "解咒" }
func (Card3021010Dispel) OnSpellHit(ctx *EffectContext) error {
	if ctx.Target != nil {
		ctx.Target.Statuses[StatusSeal]++
	}
	return nil
}

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
	if spell.TotalPower > 0 {
		spell.TotalPower--
		ctx.Engine.emit(GameEvent{
			Type:   "spell_reaction",
			Player: -1,
			Data: map[string]any{
				"player": ctx.PlayerID,
				"card":   cardToInfo(ctx.Source),
				"effect": "power_minus_one",
				"power":  spell.TotalPower,
			},
		})
	}
	return nil
}

type Card3221010WaterPhantom struct{ AlwaysActive }

func (Card3221010WaterPhantom) ID() string   { return "3221010" }
func (Card3221010WaterPhantom) Name() string { return "水幻影" }
func (Card3221010WaterPhantom) OnSpellHit(ctx *EffectContext) error {
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{Type: "next_water_copy", RemainingUses: 1, ExpiresTurn: ctx.Engine.State.TurnNumber + 1})
	return nil
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
	target := ctx.Target
	if target == nil || target.OwnerID != ctx.PlayerID || !target.Card.IsCompanion() {
		target = firstFriendlyCompanion(ctx)
	}
	if target != nil {
		target.Statuses["防止致命"] = 2
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
	affected := ctx.Engine.spellAffectedUnits(defenderID, spell.Skill, spell.Target)
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
func (Card4411002Andrew) OnTurnStart(ctx *EffectContext) error {
	if ctx.Source.Statuses["开局入组"] == 0 {
		addCardToDeck(ctx, "1401002", 1)
		ctx.Source.Statuses["开局入组"] = 1
	}
	return nil
}

type Card4411003ProfessorMaggie struct{ AlwaysActive }

func (Card4411003ProfessorMaggie) ID() string   { return "4411003" }
func (Card4411003ProfessorMaggie) Name() string { return "麦吉教授" }
func (Card4411003ProfessorMaggie) ModifyCardPlayCost(ctx *EffectContext, card *CardInstance, cost map[string]int) {
	if ctx.Source.Statuses["麦吉折扣"] == 0 && card != nil && card.Card.TotalCost() > 5 {
		reduceCost(cost, model.ElementEarth, 2)
		ctx.Source.Statuses["麦吉折扣"] = 1
	}
}

type Card4511002Ailimer struct{ AlwaysActive }

func (Card4511002Ailimer) ID() string   { return "4511002" }
func (Card4511002Ailimer) Name() string { return "神之眷子 爱里默" }
func (Card4511002Ailimer) OnTurnStart(ctx *EffectContext) error {
	if ctx.Source.Statuses["桎梏已洗入"] == 0 {
		addCardToOpponentDeck(ctx, "2501001", 5)
		ctx.Source.Statuses["桎梏已洗入"] = 1
	}
	return nil
}
func (Card4511002Ailimer) OnUltimate(ctx *EffectContext) error {
	if target := firstEnemyUnit(ctx); target != nil && !target.Card.IsHero() {
		ctx.Engine.destroyUnit(target, ctx.OpponentID)
	}
	return nil
}

type Card4511003Lexia struct{ AlwaysActive }

func (Card4511003Lexia) ID() string   { return "4511003" }
func (Card4511003Lexia) Name() string { return "骑士团长 蕾曦娅" }
func (Card4511003Lexia) OnTurnStart(ctx *EffectContext) error {
	if ctx.Source.Statuses["团结希望"] == 0 {
		addSkillToPool(ctx, "3501001")
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
