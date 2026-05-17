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

// Engine manages a single game instance
type Engine struct {
	State    *GameState
	mu       sync.Mutex
	callback EventCallback
	log      []GameEvent
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
	e.mu.Lock()
	defer e.mu.Unlock()

	// Create player states
	e.State.Players[0] = NewPlayerState(0, p1Name, p1Deck)
	e.State.Players[1] = NewPlayerState(1, p2Name, p2Deck)

	// Initialize cards
	e.State.Players[0].InitCards(0)
	e.State.Players[1].InitCards(0)

	// Draw initial hands (4 cards each; Raven starts with one extra card)
	for i := 0; i < 2; i++ {
		initialHandSize := 4
		if e.State.Players[i].Hero != nil && e.State.Players[i].Hero.Card.Number == "4311002" {
			initialHandSize++
		}
		drawn := e.State.Players[i].DrawCards(initialHandSize)
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

// HandleAction processes a player action
func (e *Engine) HandleAction(playerID int, action ActionMessage) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	log.Printf("[Game %s] Player %d action: %s", e.State.GameID, playerID, action.Action)

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
	skill := e.findSkill(ps, instanceID)
	if skill == nil {
		return fmt.Errorf("reaction skill not found")
	}
	if err := e.validateSkillForPurpose(skill, skillPurposeReaction); err != nil {
		return err
	}
	cost := e.effectiveSkillUseCost(ps, skill)
	overexertIDsRaw, _ := action.Data["overexert_ids"].([]any)
	overexertIDs := stringsFromAnySlice(overexertIDsRaw)
	overexertUnits, err := e.collectOverexertUnits(ps, overexertIDs)
	if err != nil {
		return err
	}
	if !canPayCostWithOverexert(ps, cost, overexertUnits) {
		return fmt.Errorf("not enough elements")
	}
	if !payDefenseCost(ps, cost, action, overexertUnits) {
		return fmt.Errorf("invalid payment")
	}

	skill.IsHorizontal = true
	if !e.shouldSkipCooldown(ps, skill) {
		e.ApplyKeywordOnSkillUse(skill)
	}
	e.consumeNextSkillUseModifiers(ps, skill)

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
		// Redraw: put hand back in deck, shuffle, draw 4
		ps := e.State.Players[playerID]
		ps.Deck = append(ps.Deck, ps.Hand...)
		ps.Hand = make([]*CardInstance, 0)
		shuffleDeck(ps.Deck)
		drawn := ps.DrawCards(4)
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
	e.State.CurrentTurn = 0
	e.State.TurnNumber = 1
	e.State.IsFirstTurn = true

	e.emit(GameEvent{
		Type:   "game_start",
		Player: -1,
		Data: map[string]any{
			"first_player": 0,
		},
	})

	e.startTurn()
}

// startTurn begins a new turn for the current player
func (e *Engine) startTurn() {
	ps := e.State.Players[e.State.CurrentTurn]
	ps.SpellsCastThisTurn = make(map[string]int)

	// Elements are gained by consuming vertical cards. Cards reset at end of
	// their owner's turn, so start turn only clears the spent pool.
	e.refreshElements(ps)
	e.applyTurnStartTemporaryModifiers(ps)

	// Draw a card (not on first turn of the game for first player)
	shouldDraw := !e.State.IsFirstTurn || e.State.CurrentTurn == 1
	if shouldDraw && ps.SkipNextDraw {
		ps.SkipNextDraw = false
		e.emit(GameEvent{
			Type:   "effect_trigger",
			Player: e.State.CurrentTurn,
			Data: map[string]any{
				"effect": "skip_draw",
			},
		})
	} else if shouldDraw {
		drawn := ps.DrawCards(1)
		if len(drawn) > 0 {
			e.emit(GameEvent{
				Type:   "draw_card",
				Player: e.State.CurrentTurn,
				Data: map[string]any{
					"card": cardToInfo(drawn[0]),
				},
			})
			// Notify opponent about the draw (without card info)
			e.emit(GameEvent{
				Type:   "opponent_draw",
				Player: 1 - e.State.CurrentTurn,
				Data: map[string]any{
					"count": 1,
				},
			})
		}
	}

	e.State.Phase = PhaseMain

	e.emit(GameEvent{
		Type:   "turn_start",
		Player: -1,
		Data: map[string]any{
			"current_player": e.State.CurrentTurn,
			"turn_number":    e.State.TurnNumber,
			"elements":       ps.Elements,
		},
	})

	// Trigger 回合开始 effects for all cards on the current player's field
	allCards := e.getAllFieldCards(ps)
	for _, card := range allCards {
		e.triggerEffects(TriggerOnTurnStart, card, nil, nil)
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
			}
		}
	}
	// Reset skills
	for i := 0; i < 5; i++ {
		if ps.Skills[i] != nil {
			e.resetCard(ps.Skills[i])
			ps.Skills[i].UsedThisTurn = 0
		}
	}
	// Reset equipment
	for i := 0; i < 5; i++ {
		if ps.Equipment[i] != nil {
			e.resetCard(ps.Equipment[i])
			ps.Equipment[i].UsedThisTurn = 0
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
	if ps.Units[col][row] != nil {
		return fmt.Errorf("position already occupied")
	}

	cost := e.effectiveCardPlayCost(ps, card)
	if !ps.CanPayCost(cost) {
		return fmt.Errorf("not enough elements")
	}

	// Check unit limit (max 9 including hero)
	if ps.CountUnits() >= 9 {
		return fmt.Errorf("unit area is full")
	}

	if !canPayCostForAction(ps, cost, action) {
		return fmt.Errorf("invalid payment")
	}
	if err := e.validateAndApplySummonDevour(playerID, card, action); err != nil {
		return err
	}

	// Pay cost and place
	if !payCostForAction(ps, cost, action) {
		return fmt.Errorf("invalid payment")
	}
	ps.RemoveFromHand(handIdx)
	card.Position = &Position{Col: col, Row: row}
	card.IsHorizontal = true // enters horizontal by default
	card.EnterTurn = e.State.TurnNumber
	ps.Units[col][row] = card

	// Apply keyword effects (速攻 makes it enter vertical, etc.)
	e.ApplyKeywordOnEnter(card)

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

	devourID, _ := action.Data["devour_id"].(string)
	if devourID == "" {
		return fmt.Errorf("%s requires devour before summon", card.Card.Name)
	}

	ps := e.State.Players[playerID]
	target := e.findFieldCardByInstance(ps, devourID)
	if target == nil {
		target = e.findUnitOnGrid(ps, devourID)
	}
	if target == nil || target.Card == nil || !target.Card.IsCompanion() || target.Card.IsHero() || target == card {
		return fmt.Errorf("invalid devour target")
	}
	if !cardSatisfiesDevourRequirement(target, requirement) {
		return fmt.Errorf("devour target load does not satisfy requirement")
	}

	e.destroyUnit(target, playerID)
	return nil
}

// handleConsume handles consuming a card (横置 to gain elements)
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

	// Consume: set horizontal and gain elements
	card.IsHorizontal = true
	gains := effectiveElementsGain(card)
	// First player's first turn: half load
	if e.State.IsFirstTurn && playerID == 0 {
		for k, v := range gains {
			gains[k] = v / 2
		}
	}
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

	// Find skill in skill area
	var skill *CardInstance
	for i := 0; i < 5; i++ {
		if ps.Skills[i] != nil && ps.Skills[i].InstanceID == instanceID {
			skill = ps.Skills[i]
			break
		}
	}
	if skill == nil {
		return fmt.Errorf("skill not found in skill area")
	}

	// Check if skill can be used
	if err := e.validateReadySkill(skill); err != nil {
		return err
	}
	if !canUseSkillForPurpose(skill.Card, skillPurposeAttack) {
		return fmt.Errorf("skill cannot be used to attack")
	}

	// Check cost
	cost := e.effectiveSkillUseCost(ps, skill)
	if !ps.CanPayCost(cost) {
		return fmt.Errorf("not enough elements")
	}
	if skill.Card.Number == "3021011" && !validateSingleElementPayment(ps.Elements, cost, action) {
		return fmt.Errorf("overlord sanction cost must be paid with one element")
	}

	target := SpellTarget{
		Type:     targetType,
		Position: Position{Col: int(targetColF), Row: int(targetRowF)},
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
	if !ps.CanPayCost(totalCost) {
		return fmt.Errorf("not enough elements for boost skills")
	}

	// Pay costs and set cards horizontal only after all validation succeeds.
	if !payCostForAction(ps, totalCost, action) {
		return fmt.Errorf("invalid payment")
	}
	skill.IsHorizontal = true
	tapSkills(boostSkills)

	// Apply cooldown from keyword
	if !e.shouldSkipCooldown(ps, skill) {
		e.ApplyKeywordOnSkillUse(skill)
	}
	e.consumeNextSkillUseModifiers(ps, skill)

	totalPower := e.effectiveSpellPower(playerID, skill, boostSkills, target)
	e.consumeNextElementSpellPowerBonus(ps, skill)

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
	e.triggerFieldEffectsWithData(TriggerOnSpellCast, playerID, skill, spellCastData)
	e.triggerFieldEffectsWithData(TriggerOnSpellCast, 1-playerID, skill, spellCastData)

	if isSorcery {
		// Sorcery resolves immediately
		e.resolveSpellHit(playerID, skill, target, boostSkills, extraTargets)
	} else {
		// Regular spell: open defense window
		e.State.PendingSpell = &SpellCast{
			AttackerID:   playerID,
			Skill:        skill,
			Target:       target,
			TotalPower:   totalPower,
			BoostSkills:  boostSkills,
			ExtraTargets: extraTargets,
		}
		e.State.Phase = PhaseDefenseWindow

		e.emit(GameEvent{
			Type:   "defense_window",
			Player: 1 - playerID,
			Data: map[string]any{
				"timeout": 30,
			},
		})
	}

	return nil
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
	boostIDsRaw, _ := action.Data["boost_ids"].([]any)
	overexertIDsRaw, _ := action.Data["overexert_ids"].([]any)

	ps := e.State.Players[playerID]

	defenseIDs := stringsFromAnySlice(defenseIDsRaw)
	boostIDs := stringsFromAnySlice(boostIDsRaw)
	overexertIDs := stringsFromAnySlice(overexertIDsRaw)
	defenseSkills, defenseCost, err := e.collectSkillUses(ps, defenseIDs, skillPurposeDefend, nil)
	if err != nil {
		return err
	}
	usedIDs := skillIDSet(defenseSkills)
	boostSkills, boostCost, err := e.collectSkillUses(ps, boostIDs, skillPurposeDefenseBoost, usedIDs)
	if err != nil {
		return err
	}
	overexertUnits, err := e.collectOverexertUnits(ps, overexertIDs)
	if err != nil {
		return err
	}
	totalCost := mergeElementCosts(defenseCost, boostCost)
	if !canPayCostWithOverexert(ps, totalCost, overexertUnits) {
		return fmt.Errorf("not enough elements for defense")
	}
	if len(defenseSkills)+len(boostSkills) > 0 {
		if !payDefenseCost(ps, totalCost, action, overexertUnits) {
			return fmt.Errorf("invalid payment")
		}
		tapSkills(defenseSkills)
		tapSkills(boostSkills)
	}
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
			"overexerted":   len(overexertUnits),
		},
	})

	defenseSuccess := totalDefPower >= attackPower && len(defenseSkills) > 0
	defendData := map[string]any{
		"defender":        playerID,
		"attacker":        e.State.PendingSpell.AttackerID,
		"defense_power":   totalDefPower,
		"attack_power":    attackPower,
		"defense_success": defenseSuccess,
		"attack_skill":    e.State.PendingSpell.Skill,
		"boost_skills":    e.State.PendingSpell.BoostSkills,
	}
	for _, defenseSkill := range defenseSkills {
		e.triggerEffects(TriggerOnDefend, defenseSkill, nil, defendData)
	}

	if defenseSuccess {
		// Defense successful
		e.emit(GameEvent{
			Type:   "defense_success",
			Player: -1,
			Data:   map[string]any{"defender": playerID},
		})
	} else {
		// Defense failed, spell hits
		e.resolveSpellHit(
			e.State.PendingSpell.AttackerID,
			e.State.PendingSpell.Skill,
			e.State.PendingSpell.Target,
			e.State.PendingSpell.BoostSkills,
			e.State.PendingSpell.ExtraTargets,
		)
	}

	e.State.PendingSpell = nil
	if e.State.PendingAction == nil {
		e.State.Phase = PhaseMain
	}
	e.checkWinCondition()

	return nil
}

