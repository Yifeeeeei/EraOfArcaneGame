package game

import (
	"fmt"

	"eraofarcane/model"
)

const redMoonMarkerStatus = "红月标记"
const bloodShadowBodyRedMoonMarkersStatus = "blood_shadow_body_red_moon_markers"

type Card1611101RedMoonWitchSeviana struct{ AlwaysActive }

func (Card1611101RedMoonWitchSeviana) ID() string   { return "1611101" }
func (Card1611101RedMoonWitchSeviana) Name() string { return "红月魔巫 瑟薇安娜" }
func (Card1611101RedMoonWitchSeviana) OnEnter(ctx *EffectContext) error {
	ctx.Engine.addRedMoonMarker(ctx.PlayerID, 1)
	ctx.Engine.updateRedMoonTransformations(ctx.PlayerID)
	return nil
}
func (Card1611101RedMoonWitchSeviana) IsPrayerAbility() bool { return true }
func (Card1611101RedMoonWitchSeviana) OnPerTurn(ctx *EffectContext) error {
	ctx.Engine.addRedMoonMarker(ctx.PlayerID, 1)
	ctx.Engine.updateRedMoonTransformations(ctx.PlayerID)
	return nil
}

type Card1601101BloodShadowBody struct{ AlwaysActive }

func (Card1601101BloodShadowBody) ID() string   { return "1601101" }
func (Card1601101BloodShadowBody) Name() string { return "血影之躯" }
func (Card1601101BloodShadowBody) HasActivePerTurn(card *CardInstance) bool {
	return card != nil && card.Statuses[bloodShadowBodyRedMoonMarkersStatus] > 0
}
func (Card1601101BloodShadowBody) PerTurnLabel(*CardInstance) string {
	return "移除红月标记"
}
func (Card1601101BloodShadowBody) OnPerTurn(ctx *EffectContext) error {
	redMoon := ctx.Engine.redMoonSkill(ctx.PlayerID)
	if redMoon == nil || redMoon.Statuses[redMoonMarkerStatus] <= 0 {
		return fmt.Errorf("血影之躯需要移除1个红月标记")
	}
	redMoon.Statuses[redMoonMarkerStatus]--
	ctx.Source.Statuses[bloodShadowBodyRedMoonMarkersStatus] = redMoon.Statuses[redMoonMarkerStatus]
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModNextSpellExtraTarget,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		RemainingUses:    1,
	})
	ctx.Engine.emit(GameEvent{
		Type:   "blood_shadow_body_extra_target",
		Player: -1,
		Data: map[string]any{
			"player":  ctx.PlayerID,
			"source":  cardToInfo(ctx.Source),
			"markers": redMoon.Statuses[redMoonMarkerStatus],
		},
	})
	return nil
}

type Card3611101RedMoon struct{ AlwaysActive }

func (Card3611101RedMoon) ID() string   { return "3611101" }
func (Card3611101RedMoon) Name() string { return "红月" }
func (Card3611101RedMoon) HasActiveSpellStatModifier(card *CardInstance) bool {
	return abilityDurationActive(card)
}
func (Card3611101RedMoon) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if !isAttackPurpose(ctx) || ctx.Target == nil || ctx.Target.Card == nil || ctx.Target.Card.Category != model.ElementShadow {
		return
	}
	stats.PowerBonus += 2
	if ctx.Target.Card.Number != "3611101" {
		stats.PowerBonus += ctx.Engine.redMoonMarkers(ctx.PlayerID)
	}
}

type Card1621110ScarletBeast struct{ AlwaysActive }

func (Card1621110ScarletBeast) ID() string   { return "1621110" }
func (Card1621110ScarletBeast) Name() string { return "猩红魔兽" }
func (Card1621110ScarletBeast) HasActiveSpellStatModifier(card *CardInstance) bool {
	return card != nil && card.Statuses[StatusPetrify] <= 0
}
func (Card1621110ScarletBeast) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if !isAttackPurpose(ctx) || !ctx.Engine.redMoonActive(ctx.PlayerID) || ctx.Target == nil || ctx.Target.Card == nil || ctx.Target.Card.Category != model.ElementShadow {
		return
	}
	stats.PowerBonus += 2
}

