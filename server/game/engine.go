package game

import (
	"eraofarcane/cards"
	"eraofarcane/model"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"sync"
)

// GameEvent represents something that happened in the game
type GameEvent struct {
	Type   string         `json:"type"`
	Data   map[string]any `json:"data"`
	Player int            `json:"player"` // which player this is relevant to, -1 for both
}

// ActionMessage is a player action received via WebSocket
type ActionMessage struct {
	Action string         `json:"action"`
	Data   map[string]any `json:"data"`
}

// EventCallback is called when events occur (to send to clients)
type EventCallback func(event GameEvent, targetPlayer int) // targetPlayer: 0, 1, or -1 for both

// pendingDeath records a unit that should die after the current resolution scope.
type pendingDeath struct {
	unit      *CardInstance
	ownerID   int
	deathData map[string]any
}

type pendingDeathTrigger struct {
	unit      *CardInstance
	ownerID   int
	deathData map[string]any
}

const (
	DeathCauseSacrifice = "sacrifice"
	DeathCauseDevour    = "devour"
)

// Engine manages a single game instance
type Engine struct {
	State    *GameState
	mu       sync.Mutex
	callback EventCallback
	log      []GameEvent

	resolutionDepth int
	resolvingDeaths bool
	deathQueue      []pendingDeath
	deathTriggers   []pendingDeathTrigger
}

// NewEngine creates a new game engine
func NewEngine(gameID string, callback EventCallback) *Engine {
	return &Engine{
		State:    NewGameState(gameID),
		callback: callback,
		log:      make([]GameEvent, 0),
	}
}

// emit sends an event to clients and records it
func (e *Engine) emit(event GameEvent) {
	e.log = append(e.log, event)
	if e.callback != nil {
		e.callback(event, event.Player)
	}
}

// SetupGame initializes both players and starts the game
func (e *Engine) SetupGame(p1Name string, p1Deck *model.Deck, p2Name string, p2Deck *model.Deck) error {
	return e.SetupGameWithFirstPlayer(p1Name, p1Deck, p2Name, p2Deck, rand.Intn(2))
}

// SetupGameWithFirstPlayer initializes both players with an explicit first player.
func (e *Engine) SetupGameWithFirstPlayer(p1Name string, p1Deck *model.Deck, p2Name string, p2Deck *model.Deck, firstPlayer int) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if firstPlayer < 0 || firstPlayer > 1 {
		firstPlayer = 0
	}

	// Create player states
	e.State.Players[0] = NewPlayerState(0, p1Name, p1Deck)
	e.State.Players[1] = NewPlayerState(1, p2Name, p2Deck)
	e.State.FirstPlayer = firstPlayer

	// Initialize cards
	e.State.Players[0].InitCards(0)
	e.State.Players[1].InitCards(0)
	e.triggerInitialHeroEnterEffects()
	e.triggerGameStartEffects()
	e.emit(GameEvent{
		Type:   "game_setup",
		Player: -1,
		Data: map[string]any{
			"first_player": e.State.FirstPlayer,
			"timing":       "before_initial_draw",
		},
	})

	// Draw initial hands (4 cards each; Raven starts with one extra card)
	for i := 0; i < 2; i++ {
		drawn := e.drawCards(i, e.initialHandSizeForPlayer(e.State.Players[i]))
		e.emit(GameEvent{
			Type:   "initial_draw",
			Player: i,
			Data: map[string]any{
				"cards": cardsToInfo(drawn),
				"count": len(drawn),
			},
		})
	}

	// Enter mulligan phase
	e.State.Phase = PhaseMulligan
	e.emit(GameEvent{
		Type:   "phase_change",
		Player: -1,
		Data:   map[string]any{"phase": "mulligan"},
	})

	return nil
}

func (e *Engine) triggerGameStartEffects() {
	for playerID := 0; playerID < 2; playerID++ {
		data := map[string]any{"initial_setup": true}
		for _, card := range e.getAllFieldCards(e.State.Players[playerID]) {
			e.triggerEffects(TriggerOnGameStart, card, nil, data)
		}
	}
}
func (e *Engine) triggerInitialHeroEnterEffects() {
	for playerID := 0; playerID < 2; playerID++ {
		hero := e.State.Players[playerID].Hero
		if hero == nil {
			continue
		}
		data := map[string]any{"initial_setup": true, "entered_player": playerID}
		e.triggerEffects(TriggerOnEnter, hero, nil, data)
		e.triggerFieldEffectsWithData(TriggerOnUnitEnter, playerID, hero, data)
		e.triggerFieldEffectsWithData(TriggerOnUnitEnter, 1-playerID, hero, data)
	}
}

// HandleAction processes a player action
func (e *Engine) HandleAction(playerID int, action ActionMessage) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	log.Printf("[Game %s] Player %d action: %s", e.State.GameID, playerID, action.Action)

	e.beginResolution()
	defer e.enforceAllSlotCapacities()
	defer e.endResolution()

	switch action.Action {
	case "mulligan":
		return e.handleMulligan(playerID, action)
	case "summon":
		return e.handleSummon(playerID, action)
	case "consume":
		return e.handleConsume(playerID, action)
	case "cast_spell":
		return e.handleCastSpell(playerID, action)
	case "react_spell":
		return e.handleReactSpell(playerID, action)
	case "defend":
		return e.handleDefend(playerID, action)
	case "no_defend":
		return e.handleNoDefend(playerID, action)
	case "attack":
		return e.handleAttack(playerID, action)
	case "equip":
		return e.handleEquip(playerID, action)
	case "learn_skill":
		return e.handleLearnSkill(playerID, action)
	case "use_item":
		return e.handleUseItem(playerID, action)
	case "place_terrain":
		return e.handlePlaceTerrain(playerID, action)
	case "use_ability":
		return e.handleUseAbility(playerID, action)
	case "resolve_action":
		return e.handleResolveAction(playerID, action)
	case "end_turn":
		return e.handleEndTurn(playerID, action)
	default:
		return fmt.Errorf("unknown action: %s", action.Action)
	}
}

func (e *Engine) enforceAllSlotCapacities() {
	for _, ps := range e.State.Players {
		e.enforceSlotCapacities(ps)
	}
}

func (e *Engine) enforceSlotCapacities(ps *PlayerState) {
	if ps == nil {
		return
	}

	equipmentCap := equipmentSlotCapacity(ps)
	for i := equipmentCap; i < len(ps.Equipment); i++ {
		equipment := ps.Equipment[i]
		if equipment == nil {
			continue
		}
		ps.Equipment[i] = nil
		equipment.SlotIndex = -1
		equipment.BoundSkills = nil
		ps.Graveyard = append(ps.Graveyard, equipment)
		e.emit(GameEvent{Type: "discard", Player: ps.PlayerID, Data: map[string]any{"card": cardToInfo(equipment)}})
	}

	skillCap := skillSlotCapacity(ps)
	for i := skillCap; i < len(ps.Skills); i++ {
		skill := ps.Skills[i]
		if skill == nil {
			continue
		}
		ps.Skills[i] = nil
		returnSkillToPool(skill)
		ps.SkillPool = append(ps.SkillPool, skill)
		e.emit(GameEvent{Type: "skill_returned_to_pool", Player: ps.PlayerID, Data: map[string]any{"card": cardToInfo(skill)}})
	}
}

func (e *Engine) handleReactSpell(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseDefenseWindow {
		return fmt.Errorf("not in spell reaction window")
	}
	if e.State.PendingSpell == nil {
		return fmt.Errorf("no pending spell")
	}
	if playerID == e.State.PendingSpell.AttackerID {
		return fmt.Errorf("attacker cannot react to their own spell this way")
	}

	instanceID, _ := action.Data["instance_id"].(string)
	ps := e.State.Players[playerID]
	skill := e.findReactionCard(ps, instanceID)
	if skill == nil {
		return fmt.Errorf("reaction skill not found")
	}
	if err := e.validateSkillForPurpose(skill, skillPurposeReaction); err != nil {
		return err
	}
	cost := map[string]int{}
	if skill.Card.IsSkill() {
		cost = e.effectiveSkillUseCost(ps, skill)
	}
	overexertIDsRaw, _ := action.Data["overexert_ids"].([]any)
	overexertIDs := stringsFromAnySlice(overexertIDsRaw)
	overexertUnits, err := e.collectOverexertUnits(ps, overexertIDs)
	if err != nil {
		return err
	}
	if !e.canPayCostWithOverexertOptions(ps, cost, overexertUnits, e.playerHasLightWildcard(ps)) {
		return fmt.Errorf("not enough elements")
	}
	if !e.payDefenseCostWithOptions(ps, cost, action, overexertUnits, e.playerHasLightWildcard(ps)) {
		return fmt.Errorf("invalid payment")
	}
	e.destroyFuyeDoomedAfterExert(overexertUnits)

	if skill.Card.IsSkill() {
		skill.IsHorizontal = true
		if !e.shouldSkipCooldown(ps, skill) {
			e.ApplyKeywordOnSkillUse(skill)
		}
		e.applySkillUseCooldownModifiers(ps, skill)
	}
	if skill.Card.IsSkill() {
		e.consumeNextSkillUseModifiers(ps, skill)
	}
	e.advanceMasteryForUsedSkills(playerID, skill)

	behavior := behaviorForNumber(skill.Card.Number).(SpellReactionBehavior)
	if !behavior.HasActiveSpellReaction(skill) {
		return fmt.Errorf("skill cannot react to spells")
	}
	ctx := &EffectContext{
		Engine:     e,
		Source:     skill,
		PlayerID:   playerID,
		OpponentID: 1 - playerID,
		ExtraData: map[string]any{
			"react_player": playerID,
			"spell":        e.State.PendingSpell,
		},
	}
	return behavior.OnSpellReaction(ctx, e.State.PendingSpell)
}

// handleMulligan handles the mulligan (redraw) action
func (e *Engine) handleMulligan(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseMulligan {
		return fmt.Errorf("not in mulligan phase")
	}
	if e.State.MulliganDone[playerID] {
		return fmt.Errorf("already submitted mulligan")
	}

	// Check if player wants to redraw
	keepStr, _ := action.Data["keep"].(bool)

	if !keepStr {
		// Redraw to the same starting hand size. Raven starts with one extra card,
		// so its mulligan should also redraw that extra card.
		ps := e.State.Players[playerID]
		ps.Deck = append(ps.Deck, ps.Hand...)
		ps.Hand = make([]*CardInstance, 0)
		shuffleDeck(ps.Deck)
		drawn := e.drawCards(playerID, e.initialHandSizeForPlayer(ps))
		e.emit(GameEvent{
			Type:   "mulligan_redraw",
			Player: playerID,
			Data: map[string]any{
				"cards": cardsToInfo(drawn),
				"count": len(drawn),
			},
		})
	}

	e.State.MulliganDone[playerID] = true

	e.emit(GameEvent{
		Type:   "mulligan_done",
		Player: playerID,
		Data:   map[string]any{"player": playerID},
	})

	// If both players are done, start the game
	if e.State.MulliganDone[0] && e.State.MulliganDone[1] {
		e.startGame()
	}

	return nil
}

// startGame begins the actual game
func (e *Engine) startGame() {
	e.State.CurrentTurn = e.State.FirstPlayer
	e.State.TurnNumber = 1
	e.State.IsFirstTurn = true

	e.emit(GameEvent{
		Type:   "game_start",
		Player: -1,
		Data: map[string]any{
			"first_player": e.State.FirstPlayer,
		},
	})

	e.startTurn()
}

// startTurn begins a new turn for the current player
func (e *Engine) startTurn() {
	e.clearDamageTakenThisTurn()

	ps := e.State.Players[e.State.CurrentTurn]
	ps.SpellsCastThisTurn = make(map[string]int)
	ps.DrawCountThisTurn = 0

	// Elements are cleared at the end of their owner's turn. Start turn should
	// not be the rule point that removes remaining elements.
	e.applyTurnStartTemporaryModifiers(ps)

	e.continuePreDrawTurnStartEffects(ps, append([]*CardInstance(nil), e.getAllFieldCards(ps)...), 0)
}

func (e *Engine) continueStartTurnAfterPreDraw(ps *PlayerState) {
	if ps == nil {
		return
	}
	if ps.SkipNextDraw {
		ps.SkipNextDraw = false
		e.emit(GameEvent{
			Type:   "effect_trigger",
			Player: ps.PlayerID,
			Data: map[string]any{
				"effect": "skip_draw",
			},
		})
	} else {
		drawn := e.drawCards(ps.PlayerID, 1)
		if len(drawn) > 0 {
			// Notify opponent about the draw (without card info)
			e.emit(GameEvent{
				Type:   "opponent_draw",
				Player: 1 - ps.PlayerID,
				Data: map[string]any{
					"count": 1,
				},
			})
		}
		if e.State.PendingAction != nil {
			e.wrapPendingActionContinuation(func() {
				e.continueStartTurnAfterDraw(ps)
			})
			return
		}
	}

	e.continueStartTurnAfterDraw(ps)
}

func (e *Engine) continueStartTurnAfterDraw(ps *PlayerState) {
	if ps == nil {
		return
	}
	e.State.Phase = PhaseMain

	e.emit(GameEvent{
		Type:   "turn_start",
		Player: -1,
		Data: map[string]any{
			"current_player": ps.PlayerID,
			"turn_number":    e.State.TurnNumber,
			"elements":       ps.Elements,
		},
	})

	// Trigger 回合开始 effects for all cards on the current player's field
	allCards := e.getAllFieldCards(ps)
	for _, card := range allCards {
		e.triggerEffects(TriggerOnTurnStart, card, nil, nil)
	}
	if e.State.PendingAction == nil {
		e.triggerPrayerAbilities(ps.PlayerID)
	}
}

func (e *Engine) continuePreDrawTurnStartEffects(ps *PlayerState, cards []*CardInstance, start int) {
	if ps == nil {
		return
	}
	for i := start; i < len(cards); i++ {
		card := cards[i]
		if card == nil || card.Card == nil || !isPreDrawTurnStartCard(card.Card.Number) || e.hasEffectiveStatus(card, StatusPetrify) {
			continue
		}
		if !e.cardStillOnField(card) {
			continue
		}
		behavior, ok := globalRegistry.GetBehavior(card.Card.Number).(OnTurnStartBehavior)
		if !ok {
			continue
		}
		ctx := &EffectContext{
			Engine:     e,
			Source:     card,
			PlayerID:   ps.PlayerID,
			OpponentID: 1 - ps.PlayerID,
			ExtraData:  map[string]any{"timing": "pre_draw"},
		}
		_ = behavior.OnTurnStart(ctx)
		if e.State.PendingAction != nil {
			next := i + 1
			e.wrapPendingActionContinuation(func() {
				e.continuePreDrawTurnStartEffects(ps, cards, next)
			})
			return
		}
	}
	e.continueStartTurnAfterPreDraw(ps)
}

func isPreDrawTurnStartCard(number string) bool {
	switch number {
	case "1021008", "4411001":
		return true
	default:
		return false
	}
}

func (e *Engine) triggerPrayerAbilities(playerID int) {
	ps := e.State.Players[playerID]
	if ps == nil {
		return
	}
	e.continuePrayerAbilities(playerID, append([]*CardInstance(nil), e.getAllFieldCards(ps)...), 0)
}

func (e *Engine) continuePrayerAbilities(playerID int, cards []*CardInstance, start int) {
	for i := start; i < len(cards); i++ {
		card := cards[i]
		if !cardHasActivePrayer(card) || !e.cardStillOnField(card) {
			continue
		}
		if e.executePrayerAbility(playerID, card); e.State.PendingAction != nil {
			next := i + 1
			e.wrapPendingActionContinuation(func() {
				e.continuePrayerAbilities(playerID, cards, next)
			})
			return
		}
	}
}

func (e *Engine) executePrayerAbility(playerID int, card *CardInstance) {
	behavior := cardBehavior(card)
	perTurn, ok := behavior.(PerTurnAbility)
	if !ok || !perTurn.HasActivePerTurn(card) {
		return
	}
	ctx := &EffectContext{
		Engine:     e,
		Source:     card,
		PlayerID:   playerID,
		OpponentID: 1 - playerID,
		ExtraData:  map[string]any{"prayer": true},
	}
	run := func() {
		e.emit(GameEvent{
			Type:   "effect_trigger",
			Player: -1,
			Data: map[string]any{
				"effect": "prayer",
				"card":   cardToInfo(card),
				"player": playerID,
			},
		})
		if err := perTurn.OnPerTurn(ctx); err != nil {
			e.emit(GameEvent{
				Type:   "effect_trigger",
				Player: playerID,
				Data: map[string]any{
					"effect": "prayer_error",
					"card":   cardToInfo(card),
					"error":  err.Error(),
				},
			})
		}
	}
	if optional, ok := behavior.(OptionalPrayerAbility); ok && optional.IsPrayerOptional(card) {
		e.SetPendingAction(playerID, "optional_prayer", "是否发动祈咒: "+card.Card.Name, []map[string]any{candidateInfo(card, "card", "own")}, 0, 1, func(selected []string) {
			if len(selected) > 0 {
				run()
			}
		})
		return
	}
	run()
}

func (e *Engine) wrapPendingActionContinuation(continueFn func()) {
	pa := e.State.PendingAction
	if pa == nil || continueFn == nil {
		return
	}
	after := func() {
		if e.State.PendingAction != nil {
			e.wrapPendingActionContinuation(continueFn)
			return
		}
		continueFn()
	}
	if pa.CallbackErr != nil {
		original := pa.CallbackErr
		pa.CallbackErr = func(selected []string, data map[string]any) error {
			if err := original(selected, data); err != nil {
				return err
			}
			after()
			return nil
		}
		return
	}
	if pa.CallbackData != nil {
		original := pa.CallbackData
		pa.CallbackData = func(selected []string, data map[string]any) {
			original(selected, data)
			after()
		}
		return
	}
	original := pa.Callback
	pa.Callback = func(selected []string) {
		if original != nil {
			original(selected)
		}
		after()
	}
}

func (e *Engine) cardStillOnField(card *CardInstance) bool {
	if card == nil || card.OwnerID < 0 || card.OwnerID >= len(e.State.Players) {
		return false
	}
	for _, fieldCard := range e.getAllFieldCards(e.State.Players[card.OwnerID]) {
		if fieldCard == card {
			return true
		}
	}
	return false
}

func (e *Engine) clearDamageTakenThisTurn() {
	for _, ps := range e.State.Players {
		if ps == nil {
			continue
		}
		for _, card := range e.getAllFieldCards(ps) {
			if card != nil {
				card.DamageTakenThisTurn = 0
			}
		}
		for _, card := range ps.Graveyard {
			if card != nil {
				card.DamageTakenThisTurn = 0
			}
		}
	}
}

func (e *Engine) drawCards(playerID int, n int) []*CardInstance {
	if n <= 0 {
		return nil
	}
	ps := e.State.Players[playerID]
	drawn := ps.DrawCards(n)
	for _, card := range drawn {
		e.notifyCardDrawn(playerID, card)
	}
	return drawn
}

