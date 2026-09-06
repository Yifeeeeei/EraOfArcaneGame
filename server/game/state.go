package game

import (
	"eraofarcane/model"
	"fmt"
	"sync/atomic"
)

// Position on the 3x3 unit grid
type Position struct {
	Col int `json:"col"` // 0-2, left to right
	Row int `json:"row"` // 0-2, 0=front row (closest to enemy)
}

func (p Position) Valid() bool {
	return p.Col >= 0 && p.Col < 3 && p.Row >= 0 && p.Row < 3
}

func (p Position) String() string {
	return fmt.Sprintf("(%d,%d)", p.Col, p.Row)
}

// CardInstance represents a card in play with runtime state
type CardInstance struct {
	InstanceID string      `json:"instance_id"`
	Card       *model.Card `json:"card"`
	OwnerID    int         `json:"owner_id"` // 0 or 1

	// Runtime state
	CurrentLife         int                `json:"current_life"`
	CurrentAttack       int                `json:"current_attack"`
	DamageTakenThisTurn int                `json:"damage_taken_this_turn,omitempty"`
	IsHorizontal        bool               `json:"is_horizontal"` // 横置=true, 竖置=false
	Statuses            map[string]int     `json:"statuses"`      // status -> stack count
	ElementsGainBonus   map[string]int     `json:"elements_gain_bonus,omitempty"`
	ElementsGainSet     map[string]int     `json:"elements_gain_set,omitempty"`
	PowerBonus          int                `json:"power_bonus,omitempty"`
	AttackBonus         int                `json:"attack_bonus,omitempty"`
	IsSetCounter        bool               `json:"is_set_counter,omitempty"`
	Position            *Position          `json:"position"`               // nil if not on unit grid
	SlotIndex           int                `json:"slot_index"`             // for skill/equipment slots
	EnterTurn           int                `json:"enter_turn"`             // which turn this card entered the field
	BoundSkills         []*CardInstance    `json:"bound_skills,omitempty"` // skills attached to this card, not skill slots
	UnderCards          []*CardInstance    `json:"under_cards,omitempty"`  // public cards placed under this card
	AttachedBehaviors   []AttachedBehavior `json:"-"`                      // runtime-granted behavior objects

	// Ability usage
	UsedThisTurn int  `json:"used_this_turn"` // for 回合技
	UltimateUsed bool `json:"ultimate_used"`  // for 绝技

	// Equipment uses
	UsesRemaining int `json:"uses_remaining"` // for 法宝
}

const StatusCannotUseSkillUntilTurn = "cannot_use_skill_until_turn"
const StatusEntryCostNeutralAmount = "entry_cost_neutral_amount"

// NewCardInstance creates a new card instance
func NewCardInstance(card *model.Card, ownerID int, turn int) *CardInstance {
	return newCardInstanceWithID(card, ownerID, turn, generateID())
}

func newCardInstanceWithID(card *model.Card, ownerID int, turn int, id string) *CardInstance {
	ci := &CardInstance{
		InstanceID:        id,
		Card:              card,
		OwnerID:           ownerID,
		CurrentLife:       card.Life,
		CurrentAttack:     card.Attack,
		IsHorizontal:      true, // enters horizontal (横置) by default
		Statuses:          make(map[string]int),
		ElementsGainBonus: make(map[string]int),
		EnterTurn:         turn,
	}
	return ci
}

// CanConsume checks if this card can be consumed (横置 to gain elements)
func (ci *CardInstance) CanConsume() bool {
	if ci.IsHorizontal {
		return false // already horizontal
	}
	if ci.Statuses[StatusStun] > 0 {
		return false // stunned
	}
	if ci.Statuses[StatusCooldown] > 0 {
		return false
	}
	return true
}

// Reset resets the card to vertical (竖置) state at turn start
func (ci *CardInstance) Reset() {
	if ci.Statuses[StatusFreeze] > 0 {
		return // frozen cards don't reset
	}
	if ci.Statuses[StatusCooldown] > 0 {
		return
	}
	ci.IsHorizontal = false
}

