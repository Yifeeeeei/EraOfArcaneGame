package game

import (
	"log"
	"sort"
)

type cardEffectRegistrar func(*EffectRegistry)

var cardEffectCatalog = map[string]cardEffectRegistrar{
	// Heroes
	"4011001": registerCard4011001Skadi,
	"4011002": registerCard4011002NoFace,
	"4011101": registerCard4011101PureSpirit,
	"4111002": registerCard4111002WitchVerland,
	"4111003": registerCard4111003Brahma,
	"4111101": registerCard4111101Felin,
	"4111102": registerCard4111102Kran,
	"4111201": registerCard4111201RudokClark,
	"4211001": registerCard4211001Bartel,
	"4211003": registerCard4211003CrystalHeart,
	"4211101": registerCard4211101CoralBelly,
	"4211102": registerCard4211102Sophia,
	"4311001": registerCard4311001Su,
	"4311003": registerCard4311003Muling,
	"4311101": registerCard4311101Soland,
	"4311102": registerCard4311102Fog,
	"4311201": registerCard4311201Lillian,
	"4311202": registerCard4311202Trachi,
	"4311302": registerCard4311302Yuling,
	"4411001": registerCard4411001Whitebeard,
	"4411201": registerCard4411201Hisson,
	"4411202": registerCard4411202Dorothy,
	"4511001": registerCard4511001Maris,
	"4511101": registerCard4511101Sivar,
	"4611001": registerCard4611001Alice,
	"4611002": registerCard4611002Fuye,
	"4611202": registerCard4611202Yuexi,

	// Companions
	"1221004": registerCard1221004FrostPuppet,
	"1221203": registerCard1221203AbyssFogDemon,
	"1221207": registerCard1221207Reconstructor,
	"1221208": registerCard1221208HunterCruiser,
	"1321301": registerCard1321301Technician,
	"1321304": registerCard1321304SkyWell,
	"1321306": registerCard1321306TaijiMaster,
	"1321308": registerCard1321308FloatingPilot,
	"1321309": registerCard1321309TaijiHeir,
	"1611103": registerCard1611103RobertBlackpine,
	"1621103": registerCard1621103BloodPuppet,
	"1621112": registerCard1621112SilentHunter,
	"1621113": registerCard1621113SilentPriest,
	"1621114": registerCard1621114SoulSymbiote,

	// Skills
	"3621206": registerCard3621206ChaosDevour,
	"3621301": registerCard3621301ChargeRecycle,

	// Items
	"2511101": registerCard2511101NinefoldRadiance,
	"2521101": registerCard2521101BlessedLoneStar,
	"2521104": registerCard2521104GoldenDragonbone,
}

// RegisterAllCardEffects registers all card-specific Go implementations.
//
// Card effects are intentionally organized by card number. Card instances own
// runtime state; this catalog owns the per-card behavior shared by instances of
// the same card definition.
func RegisterAllCardEffects() {
	r := GetEffectRegistry()
	ids := make([]string, 0, len(cardEffectCatalog))
	for id := range cardEffectCatalog {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		cardEffectCatalog[id](r)
	}
	log.Printf("[Effects] Registered %d card-specific effect implementations", len(ids))
}

func registerNoopActive(r *EffectRegistry, cardNumber string, trigger EffectTrigger) {
	r.RegisterActive(cardNumber, trigger, func(ctx *EffectContext) error {
		if trigger == TriggerUltimate && ctx.Source != nil {
			ctx.Source.UltimateUsed = true
		}
		return nil
	})
}
