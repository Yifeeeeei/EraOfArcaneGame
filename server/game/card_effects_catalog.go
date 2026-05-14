package game

import (
	"log"
	"sort"
)

var baseSetBehaviorFactories = map[string]func() CardBehavior{
	"1021006": func() CardBehavior { return Card1021006Grocer{} },
	"1121002": func() CardBehavior { return Card1121002LivelyHearth{} },
	"1121014": func() CardBehavior { return Card1121014Firethorn{} },
	"1221004": func() CardBehavior { return Card1221004FrostPuppet{} },
	"1221008": func() CardBehavior { return Card1221008IcefieldDemon{} },
	"1321002": func() CardBehavior { return Card1321002WindTraveler{} },
	"1321003": func() CardBehavior { return Card1321003MagicDandelion{} },
	"1321004": func() CardBehavior { return Card1321004LightningElemental{} },
	"1421001": func() CardBehavior { return Card1421001SandMage{} },
	"1421014": func() CardBehavior { return Card1421014WindbreathMerchant{} },
	"1521002": func() CardBehavior { return Card1521002LightforgedTitan{} },
	"1521014": func() CardBehavior { return Card1521014TorchWitch{} },
	"1521015": func() CardBehavior { return Card1521015EmberWitch{} },
	"1611001": func() CardBehavior { return Card1611001ObserverOkoru{} },
	"1621001": func() CardBehavior { return Card1621001UnderworldPigeon{} },
	"1621005": func() CardBehavior { return Card1621005CursedGolem{} },
	"2321007": func() CardBehavior { return Card2321007WindwhisperRing{} },
	"2601002": func() CardBehavior { return Card2601002Spellbook{} },

	"4011001": func() CardBehavior { return Card4011001Skadi{} },
	"4011002": func() CardBehavior { return Card4011002NoFace{} },
	"4111002": func() CardBehavior { return Card4111002WitchVerland{} },
	"4111003": func() CardBehavior { return Card4111003Brahma{} },
	"4211001": func() CardBehavior { return Card4211001Bartel{} },
	"4211003": func() CardBehavior { return Card4211003CrystalHeart{} },
	"4311001": func() CardBehavior { return Card4311001Su{} },
	"4311003": func() CardBehavior { return Card4311003Muling{} },
	"4411001": func() CardBehavior { return Card4411001Whitebeard{} },
	"4511001": func() CardBehavior { return Card4511001Maris{} },
	"4611001": func() CardBehavior { return Card4611001Alice{} },
	"4611002": func() CardBehavior { return Card4611002Fuye{} },
}

// RegisterAllCardEffects registers lazy factories for currently supported
// base-set card behavior. It does not instantiate every behavior object at
// startup; a behavior is built only when its card number is queried by the
// engine.
func RegisterAllCardEffects() {
	globalRegistry = NewEffectRegistry()
	r := GetEffectRegistry()

	ids := make([]string, 0, len(baseSetBehaviorFactories))
	for id := range baseSetBehaviorFactories {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		r.RegisterBehaviorFactory(id, baseSetBehaviorFactories[id])
	}
	log.Printf("[Effects] Registered %d lazy base-set card behavior factories", len(ids))
}