func (e *Engine) millTopDeckCards(playerID int, n int) []*CardInstance {
	if n <= 0 {
		return nil
	}
	ps := e.State.Players[playerID]
	if ps == nil {
		return nil
	}
	count := min(n, len(ps.Deck))
	milled := make([]*CardInstance, 0, count)
	for i := 0; i < count; i++ {
		card := ps.Deck[0]
		ps.Deck = ps.Deck[1:]
		ps.Graveyard = append(ps.Graveyard, card)
		milled = append(milled, card)
		e.emit(GameEvent{Type: "discard", Player: playerID, Data: map[string]any{"card": cardToInfo(card)}})
	}
	return milled
}

func (e *Engine) notifyCardDrawn(playerID int, card *CardInstance) {
	if card == nil {
		return
	}
	ps := e.State.Players[playerID]
	if ps.DrawnTurn == nil {
		ps.DrawnTurn = make(map[string]int)
	}
	if cardRevealsOnDraw(card) {
		ps.RevealedHand[card.InstanceID] = true
	}
	ps.DrawnTurn[card.InstanceID] = e.State.TurnNumber
	ps.DrawCountThisTurn++
	e.emit(GameEvent{
		Type:   "draw_card",
		Player: playerID,
		Data:   map[string]any{"card": cardToInfo(card)},
	})
	data := map[string]any{
		"drawn_card":           card,
		"drawn_player":         playerID,
		"draw_count_this_turn": ps.DrawCountThisTurn,
		"initial_hand":         e.State.Phase == PhaseWaitingPlayers || e.State.Phase == PhaseMulligan,
	}
	e.triggerFieldEffectsWithData(TriggerOnDraw, playerID, card, data)
	if h, ok := cardBehavior(card).(OnSelfDrawBehavior); ok && h.HasActiveDraw(card) {
		_ = h.OnSelfDraw(&EffectContext{
			Engine:     e,
			Source:     card,
			PlayerID:   playerID,
			OpponentID: 1 - playerID,
			ExtraData:  data,
		})
	}
}

// resetCards resets all cards to vertical state
func (e *Engine) resetCards(ps *PlayerState) {
	// Reset units
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if ps.Units[col][row] != nil {
				e.resetCard(ps.Units[col][row])
				ps.Units[col][row].UsedThisTurn = 0
				for _, skill := range ps.Units[col][row].BoundSkills {
					e.resetCard(skill)
					skill.UsedThisTurn = 0
				}
			}
		}
	}
	// Reset skills
	for i := range ps.Skills {
		if ps.Skills[i] != nil {
			e.resetCard(ps.Skills[i])
			ps.Skills[i].UsedThisTurn = 0
		}
	}
	// Reset equipment
	for i := range ps.Equipment {
		if ps.Equipment[i] != nil {
			e.resetCard(ps.Equipment[i])
			ps.Equipment[i].UsedThisTurn = 0
			for _, skill := range ps.Equipment[i].BoundSkills {
				e.resetCard(skill)
				skill.UsedThisTurn = 0
			}
		}
	}
}

// refreshElements calculates available elements from all vertical (竖置) cards
func (e *Engine) refreshElements(ps *PlayerState) {
	// Clear elements
	for elem := range ps.Elements {
		ps.Elements[elem] = 0
	}
	// Don't auto-gain. Elements are gained by consuming (横置) cards.
	// At turn start, elements reset to 0. Player must consume cards to gain elements.
}

// handleSummon handles summoning a companion to the field
func (e *Engine) handleSummon(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseMain {
		return fmt.Errorf("not in main phase")
	}
	if e.State.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}

	instanceID, _ := action.Data["instance_id"].(string)
	colF, _ := action.Data["col"].(float64)
	rowF, _ := action.Data["row"].(float64)
	col := int(colF)
	row := int(rowF)

	ps := e.State.Players[playerID]

	// Find card in hand
	card, handIdx := ps.FindHandCard(instanceID)
	if card == nil {
		return fmt.Errorf("card not found in hand")
	}

	// Must be a companion or item
	if !card.Card.IsCompanion() {
		return fmt.Errorf("can only summon companions to unit area")
	}

	// Check position
	pos := Position{Col: col, Row: row}
	if !pos.Valid() {
		return fmt.Errorf("invalid position")
	}

	cost := e.effectiveCardPlayCost(ps, card)
	if !e.canPayCost(ps, cost) {
		return fmt.Errorf("not enough elements")
	}

	if !e.canPayCostForCardAction(ps, card, cost, cost, paymentPurposePlay, action) {
		return fmt.Errorf("invalid payment")
	}
	if err := e.validateAndApplySummonDevour(playerID, card, action); err != nil {
		return err
	}
	if ps.Units[col][row] != nil {
		return fmt.Errorf("position already occupied")
	}

	// Check unit limit (max 9 including hero) after devour costs are paid.
	if ps.CountUnits() >= 9 {
		return fmt.Errorf("unit area is full")
	}

	// Pay cost and place
	if !e.payCostForCardAction(ps, card, cost, cost, paymentPurposePlay, action) {
		return fmt.Errorf("invalid payment")
	}
	e.notifyCardPlayCostPaid(ps, card)
	ps.RemoveFromHand(handIdx)
	card.Position = &Position{Col: col, Row: row}
	card.IsHorizontal = true // enters horizontal by default
	card.EnterTurn = e.State.TurnNumber
	ps.Units[col][row] = card

	// Apply keyword effects (速攻 makes it enter vertical, etc.)
	e.ApplyKeywordOnEnter(card)
	e.ApplySummonModifiersOnEnter(card)

	e.emit(GameEvent{
		Type:   "summon",
		Player: -1,
		Data: map[string]any{
			"player":   playerID,
			"card":     cardToInfo(card),
			"position": pos,
			"elements": ps.Elements,
		},
	})

	// Trigger 入场 (on enter) effects for the summoned card
	e.triggerEffects(TriggerOnEnter, card, nil, nil)

	enterData := map[string]any{"entered_player": playerID}
	// Notify both sides about the new unit entering; individual card behaviors
	// decide whether they care about friendly or enemy units.
	e.triggerFieldEffectsWithData(TriggerOnUnitEnter, playerID, card, enterData)
	e.triggerFieldEffectsWithData(TriggerOnUnitEnter, 1-playerID, card, enterData)

	e.checkWinCondition()
	return nil
}

func (e *Engine) validateAndApplySummonDevour(playerID int, card *CardInstance, action ActionMessage) error {
	requirement := summonDevourRequirement(card)
	if len(requirement) == 0 {
		return nil
	}

	devourIDsRaw, _ := action.Data["devour_ids"].([]any)
	devourIDs := stringsFromAnySlice(devourIDsRaw)
	if legacyID, _ := action.Data["devour_id"].(string); legacyID != "" {
		devourIDs = append(devourIDs, legacyID)
	}
	if len(devourIDs) == 0 {
		return fmt.Errorf("%s requires devour before summon", card.Card.Name)
	}

	ps := e.State.Players[playerID]
	targets := make([]*CardInstance, 0, len(devourIDs))
	total := make(map[string]int)
	seen := make(map[string]bool, len(devourIDs))
	for _, devourID := range devourIDs {
		if seen[devourID] {
			return fmt.Errorf("duplicate devour target")
		}
		seen[devourID] = true
		target := e.findFieldCardByInstance(ps, devourID)
		if target == nil {
			target = e.findUnitOnGrid(ps, devourID)
		}
		if !isValidSummonDevourTarget(target, card) {
			return fmt.Errorf("invalid devour target")
		}
		if target.CurrentLife > 0 {
			total[DevourLife] += target.CurrentLife
		}
		for elem, amount := range e.effectiveElementsGain(target) {
			if amount > 0 {
				total[elem] += amount
			}
		}
		targets = append(targets, target)
	}

	for elem, amount := range requirement {
		if total[elem] < amount {
			return fmt.Errorf("devour targets load does not satisfy requirement")
		}
	}
	for _, target := range targets {
		if target.Card.IsCompanion() {
			e.destroyUnitWithCause(target, playerID, DeathCauseDevour)
			continue
		}
		e.discardFriendlyCandidate(playerID, target.InstanceID)
	}
	return nil
}

func isValidSummonDevourTarget(target *CardInstance, summoned *CardInstance) bool {
	if target == nil || target.Card == nil || target == summoned || target.Card.IsHero() {
		return false
	}
	return target.Card.IsCompanion() || isEquipmentItem(target)
}

// handleConsume handles consuming a card (横置 to gain elements)
func (e *Engine) elementsGainedFromConsume(playerID int, card *CardInstance, action ActionMessage) (map[string]int, error) {
	gains := e.effectiveElementsGain(card)
	if !e.isFirstPlayerFirstTurnHeroConsume(playerID, card) {
		return gains, nil
	}

	total := 0
	positiveElements := 0
	lastPositiveElement := ""
	for elem, amount := range gains {
		if amount <= 0 {
			continue
		}
		total += amount
		positiveElements++
		lastPositiveElement = elem
	}
	if total <= 1 {
		return gains, nil
	}

	limit := (total + 1) / 2
	if positiveElements == 1 {
		return map[string]int{lastPositiveElement: limit}, nil
	}

	selected := elementMapFromAction(action, "gain")
	if selected == nil {
		selected = elementMapFromAction(action, "gained")
	}
	if selected == nil {
		return nil, fmt.Errorf("first turn hero load requires choosing %d elements to gain", limit)
	}

	selectedTotal := 0
	for elem, amount := range selected {
		if amount < 0 || amount > gains[elem] {
			return nil, fmt.Errorf("invalid first turn hero load choice")
		}
		selectedTotal += amount
	}
	if selectedTotal != limit {
		return nil, fmt.Errorf("first turn hero load must gain exactly %d elements", limit)
	}
	return selected, nil
}

func (e *Engine) isFirstPlayerFirstTurnHeroConsume(playerID int, card *CardInstance) bool {
	if card == nil {
		return false
	}
	ps := e.State.Players[playerID]
	return e.State.IsFirstTurn && playerID == e.State.FirstPlayer && ps != nil && card == ps.Hero
}

func elementMapFromAction(action ActionMessage, key string) map[string]int {
	raw, ok := action.Data[key]
	if !ok || raw == nil {
		return nil
	}
	result := make(map[string]int)
	switch values := raw.(type) {
	case map[string]any:
		for elem, value := range values {
			switch amount := value.(type) {
			case float64:
				result[elem] = int(amount)
			case int:
				result[elem] = amount
			}
		}
	case map[string]int:
		for elem, amount := range values {
			result[elem] = amount
		}
	default:
		return nil
	}
	return result
}

func (e *Engine) handleConsume(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseMain && e.State.Phase != PhaseDefenseWindow {
		return fmt.Errorf("cannot consume now")
	}
	// During defense window, only the defending player can consume (透支)
	if e.State.Phase == PhaseDefenseWindow {
		if e.State.PendingSpell == nil || playerID == e.State.PendingSpell.AttackerID {
			return fmt.Errorf("only defender can overdraft during defense")
		}
	}
	if e.State.Phase == PhaseMain && e.State.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}

	instanceID, _ := action.Data["instance_id"].(string)
	ps := e.State.Players[playerID]

	// Find the card on the field
	card := e.findCardOnField(ps, instanceID)
	if card == nil {
		return fmt.Errorf("card not found on field")
	}
	if !e.canConsumeCard(card) {
		return fmt.Errorf("card cannot be consumed")
	}

	gains, err := e.elementsGainedFromConsume(playerID, card, action)
	if err != nil {
		return err
	}

	// Consume: set horizontal and gain elements
	card.IsHorizontal = true
	ps.GainElements(gains)

	e.emit(GameEvent{
		Type:   "consume",
		Player: -1,
		Data: map[string]any{
			"player":      playerID,
			"instance_id": instanceID,
			"elements":    ps.Elements,
			"gained":      gains,
		},
	})

	// Trigger 消耗 effects
	e.triggerEffects(TriggerOnConsume, card, nil, map[string]any{
		"gained": gains,
	})
	e.triggerFieldEffectsWithData(TriggerOnConsume, playerID, card, map[string]any{
		"consumed_player": playerID,
		"gained":          gains,
	})
	e.triggerFieldEffectsWithData(TriggerOnConsume, 1-playerID, card, map[string]any{
		"consumed_player": playerID,
		"gained":          gains,
	})
	e.advanceMastery(card, playerID, 1)
	e.destroyFuyeDoomedCardAfterExert(card)

	return nil
}

// handleCastSpell handles casting a spell
func (e *Engine) handleCastSpell(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseMain {
		return fmt.Errorf("not in main phase")
	}
	if e.State.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}

	instanceID, _ := action.Data["instance_id"].(string)
	targetType, _ := action.Data["target_type"].(string)
	targetColF, _ := action.Data["target_col"].(float64)
	targetRowF, _ := action.Data["target_row"].(float64)
	extraTargetColF, hasExtraTargetCol := action.Data["extra_target_col"].(float64)
	extraTargetRowF, hasExtraTargetRow := action.Data["extra_target_row"].(float64)
	boostIDsRaw, _ := action.Data["boost_ids"].([]any)

	ps := e.State.Players[playerID]

	skill := e.findSkill(ps, instanceID)
	if skill == nil {
		return fmt.Errorf("skill not found in skill area or bound skills")
	}

	if err := e.validateSkillForPurpose(skill, skillPurposeAttack); err != nil {
		return err
	}

	// Check cost
	cost := e.effectiveSkillUseCost(ps, skill)
	if !e.canPayCost(ps, cost) {
		return fmt.Errorf("not enough elements")
	}
	if skill.Card.Number == "3021011" && !validateSingleElementPayment(ps.Elements, cost, action) {
		return fmt.Errorf("overlord sanction cost must be paid with one element")
	}
	if len(boostIDsRaw) > 0 && !canSkillBeBoosted(skill) {
		return fmt.Errorf("skill cannot be boosted")
	}

	target := SpellTarget{
		Type:     targetType,
		Position: Position{Col: int(targetColF), Row: int(targetRowF)},
	}
	if ownerF, ok := action.Data["target_owner"].(float64); ok {
		ownerID := int(ownerF)
		target.OwnerID = &ownerID
	}
	// Process boost skills (法术强化)
	boostIDs := stringsFromAnySlice(boostIDsRaw)
	boostSkills, boostCost, err := e.collectSkillUses(ps, boostIDs, skillPurposeAttackBoost, map[string]bool{instanceID: true})
	if err != nil {
		return err
	}
	if err := e.validateSpellTargetWithPierce(playerID, skill, target, e.spellHasPierceWithBoosts(playerID, skill, boostSkills)); err != nil {
		return err
	}
	extraTargets := make([]SpellTarget, 0, 1)
	if skill.Card.Number == "3321001" && hasExtraTargetCol && hasExtraTargetRow {
		extra := SpellTarget{Type: "unit", Position: Position{Col: int(extraTargetColF), Row: int(extraTargetRowF)}}
		if err := e.validateSpellExtraTarget(playerID, extra); err != nil {
			return err
		}
		if extra.Position != target.Position {
			extraTargets = append(extraTargets, extra)
		}
	}
	totalCost := mergeElementCosts(cost, boostCost)
	if !e.canPayCost(ps, totalCost) {
		return fmt.Errorf("not enough elements for boost skills")
	}

	// Pay costs and set cards horizontal only after all validation succeeds.
	if !e.payCostForCardAction(ps, skill, cost, totalCost, paymentPurposeUse, action) {
		return fmt.Errorf("invalid payment")
	}
	skill.IsHorizontal = true
	tapSkills(boostSkills)

	// Apply cooldown from keyword
	if !e.shouldSkipCooldown(ps, skill) {
		e.ApplyKeywordOnSkillUse(skill)
	}
	e.applySkillUseCooldownModifiers(ps, append([]*CardInstance{skill}, boostSkills...)...)
	if skill.Card.Number == "3611101" {
		e.applyNextRedMoonModifiers(playerID, skill)
		e.refreshRedMoonState(playerID)
	}
	e.consumeNextSkillUseModifiers(ps, skill)
	e.advanceMasteryForUsedSkills(playerID, append([]*CardInstance{skill}, boostSkills...)...)

	powerTargets := append([]SpellTarget{target}, extraTargets...)
	totalPower := e.effectiveSpellPower(playerID, skill, boostSkills, powerTargets...)
	powerSources := e.spellPowerSources(playerID, skill, boostSkills, totalPower, powerTargets...)
	e.consumeNextSpellPowerBonuses(ps, skill)

	// Check if it's a 咒术 (sorcery - unblockable)
	isSorcery := isSorcerySkill(skill.Card)
	if ps.SpellsCastThisTurn == nil {
		ps.SpellsCastThisTurn = make(map[string]int)
	}
	ps.SpellsCastThisTurn[skill.Card.Category]++
	spellCastData := map[string]any{
		"cast_player": playerID,
		"attacker":    playerID,
		"skill":       cardToInfo(skill),
		"target":      target,
		"power":       totalPower,
		"boost_count": len(boostSkills),
		"is_sorcery":  isSorcery,
	}
	e.emit(GameEvent{
		Type:   "spell_cast",
		Player: -1,
		Data:   spellCastData,
	})
	e.triggerEffects(TriggerOnSpellCast, skill, nil, spellCastData)

	if isSorcery {
		resolveSorcery := func() {
			if e.shouldResolveSorceryHit(skill) {
				e.resolveSpellHit(playerID, skill, target, boostSkills, extraTargets)
			}
			e.removeStoredArchmageStaffSkillAfterUse(playerID, skill)
		}
		if e.triggerSpellCastFieldEffectsWithContinuation(playerID, skill, spellCastData, resolveSorcery) {
			return nil
		}
		resolveSorcery()
	} else {
		e.State.PendingSpell = &SpellCast{
			AttackerID:   playerID,
			Skill:        skill,
			Target:       target,
			TotalPower:   totalPower,
			PowerSources: powerSources,
			BoostSkills:  boostSkills,
			ExtraTargets: extraTargets,
		}
		resolveWithoutDefense := func() {
			e.resolvePendingSpellHit()
		}
		openDefenseWindow := func() {
			if e.State.PendingSpell == nil {
				return
			}
			e.State.ResumePhase = PhaseDefenseWindow
			e.State.Phase = PhaseDefenseWindow
			e.emit(GameEvent{
				Type:   "defense_window",
				Player: 1 - playerID,
				Data: map[string]any{
					"timeout": 30,
				},
			})
		}
		continueSpell := openDefenseWindow
		if !e.spellAllowsDefense(playerID, skill, target) {
			continueSpell = resolveWithoutDefense
		}
		if e.triggerSpellCastFieldEffectsWithContinuation(playerID, skill, spellCastData, continueSpell) {
			if e.spellAllowsDefense(playerID, skill, target) {
				e.State.ResumePhase = PhaseDefenseWindow
			}
			return nil
		}
		continueSpell()
	}

	return nil
}