func (e *Engine) collectOverexertUnits(ps *PlayerState, ids []string) ([]*CardInstance, error) {
	units := make([]*CardInstance, 0, len(ids))
	seen := make(map[string]bool)
	for _, id := range ids {
		if seen[id] {
			return nil, fmt.Errorf("unit %s selected more than once", id)
		}
		seen[id] = true
		unit := e.findUnitOnGrid(ps, id)
		if unit == nil {
			return nil, fmt.Errorf("overexert unit not found: %s", id)
		}
		if !e.canConsumeCard(unit) {
			return nil, fmt.Errorf("unit cannot be overexerted: %s", id)
		}
		units = append(units, unit)
	}
	return units, nil
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

	// Spell hits
	e.resolveSpellHit(
		e.State.PendingSpell.AttackerID,
		e.State.PendingSpell.Skill,
		e.State.PendingSpell.Target,
		e.State.PendingSpell.BoostSkills,
		e.State.PendingSpell.ExtraTargets,
	)

	e.State.PendingSpell = nil
	if e.State.PendingAction == nil {
		e.State.Phase = PhaseMain
	}
	e.checkWinCondition()

	return nil
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
	e.State.PendingSpell = nil
	if e.State.PendingAction == nil {
		e.State.Phase = PhaseMain
	}
}

