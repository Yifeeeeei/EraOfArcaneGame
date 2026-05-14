package game

import (
	"log"
	"sort"
)

var baseSetBehaviors = []CardBehavior{
	Card1021006Grocer{},
	Card1121002LivelyHearth{},
	Card1121014Firethorn{},
	Card1221004FrostPuppet{},
	Card1221008IcefieldDemon{},
	Card1321002WindTraveler{},
	Card1321003MagicDandelion{},
	Card1321004LightningElemental{},
	Card1421001SandMage{},
	Card1421014WindbreathMerchant{},
	Card1521002LightforgedTitan{},
	Card1521014TorchWitch{},
	Card1521015EmberWitch{},
	Card1611001ObserverOkoru{},
	Card1621001UnderworldPigeon{},
	Card1621005CursedGolem{},
	Card2321007WindwhisperRing{},
	Card2601002Spellbook{},

	Card4011001Skadi{},
	Card4011002NoFace{},
	Card4111002WitchVerland{},
	Card4111003Brahma{},
	Card4211001Bartel{},
	Card4211003CrystalHeart{},
	Card4311001Su{},
	Card4311003Muling{},
	Card4411001Whitebeard{},
	Card4511001Maris{},
	Card4611001Alice{},
	Card4611002Fuye{},
}

// RegisterAllCardEffects registers every currently supported base-set card
// behavior. The effect registry remains as an engine adapter, but card logic is
// owned by concrete card structs.
func RegisterAllCardEffects() {
	globalRegistry = NewEffectRegistry()
	r := GetEffectRegistry()

	behaviors := append([]CardBehavior(nil), baseSetBehaviors...)
	sort.Slice(behaviors, func(i, j int) bool { return behaviors[i].ID() < behaviors[j].ID() })

	for _, behavior := range behaviors {
		registerBehavior(r, behavior)
	}
	log.Printf("[Effects] Registered %d base-set card behavior objects", len(behaviors))
}
