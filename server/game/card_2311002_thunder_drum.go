package game

import (
	"eraofarcane/model"
)

type Card2311002ThunderDrum struct{ AlwaysActive }

func (Card2311002ThunderDrum) ID() string { return "2311002" }

func (Card2311002ThunderDrum) Name() string { return "唤雷震鼓" }

func (Card2311002ThunderDrum) OnDraw(ctx *EffectContext) error {
	if ctx.ExtraData == nil {
		return nil
	}
	drawnPlayer, _ := ctx.ExtraData["drawn_player"].(int)
	if drawnPlayer != ctx.PlayerID {
		return nil
	}
	candidates := []map[string]any{candidateInfo(ctx.Source, "equipment", "own")}
	if drawn, ok := ctx.ExtraData["drawn_card"].(*CardInstance); ok && drawn != nil {
		info := cardToInfo(drawn)
		info["zone"] = "hand"
		info["side"] = "own"
		candidates = append(candidates, info)
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "thunder_drum_mark",
		"唤雷震鼓:是否展示抽到的牌并放置1个标记?", candidates, 0, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			ctx.Source.Statuses["雷鼓标记"]++
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source": cardToInfo(ctx.Source),
				"effect": "thunder_drum_mark",
				"amount": 1,
			}})
		})
	return nil
}

func (Card2311002ThunderDrum) OnPerTurn(ctx *EffectContext) error {
	if thunderDrumMarks(ctx.Source) < 3 || ctx.Engine.State.PendingAction != nil {
		return nil
	}
	choices := []map[string]any{
		{"instance_id": "attack", "name": "本回合你的大气法术+1攻", "zone": "choice"},
		{"instance_id": "stun", "name": "本回合你的大气法术获得晕眩1", "zone": "choice"},
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "thunder_drum_bonus",
		"唤雷震鼓:移除3个标记,选择本回合大气法术获得的效果", choices, 1, 1,
		func(selected []string) {
			if len(selected) == 0 || thunderDrumMarks(ctx.Source) < 3 {
				return
			}
			spendThunderDrumMarks(ctx.Source, 3)
			switch selected[0] {
			case "attack":
				ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{Type: TempModCurrentTurnElementDamage, Element: model.ElementAir, Amount: 1, ExpiresTurn: ctx.Engine.State.TurnNumber})
			case "stun":
				ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{Type: TempModCurrentTurnElementHitStatus, Element: model.ElementAir, Status: StatusStun, Amount: 1, ExpiresTurn: ctx.Engine.State.TurnNumber})
			}
		})
	return nil
}
