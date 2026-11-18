package bot

import (
	"encoding/json"
	"eraofarcane/cards"
	"eraofarcane/game"
	"eraofarcane/model"
	"fmt"
	"os"
	"testing"
)

func profileFixture(t *testing.T, a, b string, seed int64) (*game.Engine, [2]*Policy) {
	t.Helper()
	fixture(t, seed)
	var decks [2]*model.Deck
	var policies [2]*Policy
	for i, id := range []string{a, b} {
		profile, err := LookupProfile(id)
		if err != nil {
			t.Fatal(err)
		}
		decks[i], err = model.ParseDeckCode(profile.DeckCode)
		if err != nil {
			t.Fatal(err)
		}
		if err = decks[i].Validate(cards.PlayableCardDB); err != nil {
			t.Fatal(err)
		}
		policies[i], err = NewForProfile(id)
		if err != nil {
			t.Fatal(err)
		}
	}
	e := game.NewEngineWithSeed(fmt.Sprintf("profiles-%d", seed), nil, seed)
	if err := e.SetupGameWithFirstPlayer("A", decks[0], "B", decks[1], 0); err != nil {
		t.Fatal(err)
	}
	return e, policies
}

func TestMixedProfilesFinishWithoutRejectedActions(t *testing.T) {
	for _, pair := range [][2]string{{ThunderlightProfileID, ProfileID}, {ProfileID, ThunderlightProfileID}, {ThunderlightProfileID, ThunderlightProfileID}} {
		for _, seed := range []int64{20261117, 1, 7, 42} {
			t.Run(fmt.Sprintf("%s-vs-%s/%d", pair[0], pair[1], seed), func(t *testing.T) {
				e, policies := profileFixture(t, pair[0], pair[1], seed)
				for step := 0; step < 1500; step++ {
					if e.State.Phase == game.PhaseGameOver {
						return
					}
					acted := false
					for player, p := range policies {
						obs, err := e.ObserveDecision(player)
						if err != nil {
							t.Fatal(err)
						}
						d, ok, err := p.Next(obs)
						if err != nil {
							t.Fatal(err)
						}
						if !ok {
							continue
						}
						if d.Action.Action == "surrender" {
							t.Fatalf("fallback surrender: %s", d.Reason)
						}
						if err = e.HandleAction(player, d.Action); err != nil {
							t.Fatalf("step %d player %d rejected %+v: %v", step, player, d, err)
						}
						p.Accept()
						acted = true
						break
					}
					if !acted {
						t.Fatalf("no player acted at step %d phase %s", step, e.State.Phase)
					}
				}
				t.Fatal("match did not finish")
			})
		}
	}
}

func TestThunderlightPinnedDeckMatchesCatalog(t *testing.T) {
	fixture(t, 1)
	data, err := os.ReadFile("../../web/assets/preconstructed-decks.json")
	if err != nil {
		t.Fatal(err)
	}
	var catalog struct {
		Decks    []struct{ ID, DeckCode string }
		Revision string
	}
	if err = json.Unmarshal(data, &catalog); err != nil {
		t.Fatal(err)
	}
	for _, recipe := range catalog.Decks {
		if recipe.ID == "王权纷争/雷光气" {
			if recipe.DeckCode != ThunderlightDeckCode || catalog.Revision != ThunderlightDeckRevision {
				t.Fatal("pinned AI deck drifted from catalog")
			}
			return
		}
	}
	t.Fatal("thunderlight precon missing")
}

func thunderScenario(t *testing.T) (*game.Engine, *Policy) {
	e, policies := profileFixture(t, ThunderlightProfileID, ProfileID, 17)
	e.State.Phase = game.PhaseMain
	e.State.CurrentTurn = 0
	e.State.TurnNumber = 2
	e.State.IsFirstTurn = false
	for _, player := range e.State.Players {
		player.Hand = nil
		player.SkillPool = nil
		for key := range player.Elements {
			player.Elements[key] = 0
		}
		player.Hero.IsHorizontal = false
	}
	return e, policies[0]
}

