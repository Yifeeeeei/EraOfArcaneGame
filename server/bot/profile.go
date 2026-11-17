// Package bot contains deterministic player policies, independent of rooms,
// WebSockets and the engine's private mutable state.
package bot

import "eraofarcane/model"

const ProfileID = "royal-fire-aggro-v1"
const Name = "快攻火 · 维兰德"
const DeckRevision = "9be680f0afa2ab9f0eda65dd95f99a38c2d61fa9"

// Pinned ArcaneDecks/decks/王权纷争/快攻火.json. Changing this deck requires
// reviewing the matching policy, not fetching a remote list at runtime.
const DeckCode = "4111002 // 1021001 1021001 1021102 1111001 1111103 1121001 1121001 1121005 1121012 1121012 1121014 1121014 1121108 1121108 1121110 1121111 1121111 1121112 2021107 2111002 2121001 2121001 2121002 2121002 2121003 2121004 2121006 2121014 2121014 2121111 // 3021006 3121001 3121002 3121003 3121007 3121008 3121015 3121105 3121107 3121109 // 3101001"

func Deck() (*model.Deck, error) { return model.ParseDeckCode(DeckCode) }
