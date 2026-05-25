package game

import (
	"strings"

	"eraofarcane/cards"
	"eraofarcane/model"
)

type SpellArea string

const (
	SpellAreaSingle      SpellArea = "single"
	SpellAreaSquare      SpellArea = "square"
	SpellAreaAll         SpellArea = "all"
	SpellAreaColumn      SpellArea = "column"
	SpellAreaFrontRow    SpellArea = "front_row"
	SpellAreaSplashCross SpellArea = "splash_cross"
)

type cardTraits struct {
	rush         bool
	pierce       bool
	temporary    bool
	taunt        bool
	stealth      int
	shield       int
	shielding    bool
	cooldown     int
	perTurnLimit int
	overload     int

	needsTarget *bool
	area        SpellArea
	defenseOnly bool
	sorcery     bool
	noBoost     bool
	noAttack    bool
	noDefend    bool

	statuses     map[string]int
	elementGains map[string]int
}

func truePtr() *bool {
	v := true
	return &v
}

func falsePtr() *bool {
	v := false
	return &v
}

func traitsForCardNumber(number string) cardTraits {
	t := cardTraits{area: SpellAreaSingle}
	if card, ok := cards.GetPlayableCard(number); ok && card.IsSkill() && card.Power == -1 {
		t.sorcery = true
	}

	switch number {
	case "1021011", "3021001", "3021008", "3021009", "3121007", "3121011", "3321012", "3321013", "3421015", "3521011", "3621009":
		t.rush = true
	}
	switch number {
	case "1011001", "1011002", "1111001", "1111003", "1211003", "1321010", "1511003":
		t.taunt = true
	}
	switch number {
	case "2121011", "2521008", "3021011", "3121006", "3121010", "3121015", "3201002", "3221011", "3321004", "3321009", "3321013", "3421009", "3421010", "3421012", "3521005", "3621004", "3621011":
		t.pierce = true
	}
	switch number {
	case "3021004", "3021006", "3021008", "3021010", "3021012", "3121007", "3121008", "3221006", "3221007", "3221008", "3221010", "3221015", "3321007", "3321008", "3321012", "3321014", "3421009", "3421013", "3421014", "3521014", "3621010":
		t.cooldown = 1
	case "3421015", "3521011", "3621015":
		t.cooldown = 2
	}
	switch number {
	case "1211003", "2321001":
		t.perTurnLimit = 3
	}

	switch number {
	case "2121008", "2321003", "2421008", "3121005", "3121010", "3221001", "3321006", "3421007", "3421014", "3621011":
		t.area = SpellAreaSquare
	case "3221006", "3421013":
		t.area = SpellAreaAll
	case "3321011", "3521012":
		t.area = SpellAreaColumn
	case "2521009", "3001001", "3121004", "3521008", "3621003":
		t.area = SpellAreaFrontRow
	case "2221009", "2621009", "3201002", "3221005":
		t.area = SpellAreaSplashCross
	}

	switch number {
	case "2121009", "2521013", "3121012", "3121013", "3201001", "3221004", "3221014", "3321010", "3321015", "3421001", "3421005", "3521003", "3621013", "3621014":
		t.defenseOnly = true
	}
	switch number {
	case "3001001":
		t.noBoost = true
	}
	switch number {
	case "3221008":
		t.noAttack = true
	}

	switch number {
	case "3021001", "3021004", "3021006", "3021007", "3021010", "3021012", "3221007", "3221010", "3321007", "3321014", "3621012":
		t.needsTarget = falsePtr()
	case "3021005", "3021008", "3021009", "3121003", "3121005", "3121006", "3121010", "3121011", "3221001", "3221005", "3221006", "3221011", "3221012", "3321003", "3321004", "3321006", "3321009", "3321011", "3321013", "3421002", "3421007", "3421009", "3421010", "3421012", "3421013", "3421014", "3521004", "3521005", "3521008", "3521012", "3621003", "3621004", "3621011":
		t.needsTarget = truePtr()
	}

	switch number {
	case "3021009", "3321003", "3321006", "3421007", "3521004", "3521008":
		t.statuses = map[string]int{StatusStun: 1}
	case "2221003", "2221009", "3221005", "3201002", "3221014":
		t.statuses = map[string]int{StatusFreeze: 1}
	case "3221012":
		t.statuses = map[string]int{StatusFreeze: 2}
	case "3121003":
		t.statuses = map[string]int{StatusBurn: 2}
	case "3121005", "3121010", "3121011", "3121013":
		t.statuses = map[string]int{StatusBurn: 1}
	case "3421002":
		t.statuses = map[string]int{StatusPetrify: 1}
	case "2421005", "3421009":
		t.statuses = map[string]int{StatusPetrify: 2}
	case "2621001", "3621009", "3621014":
		t.statuses = map[string]int{StatusWeaken: 2}
	}

	return t
}

