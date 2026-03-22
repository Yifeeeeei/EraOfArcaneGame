package main

import (
	"eraofarcane/api"
	"eraofarcane/cards"
	"eraofarcane/game"
	"eraofarcane/match"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Load card data
	cardDataPath := "../data/all_card_infos.json"
	if err := cards.LoadCards(cardDataPath); err != nil {
		log.Fatalf("Failed to load cards: %v", err)
	}

	// Set card DB reference for game engine
	game.SetCardDB(cards.CardDB)

	// Register card effects (manual + auto-parsed)
	game.RegisterAllCardEffects()
	game.AutoParseAndRegister()

	// Create room manager
	rm := match.NewRoomManager()

	// Setup routes
	mux := http.NewServeMux()
	api.SetupRoutes(mux, rm)

	port := 9090
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server starting on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
