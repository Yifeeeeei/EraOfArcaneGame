package game

import "eraofarcane/model"

// DamageRequest is the only production entry for damage. Target owns the
// affected player; Source or an explicit SourceKnown attributes the cause.
// A zero request has an unknown source, never an implicit player 0 attacker.
type DamageRequest struct {
	Target               *CardInstance
	Amount               int
	Source               *CardInstance
	SourcePlayer         int
	SourceKnown          bool
	Kind                 string
	Element              string
	Status               string
	Spell                *model.Card
	BoostCount           int
	SkipPlayerShield     bool
	ForcedAttack         bool
	Reason               string
	ActualDamage         map[string]int
	ActualFriendlyDamage map[string]int
}

func (r DamageRequest) triggerData() map[string]any {
	data := map[string]any{"damage_source": r.Kind, "damage_element": r.Element, "status_damage": r.Status}
	if r.Source != nil {
		data["source_card"] = r.Source
		data["attacker"] = r.Source.OwnerID
	} else if r.SourceKnown {
		data["attacker"] = r.SourcePlayer
	}
	if r.Spell != nil {
		data["skill"] = r.Spell.Number
		data["boost_count"] = r.BoostCount
	}
	if r.SkipPlayerShield {
		data["skip_player_shield"] = true
	}
	if r.ForcedAttack {
		data["forced_attack"] = true
	}
	if r.Reason != "" {
		data["reason"] = r.Reason
	}
	if r.ActualDamage != nil {
		data["actual_damage_by_instance"] = r.ActualDamage
	}
	if r.ActualFriendlyDamage != nil {
		data["actual_friendly_damage_by_instance"] = r.ActualFriendlyDamage
	}
	return data
}

func (e *Engine) ApplyDamage(request DamageRequest) {
	if request.Target == nil {
		return
	}
	e.resolveDamage(request.Target, request.Amount, request.Target.OwnerID, request.triggerData())
}
