package bot

import (
	"encoding/json"
	"eraofarcane/cards"
	"eraofarcane/game"
	"fmt"
	"sync"
	"testing"
)

var loadOnce sync.Once

func fixture(t *testing.T, seed int64) *game.Engine {
	t.Helper()
	loadOnce.Do(func() {
		if err := cards.LoadCards(); err != nil {
			panic(err)
		}
		game.SetCardDB(cards.PlayableCardDB)
		game.RegisterAllCardEffects()
	})
	d, err := Deck()
	if err != nil {
		t.Fatal(err)
	}
	if err = d.Validate(cards.PlayableCardDB); err != nil {
		t.Fatal(err)
	}
	e := game.NewEngineWithSeed(fmt.Sprintf("bot-%d", seed), nil, seed)
	if err = e.SetupGameWithFirstPlayer("A", d, "B", d, 0); err != nil {
		t.Fatal(err)
	}
	return e
}
func scenario(t *testing.T) *game.Engine {
	e := fixture(t, 17)
	e.State.Phase = game.PhaseMain
	e.State.CurrentTurn = 0
	e.State.TurnNumber = 2
	e.State.IsFirstTurn = false
	for _, p := range e.State.Players {
		p.Hand = nil
		p.SkillPool = nil
		for k := range p.Elements {
			p.Elements[k] = 0
		}
	}
	return e
}
func put(t *testing.T, e *game.Engine, number string, player, col, row int) *game.CardInstance {
	t.Helper()
	c := game.NewCardInstance(cards.PlayableCardDB[number], player, 1)
	c.IsHorizontal = false
	if c.Card.IsSkill() {
		e.State.Players[player].Skills[0] = c
	} else {
		c.Position = &game.Position{Col: col, Row: row}
		e.State.Players[player].Units[col][row] = c
	}
	return c
}
func next(t *testing.T, e *game.Engine, p *Policy, player int) Decision {
	t.Helper()
	obs, err := e.ObserveDecision(player)
	if err != nil {
		t.Fatal(err)
	}
	d, ok, err := p.Next(obs)
	if err != nil || !ok {
		t.Fatalf("no decision: %v", err)
	}
	return d
}

func TestMirrorCompletesAndLearnsFromSecondOwnTurn(t *testing.T) {
	for _, seed := range []int64{20261117, 1, 2, 3, 4, 5, 19, 42} {
		t.Run(fmt.Sprint(seed), func(t *testing.T) {
			e := fixture(t, seed)
			policies := [2]*Policy{New(), New()}
			learnedSecond := false
			rejected := 0
			for step := 0; step < 1800; step++ {
				if e.State.Phase == game.PhaseGameOver {
					if !learnedSecond {
						t.Fatal("never learned on own turn two")
					}
					if rejected > 8 {
						t.Fatalf("too many invalid actions: %d", rejected)
					}
					return
				}
				player := e.State.CurrentTurn
				if e.State.Phase == game.PhaseMulligan {
					if e.State.MulliganDone[0] {
						player = 1
					} else {
						player = 0
					}
				} else {
					view := e.GetStateForPlayer(0)
					player = view["decision_player"].(int)
				}
				d := next(t, e, policies[player], player)
				if d.Action.Action == "surrender" {
					t.Fatalf("choice fallback surrendered at step %d: %s", step, d.Reason)
				}
				if d.Action.Action == "learn_skill" {
					if e.State.TurnNumber == 1 {
						t.Fatal("learned on first own turn")
					}
					if e.State.TurnNumber == 2 {
						learnedSecond = true
					}
				}
				if err := e.HandleAction(player, d.Action); err != nil {
					policies[player].Reject(d)
					rejected++
					t.Logf("rejected step %d: %s: %v", step, d.Action.Action, err)
				} else {
					policies[player].Accept()
				}
			}
			t.Fatalf("mirror did not finish, phase=%s round=%d", e.State.Phase, e.State.TurnNumber)
		})
	}
}

