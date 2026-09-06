package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
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
	randomSeed      int64
	rng             *rand.Rand
	objectSequence  uint64
	actionSequence  uint64
	frameSequence   uint64
	activeFrame     *resolutionFrame
	traceSequence   uint64
	resolutionTrace [resolutionTraceCapacity]ResolutionTraceEntry
}

// NewEngine creates a new game engine
func NewEngine(gameID string, callback EventCallback) *Engine {
	return NewEngineWithSeed(gameID, callback, newRandomSeed())
}

// emit sends an event to clients and records it
func (e *Engine) emit(event GameEvent) {
	e.traceResolution("event", event.Type, e.activeFrame)
	e.log = append(e.log, event)
	if e.callback != nil {
		e.callback(event, event.Player)
	}
}

// HandleAction processes a player action
func (e *Engine) HandleAction(playerID int, action ActionMessage) (actionErr error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.actionSequence++
	e.traceResolution("action_start", action.Action, nil)
	defer func() {
		kind := "action_complete"
		if actionErr != nil {
			kind = "action_rejected"
		}
		e.traceResolution(kind, action.Action, nil)
	}()

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
	case "surrender":
		return e.handleSurrender(playerID)
	case "offer_draw":
		return e.handleOfferDraw(playerID)
	case "respond_draw_offer":
		return e.handleRespondDrawOffer(playerID, action)
	default:
		return fmt.Errorf("unknown action: %s", action.Action)
	}
}
