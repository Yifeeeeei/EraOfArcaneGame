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
		"sendAction(req.action, { ...req.data, payment: { ...paymentSelection.value } })",
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("pending action costs should route through payment request, missing %q", want)
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