func TestWindUsesHeroConversionForRearOneLifeTarget(t *testing.T) {
	e := scenario(t)
	put(t, e, "3121015", 0, 0, 0)
	front := put(t, e, "1121005", 1, 0, 0)
	front.CurrentLife = 4
	rear := put(t, e, "1021001", 1, 2, 2)
	p := New()
	converted := false
	for i := 0; i < 6; i++ {
		d := next(t, e, p, 0)
		if d.Action.Action == "use_ability" && d.Action.Data["instance_id"] == e.State.Players[0].Hero.InstanceID {
			converted = true
		}
		if err := e.HandleAction(0, d.Action); err != nil {
			t.Fatal(err)
		}
		p.Accept()
		if d.Action.Action == "cast_spell" {
			if !converted {
				t.Fatal("did not convert hero load to obtain air payment")
			}
			if number(d.Action.Data["target_col"]) != rear.Position.Col || number(d.Action.Data["target_row"]) != rear.Position.Row {
				t.Fatal("wind did not target rear one-life unit")
			}
			return
		}
	}
	t.Fatal("wind never cast")
}

func TestPegasusIsTargetedBeforeProtectedHero(t *testing.T) {
	e := scenario(t)
	put(t, e, "3121001", 0, 0, 0)
	pegasus := put(t, e, "1511003", 1, 0, 2)
	e.State.Players[1].Hero.CurrentLife = 1
	p := New()
	for i := 0; i < 4; i++ {
		d := next(t, e, p, 0)
		if d.Action.Action == "cast_spell" {
			if number(d.Action.Data["target_col"]) != pegasus.Position.Col || number(d.Action.Data["target_row"]) != pegasus.Position.Row {
				t.Fatal("wasted a spell on a Pegasus-protected ally")
			}
			return
		}
		if err := e.HandleAction(0, d.Action); err != nil {
			t.Fatal(err)
		}
		p.Accept()
	}
	t.Fatal("never attacked Pegasus")
}

func TestConvertedFirstTurnHeroChoosesReducedMixedGain(t *testing.T) {
	e := scenario(t)
	e.State.TurnNumber = 1
	e.State.IsFirstTurn = true
	put(t, e, "3121015", 0, 0, 0)
	put(t, e, "1021001", 1, 2, 2)
	p := New()
	converted := false
	for step := 0; step < 5; step++ {
		d := next(t, e, p, 0)
		if d.Action.Action == "use_ability" {
			converted = true
		}
		if d.Action.Action == "consume" {
			gain := elements(d.Action.Data["gain"])
			if !converted || gain["火"] != 1 || gain["无"] != 1 || sum(gain) != 2 {
				t.Fatalf("wrong reduced first gain: %+v", d.Action.Data)
			}
		}
		if err := e.HandleAction(0, d.Action); err != nil {
			t.Fatal(err)
		}
		p.Accept()
		if d.Action.Action == "cast_spell" {
			return
		}
	}
	t.Fatal("did not finish first-turn conversion/payment")
}

func TestUnknownMandatoryChoiceRetriesDifferentDefault(t *testing.T) {
	e := scenario(t)
	e.SetPendingActionWithError(0, "unfamiliar_effect", "opaque", []map[string]any{{"instance_id": "a"}, {"instance_id": "b"}}, 1, 1, nil, false, func(ids []string, _ map[string]any) error {
		if ids[0] != "b" {
			return fmt.Errorf("choose another candidate")
		}
		return nil
	})
	p := New()
	first := next(t, e, p, 0)
	if err := e.HandleAction(0, first.Action); err == nil {
		t.Fatal("expected first default to be rejected")
	}
	p.Reject(first)
	second := next(t, e, p, 0)
	if second.Action.Action != "resolve_action" {
		t.Fatal("did not retry mandatory choice")
	}
	if err := e.HandleAction(0, second.Action); err != nil {
		t.Fatal(err)
	}
	if e.State.PendingAction != nil {
		t.Fatal("choice did not complete")
	}
}