// Status constants
const (
	StatusBurn     = "点燃"
	StatusFreeze   = "冻结"
	StatusStun     = "眩晕"
	StatusPetrify  = "石化"
	StatusWeaken   = "虚弱"
	StatusCooldown = "冷却"
	StatusSeal     = "封印"
	StatusMastery  = "精通"
	StatusStealth  = "隐蔽"
)

const (
	BaseSkillSlots     = 5
	MaxSkillSlots      = 8
	BaseEquipmentSlots = 5
	MaxEquipmentSlots  = 8
)

// PlayerState holds all state for one player
type PlayerState struct {
	PlayerID   int                              `json:"player_id"` // 0 or 1
	PlayerName string                           `json:"player_name"`
	Hero       *CardInstance                    `json:"hero"`
	Units      [3][3]*CardInstance              `json:"units"`   // [col][row], hero is at [1][1]
	Terrain    [3][3]*CardInstance              `json:"terrain"` // Terrain cards placed on grid
	Skills     [MaxSkillSlots]*CardInstance     `json:"skills"`
	Equipment  [MaxEquipmentSlots]*CardInstance `json:"equipment"`
	Hand       []*CardInstance                  `json:"hand"`
	Deck       []*CardInstance                  `json:"deck"` // remaining draw pile
	SkillPool  []*CardInstance                  `json:"skill_pool"`
	Graveyard  []*CardInstance                  `json:"graveyard"`
	Exile      []*CardInstance                  `json:"exile,omitempty"`
	ExtraDeck  []*CardInstance                  `json:"extra_deck"`

	// Element pool - available elements this turn
	Elements                    map[string]int            `json:"elements"`
	StrictArcane                int                       `json:"strict_arcane,omitempty"`
	Shield                      int                       `json:"shield,omitempty"`
	ShieldBrokenThisTurn        bool                      `json:"shield_broken_this_turn,omitempty"`
	CannotGainShield            bool                      `json:"cannot_gain_shield,omitempty"`
	NextCompanionStealth        int                       `json:"next_companion_stealth,omitempty"`
	NextRedMoonDuration         int                       `json:"next_red_moon_duration,omitempty"`
	NextRedMoonCooldown         int                       `json:"next_red_moon_cooldown,omitempty"`
	SkipNextDraw                bool                      `json:"skip_next_draw"`
	TempModifiers               []TemporaryModifier       `json:"temp_modifiers"`
	SpellsCastThisTurn          map[string]int            `json:"spells_cast_this_turn,omitempty"`
	SpellsCastByNumberThisTurn  map[string]int            `json:"spells_cast_by_number_this_turn,omitempty"`
	LastLowCostWaterSpell       *CardInstance             `json:"last_low_cost_water_spell,omitempty"`
	SpellHitsThisTurn           int                       `json:"spell_hits_this_turn,omitempty"`
	SpellHitsLastTurn           int                       `json:"spell_hits_last_turn,omitempty"`
	SpellHitTargetsThisTurn     int                       `json:"spell_hit_targets_this_turn,omitempty"`
	SpellHitTargetsLastTurn     int                       `json:"spell_hit_targets_last_turn,omitempty"`
	SpellDamageThisTurn         int                       `json:"spell_damage_this_turn,omitempty"`
	SpellDamageLastTurn         int                       `json:"spell_damage_last_turn,omitempty"`
	DiscardAtTurnEnd            map[string]bool           `json:"discard_at_turn_end,omitempty"`
	LoadGainAtTurnEnd           map[string]map[string]int `json:"load_gain_at_turn_end,omitempty"`
	RevealedHand                map[string]bool           `json:"revealed_hand,omitempty"`
	DrawnTurn                   map[string]int            `json:"drawn_turn,omitempty"`
	DrawCountThisTurn           int                       `json:"draw_count_this_turn,omitempty"`
	DiscardedHandCountThisTurn  int                       `json:"discarded_hand_count_this_turn,omitempty"`
	FriendlyUnitDamagedThisTurn bool                      `json:"friendly_unit_damaged_this_turn,omitempty"`
	FriendlyUnitDamagedLastTurn bool                      `json:"friendly_unit_damaged_last_turn,omitempty"`
	HeroDamageTakenThisTurn     int                       `json:"hero_damage_taken_this_turn,omitempty"`
	HeroDamageTakenLastTurn     int                       `json:"hero_damage_taken_last_turn,omitempty"`

	// Legacy charge pool. Do not use this for 精通; mastery is a per-card
	// instance mark stored in CardInstance.Statuses[StatusMastery].
	Charge int `json:"charge"`

	// Deck definition
	DeckDef *model.Deck `json:"-"`
}

