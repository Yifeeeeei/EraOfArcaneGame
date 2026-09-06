package game

import (
	"eraofarcane/model"
	"strings"
)

func (e *Engine) resolveLavaArmorYeYanShieldBreak(playerID int) {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return
	}
	ps := e.State.Players[playerID]
	if ps == nil {
		return
	}
	kept := ps.TempModifiers[:0]
	for _, modifier := range ps.TempModifiers {
		if modifier.Type != TempModLavaArmorYeYanShieldBreak || modifier.RemainingUses == 0 {
			kept = append(kept, modifier)
			continue
		}
		source := &CardInstance{
			InstanceID: modifier.TargetInstanceID,
			OwnerID:    playerID,
			Card: &model.Card{
				Number: modifier.SourceCardNumber,
				Name:   modifier.SourceName,
			},
		}
		e.equipMoltenArmorForLavaArmorYeYan(playerID, source)
	}
	ps.TempModifiers = kept
}

func (e *Engine) triggerErebosSoulChainMarkedOverexert(playerID int, units []*CardInstance) {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) || len(units) == 0 {
		return
	}
	opponentID := 1 - playerID
	if opponentID < 0 || opponentID >= len(e.State.Players) || !e.playerHasActiveCard(e.State.Players[opponentID], "2611101") {
		return
	}
	for _, unit := range units {
		e.weakenErebosSoulChainMarkedSpellsForUnit(unit)
	}
}

func (e *Engine) triggerAshKeltAfterOpponentShieldBreak(playerID int, data map[string]any) {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) || data == nil {
		return
	}
	attacker, _ := data["attacker"].(int)
	if attacker == playerID {
		return
	}
	ps := e.State.Players[playerID]
	for _, source := range e.getAllFieldCards(ps) {
		if source == nil || source.Card == nil || source.Card.Number != "1511101" || e.hasEffectiveStatus(source, StatusPetrify) {
			continue
		}
		e.drawCards(playerID, 2)
		ps.GainElements(map[string]int{model.ElementLight: 2})
		e.emit(GameEvent{Type: "effect_trigger", Player: playerID, Data: map[string]any{
			"source": cardToInfo(source),
			"effect": "ash_kelt_shield_break",
			"amount": 2,
		}})
	}
}

func (e *Engine) promptLampusSwordDelayedDamage(playerID int, modifier TemporaryModifier) {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) || modifier.Amount <= 0 {
		return
	}
	candidates := e.enemyCompanionsInSpellRange(playerID)
	if len(candidates) == 0 {
		e.removeTemporaryModifier(playerID, modifier.ID)
		return
	}
	maxSelect := min(modifier.Amount, len(candidates))
	e.SetPendingAction(playerID, "lampus_sword_distribute_damage",
		"兰普斯之剑:分配延迟伤害", candidates, 1, maxSelect,
		func(selected []string) {
			allocations := map[string]int{}
			order := make([]string, 0, len(selected))
			for _, id := range selected {
				if allocations[id] == 0 {
					order = append(order, id)
				}
				allocations[id]++
			}
			for remaining := modifier.Amount - len(selected); remaining > 0 && len(order) > 0; remaining-- {
				allocations[order[0]]++
			}
			for _, id := range order {
				target := findEnemyCardCandidate(e, playerID, id, candidates)
				if target == nil || !e.unitStillOnField(target) {
					continue
				}
				e.ApplyDamage(DamageRequest{Target: target, Amount: allocations[id], Kind: "effect", SourcePlayer: playerID, SourceKnown: true})
			}
			e.removeTemporaryModifier(playerID, modifier.ID)
		})
}

func (e *Engine) promptPainScreamWeakenAfterFriendlyDamage(playerID int, target *CardInstance, damage int) {
	if e == nil || target == nil || damage <= 0 || playerID < 0 || playerID >= len(e.State.Players) || e.State.PendingAction != nil {
		return
	}
	ps := e.State.Players[playerID]
	if ps == nil || !playerHasPainScreamModifier(ps) {
		return
	}
	candidates := enemySpellCandidatesWithoutWeaken(e, playerID)
	if len(candidates) == 0 {
		return
	}
	e.SetPendingAction(playerID, "pain_scream_weaken_enemy_spells",
		"苦痛尖啸卷轴:选择没有虚弱的敌方法术获得虚弱2", candidates, 1, min(damage, len(candidates)),
		func(selected []string) {
			weakened := 0
			for _, id := range selected {
				if weakened >= damage {
					break
				}
				skill := findEnemySkillIncludingBound(e, playerID, id)
				if skill == nil || !canInstanceBeWeakened(skill) || skill.Statuses[StatusWeaken] > 0 {
					continue
				}
				e.addStatus(skill, StatusWeaken, 2)
				weakened++
			}
		})
}

func enemySpellCandidatesWithoutWeaken(e *Engine, playerID int) []map[string]any {
	candidates := make([]map[string]any, 0)
	for _, skill := range enemySpellInstancesIncludingBound(e, playerID) {
		if skill != nil && skill.Card != nil && canInstanceBeWeakened(skill) && skill.Statuses[StatusWeaken] <= 0 {
			candidates = append(candidates, candidateInfo(skill, "skill", "enemy"))
		}
	}
	return candidates
}

func clearStatusPrefix(card *CardInstance, prefix string) {
	if card == nil || prefix == "" {
		return
	}
	for status := range card.Statuses {
		if strings.HasPrefix(status, prefix) {
			delete(card.Statuses, status)
		}
	}
}

const coralBellyFirstSpellAttackUsedStatus = "海神之使首次法术攻击已触发"

const soulMarkerStatus = "灵魂标记物"

const curseBoxMarkerStatus = "诅咒魔盒标记物"
