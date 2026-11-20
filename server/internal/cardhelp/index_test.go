package cardhelp

import (
	"eraofarcane/cards"
	"eraofarcane/model"
	"slices"
	"testing"
)

func TestPrintedTermsDoNotGrantRestrictedUses(t *testing.T) {
	if err := cards.LoadCards(); err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		number, want, absent string
	}{
		{"3121012", "K13", "A07"}, // Standalone 防御.
		{"3121013", "K13", "A07"}, // Defense plus success effect, one explanation.
		{"3521013", "A07", "K13"}, // A conditional bonus doesn't restrict attacking.
		{"3121101", "A07", "K13"},
		{"3121105", "A07", "K13"}, // Explicit prohibition must not become defense-only.
		{"1121112", "K04", "K01"}, // Entry cost reduction is not an enter trigger.
		{"3121101", "A06", "K24"}, // 法力范围 is not an area spell.
		{"3321015", "K32", ""},    // Printed 晕眩 alias.
	} {
		t.Run(tt.number+tt.want, func(t *testing.T) {
			refs := References(cards.PlayableCardDB[tt.number])
			if !slices.Contains(refs, tt.want) || tt.absent != "" && slices.Contains(refs, tt.absent) {
				t.Fatalf("unexpected printed references: %v", refs)
			}
		})
	}
	if refs := References(cards.PlayableCardDB["3121001"]); len(refs) != 0 {
		t.Fatalf("blank Fireball should not gain terms: %v", refs)
	}
}

func TestOnlyPrintedEffectAndRelevantTagAreScanned(t *testing.T) {
	c := &model.Card{Name: "防御冻结", Quote: "精通", Description: "", Tag: "衍生-伙伴"}
	if got := References(c); !slices.Equal(got, []string{"K21"}) {
		t.Fatalf("name or flavor text leaked into references: %v", got)
	}
	c.Description = "入场,遗言:使1个单位晕眩1.范围:方阵"
	for _, id := range []string{"K01", "K02", "K24", "K27", "K32"} {
		if !slices.Contains(References(c), id) {
			t.Fatalf("missing compound printed reference %s", id)
		}
	}
}