type Card1621111RedMoonProphet struct{ AlwaysActive }

func (Card1621111RedMoonProphet) ID() string   { return "1621111" }
func (Card1621111RedMoonProphet) Name() string { return "红月先知" }
func (Card1621111RedMoonProphet) OnEnter(ctx *EffectContext) error {
	ctx.Engine.reduceCurrentOrNextRedMoonCooldown(ctx.PlayerID, 1)
	return nil
}
func (Card1621111RedMoonProphet) OnDeath(ctx *EffectContext) error {
	ctx.Engine.reduceCurrentOrNextRedMoonCooldown(ctx.PlayerID, 1)
	return nil
}

type Card1621109ScarletWings struct{ AlwaysActive }

func (Card1621109ScarletWings) ID() string   { return "1621109" }
func (Card1621109ScarletWings) Name() string { return "猩红之翼" }

func (e *Engine) triggerScarletWingsAfterRedMoon(playerID int) {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	ps := e.State.Players[playerID]
	for _, source := range append([]*CardInstance(nil), e.getAllFieldCards(ps)...) {
		if source == nil || source.Card == nil || source.Card.Number != "1621109" || source.Position == nil || e.hasEffectiveStatus(source, StatusPetrify) {
			continue
		}
		candidates := e.enemyUnits(playerID, true, func(card *CardInstance) bool {
			return card != nil && card.Position != nil && e.IsInSpellRange(playerID, card.Position.Col, card.Position.Row, false)
		})
		if len(candidates) == 0 {
			continue
		}
		wing := source
		e.SetPendingActionWithError(playerID, "scarlet_wings_red_moon_damage",
			"猩红之翼:选择法力范围内1个单位造成1点伤害", candidates, 1, 1,
			nil, false, func(selected []string, _ map[string]any) error {
				target := selectedUnitFromCandidates(e, selected, candidates)
				if target == nil || target.Position == nil || !e.IsInSpellRange(playerID, target.Position.Col, target.Position.Row, false) {
					return fmt.Errorf("invalid scarlet wings target")
				}
				e.dealDamageWithExtra(target, 1, target.OwnerID, map[string]any{
					"damage_source":  "scarlet_wings",
					"damage_element": model.ElementShadow,
					"source_card":    wing,
					"attacker":       playerID,
				})
				wing.CurrentLife++
				e.emit(GameEvent{
					Type:   "scarlet_wings_red_moon_damage",
					Player: -1,
					Data: map[string]any{
						"player": playerID,
						"source": cardToInfo(wing),
						"target": cardToInfo(target),
						"damage": 1,
					},
				})
				return nil
			})
		return
	}
}

type Card2621105RedMoonPendant struct{ AlwaysActive }

func (Card2621105RedMoonPendant) ID() string   { return "2621105" }
func (Card2621105RedMoonPendant) Name() string { return "红月吊坠" }
func (Card2621105RedMoonPendant) PerTurnLabel(*CardInstance) string {
	return "主动"
}
func (Card2621105RedMoonPendant) OnPerTurn(ctx *EffectContext) error {
	if !ctx.Engine.sacrificeEquipment(ctx.PlayerID, ctx.Source.InstanceID) {
		return fmt.Errorf("red moon pendant must be sacrificed from equipment")
	}
	ctx.Engine.State.Players[ctx.PlayerID].NextRedMoonDuration++
	return nil
}

type Card3621107WillErosion struct{ AlwaysActive }

func (Card3621107WillErosion) ID() string   { return "3621107" }
func (Card3621107WillErosion) Name() string { return "意志侵蚀" }
func (Card3621107WillErosion) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	if !isAttackPurpose(ctx) || !ctx.Engine.redMoonActive(ctx.PlayerID) {
		return
	}
	stats.PowerBonus++
	stats.Pierce = true
}

