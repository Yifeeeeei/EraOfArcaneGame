package bot

import (
	"encoding/json"
	"eraofarcane/game"
)

type card struct {
	ID             string         `json:"instance_id"`
	Number         string         `json:"number"`
	Name           string         `json:"name"`
	Type           string         `json:"type"`
	Category       string         `json:"category"`
	Life           int            `json:"current_life"`
	MaxLife        int            `json:"life"`
	Attack         int            `json:"current_attack"`
	PrintedAttack  int            `json:"attack"`
	Power          int            `json:"power"`
	Horizontal     bool           `json:"is_horizontal"`
	Position       *game.Position `json:"position"`
	Statuses       map[string]int `json:"statuses"`
	Gain           map[string]int `json:"elements_gain"`
	Cost           map[string]int `json:"effective_elements_cost"`
	LearnCost      map[string]int `json:"effective_learn_cost"`
	Equipment      bool           `json:"is_equipment"`
	Scroll         bool           `json:"is_spell_scroll"`
	Counter        bool           `json:"is_counter_trap"`
	SetCounter     bool           `json:"is_set_counter"`
	CanPlace       bool           `json:"can_place_without_replacement"`
	ReplacementIDs []string       `json:"equipment_replacement_ids"`
	Subtype        string         `json:"equipment_subtype"`
	Sorcery        bool           `json:"is_sorcery"`
	Pierce         bool           `json:"has_pierce"`
	Stealth        bool           `json:"has_stealth"`
	Used           int            `json:"used_this_turn"`
	UltimateUsed   bool           `json:"ultimate_used"`
	DefensePower   int            `json:"effective_defense_power"`
	Bound          []*card        `json:"bound_skills"`
	Area           string         `json:"spell_area"`
}

type player struct {
	ID         int            `json:"player_id"`
	Hero       *card          `json:"hero"`
	Units      [3][3]*card    `json:"units"`
	Hand       []*card        `json:"hand"`
	Skills     []*card        `json:"skills"`
	Pool       []*card        `json:"skill_pool"`
	Equipment  []*card        `json:"equipment"`
	Elements   map[string]int `json:"elements"`
	Strict     int            `json:"strict_arcane"`
	Light      bool           `json:"light_wildcard"`
	LightCost  bool           `json:"light_cost_wildcard"`
	Shield     int            `json:"shield"`
	SkillSlots int            `json:"skill_slot_capacity"`
}

type pending struct {
	Type       string           `json:"type"`
	ID         uint64           `json:"decision_id"`
	Candidates []map[string]any `json:"candidates"`
	Min        int              `json:"min_select"`
	Max        int              `json:"max_select"`
	UI         map[string]any   `json:"ui"`
	Context    map[string]any   `json:"context"`
	Cost       map[string]int   `json:"cost"`
	Overexert  bool             `json:"can_overexert"`
}

type state struct {
	Revision uint64   `json:"revision"`
	Phase    string   `json:"phase"`
	Turn     int      `json:"current_turn"`
	Number   int      `json:"turn_number"`
	First    int      `json:"first_player"`
	Decision int      `json:"decision_player"`
	Mulligan [2]bool  `json:"mulligan_done"`
	You      player   `json:"you"`
	Opponent player   `json:"opponent"`
	Pending  *pending `json:"pending_action"`
	Spell    *struct {
		Power    int   `json:"power"`
		Skill    *card `json:"skill"`
		Affected []struct {
			Owner int `json:"owner_id"`
			Col   int `json:"col"`
			Row   int `json:"row"`
		} `json:"affected_positions"`
	} `json:"pending_spell"`
	DrawOffer    int `json:"draw_offer_by"`
	IncomingBurn int `json:"-"`
}

type Decision struct {
	Action game.ActionMessage `json:"action"`
	Reason string             `json:"reason"`
}

type Policy struct {
	turn    int
	blocked map[string]bool
	steps   int
}

func New() *Policy { return &Policy{blocked: map[string]bool{}} }

func Fingerprint(a game.ActionMessage) string {
	data, _ := json.Marshal(struct {
		Action string
		Data   map[string]any
	}{a.Action, a.Data})
	return string(data)
}

// Reject blacklists the attempted payload, not the revision: rejected actions
// can advance revision without changing the facts that made them invalid.
func (p *Policy) Reject(d Decision) { p.blocked[Fingerprint(d.Action)] = true }

func (p *Policy) Accept() { p.blocked = map[string]bool{} }

func action(name string, data map[string]any) game.ActionMessage {
	if data == nil {
		data = map[string]any{}
	}
	// Normalize exactly as WebSocket JSON decoding does (arrays and numbers).
	raw, _ := json.Marshal(game.ActionMessage{Action: name, Data: data})
	var a game.ActionMessage
	_ = json.Unmarshal(raw, &a)
	return a
}

func sum(m map[string]int) int {
	n := 0
	for _, v := range m {
		if v > 0 {
			n += v
		}
	}
	return n
}
func plus(a, b map[string]int) map[string]int {
	c := map[string]int{}
	for k, v := range a {
		c[k] = v
	}
	for k, v := range b {
		c[k] += v
	}
	return c
}
func number(v any) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case uint64:
		return int(n)
	}
	return 0
}
func text(v any) string          { s, _ := v.(string); return s }
func asMap(v any) map[string]any { m, _ := v.(map[string]any); return m }
func elements(v any) map[string]int {
	m := map[string]int{}
	for k, v := range asMap(v) {
		m[k] = number(v)
	}
	return m
}
func (p player) units() []*card {
	var result []*card
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if c := p.Units[col][row]; c != nil {
				result = append(result, c)
			}
		}
	}
	return result
}
func (p player) field() []*card { return append(p.units(), p.Equipment...) }
func (p player) spells() []*card {
	result := append([]*card{}, p.Skills...)
	for _, c := range p.field() {
		if c != nil {
			result = append(result, c.Bound...)
		}
	}
	return result
}
func (p player) find(id string) *card {
	for _, list := range [][]*card{p.field(), p.Hand, p.spells(), p.Pool} {
		for _, c := range list {
			if c != nil && c.ID == id {
				return c
			}
		}
	}
	return nil
}
func (p player) front() int {
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if c := p.Units[col][row]; c != nil && !c.Stealth {
				return row
			}
		}
	}
	return -1
}

// TurnNumber counts full rounds, so it is also each player's own turn count.
func ownTurn(s state) int { return s.Number }
func payment(s state, available, cost map[string]int) (map[string]int, bool) {
	return game.DecisionPayment(available, cost, s.You.Strict, s.You.Light, s.You.LightCost)
}