// resolveSpellHit applies spell damage to the target
func (e *Engine) resolveSpellHit(attackerID int, skill *CardInstance, target SpellTarget, boostSkills []*CardInstance, extraTargets []SpellTarget) {
	defenderID := 1 - attackerID
	if friendly, ok := behaviorForNumber(skill.Card.Number).(FriendlySpellTargetBehavior); ok && friendly.HasActiveFriendlySpellTarget(skill) && friendly.AllowsFriendlySpellTarget() && target.Type == "unit" && target.Position.Valid() {
		if e.State.Players[attackerID].Units[target.Position.Col][target.Position.Row] != nil {
			defenderID = attackerID
		}
	}
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
	var targetUnit *CardInstance
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

	if dmg > 0 {
		for _, unit := range affectedUnits {
			e.dealDamageWithExtra(unit, dmg, defenderID, map[string]any{
				"damage_source":  "spell",
				"damage_element": skill.Card.Category,
				"skill":          skill.Card.Number,
				"attacker":       attackerID,
				"boost_count":    len(boostSkills),
			})
		}
	}
	e.applyGenericSpellEffects(attackerID, defenderID, skill, affectedUnits, target)
	e.applyTemporarySpellHitStatus(attackerID, skill, affectedUnits)

	// Trigger spell hit effects on the skill card itself
	hitData := map[string]any{
		"damage":         dmg,
		"target":         target,
		"affected_units": affectedUnits,
	}
	e.triggerEffects(TriggerOnSpellHit, skill, targetUnit, hitData)
	e.triggerFieldEffectsWithData(TriggerOnSpellHit, attackerID, skill, hitData)
	if skill.Statuses["下一次范围前排"] > 0 {
		skill.Statuses["下一次范围前排"]--
	}
}

