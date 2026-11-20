package api

import (
	"encoding/json"
	"eraofarcane/game"
	"eraofarcane/match"
	"net/http"
)

func handleCreateTutorialRoom(rm *match.RoomManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		room, owner, err := rm.CreateTutorialRoom()
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		deck, _ := game.TutorialDecks()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"room_id": room.ID, "player_id": owner, "deck_code": deck.ToDeckCode(), "lesson_id": game.TutorialLessonID})
	}
}