// NewPlayerState creates initial player state from a deck
func NewPlayerState(id int, name string, deck *model.Deck) *PlayerState {
	ps := &PlayerState{
		PlayerID:                   id,
		PlayerName:                 name,
		Elements:                   make(map[string]int),
		SpellsCastThisTurn:         make(map[string]int),
		SpellsCastByNumberThisTurn: make(map[string]int),
		DiscardAtTurnEnd:           make(map[string]bool),
		LoadGainAtTurnEnd:          make(map[string]map[string]int),
		RevealedHand:               make(map[string]bool),
		DrawnTurn:                  make(map[string]int),
		DeckDef:                    deck,
	}

	// Initialize elements to 0
	for _, e := range model.AllElements {
		ps.Elements[e] = 0
	}

	return ps
}

// InitCards creates all card instances from the deck definition
func (e *Engine) initPlayerCards(ps *PlayerState, turn int) {
	cardDB := make(map[string]*model.Card)
	for k, v := range getCardDB() {
		cardDB[k] = v
	}

	// Create hero
	heroCard := cardDB[ps.DeckDef.HeroID]
	ps.Hero = e.newCardInstance(heroCard, ps.PlayerID, turn)
	ps.Hero.IsHorizontal = false // hero starts vertical
	ps.Hero.Position = &Position{Col: 1, Row: 1}
	ps.Units[1][1] = ps.Hero

	// Create main deck cards and shuffle
	ps.Deck = make([]*CardInstance, 0, len(ps.DeckDef.MainDeck))
	for _, id := range ps.DeckDef.MainDeck {
		card := cardDB[id]
		ci := e.newCardInstance(card, ps.PlayerID, 0)
		ps.Deck = append(ps.Deck, ci)
	}
	e.shuffleCards(ps.Deck)

	// Create skill pool cards
	ps.SkillPool = make([]*CardInstance, 0, len(ps.DeckDef.SkillPool))
	for _, id := range ps.DeckDef.SkillPool {
		card := cardDB[id]
		ci := e.newCardInstance(card, ps.PlayerID, 0)
		ps.SkillPool = append(ps.SkillPool, ci)
	}

	// Create extra deck cards
	ps.ExtraDeck = make([]*CardInstance, 0, len(ps.DeckDef.ExtraDeck))
	for _, id := range ps.DeckDef.ExtraDeck {
		card := cardDB[id]
		ci := e.newCardInstance(card, ps.PlayerID, 0)
		ps.ExtraDeck = append(ps.ExtraDeck, ci)
	}

	ps.Hand = make([]*CardInstance, 0)
	ps.Graveyard = make([]*CardInstance, 0)
}

// DrawCards draws n cards from the deck to hand
func (ps *PlayerState) DrawCards(n int) []*CardInstance {
	drawn := make([]*CardInstance, 0, n)
	for i := 0; i < n && len(ps.Deck) > 0; i++ {
		card := ps.Deck[0]
		ps.Deck = ps.Deck[1:]
		ps.Hand = append(ps.Hand, card)
		if cardRevealsOnDraw(card) {
			ps.RevealedHand[card.InstanceID] = true
		}
		drawn = append(drawn, card)
	}
	return drawn
}

