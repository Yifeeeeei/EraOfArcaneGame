package match

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"eraofarcane/bot"
	"eraofarcane/cards"
	"eraofarcane/model"
	"time"
)

// CreateSoloRoom reserves only the AI seat. The human joins through the normal
// WebSocket path, which starts the game after both seats have been filled.
func (rm *RoomManager) CreateSoloRoom() (*Room, error) {
	return rm.CreateSoloRoomWithProfile(bot.ProfileID)
}

func (rm *RoomManager) CreateSoloRoomWithProfile(id string) (*Room, error) {
	profile, err := bot.LookupProfile(id)
	if err != nil {
		return nil, err
	}
	deck, err := model.ParseDeckCode(profile.DeckCode)
	if err != nil {
		return nil, err
	}
	if err = deck.Validate(cards.PlayableCardDB); err != nil {
		return nil, err
	}
	var token [16]byte
	if _, err = rand.Read(token[:]); err != nil {
		return nil, err
	}
	r := rm.CreateRoom()
	r.mu.Lock()
	r.botProfile = profile.ID
	r.Players[1] = &RoomPlayer{ID: "bot_" + hex.EncodeToString(token[:]), Name: "AI · " + profile.Name, Deck: deck, IsConnected: true, IsBot: true}
	r.mu.Unlock()
	r.LogRoomEvent("solo_room_created", map[string]any{"profile": profile.ID, "deck_revision": profile.DeckRevision})
	return r, nil
}

func (r *Room) PlayerStateView(player int) map[string]any {
	if player < 0 || player > 1 {
		return nil
	}
	r.mu.Lock()
	engine := r.Engine
	opponentBot := r.Players[1-player] != nil && r.Players[1-player].IsBot
	failure := r.botError
	r.mu.Unlock()
	if engine == nil {
		return nil
	}
	view := engine.GetStateForPlayer(player)
	view["bot_opponent"] = opponentBot
	if opponentBot && failure != "" {
		view["bot_error"] = failure
	}
	return view
}

func (r *Room) startBot() {
	r.mu.Lock()
	if r.botCancel != nil || r.Engine == nil || r.Players[1] == nil || !r.Players[1].IsBot {
		r.mu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	r.botCancel = cancel
	done := make(chan struct{})
	r.botDone = done
	r.mu.Unlock()
	go func() { defer close(done); r.runBot(ctx, 450*time.Millisecond) }()
}

func (r *Room) stopBot() {
	r.mu.Lock()
	cancel := r.botCancel
	r.botCancel = nil
	r.mu.Unlock()
	if cancel != nil {
		cancel()
	}
}

func (r *Room) runBot(ctx context.Context, delay time.Duration) {
	ticker := time.NewTicker(delay)
	defer ticker.Stop()
	r.mu.Lock()
	profile := r.botProfile
	r.mu.Unlock()
	policy, err := bot.NewForProfile(profile)
	if err != nil {
		return
	}
	// Polling avoids doing any rule queries or submitting actions from an
	// engine callback (callbacks can run while Engine.mu is held).
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
		r.mu.Lock()
		engine := r.Engine
		humanPresent := r.Players[0] != nil && r.Players[0].IsConnected
		id := ""
		if r.Players[1] != nil {
			id = r.Players[1].ID
		}
		r.mu.Unlock()
		if !humanPresent || engine == nil {
			continue
		}
		observation, err := engine.ObserveDecision(1)
		if err != nil {
			r.LogRoomEvent("bot_observation_error", map[string]any{"error": err.Error()})
			continue
		}
		var phase struct {
			Phase string `json:"phase"`
		}
		_ = json.Unmarshal(observation.State, &phase)
		if phase.Phase == "game_over" {
			return
		}
		decision, act, err := policy.Next(observation)
		if err != nil {
			r.LogRoomEvent("bot_policy_error", map[string]any{"error": err.Error()})
			continue
		}
		if !act {
			continue
		}
		if decision.Action.Action == "surrender" {
			r.mu.Lock()
			r.botError = decision.Reason
			r.mu.Unlock()
		}
		r.LogRoomEvent("bot_decision", map[string]any{"profile": bot.ProfileID, "reason": decision.Reason})
		r.LogClientAction(1, id, decision.Action)
		_, _, err = engine.HandleActionWithRevision(1, decision.Action)
		if err != nil {
			policy.Reject(decision)
			r.LogActionError(1, id, decision.Action, err)
		} else {
			policy.Accept()
		}
		r.BroadcastState()
	}
}