func hasPerTurnAbilityNumber(number string) bool {
	switch number {
	case "1211001", "1211003", "1221005", "1221014", "1221015",
		"1321001", "1321013", "1321015", "1421009", "1421010",
		"1421012", "1521001", "1621009", "2111001", "2111002", "2121001",
		"2311002", "2411001", "2421011", "2511002", "2601001", "2621013",
		"4111002":
		return true
	default:
		return false
	}
}

func hasUltimateAbilityNumber(number string) bool {
	switch number {
	case "1021012", "1121010", "1221011", "1321005", "1411001", "1511001", "1521011",
		"1611002", "1621004", "1621012", "2011003", "2021006", "2121007", "2211001", "2321012",
		"2521007", "4311001", "4311003", "4511002", "4611002":
		return true
	default:
		return false
	}
}

func devourRequirementForNumber(number string) map[string]int {
	switch number {
	case "1111001":
		return map[string]int{model.ElementFire: 3}
	case "1321010":
		return map[string]int{model.ElementAir: 3}
	case "1621010":
		return map[string]int{model.ElementShadow: 4}
	default:
		return nil
	}
}

func behaviorForNumber(number string) CardBehavior {
	return globalRegistry.GetBehavior(number)
}

func cardBehavior(card *CardInstance) CardBehavior {
	if card == nil || card.Card == nil {
		return nil
	}
	return behaviorForNumber(card.Card.Number)
}

func cardHasActiveDeathrattle(card *CardInstance) bool {
	if h, ok := cardBehavior(card).(OnDeathBehavior); ok {
		if h.HasActiveDeathrattle(card) {
			return true
		}
	}
	return len(attachedDeathrattles(card)) > 0
}

func cardHasActivePerTurn(card *CardInstance) bool {
	if h, ok := cardBehavior(card).(PerTurnAbility); ok {
		return h.HasActivePerTurn(card)
	}
	return false
}

func cardHasActiveUltimate(card *CardInstance) bool {
	if h, ok := cardBehavior(card).(UltimateAbility); ok {
		return h.HasActiveUltimate(card)
	}
	return false
}

func cardHasActiveSpellReaction(card *CardInstance) bool {
	if h, ok := cardBehavior(card).(SpellReactionBehavior); ok {
		return h.HasActiveSpellReaction(card)
	}
	return false
}

func cardHasRush(card *CardInstance) bool {
	if card == nil || card.Card == nil {
		return false
	}
	if h, ok := behaviorForNumber(card.Card.Number).(RushBehavior); ok && h.HasActiveRush(card) {
		return h.HasRush()
	}
	return traitsForCardNumber(card.Card.Number).rush
}

func cardHasPierce(card *CardInstance) bool {
	if card == nil || card.Card == nil {
		return false
	}
	if h, ok := behaviorForNumber(card.Card.Number).(PierceBehavior); ok && h.HasActivePierce(card) {
		return h.HasPierce()
	}
	return traitsForCardNumber(card.Card.Number).pierce
}

func cardIsTemporary(card *CardInstance) bool {
	if card == nil || card.Card == nil {
		return false
	}
	if h, ok := behaviorForNumber(card.Card.Number).(TemporaryBehavior); ok && h.HasActiveTemporary(card) {
		return h.IsTemporary()
	}
	return traitsForCardNumber(card.Card.Number).temporary
}

func cardHasTaunt(card *CardInstance) bool {
	if card == nil || card.Card == nil {
		return false
	}
	if h, ok := behaviorForNumber(card.Card.Number).(TauntBehavior); ok && h.HasActiveTaunt(card) {
		return h.HasTaunt()
	}
	return traitsForCardNumber(card.Card.Number).taunt
}

func cardStealthLayers(card *CardInstance) int {
	if card == nil || card.Card == nil {
		return 0
	}
	if h, ok := behaviorForNumber(card.Card.Number).(StealthBehavior); ok && h.HasActiveStealth(card) {
		return h.StealthLayers()
	}
	return traitsForCardNumber(card.Card.Number).stealth
}