// FindHandCard finds a card in hand by instance ID
func (ps *PlayerState) FindHandCard(instanceID string) (*CardInstance, int) {
	for i, c := range ps.Hand {
		if c.InstanceID == instanceID {
			return c, i
		}
	}
	return nil, -1
}

// RemoveFromHand removes a card from hand by index
func (ps *PlayerState) RemoveFromHand(index int) {
	if index >= 0 && index < len(ps.Hand) && ps.RevealedHand != nil {
		delete(ps.RevealedHand, ps.Hand[index].InstanceID)
	}
	ps.Hand = append(ps.Hand[:index], ps.Hand[index+1:]...)
}

// CountUnits counts non-nil units on the grid (including hero)
func (ps *PlayerState) CountUnits() int {
	count := 0
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if ps.Units[col][row] != nil {
				count++
			}
		}
	}
	return count
}

// FindEmptyPosition finds an empty position on the unit grid
func (ps *PlayerState) FindEmptyPosition() *Position {
	// Prefer positions near hero: adjacent first, then corners
	order := []Position{
		{0, 1}, {2, 1}, {1, 0}, {1, 2},
		{0, 0}, {2, 0}, {0, 2}, {2, 2},
	}
	for _, p := range order {
		if ps.Units[p.Col][p.Row] == nil {
			return &p
		}
	}
	return nil
}

// GetFrontRow returns the front row index (closest row with units, from enemy's perspective)
// Row 0 = front (closest to enemy)
func (ps *PlayerState) GetFrontRow() int {
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if ps.Units[col][row] != nil {
				return row
			}
		}
	}
	return -1 // no units
}

// CanPayCost checks if the player can pay the given element cost
// Arcane (无) in cost can be paid with any element, and arcane elements can pay for any cost
func (ps *PlayerState) CanPayCost(cost map[string]int) bool {
	_, ok := calculateElementPayment(ps.Elements, cost)
	return ok
}

// PayCost deducts elements from the pool. Returns false if insufficient.
func (ps *PlayerState) PayCost(cost map[string]int) bool {
	payment, ok := calculateElementPayment(ps.Elements, cost)
	if !ok {
		return false
	}

	return ps.PayCostWithPayment(cost, payment)
}

func (ps *PlayerState) PayCostWithPayment(cost map[string]int, payment map[string]int) bool {
	if !validateElementPayment(ps.Elements, cost, payment) {
		return false
	}
	for elem, amount := range payment {
		ps.Elements[elem] -= amount
	}
	return true
}

// GainElements adds elements to the pool
func (ps *PlayerState) GainElements(gains map[string]int) {
	for elem, amount := range gains {
		ps.Elements[elem] += amount
	}
}

// GamePhase represents the current phase of the game
type GamePhase int

const (
	PhaseWaitingPlayers GamePhase = iota
	PhaseMulligan
	PhaseTurnStart
	PhaseMain
	PhaseSpellCast // spell has been cast, waiting for defense
	PhaseDefenseWindow
	PhaseWaitingAction // waiting for a player to resolve a pending action
	PhaseTurnEnd
	PhaseGameOver
)

func (p GamePhase) String() string {
	switch p {
	case PhaseWaitingPlayers:
		return "waiting_players"
	case PhaseMulligan:
		return "mulligan"
	case PhaseTurnStart:
		return "turn_start"
	case PhaseMain:
		return "main"
	case PhaseSpellCast:
		return "spell_cast"
	case PhaseDefenseWindow:
		return "defense_window"
	case PhaseWaitingAction:
		return "waiting_action"
	case PhaseTurnEnd:
		return "turn_end"
	case PhaseGameOver:
		return "game_over"
	default:
		return "unknown"
	}
}