func (e *Engine) shouldResolveSorceryHit(skill *CardInstance) bool {
	if skill == nil || skill.Card == nil {
		return false
	}
	switch skill.Card.Number {
	case "3001002":
		return false
	default:
		return true
	}
}

// handleDefend handles the defender's response to a spell
func (e *Engine) handleDefend(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseDefenseWindow {
		return fmt.Errorf("not in defense window")
	}
	if e.State.PendingSpell == nil {
		return fmt.Errorf("no pending spell")
	}
	if playerID == e.State.PendingSpell.AttackerID {
		return fmt.Errorf("attacker cannot defend")
	}

	defenseIDsRaw, _ := action.Data["skill_ids"].([]any)
	defenseScrollIDsRaw, _ := action.Data["scroll_ids"].([]any)
	boostIDsRaw, _ := action.Data["boost_ids"].([]any)
	overexertIDsRaw, _ := action.Data["overexert_ids"].([]any)

	ps := e.State.Players[playerID]

	defenseIDs := stringsFromAnySlice(defenseIDsRaw)
	defenseScrollIDs := stringsFromAnySlice(defenseScrollIDsRaw)
	boostIDs := stringsFromAnySlice(boostIDsRaw)
	overexertIDs := stringsFromAnySlice(overexertIDsRaw)
	defenseSkills, defenseCost, err := e.collectSkillUses(ps, defenseIDs, skillPurposeDefend, nil)
	if err != nil {
		return err
	}
	usedIDs := skillIDSet(defenseSkills)
	defenseScrolls, scrollCost, err := e.collectDefenseScrollUses(ps, defenseScrollIDs, usedIDs)
	if err != nil {
		return err
	}
	usedIDs = mergeSkillIDSet(usedIDs, skillIDSet(defenseScrolls))
	boostSkillIDs, boostScrollIDs := e.splitDefenseBoostIDs(ps, boostIDs)
	boostSkills, boostCost, err := e.collectSkillUses(ps, boostSkillIDs, skillPurposeDefenseBoost, usedIDs)
	if err != nil {
		return err
	}
	usedIDs = mergeSkillIDSet(usedIDs, skillIDSet(boostSkills))
	boostScrolls, boostScrollCost, err := e.collectDefenseBoostScrollUses(ps, boostScrollIDs, usedIDs)
	if err != nil {
		return err
	}
	overexertUnits, err := e.collectOverexertUnits(ps, overexertIDs)
	if err != nil {
		return err
	}
	totalCost := mergeElementCosts(defenseCost, scrollCost, boostCost, boostScrollCost)
	if !e.canPayCostWithOverexertOptions(ps, totalCost, overexertUnits, e.playerHasLightWildcard(ps)) {
		return fmt.Errorf("not enough elements for defense")
	}
	if len(defenseSkills)+len(defenseScrolls)+len(boostSkills) > 0 {
		if !e.payDefenseCostWithOptions(ps, totalCost, action, overexertUnits, e.playerHasLightWildcard(ps)) {
			return fmt.Errorf("invalid payment")
		}
		e.destroyFuyeDoomedAfterExert(overexertUnits)
		tapSkills(defenseSkills)
		tapSkills(boostSkills)
		e.moveHandConsumablesToGraveyard(ps, append(defenseScrolls, boostScrolls...))
		usedSkills := append([]*CardInstance{}, defenseSkills...)
		usedSkills = append(usedSkills, boostSkills...)
		e.advanceMasteryForUsedSkills(playerID, usedSkills...)
	}

	defenseSources := append([]*CardInstance{}, defenseSkills...)
	defenseSources = append(defenseSources, defenseScrolls...)
	boostSources := append([]*CardInstance{}, boostSkills...)
	boostSources = append(boostSources, boostScrolls...)

	continueAfterDefenseSpellCounters := func() {
		if e.State.PendingSpell == nil {
			return
		}
		if e.promptDispelDefenseSpellIfEligible(e.State.PendingSpell.AttackerID, playerID, defenseSources, boostSources, len(overexertUnits)) {
			return
		}
		e.finishDefenseResolution(playerID, defenseSources, boostSources, len(overexertUnits))
	}
	if e.promptDefenseSpellCastCounters(e.State.PendingSpell.AttackerID, playerID, defenseSources, boostSources, continueAfterDefenseSpellCounters) {
		return nil
	}

	continueAfterDefenseSpellCounters()
	return nil
}

func (e *Engine) promptDefenseSpellCastCounters(attackerID int, defenderID int, defenseSources []*CardInstance, boostSources []*CardInstance, afterDone func()) bool {
	type defenseSpellSource struct {
		card    *CardInstance
		purpose skillPurpose
	}
	sources := make([]defenseSpellSource, 0, len(defenseSources)+len(boostSources))
	for _, source := range defenseSources {
		if source != nil && source.Card != nil && isSpellLikeCard(source.Card) {
			sources = append(sources, defenseSpellSource{card: source, purpose: skillPurposeDefend})
		}
	}
	for _, source := range boostSources {
		if source != nil && source.Card != nil && isSpellLikeCard(source.Card) {
			sources = append(sources, defenseSpellSource{card: source, purpose: skillPurposeDefenseBoost})
		}
	}
	var promptNext func(int, bool)
	promptNext = func(index int, continuing bool) {
		for index < len(sources) {
			source := sources[index]
			index++
			power := e.effectiveSkillPowerForPurpose(defenderID, source.card, source.purpose)
			data := map[string]any{
				"cast_player": defenderID,
				"attacker":    defenderID,
				"skill":       cardToInfo(source.card),
				"power":       power,
				"is_sorcery":  isSorcerySkill(source.card.Card),
				"defense_use": true,
			}
			counters := e.eligibleCounterTraps(attackerID, TriggerOnSpellCast, source.card, data)
			if e.promptCounterTrapQueue(counters, TriggerOnSpellCast, source.card, data, func() {
				promptNext(index, true)
			}) {
				return
			}
		}
		if continuing && afterDone != nil {
			afterDone()
		}
	}
	promptNext(0, false)
	return e.State.PendingAction != nil && e.State.PendingAction.Type == "counter_trigger"
}

func mergeSkillIDSet(a map[string]bool, b map[string]bool) map[string]bool {
	merged := make(map[string]bool, len(a)+len(b))
	for id := range a {
		merged[id] = true
	}
	for id := range b {
		merged[id] = true
	}
	return merged
}

func (e *Engine) collectDefenseScrollUses(ps *PlayerState, ids []string, reserved map[string]bool) ([]*CardInstance, map[string]int, error) {
	scrolls := make([]*CardInstance, 0, len(ids))
	totalCost := make(map[string]int)
	seen := make(map[string]bool)
	for id := range reserved {
		seen[id] = true
	}
	for _, id := range ids {
		if seen[id] {
			return nil, nil, fmt.Errorf("defense source %s selected more than once", id)
		}
		seen[id] = true
		card, _ := ps.FindHandCard(id)
		if card == nil {
			return nil, nil, fmt.Errorf("defense scroll not found: %s", id)
		}
		if !isSpellScrollCard(card.Card) || !canUseSkillForPurpose(card.Card, skillPurposeDefend) {
			return nil, nil, fmt.Errorf("card %s is not a defense spell scroll", id)
		}
		if err := e.validateHandSpellScrollForPurpose(card, skillPurposeDefend); err != nil {
			return nil, nil, fmt.Errorf("defense scroll %s cannot be used for %s: %w", id, skillPurposeDefend, err)
		}
		scrolls = append(scrolls, card)
		for elem, amount := range e.effectiveCardPlayCost(ps, card) {
			totalCost[elem] += amount
		}
	}
	return scrolls, totalCost, nil
}

func (e *Engine) splitDefenseBoostIDs(ps *PlayerState, ids []string) ([]string, []string) {
	skillIDs := make([]string, 0, len(ids))
	scrollIDs := make([]string, 0)
	for _, id := range ids {
		if card, _ := ps.FindHandCard(id); card != nil && isSpellScrollCard(card.Card) {
			scrollIDs = append(scrollIDs, id)
			continue
		}
		skillIDs = append(skillIDs, id)
	}
	return skillIDs, scrollIDs
}

func (e *Engine) collectDefenseBoostScrollUses(ps *PlayerState, ids []string, reserved map[string]bool) ([]*CardInstance, map[string]int, error) {
	scrolls := make([]*CardInstance, 0, len(ids))
	totalCost := make(map[string]int)
	seen := make(map[string]bool)
	for id := range reserved {
		seen[id] = true
	}
	for _, id := range ids {
		if seen[id] {
			return nil, nil, fmt.Errorf("defense source %s selected more than once", id)
		}
		seen[id] = true
		card, _ := ps.FindHandCard(id)
		if card == nil {
			return nil, nil, fmt.Errorf("defense boost scroll not found: %s", id)
		}
		if !isSpellScrollCard(card.Card) {
			return nil, nil, fmt.Errorf("card %s is not a spell scroll", id)
		}
		if err := e.validateHandSpellScrollForPurpose(card, skillPurposeDefenseBoost); err != nil {
			return nil, nil, fmt.Errorf("defense boost scroll %s cannot be used for %s: %w", id, skillPurposeDefenseBoost, err)
		}
		scrolls = append(scrolls, card)
		for elem, amount := range e.effectiveCardPlayCost(ps, card) {
			totalCost[elem] += amount
		}
	}
	return scrolls, totalCost, nil
}

func (e *Engine) validateHandSpellScrollForPurpose(card *CardInstance, purpose skillPurpose) error {
	if card == nil || card.Card == nil || !isSpellScrollCard(card.Card) {
		return fmt.Errorf("card is not a spell scroll")
	}
	if e.hasEffectiveStatus(card, StatusPetrify) {
		return fmt.Errorf("card is petrified")
	}
	if !canUseSkillForPurpose(card.Card, purpose) {
		return fmt.Errorf("spell scroll cannot be used for %s", purpose)
	}
	return e.validateSkillUsePermissionModifiers(card, purpose)
}

func (e *Engine) moveHandConsumablesToGraveyard(ps *PlayerState, cards []*CardInstance) {
	for _, card := range cards {
		if card == nil {
			continue
		}
		_, idx := ps.FindHandCard(card.InstanceID)
		if idx < 0 {
			continue
		}
		ps.RemoveFromHand(idx)
		ps.Graveyard = append(ps.Graveyard, card)
		e.emit(GameEvent{
			Type:   "use_item",
			Player: -1,
			Data: map[string]any{
				"player":         ps.PlayerID,
				"card":           cardToInfo(card),
				"elements":       ps.Elements,
				"defense_scroll": true,
			},
		})
	}
}
func (e *Engine) finishDefenseResolution(playerID int, defenseSkills []*CardInstance, boostSkills []*CardInstance, overexerted int) {
	totalDefPower := e.totalEffectiveSkillPower(playerID, defenseSkills, skillPurposeDefend) +
		e.totalEffectiveSkillPower(playerID, boostSkills, skillPurposeDefenseBoost)

	attackPower := e.State.PendingSpell.TotalPower

	e.emit(GameEvent{
		Type:   "defense_attempt",
		Player: -1,
		Data: map[string]any{
			"defender":      playerID,
			"defense_power": totalDefPower,
			"attack_power":  attackPower,
			"skills_used":   len(defenseSkills) + len(boostSkills),
			"overexerted":   overexerted,
		},
	})

	defenseSuccess := attackPower <= 0 || (totalDefPower >= attackPower && len(defenseSkills) > 0)
	defendData := map[string]any{
		"defender":        playerID,
		"attacker":        e.State.PendingSpell.AttackerID,
		"defense_power":   totalDefPower,
		"attack_power":    attackPower,
		"defense_success": defenseSuccess,
		"attack_skill":    e.State.PendingSpell.Skill,
		"boost_skills":    e.State.PendingSpell.BoostSkills,
		"defense_skills":  defenseSkills,
		"defense_boosts":  boostSkills,
	}
	for _, defenseSkill := range defenseSkills {
		e.triggerEffects(TriggerOnDefend, defenseSkill, nil, defendData)
	}
	e.triggerFieldEffectsWithData(TriggerOnDefend, playerID, e.State.PendingSpell.Skill, defendData)
	e.triggerFieldEffectsWithData(TriggerOnDefend, e.State.PendingSpell.AttackerID, e.State.PendingSpell.Skill, defendData)

	if defenseSuccess {
		// Defense successful
		e.emit(GameEvent{
			Type:   "defense_success",
			Player: -1,
			Data:   map[string]any{"defender": playerID},
		})
		e.removeStoredArchmageStaffSkillAfterUse(e.State.PendingSpell.AttackerID, e.State.PendingSpell.Skill)
	} else {
		// Defense failed, spell hits
		if e.resolveSpellHit(
			e.State.PendingSpell.AttackerID,
			e.State.PendingSpell.Skill,
			e.State.PendingSpell.Target,
			e.State.PendingSpell.BoostSkills,
			e.State.PendingSpell.ExtraTargets,
		) {
			return
		}
		e.removeStoredArchmageStaffSkillAfterUse(e.State.PendingSpell.AttackerID, e.State.PendingSpell.Skill)
	}

	e.State.PendingSpell = nil
	if e.State.PendingAction == nil {
		e.State.Phase = PhaseMain
	}
	e.checkWinCondition()
}

func (e *Engine) promptDispelDefenseSpellIfEligible(attackerID int, defenderID int, defenseSkills []*CardInstance, boostSkills []*CardInstance, overexerted int) bool {
	defenseOnlySkills := make([]*CardInstance, 0, len(defenseSkills))
	for _, skill := range defenseSkills {
		if skill != nil && skill.Card != nil && isDefenseOnlySkill(skill.Card) {
			defenseOnlySkills = append(defenseOnlySkills, skill)
		}
	}
	if len(defenseOnlySkills) == 0 {
		return false
	}

	dispel := e.findReadyDispelSkill(attackerID)
	if dispel == nil {
		return false
	}

	candidates := make([]map[string]any, 0, len(defenseOnlySkills))
	validTargets := make(map[string]*CardInstance, len(defenseOnlySkills))
	for _, skill := range defenseOnlySkills {
		candidate := candidateInfo(skill, "skill", "enemy")
		candidates = append(candidates, candidate)
		validTargets[skill.InstanceID] = skill
	}

	cost := e.effectiveSkillUseCost(e.State.Players[attackerID], dispel)
	e.SetPendingActionWithError(attackerID, "dispel_defense_spell",
		"解咒:选择1个防御法术无效",
		candidates, 0, 1, cost, true,
		func(selected []string, data map[string]any) error {
			if len(selected) == 0 {
				e.finishDefenseResolution(defenderID, defenseSkills, boostSkills, overexerted)
				return nil
			}
			cancelledSkill := validTargets[selected[0]]
			if cancelledSkill == nil {
				return fmt.Errorf("invalid defense spell selection")
			}
			if err := e.payAndUseDispel(attackerID, dispel, cost, data); err != nil {
				return err
			}
			e.emit(GameEvent{
				Type:   "spell_reaction",
				Player: -1,
				Data: map[string]any{
					"player":    attackerID,
					"card":      cardToInfo(dispel),
					"effect":    "cancel_defense_spell",
					"cancelled": cardToInfo(cancelledSkill),
				},
			})
			e.finishDefenseResolution(defenderID, withoutCardInstance(defenseSkills, cancelledSkill), boostSkills, overexerted)
			return nil
		})
	return e.State.PendingAction != nil && e.State.PendingAction.Type == "dispel_defense_spell"
}

func (e *Engine) findReadyDispelSkill(playerID int) *CardInstance {
	ps := e.State.Players[playerID]
	for _, skill := range ps.Skills {
		if skill == nil || skill.Card == nil || skill.Card.Number != "3021010" {
			continue
		}
		if e.validateReadySkill(skill) == nil {
			return skill
		}
	}
	return nil
}

func (e *Engine) payAndUseDispel(playerID int, dispel *CardInstance, cost map[string]int, data map[string]any) error {
	if dispel == nil {
		return fmt.Errorf("dispel not found")
	}
	if err := e.validateReadySkill(dispel); err != nil {
		return err
	}
	overexertIDs := stringsFromAnySlice(anySliceFromData(data, "overexert_ids"))
	overexertUnits, err := e.collectOverexertUnits(e.State.Players[playerID], overexertIDs)
	if err != nil {
		return err
	}
	if !e.canPayCostWithOverexertOptions(e.State.Players[playerID], cost, overexertUnits, e.playerHasLightWildcard(e.State.Players[playerID])) {
		return fmt.Errorf("not enough elements")
	}
	if !e.payDefenseCostWithOptions(e.State.Players[playerID], cost, ActionMessage{Data: data}, overexertUnits, e.playerHasLightWildcard(e.State.Players[playerID])) {
		return fmt.Errorf("invalid payment")
	}
	e.destroyFuyeDoomedAfterExert(overexertUnits)
	dispel.IsHorizontal = true
	if !e.shouldSkipCooldown(e.State.Players[playerID], dispel) {
		e.ApplyKeywordOnSkillUse(dispel)
	}
	e.applySkillUseCooldownModifiers(e.State.Players[playerID], dispel)
	e.consumeNextSkillUseModifiers(e.State.Players[playerID], dispel)
	e.advanceMasteryForUsedSkills(playerID, dispel)
	return nil
}

func withoutCardInstance(cards []*CardInstance, removed *CardInstance) []*CardInstance {
	result := make([]*CardInstance, 0, len(cards))
	for _, card := range cards {
		if card != nil && removed != nil && card.InstanceID == removed.InstanceID {
			continue
		}
		result = append(result, card)
	}
	return result
}

func (e *Engine) collectOverexertUnits(ps *PlayerState, ids []string) ([]*CardInstance, error) {
	cards := make([]*CardInstance, 0, len(ids))
	seen := make(map[string]bool)
	for _, id := range ids {
		if seen[id] {
			return nil, fmt.Errorf("card %s selected more than once", id)
		}
		seen[id] = true
		card := e.findUnitOnGrid(ps, id)
		if card == nil {
			for _, equipment := range ps.Equipment {
				if equipment != nil && equipment.InstanceID == id {
					card = equipment
					break
				}
			}
		}
		if card == nil || card.Card == nil {
			return nil, fmt.Errorf("overexert card not found: %s", id)
		}
		if !e.canConsumeCard(card) {
			return nil, fmt.Errorf("card cannot be overexerted: %s", id)
		}
		cards = append(cards, card)
	}
	return cards, nil
}

// handleNoDefend handles when the defender chooses not to defend
func (e *Engine) handleNoDefend(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseDefenseWindow {
		return fmt.Errorf("not in defense window")
	}
	if e.State.PendingSpell == nil {
		return fmt.Errorf("no pending spell")
	}
	if playerID == e.State.PendingSpell.AttackerID {
		return fmt.Errorf("attacker cannot respond here")
	}

	e.resolvePendingSpellHit()
	return nil
}

