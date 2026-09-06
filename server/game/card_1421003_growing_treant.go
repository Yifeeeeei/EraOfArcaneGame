package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card1421003GrowingTreant struct{ AlwaysActive }

func (Card1421003GrowingTreant) ID() string { return "1421003" }

func (Card1421003GrowingTreant) Name() string { return "成长的树人" }

func (Card1421003GrowingTreant) MasteryMax() int { return 4 }

func (Card1421003GrowingTreant) OnMastery(ctx *EffectContext, level int) error {
	if level != 2 && level != 4 {
		return nil
	}
	loadID := fmt.Sprintf("%s:load:%d", ctx.Source.InstanceID, level)
	lifeID := fmt.Sprintf("%s:life:%d", ctx.Source.InstanceID, level)
	choices := []map[string]any{
		{"instance_id": loadID, "name": "+1地负载", "zone": "choice", "side": "own"},
		{"instance_id": lifeID, "name": "+1生命", "zone": "choice", "side": "own"},
	}
	source := ctx.Source
	ctx.Engine.SetPendingAction(ctx.PlayerID, "growing_treant_mastery_choice", "成长的树人:选择精通奖励", choices, 1, 1, func(selected []string) {
		if len(selected) == 0 || !ctx.Engine.cardStillOnField(source) {
			return
		}
		if selected[0] == lifeID {
			ctx.Engine.gainLife(source, 1, source)
			return
		}
		if selected[0] == loadID {
			ctx.Engine.addElementsGainBonus(source, ctx.PlayerID, model.ElementEarth, 1, source)
		}
	})
	return nil
}
