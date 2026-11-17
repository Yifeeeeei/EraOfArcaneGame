package game

import "eraofarcane/model"

type Card1211003SnowWoman struct{ AlwaysActive }

func (Card1211003SnowWoman) ID() string { return "1211003" }

func (Card1211003SnowWoman) Name() string { return "\"雪女\" 天户凌" }

func (Card1211003SnowWoman) OnEnter(ctx *EffectContext) error {
	ctx.Source.Statuses["引魔"] = 1
	return nil
}

func (Card1211003SnowWoman) OnCardSearched(ctx *EffectContext, event CardSearchedEvent) {
	if event.Card == nil || event.Card.Card == nil || event.Card.Card.Category != model.ElementWater || event.PlayerID != ctx.PlayerID {
		return
	}
	e, source, playerID := ctx.Engine, ctx.Source, ctx.PlayerID
	if !triggeredTurnAvailable(source) {
		return
	}
	targets := e.enemyUnits(playerID, true, func(target *CardInstance) bool {
		return target.Position != nil && e.IsInSpellRange(playerID, target.Position.Col, target.Position.Row, false)
	})
	if len(targets) == 0 || !useTriggeredTurn(source) {
		return
	}
	e.SetPendingAction(playerID, "snow_woman_freeze_after_search",
		"雪女:你检索了水纹卡牌,选择1个法力范围内的敌人冻结1", targets, 1, 1, func(selected []string) {
			target := selectedUnitFromCandidates(e, selected, targets)
			if target != nil {
				e.addStatus(target, StatusFreeze, 1)
			}
		})
}