func (e *Engine) resolvePendingSpellHit() {
	if e.State.PendingSpell == nil {
		return
	}
	spell := e.State.PendingSpell
	if e.resolveSpellHit(
		spell.AttackerID,
		spell.Skill,
		spell.Target,
		spell.BoostSkills,
		spell.ExtraTargets,
	) {
		return
	}

	e.removeStoredArchmageStaffSkillAfterUse(spell.AttackerID, spell.Skill)
	if e.State.PendingSpell == spell {
		e.State.PendingSpell = nil
	}
	if e.State.PendingAction == nil {
		e.State.Phase = PhaseMain
	}
	e.checkWinCondition()
}

func (e *Engine) cancelPendingSpell(playerID int, source *CardInstance, reason string) {
	if e.State.PendingSpell == nil {
		return
	}
	e.emit(GameEvent{
		Type:   "spell_cancelled",
		Player: -1,
		Data: map[string]any{
			"player": playerID,
			"card":   cardToInfo(source),
			"reason": reason,
		},
	})
	e.removeStoredArchmageStaffSkillAfterUse(e.State.PendingSpell.AttackerID, e.State.PendingSpell.Skill)
	e.State.PendingSpell = nil
	if e.State.PendingAction == nil {
		e.State.Phase = PhaseMain
	}
}

func (e *Engine) spellAllowsDefense(attackerID int, skill *CardInstance, target SpellTarget) bool {
	return e.spellDefenderID(attackerID, skill, target) != attackerID
}

func (e *Engine) spellDefenderID(attackerID int, skill *CardInstance, target SpellTarget) int {
	defenderID := 1 - attackerID
	if target.OwnerID != nil && *target.OwnerID == attackerID {
		defenderID = attackerID
	}
	if skill != nil && skill.Card != nil {
		if friendly, ok := behaviorForNumber(skill.Card.Number).(FriendlySpellTargetBehavior); ok && friendly.HasActiveFriendlySpellTarget(skill) && friendly.AllowsFriendlySpellTarget() && target.Type == "unit" && target.Position.Valid() {
			if e.State.Players[attackerID].Units[target.Position.Col][target.Position.Row] != nil {
				defenderID = attackerID
			}
		}
	}
	if target.Type == "hero" {
		defenderID = attackerID
	}
	return defenderID
}

// resolveSpellHit applies spell damage to the target. It returns true when a
// pre-hit counter prompt delayed resolution.
func (e *Engine) resolveSpellHit(attackerID int, skill *CardInstance, target SpellTarget, boostSkills []*CardInstance, extraTargets []SpellTarget) bool {
	e.beginResolution()
	defer e.endResolution()

	defenderID := e.spellDefenderID(attackerID, skill, target)
	affectedUnits := e.spellAffectedUnits(defenderID, skill, target)
	for _, extraTarget := range extraTargets {
		if extraTarget.Type != "unit" || !extraTarget.Position.Valid() {
			continue
		}
		extraUnit := e.State.Players[defenderID].Units[extraTarget.Position.Col][extraTarget.Position.Row]
		if extraUnit == nil {
			continue
		}
		alreadyIncluded := false
		for _, unit := range affectedUnits {
			if unit == extraUnit {
				alreadyIncluded = true
				break
			}
		}
		if !alreadyIncluded {
			affectedUnits = append(affectedUnits, extraUnit)
		}
	}
	if target.Type == "unit" && len(affectedUnits) == 0 {
		e.emit(GameEvent{
			Type:   "spell_miss",
			Player: -1,
			Data: map[string]any{
				"attacker": attackerID,
				"skill":    cardToInfo(skill),
				"target":   target,
				"reason":   "target_lost",
			},
		})
		return false
	}
	var targetUnit *CardInstance
	if target.Type == "hero" {
		targetUnit = e.State.Players[attackerID].Hero
	}
	if len(affectedUnits) > 0 {
		targetUnit = affectedUnits[0]
	}
	ctx := &EffectContext{
		Engine:     e,
		Source:     skill,
		Target:     targetUnit,
		PlayerID:   attackerID,
		OpponentID: defenderID,
		ExtraData:  map[string]any{"target": target},
	}
	dmg := max(skill.Card.Attack+skill.AttackBonus, 0)
	if override, ok := globalRegistry.SpellDamage(skill.Card.Number, ctx); ok {
		dmg = max(override, 0)
	}
	dmg = e.effectiveSpellDamage(attackerID, skill, dmg, boostSkills)
	e.consumeNextSpellAttackBonuses(e.State.Players[attackerID], skill)
	e.consumeNextElementSpellDamageBonus(e.State.Players[attackerID], skill)
	e.consumeAllSpellDamageZero(e.State.Players[attackerID], skill)
	e.consumeAllSpellDamageZero(e.State.Players[defenderID], skill)
	e.consumeFriendlySpellDamageMinus(e.State.Players[defenderID], skill)

	{
		totalPower := e.effectiveSpellPower(attackerID, skill, boostSkills, target)
		if e.State.PendingSpell != nil && e.State.PendingSpell.Skill == skill {
			totalPower = e.State.PendingSpell.TotalPower
		}
		hitCancelled := false
		hitData := map[string]any{
			"damage":           dmg,
			"power":            totalPower,
			"attacker":         attackerID,
			"spell_source":     skill,
			"target":           target,
			"affected_units":   affectedUnits,
			"boost_skills":     boostSkills,
			"cancel_spell_hit": &hitCancelled,
			"damage_ptr":       &dmg,
		}
		finishHit := func() {
			if hitCancelled {
				return
			}
			e.emit(GameEvent{
				Type:   "spell_hit",
				Player: -1,
				Data: map[string]any{
					"attacker": attackerID,
					"skill":    cardToInfo(skill),
					"target":   target,
					"damage":   dmg,
				},
			})

			hitData["skip_counter_traps"] = true
			hitData["timing"] = "before_damage"
			e.triggerEffects(TriggerOnSpellHitBeforeDamage, skill, targetUnit, hitData)
			e.triggerFieldEffectsWithData(TriggerOnSpellHitBeforeDamage, attackerID, skill, hitData)
			e.triggerFieldEffectsWithData(TriggerOnSpellHitBeforeDamage, defenderID, skill, hitData)
			if hitCancelled {
				return
			}

			if dmg > 0 {
				spellDamageData := map[string]any{
					"damage_source":  "spell",
					"damage_element": skill.Card.Category,
					"skill":          skill.Card.Number,
					"attacker":       attackerID,
					"boost_count":    len(boostSkills),
				}
				spellDamage := dmg
				if len(affectedUnits) > 1 {
					shieldTarget := targetUnit
					if shieldTarget == nil && len(affectedUnits) > 0 {
						shieldTarget = affectedUnits[0]
					}
					spellDamage = e.applyPlayerShieldDamage(shieldTarget, dmg, spellDamageData)
					spellDamageData["skip_player_shield"] = true
				}
				for _, unit := range affectedUnits {
					e.dealDamageWithExtra(unit, spellDamage, defenderID, spellDamageData)
				}
			}
			resolvedUnits := e.unitsStillOnField(affectedUnits)
			resolvedTargetUnit := targetUnit
			if target.Type != "hero" && !e.unitStillOnField(resolvedTargetUnit) {
				resolvedTargetUnit = nil
			}
			hitData["affected_units"] = resolvedUnits
			e.applyGenericSpellEffects(attackerID, defenderID, skill, resolvedUnits, target)
			e.applyTemporarySpellHitStatus(attackerID, skill, resolvedUnits)

			hitData["timing"] = "after_damage"
			e.triggerEffects(TriggerOnSpellHit, skill, resolvedTargetUnit, hitData)
			e.triggerFieldEffectsWithData(TriggerOnSpellHit, attackerID, skill, hitData)
			e.triggerFieldEffectsWithData(TriggerOnSpellHit, defenderID, skill, hitData)
			if skill.Statuses[StatusNextFrontRowRange] > 0 {
				skill.Statuses[StatusNextFrontRowRange]--
			}
		}
		afterCounterWindow := func() {
			finishHit()
			if e.State.PendingSpell != nil && e.State.PendingSpell.Skill == skill {
				e.removeStoredArchmageStaffSkillAfterUse(attackerID, skill)
				e.State.PendingSpell = nil
				if e.State.PendingAction == nil {
					e.State.Phase = PhaseMain
				}
				e.checkWinCondition()
			}
		}
		if e.promptCounterTrapQueue(e.eligibleCounterTraps(defenderID, TriggerOnSpellHitBeforeDamage, skill, hitData), TriggerOnSpellHitBeforeDamage, skill, hitData, afterCounterWindow) {
			return true
		}
		finishHit()
		return false
	}
}

func (e *Engine) unitsStillOnField(units []*CardInstance) []*CardInstance {
	result := make([]*CardInstance, 0, len(units))
	for _, unit := range units {
		if e.unitStillOnField(unit) {
			result = append(result, unit)
		}
	}
	return result
}

func (e *Engine) unitStillOnField(unit *CardInstance) bool {
	if unit == nil || unit.CurrentLife <= 0 || unit.OwnerID < 0 || unit.OwnerID >= len(e.State.Players) {
		return false
	}
	ps := e.State.Players[unit.OwnerID]
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if ps.Units[col][row] == unit {
				return true
			}
		}
	}
	return false
}

func (e *Engine) beginResolution() {
	e.resolutionDepth++
}

func (e *Engine) endResolution() {
	if e.resolutionDepth > 0 {
		e.resolutionDepth--
	}
	if e.resolutionDepth == 0 {
		e.resolvePendingDeaths()
	}
}

func (e *Engine) queueDeath(unit *CardInstance, ownerID int) {
	e.queueDeathWithData(unit, ownerID, nil)
}

func (e *Engine) queueDeathWithData(unit *CardInstance, ownerID int, deathData map[string]any) {
	if unit == nil {
		return
	}
	for i := range e.deathQueue {
		if e.deathQueue[i].unit == unit {
			if len(deathData) > 0 && len(e.deathQueue[i].deathData) == 0 {
				e.deathQueue[i].deathData = cloneExtraData(deathData)
			}
			return
		}
	}
	e.deathQueue = append(e.deathQueue, pendingDeath{unit: unit, ownerID: ownerID, deathData: cloneExtraData(deathData)})
}

func (e *Engine) removeQueuedDeath(unit *CardInstance) {
	if unit == nil || len(e.deathQueue) == 0 {
		return
	}
	kept := e.deathQueue[:0]
	for _, pending := range e.deathQueue {
		if pending.unit != unit {
			kept = append(kept, pending)
		}
	}
	e.deathQueue = kept
}

func (e *Engine) resolvePendingDeaths() {
	if e.resolvingDeaths {
		return
	}
	e.resolvingDeaths = true

	for len(e.deathQueue) > 0 {
		pending := e.deathQueue[0]
		e.deathQueue = e.deathQueue[1:]
		if pending.unit == nil || pending.unit.CurrentLife > 0 || !e.unitInOwnerGrid(pending.unit, pending.ownerID) {
			continue
		}
		e.destroyUnitWithData(pending.unit, pending.ownerID, pending.deathData)
	}

	e.resolvingDeaths = false
	e.resolveQueuedDeathTriggers()
	e.checkWinCondition()
}

func (e *Engine) queueDeathTriggers(unit *CardInstance, ownerID int) {
	e.queueDeathTriggersWithData(unit, ownerID, nil)
}

func (e *Engine) queueDeathTriggersWithData(unit *CardInstance, ownerID int, deathData map[string]any) {
	if unit == nil {
		return
	}
	e.deathTriggers = append(e.deathTriggers, pendingDeathTrigger{unit: unit, ownerID: ownerID, deathData: cloneExtraData(deathData)})
}

func (e *Engine) resolveQueuedDeathTriggers() {
	if len(e.deathTriggers) == 0 {
		return
	}
	e.beginResolution()
	defer e.endResolution()

	sort.SliceStable(e.deathTriggers, func(i, j int) bool {
		leftCurrent := e.deathTriggers[i].ownerID == e.State.CurrentTurn
		rightCurrent := e.deathTriggers[j].ownerID == e.State.CurrentTurn
		return leftCurrent && !rightCurrent
	})
	for len(e.deathTriggers) > 0 {
		pending := e.deathTriggers[0]
		e.deathTriggers = e.deathTriggers[1:]
		e.resolveDeathTriggersWithData(pending.unit, pending.ownerID, pending.deathData)
	}
}

func (e *Engine) resolveDeathTriggers(unit *CardInstance, ownerID int) {
	e.resolveDeathTriggersWithData(unit, ownerID, nil)
}

func (e *Engine) resolveDeathTriggersWithData(unit *CardInstance, ownerID int, deathData map[string]any) {
	if unit == nil || unit.Card == nil {
		return
	}

	// Trigger 遗言 (on death) effects
	e.triggerEffects(TriggerOnDeath, unit, nil, deathData)

	// Notify friendly cards about the death
	e.triggerFieldEffectsWithData(TriggerOnFriendlyDeath, ownerID, unit, deathData)

	// Notify enemy cards about the death
	e.triggerFieldEffectsWithData(TriggerOnEnemyDeath, 1-ownerID, unit, deathData)
}

func (e *Engine) unitInOwnerGrid(unit *CardInstance, ownerID int) bool {
	if unit == nil || ownerID < 0 || ownerID >= len(e.State.Players) {
		return false
	}
	ps := e.State.Players[ownerID]
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if ps.Units[col][row] == unit {
				return true
			}
		}
	}
	return false
}
func (e *Engine) spellAffectedUnits(defenderID int, skill *CardInstance, target SpellTarget) []*CardInstance {
	if target.Type != "unit" {
		return nil
	}
	defender := e.State.Players[defenderID]
	units := make([]*CardInstance, 0, 9)

	if skill != nil && skill.Card != nil && skill.Card.Number == "3511010" && defenderID != skill.OwnerID && target.Position.Valid() {
		targetUnit := defender.Units[target.Position.Col][target.Position.Row]
		if targetUnit != nil && targetUnit.Card != nil && targetUnit.Card.IsCompanion() {
			for col := 0; col < 3; col++ {
				for row := 0; row < 3; row++ {
					unit := defender.Units[col][row]
					if unit != nil && unit.Card != nil && unit.Card.Category == targetUnit.Card.Category {
						units = append(units, unit)
					}
				}
			}
			return units
		}
	}

	switch e.effectiveSpellArea(skill) {
	case SpellAreaSquare:
		startCol := min(max(target.Position.Col, 0), 1)
		startRow := min(max(target.Position.Row, 0), 1)
		for col := startCol; col < startCol+2; col++ {
			for row := startRow; row < startRow+2; row++ {
				if defender.Units[col][row] != nil {
					units = append(units, defender.Units[col][row])
				}
			}
		}
		return units
	case SpellAreaAll:
		for col := 0; col < 3; col++ {
			for row := 0; row < 3; row++ {
				if defender.Units[col][row] != nil {
					units = append(units, defender.Units[col][row])
				}
			}
		}
		return units
	case SpellAreaColumn:
		for row := 0; row < 3; row++ {
			if defender.Units[target.Position.Col][row] != nil {
				units = append(units, defender.Units[target.Position.Col][row])
			}
		}
		return units
	case SpellAreaFrontRow:
		frontRow := defender.GetFrontRow()
		if frontRow >= 0 {
			for col := 0; col < 3; col++ {
				if defender.Units[col][frontRow] != nil {
					units = append(units, defender.Units[col][frontRow])
				}
			}
		}
		return units
	case SpellAreaSplashCross:
		for _, delta := range []struct{ col, row int }{{0, 0}, {-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			col := target.Position.Col + delta.col
			row := target.Position.Row + delta.row
			if col < 0 || col >= 3 || row < 0 || row >= 3 {
				continue
			}
			if defender.Units[col][row] != nil {
				units = append(units, defender.Units[col][row])
			}
		}
		return units
	}

	if target.Position.Valid() && defender.Units[target.Position.Col][target.Position.Row] != nil {
		return []*CardInstance{defender.Units[target.Position.Col][target.Position.Row]}
	}
	return nil
}

func (e *Engine) effectiveSpellArea(skill *CardInstance) SpellArea {
	area := spellArea(skill)
	if skill == nil || skill.Card == nil || skill.OwnerID < 0 || skill.OwnerID >= len(e.State.Players) {
		return area
	}
	ps := e.State.Players[skill.OwnerID]
	ctx := &EffectContext{
		Engine:     e,
		Source:     skill,
		PlayerID:   skill.OwnerID,
		OpponentID: 1 - skill.OwnerID,
	}
	for _, fieldCard := range e.getAllFieldCards(ps) {
		if fieldCard == nil || fieldCard.Card == nil || e.hasEffectiveStatus(fieldCard, StatusPetrify) {
			continue
		}
		behavior := globalRegistry.GetBehavior(fieldCard.Card.Number)
		if modifier, ok := behavior.(SpellAreaModifier); ok && modifier.HasActiveSpellAreaModifier(fieldCard) {
			ctx.Target = fieldCard
			modifier.ModifySpellArea(ctx, &area)
		}
	}
	return area
}

func (e *Engine) applyGenericSpellEffects(attackerID int, defenderID int, skill *CardInstance, targets []*CardInstance, target SpellTarget) {
	e.applyGenericElementGain(attackerID, skill)
	for _, unit := range targets {
		e.applyExplicitSpellHitStatuses(skill, unit)
	}
}

// handleAttack handles a unit attacking another unit
func (e *Engine) handleAttack(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseMain {
		return fmt.Errorf("not in main phase")
	}
	if e.State.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}

	attackerID, _ := action.Data["attacker_id"].(string)
	targetColF, _ := action.Data["target_col"].(float64)
	targetRowF, _ := action.Data["target_row"].(float64)
	targetCol := int(targetColF)
	targetRow := int(targetRowF)

	ps := e.State.Players[playerID]
	opponent := e.State.Players[1-playerID]

	// Find attacker
	attacker := e.findUnitOnGrid(ps, attackerID)
	attackerIsEquipment := false
	if attacker == nil {
		attacker = e.findEquipment(ps, attackerID)
		attackerIsEquipment = attacker != nil
	}
	if attacker == nil {
		return fmt.Errorf("attacker not found")
	}
	if attacker.Card.Attack <= 0 {
		return fmt.Errorf("attacker has no attack")
	}
	if attacker.IsHorizontal {
		return fmt.Errorf("attacker is horizontal")
	}
	if e.hasEffectiveStatus(attacker, StatusStun) {
		return fmt.Errorf("attacker is stunned")
	}

	if !attackerIsEquipment {
		// Check attacker is in front row
		frontRow := ps.GetFrontRow()
		if attacker.Position == nil || attacker.Position.Row != frontRow {
			return fmt.Errorf("attacker is not in front row")
		}
	}

	// Check target is in attacker's range (default: enemy front row)
	targetPos := Position{Col: targetCol, Row: targetRow}
	if !targetPos.Valid() {
		return fmt.Errorf("invalid target position")
	}
	target := opponent.Units[targetCol][targetRow]
	if target == nil {
		return fmt.Errorf("no unit at target position")
	}
	if !e.isInDirectAttackRange(playerID, attacker, attackerIsEquipment, targetCol, targetRow) {
		return fmt.Errorf("target is not in attack range")
	}

	// Consume attacker (横置)
	attacker.IsHorizontal = true

	attackData := map[string]any{
		"attacker_player": playerID,
		"attacker":        attacker,
		"attack_source":   attackSourceKind(attackerIsEquipment),
		"target":          target,
		"target_pos":      targetPos,
	}

	// Trigger 攻击时 effects
	e.triggerEffects(TriggerOnAttack, attacker, target, attackData)

	// Trigger 受攻击时 effects before damage is dealt.
	e.triggerFieldEffectsWithData(TriggerOnAttacked, 1-playerID, attacker, attackData)
	e.triggerFieldEffectsWithData(TriggerOnAttacked, playerID, attacker, attackData)

	dmg := attacker.CurrentAttack

	e.emit(GameEvent{
		Type:   "unit_attack",
		Player: -1,
		Data: map[string]any{
			"attacker_player": playerID,
			"attacker":        cardToInfo(attacker),
			"attack_source":   attackSourceKind(attackerIsEquipment),
			"target":          cardToInfo(target),
			"target_pos":      targetPos,
			"damage":          dmg,
		},
	})

	// Deal damage (unit attacks cannot be defended)
	if dmg > 0 {
		e.dealDamageWithExtra(target, dmg, 1-playerID, map[string]any{"damage_source": "attack", "attacker": playerID})
		// Trigger 命中 effects
	}

	e.checkWinCondition()
	return nil
}

