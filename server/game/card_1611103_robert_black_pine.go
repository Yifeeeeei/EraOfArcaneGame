package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card1611103RobertBlackPine struct{ AlwaysActive }

func (Card1611103RobertBlackPine) ID() string { return "1611103" }

func (Card1611103RobertBlackPine) Name() string { return "鲜血贵公子 罗伯特 黑松" }

func (Card1611103RobertBlackPine) DamageScope() DamageScope { return DamageOtherFriendly }

func (Card1611103RobertBlackPine) OnDamaged(ctx *EffectContext, event DamageEvent) error {
	if ctx == nil || ctx.Source == nil || event.Target == nil {
		return nil
	}
	if event.Target.OwnerID != ctx.PlayerID {
		return nil
	}
	if attacker, ok := event.SourcePlayer, event.SourcePlayer >= 0; !ok || attacker != ctx.PlayerID {
		return nil
	}
	if damage := event.Amount; damage <= 0 {
		return nil
	}
	ctx.Source.Statuses[robertBlackPineMarkerStatus]++
	return nil
}

func (Card1611103RobertBlackPine) OnFriendlyDeath(ctx *EffectContext) error {
	if ctx == nil || ctx.Source == nil || ctx.Target == nil || !bloodThornKilledByFriendlyCard(ctx) {
		return nil
	}
	ctx.Source.Statuses[robertBlackPineMarkerStatus] += 2
	return nil
}

func (Card1611103RobertBlackPine) HasActivePerTurn(card *CardInstance) bool {
	return card != nil && card.Statuses[robertBlackPineMarkerStatus] >= 3
}

func (Card1611103RobertBlackPine) PerTurnLabel(*CardInstance) string {
	return "移除标记"
}

func (Card1611103RobertBlackPine) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Statuses[robertBlackPineMarkerStatus] < 3 {
		return fmt.Errorf("罗伯特需要移除3个标记物")
	}
	choices := []map[string]any{
		{"instance_id": "life", "name": "+1血", "zone": "choice", "side": "own"},
		{"instance_id": "load", "name": "负载+1暗", "zone": "choice", "side": "own"},
		{"instance_id": "attack", "name": "+1攻", "zone": "choice", "side": "own"},
	}
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "robert_black_pine_reward",
		"鲜血贵公子 罗伯特 黑松:选择奖励", choices, 1, 1,
		nil, false, func(selected []string, _ map[string]any) error {
			if ctx.Source.Statuses[robertBlackPineMarkerStatus] < 3 || !ctx.Engine.cardStillOnField(ctx.Source) {
				return fmt.Errorf("invalid Robert reward")
			}
			ctx.Source.Statuses[robertBlackPineMarkerStatus] -= 3
			switch firstSelected(selected) {
			case "life":
				ctx.Engine.gainLife(ctx.Source, 1, ctx.Source)
			case "load":
				ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementShadow, 1, ctx.Source)
			case "attack":
				ctx.Source.AttackBonus++
			default:
				return fmt.Errorf("invalid Robert reward")
			}
			ctx.Engine.emit(GameEvent{
				Type:   "robert_black_pine_reward",
				Player: -1,
				Data: map[string]any{
					"player": ctx.PlayerID,
					"source": cardToInfo(ctx.Source),
					"choice": firstSelected(selected),
				},
			})
			return nil
		})
	return nil
}

const robertBlackPineMarkerStatus = "robert_black_pine_markers"
