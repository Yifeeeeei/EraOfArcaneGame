package game

import (
	"fmt"

	"eraofarcane/model"
)

type Card1011103Gambler struct{ AlwaysActive }

func (Card1011103Gambler) ID() string   { return "1011103" }
func (Card1011103Gambler) Name() string { return "\"弈者\"" }
func (Card1011103Gambler) OnEnter(ctx *EffectContext) error {
	bindSkillToHost(ctx, "3001101")
	return nil
}

func (Card1011103Gambler) OnUltimate(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Source.Card == nil {
		return nil
	}
	for playerID, ps := range ctx.Engine.State.Players {
		if ps == nil {
			continue
		}
		for col := 0; col < 3; col++ {
			for row := 0; row < 3; row++ {
				target := ps.Units[col][row]
				if target == nil || target.Card == nil || target.Card.Number != "1001101" {
					continue
				}
				ctx.Engine.ApplyDamage(DamageRequest{Target: target, Amount: 1, Kind: "effect", Source: ctx.Source, SourcePlayer: ctx.PlayerID, SourceKnown: true})
				if target.CurrentLife <= 0 && ctx.Engine.unitInOwnerGrid(target, playerID) {
					ctx.Engine.destroyUnitWithData(target, playerID, map[string]any{
						"death_cause": "gambler_ultimate",
						"attacker":    ctx.PlayerID,
						"source_card": ctx.Source.Card.Number,
					})
				}
			}
		}
	}
	return nil
}

type Card2511102FiveRainbowRing struct{ AlwaysActive }

func (Card2511102FiveRainbowRing) ID() string   { return "2511102" }
func (Card2511102FiveRainbowRing) Name() string { return "五虹之环" }
func (Card2511102FiveRainbowRing) OnEnter(ctx *EffectContext) error {
	bindSkillToHost(ctx, "3501101")
	return nil
}
func (Card2511102FiveRainbowRing) PerTurnLabel() string { return "回合技" }
func (Card2511102FiveRainbowRing) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.PlayerID < 0 || ctx.PlayerID >= len(ctx.Engine.State.Players) {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	candidates := []map[string]any{}
	for _, elem := range []string{model.ElementFire, model.ElementWater, model.ElementEarth, model.ElementAir, model.ElementLight} {
		if ps.Elements[elem] <= 0 {
			continue
		}
		candidates = append(candidates, map[string]any{
			"instance_id": elem,
			"name":        elem,
			"zone":        "element",
			"side":        "own",
		})
	}
	if len(candidates) == 0 {
		return nil
	}
	sourceID := ctx.Source.InstanceID
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "five_rainbow_ring_marker",
		"五虹之环:支付1点元素并放置同类标记", candidates, 1, 1, nil, false,
		func(selected []string, _ map[string]any) error {
			elem := firstSelected(selected)
			if elem == "" || ps.Elements[elem] <= 0 {
				return fmt.Errorf("invalid rainbow marker element")
			}
			source := ctx.Engine.findEquipment(ps, sourceID)
			if source == nil {
				return nil
			}
			ps.Elements[elem]--
			source.Statuses[fiveRainbowMarkerStatus(elem)]++
			return nil
		})
	return nil
}

func fiveRainbowMarkerStatus(elem string) string {
	return "五虹标记:" + elem
}