func (e *Engine) resolveForcedUnitAttack(attackerOwnerID int, attacker *CardInstance, target *CardInstance, reason string) {
	if attacker == nil || target == nil || attacker.CurrentAttack <= 0 {
		return
	}
	attackData := map[string]any{
		"attacker_player": attackerOwnerID,
		"attacker":        attacker,
		"attack_source":   "unit",
		"target":          target,
		"target_pos":      target.Position,
		"forced":          true,
		"reason":          reason,
	}
	e.triggerEffects(TriggerOnAttack, attacker, target, attackData)
	triggered := map[int]bool{}
	for _, ownerID := range []int{target.OwnerID, attackerOwnerID} {
		if ownerID < 0 || ownerID >= len(e.State.Players) || triggered[ownerID] {
			continue
		}
		triggered[ownerID] = true
		e.triggerFieldEffectsWithData(TriggerOnAttacked, ownerID, attacker, attackData)
	}

	dmg := attacker.CurrentAttack
	e.emit(GameEvent{
		Type:   "unit_attack",
		Player: -1,
		Data: map[string]any{
			"attacker_player": attackerOwnerID,
			"attacker":        cardToInfo(attacker),
			"attack_source":   "unit",
			"target":          cardToInfo(target),
			"target_pos":      target.Position,
			"damage":          dmg,
			"forced":          true,
			"reason":          reason,
		},
	})
	e.dealDamageWithExtra(target, dmg, target.OwnerID, map[string]any{
		"damage_source": "attack",
		"attacker":      attackerOwnerID,
		"forced_attack": true,
		"reason":        reason,
	})
}

func attackSourceKind(isEquipment bool) string {
	if isEquipment {
		return "equipment"
	}
	return "unit"
}

func (e *Engine) isInDirectAttackRange(playerID int, attacker *CardInstance, attackerIsEquipment bool, targetCol, targetRow int) bool {
	if attackerIsEquipment {
		return e.isEnemyFrontRowAttackTarget(playerID, attacker, targetCol, targetRow)
	}
	return e.IsInAttackRange(playerID, attacker, targetCol, targetRow)
}

func (e *Engine) isEnemyFrontRowAttackTarget(playerID int, attacker *CardInstance, targetCol, targetRow int) bool {
	opponent := e.State.Players[1-playerID]
	enemyFront := opponent.GetFrontRow()
	if enemyFront == -1 || targetRow != enemyFront {
		return false
	}
	target := opponent.Units[targetCol][targetRow]
	if target == nil {
		return false
	}
	if e.hasStealthFromOpponent(playerID, target) {
		return false
	}
	return true
}

// dealDamage deals damage to a card instance
func (e *Engine) dealDamage(target *CardInstance, amount int, ownerID int) {
	e.dealDamageWithExtra(target, amount, ownerID, nil)
}

func (e *Engine) fieldDamagePreventionSource(target *CardInstance, ownerID int, damageData map[string]any) *CardInstance {
	if target == nil || ownerID < 0 || ownerID >= len(e.State.Players) {
		return nil
	}
	ps := e.State.Players[ownerID]
	for _, source := range e.getAllFieldCards(ps) {
		if source == nil || source.Card == nil || source == target {
			continue
		}
		behavior, ok := globalRegistry.GetBehavior(source.Card.Number).(FieldDamagePreventionBehavior)
		if !ok || !behavior.HasActiveFieldDamagePrevention(source) {
			continue
		}
		ctx := &EffectContext{
			Engine:     e,
			Source:     source,
			Target:     target,
			PlayerID:   ownerID,
			OpponentID: 1 - ownerID,
			ExtraData:  damageData,
		}
		if behavior.PreventsFieldDamage(ctx) {
			return source
		}
	}
	return nil
}

func (e *Engine) dealDamageWithExtra(target *CardInstance, amount int, ownerID int, extraData map[string]any) {
	damageData := map[string]any{
		"damage": amount,
	}
	for key, value := range extraData {
		damageData[key] = value
	}
	if behavior, ok := globalRegistry.GetBehavior(target.Card.Number).(DamagePreventionBehavior); ok && behavior.HasActiveDamagePrevention(target) {
		ctx := &EffectContext{
			Engine:     e,
			Source:     target,
			Target:     target,
			PlayerID:   ownerID,
			OpponentID: 1 - ownerID,
			ExtraData:  damageData,
		}
		if behavior.PreventsDamage(ctx) {
			e.emit(GameEvent{
				Type:   "damage_prevented",
				Player: -1,
				Data: map[string]any{
					"target": cardToInfo(target),
					"amount": amount,
				},
			})
			return
		}
	}
	if source := e.fieldDamagePreventionSource(target, ownerID, damageData); source != nil {
		e.emit(GameEvent{
			Type:   "damage_prevented",
			Player: -1,
			Data: map[string]any{
				"source": cardToInfo(source),
				"target": cardToInfo(target),
				"amount": amount,
				"reason": "field_prevention",
			},
		})
		return
	}

	if target.Statuses[sturdyScrollShieldStatus] > 0 && target.Statuses[sturdyScrollShieldUntilStatus] >= e.State.TurnNumber {
		prevented := min(amount, target.Statuses[sturdyScrollShieldStatus])
		target.Statuses[sturdyScrollShieldStatus] -= prevented
		if target.Statuses[sturdyScrollShieldStatus] <= 0 {
			delete(target.Statuses, sturdyScrollShieldStatus)
			delete(target.Statuses, sturdyScrollShieldUntilStatus)
		}
		amount -= prevented
		e.emit(GameEvent{
			Type:   "damage_prevented",
			Player: -1,
			Data: map[string]any{
				"target": cardToInfo(target),
				"amount": prevented,
				"reason": "sturdy_scroll",
			},
		})
		if amount <= 0 {
			return
		}
	}

	amount = e.applyPlayerShieldDamage(target, amount, damageData)
	if amount <= 0 {
		return
	}
	if target.Statuses["防止致命"] > 0 && target.CurrentLife-amount <= 0 {
		target.Statuses["防止致命"]--
		if target.Statuses["防止致命"] <= 0 {
			delete(target.Statuses, "防止致命")
		}
		e.emit(GameEvent{
			Type:   "damage_prevented",
			Player: -1,
			Data: map[string]any{
				"target": cardToInfo(target),
				"amount": damageData["damage"],
				"reason": "prevent_lethal",
			},
		})
		return
	}

	if e.promptDolphinPartnerPrevention(target, amount, ownerID, damageData) {
		return
	}

	target.CurrentLife -= amount
	target.DamageTakenThisTurn += amount

	e.emit(GameEvent{
		Type:   "damage",
		Player: -1,
		Data: map[string]any{
			"target":    cardToInfo(target),
			"amount":    amount,
			"remaining": target.CurrentLife,
		},
	})

	damageData["damage"] = amount
	damageData["damage_taken_this_turn"] = target.DamageTakenThisTurn

	// Trigger 受伤 effects
	e.triggerEffects(TriggerOnDamaged, target, nil, damageData)
	fieldDamageData := map[string]any{
		"damaged_player": ownerID,
		"damage":         amount,
	}
	for key, value := range extraData {
		fieldDamageData[key] = value
	}
	e.triggerFieldEffectsWithData(TriggerOnDamaged, ownerID, target, fieldDamageData)
	enemyDamageData := map[string]any{"damaged_player": ownerID, "damage": amount}
	for key, value := range extraData {
		enemyDamageData[key] = value
	}
	e.triggerFieldEffectsWithData(TriggerOnDamaged, 1-ownerID, target, enemyDamageData)
	e.triggerHiddenFriendlyDamaged(ownerID, target, fieldDamageData)

	if target.CurrentLife <= 0 {
		if attacker, ok := damageData["attacker"].(int); ok {
			target.Statuses["lethal_source_player"] = attacker + 1
		} else {
			delete(target.Statuses, "lethal_source_player")
		}
		e.queueDeathWithData(target, ownerID, extraData)
		if e.resolutionDepth == 0 && !e.resolvingDeaths {
			e.resolvePendingDeaths()
		}
	}
}

func (e *Engine) promptDolphinPartnerPrevention(target *CardInstance, amount int, ownerID int, damageData map[string]any) bool {
	if e == nil || target == nil || amount <= 0 || target.CurrentLife-amount > 0 {
		return false
	}
	if skip, _ := damageData["skip_dolphin_prevention"].(bool); skip {
		return false
	}
	ps := e.State.Players[ownerID]
	if ps == nil {
		return false
	}
	candidates := make([]map[string]any, 0)
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := ps.Units[col][row]
			if unit == nil || unit == target || unit.Card.Number != "1221001" {
				continue
			}
			candidates = append(candidates, candidateInfo(unit, "unit", "own"))
		}
	}
	if len(candidates) == 0 {
		return false
	}
	retryData := make(map[string]any, len(damageData)+1)
	for key, value := range damageData {
		retryData[key] = value
	}
	retryData["skip_dolphin_prevention"] = true
	prompt := fmt.Sprintf("海豚伙伴:是否牺牲1个海豚伙伴，防止%s受到的%d点伤害", target.Card.Name, amount)
	e.SetPendingAction(ownerID, "dolphin_prevent_lethal", prompt, candidates, 0, 1, func(selected []string) {
		dolphin := selectedUnitFromCandidates(e, selected, candidates)
		if dolphin == nil || dolphin.OwnerID != ownerID || dolphin == target || dolphin.Card.Number != "1221001" {
			e.dealDamageWithExtra(target, amount, ownerID, retryData)
			return
		}
		e.destroyUnitWithCause(dolphin, ownerID, DeathCauseSacrifice)
		e.emit(GameEvent{
			Type:   "damage_prevented",
			Player: -1,
			Data: map[string]any{
				"source": cardToInfo(dolphin),
				"target": cardToInfo(target),
				"amount": amount,
				"reason": "dolphin_partner",
			},
		})
	})
	return true
}

func (e *Engine) triggerHiddenFriendlyDamaged(playerID int, target *CardInstance, extraData map[string]any) {
	ps := e.State.Players[playerID]
	hidden := append([]*CardInstance{}, ps.Hand...)
	hidden = append(hidden, ps.Deck...)
	for _, card := range hidden {
		if card == nil || card.Card == nil {
			continue
		}
		behavior, ok := globalRegistry.GetBehavior(card.Card.Number).(OnFriendlyDamagedFromHiddenBehavior)
		if !ok || !behavior.HasActiveFriendlyDamagedFromHidden(card) {
			continue
		}
		ctx := &EffectContext{
			Engine:     e,
			Source:     card,
			Target:     target,
			PlayerID:   playerID,
			OpponentID: 1 - playerID,
			ExtraData:  extraData,
		}
		_ = behavior.OnFriendlyDamagedFromHidden(ctx)
	}
}

// destroyUnit removes a unit from the field and sends it to graveyard
func (e *Engine) destroyUnit(unit *CardInstance, ownerID int) {
	e.destroyUnitWithData(unit, ownerID, nil)
}

func (e *Engine) destroyUnitWithCause(unit *CardInstance, ownerID int, cause string) {
	data := map[string]any{}
	if cause != "" {
		data["death_cause"] = cause
	}
	e.destroyUnitWithData(unit, ownerID, data)
}

func (e *Engine) destroyUnitWithData(unit *CardInstance, ownerID int, deathData map[string]any) {
	e.removeQueuedDeath(unit)
	ps := e.State.Players[ownerID]

	// Remove from grid
	if unit.Position != nil {
		ps.Units[unit.Position.Col][unit.Position.Row] = nil
	}

	// Bound skills live only while their host is on the battlefield. They do not
	// enter the graveyard as independent cards.
	e.releaseUnderCardsToGraveyard(ownerID, unit)
	unit.BoundSkills = nil

	// Add to graveyard
	ps.Graveyard = append(ps.Graveyard, unit)

	e.emit(GameEvent{
		Type:   "unit_destroyed",
		Player: -1,
		Data: map[string]any{
			"player": ownerID,
			"card":   cardToInfo(unit),
		},
	})

	if e.resolvingDeaths || e.resolutionDepth > 0 {
		e.queueDeathTriggersWithData(unit, ownerID, deathData)
	} else {
		e.resolveDeathTriggersWithData(unit, ownerID, deathData)
	}

	// Check if hero died. During queued death resolution, wait until all pending
	// deaths are processed so simultaneous hero deaths can become a draw.
	if unit.Card.IsHero() && !e.resolvingDeaths && e.resolutionDepth == 0 {
		e.checkWinCondition()
	}
}

func (e *Engine) initialHandSizeForPlayer(ps *PlayerState) int {
	if ps == nil {
		return 4
	}
	size := 4
	if ps.Hero != nil && ps.Hero.Card != nil && ps.Hero.Card.Number == "4311002" {
		size++
	}
	return size
}

// handleEquip handles equipping an item
func (e *Engine) handleEquip(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseMain {
		return fmt.Errorf("not in main phase")
	}
	if e.State.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}

	instanceID, _ := action.Data["instance_id"].(string)
	replaceID, _ := action.Data["replace_id"].(string)
	ps := e.State.Players[playerID]

	card, handIdx := ps.FindHandCard(instanceID)
	if card == nil {
		return fmt.Errorf("card not found in hand")
	}
	if !card.Card.IsItem() {
		return fmt.Errorf("card is not an item")
	}
	if !isEquipmentCard(card.Card) {
		return fmt.Errorf("card is not equipment")
	}
	cost := e.effectiveCardPlayCost(ps, card)
	if !e.canPayCost(ps, cost) {
		return fmt.Errorf("not enough elements")
	}

	slotIdx := -1
	var replacedEquipment *CardInstance
	newSubtype := restrictedEquipmentSubtype(card.Card)
	if replaceID != "" {
		for i := 0; i < equipmentSlotCapacity(ps); i++ {
			if ps.Equipment[i] != nil && ps.Equipment[i].InstanceID == replaceID {
				if ps.Equipment[i].IsHorizontal {
					return fmt.Errorf("can only replace vertical equipment")
				}
				if newSubtype != "" && restrictedEquipmentSubtype(ps.Equipment[i].Card) != newSubtype {
					return fmt.Errorf("restricted equipment can only replace same subtype")
				}
				replacedEquipment = ps.Equipment[i]
				slotIdx = i
				break
			}
		}
		if slotIdx == -1 {
			return fmt.Errorf("replacement equipment not found")
		}
	} else {
		if newSubtype != "" {
			for _, equipment := range ps.Equipment {
				if equipment != nil && restrictedEquipmentSubtype(equipment.Card) == newSubtype {
					if equipment.IsHorizontal {
						return fmt.Errorf("same subtype equipment is horizontal and cannot be replaced")
					}
					return fmt.Errorf("same subtype equipment must be replaced")
				}
			}
		}
		// Find empty equipment slot
		for i := 0; i < equipmentSlotCapacity(ps); i++ {
			if ps.Equipment[i] == nil {
				slotIdx = i
				break
			}
		}
		if slotIdx == -1 {
			return fmt.Errorf("equipment area is full")
		}
	}

	if !e.payCostForCardAction(ps, card, cost, cost, paymentPurposePlay, action) {
		return fmt.Errorf("invalid payment")
	}
	e.notifyCardPlayCostPaid(ps, card)
	ps.RemoveFromHand(handIdx)
	if replacedEquipment != nil {
		e.moveEquipmentToGraveyard(playerID, slotIdx, replacedEquipment)
	}
	card.IsHorizontal = true
	card.SlotIndex = slotIdx
	card.EnterTurn = e.State.TurnNumber
	ps.Equipment[slotIdx] = card

	e.emit(GameEvent{
		Type:   "equip",
		Player: -1,
		Data: map[string]any{
			"player":   playerID,
			"card":     cardToInfo(card),
			"slot":     slotIdx,
			"elements": ps.Elements,
		},
	})

	e.triggerEffects(TriggerOnEquip, card, nil, nil)
	e.triggerEffects(TriggerOnEnter, card, nil, nil)

	return nil
}

