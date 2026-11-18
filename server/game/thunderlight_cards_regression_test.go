package game

import (
	"eraofarcane/model"
	"testing"
)

func TestDawnThresholdIncludesAurasAndRemovesThem(t *testing.T) {
	e := setupReportedBugEngine(t)
	p := e.State.Players[0]
	p.Hero = NewCardInstance(baseCard(t, "4311101"), 0, 1)
	for i, n := range []string{"2021105", "2321101", "2321105", "2321105"} {
		p.Equipment[i] = NewCardInstance(baseCard(t, n), 0, 1)
	}
	skill := readySkill(baseCard(t, "3511010"), 0)
	p.Skills[0] = skill
	for _, bonus := range []int{2, 3} {
		skill.PowerBonus = bonus
		err := e.validateSkillForPurpose(skill, skillPurposeAttack)
		if (err == nil) != (bonus == 3) {
			t.Fatalf("threshold power=%d error=%v", e.effectiveSpellPower(0, skill, nil), err)
		}
	}
	p.Equipment[3].Statuses[StatusPetrify] = 1
	if err := e.validateSkillForPurpose(skill, skillPurposeAttack); err == nil {
		t.Fatal("lost armor aura still satisfied Dawn attack threshold")
	}
	if err := e.validateSkillForPurpose(skill, skillPurposeDefend); err != nil {
		t.Fatalf("attack threshold restricted defense: %v", err)
	}
}

func TestRocDrawsNotifyDandelionAndKeepExpiry(t *testing.T) {
	e := setupReportedBugEngine(t)
	p := e.State.Players[0]
	dandelion := NewCardInstance(baseCard(t, "1321003"), 0, e.State.TurnNumber)
	breath := NewCardInstance(baseCard(t, "2321103"), 0, e.State.TurnNumber)
	beast := NewCardInstance(baseCard(t, "1321006"), 0, e.State.TurnNumber)
	p.Deck = []*CardInstance{dandelion, breath, beast}
	roc := NewCardInstance(baseCard(t, "1311001"), 0, e.State.TurnNumber)
	if err := (Card1311001Roc{}).OnEnter(e.skillContext(0, roc)); err != nil {
		t.Fatal(err)
	}
	if !p.RevealedHand[dandelion.InstanceID] || p.DrawnTurn[dandelion.InstanceID] != e.State.TurnNumber || p.DrawCountThisTurn != 2 {
		t.Fatal("roc draws missed ordinary draw bookkeeping")
	}
	p.Hand = []*CardInstance{breath}
	if err := (Card1321003MagicDandelion{}).OnEnter(e.skillContext(0, dandelion)); err != nil {
		t.Fatal(err)
	}
	if len(p.Deck) != 0 || len(p.Hand) != 2 {
		t.Fatal("roc-drawn dandelion did not draw")
	}
	e.discardMarkedEndOfTurnCards(p)
	if len(p.Hand) != 1 || p.Hand[0] != beast || p.Elements[model.ElementAir] != 1 {
		t.Fatal("roc expiry discarded a played card or missed discarded breath")
	}
}

func TestDisarmUsesHitTargetController(t *testing.T) {
	for _, owner := range []int{0, 1} {
		e := setupReportedBugEngine(t)
		for player := 0; player < 2; player++ {
			placeUnit(baseCard(t, "1021001"), player, 1, 1, e)
			e.State.Players[player].Equipment[0] = NewCardInstance(baseCard(t, "2321101"), player, 1)
		}
		p := e.State.Players[0]
		skill := readySkill(baseCard(t, "3021008"), 0)
		p.Skills[0] = skill
		p.Elements[model.ElementAir] = 2
		if err := e.HandleAction(0, ActionMessage{Action: "cast_spell", Data: map[string]any{"instance_id": skill.InstanceID, "target_type": "unit", "target_col": float64(1), "target_row": float64(1), "target_owner": float64(owner)}}); err != nil {
			t.Fatal(err)
		}
		if e.State.Phase == PhaseDefenseWindow {
			if err := e.HandleAction(1, ActionMessage{Action: "no_defend"}); err != nil {
				t.Fatal(err)
			}
		}
		pa := e.State.PendingAction
		item := e.State.Players[owner].Equipment[0]
		if pa == nil || len(pa.Candidates) != 1 || pa.Candidates[0]["instance_id"] != item.InstanceID {
			t.Fatalf("wrong controller equipment: owner=%d pending=%+v", owner, pa)
		}
		resolvePendingSelection(t, e, 0, item.InstanceID)
		if e.State.Players[owner].Equipment[0] != nil || e.State.Players[1-owner].Equipment[0] == nil {
			t.Fatal("destroyed the wrong controller's equipment")
		}
	}
}

