package game

import "eraofarcane/model"

func (e *Engine) notifyCardSearched(playerID int, card *CardInstance) {
	if e == nil || card == nil || card.Card == nil || card.Card.Category != model.ElementWater || e.State.PendingAction != nil {
		return
	}
	e.promptSnowWomanAfterWaterSearch(playerID)
}

func (e *Engine) promptSnowWomanAfterWaterSearch(playerID int) {
	ps := e.State.Players[playerID]
	if ps == nil {
		return
	}
	for _, source := range e.getAllFieldCards(ps) {
		if source == nil || source.Card == nil || source.Card.Number != "1211003" || e.hasEffectiveStatus(source, StatusPetrify) {
			continue
		}
		if source.UsedThisTurn >= 3 {
			continue
		}
		targets := e.enemyUnits(playerID, true, func(target *CardInstance) bool {
			return target.Position != nil && e.IsInSpellRange(playerID, target.Position.Col, target.Position.Row, false)
		})
		if len(targets) == 0 {
			return
		}
		e.SetPendingAction(playerID, "snow_woman_freeze_after_search",
			"雪女:你检索了水纹卡牌,选择1个法力范围内的敌人冻结1",
			targets, 1, 1, func(selected []string) {
				target := selectedUnitFromCandidates(e, selected, targets)
				if target == nil || source.UsedThisTurn >= 3 {
					return
				}
				e.addStatus(target, StatusFreeze, 1)
				source.UsedThisTurn++
			})
		return
	}
}
