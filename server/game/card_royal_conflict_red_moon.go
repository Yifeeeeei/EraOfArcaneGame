package game

import "eraofarcane/model"

const redMoonMarkerStatus = "红月标记"

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

func (e *Engine) updateRedMoonTransformations(playerID int) {
	if playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	active := e.redMoonActive(playerID)
	for _, unit := range e.getAllFieldCards(e.State.Players[playerID]) {
		if unit == nil || unit.Card == nil || unit.Position == nil {
			continue
		}
		switch unit.Card.Number {
		case "1611101":
			if active {
				e.replaceUnitCard(unit, "1601101", false)
			}
		case "1601101":
			if !active {
				e.replaceUnitCard(unit, "1611101", true)
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
