package game

import (
	"os"
	"strings"
	"testing"
)

func TestGameHTMLFriendlyUnitImageCanSelectFriendlySpellTarget(t *testing.T) {
	content, err := os.ReadFile("../../web/game.html")
	if err != nil {
		t.Fatalf("read game.html: %v", err)
	}
	html := string(content)
	want := `@click.stop="handleUnitImageClick(mySlot, col, row, getMyUnit(col, row))"`
	if !strings.Contains(html, want) {
		t.Fatalf("friendly unit image should route through handleUnitImageClick so friendly spell targets work")
	}
	bad := `@click.stop="showCardDetail(getMyUnit(col, row))"`
	if strings.Contains(html, bad) {
		t.Fatalf("friendly unit image should not stop clicks by directly opening detail")
	}
}

func TestGameHTMLPendingActionCostUsesPaymentRequest(t *testing.T) {
	content, err := os.ReadFile("../../web/game.html")
	if err != nil {
		t.Fatalf("read game.html: %v", err)
	}
	html := string(content)
	for _, want := range []string{
		"pendingRequiresPaymentChoice()",
		"paymentRequest.value = {",
		"action: 'resolve_action'",
		"afterSend: resetPendingSelectionState",
		"sendAction(req.action, { ...req.data, payment: { ...paymentSelection.value } }, {",
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("pending action costs should route through payment request, missing %q", want)
		}
	}
}

func TestGameHTMLRainbowAngelPaymentChoice(t *testing.T) {
	content, err := os.ReadFile("../../web/game.html")
	if err != nil {
		t.Fatalf("read game.html: %v", err)
	}
	html := string(content)
	for _, want := range []string{
		"function hasArcanePaymentChoice(cost)",
		"function hasLightWildcardPaymentChoice(cost)",
		"calculatePaymentWithAvailable(myState.value.elements || {}, cost, false) === null",
		"remainingCost[elem] = Number(remainingCost[elem] || 0) - 1",
		"remainingAvailable['光'] = Number(remainingAvailable['光'] || 0) - 1",
		"return hasArcanePaymentChoice(cost) || hasLightWildcardPaymentChoice(cost);",
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("Rainbow Angel should create explicit payment choices when light can substitute ambiguously, missing %q", want)
		}
	}
}

func TestGameHTMLDefenseWindowIncludesBoundSkills(t *testing.T) {
	content, err := os.ReadFile("../../web/game.html")
	if err != nil {
		t.Fatalf("read game.html: %v", err)
	}
	html := string(content)
	for _, want := range []string{
		"return allMySkills.value.filter(canUseSkillAsDefense);",
		"...allMySkills.value,",
		"...allMySkills.value.filter(s => s && defenseSelected.value.includes(s.instance_id)),",
		"const learnedBoosts = allMySkills.value.filter(s =>",
		"const scrollBoosts = myHand.value.filter(s =>",
		"canUseScrollAsDefenseBoost(s)",
		"function canUseScrollAsDefenseBoost(scroll)",
		"const availableDefenseBoostCandidates = computed(() =>",
		"function isDefenseBoostDisabled(skill)",
		"const skills = allMySkills.value;",
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("defense window should include bound skills, missing %q", want)
		}
	}
	if strings.Contains(html, "for (const equipment of myEquipment.value.filter(Boolean))") {
		t.Fatalf("bound skills should not be collected from equipment twice")
	}
}

