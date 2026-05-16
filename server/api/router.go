package api

import (
	"encoding/json"
	"eraofarcane/cards"
	"eraofarcane/game"
	"eraofarcane/match"
	"eraofarcane/model"
	"fmt"
	"math/rand"
	"net/http"
)

// SetupRoutes configures all HTTP routes
func SetupRoutes(mux *http.ServeMux, rm *match.RoomManager) {
	// API routes (registered before static file handler)
	mux.HandleFunc("/api/cards", handleGetCards)
	mux.HandleFunc("/api/deck/validate", handleValidateDeck)
	mux.HandleFunc("/api/room/create", handleCreateRoom(rm))
	mux.HandleFunc("/api/room/list", handleListRooms(rm))
	mux.HandleFunc("/api/room/info", handleRoomInfo(rm))

	// WebSocket
	mux.HandleFunc("/ws", HandleWebSocket(rm))

	// Serve static files (last, as "/" catches everything)
	mux.Handle("/", http.FileServer(http.Dir("../web")))
}

func handleGetCards(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := make(map[string]map[string]any, len(cards.PlayableCardDB))
	for id, card := range cards.PlayableCardDB {
		entry := map[string]any{
			"number":           card.Number,
			"type":             card.Type,
			"name":             card.Name,
			"category":         card.Category,
			"tag":              card.Tag,
			"description":      card.Description,
			"quote":            card.Quote,
			"elements_cost":    card.ElementsCost,
			"elements_gain":    card.ElementsGain,
			"elements_expense": card.ElementsExpense,
			"version_number":   card.VersionNum,
			"version_name":     card.VersionName,
			"attack":           card.Attack,
			"life":             card.Life,
			"duration":         card.Duration,
			"power":            card.Power,
			"spawns":           card.Spawns,
			"output_path":      card.OutputPath,
		}
		for key, value := range game.CardRuleInfo(card) {
			entry[key] = value
		}
		response[id] = entry
	}
	json.NewEncoder(w).Encode(response)
}

func handleValidateDeck(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		DeckCode string `json:"deck_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	deck, err := model.ParseDeckCode(req.DeckCode)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := deck.Validate(cards.PlayableCardDB); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return deck info
	heroCard := cards.PlayableCardDB[deck.HeroID]
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"valid":       true,
		"card_pool":   cards.BaseVersionName,
		"hero_name":   heroCard.Name,
		"hero_id":     deck.HeroID,
		"main_count":  len(deck.MainDeck),
		"skill_count": len(deck.SkillPool),
		"extra_count": len(deck.ExtraDeck),
	})
}

func handleCreateRoom(rm *match.RoomManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		room := rm.CreateRoom()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"room_id": room.ID,
		})
	}
}

func handleListRooms(rm *match.RoomManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms := rm.ListRooms()
		result := make([]map[string]any, 0, len(rooms))
		for _, room := range rooms {
			result = append(result, room.RoomInfo())
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

func handleRoomInfo(rm *match.RoomManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := r.URL.Query().Get("id")
		if roomID == "" {
			jsonError(w, "missing room id", http.StatusBadRequest)
			return
		}
		room := rm.GetRoom(roomID)
		if room == nil {
			jsonError(w, "room not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(room.RoomInfo())
	}
}

func jsonError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// GeneratePlayerID generates a random player ID
func GeneratePlayerID() string {
	return fmt.Sprintf("p_%d", rand.Intn(1000000))
}
