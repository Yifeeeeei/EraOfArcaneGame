package game

import (
	"eraofarcane/model"
	"fmt"
	"strings"
)

type Card1321111ThunderlightWarrior struct{ AlwaysActive }

func (Card1321111ThunderlightWarrior) ID() string { return "1321111" }

func (Card1321111ThunderlightWarrior) Name() string { return "雷光战士" }

func (Card1321111ThunderlightWarrior) OnEnter(ctx *EffectContext) error {
	count := 0
	for _, item := range ctx.Engine.State.Players[ctx.PlayerID].Equipment {
		if isThunderlightItem(item) {
			count++
		}
	}
	if count == 0 {
		return nil
	}
	choices := make([]map[string]any, 0, count*4)
	for i := 0; i < count; i++ {
		for _, choice := range []struct {
			id   string
			name string
		}{
			{id: "life", name: "+2血"},
			{id: "attack", name: "+1攻"},
			{id: "air", name: "负载+1气"},
			{id: "light", name: "负载+1光"},
		} {
			choices = append(choices, map[string]any{
				"instance_id": fmt.Sprintf("%s#%d", choice.id, i),
				"number":      "1321111",
				"name":        choice.name,
				"type":        "选择",
				"zone":        "choice",
				"side":        "own",
			})
		}
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "thunderlight_warrior_rewards",
		"雷光战士:每件雷光道具选择1项奖励", choices, count, count,
		func(selected []string) {
			for _, id := range selected {
				reward, _, _ := strings.Cut(id, "#")
				switch reward {
				case "life":
					ctx.Engine.gainLife(ctx.Source, 2, ctx.Source)
				case "attack":
					ctx.Source.AttackBonus++
				case "air":
					ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementAir, 1, ctx.Source)
				case "light":
					ctx.Engine.addElementsGainBonus(ctx.Source, ctx.PlayerID, model.ElementLight, 1, ctx.Source)
				}
			}
		})
	return nil
}