func cardShieldLayers(card *CardInstance) int {
	if card == nil || card.Card == nil {
		return 0
	}
	if h, ok := behaviorForNumber(card.Card.Number).(ShieldBehavior); ok && h.HasActiveShield(card) {
		return h.ShieldLayers()
	}
	return traitsForCardNumber(card.Card.Number).shield
}

func cardHasShielding(card *CardInstance) bool {
	if card == nil || card.Card == nil {
		return false
	}
	if h, ok := behaviorForNumber(card.Card.Number).(ShieldingBehavior); ok && h.HasActiveShielding(card) {
		return h.HasShielding()
	}
	return traitsForCardNumber(card.Card.Number).shielding
}

func skillCooldown(skill *CardInstance) int {
	if skill == nil || skill.Card == nil {
		return 0
	}
	if h, ok := behaviorForNumber(skill.Card.Number).(CooldownBehavior); ok && h.HasActiveCooldown(skill) {
		return h.Cooldown()
	}
	return traitsForCardNumber(skill.Card.Number).cooldown
}

func perTurnLimit(card *CardInstance) int {
	if card == nil || card.Card == nil {
		return 1
	}
	if h, ok := behaviorForNumber(card.Card.Number).(PerTurnLimitBehavior); ok && h.HasActivePerTurnLimit(card) {
		if n := h.PerTurnLimit(); n > 0 {
			return n
		}
	}
	if n := traitsForCardNumber(card.Card.Number).perTurnLimit; n > 0 {
		return n
	}
	return 1
}

func spellArea(skill *CardInstance) SpellArea {
	if skill == nil || skill.Card == nil {
		return SpellAreaSingle
	}
	if h, ok := behaviorForNumber(skill.Card.Number).(SpellAreaBehavior); ok && h.HasActiveSpellArea(skill) {
		if area := h.SpellArea(); area != "" {
			return area
		}
	}
	if skill.Statuses["下一次范围前排"] > 0 {
		return SpellAreaFrontRow
	}
	return traitsForCardNumber(skill.Card.Number).area
}

func cardRevealsOnDraw(card *CardInstance) bool {
	if card == nil || card.Card == nil {
		return false
	}
	reveal, ok := behaviorForNumber(card.Card.Number).(DrawRevealBehavior)
	return ok && reveal.HasActiveDrawReveal(card) && reveal.RevealsOnDraw()
}

func skillNeedsTargetInstance(skill *CardInstance) bool {
	if skill == nil || skill.Card == nil || !skill.Card.IsSkill() {
		return false
	}
	if h, ok := behaviorForNumber(skill.Card.Number).(SpellTargetingBehavior); ok && h.HasActiveSpellTargeting(skill) {
		return h.NeedsSpellTarget()
	}
	t := traitsForCardNumber(skill.Card.Number)
	if t.needsTarget != nil {
		return *t.needsTarget
	}
	return max(skill.Card.Attack, 0) > 0 || max(skill.Card.Power, 0) > 0
}

func skillNeedsTargetCard(card *model.Card) bool {
	if card == nil || !card.IsSkill() {
		return false
	}
	return skillNeedsTargetInstance(&CardInstance{Card: card})
}

func canUseSkillForPurpose(card *model.Card, purpose skillPurpose) bool {
	if card == nil || !card.IsSkill() {
		return false
	}
	instance := &CardInstance{Card: card}
	if h, ok := behaviorForNumber(card.Number).(SkillUsabilityBehavior); ok && h.HasActiveSkillUsability(instance) {
		return h.CanUseForSkillPurpose(purpose)
	}
	return staticCanUseSkillForPurpose(card, traitsForCardNumber(card.Number), purpose)
}

func staticCanUseSkillForPurpose(card *model.Card, t cardTraits, purpose skillPurpose) bool {
	if t.sorcery {
		return purpose == skillPurposeAttack && !t.noAttack
	}
	switch purpose {
	case skillPurposeAttack:
		return !t.defenseOnly && !t.noAttack
	case skillPurposeDefend:
		return !t.noDefend && card.Power > 0
	case skillPurposeBoost, skillPurposeAttackBoost:
		return !t.defenseOnly && !t.noBoost && card.Power > 0
	case skillPurposeDefenseBoost:
		return !t.noBoost && card.Power > 0
	default:
		return false
	}
}

