package game

import (
	"eraofarcane/model"
	"fmt"
)

func (Card2521109PunishmentRune) CounterTriggers() []EffectTrigger {
	return []EffectTrigger{TriggerOnSpellCast}
}

func (Card2521109PunishmentRune) CanTriggerCounter(ctx *CounterContext) bool {
	return ctx.Event.Trigger == TriggerOnSpellCast && ctx.Event.PlayerID != ctx.Source.OwnerID && ctx.Event.Card != nil &&
		ctx.Event.Card.Card.IsSkill() && canUseSkillForPurpose(ctx.Event.Card.Card, skillPurposeAttack) &&
		totalSpellsCastThisTurn(ctx.Engine.State.Players[ctx.Event.PlayerID]) > 2 &&
		len(ctx.Engine.friendlyUnits(ctx.Event.PlayerID, false, func(unit *CardInstance) bool {
			return unit != nil && unit.Card != nil && unit.Card.IsCompanion()
		})) > 0
}

type Card2521109PunishmentRune struct{ AlwaysActive }

func (Card2521109PunishmentRune) ID() string { return "2521109" }

func (Card2521109PunishmentRune) Name() string { return "惩戒符文" }

func (Card2521109PunishmentRune) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.OwnerID == ctx.PlayerID {
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.Target.OwnerID, false, func(unit *CardInstance) bool {
		return unit != nil && unit.Card != nil && unit.Card.IsCompanion()
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "punishment_rune_damage",
		"惩戒符文:选择1个敌方伙伴造成2点伤害", candidates, 1, 1,
		nil, false, func(selected []string, _ map[string]any) error {
			target := selectedUnitFromCandidates(ctx.Engine, selected, candidates)
			if target == nil || target.Card == nil || target.OwnerID != ctx.Target.OwnerID || !target.Card.IsCompanion() {
				return fmt.Errorf("invalid punishment rune target")
			}
			ctx.Engine.ApplyDamage(DamageRequest{Target: target, Amount: 2, Kind: "punishment_rune", Element: model.ElementLight, Source: ctx.Source})
			ctx.Engine.emit(GameEvent{
				Type:   "punishment_rune_damage",
				Player: -1,
				Data: map[string]any{
					"player": ctx.PlayerID,
					"source": cardToInfo(ctx.Source),
					"target": cardToInfo(target),
					"damage": 2,
				},
			})
			return nil
		})
	return nil
}