func handCard(t *testing.T, e *game.Engine, n string) *game.CardInstance {
	t.Helper()
	c := game.NewCardInstance(cards.PlayableCardDB[n], 0, e.State.TurnNumber)
	e.State.Players[0].Hand = append(e.State.Players[0].Hand, c)
	return c
}

func applyNext(t *testing.T, e *game.Engine, p *Policy) Decision {
	t.Helper()
	d := next(t, e, p, 0)
	if err := e.HandleAction(0, d.Action); err != nil {
		t.Fatalf("rejected %+v: %v", d, err)
	}
	p.Accept()
	return d
}

func TestThunderlightArmorerAndCoralSequencing(t *testing.T) {
	e, p := thunderScenario(t)
	armorer := handCard(t, e, "1021016")
	handCard(t, e, "2021104")
	e.State.Players[0].Deck = []*game.CardInstance{game.NewCardInstance(cards.PlayableCardDB["2321105"], 0, 1), game.NewCardInstance(cards.PlayableCardDB["2021104"], 0, 1)}
	for step := 0; step < 5; step++ {
		d := applyNext(t, e, p)
		if d.Action.Action == "equip" {
			t.Fatal("equipped before available armorer search")
		}
		if d.Action.Action == "summon" {
			if text(d.Action.Data["instance_id"]) != armorer.InstanceID || number(d.Action.Data["row"]) != 0 {
				t.Fatal("armorer did not become a front blocker")
			}
			break
		}
	}
	if e.State.PendingAction == nil {
		t.Fatal("armorer never searched")
	}
	d := applyNext(t, e, p)
	if d.Action.Action != "resolve_action" {
		t.Fatal("skipped equipment search")
	}
	selected := d.Action.Data["selected"].([]any)
	for _, c := range e.State.Players[0].Hand {
		if c.InstanceID == selected[0] && c.Card.Number == "2021104" {
			return
		}
	}
	t.Fatal("did not tutor missing air/light source")
}

func TestThunderlightCoralIsNotHeldForMissingArmorer(t *testing.T) {
	e, p := thunderScenario(t)
	handCard(t, e, "2021104")
	for step := 0; step < 4 && e.State.PendingAction == nil; step++ {
		applyNext(t, e, p)
	}
	if e.State.PendingAction == nil || e.State.PendingAction.Type != "five_color_coral_load" {
		t.Fatal("did not equip available coral")
	}
	applyNext(t, e, p)
	for _, c := range e.State.Players[0].Equipment {
		if c != nil {
			if c.ElementsGainBonus["气"] != 1 || c.ElementsGainBonus["光"] != 1 {
				t.Fatal("coral did not choose air and light")
			}
			return
		}
	}
	t.Fatal("coral missing")
}

func TestThunderlightWarriorBalancesFrontAndLoad(t *testing.T) {
	for _, haveFront := range []bool{false, true} {
		e, p := thunderScenario(t)
		player := e.State.Players[0]
		for i, n := range []string{"2321101", "2321104", "2321105"} {
			player.Equipment[i] = game.NewCardInstance(cards.PlayableCardDB[n], 0, 1)
		}
		if haveFront {
			put(t, e, "1021016", 0, 0, 0)
			put(t, e, "1021016", 0, 2, 0)
		}
		warrior := handCard(t, e, "1321111")
		player.Elements["气"] = 10
		player.Elements["光"] = 10
		applyNext(t, e, p)
		applyNext(t, e, p)
		if haveFront && warrior.CurrentLife != 2 {
			t.Fatal("ignored adequate front line and did not prioritize load")
		}
		if !haveFront && warrior.CurrentLife != 4 {
			t.Fatal("did not add life for missing front line")
		}
	}
}