func TestEmbersDeclarationResolvesBeforeCleanup(t *testing.T) {
	e := scenario(t)
	embers := put(t, e, "3121105", 0, 0, 0)
	e.State.Players[0].Hero.IsHorizontal = true
	p := New()
	d := next(t, e, p, 0)
	if d.Action.Action != "end_turn" {
		t.Fatalf("expected declaration, got %s", d.Action.Action)
	}
	if err := e.HandleAction(0, d.Action); err != nil {
		t.Fatal(err)
	}
	p.Accept()
	d = next(t, e, p, 0)
	if d.Action.Action != "resolve_action" {
		t.Fatal("did not select free Embers target")
	}
	if err := e.HandleAction(0, d.Action); err != nil {
		t.Fatal(err)
	}
	if e.State.PendingSpell == nil || !embers.IsHorizontal || e.State.CurrentTurn != 0 {
		t.Fatal("did not wait for free spell before turn cleanup")
	}
}

func TestObservationIsPrivateDetachedAndPure(t *testing.T) {
	e := scenario(t)
	e.State.Players[1].Hand = append(e.State.Players[1].Hand, game.NewCardInstance(cards.PlayableCardDB["1111001"], 1, 1))
	before, _ := json.Marshal(e.GetStateForPlayer(0))
	obs, err := e.ObserveDecision(0)
	if err != nil {
		t.Fatal(err)
	}
	after, _ := json.Marshal(e.GetStateForPlayer(0))
	if string(before) != string(after) {
		t.Fatal("observation changed game state")
	}
	var s map[string]any
	_ = json.Unmarshal(obs.State, &s)
	op := s["opponent"].(map[string]any)
	if _, ok := op["hand"]; ok {
		t.Fatal("opponent hand leaked")
	}
	if _, ok := obs.Cards[e.State.Players[1].Hand[0].InstanceID]; ok {
		t.Fatal("opponent private options leaked")
	}
	o := obs.Cards[e.State.Players[0].Hero.InstanceID]
	o.ConsumeGain["火"] = 999
	if e.State.Players[0].Hero.Card.ElementsGain["火"] != 4 {
		t.Fatal("snapshot retained engine-owned mutable maps")
	}
}

func TestCheapNonLethalDamageDoesNotSpendOffenseOnDefense(t *testing.T) {
	e := scenario(t)
	put(t, e, "3121001", 0, 0, 0)
	var s state
	obs, _ := e.ObserveDecision(0)
	_ = json.Unmarshal(obs.State, &s)
	s.Phase = "defense_window"
	s.Decision = 0
	raw, _ := json.Marshal(map[string]any{"power": 1, "skill": map[string]any{"attack": 1}, "affected_positions": []map[string]any{{"owner_id": 0, "col": 1, "row": 1}}})
	_ = json.Unmarshal(raw, &s.Spell)
	d := New().defend(s, obs.Cards)
	if d.Action.Action != "no_defend" {
		t.Fatal("spent offense defending minor hero damage")
	}
}

func incomingHeroSpell(t *testing.T, e *game.Engine, number string) {
	t.Helper()
	e.State.CurrentTurn = 1
	skill := put(t, e, number, 1, 0, 0)
	e.State.Players[1].Elements["火"] = 10
	if err := e.HandleAction(1, action("cast_spell", map[string]any{"instance_id": skill.InstanceID, "target_type": "unit", "target_col": 1, "target_row": 1})); err != nil {
		t.Fatal(err)
	}
	if e.State.Phase != game.PhaseDefenseWindow {
		t.Fatalf("expected defense window, got %s", e.State.Phase)
	}
}