func TestGameHTMLDefenseWindowClickHandlersAreReturned(t *testing.T) {
	content, err := os.ReadFile("../../web/game.html")
	if err != nil {
		t.Fatalf("read game.html: %v", err)
	}
	html := string(content)
	for _, want := range []string{
		`@click.capture="handleDefensePanelClick"`,
		`data-defense-action="scroll"`,
		`data-defense-action="skill"`,
		`data-defense-action="boost"`,
		`data-defense-action="overexert"`,
		`data-defense-action="submit"`,
		`data-defense-action="no_defend"`,
		`data-defense-action="detail"`,
		`:class="{ disabled: !canSubmitDefense }"`,
		"function defenseCardById(id)",
		"function handleDefensePanelClick(event)",
		"function syncDefensePanelVisualState()",
		"el.classList.toggle('selected', selected)",
		"submit.classList.toggle('disabled', !canSubmitDefense.value)",
		"const CLIENT_BUILD = 'issue133-royal-feedback-20260810-1'",
		"const pendingSpellView = computed(() => normalizePendingSpell(gameState.value?.pending_spell))",
		"function normalizePendingSpell(spell)",
		"spell.power ?? spell.total_power",
		"spell.defender_id === currentPlayerSlot()",
		"data-defense-summary=\"power\"",
		"function syncDefenseSummary(panel)",
		"defensePower, defenseBoostPower, defensePaymentCost,",
		"function logDefenseWindowStateSync(state)",
		"function logDefenseSelection(message, card = null)",
		"defense_submit_blocked",
		"if (!defenseSelected.value.includes(skill.instance_id))",
		"if (!defenseScrollSelected.value.includes(scroll.instance_id))",
		"handleDefensePanelClick, showInteractionCardContextMenu, toggleDefenseSkill, toggleDefenseScroll, toggleDefenseBoost, toggleDefenseOverexert, toggleDefensePowerSacrifice, reactSpell, submitDefend",
		"canSubmitDefense,",
		"noDefend,",
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("defense window click handlers should be exposed to the template, missing %q", want)
		}
	}
	if strings.Contains(html, `:disabled="!canSubmitDefense"`) {
		t.Fatalf("defense submit should not be silently disabled; submitDefend should log why it cannot submit")
	}
	for _, bad := range []string{
		`@pointerdown.capture="handleDefensePanelPointerDown"`,
		`@pointerdown.stop.prevent="toggleDefenseScroll(scroll)"`,
		`@pointerdown.stop.prevent="toggleDefenseSkill(skill)"`,
		`@pointerdown.stop.prevent="toggleDefenseBoost(skill)"`,
		`@pointerdown.stop.prevent="toggleDefenseOverexert(unit)"`,
		`@pointerdown.stop.prevent="submitDefend()"`,
		`@pointerdown.stop.prevent="noDefend()"`,
	} {
		if strings.Contains(html, bad) {
			t.Fatalf("defense card actions should be handled by panel-level delegation, found %q", bad)
		}
	}
}

func TestGameHTMLInteractionWindowsUseUnifiedReadableCards(t *testing.T) {
	content, err := os.ReadFile("../../web/game.html")
	if err != nil {
		t.Fatalf("read game.html: %v", err)
	}
	html := string(content)
	for _, want := range []string{
		"css/game.css?v=20260703-interaction-card-meta",
		"class=\"overlay interaction-overlay defense-overlay\"",
		"class=\"overlay-content interaction-panel defense-panel\"",
		"class=\"overlay interaction-overlay pending-action-overlay\"",
		"class=\"overlay-content interaction-panel pending-action-panel\"",
		"class=\"defense-skill-card interaction-card",
		"class=\"pending-card interaction-card",
		"class=\"interaction-detail-link\"",
		"class=\"overlay card-detail-overlay\"",
		"@mouseenter=\"candidate.number && previewCard(candidate)\"",
		"@contextmenu.prevent=\"showInteractionCardContextMenu($event)\"",
		"@contextmenu.prevent.stop=\"candidate.number && showCardContextMenu($event, candidate, 'interaction')\"",
		"showInteractionCardContextMenu,",
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("interaction windows should share readable interaction UI, missing %q", want)
		}
	}
}

func TestGameCSSInteractionCardsKeepMetadataReadable(t *testing.T) {
	content, err := os.ReadFile("../../web/css/game.css")
	if err != nil {
		t.Fatalf("read game.css: %v", err)
	}
	css := string(content)
	for _, want := range []string{
		"-webkit-line-clamp: 2;",
		"grid-template-rows: minmax(32px, auto) auto;",
		".pending-card-side {",
		"z-index: 11950;",
		"z-index: 11940;",
	} {
		if !strings.Contains(css, want) {
			t.Fatalf("interaction cards should preserve metadata/readability, missing %q", want)
		}
	}
}