func TestThunderlightPierceThreatAndLowHeroPriorities(t *testing.T) {
	for _, life := range []int{6, 2} {
		e, p := thunderScenario(t)
		player := e.State.Players[0]
		put(t, e, "3321004", 0, 0, 0)
		put(t, e, "1321006", 0, 0, 2)
		put(t, e, "1021001", 1, 0, 0)
		threat := put(t, e, "1321006", 1, 2, 2)
		e.State.Players[1].Hero.CurrentLife = life
		player.Elements["气"] = 2
		d := next(t, e, p, 0)
		if d.Action.Action != "cast_spell" {
			t.Fatalf("not attacking: %+v", d)
		}
		want := threat.Position
		if life == 2 {
			want = e.State.Players[1].Hero.Position
		}
		if number(d.Action.Data["target_col"]) != want.Col || number(d.Action.Data["target_row"]) != want.Row {
			t.Fatalf("wrong threat target at life=%d: %+v", life, d)
		}
	}
}

func TestThunderlightChimeraConsumesFodderBeforeDevour(t *testing.T) {
	e, p := thunderScenario(t)
	player := e.State.Players[0]
	put(t, e, "1321006", 0, 2, 2)
	put(t, e, "1321011", 0, 0, 2)
	var fodder []*game.CardInstance
	for col := 0; col < 3; col++ {
		fodder = append(fodder, put(t, e, "1321003", 0, col, 0))
	}
	put(t, e, "3321002", 0, 0, 0)
	player.Skills[1] = game.NewCardInstance(cards.PlayableCardDB["3321003"], 0, 1)
	player.Skills[1].IsHorizontal = false
	chimera := handCard(t, e, "1321010")
	consumed := map[string]bool{}
	for step := 0; step < 10; step++ {
		d := applyNext(t, e, p)
		if d.Action.Action == "consume" {
			consumed[text(d.Action.Data["instance_id"])] = true
		}
		if d.Action.Action == "summon" && text(d.Action.Data["instance_id"]) == chimera.InstanceID {
			for _, c := range fodder {
				if !consumed[c.InstanceID] {
					t.Fatal("lost usable fodder load before devour")
				}
			}
			if player.Units[2][2] == nil || player.Units[0][2] == nil {
				t.Fatal("devoured core damage/light engine")
			}
			return
		}
	}
	t.Fatal("never completed legal chimera summon")
}

func TestThunderlightChainFindsOutOfRangeExtraHero(t *testing.T) {
	e, p := thunderScenario(t)
	player := e.State.Players[0]
	chain := game.NewCardInstance(cards.PlayableCardDB["2321101"], 0, 1)
	chain.IsHorizontal = false
	player.Equipment[0] = chain
	put(t, e, "3321002", 0, 0, 0)
	put(t, e, "1021001", 1, 0, 0)
	player.Elements["气"] = 1
	d := applyNext(t, e, p)
	if d.Action.Action != "use_ability" {
		t.Fatalf("did not charge chain: %+v", d)
	}
	d = applyNext(t, e, p)
	if d.Action.Action != "cast_spell" || d.Action.Data["extra_target_col"] == nil {
		t.Fatalf("missing extra target: %+v", d)
	}
	if number(d.Action.Data["extra_target_col"]) != e.State.Players[1].Hero.Position.Col || number(d.Action.Data["extra_target_row"]) != e.State.Players[1].Hero.Position.Row {
		t.Fatal("did not reach extra hero outside normal range")
	}
}