func isDefenseOnlySkill(card *model.Card) bool {
	if card == nil || !card.IsSkill() {
		return false
	}
	instance := &CardInstance{Card: card}
	if h, ok := behaviorForNumber(card.Number).(DefenseOnlySkillBehavior); ok && h.HasActiveDefenseOnlySkill(instance) {
		return h.IsDefenseOnlySkill()
	}
	return traitsForCardNumber(card.Number).defenseOnly
}

func isSorcerySkill(card *model.Card) bool {
	if card == nil || !card.IsSkill() {
		return false
	}
	instance := &CardInstance{Card: card}
	if h, ok := behaviorForNumber(card.Number).(SorcerySkillBehavior); ok && h.HasActiveSorcerySkill(instance) {
		return h.IsSorcerySkill()
	}
	return traitsForCardNumber(card.Number).sorcery
}

func isConsumableCardInstance(card *CardInstance) bool {
	return card != nil && card.Card != nil && cards.IsConsumable(card.Card.Number)
}

func isEquipmentItem(card *CardInstance) bool {
	return card != nil && card.Card != nil && cards.IsEquipment(card.Card.Number)
}

func isRuneOrScroll(card *CardInstance) bool {
	if card == nil || card.Card == nil {
		return false
	}
	number := card.Card.Number
	return cards.IsConsumable(number) && !cards.IsEquipment(number)
}

func ruleInfoForCard(card *model.Card) map[string]any {
	if card == nil {
		return map[string]any{}
	}
	t := traitsForCardNumber(card.Number)
	info := map[string]any{
		"is_terrain":    cards.IsTerrain(card.Number),
		"is_consumable": cards.IsConsumable(card.Number),
		"is_equipment":  cards.IsEquipment(card.Number),
		"is_weapon":     cards.IsWeapon(card.Number),
		"has_per_turn":  hasPerTurnAbilityNumber(card.Number),
		"has_ultimate":  hasUltimateAbilityNumber(card.Number),
	}
	if requirement := devourRequirementForNumber(card.Number); len(requirement) > 0 {
		info["devour_requirement"] = requirement
	}
	if card.IsSkill() {
		needsTarget := max(card.Attack, 0) > 0 || max(card.Power, 0) > 0
		if t.needsTarget != nil {
			needsTarget = *t.needsTarget
		}
		info["is_defense_only"] = t.defenseOnly
		info["is_sorcery"] = t.sorcery
		info["needs_target"] = needsTarget
		info["has_pierce"] = t.pierce
		info["can_attack"] = staticCanUseSkillForPurpose(card, t, skillPurposeAttack)
		info["can_defend"] = staticCanUseSkillForPurpose(card, t, skillPurposeDefend)
		info["can_attack_boost"] = staticCanUseSkillForPurpose(card, t, skillPurposeAttackBoost)
		info["can_defense_boost"] = staticCanUseSkillForPurpose(card, t, skillPurposeDefenseBoost)
		info["can_boost"] = info["can_attack_boost"]
		info["spell_area"] = t.area
		info["can_react"] = hasSpellReactionNumber(card.Number)
	}
	return info
}

func CardRuleInfo(card *model.Card) map[string]any {
	return ruleInfoForCard(card)
}

func isBeastPlantOrSpirit(card *CardInstance) bool {
	if card == nil || card.Card == nil {
		return false
	}
	if hasAnyTag(card.Card.Tag, "野兽", "植物", "精灵") {
		return true
	}
	switch card.Card.Number {
	case "1021007", "1021008", "1121001", "1121004", "1121006", "1121009", "1121014",
		"1221001", "1221006", "1221007", "1221014", "1311001", "1321001", "1321002",
		"1321003", "1321008", "1321011", "1321015", "1401001", "1401002", "1411002",
		"1421003", "1421006", "1421007", "1421008", "1421009", "1421010", "1421011",
		"1421012", "1421013", "1501001", "1521002", "1521005", "1521006", "1521007",
		"1521008", "1621001":
		return true
	default:
		return false
	}
}

func hasAnyTag(tag string, needles ...string) bool {
	for _, needle := range needles {
		if strings.Contains(tag, needle) {
			return true
		}
	}
	return false
}

func isConstructOrDemon(card *CardInstance) bool {
	if card == nil || card.Card == nil {
		return false
	}
	switch card.Card.Number {
	case "1011002", "1021009", "1121002", "1121005", "1121007", "1121008", "1221004",
		"1221008", "1321004", "1321009", "1421004", "1421005", "1521004", "1611001",
		"1621002", "1621003", "1621005", "1621006", "1621007", "1621010", "1621011",
		"1621013", "1621014":
		return true
	default:
		return false
	}
}