func (e *Engine) spellAffectedUnits(defenderID int, skill *CardInstance, target SpellTarget) []*CardInstance {
	if target.Type != "unit" {
		return nil
	}
	defender := e.State.Players[defenderID]
	units := make([]*CardInstance, 0, 9)

	switch e.effectiveSpellArea(skill) {
	case SpellAreaSquare:
		for col := 0; col < 3; col++ {
			for row := 0; row < 3; row++ {
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
		e.applyGenericStatusFromDescription(skill, unit)
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
	if attacker == nil {
		return fmt.Errorf("attacker not found")
	}
	if attacker.Card.Attack <= 0 {
		return fmt.Errorf("unit has no attack")
	}
	if attacker.IsHorizontal {
		return fmt.Errorf("attacker is horizontal")
	}
	if e.hasEffectiveStatus(attacker, StatusStun) {
		return fmt.Errorf("attacker is stunned")
	}

	// Check attacker is in front row
	frontRow := ps.GetFrontRow()
	if attacker.Position.Row != frontRow {
		return fmt.Errorf("attacker is not in front row")
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
	if !e.IsInAttackRange(playerID, attacker, targetCol, targetRow) {
		return fmt.Errorf("target is not in attack range")
	}

	// Consume attacker (横置)
	attacker.IsHorizontal = true

	// Trigger 攻击时 effects
	e.triggerEffects(TriggerOnAttack, attacker, target, nil)

	dmg := attacker.CurrentAttack

	e.emit(GameEvent{
		Type:   "unit_attack",
		Player: -1,
		Data: map[string]any{
			"attacker_player": playerID,
			"attacker":        cardToInfo(attacker),
			"target":          cardToInfo(target),
			"target_pos":      targetPos,
			"damage":          dmg,
		},
	})

	// Deal damage (unit attacks cannot be defended)
	if dmg > 0 {
		e.dealDamageWithExtra(target, dmg, 1-playerID, map[string]any{"damage_source": "attack"})
		// Trigger 命中 effects
		e.triggerEffects(TriggerOnHit, attacker, target, map[string]any{"damage": dmg})
	}

	e.checkWinCondition()
	return nil
}

// dealDamage deals damage to a card instance
func (e *Engine) dealDamage(target *CardInstance, amount int, ownerID int) {
	e.dealDamageWithExtra(target, amount, ownerID, nil)
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

	// Apply shield damage reduction
	amount = ApplyShieldDamage(target, amount)
	if amount <= 0 {
		e.emit(GameEvent{
			Type:   "shield_block",
			Player: -1,
			Data: map[string]any{
				"target": cardToInfo(target),
				"shield": target.Statuses["护盾"],
			},
		})
		return
	}
	if target.Statuses["防止致命"] > 0 && target.CurrentLife-amount <= 0 {
		amount = max(target.CurrentLife-1, 0)
		if amount <= 0 {
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
	}

	target.CurrentLife -= amount

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
		e.destroyUnit(target, ownerID)
	}
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
	ps := e.State.Players[ownerID]

	// Remove from grid
	if unit.Position != nil {
		ps.Units[unit.Position.Col][unit.Position.Row] = nil
	}

	// Bound skills live only while their host is on the battlefield. They do not
	// enter the graveyard as independent cards.
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

	// Trigger 遗言 (on death) effects
	e.triggerEffects(TriggerOnDeath, unit, nil, nil)

	// Notify friendly cards about the death
	e.triggerFieldEffects(TriggerOnFriendlyDeath, ownerID, unit)

	// Notify enemy cards about the death
	e.triggerFieldEffects(TriggerOnEnemyDeath, 1-ownerID, unit)

	// Check if hero died
	if unit.Card.IsHero() {
		e.State.Winner = 1 - ownerID
		e.State.Phase = PhaseGameOver
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

// handleEquip handles equipping an item
func (e *Engine) handleEquip(playerID int, action ActionMessage) error {
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
	cost := e.effectiveCardPlayCost(ps, card)
	if !ps.CanPayCost(cost) {
		return fmt.Errorf("not enough elements")
	}

	// Find empty equipment slot
	slotIdx := -1
	for i := 0; i < 5; i++ {
		if ps.Equipment[i] == nil {
			slotIdx = i
			break
		}
	}
	if slotIdx == -1 {
		return fmt.Errorf("equipment area is full")
	}

	if !payCostForAction(ps, cost, action) {
		return fmt.Errorf("invalid payment")
	}
	ps.RemoveFromHand(handIdx)
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
	if !ps.CanPayCost(cost) {
		return fmt.Errorf("not enough elements")
	}
	if skill.Card.Number == "3021011" && !validateSingleElementPayment(ps.Elements, cost, action) {
		return fmt.Errorf("overlord sanction cost must be paid with one element")
	}

	// Find slot
	slotIdx := -1
	if replaceID != "" {
		// Replace existing skill
		for i := 0; i < 5; i++ {
			if ps.Skills[i] != nil && ps.Skills[i].InstanceID == replaceID {
				if !ps.Skills[i].IsHorizontal {
					return fmt.Errorf("can only replace horizontal skills")
				}
				// Send replaced skill to graveyard
				ps.Graveyard = append(ps.Graveyard, ps.Skills[i])
				slotIdx = i
				break
			}
		}
		if slotIdx == -1 {
			return fmt.Errorf("replacement skill not found")
		}
	} else {
		// Find empty slot
		for i := 0; i < 5; i++ {
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
	if !payCostForAction(ps, cost, action) {
		return fmt.Errorf("invalid payment")
	}
	e.consumeEarthSkillLearnCostModifier(ps, skill)
	ps.SkillPool = append(ps.SkillPool[:poolIdx], ps.SkillPool[poolIdx+1:]...)
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

	// Regular consumable item
	cost := e.effectiveCardPlayCost(ps, card)
	if !ps.CanPayCost(cost) {
		return fmt.Errorf("not enough elements")
	}

	// Pay and use
	if !payCostForAction(ps, cost, action) {
		return fmt.Errorf("invalid payment")
	}
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

	e.triggerEffects(TriggerOnUseItem, card, nil, nil)

	return nil
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
	if !ps.CanPayCost(cost) {
		return fmt.Errorf("not enough elements")
	}

	// Pay cost and place
	if !payCostForAction(ps, cost, action) {
		return fmt.Errorf("invalid payment")
	}
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

	// Execute callback
	callback := pa.Callback
	e.State.PendingAction = nil
	e.State.Phase = e.State.ResumePhase

	if callback != nil {
		callback(selected)
	}

	if e.State.PendingAction == nil && e.State.Phase == PhaseDefenseWindow && e.State.PendingSpell == nil {
		e.State.Phase = PhaseMain
	}

	e.checkWinCondition()
	return nil
}

// SetPendingAction sets a pending player action and pauses the game
func (e *Engine) SetPendingAction(playerID int, actionType string, prompt string, candidates []map[string]any, minSelect, maxSelect int, callback func([]string)) {
	if minSelect > 0 && len(candidates) == 0 {
		return
	}
	e.State.ResumePhase = e.State.Phase
	e.State.Phase = PhaseWaitingAction
	e.State.PendingAction = &PendingAction{
		Type:       actionType,
		PlayerID:   playerID,
		Prompt:     prompt,
		Candidates: candidates,
		MinSelect:  minSelect,
		MaxSelect:  maxSelect,
		Callback:   callback,
	}

	e.emit(GameEvent{
		Type:   "pending_action",
		Player: playerID,
		Data: map[string]any{
			"type":       actionType,
			"prompt":     prompt,
			"candidates": candidates,
			"min_select": minSelect,
			"max_select": maxSelect,
		},
	})
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
	if len(ps.Hand) > e.State.HandLimit {
		discardCount := len(ps.Hand) - e.State.HandLimit
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
						ps.Graveyard = append(ps.Graveyard, c)
						e.emit(GameEvent{
							Type:   "discard",
							Player: currentTurn,
							Data:   map[string]any{"card": cardToInfo(c)},
						})
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
	e.applyOpponentTurnEndTemporaryModifiers(ps.PlayerID)
	e.discardMarkedEndOfTurnCards(ps)
	e.applyLoadGainAtTurnEnd(ps)

	e.clearExpiredTemporaryModifiers(ps.PlayerID)

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
	if e.State.IsFirstTurn && e.State.CurrentTurn == 0 {
		e.State.IsFirstTurn = false
	}
	e.State.CurrentTurn = 1 - e.State.CurrentTurn
	if e.State.CurrentTurn == 0 {
		e.State.TurnNumber++
	}

	if e.State.Phase != PhaseGameOver {
		e.startTurn()
	}
}

// processEndOfTurnStatuses processes status marks at end of turn
func (e *Engine) processEndOfTurnStatuses(ps *PlayerState) {
	allCards := e.getAllFieldCards(ps)

	for _, card := range allCards {
		// 点燃: remove 1 stack, deal 1 fire damage
		if card.Statuses[StatusBurn] > 0 {
			card.Statuses[StatusBurn]--
			e.dealDamageWithExtra(card, 1, ps.PlayerID, map[string]any{"status_damage": StatusBurn})
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
			card.Statuses[StatusPetrify]--
		}
		// 冷却: remove 1 stack
		if card.Statuses[StatusCooldown] > 0 {
			card.Statuses[StatusCooldown]--
		}
	}

	// 虚弱 is on skills, handled separately
	for i := 0; i < 5; i++ {
		if ps.Skills[i] != nil && ps.Skills[i].Statuses[StatusWeaken] > 0 {
			ps.Skills[i].Statuses[StatusWeaken]--
		}
		if ps.Skills[i] != nil && ps.Skills[i].Statuses[StatusSeal] > 0 {
			ps.Skills[i].Statuses[StatusSeal]--
		}
	}
}

// checkWinCondition checks if the game is over
func (e *Engine) checkWinCondition() {
	if e.State.Phase == PhaseGameOver {
		return
	}
	for i := 0; i < 2; i++ {
		if e.State.Players[i].Hero != nil && e.State.Players[i].Hero.CurrentLife <= 0 {
			e.State.Winner = 1 - i
			e.State.Phase = PhaseGameOver
			e.emit(GameEvent{
				Type:   "game_over",
				Player: -1,
				Data: map[string]any{
					"winner": e.State.Winner,
					"reason": "hero_killed",
				},
			})
			return
		}
	}
}

func payCostForAction(ps *PlayerState, cost map[string]int, action ActionMessage) bool {
	if payment := paymentFromAction(action); payment != nil {
		return ps.PayCostWithPayment(cost, payment)
	}
	return ps.PayCost(cost)
}

func canPayCostForAction(ps *PlayerState, cost map[string]int, action ActionMessage) bool {
	if payment := paymentFromAction(action); payment != nil {
		return validateElementPayment(ps.Elements, cost, payment)
	}
	return ps.CanPayCost(cost)
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
		"turn_number":  state.TurnNumber,
		"winner":       state.Winner,
		"you":          playerStateToInfo(ps, true),
		"opponent":     playerStateToInfo(op, false),
		"pending_spell": func() any {
			if state.PendingSpell != nil {
				return map[string]any{
					"attacker": state.PendingSpell.AttackerID,
					"skill":    cardToInfo(state.PendingSpell.Skill),
					"target":   state.PendingSpell.Target,
					"power":    state.PendingSpell.TotalPower,
				}
			}
			return nil
		}(),
		"pending_action": func() any {
			if state.PendingAction != nil && state.PendingAction.PlayerID == playerID {
				return map[string]any{
					"type":       state.PendingAction.Type,
					"prompt":     state.PendingAction.Prompt,
					"candidates": state.PendingAction.Candidates,
					"min_select": state.PendingAction.MinSelect,
					"max_select": state.PendingAction.MaxSelect,
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
	for i := 0; i < 5; i++ {
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

func (e *Engine) findSkill(ps *PlayerState, instanceID string) *CardInstance {
	for i := 0; i < 5; i++ {
		if ps.Skills[i] != nil && ps.Skills[i].InstanceID == instanceID {
			return ps.Skills[i]
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
	for i := 0; i < 5; i++ {
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
	for i := 0; i < 5; i++ {
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
		"instance_id":      ci.InstanceID,
		"number":           ci.Card.Number,
		"name":             ci.Card.Name,
		"type":             ci.Card.Type,
		"category":         ci.Card.Category,
		"tag":              ci.Card.Tag,
		"description":      ci.Card.Description,
		"attack":           ci.Card.Attack + ci.AttackBonus,
		"life":             ci.Card.Life,
		"power":            ci.Card.Power + ci.PowerBonus,
		"duration":         ci.Card.Duration,
		"elements_cost":    ci.Card.ElementsCost,
		"elements_gain":    effectiveElementsGain(ci),
		"elements_expense": ci.Card.ElementsExpense,
		"current_life":     ci.CurrentLife,
		"current_attack":   ci.CurrentAttack,
		"is_horizontal":    ci.IsHorizontal,
		"is_terrain":       cards.IsTerrain(ci.Card.Number),
		"is_consumable":    cards.IsConsumable(ci.Card.Number),
		"is_equipment":     cards.IsEquipment(ci.Card.Number),
		"is_weapon":        cards.IsWeapon(ci.Card.Number),
		"statuses":         ci.Statuses,
		"position":         ci.Position,
		"output_path":      ci.Card.OutputPath,
		"used_this_turn":   ci.UsedThisTurn,
		"ultimate_used":    ci.UltimateUsed,
		"uses_remaining":   ci.UsesRemaining,
	}
	if len(ci.BoundSkills) > 0 {
		info["bound_skills"] = cardsToInfo(ci.BoundSkills)
	}
	if attached := attachedBehaviorsInfo(ci); len(attached) > 0 {
		info["attached_behaviors"] = attached
	}

	hasPerTurn := cardHasActivePerTurn(ci)
	hasUltimate := cardHasActiveUltimate(ci)
	info["has_per_turn"] = hasPerTurn
	info["has_ultimate"] = hasUltimate

	if hasPerTurn {
		info["per_turn_limit"] = perTurnLimit(ci)
	}
	if requirement := summonDevourRequirement(ci); len(requirement) > 0 {
		info["devour_requirement"] = requirement
	}

	// Mark defense-only skills
	if ci.Card.IsSkill() {
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
		result = append(result, map[string]any{
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
		})
	}

	return result
}

func playerStateToInfo(ps *PlayerState, isOwner bool) map[string]any {
	info := map[string]any{
		"player_id":      ps.PlayerID,
		"player_name":    ps.PlayerName,
		"hero":           cardToInfo(ps.Hero),
		"elements":       ps.Elements,
		"charge":         ps.Charge,
		"temp_modifiers": ps.TempModifiers,
		"deck_count":     len(ps.Deck),
		"graveyard":      cardsToInfo(ps.Graveyard),
	}

	// Units grid
	units := [3][3]any{}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			units[col][row] = cardToInfo(ps.Units[col][row])
		}
	}
	info["units"] = units

	// Terrain grid
	terrain := [3][3]any{}
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			terrain[col][row] = cardToInfo(ps.Terrain[col][row])
		}
	}
	info["terrain"] = terrain

	// Skills
	skills := [5]any{}
	for i := 0; i < 5; i++ {
		skills[i] = cardToInfo(ps.Skills[i])
	}
	info["skills"] = skills

	// Equipment
	equipment := [5]any{}
	for i := 0; i < 5; i++ {
		equipment[i] = cardToInfo(ps.Equipment[i])
	}
	info["equipment"] = equipment

	if isOwner {
		// Show full hand
		info["hand"] = cardsToInfo(ps.Hand)
		info["deck_summary"] = deckSummaryToInfo(ps.Deck)
		info["skill_pool"] = cardsToInfo(ps.SkillPool)
	} else {
		// Only show count
		info["hand_count"] = len(ps.Hand)
		revealed := make([]*CardInstance, 0)
		for _, card := range ps.Hand {
			if ps.RevealedHand[card.InstanceID] {
				revealed = append(revealed, card)
			}
		}
		info["revealed_hand"] = cardsToInfo(revealed)
		info["skill_pool_count"] = len(ps.SkillPool)
	}

	return info
}