// handleLearnSkill handles learning a skill from the skill pool
func (e *Engine) handleLearnSkill(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseMain {
		return fmt.Errorf("not in main phase")
	}
	if e.State.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}

	instanceID, _ := action.Data["instance_id"].(string)
	replaceID, _ := action.Data["replace_id"].(string) // optional: which skill to replace

	ps := e.State.Players[playerID]

	// Find skill in skill pool
	var skill *CardInstance
	var poolIdx int
	for i, s := range ps.SkillPool {
		if s.InstanceID == instanceID {
			skill = s
			poolIdx = i
			break
		}
	}
	if skill == nil {
		return fmt.Errorf("skill not found in skill pool")
	}

	// Check cost
	cost := e.effectiveSkillLearnCost(ps, skill)
	if !e.canPayCost(ps, cost) {
		return fmt.Errorf("not enough elements")
	}
	if skill.Card.Number == "3021011" && !validateSingleElementPayment(ps.Elements, cost, action) {
		return fmt.Errorf("overlord sanction cost must be paid with one element")
	}

	// Find slot
	slotIdx := -1
	var replacedSkill *CardInstance
	if replaceID != "" {
		// Replace existing skill
		for i := 0; i < skillSlotCapacity(ps); i++ {
			if ps.Skills[i] != nil && ps.Skills[i].InstanceID == replaceID {
				if ps.Skills[i].IsHorizontal {
					return fmt.Errorf("can only replace vertical skills")
				}
				replacedSkill = ps.Skills[i]
				slotIdx = i
				break
			}
		}
		if slotIdx == -1 {
			return fmt.Errorf("replacement skill not found")
		}
	} else {
		// Find empty slot
		for i := 0; i < skillSlotCapacity(ps); i++ {
			if ps.Skills[i] == nil {
				slotIdx = i
				break
			}
		}
		if slotIdx == -1 {
			return fmt.Errorf("skill area is full, must replace an existing skill")
		}
	}

	// Pay cost and place
	if !e.payCostForCardAction(ps, skill, cost, cost, paymentPurposeLearn, action) {
		return fmt.Errorf("invalid payment")
	}
	e.notifyCardPlayCostPaid(ps, skill)
	e.consumeEarthSkillLearnCostModifier(ps, skill)
	ps.SkillPool = append(ps.SkillPool[:poolIdx], ps.SkillPool[poolIdx+1:]...)
	if replacedSkill != nil {
		ps.Skills[slotIdx] = nil
		returnSkillToPool(replacedSkill)
		ps.SkillPool = append(ps.SkillPool, replacedSkill)
		if replacedSkill.Card.Number == "3611101" {
			e.refreshRedMoonState(playerID)
		}
	}
	skill.IsHorizontal = true
	skill.SlotIndex = slotIdx
	skill.EnterTurn = e.State.TurnNumber
	e.ApplyKeywordOnEnter(skill)
	ps.Skills[slotIdx] = skill
	for _, modifier := range append([]TemporaryModifier(nil), ps.TempModifiers...) {
		if modifier.Type == TempModNextLearnedSkillHaste && modifier.RemainingUses != 0 {
			skill.IsHorizontal = false
			e.removeTemporaryModifier(playerID, modifier.ID)
			break
		}
	}

	e.emit(GameEvent{
		Type:   "learn_skill",
		Player: -1,
		Data: map[string]any{
			"player":   playerID,
			"card":     cardToInfo(skill),
			"slot":     slotIdx,
			"elements": ps.Elements,
		},
	})

	return nil
}

func (e *Engine) learnSkillFromPoolWithoutCost(playerID int, instanceID string, replaceID string) bool {
	ps := e.State.Players[playerID]
	if ps == nil {
		return false
	}
	var skill *CardInstance
	poolIdx := -1
	for i, s := range ps.SkillPool {
		if s != nil && s.InstanceID == instanceID {
			skill = s
			poolIdx = i
			break
		}
	}
	if skill == nil || skill.Card == nil || !skill.Card.IsSkill() {
		return false
	}

	slotIdx := -1
	var replacedSkill *CardInstance
	if replaceID != "" {
		for i := 0; i < skillSlotCapacity(ps); i++ {
			if ps.Skills[i] != nil && ps.Skills[i].InstanceID == replaceID && !ps.Skills[i].IsHorizontal {
				replacedSkill = ps.Skills[i]
				slotIdx = i
				break
			}
		}
	} else {
		for i := 0; i < skillSlotCapacity(ps); i++ {
			if ps.Skills[i] == nil {
				slotIdx = i
				break
			}
		}
	}
	if slotIdx == -1 {
		return false
	}

	ps.SkillPool = append(ps.SkillPool[:poolIdx], ps.SkillPool[poolIdx+1:]...)
	if replacedSkill != nil {
		ps.Skills[slotIdx] = nil
		returnSkillToPool(replacedSkill)
		ps.SkillPool = append(ps.SkillPool, replacedSkill)
		if replacedSkill.Card.Number == "3611101" {
			e.refreshRedMoonState(playerID)
		}
	}
	skill.IsHorizontal = true
	skill.SlotIndex = slotIdx
	skill.EnterTurn = e.State.TurnNumber
	e.ApplyKeywordOnEnter(skill)
	ps.Skills[slotIdx] = skill
	for _, modifier := range append([]TemporaryModifier(nil), ps.TempModifiers...) {
		if modifier.Type == TempModNextLearnedSkillHaste && modifier.RemainingUses != 0 {
			skill.IsHorizontal = false
			e.removeTemporaryModifier(playerID, modifier.ID)
			break
		}
	}

	e.emit(GameEvent{
		Type:   "learn_skill",
		Player: -1,
		Data: map[string]any{
			"player":   playerID,
			"card":     cardToInfo(skill),
			"slot":     slotIdx,
			"elements": ps.Elements,
		},
	})
	return true
}

func returnSkillToPool(skill *CardInstance) {
	if skill == nil {
		return
	}
	skill.IsHorizontal = true
	skill.Position = nil
	skill.SlotIndex = -1
	skill.EnterTurn = 0
	skill.UsedThisTurn = 0
	skill.UltimateUsed = false
	skill.Statuses = make(map[string]int)
	skill.ElementsGainBonus = make(map[string]int)
	skill.ElementsGainSet = nil
	skill.PowerBonus = 0
	skill.AttackBonus = 0
	skill.AttachedBehaviors = nil
}

// handleUseItem handles using a consumable item from hand
func (e *Engine) handleUseItem(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseMain {
		return fmt.Errorf("not in main phase")
	}
	if e.State.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}

	instanceID, _ := action.Data["instance_id"].(string)
	ps := e.State.Players[playerID]

	card, handIdx := ps.FindHandCard(instanceID)
	if card == nil {
		return fmt.Errorf("card not found in hand")
	}
	if !card.Card.IsItem() {
		return fmt.Errorf("card is not an item")
	}
	if isCounterTrapCard(card.Card.Number) {
		return e.placeCounterTrap(playerID, card, handIdx)
	}

	// Check if this is a terrain card - terrain cards go to battlefield
	if cards.IsTerrain(card.Card.Number) {
		// Redirect to terrain placement handler
		// Re-use the action but call the terrain handler
		colF, _ := action.Data["col"].(float64)
		rowF, _ := action.Data["row"].(float64)
		return e.handlePlaceTerrain(playerID, ActionMessage{
			Action: "place_terrain",
			Data: map[string]any{
				"instance_id": instanceID,
				"col":         colF,
				"row":         rowF,
				"payment":     action.Data["payment"],
			},
		})
	}
	if isSpellScrollCard(card.Card) {
		return e.handleUseSpellScrollItem(playerID, action, card, handIdx)
	}
	if cards.IsEquipment(card.Card.Number) {
		return fmt.Errorf("equipment cannot be used as a consumable item")
	}
	if !cards.IsConsumable(card.Card.Number) {
		if behavior, ok := globalRegistry.GetBehavior(card.Card.Number).(OnUseItemBehavior); !ok || !behavior.HasActiveUseItem(card) {
			return fmt.Errorf("item is not consumable")
		}
	}
	if err := e.validateConsumableItemUse(playerID, card); err != nil {
		return err
	}

	// Regular consumable item
	cost := e.effectiveCardPlayCost(ps, card)
	if !e.canPayCost(ps, cost) {
		return fmt.Errorf("not enough elements")
	}

	// Pay and use
	if !e.payCostForAction(ps, cost, action) {
		return fmt.Errorf("invalid payment")
	}
	e.notifyCardPlayCostPaid(ps, card)
	ps.RemoveFromHand(handIdx)
	ps.Graveyard = append(ps.Graveyard, card)

	e.emit(GameEvent{
		Type:   "use_item",
		Player: -1,
		Data: map[string]any{
			"player":   playerID,
			"card":     cardToInfo(card),
			"elements": ps.Elements,
		},
	})

	cancelled := false
	useData := map[string]any{"used_player": playerID, "cancel_item": &cancelled}
	resolveItem := func() {
		if cancelled {
			return
		}
		e.triggerEffects(TriggerOnUseItem, card, nil, nil)
		e.triggerFieldEffectsWithData(TriggerOnUseItem, playerID, card, useData)
		e.triggerFieldEffectsWithData(TriggerOnUseItem, 1-playerID, card, useData)
	}
	if e.promptOpponentCounterTrap(playerID, TriggerOnUseItem, card, useData, resolveItem) {
		return nil
	}
	resolveItem()

	return nil
}

func (e *Engine) validateConsumableItemUse(playerID int, card *CardInstance) error {
	if card == nil || card.Card == nil {
		return nil
	}
	switch card.Card.Number {
	case "2021010":
		if len(e.enemySkills(playerID, nil)) < 4 {
			return fmt.Errorf("Sealing Scroll requires the enemy to have at least four skills")
		}
	case "2611002":
		if !e.demonContractHasPayablePathAfterEntryCost(playerID, card) {
			return fmt.Errorf("Demon Contract requires a payable sacrifice and target")
		}
	case "2021012":
		if len(e.sketchScrollSkillCandidates(playerID)) == 0 {
			return fmt.Errorf("Sketch Scroll requires a payable learned attack spell")
		}
	case "2221101":
		if len(e.friendlySkillsIncludingBound(playerID, func(skill *CardInstance) bool {
			return skill != nil && skill.Card != nil && isSpellLikeCard(skill.Card)
		})) == 0 {
			return fmt.Errorf("Mirrorsea Spring requires a friendly spell")
		}
	case "2521102":
		if !e.hasEnemySetCounter(playerID) && !e.hasEnemyFrontStealth(playerID) {
			return fmt.Errorf("Moonlight Dust requires enemy set counters or stealthy front enemies")
		}
	case "2621111":
		if countShadowCompanionsInGraveyard(e.State.Players[playerID]) < 5 {
			return fmt.Errorf("Dark Burst Scroll requires at least five shadow companions in graveyard")
		}
	}
	return nil
}

func (e *Engine) handleUseSpellScrollItem(playerID int, action ActionMessage, card *CardInstance, handIdx int) error {
	ps := e.State.Players[playerID]
	if isDefenseOnlySkill(card.Card) {
		return fmt.Errorf("defense spell scroll can only be used during a defense window")
	}

	target := SpellTarget{Type: "none"}
	if skillNeedsTargetInstance(card) {
		targetType, _ := action.Data["target_type"].(string)
		if targetType == "hero" {
			target = SpellTarget{Type: "hero"}
			if err := e.validateSpellTarget(playerID, card, target); err != nil {
				return err
			}
		} else {
			colF, hasCol := action.Data["target_col"].(float64)
			rowF, hasRow := action.Data["target_row"].(float64)
			if !hasCol || !hasRow {
				return fmt.Errorf("spell scroll requires a target")
			}
			target = SpellTarget{Type: "unit", Position: Position{Col: int(colF), Row: int(rowF)}}
			if ownerF, ok := action.Data["target_owner"].(float64); ok {
				ownerID := int(ownerF)
				target.OwnerID = &ownerID
			}
			if err := e.validateSpellTarget(playerID, card, target); err != nil {
				return err
			}
		}
	}

	cost := e.effectiveCardPlayCost(ps, card)
	if !e.canPayCost(ps, cost) {
		return fmt.Errorf("not enough elements")
	}
	if !e.payCostForAction(ps, cost, action) {
		return fmt.Errorf("invalid payment")
	}
	e.notifyCardPlayCostPaid(ps, card)
	ps.RemoveFromHand(handIdx)
	ps.Graveyard = append(ps.Graveyard, card)

	e.emit(GameEvent{
		Type:   "use_item",
		Player: -1,
		Data: map[string]any{
			"player":   playerID,
			"card":     cardToInfo(card),
			"elements": ps.Elements,
		},
	})

	cancelled := false
	useData := map[string]any{"used_player": playerID, "cancel_item": &cancelled, "spell_scroll": true}
	resolveItem := func() {
		if cancelled {
			return
		}
		e.startSpellScrollCast(playerID, card, target)
	}
	if e.promptOpponentCounterTrap(playerID, TriggerOnUseItem, card, useData, resolveItem) {
		return nil
	}
	resolveItem()
	return nil
}

func (e *Engine) startSpellScrollCast(playerID int, scroll *CardInstance, target SpellTarget) {
	ps := e.State.Players[playerID]
	boostSkills := []*CardInstance{}
	totalPower := e.effectiveSpellPower(playerID, scroll, boostSkills, target)
	powerSources := e.spellPowerSources(playerID, scroll, boostSkills, totalPower, target)
	e.consumeNextSpellPowerBonuses(ps, scroll)

	if ps.SpellsCastThisTurn == nil {
		ps.SpellsCastThisTurn = make(map[string]int)
	}
	ps.SpellsCastThisTurn[scroll.Card.Category]++
	spellCastData := map[string]any{
		"cast_player":  playerID,
		"attacker":     playerID,
		"skill":        cardToInfo(scroll),
		"target":       target,
		"power":        totalPower,
		"boost_count":  0,
		"is_sorcery":   false,
		"spell_scroll": true,
	}
	e.emit(GameEvent{Type: "spell_cast", Player: -1, Data: spellCastData})
	e.triggerEffects(TriggerOnSpellCast, scroll, nil, spellCastData)

	e.State.PendingSpell = &SpellCast{
		AttackerID:   playerID,
		Skill:        scroll,
		Target:       target,
		TotalPower:   totalPower,
		PowerSources: powerSources,
		BoostSkills:  boostSkills,
	}
	openDefenseWindow := func() {
		if e.State.PendingSpell == nil {
			return
		}
		e.State.ResumePhase = PhaseDefenseWindow
		e.State.Phase = PhaseDefenseWindow
		e.emit(GameEvent{Type: "defense_window", Player: 1 - playerID, Data: map[string]any{"timeout": 30}})
	}
	continueSpell := openDefenseWindow
	if !e.spellAllowsDefense(playerID, scroll, target) {
		continueSpell = func() {
			e.resolvePendingSpellHit()
		}
	}
	if e.triggerSpellCastFieldEffectsWithContinuation(playerID, scroll, spellCastData, continueSpell) {
		if e.spellAllowsDefense(playerID, scroll, target) {
			e.State.ResumePhase = PhaseDefenseWindow
		}
		return
	}
	continueSpell()
}

// handlePlaceTerrain handles placing a terrain card (地形牌) on the battlefield
func (e *Engine) handlePlaceTerrain(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseMain {
		return fmt.Errorf("not in main phase")
	}
	if e.State.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}

	instanceID, _ := action.Data["instance_id"].(string)
	colF, _ := action.Data["col"].(float64)
	rowF, _ := action.Data["row"].(float64)
	col := int(colF)
	row := int(rowF)

	ps := e.State.Players[playerID]

	// Find card in hand
	card, handIdx := ps.FindHandCard(instanceID)
	if card == nil {
		return fmt.Errorf("card not found in hand")
	}

	// Must be an item with terrain keyword
	if !card.Card.IsItem() {
		return fmt.Errorf("card is not an item")
	}
	if !cards.IsTerrain(card.Card.Number) {
		return fmt.Errorf("card is not a terrain")
	}

	// Check position
	pos := Position{Col: col, Row: row}
	if !pos.Valid() {
		return fmt.Errorf("invalid position")
	}
	if ps.Terrain[col][row] != nil {
		return fmt.Errorf("position already has terrain")
	}

	// Check cost
	cost := e.effectiveCardPlayCost(ps, card)
	if !e.canPayCost(ps, cost) {
		return fmt.Errorf("not enough elements")
	}

	// Pay cost and place
	if !e.payCostForAction(ps, cost, action) {
		return fmt.Errorf("invalid payment")
	}
	e.notifyCardPlayCostPaid(ps, card)
	ps.RemoveFromHand(handIdx)
	card.Position = &Position{Col: col, Row: row}
	card.EnterTurn = e.State.TurnNumber
	ps.Terrain[col][row] = card

	e.emit(GameEvent{
		Type:   "place_terrain",
		Player: -1,
		Data: map[string]any{
			"player":   playerID,
			"card":     cardToInfo(card),
			"position": pos,
			"elements": ps.Elements,
		},
	})

	// Trigger 入场 (on enter) effects for the terrain
	e.triggerEffects(TriggerOnEnter, card, nil, nil)

	e.checkWinCondition()
	return nil
}

// handleUseAbility handles using a card's activated ability (回合技/绝技)
func (e *Engine) handleUseAbility(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseMain {
		return fmt.Errorf("not in main phase")
	}
	if e.State.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}

	instanceID, _ := action.Data["instance_id"].(string)
	abilityType, _ := action.Data["ability_type"].(string) // "per_turn" or "ultimate"
	targetID, _ := action.Data["target_id"].(string)

	ps := e.State.Players[playerID]

	// Find the card with the ability
	card := e.findCardOnField(ps, instanceID)
	if card == nil {
		return fmt.Errorf("card not found on field")
	}

	if e.hasEffectiveStatus(card, StatusPetrify) {
		return fmt.Errorf("card is petrified")
	}
	if e.hasEffectiveStatus(card, StatusStun) {
		return fmt.Errorf("card is stunned")
	}

	var trigger EffectTrigger
	if abilityType == "ultimate" {
		trigger = TriggerUltimate
		if !cardHasActiveUltimate(card) {
			return fmt.Errorf("card has no active ultimate ability")
		}
		if card.UltimateUsed {
			return fmt.Errorf("ultimate already used")
		}
		if err := e.validateUltimatePreconditions(card); err != nil {
			return err
		}
	} else {
		trigger = TriggerPerTurn
		if !cardHasActivePerTurn(card) {
			return fmt.Errorf("card has no active per-turn ability")
		}
		// Check if already used this turn (回合技 limit)
		maxUses := perTurnLimit(card)
		if card.UsedThisTurn >= maxUses {
			return fmt.Errorf("ability already used this turn")
		}
	}

	// Find target if specified
	var target *CardInstance
	if targetID != "" {
		// Search both player fields for target
		for i := 0; i < 2; i++ {
			t := e.findCardOnField(e.State.Players[i], targetID)
			if t == nil {
				t = e.findUnitOnGrid(e.State.Players[i], targetID)
			}
			if t != nil {
				target = t
				break
			}
		}
	}

	// Execute the ability
	effects := globalRegistry.GetEffects(card.Card.Number, trigger)
	if len(effects) == 0 {
		return fmt.Errorf("card has no %s ability", abilityType)
	}

	ctx := &EffectContext{
		Engine:     e,
		Source:     card,
		Target:     target,
		PlayerID:   playerID,
		OpponentID: 1 - playerID,
	}

	for _, eff := range effects {
		if eff.IsActive {
			if err := eff.Handler(ctx); err != nil {
				return err
			}
		}
	}

	if abilityType == "ultimate" {
		card.UltimateUsed = true
	} else {
		card.UsedThisTurn++
	}

	e.emit(GameEvent{
		Type:   "ability_used",
		Player: -1,
		Data: map[string]any{
			"player":  playerID,
			"card":    cardToInfo(card),
			"ability": abilityType,
		},
	})

	e.checkWinCondition()
	return nil
}

