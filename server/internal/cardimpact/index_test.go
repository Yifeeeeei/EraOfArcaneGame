package cardimpact

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestIndexFollowsHelpersAndSeparatesRelatedHooks(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"cards.go": `package game
type Card1111001First struct{}
func(Card1111001First) OnEnter(){ firstHelper() }
type Card1111002Second struct{}
func(Card1111002Second) OnEnter(){}
type Card1111003Third struct{}
func(Card1111003Third) OnDeath(){ firstHelper() }
`,
		"helpers.go": `package game
func firstHelper(){ secondHelper() }
func secondHelper(){ firstHelper() }
`,
		"cards_test.go": `package game
func TestFirst(){ useCard("1111001") }
func TestHelper(){ secondHelper() }
func TestUnrelated(){ useCard("9999999") }
`,
	}
	for name, source := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(source), 0600); err != nil {
			t.Fatal(err)
		}
	}
	index, err := Build(dir)
	if err != nil {
		t.Fatal(err)
	}
	report, err := index.Analyze("1111001", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Dependencies) != 2 {
		t.Fatalf("cycle/transitive dependencies lost: %+v", report.Dependencies)
	}
	if len(report.RelatedCards) != 1 || report.RelatedCards[0].Number != "1111002" {
		t.Fatalf("unexpected shared hooks: %+v", report.RelatedCards)
	}
	if len(report.Tests) != 2 {
		t.Fatalf("literal and helper tests must both be found: %+v", report.Tests)
	}
	if report.Sources[0].Line != 3 {
		t.Fatalf("wrong source line: %+v", report.Sources[0])
	}
	changed, err := index.Analyze("", filepath.Join(dir, "helpers.go"))
	if err != nil {
		t.Fatal(err)
	}
	if len(changed.RelatedCards) != 2 || changed.RelatedCards[0].Number != "1111001" || changed.RelatedCards[1].Number != "1111003" {
		t.Fatalf("reverse helper dependencies lost: %+v", changed.RelatedCards)
	}
	again, err := index.Analyze("1111001", "")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(report, again) {
		t.Fatal("report order is nondeterministic")
	}
}
