package game

import (
	"encoding/json"
	"eraofarcane/model"
	"reflect"
	"strings"
	"testing"
)

func TestSeededEngineIsIndependentOfOtherMatches(t *testing.T) {
	setupBaseCardSmokeSuite(t)
	deck, err := model.ParseDeckCode(testDeckCode)
	if err != nil {
		t.Fatal(err)
	}
	const seed int64 = 91735028814231
	left := NewEngineWithSeed("replay", nil, seed)
	right := NewEngineWithSeed("replay", nil, seed)
	noise := NewEngineWithSeed("other-room", nil, 42)
	if err := left.SetupGame("a", deck, "b", deck); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 100; i++ {
		noise.randomIntn(31)
		noise.newCardInstance(baseCard(t, "1021001"), 0, 0)
		NewCardInstance(baseCard(t, "1021001"), 0, 0)
	}
	if err := right.SetupGame("a", deck, "b", deck); err != nil {
		t.Fatal(err)
	}
	for player := 0; player < 2; player++ {
		left.shuffleDeck(player)
		noise.randomIntn(17)
		right.shuffleDeck(player)
		if !reflect.DeepEqual(left.State.Players[player].Deck, right.State.Players[player].Deck) {
			t.Fatal("deck order or IDs depend on another game")
		}
		if !reflect.DeepEqual(left.GetStateForPlayer(player), right.GetStateForPlayer(player)) {
			t.Fatal("same seed did not reproduce player view")
		}
	}
	if !reflect.DeepEqual(left.DebugResolutionTrace(), right.DebugResolutionTrace()) {
		t.Fatal("same execution did not reproduce trace")
	}
	for _, view := range []map[string]any{left.GetStateForPlayer(0), left.GetStateForSpectator()} {
		encoded, err := json.Marshal(view)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(encoded), "91735028814231") || strings.Contains(string(encoded), "random_seed") || strings.Contains(string(encoded), "resolution_trace") {
			t.Fatal("private replay data leaked to a game view")
		}
	}
}

func TestResolutionTraceRecordsParentAndBoundsMemory(t *testing.T) {
	e := NewEngineWithSeed("trace", nil, 1)
	e.runResolution("parent", func() { e.runResolution("child", func() {}) })
	trace := e.DebugResolutionTrace()
	var parent, child uint64
	for _, entry := range trace {
		if entry.Kind == "frame_start" && entry.Name == "parent" {
			parent = entry.Frame
		}
		if entry.Kind == "frame_start" && entry.Name == "child" {
			child = entry.Frame
			if entry.Parent != parent {
				t.Fatalf("wrong parent: %+v", entry)
			}
		}
	}
	if parent == 0 || child == 0 || parent == child || e.activeFrame != nil {
		t.Fatal("invalid frame ownership")
	}
	for i := 0; i < resolutionTraceCapacity+20; i++ {
		e.traceResolution("test", "", nil)
	}
	trace = e.DebugResolutionTrace()
	if len(trace) != resolutionTraceCapacity {
		t.Fatalf("unbounded trace: %d", len(trace))
	}
	for i := 1; i < len(trace); i++ {
		if trace[i].Sequence != trace[i-1].Sequence+1 {
			t.Fatal("trace lost execution order")
		}
	}
	trace[0].Kind = "mutated"
	if e.DebugResolutionTrace()[0].Kind == "mutated" {
		t.Fatal("trace accessor exposed internal memory")
	}
}