func TestGameHTMLDefenseWindowDoesNotReadRawPendingSpellPower(t *testing.T) {
	content, err := os.ReadFile("../../web/game.html")
	if err != nil {
		t.Fatalf("read game.html: %v", err)
	}
	html := string(content)
	if strings.Contains(html, "gameState.pending_spell.power") {
		t.Fatalf("defense UI should use normalized pendingSpellView power, not raw pending_spell.power")
	}
	if strings.Contains(html, "gameState.value?.pending_spell?.power") {
		t.Fatalf("defense logging should use pendingAttackPower, not raw pending_spell.power")
	}
}

func TestGameHTMLIssue156SpellRangeAndScrollBoostControls(t *testing.T) {
	content, err := os.ReadFile("../../web/game.html")
	if err != nil {
		t.Fatalf("read game.html: %v", err)
	}
	html := string(content)
	for _, want := range []string{
		"'pending-spell-center': isPendingSpellCenter(opponentBoardOwnerID(), col, row)",
		"'pending-spell-affected': isPendingSpellAffected(myBoardOwnerID(), col, row)",
		"affected_positions: Array.isArray(spell.affected_positions)",
		"function isPendingSpellCenter(ownerID, col, row)",
		"function isPendingSpellAffected(ownerID, col, row)",
		"myBoardOwnerID, opponentBoardOwnerID, isPendingSpellCenter, isPendingSpellAffected,",
		"selectedSkill && canBoostWith(card)",
		"isSpellScrollItem(skill) || skill.statuses?.['石化'] > 0 || !skill.can_attack_boost",
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("issue #156 frontend behavior missing %q", want)
		}
	}

	cssContent, err := os.ReadFile("../../web/css/game.css")
	if err != nil {
		t.Fatalf("read game.css: %v", err)
	}
	css := string(cssContent)
	for _, want := range []string{".unit-cell.pending-spell-affected", ".unit-cell.pending-spell-center", ".hand-card.boost-selected"} {
		if !strings.Contains(css, want) {
			t.Fatalf("issue #156 styling missing %q", want)
		}
	}
}

func TestGameHTMLCompositeCostsUseSequentialPlannerAndActionAcknowledgement(t *testing.T) {
	content, err := os.ReadFile("../../web/game.html")
	if err != nil {
		t.Fatalf("read game.html: %v", err)
	}
	html := string(content)
	for _, want := range []string{
		"/* ACTION_COST_PLAN_START */",
		"function attackActionCost(mainSkill, boostIDs)",
		"function defenseActionCost()",
		"ArcaneActionCostPlan.planActionCost(entries, myState.value.temp_modifiers || [])",
		"missingCostText(attackActionCost(mySkills[i-1], []))",
		"missingCostText(attackActionCost(skill, []))",
		"costLabel(reactionActionCost(skill))",
		"skillUseCost, attackActionCost, reactionActionCost, canLearnSkill",
		"const pendingActionCallbacks = new Map();",
		"message.request_id = requestID;",
		"else if (msg.type === 'action_result')",
		"sendAction('defend', data, {",
		"onSuccess: () => {",
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("composite actions should use sequential pricing and wait for acknowledgement, missing %q", want)
		}
	}
	if strings.Contains(html, "mergedDefenseCost(") {
		t.Fatalf("frontend should not merge independently discounted card costs")
	}
	if strings.Contains(html, "missingCostText(skillUseCost(") {
		t.Fatalf("skill affordability hints should use the same action cost planner as button availability")
	}
}

func TestGameHTMLRoom5543TargetingRegressions(t *testing.T) {
	content, err := os.ReadFile("../../web/game.html")
	if err != nil {
		t.Fatalf("read game.html: %v", err)
	}
	html := string(content)
	for _, want := range []string{
		"target_owner: targetOwner",
		"hasActiveTargeting, activeTargetSpell,",
		"if (area === 'front_row') {",
		"return !!unit && row === opponentFrontRow.value;",
		"function canAttackSpellTargetFriendlyUnit(spell)",
		"selectedItemSpellScroll.value && isFriendlySpellTarget(col, row)",
		"String(card.number) === '2021010' && opponentSkills.value.filter(Boolean).length < 4",
		"需要敌方至少有4个法术",
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("room 5543 frontend targeting fix missing %q", want)
		}
	}
}
