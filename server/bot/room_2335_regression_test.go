package bot

import (
	"encoding/json"
	"eraofarcane/cards"
	"eraofarcane/game"
	"testing"
)

func TestThunderlightBuildsLoadBeforeMoreSkills(t *testing.T) {
	e, p := thunderScenario(t)
	player := e.State.Players[0]
	apprentice := handCard(t, e, "1021001")
	player.Elements["气"], player.Elements["光"] = 6, 2
	for _, n := range []string{"3321002", "3321003", "3321103", "3021101"} {
		player.SkillPool = append(player.SkillPool, game.NewCardInstance(cards.PlayableCardDB[n], 0, 1))
	}
	d := applyNext(t, e, p)
	if d.Action.Action != "summon" || text(d.Action.Data["instance_id"]) != apprentice.InstanceID {
		t.Fatalf("learning displaced affordable development: %+v", d)
	}
}

func TestThunderlightLearnsCheapAirBeforeLargeSpell(t *testing.T) {
	e, p := thunderScenario(t)
	player := e.State.Players[0]
	put(t, e, "1321011", 0, 0, 0)
	put(t, e, "1321011", 0, 2, 0)
	player.Elements["气"], player.Elements["光"] = 6, 2
	player.Hero.IsHorizontal = true
	var lightning *game.CardInstance
	for _, n := range []string{"3321103", "3021101", "3321004", "3321002"} {
		c := game.NewCardInstance(cards.PlayableCardDB[n], 0, 1)
		player.SkillPool = append(player.SkillPool, c)
		if n == "3321002" {
			lightning = c
		}
	}
	d := applyNext(t, e, p)
	if d.Action.Action != "learn_skill" || text(d.Action.Data["instance_id"]) != lightning.InstanceID {
		t.Fatalf("large spell displaced efficient first air spell: %+v", d)
	}
}

func TestThunderlightDevelopsAfterDefenseBeforeFourthLargeSkill(t *testing.T) {
	e, p := thunderScenario(t)
	player := e.State.Players[0]
	player.Hero.IsHorizontal = true
	for _, col := range []int{0, 2} {
		put(t, e, "1321011", 0, col, 0).IsHorizontal = true
	}
	for i, n := range []string{"3321002", "3321003", "3321004"} {
		player.Skills[i] = game.NewCardInstance(cards.PlayableCardDB[n], 0, 1)
	}
	player.Elements["气"] = 5
	apprentice := handCard(t, e, "1021001")
	player.SkillPool = []*game.CardInstance{game.NewCardInstance(cards.PlayableCardDB["3321106"], 0, 1)}
	d := applyNext(t, e, p)
	if d.Action.Action != "summon" || text(d.Action.Data["instance_id"]) != apprentice.InstanceID {
		t.Fatalf("bought a fourth large skill instead of developing after defense: %+v", d)
	}
}

func TestThunderlightArcaneImpactNeverSpendsAirOrReservedLight(t *testing.T) {
	for _, learning := range []bool{false, true} {
		for _, light := range []int{0, 1, 2, 3, 4} {
			e, p := thunderScenario(t)
			player := e.State.Players[0]
			player.Hero.IsHorizontal = true
			player.Elements["气"], player.Elements["光"] = 8, light
			put(t, e, "1321011", 0, 0, 0).IsHorizontal = true
			put(t, e, "1321011", 0, 2, 0).IsHorizontal = true
			put(t, e, "1321002", 0, 0, 2).IsHorizontal = true
			for i, n := range []string{"3321002", "3321003"} {
				player.Skills[i] = game.NewCardInstance(cards.PlayableCardDB[n], 0, 1)
				player.Skills[i].IsHorizontal = true
			}
			handCard(t, e, "1321006") // Reserves one light for development.
			var impact *game.CardInstance
			cost := 2
			if learning {
				impact = game.NewCardInstance(cards.PlayableCardDB["3021101"], 0, 1)
				player.SkillPool = []*game.CardInstance{impact}
				cost = 3
			} else {
				impact = put(t, e, "3021101", 0, 0, 0)
			}
			obs, err := e.ObserveDecision(0)
			if err != nil {
				t.Fatal(err)
			}
			var s state
			if err := json.Unmarshal(obs.State, &s); err != nil {
				t.Fatal(err)
			}
			move := action("cast_spell", map[string]any{"instance_id": impact.InstanceID})
			if learning {
				move.Action = "learn_skill"
			}
			funded, ok := p.thunderFund(s, obs.Cards, move, map[string]int{"无": cost})
			if light < cost+1 {
				if ok {
					t.Fatalf("learning=%v light=%d improperly funded impact: %+v", learning, light, funded)
				}
			} else if !ok || elements(funded.Data["payment"])["光"] != cost {
				t.Fatalf("spare light not used: %+v", funded)
			}
		}
	}
}

func TestFireDefenseIncludesElementalEnchantBurn(t *testing.T) {
	e := scenario(t)
	hero := e.State.Players[0].Hero
	hero.CurrentLife = 2
	put(t, e, "3121001", 0, 0, 0)
	enchant := game.NewCardInstance(cards.PlayableCardDB["3021007"], 1, 1)
	if err := (game.Card3021007ElementalEnchant{}).OnSpellCast(&game.EffectContext{Engine: e, Source: enchant, PlayerID: 1, OpponentID: 0}); err != nil {
		t.Fatal(err)
	}
	if err := e.HandleAction(1, action("resolve_action", map[string]any{"selected": []string{game.StatusBurn}})); err != nil {
		t.Fatal(err)
	}
	incomingHeroSpell(t, e, "3121001")
	before, _ := json.Marshal(e.State.Players[1].TempModifiers)
	for i := 0; i < 2; i++ {
		obs, err := e.ObserveDecision(0)
		if err != nil {
			t.Fatal(err)
		}
		if obs.IncomingBurn != 1 || obs.Cards[hero.InstanceID].Incoming.Burn != 1 {
			t.Fatalf("enchantment omitted: %+v", obs.Cards[hero.InstanceID])
		}
	}
	after, _ := json.Marshal(e.State.Players[1].TempModifiers)
	if string(before) != string(after) {
		t.Fatal("preview spent enchantment")
	}
	d := next(t, e, New(), 0)
	if d.Action.Action != "defend" {
		t.Fatalf("ignored lethal enchantment: %+v", d)
	}
	if err := e.HandleAction(0, d.Action); err != nil {
		t.Fatal(err)
	}
	if hero.CurrentLife != 2 || hero.Statuses[game.StatusBurn] != 0 {
		t.Fatal("defense did not prevent lethal hit and burn")
	}
}