func TestDefenseIncludesOutstandingAndIncomingBurn(t *testing.T) {
	for _, tc := range []struct {
		name                                string
		life, burn, shield                  int
		spell                               string
		cheapSource, exhausted, wantDefense bool
	}{
		{"existing burn makes one damage lethal", 2, 1, 0, "3121001", false, false, true},
		{"ray has three total damage", 3, 0, 0, "3121003", false, false, true},
		{"ray adds to existing burn", 4, 1, 0, "3121003", false, false, true},
		{"healthy hero blocks efficient ray cheaply", 5, 0, 0, "3121003", true, false, true},
		{"shield does not absorb ray burn", 2, 0, 1, "3121003", false, false, true},
		{"safe ray does not justify hero overexertion", 5, 0, 0, "3121003", false, false, false},
		{"ordinary chip still preserves offense", 5, 0, 0, "3121001", true, false, false},
		{"exhausted resources cannot defend room 2683", 2, 1, 0, "3121001", false, true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			e := scenario(t)
			hero := e.State.Players[0].Hero
			hero.CurrentLife = tc.life
			hero.Statuses[game.StatusBurn] = tc.burn
			hero.IsHorizontal = tc.exhausted
			e.State.Players[0].Shield = tc.shield
			put(t, e, tc.spell, 0, 0, 0)
			if tc.cheapSource {
				seed := put(t, e, "1121111", 0, 0, 2)
				seed.ElementsGainBonus["火"] = 1
				seed.UltimateUsed = true
				// Room 2683 also had two one-fire sources. Paying with the
				// two-fire seed alone keeps this defense within the same budget.
				put(t, e, "1121001", 0, 0, 1)
				put(t, e, "1121012", 0, 2, 2)
			}
			incomingHeroSpell(t, e, tc.spell)
			d := next(t, e, New(), 0)
			want := "no_defend"
			if tc.wantDefense {
				want = "defend"
			}
			if d.Action.Action != want {
				t.Fatalf("got %s, want %s: %s", d.Action.Action, want, d.Reason)
			}
			if err := e.HandleAction(0, d.Action); err != nil {
				t.Fatalf("decision was not legal: %v", err)
			}
			if tc.cheapSource && tc.wantDefense {
				if !e.State.Players[0].Units[0][2].IsHorizontal || e.State.Players[0].Units[0][1].IsHorizontal || e.State.Players[0].Units[2][2].IsHorizontal || hero.IsHorizontal {
					t.Fatal("cheap defense must overexert only the two-fire seed")
				}
			}
			if tc.wantDefense && (hero.CurrentLife != tc.life || hero.Statuses[game.StatusBurn] != tc.burn) {
				t.Fatal("defense must block new damage and burn while preserving existing burn")
			}
		})
	}
}

func TestIncomingBurnObservationIsPureAndTicksRemainSequential(t *testing.T) {
	e := scenario(t)
	hero := e.State.Players[0].Hero
	hero.CurrentLife = 5
	incomingHeroSpell(t, e, "3121003")
	before, _ := json.Marshal(e.GetStateForPlayer(0))
	for i := 0; i < 2; i++ {
		obs, err := e.ObserveDecision(0)
		if err != nil {
			t.Fatal(err)
		}
		if obs.IncomingBurn != 2 {
			t.Fatalf("ray preview burn=%d, want 2", obs.IncomingBurn)
		}
	}
	after, _ := json.Marshal(e.GetStateForPlayer(0))
	if string(before) != string(after) {
		t.Fatal("reading incoming burn changed game state")
	}
	if err := e.HandleAction(0, action("no_defend", nil)); err != nil {
		t.Fatal(err)
	}
	if hero.CurrentLife != 4 || hero.Statuses[game.StatusBurn] != 2 {
		t.Fatal("ray must initially deal one damage and add two burn, not apply three immediate damage")
	}
	if err := e.HandleAction(1, action("end_turn", nil)); err != nil {
		t.Fatal(err)
	}
	if err := e.HandleAction(0, action("end_turn", nil)); err != nil {
		t.Fatal(err)
	}
	if hero.CurrentLife != 3 || hero.Statuses[game.StatusBurn] != 1 {
		t.Fatal("burn must settle one point per owner's turn end")
	}
}
