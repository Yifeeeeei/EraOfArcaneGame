package game

import (
	"eraofarcane/model"
)

type Card2211101DeepSword struct{ AlwaysActive }

func (Card2211101DeepSword) ID() string { return "2211101" }

func (Card2211101DeepSword) Name() string { return "珊瑚秘宝 深邃之剑" }

func (Card2211101DeepSword) OnDraw(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	if drawnPlayer, ok := ctx.ExtraData["drawn_player"].(int); !ok || drawnPlayer != ctx.PlayerID {
		return nil
	}
	if ctx.Engine.currentLearnedSpellPower(ctx.OpponentID) <= ctx.Engine.currentLearnedSpellPower(ctx.PlayerID) {
		return nil
	}
	targets := ctx.Engine.enemyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() && card.Position != nil &&
			ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false)
	})
	if len(targets) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "deep_sword_reveal_damage",
		"珊瑚秘宝 深邃之剑:是否展示并对法力范围内所有敌人造成2点伤害",
		[]map[string]any{candidateInfo(ctx.Source, "hand", "own")}, 0, 1, func(selected []string) {
			if len(selected) == 0 || ctx.Engine.findFriendlyHandCard(ctx.PlayerID, ctx.Source.InstanceID) == nil {
				return
			}
			ps := ctx.Engine.State.Players[ctx.PlayerID]
			if ps.RevealedHand == nil {
				ps.RevealedHand = make(map[string]bool)
			}
			ps.RevealedHand[ctx.Source.InstanceID] = true
			currentTargets := ctx.Engine.enemyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
				return card != nil && card.Card != nil && card.Card.IsCompanion() && card.Position != nil &&
					ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false)
			})
			for _, info := range currentTargets {
				id, _ := info["instance_id"].(string)
				target := findEnemyCardCandidate(ctx.Engine, ctx.PlayerID, id, currentTargets)
				if target == nil || !ctx.Engine.unitStillOnField(target) {
					continue
				}
				ctx.Engine.ApplyDamage(DamageRequest{Target: target, Amount: 2, Kind: "effect", Element: model.ElementWater, Source: ctx.Source, SourcePlayer: ctx.PlayerID, SourceKnown: true})
				if target.CurrentLife <= 0 && !target.Card.IsHero() && ctx.Engine.unitInOwnerGrid(target, target.OwnerID) {
					ctx.Engine.destroyUnitWithData(target, target.OwnerID, map[string]any{
						"death_cause": "deep_sword",
						"source_card": ctx.Source,
						"attacker":    ctx.PlayerID,
					})
				}
			}
		})
	return nil
}

func (Card2211101DeepSword) CanBeFlippedOrSearched(*CardInstance) bool { return false }

func (e *Engine) currentLearnedSpellPower(playerID int) int {
	if e == nil || playerID < 0 || playerID >= len(e.State.Players) {
		return 0
	}
	total := 0
	for _, skill := range e.State.Players[playerID].Skills {
		if skill == nil || skill.Card == nil || !isSpellLikeCard(skill.Card) {
			continue
		}
		total += e.effectiveSkillPowerForPurpose(playerID, skill, skillPurposeAttack)
	}
	return total
}