func (e *Engine) validateUltimatePreconditions(card *CardInstance) error {
	if card == nil || card.Card == nil {
		return nil
	}
	switch card.Card.Number {
	case "4211001":
		if len(e.State.Players[card.OwnerID].Hand) == 0 {
			return fmt.Errorf("Bartel ultimate requires a hand card")
		}
	case "4311001":
		if len(e.friendlyHandCards(card.OwnerID, func(candidate *CardInstance) bool {
			return candidate.Card.Category == model.ElementAir
		})) < 2 {
			return fmt.Errorf("Su ultimate requires two air cards in hand")
		}
		if len(e.enemyUnits(card.OwnerID, true, nil)) == 0 {
			return fmt.Errorf("Su ultimate requires an enemy target")
		}
	}
	return nil
}

// handleResolveAction handles the player's response to a pending action
func (e *Engine) handleResolveAction(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseWaitingAction {
		return fmt.Errorf("no pending action")
	}
	pa := e.State.PendingAction
	if pa == nil {
		return fmt.Errorf("no pending action")
	}
	if pa.PlayerID != playerID {
		return fmt.Errorf("not your pending action")
	}

	selectedRaw, _ := action.Data["selected"].([]any)
	var selected []string
	for _, s := range selectedRaw {
		if str, ok := s.(string); ok {
			selected = append(selected, str)
		}
	}

	if len(selected) < pa.MinSelect {
		return fmt.Errorf("must select at least %d", pa.MinSelect)
	}
	if len(selected) > pa.MaxSelect {
		return fmt.Errorf("can select at most %d", pa.MaxSelect)
	}
	allowed := make(map[string]bool, len(pa.Candidates))
	selectable := make(map[string]bool, len(pa.Candidates))
	for _, candidate := range pa.Candidates {
		if id, ok := candidate["instance_id"].(string); ok && id != "" {
			allowed[id] = true
			selectable[id] = candidate["can_select"] != false
		}
	}
	seen := make(map[string]bool, len(selected))
	for _, id := range selected {
		if !allowed[id] {
			return fmt.Errorf("invalid selection")
		}
		if !selectable[id] {
			return fmt.Errorf("invalid selection")
		}
		if seen[id] {
			return fmt.Errorf("duplicate selection")
		}
		seen[id] = true
	}

	// Execute callback
	callback := pa.Callback
	callbackData := pa.CallbackData
	callbackErr := pa.CallbackErr
	data := action.Data
	e.State.PendingAction = nil
	e.State.Phase = e.State.ResumePhase

	if callbackErr != nil {
		if err := callbackErr(selected, data); err != nil {
			e.State.PendingAction = pa
			e.State.Phase = PhaseWaitingAction
			return err
		}
	} else if callbackData != nil {
		callbackData(selected, data)
	} else if callback != nil {
		callback(selected)
	}

	e.emitPendingActionCleared(pa)
	e.advancePendingActionQueue()
	if e.State.PendingAction == nil && e.State.Phase == PhaseDefenseWindow && e.State.PendingSpell == nil {
		e.State.Phase = PhaseMain
	}

	e.checkWinCondition()
	return nil
}

// SetPendingAction sets a pending player action and pauses the game
func (e *Engine) SetPendingAction(playerID int, actionType string, prompt string, candidates []map[string]any, minSelect, maxSelect int, callback func([]string)) {
	e.setPendingAction(playerID, actionType, prompt, candidates, minSelect, maxSelect, callback, nil)
}

func (e *Engine) SetPendingActionWithData(playerID int, actionType string, prompt string, candidates []map[string]any, minSelect, maxSelect int, callback func([]string, map[string]any)) {
	e.setPendingAction(playerID, actionType, prompt, candidates, minSelect, maxSelect, nil, callback)
}

func (e *Engine) SetPendingActionWithError(playerID int, actionType string, prompt string, candidates []map[string]any, minSelect, maxSelect int, cost map[string]int, canOverexert bool, callback func([]string, map[string]any) error) {
	e.setPendingActionWithOptions(playerID, actionType, prompt, candidates, minSelect, maxSelect, cost, canOverexert, nil, nil, callback, nil)
}

func (e *Engine) SetPendingActionWithErrorAndContext(playerID int, actionType string, prompt string, candidates []map[string]any, minSelect, maxSelect int, cost map[string]int, canOverexert bool, context map[string]any, callback func([]string, map[string]any) error) {
	e.setPendingActionWithOptions(playerID, actionType, prompt, candidates, minSelect, maxSelect, cost, canOverexert, nil, nil, callback, context)
}

func (e *Engine) setPendingAction(playerID int, actionType string, prompt string, candidates []map[string]any, minSelect, maxSelect int, callback func([]string), callbackData func([]string, map[string]any)) {
	e.setPendingActionWithOptions(playerID, actionType, prompt, candidates, minSelect, maxSelect, nil, false, callback, callbackData, nil, nil)
}

func (e *Engine) setPendingActionWithOptions(playerID int, actionType string, prompt string, candidates []map[string]any, minSelect, maxSelect int, cost map[string]int, canOverexert bool, callback func([]string), callbackData func([]string, map[string]any), callbackErr func([]string, map[string]any) error, context map[string]any) {
	if minSelect > 0 && len(candidates) == 0 {
		return
	}
	resumePhase := e.State.Phase
	if e.State.PendingAction != nil {
		resumePhase = e.State.ResumePhase
	}
	action := &PendingAction{
		Type:         actionType,
		PlayerID:     playerID,
		Prompt:       prompt,
		Candidates:   candidates,
		MinSelect:    minSelect,
		MaxSelect:    maxSelect,
		Context:      context,
		Cost:         cost,
		CanOverexert: canOverexert,
		Callback:     callback,
		CallbackData: callbackData,
		CallbackErr:  callbackErr,
	}
	if e.State.PendingAction != nil {
		e.State.PendingActionQueue = append(e.State.PendingActionQueue, action)
		return
	}
	e.activatePendingAction(action, resumePhase)
}

func (e *Engine) activatePendingAction(action *PendingAction, resumePhase GamePhase) {
	if action == nil {
		return
	}
	e.State.ResumePhase = resumePhase
	e.State.Phase = PhaseWaitingAction
	e.State.PendingAction = action
	data := map[string]any{
		"type":       action.Type,
		"player_id":  action.PlayerID,
		"prompt":     action.Prompt,
		"candidates": action.Candidates,
		"min_select": action.MinSelect,
		"max_select": action.MaxSelect,
	}
	if action.Context != nil {
		data["context"] = action.Context
	}
	if action.Cost != nil {
		data["cost"] = action.Cost
	}
	if action.CanOverexert {
		data["can_overexert"] = true
	}
	e.emit(GameEvent{Type: "pending_action", Player: action.PlayerID, Data: data})
}

func (e *Engine) advancePendingActionQueue() bool {
	if e.State.PendingAction != nil || len(e.State.PendingActionQueue) == 0 {
		return false
	}
	next := e.State.PendingActionQueue[0]
	e.State.PendingActionQueue = e.State.PendingActionQueue[1:]
	e.activatePendingAction(next, e.State.Phase)
	return true
}

func (e *Engine) emitPendingActionCleared(action *PendingAction) {
	if action == nil {
		return
	}
	e.emit(GameEvent{Type: "pending_action_cleared", Player: action.PlayerID, Data: map[string]any{
		"type":      action.Type,
		"player_id": action.PlayerID,
	}})
}

// handleEndTurn handles ending the current turn
func (e *Engine) handleEndTurn(playerID int, action ActionMessage) error {
	if e.State.Phase != PhaseMain {
		return fmt.Errorf("not in main phase")
	}
	if e.State.CurrentTurn != playerID {
		return fmt.Errorf("not your turn")
	}

	e.endTurn()
	return nil
}

// endTurn processes end-of-turn effects and switches to next player
func (e *Engine) endTurn() {
	ps := e.State.Players[e.State.CurrentTurn]

	// Discard to hand limit
	handLimit := e.handLimitForPlayer(ps)
	if len(ps.Hand) > handLimit {
		discardCount := len(ps.Hand) - handLimit
		// Build candidates from hand cards
		candidates := make([]map[string]any, len(ps.Hand))
		for i, c := range ps.Hand {
			candidates[i] = cardToInfo(c)
		}
		currentTurn := e.State.CurrentTurn
		e.SetPendingAction(currentTurn, "discard",
			fmt.Sprintf("弃牌至手牌上限（需弃%d张）", discardCount),
			candidates, discardCount, discardCount,
			func(selected []string) {
				// Discard selected cards
				toDiscard := make(map[string]bool)
				for _, id := range selected {
					toDiscard[id] = true
				}
				remaining := make([]*CardInstance, 0, len(ps.Hand)-len(selected))
				for _, c := range ps.Hand {
					if toDiscard[c.InstanceID] {
						e.discardHandCardToGraveyard(currentTurn, c)
					} else {
						remaining = append(remaining, c)
					}
				}
				ps.Hand = remaining
				// Continue end turn processing
				e.finishEndTurn(ps)
			})
		return // Wait for player to choose
	}

	e.finishEndTurn(ps)
}

// finishEndTurn completes end-of-turn processing (after optional discard)
func (e *Engine) finishEndTurn(ps *PlayerState) {
	// Trigger 回合结束 effects for all cards on the current player's field
	allCards := e.getAllFieldCards(ps)
	for _, card := range allCards {
		e.triggerEffects(TriggerOnTurnEnd, card, nil, nil)
	}
	e.triggerFieldEffectsWithData(TriggerOnTurnEnd, 1-ps.PlayerID, nil, map[string]any{"ended_player": ps.PlayerID})
	e.applyOpponentTurnEndTemporaryModifiers(ps.PlayerID)
	e.discardMarkedEndOfTurnCards(ps)
	e.applyLoadGainAtTurnEnd(ps)

	e.clearExpiredTemporaryModifiers(ps.PlayerID)
	e.processAbilityDurations(ps)

	// Remove 临时 (temporary) units before the cleanup/reset steps.
	e.HandleTemporaryUnits(ps)

	// Discard phase has already happened above. The cleanup order is:
	// reset cards first, then settle marks. A skill with 冷却1 therefore remains
	// horizontal through the next turn, because 冷却 blocks this reset before it
	// is removed by mark settlement.
	e.resetCards(ps)

	// Process status marks (点燃, 冻结, 冷却, etc.) after reset.
	e.processEndOfTurnStatuses(ps)

	// Decay 护盾 and 隐蔽 as part of mark settlement.
	e.HandleShieldDecay(ps)

	// 精通 is a card-instance mark. It advances during the unified mark
	// settlement step, after reset has already considered 冷却/冻结/etc.
	e.settleMastery(ps)

	// Clear elements
	for elem := range ps.Elements {
		ps.Elements[elem] = 0
	}

	e.emit(GameEvent{
		Type:   "turn_end",
		Player: -1,
		Data: map[string]any{
			"player": e.State.CurrentTurn,
		},
	})

	// Switch turns
	if e.State.IsFirstTurn && e.State.CurrentTurn == e.State.FirstPlayer {
		e.State.IsFirstTurn = false
	}
	e.State.CurrentTurn = 1 - e.State.CurrentTurn
	if e.State.CurrentTurn == e.State.FirstPlayer {
		e.State.TurnNumber++
	}

	if e.State.Phase != PhaseGameOver {
		e.startTurn()
	}
}

func (e *Engine) processAbilityDurations(ps *PlayerState) {
	changedRedMoon := false
	for _, card := range e.getAllFieldCards(ps) {
		if card == nil || card.Statuses[StatusAbilityDuration] <= 0 {
			continue
		}
		wasRedMoon := card.Card != nil && card.Card.Number == "3611101"
		card.Statuses[StatusAbilityDuration]--
		if card.Statuses[StatusAbilityDuration] <= 0 {
			delete(card.Statuses, StatusAbilityDuration)
			if wasRedMoon {
				changedRedMoon = true
			}
		}
	}
	if changedRedMoon {
		e.updateRedMoonTransformations(ps.PlayerID)
	}
}

// processEndOfTurnStatuses processes status marks at end of turn
func (e *Engine) processEndOfTurnStatuses(ps *PlayerState) {
	allCards := e.getAllFieldCards(ps)
	redMoonPetrifyChanged := false

	for _, card := range allCards {
		// 点燃: remove 1 stack, deal 1 fire damage
		if card.Statuses[StatusBurn] > 0 {
			effective := e.hasEffectiveStatus(card, StatusBurn)
			card.Statuses[StatusBurn]--
			if effective {
				e.dealDamageWithExtra(card, 1, ps.PlayerID, map[string]any{"status_damage": StatusBurn})
			}
		}
		// 冻结: remove 1 stack
		if card.Statuses[StatusFreeze] > 0 {
			card.Statuses[StatusFreeze]--
		}
		// 眩晕: remove 1 stack
		if card.Statuses[StatusStun] > 0 {
			card.Statuses[StatusStun]--
		}
		// 石化: remove 1 stack
		if card.Statuses[StatusPetrify] > 0 {
			wasRedMoon := card.Card != nil && card.Card.Number == "3611101"
			card.Statuses[StatusPetrify]--
			if wasRedMoon && card.Statuses[StatusPetrify] <= 0 {
				redMoonPetrifyChanged = true
			}
		}
		// 冷却: remove 1 stack
		if card.Statuses[StatusCooldown] > 0 {
			card.Statuses[StatusCooldown]--
		}
	}

	// 虚弱 is on skills, handled separately
	for i := range ps.Skills {
		if ps.Skills[i] != nil && ps.Skills[i].Statuses[StatusWeaken] > 0 {
			ps.Skills[i].Statuses[StatusWeaken]--
		}
		if ps.Skills[i] != nil && ps.Skills[i].Statuses[StatusSeal] > 0 {
			ps.Skills[i].Statuses[StatusSeal]--
		}
	}
	if redMoonPetrifyChanged {
		e.refreshRedMoonState(ps.PlayerID)
	}
}

// checkWinCondition checks if the game is over
func (e *Engine) checkWinCondition() {
	if e.State.Phase == PhaseGameOver {
		return
	}
	if e.resolutionDepth > 0 || e.resolvingDeaths || len(e.deathQueue) > 0 {
		return
	}

	p0Dead := e.State.Players[0].Hero != nil && e.State.Players[0].Hero.CurrentLife <= 0
	p1Dead := e.State.Players[1].Hero != nil && e.State.Players[1].Hero.CurrentLife <= 0
	switch {
	case p0Dead && p1Dead:
		e.State.Winner = -2
		e.State.Phase = PhaseGameOver
		e.clearPendingForGameOver()
		e.emit(GameEvent{
			Type:   "game_over",
			Player: -1,
			Data: map[string]any{
				"winner": e.State.Winner,
				"reason": "both_heroes_killed",
			},
		})
	case p0Dead:
		e.State.Winner = 1
		e.State.Phase = PhaseGameOver
		e.clearPendingForGameOver()
		e.emit(GameEvent{
			Type:   "game_over",
			Player: -1,
			Data: map[string]any{
				"winner": e.State.Winner,
				"reason": "hero_killed",
			},
		})
	case p1Dead:
		e.State.Winner = 0
		e.State.Phase = PhaseGameOver
		e.clearPendingForGameOver()
		e.emit(GameEvent{
			Type:   "game_over",
			Player: -1,
			Data: map[string]any{
				"winner": e.State.Winner,
				"reason": "hero_killed",
			},
		})
	}
}

func (e *Engine) clearPendingForGameOver() {
	e.State.PendingAction = nil
	e.State.PendingSpell = nil
	e.State.ResumePhase = PhaseGameOver
}

func (e *Engine) payCostForAction(ps *PlayerState, cost map[string]int, action ActionMessage) bool {
	if payment := paymentFromAction(action); payment != nil {
		if !validateElementPaymentWithOptions(ps.Elements, cost, payment, e.playerHasLightWildcard(ps)) {
			return false
		}
		for elem, amount := range payment {
			ps.Elements[elem] -= amount
		}
		return true
	}
	payment, ok := calculateElementPaymentWithOptions(ps.Elements, cost, e.playerHasLightWildcard(ps))
	if !ok {
		return false
	}
	for elem, amount := range payment {
		ps.Elements[elem] -= amount
	}
	return true
}

func (e *Engine) payCostForCardAction(ps *PlayerState, card *CardInstance, strictCost map[string]int, totalCost map[string]int, purpose paymentPurpose, action ActionMessage) bool {
	payment := paymentFromAction(action)
	if payment == nil {
		var ok bool
		payment, ok = calculateElementPaymentWithOptions(ps.Elements, totalCost, e.playerHasLightWildcard(ps))
		if !ok {
			return false
		}
	} else if !validateElementPaymentWithOptions(ps.Elements, totalCost, payment, e.playerHasLightWildcard(ps)) {
		return false
	}
	if !strictPaymentSatisfied(card, purpose, strictCost, payment) {
		return false
	}
	for elem, amount := range payment {
		ps.Elements[elem] -= amount
	}
	return true
}

func (e *Engine) canPayCostForAction(ps *PlayerState, cost map[string]int, action ActionMessage) bool {
	if payment := paymentFromAction(action); payment != nil {
		return validateElementPaymentWithOptions(ps.Elements, cost, payment, e.playerHasLightWildcard(ps))
	}
	return e.canPayCost(ps, cost)
}

func (e *Engine) canPayCostForCardAction(ps *PlayerState, card *CardInstance, strictCost map[string]int, totalCost map[string]int, purpose paymentPurpose, action ActionMessage) bool {
	payment := paymentFromAction(action)
	if payment == nil {
		var ok bool
		payment, ok = calculateElementPaymentWithOptions(ps.Elements, totalCost, e.playerHasLightWildcard(ps))
		if !ok {
			return false
		}
	} else if !validateElementPaymentWithOptions(ps.Elements, totalCost, payment, e.playerHasLightWildcard(ps)) {
		return false
	}
	return strictPaymentSatisfied(card, purpose, strictCost, payment)
}

func (e *Engine) canPayCost(ps *PlayerState, cost map[string]int) bool {
	_, ok := calculateElementPaymentWithOptions(ps.Elements, cost, e.playerHasLightWildcard(ps))
	return ok
}

func (e *Engine) playerHasLightWildcard(ps *PlayerState) bool {
	if e == nil || ps == nil {
		return false
	}
	for _, card := range e.getAllFieldCards(ps) {
		if card != nil && card.Card != nil && card.Card.Number == "1521007" && !e.hasEffectiveStatus(card, StatusPetrify) {
			return true
		}
	}
	return false
}

func paymentFromAction(action ActionMessage) map[string]int {
	raw, ok := action.Data["payment"]
	if !ok || raw == nil {
		return nil
	}
	payment := make(map[string]int)
	switch values := raw.(type) {
	case map[string]any:
		for elem, value := range values {
			switch amount := value.(type) {
			case float64:
				payment[elem] = int(amount)
			case int:
				payment[elem] = amount
			}
		}
	case map[string]int:
		for elem, amount := range values {
			payment[elem] = amount
		}
	}
	return payment
}

