package api

import (
	"encoding/json"
	"eraofarcane/bot"
	"eraofarcane/match"
	"net/http"
)

func handleCreateSoloRoom(rm *match.RoomManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Profile string `json:"profile"`
		}
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&req); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}
		profile, err := bot.LookupProfile(req.Profile)
		if err != nil {
			jsonError(w, "unknown AI profile", http.StatusBadRequest)
			return
		}
		room, err := rm.CreateSoloRoomWithProfile(profile.ID)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"room_id": room.ID, "profile": profile.ID, "opponent": profile.Name})
	}
}
