package model_test

import (
	"eraofarcane/cards"
	"eraofarcane/model"
	"sort"
	"testing"
)

func validConstruction(t *testing.T) *model.Deck {
	t.Helper()
	if err := cards.LoadCards(); err != nil {
		t.Fatal(err)
	}
	keys := []string{}
	for key := range cards.PlayableCardDB {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	d := &model.Deck{HeroID: "4011001", ExtraDeck: []string{"1001101"}}
	for _, key := range keys {
		c := cards.PlayableCardDB[key]
		if c.HasTag("传奇") || c.HasTag("衍生") {
			continue
		}
		if (c.IsCompanion() || c.IsItem()) && len(d.MainDeck) < 30 {
			d.MainDeck = append(d.MainDeck, key, key)
		}
		if c.IsSkill() && len(d.SkillPool) < 10 {
			d.SkillPool = append(d.SkillPool, key)
		}
	}
	return d
}

func TestConstructionRestrictions166(t *testing.T) {
	for _, scenario := range []string{"valid", "legendary_pair", "legendary_skills", "derived_main", "derived_skill"} {
		t.Run(scenario, func(t *testing.T) {
			d := validConstruction(t)
			switch scenario {
			case "legendary_pair":
				d.MainDeck[0], d.MainDeck[1] = "1011001", "1011001"
			case "legendary_skills":
				d.SkillPool[0], d.SkillPool[1] = "3011101", "3111101"
			case "derived_main":
				d.MainDeck[0] = "1001101"
			case "derived_skill":
				d.SkillPool[0] = "3001101"
			}
			err := d.Validate(cards.PlayableCardDB)
			if (err == nil) != (scenario == "valid") {
				t.Fatalf("scenario %s: %v", scenario, err)
			}
		})
	}
}