func TestChimeraDiscountExcludesSorceriesAndLearnCosts(t *testing.T) {
	e := setupReportedBugEngine(t)
	placeUnit(baseCard(t, "1321010"), 0, 0, 0, e)
	p := e.State.Players[0]
	for _, tc := range []struct {
		number string
		want   int
	}{{"3321002", 0}, {"3321007", 1}, {"3521105", 2}} {
		skill := readySkill(baseCard(t, tc.number), 0)
		if got := totalCost(e.effectiveSkillUseCost(p, skill)); got != tc.want {
			t.Fatalf("%s cost=%d want %d", tc.number, got, tc.want)
		}
		if got := totalCost(e.effectiveSkillLearnCost(p, skill)); got != totalCost(skill.Card.ElementsCost) {
			t.Fatalf("discounted learn cost: %s", tc.number)
		}
	}
}

func TestDurationDamageProtectionsAreContinuous(t *testing.T) {
	for _, source := range []string{"1221010", "2421010"} {
		for _, first := range []int{0, 1} {
			e := setupReportedBugEngine(t)
			e.State.FirstPlayer = first
			e.State.CurrentTurn = first
			c := NewCardInstance(baseCard(t, source), first, 1)
			ctx := e.skillContext(first, c)
			if source == "1221010" {
				_ = (Card1221010WallKeeper{}).OnEnter(ctx)
			} else {
				_ = (Card2421010NatureSealScroll{}).OnUseItem(ctx)
			}
			spell := readySkill(baseCard(t, "3321002"), first)
			for tick := 0; tick < 2; tick++ {
				if got := e.effectiveSpellDamage(first, spell, 1, nil); got != 0 {
					t.Fatalf("%s stopped protecting after %d casts", source, tick)
				}
				e.consumeAllSpellDamageZero(e.State.Players[first], spell)
			}
			mod := e.State.Players[first].TempModifiers[0]
			want := e.turnOrdinal() + 1
			if source == "2421010" {
				want++
			}
			if mod.ExpiresAtOrdinal != want {
				t.Fatalf("wrong %s expiry: %+v", source, mod)
			}
			e.clearTemporaryModifiersAtTurnEnd(first)
			if len(e.State.Players[first].TempModifiers) != 1 {
				t.Fatal("expired in the casting turn")
			}
			e.State.CurrentTurn = 1 - first
			e.clearTemporaryModifiersAtTurnEnd(1 - first)
			if source == "1221010" && len(e.State.Players[first].TempModifiers) != 0 {
				t.Fatal("wall keeper outlived immediately following turn")
			}
			if source == "2421010" {
				if len(e.State.Players[first].TempModifiers) != 1 {
					t.Fatal("seal expired before caster's next turn")
				}
				e.State.CurrentTurn = first
				e.State.TurnNumber++
				e.clearTemporaryModifiersAtTurnEnd(first)
				if len(e.State.Players[first].TempModifiers) != 0 {
					t.Fatal("seal outlived caster's next turn")
				}
			}
		}
	}
	e := setupReportedBugEngine(t)
	u := placeUnit(baseCard(t, "1321006"), 0, 0, 0, e)
	c := readySkill(baseCard(t, "3521011"), 0)
	ctx := &EffectContext{Engine: e, Source: c, Target: u, PlayerID: 0, OpponentID: 1, ExtraData: map[string]any{"spell_source": c}}
	_ = (Card3521011LightShelter{}).OnSpellHit(ctx)
	before := u.CurrentLife
	for i := 0; i < 2; i++ {
		e.ApplyDamage(DamageRequest{Target: u, Amount: 20})
		if u.CurrentLife != before {
			t.Fatal("shelter failed repeated lethal prevention")
		}
	}
	e.State.CurrentTurn = 1
	e.expireTimedDamageProtection()
	if e.lightShelterActive(u) {
		t.Fatal("shelter did not expire after next turn")
	}
	e.ApplyDamage(DamageRequest{Target: u, Amount: 20})
	if u.CurrentLife > 0 {
		t.Fatal("expired shelter still prevents lethal")
	}
}
