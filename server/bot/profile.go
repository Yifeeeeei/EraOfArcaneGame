// Package bot contains deterministic player policies, independent of rooms,
// WebSockets and the engine's private mutable state.
package bot

import (
	"eraofarcane/model"
	"fmt"
)

const ProfileID = "royal-fire-aggro-v1"
const Name = "快攻火 · 维兰德"
const DeckRevision = "9be680f0afa2ab9f0eda65dd95f99a38c2d61fa9"

// Pinned ArcaneDecks/decks/王权纷争/快攻火.json. Changing this deck requires
// reviewing the matching policy, not fetching a remote list at runtime.
const DeckCode = "4111002 // 1021001 1021001 1021102 1111001 1111103 1121001 1121001 1121005 1121012 1121012 1121014 1121014 1121108 1121108 1121110 1121111 1121111 1121112 2021107 2111002 2121001 2121001 2121002 2121002 2121003 2121004 2121006 2121014 2121014 2121111 // 3021006 3121001 3121002 3121003 3121007 3121008 3121015 3121105 3121107 3121109 // 3101001"

func Deck() (*model.Deck, error) { return model.ParseDeckCode(DeckCode) }

const ThunderlightProfileID = "royal-thunderlight-midrange-v1"
const ThunderlightDeckRevision = "a108ed593bbb234dadcdf291e4a1b72da672e706"
const ThunderlightDeckCode = "4311101 // 1021001 1021001 1021016 1021016 1311001 1321002 1321002 1321003 1321003 1321006 1321006 1321010 1321011 1321011 1321102 1321102 1321110 1321110 1321111 2021101 2021101 2021104 2021104 2021105 2321101 2321103 2321103 2321104 2321105 2321105 // 3021008 3321002 3321003 3321004 3321006 3321013 3321103 3321106 3021101 3521105 //"

type Profile struct {
	ID, Name, DeckCode, DeckRevision string
}

func LookupProfile(id string) (Profile, error) {
	switch id {
	case "", ProfileID:
		return Profile{ProfileID, Name, DeckCode, DeckRevision}, nil
	case ThunderlightProfileID:
		return Profile{ThunderlightProfileID, "雷光气 · 索兰德", ThunderlightDeckCode, ThunderlightDeckRevision}, nil
	default:
		return Profile{}, fmt.Errorf("unknown AI profile")
	}
}

func NewForProfile(id string) (*Policy, error) {
	profile, err := LookupProfile(id)
	if err != nil {
		return nil, err
	}
	p := New()
	p.profile = profile.ID
	return p, nil
}