func isAttackPurpose(ctx *EffectContext) bool {
	if ctx == nil || ctx.ExtraData == nil {
		return true
	}
	purpose, _ := ctx.ExtraData["purpose"].(string)
	return purpose == "" || purpose == string(skillPurposeAttack)
}

func (e *Engine) redMoonSkill(playerID int) *CardInstance {
	if playerID < 0 || playerID >= len(e.State.Players) {
		return nil
	}
	for _, skill := range e.State.Players[playerID].Skills {
		if skill != nil && skill.Card != nil && skill.Card.Number == "3611101" {
			return skill
		}
	}
	return nil
}

func (e *Engine) redMoonMarkers(playerID int) int {
	if redMoon := e.redMoonSkill(playerID); redMoon != nil {
		return redMoon.Statuses[redMoonMarkerStatus]
	}
	return 0
}

func (e *Engine) addRedMoonMarker(playerID int, amount int) {
	if amount <= 0 {
		return
	}
	if redMoon := e.redMoonSkill(playerID); redMoon != nil {
		redMoon.Statuses[redMoonMarkerStatus] += amount
	}
}

func (e *Engine) reduceCurrentOrNextRedMoonCooldown(playerID int, amount int) {
	if amount <= 0 || playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	if redMoon := e.redMoonSkill(playerID); redMoon != nil && e.redMoonActive(playerID) {
		redMoon.Statuses[StatusCooldown] = max(0, redMoon.Statuses[StatusCooldown]-amount)
		if redMoon.Statuses[StatusCooldown] == 0 {
			delete(redMoon.Statuses, StatusCooldown)
		}
		return
	}
	e.State.Players[playerID].NextRedMoonCooldown += amount
}

func (e *Engine) applyNextRedMoonModifiers(playerID int, redMoon *CardInstance) {
	if redMoon == nil || redMoon.Card == nil || redMoon.Card.Number != "3611101" || playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	ps := e.State.Players[playerID]
	if ps.NextRedMoonDuration > 0 {
		redMoon.Statuses[StatusAbilityDuration] += ps.NextRedMoonDuration
		ps.NextRedMoonDuration = 0
	}
	if ps.NextRedMoonCooldown > 0 {
		redMoon.Statuses[StatusCooldown] = max(0, redMoon.Statuses[StatusCooldown]-ps.NextRedMoonCooldown)
		if redMoon.Statuses[StatusCooldown] == 0 {
			delete(redMoon.Statuses, StatusCooldown)
		}
		ps.NextRedMoonCooldown = 0
	}
}

func (e *Engine) updateRedMoonTransformations(playerID int) {
	if playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	active := e.redMoonActive(playerID)
	markers := e.redMoonMarkers(playerID)
	for _, unit := range e.getAllFieldCards(e.State.Players[playerID]) {
		if unit == nil || unit.Card == nil || unit.Position == nil {
			continue
		}
		switch unit.Card.Number {
		case "1611101":
			if active {
				e.replaceUnitCard(unit, "1601101", false)
				unit.Statuses[bloodShadowBodyRedMoonMarkersStatus] = markers
			}
		case "1601101":
			if !active {
				e.replaceUnitCard(unit, "1611101", true)
			} else {
				unit.Statuses[bloodShadowBodyRedMoonMarkersStatus] = markers
			}
		}
	}
}

func (e *Engine) refreshRedMoonState(playerID int) {
	e.updateRedMoonTransformations(playerID)
}

func (e *Engine) replaceUnitCard(unit *CardInstance, number string, reset bool) {
	card := getCardDB()[number]
	if unit == nil || card == nil {
		return
	}
	unit.Card = card
	if reset {
		unit.CurrentLife = card.Life
		unit.CurrentAttack = card.Attack
		unit.IsHorizontal = false
		unit.Statuses = make(map[string]int)
		unit.UsedThisTurn = 0
		unit.UltimateUsed = false
		return
	}
	unit.CurrentLife = min(max(unit.CurrentLife, 1), card.Life)
	unit.CurrentAttack = card.Attack
}
