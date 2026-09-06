package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card3421104NaturalEcho struct{ AlwaysActive }

func (Card3421104NaturalEcho) ID() string { return "3421104" }

func (Card3421104NaturalEcho) Name() string { return "自然回响" }

func (Card3421104NaturalEcho) PerTurnLabel(*CardInstance) string {
	return "回响"
}

func (Card3421104NaturalEcho) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	if ctx.Source.UsedThisTurn > 0 {
		return fmt.Errorf("自然回响本回合已经发动")
	}
	candidates := make([]map[string]any, 0)
	for _, card := range ctx.Engine.getAllFieldCards(ctx.Engine.State.Players[ctx.PlayerID]) {
		if card == nil || card.Card == nil || reducibleElementLoad(card, model.ElementEarth) <= 0 {
			continue
		}
		info := candidateInfo(card, "field", "own")
		info["load_element"] = model.ElementEarth
		candidates = append(candidates, info)
	}
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "natural_echo_remove_load",
		"自然回响:移除1点友方卡牌地脉负载", candidates, 1, 1,
		nil, false, func(selected []string, _ map[string]any) error {
			target := ctx.Engine.findFieldCardByInstance(ctx.Engine.State.Players[ctx.PlayerID], firstSelected(selected))
			if target == nil || reducibleElementLoad(target, model.ElementEarth) <= 0 {
				return fmt.Errorf("invalid natural echo target")
			}
			reduceCardElementLoad(target, model.ElementEarth, 1)
			ctx.Source.UsedThisTurn++
			resetCard(ctx.Source)
			ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
				Type:             TempModSkillPowerBonus,
				SourceCardNumber: ctx.Source.Card.Number,
				SourceName:       ctx.Source.Card.Name,
				TargetInstanceID: ctx.Source.InstanceID,
				Amount:           2,
				RemainingUses:    1,
			})
			ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
				Type:             TempModNextSpellExtraTarget,
				SourceCardNumber: ctx.Source.Card.Number,
				SourceName:       ctx.Source.Card.Name,
				TargetInstanceID: ctx.Source.InstanceID,
				RemainingUses:    1,
				AllowSameTarget:  true,
			})
			ctx.Engine.emit(GameEvent{
				Type:   "natural_echo",
				Player: -1,
				Data: map[string]any{
					"player": ctx.PlayerID,
					"source": cardToInfo(ctx.Source),
					"target": cardToInfo(target),
				},
			})
			return nil
		})
	return nil
}

var _ PerTurnAbility = Card3421104NaturalEcho{}