// GetStateForPlayer returns a filtered game state visible to the specified player
func (e *Engine) GetStateForPlayer(playerID int) map[string]any {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.State
	opponentID := 1 - playerID
	ps := state.Players[playerID]
	op := state.Players[opponentID]

	return map[string]any{
		"game_id":      state.GameID,
		"phase":        state.Phase.String(),
		"current_turn": state.CurrentTurn,
		"first_player": state.FirstPlayer,
		"turn_order":   map[string]string{"you": turnOrderLabel(playerID, state.FirstPlayer), "opponent": turnOrderLabel(opponentID, state.FirstPlayer)},
		"turn_number":  state.TurnNumber,
		"winner":       state.Winner,
		"you":          e.playerStateToInfo(ps, true),
		"opponent":     e.playerStateToInfo(op, false),
		"pending_spell": func() any {
			if state.PendingSpell != nil {
				return map[string]any{
					"attacker":      state.PendingSpell.AttackerID,
					"skill":         cardToInfo(state.PendingSpell.Skill),
					"target":        state.PendingSpell.Target,
					"power":         state.PendingSpell.TotalPower,
					"power_sources": state.PendingSpell.PowerSources,
					"boost_skills":  cardsToInfo(state.PendingSpell.BoostSkills),
					"extra_targets": state.PendingSpell.ExtraTargets,
				}
			}
			return nil
		}(),
		"pending_action": func() any {
			if state.PendingAction != nil && state.PendingAction.PlayerID == playerID {
				return map[string]any{
					"type":          state.PendingAction.Type,
					"player_id":     state.PendingAction.PlayerID,
					"prompt":        state.PendingAction.Prompt,
					"candidates":    state.PendingAction.Candidates,
					"min_select":    state.PendingAction.MinSelect,
					"max_select":    state.PendingAction.MaxSelect,
					"context":       state.PendingAction.Context,
					"cost":          state.PendingAction.Cost,
					"can_overexert": state.PendingAction.CanOverexert,
				}
			}
			return nil
		}(),
	}
}

// GetStateForSpectator returns the public game state without either player's
// hidden hand, skill-pool, or deck contents.
func (e *Engine) GetStateForSpectator() map[string]any {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.State
	return map[string]any{
		"game_id":      state.GameID,
		"phase":        state.Phase.String(),
		"current_turn": state.CurrentTurn,
		"first_player": state.FirstPlayer,
		"turn_order":   map[string]string{"you": turnOrderLabel(0, state.FirstPlayer), "opponent": turnOrderLabel(1, state.FirstPlayer)},
		"turn_number":  state.TurnNumber,
		"winner":       state.Winner,
		"is_spectator": true,
		"you":          e.playerStateToInfo(state.Players[0], false),
		"opponent":     e.playerStateToInfo(state.Players[1], false),
		"pending_spell": func() any {
			if state.PendingSpell != nil {
				return map[string]any{
					"attacker":      state.PendingSpell.AttackerID,
					"skill":         cardToInfo(state.PendingSpell.Skill),
					"target":        state.PendingSpell.Target,
					"power":         state.PendingSpell.TotalPower,
					"power_sources": state.PendingSpell.PowerSources,
					"boost_skills":  cardsToInfo(state.PendingSpell.BoostSkills),
					"extra_targets": state.PendingSpell.ExtraTargets,
				}
			}
			return nil
		}(),
		"pending_action": func() any {
			if state.PendingAction != nil {
				return map[string]any{
					"type":      state.PendingAction.Type,
					"player_id": state.PendingAction.PlayerID,
					"prompt":    state.PendingAction.Prompt,
				}
			}
			return nil
		}(),
	}
}

// Helper functions

func (e *Engine) findCardOnField(ps *PlayerState, instanceID string) *CardInstance {
	if ps.Hero != nil && ps.Hero.InstanceID == instanceID {
		return ps.Hero
	}
	// Check units
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if ps.Units[col][row] != nil && ps.Units[col][row].InstanceID == instanceID {
				return ps.Units[col][row]
			}
		}
	}
	// Check terrain
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if ps.Terrain[col][row] != nil && ps.Terrain[col][row].InstanceID == instanceID {
				return ps.Terrain[col][row]
			}
		}
	}
	// Check equipment
	for i := range ps.Equipment {
		if ps.Equipment[i] != nil && ps.Equipment[i].InstanceID == instanceID {
			return ps.Equipment[i]
		}
	}
	return nil
}

func (e *Engine) findUnitOnGrid(ps *PlayerState, instanceID string) *CardInstance {
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if ps.Units[col][row] != nil && ps.Units[col][row].InstanceID == instanceID {
				return ps.Units[col][row]
			}
		}
	}
	return nil
}

func (e *Engine) findEquipment(ps *PlayerState, instanceID string) *CardInstance {
	for i := range ps.Equipment {
		if ps.Equipment[i] != nil && ps.Equipment[i].InstanceID == instanceID {
			return ps.Equipment[i]
		}
	}
	return nil
}

func (e *Engine) findSkill(ps *PlayerState, instanceID string) *CardInstance {
	for i := range ps.Skills {
		if ps.Skills[i] != nil && ps.Skills[i].InstanceID == instanceID {
			return ps.Skills[i]
		}
	}
	for _, card := range e.getAllFieldCards(ps) {
		if card == nil {
			continue
		}
		for _, skill := range card.BoundSkills {
			if skill != nil && skill.InstanceID == instanceID {
				return skill
			}
		}
	}
	return nil
}

func (e *Engine) findReactionCard(ps *PlayerState, instanceID string) *CardInstance {
	if skill := e.findSkill(ps, instanceID); skill != nil {
		return skill
	}
	for _, equipment := range ps.Equipment {
		if equipment != nil && equipment.InstanceID == instanceID {
			return equipment
		}
	}
	return nil
}

func (e *Engine) getAllFieldCards(ps *PlayerState) []*CardInstance {
	var cards []*CardInstance
	if ps.Hero != nil {
		cards = append(cards, ps.Hero)
	}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if ps.Units[col][row] != nil && ps.Units[col][row] != ps.Hero {
				cards = append(cards, ps.Units[col][row])
			}
		}
	}
	for i := range ps.Skills {
		if ps.Skills[i] != nil {
			cards = append(cards, ps.Skills[i])
		}
	}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if ps.Terrain[col][row] != nil {
				cards = append(cards, ps.Terrain[col][row])
			}
		}
	}
	for i := range ps.Equipment {
		if ps.Equipment[i] != nil {
			cards = append(cards, ps.Equipment[i])
		}
	}
	return cards
}

func shuffleDeck(deck []*CardInstance) {
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
}

func cardToInfo(ci *CardInstance) map[string]any {
	if ci == nil {
		return nil
	}
	info := map[string]any{
		"instance_id":            ci.InstanceID,
		"owner":                  ci.OwnerID,
		"number":                 ci.Card.Number,
		"name":                   ci.Card.Name,
		"type":                   ci.Card.Type,
		"category":               ci.Card.Category,
		"tag":                    ci.Card.Tag,
		"description":            ci.Card.Description,
		"attack":                 ci.Card.Attack + ci.AttackBonus,
		"life":                   maxLife(ci),
		"power":                  ci.Card.Power + ci.PowerBonus,
		"duration":               ci.Card.Duration,
		"elements_cost":          ci.Card.ElementsCost,
		"elements_gain":          effectiveElementsGain(ci),
		"elements_expense":       ci.Card.ElementsExpense,
		"current_life":           ci.CurrentLife,
		"current_attack":         ci.CurrentAttack,
		"is_horizontal":          ci.IsHorizontal,
		"is_terrain":             cards.IsTerrain(ci.Card.Number),
		"is_companion":           ci.Card.IsCompanion(),
		"is_consumable":          cards.IsConsumable(ci.Card.Number),
		"is_equipment":           cards.IsEquipment(ci.Card.Number),
		"is_weapon":              cards.IsWeapon(ci.Card.Number),
		"has_taunt":              cardHasTaunt(ci),
		"has_global_spell_range": cardHasActiveGlobalSpellRange(ci),
		"is_counter_trap":        isCounterTrapCard(ci.Card.Number),
		"is_set_counter":         ci.IsSetCounter,
		"statuses":               ci.Statuses,
		"position":               ci.Position,
		"output_path":            ci.Card.OutputPath,
		"used_this_turn":         ci.UsedThisTurn,
		"ultimate_used":          ci.UltimateUsed,
		"uses_remaining":         ci.UsesRemaining,
	}
	addCardEffectMetadata(info, ci.Card)
	if len(ci.BoundSkills) > 0 {
		info["bound_skills"] = cardsToInfo(ci.BoundSkills)
	}
	if len(ci.UnderCards) > 0 {
		info["under_cards"] = cardsToInfo(ci.UnderCards)
	}
	if attached := attachedBehaviorsInfo(ci); len(attached) > 0 {
		info["attached_behaviors"] = attached
	}

	hasPerTurn := cardHasActivePerTurn(ci)
	hasUltimate := cardHasActiveUltimate(ci)
	info["has_per_turn"] = hasPerTurn
	info["has_prayer"] = cardHasActivePrayer(ci)
	info["has_ultimate"] = hasUltimate
	if cardHasActiveSpellReaction(ci) {
		info["can_react"] = true
	}

	if hasPerTurn {
		info["per_turn_limit"] = perTurnLimit(ci)
		info["per_turn_label"] = "回合技"
		behavior := behaviorForNumber(ci.Card.Number)
		if labeler, ok := behavior.(PerTurnLabelBehavior); ok {
			if label := labeler.PerTurnLabel(ci); label != "" {
				info["per_turn_label"] = label
			}
		}
	}
	if requirement := summonDevourRequirement(ci); len(requirement) > 0 {
		info["devour_requirement"] = requirement
	}

	// Mark spell-like skills and spell scrolls.
	if isSpellLikeCard(ci.Card) {
		info["is_defense_only"] = isDefenseOnlySkill(ci.Card)
		info["is_sorcery"] = isSorcerySkill(ci.Card)
		info["needs_target"] = skillNeedsTargetInstance(ci)
		info["has_pierce"] = cardHasPierce(ci)
		if friendly, ok := behaviorForNumber(ci.Card.Number).(FriendlySpellTargetBehavior); ok && friendly.HasActiveFriendlySpellTarget(ci) && friendly.AllowsFriendlySpellTarget() {
			info["allows_friendly_target"] = true
		}
		info["can_attack"] = canUseSkillForPurpose(ci.Card, skillPurposeAttack)
		info["can_defend"] = canUseSkillForPurpose(ci.Card, skillPurposeDefend)
		info["can_attack_boost"] = canUseSkillForPurpose(ci.Card, skillPurposeAttackBoost)
		info["can_defense_boost"] = canUseSkillForPurpose(ci.Card, skillPurposeDefenseBoost)
		info["can_react"] = cardHasActiveSpellReaction(ci)
		info["can_boost"] = info["can_attack_boost"]
		info["spell_area"] = spellArea(ci)
	}

	return info
}

func cardsToInfo(cards []*CardInstance) []map[string]any {
	result := make([]map[string]any, len(cards))
	for i, c := range cards {
		result[i] = cardToInfo(c)
	}
	return result
}

func (e *Engine) cardToInfo(ci *CardInstance) map[string]any {
	info := cardToInfo(ci)
	if info == nil {
		return nil
	}
	info["elements_gain"] = e.effectiveElementsGain(ci)
	if len(ci.BoundSkills) > 0 {
		info["bound_skills"] = e.cardsToInfo(ci.BoundSkills)
	}
	if len(ci.UnderCards) > 0 {
		info["under_cards"] = e.cardsToInfo(ci.UnderCards)
	}
	return info
}

func (e *Engine) cardsToInfo(cards []*CardInstance) []map[string]any {
	result := make([]map[string]any, len(cards))
	for i, c := range cards {
		result[i] = e.cardToInfo(c)
	}
	return result
}

func deckSummaryToInfo(deck []*CardInstance) []map[string]any {
	type summary struct {
		card  *model.Card
		count int
	}

	byNumber := map[string]*summary{}
	for _, ci := range deck {
		if ci == nil || ci.Card == nil {
			continue
		}
		number := ci.Card.Number
		entry := byNumber[number]
		if entry == nil {
			entry = &summary{card: ci.Card}
			byNumber[number] = entry
		}
		entry.count++
	}

	numbers := make([]string, 0, len(byNumber))
	for number := range byNumber {
		numbers = append(numbers, number)
	}
	sort.Strings(numbers)

	result := make([]map[string]any, 0, len(numbers))
	for _, number := range numbers {
		entry := byNumber[number]
		card := entry.card
		info := map[string]any{
			"number":           card.Number,
			"name":             card.Name,
			"type":             card.Type,
			"category":         card.Category,
			"tag":              card.Tag,
			"description":      card.Description,
			"attack":           card.Attack,
			"life":             card.Life,
			"power":            card.Power,
			"duration":         card.Duration,
			"elements_cost":    card.ElementsCost,
			"elements_gain":    card.ElementsGain,
			"elements_expense": card.ElementsExpense,
			"output_path":      card.OutputPath,
			"count":            entry.count,
			"is_terrain":       cards.IsTerrain(card.Number),
			"is_consumable":    cards.IsConsumable(card.Number),
			"is_equipment":     cards.IsEquipment(card.Number),
			"is_weapon":        cards.IsWeapon(card.Number),
		}
		addCardEffectMetadata(info, card)
		result = append(result, info)
	}

	return result
}

func addCardEffectMetadata(info map[string]any, card *model.Card) {
	if card == nil {
		return
	}
	if len(card.EffectCategories) > 0 {
		info["effect_categories"] = card.EffectCategories
	}
	if len(card.EffectOptionality) > 0 {
		info["effect_optionality"] = card.EffectOptionality
	}
}

func turnOrderLabel(playerID int, firstPlayer int) string {
	if playerID == firstPlayer {
		return "先手"
	}
	return "后手"
}

func (e *Engine) playerStateToInfo(ps *PlayerState, isOwner bool) map[string]any {
	info := map[string]any{
		"player_id":              ps.PlayerID,
		"player_name":            ps.PlayerName,
		"hero":                   e.cardToInfo(ps.Hero),
		"elements":               ps.Elements,
		"strict_arcane":          ps.StrictArcane,
		"shield":                 ps.Shield,
		"cannot_gain_shield":     ps.CannotGainShield,
		"next_red_moon_duration": ps.NextRedMoonDuration,
		"next_red_moon_cooldown": ps.NextRedMoonCooldown,
		"charge":                 ps.Charge,
		"temp_modifiers":         ps.TempModifiers,
		"deck_count":             len(ps.Deck),
		"graveyard":              e.cardsToInfo(ps.Graveyard),
		"exile":                  e.cardsToInfo(ps.Exile),
	}

	// Units grid
	units := [3][3]any{}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			units[col][row] = e.cardToInfo(ps.Units[col][row])
		}
	}
	info["units"] = units

	// Terrain grid
	terrain := [3][3]any{}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			terrain[col][row] = e.cardToInfo(ps.Terrain[col][row])
		}
	}
	info["terrain"] = terrain

	// Skills
	skills := make([]any, skillSlotCapacity(ps))
	for i := range skills {
		skills[i] = e.cardToInfoForPlayer(ps, ps.Skills[i])
	}
	info["skills"] = skills
	info["skill_slot_capacity"] = len(skills)

	// Equipment
	equipment := make([]any, equipmentSlotCapacity(ps))
	for i := range equipment {
		if !isOwner && ps.Equipment[i] != nil && ps.Equipment[i].IsSetCounter {
			equipment[i] = hiddenCounterInfo(ps.Equipment[i])
		} else {
			equipment[i] = e.cardToInfo(ps.Equipment[i])
		}
	}
	info["equipment"] = equipment
	info["equipment_slot_capacity"] = len(equipment)

	if isOwner {
		// Show full hand
		info["hand"] = e.cardsToInfoWithEffectiveCosts(ps, ps.Hand, false)
		info["deck_summary"] = deckSummaryToInfo(ps.Deck)
		info["skill_pool"] = e.cardsToInfoWithEffectiveCosts(ps, ps.SkillPool, true)
	} else {
		// Only show count
		info["hand_count"] = len(ps.Hand)
		revealed := make([]*CardInstance, 0)
		for _, card := range ps.Hand {
			if ps.RevealedHand[card.InstanceID] {
				revealed = append(revealed, card)
			}
		}
		info["revealed_hand"] = e.cardsToInfo(revealed)
		info["skill_pool_count"] = len(ps.SkillPool)
	}

	return info
}

func (e *Engine) cardsToInfoWithEffectiveCosts(ps *PlayerState, cards []*CardInstance, learn bool) []map[string]any {
	result := make([]map[string]any, len(cards))
	for i, c := range cards {
		info := e.cardToInfoForPlayer(ps, c)
		if c != nil && c.Card != nil {
			if learn {
				info["effective_learn_cost"] = e.effectiveSkillLearnCost(ps, c)
			} else {
				info["effective_elements_cost"] = e.effectiveCardPlayCost(ps, c)
			}
		}
		result[i] = info
	}
	return result
}

func (e *Engine) cardToInfoForPlayer(ps *PlayerState, card *CardInstance) map[string]any {
	info := e.cardToInfo(card)
	if ps == nil || card == nil || card.Card == nil {
		return info
	}
	if isSpellLikeCard(card.Card) {
		info["has_pierce"] = e.skillHasPierce(ps.PlayerID, card)
		info["spell_area"] = e.effectiveSpellArea(card)
		info["effective_defense_power"] = e.effectiveSkillPowerForPurpose(ps.PlayerID, card, skillPurposeDefend)
		info["effective_defense_boost_power"] = e.effectiveSkillPowerForPurpose(ps.PlayerID, card, skillPurposeDefenseBoost)
		info["effective_attack_power"] = e.effectiveSkillPowerForPurpose(ps.PlayerID, card, skillPurposeAttack)
		info["effective_attack_boost_power"] = e.effectiveSkillPowerForPurpose(ps.PlayerID, card, skillPurposeAttackBoost)
	}
	return info
}

func (e *Engine) effectiveSkillPowerForPurpose(playerID int, skill *CardInstance, purpose skillPurpose) int {
	if skill == nil || skill.Card == nil {
		return 0
	}
	return e.effectiveSkillPowerForPurposeWithData(playerID, skill, nil, purpose, map[string]any{"stat": "power"})
}

func (e *Engine) refreshPendingSpellPowerForModifiedSkill(playerID int, skill *CardInstance) {
	if e == nil || e.State.PendingSpell == nil || skill == nil {
		return
	}
	spell := e.State.PendingSpell
	if spell.AttackerID != playerID {
		return
	}
	if spell.Skill != skill {
		foundBoost := false
		for _, boost := range spell.BoostSkills {
			if boost == skill {
				foundBoost = true
				break
			}
		}
		if !foundBoost {
			return
		}
	}
	powerTargets := append([]SpellTarget{spell.Target}, spell.ExtraTargets...)
	spell.TotalPower = e.effectiveSpellPower(playerID, spell.Skill, spell.BoostSkills, powerTargets...)
	spell.PowerSources = e.spellPowerSources(playerID, spell.Skill, spell.BoostSkills, spell.TotalPower, powerTargets...)
}
