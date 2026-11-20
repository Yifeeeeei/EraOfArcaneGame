package match

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"eraofarcane/game"
	"time"
)

func (rm *RoomManager) CreateTutorialRoom() (*Room, string, error) {
	var token [16]byte
	if _, err := rand.Read(token[:]); err != nil {
		return nil, "", err
	}
	owner := "learner_" + hex.EncodeToString(token[:])
	r := rm.createRoomWithTutorialOwner(false, owner)
	_, deck := game.TutorialDecks()
	r.mu.Lock()
	r.Players[1] = &RoomPlayer{ID: "lesson_opponent_" + r.ID, Name: "教学对手 · 巴特尔", Deck: deck, IsConnected: true, IsBot: true}
	r.mu.Unlock()
	r.LogRoomEvent("tutorial_room_created", map[string]any{"lesson_id": game.TutorialLessonID})
	return r, owner, nil
}

func (r *Room) IsTutorial() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.tutorialPlayerID != ""
}

func (r *Room) startTutorial() {
	r.mu.Lock()
	if r.botCancel != nil || r.Engine == nil {
		r.mu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	r.botCancel = cancel
	done := make(chan struct{})
	r.botDone = done
	r.mu.Unlock()
	go func() { defer close(done); r.runTutorial(ctx, 450*time.Millisecond) }()
}

func (r *Room) runTutorial(ctx context.Context, delay time.Duration) {
	ticker := time.NewTicker(delay)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
		r.mu.Lock()
		engine := r.Engine
		present := r.Players[0] != nil && r.Players[0].IsConnected
		r.mu.Unlock()
		if engine == nil || !present {
			continue
		}
		action := engine.NextTutorialOpponentAction()
		if action == nil {
			continue
		}
		if _, _, err := engine.HandleActionWithRevision(1, *action); err != nil {
			engine.PauseTutorial("教学对手行动无法完成，请重新开始")
			r.LogRoomEvent("tutorial_action_error", map[string]any{"error": err.Error(), "action": action.Action})
			r.BroadcastState()
			return
		}
		r.BroadcastState()
	}
}
