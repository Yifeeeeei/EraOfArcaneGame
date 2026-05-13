package main

import (
	"eraofarcane/cards"
	"log"
)

func main() {
	if err := cards.LoadCards("../data/all_card_infos.json"); err != nil {
		log.Fatalf("load cards: %v", err)
	}
	if err := cards.WriteSupportedCardInfoSnapshot(cards.SupportedCardInfoSnapshotPath); err != nil {
		log.Fatalf("write supported card snapshot: %v", err)
	}
}