func TestBothPoliciesDefendWhenOffenseIsSuppressed(t *testing.T) {
	for _, id := range []string{ProfileID, ThunderlightProfileID} {
		for _, shelter := range []bool{false, true} {
			e, pair := profileFixture(t, id, ProfileID, 17)
			p := pair[0]
			e.State.Phase = game.PhaseMain
			e.State.CurrentTurn = 0
			e.State.TurnNumber = 2
			e.State.IsFirstTurn = false
			for _, player := range e.State.Players {
				player.Hand = nil
				player.SkillPool = nil
				for k := range player.Elements {
					player.Elements[k] = 0
				}
				player.Hero.IsHorizontal = false
			}
			ownSpell := "3121001"
			incomingSpell := "3121001"
			if !shelter {
				ownSpell = "3121003"
				incomingSpell = "3121003"
			}
			if id == ThunderlightProfileID {
				ownSpell = "3321002"
			}
			put(t, e, ownSpell, 0, 0, 0)
			e.State.Players[0].Hero.CurrentLife = 5
			if shelter {
				u := put(t, e, "1021001", 1, 0, 0)
				u.Statuses["防止致命"] = 1
				u.Statuses["light_shelter_until"] = 99
			} else {
				e.State.Players[1].TempModifiers = append(e.State.Players[1].TempModifiers, game.TemporaryModifier{Type: game.TempModAllSpellDamageZeroUntilEnd})
			}
			incomingHeroSpell(t, e, incomingSpell)
			d := next(t, e, p, 0)
			if d.Action.Action != "defend" {
				t.Fatalf("%s did not defend under protection, shelter=%v: %+v", id, shelter, d)
			}
			if err := e.HandleAction(0, d.Action); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestThunderlightReservesAttackBeforeFurtherLearning(t *testing.T) {
	e, p := thunderScenario(t)
	player := e.State.Players[0]
	put(t, e, "3321002", 0, 0, 0)
	player.Elements["气"] = 3
	player.Hero.IsHorizontal = true
	player.SkillPool = []*game.CardInstance{game.NewCardInstance(cards.PlayableCardDB["3021101"], 0, 1)}
	d := applyNext(t, e, p)
	if d.Action.Action != "cast_spell" {
		t.Fatalf("spent the full attack budget on learning: %+v", d)
	}
}

func TestThunderlightUsesLimitedBoostForLethal(t *testing.T) {
	e, p := thunderScenario(t)
	player := e.State.Players[0]
	put(t, e, "3321002", 0, 0, 0)
	player.Skills[1] = game.NewCardInstance(cards.PlayableCardDB["3321003"], 0, 1)
	player.Skills[1].IsHorizontal = false
	other := e.State.Players[1]
	other.Hero.CurrentLife = 1
	other.Hero.IsHorizontal = false
	for i, n := range []string{"3121001", "3121003"} {
		other.Skills[i] = game.NewCardInstance(cards.PlayableCardDB[n], 1, 1)
		other.Skills[i].IsHorizontal = false
	}
	player.Elements["气"] = 2
	d := applyNext(t, e, p)
	if d.Action.Action != "cast_spell" || len(d.Action.Data["boost_ids"].([]any)) != 1 {
		t.Fatalf("did not choose one affordable boost: %+v", d)
	}
	if e.State.PendingSpell == nil || e.State.PendingSpell.TotalPower <= 7 {
		t.Fatal("boost did not exceed public payable defense")
	}
}

func TestThunderlightArcaneImpactUsesSpareLight(t *testing.T) {
	e, p := thunderScenario(t)
	player := e.State.Players[0]
	skill := put(t, e, "3021101", 0, 0, 0)
	player.Elements["气"] = 2
	player.Elements["光"] = 2
	d := applyNext(t, e, p)
	if d.Action.Action != "cast_spell" || text(d.Action.Data["instance_id"]) != skill.InstanceID {
		t.Fatalf("did not cast arcane impact: %+v", d)
	}
	if player.Elements["气"] != 2 || player.Elements["光"] != 0 {
		t.Fatalf("did not preserve air with spare light: %v", player.Elements)
	}
}

func TestDecisionPreviewDoesNotConsumeOneShotProtection(t *testing.T) {
	e, _ := thunderScenario(t)
	skill := put(t, e, "3321002", 0, 0, 0)
	e.State.Players[1].TempModifiers = append(e.State.Players[1].TempModifiers, game.TemporaryModifier{Type: game.TempModAllSpellDamageZero, TargetInstanceID: skill.InstanceID, RemainingUses: 1})
	before, _ := json.Marshal(e.GetStateForPlayer(0))
	for i := 0; i < 2; i++ {
		obs, err := e.ObserveDecision(0)
		if err != nil {
			t.Fatal(err)
		}
		for _, target := range obs.Cards[skill.InstanceID].Targets {
			for _, hit := range target.Impacts {
				if hit.Damage != 0 {
					t.Fatal("ignored one-shot damage prevention")
				}
			}
		}
	}
	after, _ := json.Marshal(e.GetStateForPlayer(0))
	if string(before) != string(after) || e.State.Players[1].TempModifiers[0].RemainingUses != 1 {
		t.Fatal("decision query consumed protection")
	}
}
