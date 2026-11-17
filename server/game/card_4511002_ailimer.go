package game

type Card4511002Ailimer struct{ AlwaysActive }

func (Card4511002Ailimer) ID() string { return "4511002" }

func (Card4511002Ailimer) Name() string { return "神之眷子 爱里默" }

func (Card4511002Ailimer) OnEnter(ctx *EffectContext) error {
	ailimerShuffleShacklesOnce(ctx)
	ailimerUnlockIfShacklesCleared(ctx)
	return nil
}

func (Card4511002Ailimer) OnTurnStart(ctx *EffectContext) error {
	ailimerShuffleShacklesOnce(ctx)
	ailimerUnlockIfShacklesCleared(ctx)
	return nil
}

func (Card4511002Ailimer) OnUseItem(ctx *EffectContext) error {
	if ctx.Target != nil && ctx.Target.Card != nil && ctx.Target.Card.Number == "2501001" {
		ailimerUnlockIfShacklesCleared(ctx)
	}
	return nil
}

func (Card4511002Ailimer) HasActiveUltimate(card *CardInstance) bool {
	return card != nil && card.Statuses["爱里默解放"] > 0
}

func (Card4511002Ailimer) OnUltimate(ctx *EffectContext) error {
	candidates := ctx.Engine.nonHeroFieldCardCandidates(ctx.PlayerID)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "ailimer_remove_cards", "爱里默:移除最多3张非人物场上卡牌", candidates, 0, min(3, len(candidates)), func(selected []string) {
		for _, id := range selected {
			ctx.Engine.removeFieldCardFromGameByID(id)
		}
	})
	return nil
}
