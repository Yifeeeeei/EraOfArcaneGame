package game

import (
	"fmt"

	"eraofarcane/model"
)

type Card1011002WizardTowerV20260619 struct{ AlwaysActive }

func (Card1011002WizardTowerV20260619) ID() string   { return "1011002" }
func (Card1011002WizardTowerV20260619) Name() string { return "巫师之塔 通天阁" }
func (Card1011002WizardTowerV20260619) HasGlobalSpellRange() bool {
	return true
}
func (Card1011002WizardTowerV20260619) OnEnter(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	for _, skill := range ps.Skills {
		if skill == nil || skill.Card == nil {
			continue
		}
		for _, elem := range model.AllElements {
			if skill.Card.Category == elem {
				ps.Elements[elem]++
				break
			}
		}
	}
	ctx.Source.Statuses["全场法力范围"] = 1
	return nil
}

type Card2011003KingRobeV20260619 struct{ AlwaysActive }

func (Card2011003KingRobeV20260619) ID() string   { return "2011003" }
func (Card2011003KingRobeV20260619) Name() string { return "君王法袍 至贤" }
func (Card2011003KingRobeV20260619) CanReactToSpell(ctx *EffectContext, spell *SpellCast) bool {
	if ctx == nil || spell == nil || spell.AttackerID == ctx.PlayerID {
		return false
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	return len(ps.SkillPool) > 0
}
func (Card2011003KingRobeV20260619) OnSpellReaction(ctx *EffectContext, spell *SpellCast) error {
	robe := Card2011003KingRobeV20260619{}
	if ctx == nil || spell == nil || !robe.CanReactToSpell(ctx, spell) {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	candidates := make([]map[string]any, 0, len(ps.SkillPool))
	for _, skill := range ps.SkillPool {
		if skill != nil {
			candidates = append(candidates, candidateInfo(skill, "skill_pool", "own"))
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "king_robe_remove_skill",
		"君王法袍 至贤:可以将技能池1张技能移出游戏，使本次敌方法术伤害-2", candidates, 0, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			for i, skill := range ps.SkillPool {
				if skill == nil || skill.InstanceID != selected[0] {
					continue
				}
				ps.SkillPool = append(ps.SkillPool[:i], ps.SkillPool[i+1:]...)
				ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
					Type:             TempModFriendlySpellDamageMinus,
					TargetInstanceID: spell.Skill.InstanceID,
					Amount:           2,
					RemainingUses:    1,
					ExpiresTurn:      ctx.Engine.State.TurnNumber + 1,
				})
				ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
					"source": cardToInfo(ctx.Source),
					"effect": "enemy_spell_damage_modifier",
					"amount": -2,
				}})
				return
			}
		})
	return nil
}

type Card2111001FireDragonHeartV20260619 struct{ AlwaysActive }

func (Card2111001FireDragonHeartV20260619) ID() string   { return "2111001" }
func (Card2111001FireDragonHeartV20260619) Name() string { return "火龙之心" }
func (Card2111001FireDragonHeartV20260619) OnPerTurn(ctx *EffectContext) error {
	candidates := make([]map[string]any, 0)
	for _, card := range ctx.Engine.getAllFieldCards(ctx.Engine.State.Players[ctx.PlayerID]) {
		if card == nil || card == ctx.Source || card.Card.IsHero() || effectiveElementsGain(card)[model.ElementFire] <= 0 {
			continue
		}
		candidates = append(candidates, candidateInfo(card, "field", "own"))
	}
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "fire_dragon_heart_sacrifice",
		"火龙之心:献祭包含火负载的卡牌，最多按3点火负载结算", candidates, 1, len(candidates),
		func(selected []string) {
			totalFire := 0
			selectedCards := make([]*CardInstance, 0, len(selected))
			for _, id := range selected {
				card, _ := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, id)
				if card == nil || card == ctx.Source {
					continue
				}
				fireLoad := effectiveElementsGain(card)[model.ElementFire]
				if fireLoad <= 0 {
					continue
				}
				totalFire += fireLoad
				selectedCards = append(selectedCards, card)
			}
			points := min(totalFire, 3)
			if points <= 0 {
				return
			}
			for _, card := range selectedCards {
				if card.Card.IsCompanion() {
					ctx.Engine.destroyUnitWithCause(card, ctx.PlayerID, DeathCauseSacrifice)
				} else {
					ctx.Engine.discardFriendlyCandidate(ctx.PlayerID, card.InstanceID)
				}
			}
			choices := make([]map[string]any, 0, points)
			for i := 1; i <= points; i++ {
				choices = append(choices, map[string]any{
					"instance_id": fmt.Sprintf("attack_%d", i),
					"name":        fmt.Sprintf("第%d点火负载:+1攻", i),
					"zone":        "choice",
				})
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "fire_dragon_heart_bonus",
				"火龙之心:选择哪些火负载给下一次火焰法术+1攻，未选择的改为+3威", choices, 0, points,
				func(selectedBonus []string) {
					attackPoints := len(selectedBonus)
					powerPoints := points - attackPoints
					ctx.Engine.addNextElementSpellDamageBonus(ctx.PlayerID, model.ElementFire, attackPoints)
					ctx.Engine.addNextElementSpellPowerBonus(ctx.PlayerID, model.ElementFire, powerPoints*3)
				})
		})
	return nil
}

type Card2211002WinterBowV20260619 struct{ AlwaysActive }

func (Card2211002WinterBowV20260619) ID() string   { return "2211002" }
func (Card2211002WinterBowV20260619) Name() string { return "嗜魔弓 凛冬" }
func (Card2211002WinterBowV20260619) OnEnter(ctx *EffectContext) error {
	bindSkillToHost(ctx, "3201002")
	return nil
}
func (Card2211002WinterBowV20260619) OnSpellCast(ctx *EffectContext) error {
	if ctx.ExtraData == nil {
		return nil
	}
	if _, ok := ctx.ExtraData["cast_player"].(int); !ok {
		return nil
	}
	ctx.Source.Statuses[winterBowWaterMark]++
	ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
		"source": cardToInfo(ctx.Source),
		"effect": "winter_bow_water_mark",
		"amount": 1,
	}})
	return nil
}
