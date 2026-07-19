package main

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestWebsocketURL(t *testing.T) {
	got, err := websocketURL("http://127.0.0.1:9090", "0042", "agent a", "风 Agent", "hero // main // skills")
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := url.Parse(got)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Scheme != "ws" || parsed.Host != "127.0.0.1:9090" || parsed.Path != "/ws" {
		t.Fatalf("unexpected URL: %s", got)
	}
	query := parsed.Query()
	if query.Get("room") != "0042" || query.Get("player_id") != "agent a" ||
		query.Get("player_name") != "风 Agent" || query.Get("deck_code") != "hero // main // skills" {
		t.Fatalf("query was not preserved: %#v", query)
	}
}

func TestInitDataPreservesExistingKnowledge(t *testing.T) {
	root := t.TempDir()
	existing := filepath.Join(root, "knowledge", "core-rules.md")
	if err := os.MkdirAll(filepath.Dir(existing), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(existing, []byte("keep me"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := initData([]string{"-root", root}); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(existing)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "keep me" {
		t.Fatalf("existing knowledge was overwritten: %q", content)
	}
	for _, relative := range []string{
		"knowledge/gameplay-principles.md",
		"knowledge/deck-lab.md",
		"knowledge/open-questions.md",
		"knowledge/retired-lessons.md",
		"context-packs/next-match.md",
	} {
		if _, err := os.Stat(filepath.Join(root, relative)); err != nil {
			t.Errorf("%s was not initialized: %v", relative, err)
		}
	}
}

func TestValidateAction(t *testing.T) {
	valid := []string{
		`{"action":"mulligan","data":{"keep":true}}`,
		`{"action":"end_turn","data":{}}`,
	}
	for _, input := range valid {
		if err := validateAction([]byte(input)); err != nil {
			t.Errorf("validateAction(%s): %v", input, err)
		}
	}
	invalid := []string{
		`not json`,
		`{"data":{}}`,
		`{"action":"end_turn"}`,
	}
	for _, input := range invalid {
		if err := validateAction([]byte(input)); err == nil {
			t.Errorf("validateAction(%s) unexpectedly succeeded", input)
		}
	}
}

func TestDecisionRelevant(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{`{"type":"joined","data":{}}`, true},
		{`{"type":"error","message":"bad action"}`, true},
		{`{"type":"game_event","event":{"type":"state_sync","data":{}}}`, true},
		{`{"type":"game_event","event":{"type":"draw_card","data":{}}}`, false},
	}
	for _, test := range tests {
		if got := decisionRelevant([]byte(test.input)); got != test.want {
			t.Errorf("decisionRelevant(%s) = %v, want %v", test.input, got, test.want)
		}
	}
}