// PendingAction represents a player choice that must be resolved
type PendingAction struct {
	Type         string                                             `json:"type"`       // "select_target", "select_card", "discard", "select_position"
	PlayerID     int                                                `json:"player_id"`  // which player must choose
	Prompt       string                                             `json:"prompt"`     // display text
	Candidates   []map[string]any                                   `json:"candidates"` // selectable options (cards or positions)
	MinSelect    int                                                `json:"min_select"` // minimum selections required
	MaxSelect    int                                                `json:"max_select"` // maximum selections allowed
	Context      map[string]any                                     `json:"context,omitempty"`
	Cost         map[string]int                                     `json:"cost,omitempty"`
	CanOverexert bool                                               `json:"can_overexert,omitempty"`
	Callback     func(selected []string)                            `json:"-"` // called when resolved
	CallbackData func(selected []string, data map[string]any)       `json:"-"`
	CallbackErr  func(selected []string, data map[string]any) error `json:"-"`
	Available    func() bool                                        `json:"-"` // checked before a queued action is shown
	resolutions  []*resolutionFrame
}

// GameState holds the entire game state
type GameState struct {
	GameID                       string          `json:"game_id"`
	Players                      [2]*PlayerState `json:"players"`
	CurrentTurn                  int             `json:"current_turn"` // 0 or 1 (which player's turn)
	FirstPlayer                  int             `json:"first_player"` // 0 or 1
	TurnNumber                   int             `json:"turn_number"`
	Phase                        GamePhase       `json:"phase"`
	Winner                       int             `json:"winner"` // -2 = draw, -1 = no winner, 0 or 1
	HandLimit                    int             `json:"hand_limit"`
	IsFirstTurn                  bool            `json:"is_first_turn"` // first player's first turn
	CardEnteredGraveyardThisTurn bool            `json:"card_entered_graveyard_this_turn,omitempty"`
	DrawOfferBy                  int             `json:"draw_offer_by"` // -1 when there is no active draw offer

	// Combat state
	PendingSpell *SpellCast `json:"pending_spell,omitempty"`

	// Player choice state
	PendingAction *PendingAction `json:"pending_action,omitempty"`

	// Mulligan state
	MulliganDone [2]bool `json:"mulligan_done"`

	// Phase to resume after pending action resolves
	ResumePhase        GamePhase        `json:"-"`
	PendingActionQueue []*PendingAction `json:"-"`
}

// SpellCast represents an ongoing spell combat
type SpellCast struct {
	AttackerID   int                `json:"attacker_id"`
	Skill        *CardInstance      `json:"skill"`
	Target       SpellTarget        `json:"target"`
	TotalPower   int                `json:"total_power"`
	PowerSources []SpellPowerSource `json:"power_sources,omitempty"`
	BoostSkills  []*CardInstance    `json:"boost_skills"` // skills used to boost
	ExtraTargets []SpellTarget      `json:"extra_targets,omitempty"`
	resolutions  []*resolutionFrame
}

type SpellPowerSource struct {
	InstanceID string `json:"instance_id"`
	CardName   string `json:"card_name"`
	Power      int    `json:"power"`
	IsMain     bool   `json:"is_main"`
}

// SpellTarget represents the target of a spell
type SpellTarget struct {
	Type     string   `json:"type"`     // "unit" or "area"
	Position Position `json:"position"` // for unit targets
	OwnerID  *int     `json:"owner_id,omitempty"`
}

// NewGameState creates a new game state
func NewGameState(gameID string) *GameState {
	return &GameState{
		GameID:      gameID,
		Winner:      -1,
		HandLimit:   5,
		Phase:       PhaseWaitingPlayers,
		DrawOfferBy: -1,
	}
}

var idCounter atomic.Uint64

func generateID() string {
	return fmt.Sprintf("ci_fixture_%d", idCounter.Add(1))
}

func getCardDB() map[string]*model.Card {
	// This will be set by the engine
	return cardDBRef
}

var cardDBRef map[string]*model.Card

func SetCardDB(db map[string]*model.Card) {
	cardDBRef = db
}
