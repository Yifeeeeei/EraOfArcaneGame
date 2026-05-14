package cards

import "testing"

func TestCompiledDefinitionsImplementCategoryInterfaces(t *testing.T) {
	for _, definition := range compiledCardDefinitions {
		switch definition.Kind() {
		case "人物":
			if _, ok := definition.(HeroCard); !ok {
				t.Fatalf("%s %s does not implement HeroCard", definition.ID(), definition.Name())
			}
		case "伙伴":
			if _, ok := definition.(CompanionCard); !ok {
				t.Fatalf("%s %s does not implement CompanionCard", definition.ID(), definition.Name())
			}
		case "技能":
			if _, ok := definition.(SkillCard); !ok {
				t.Fatalf("%s %s does not implement SkillCard", definition.ID(), definition.Name())
			}
		case "道具":
			if _, ok := definition.(ItemCard); !ok {
				t.Fatalf("%s %s does not implement ItemCard", definition.ID(), definition.Name())
			}
		default:
			t.Fatalf("%s %s has unsupported kind %q", definition.ID(), definition.Name(), definition.Kind())
		}
	}
}
