package game

type Card2511101JiuxiaoRadiance struct{ AlwaysActive }

func (Card2511101JiuxiaoRadiance) ID() string { return "2511101" }

func (Card2511101JiuxiaoRadiance) Name() string { return "九霄辉迹" }

func (Card2511101JiuxiaoRadiance) OnUltimate(ctx *EffectContext) error {
	counts := make([]int, len(ctx.Engine.State.Players))
	for playerID, ps := range ctx.Engine.State.Players {
		counts[playerID] = len(ps.Hand)
	}
	for playerID := range ctx.Engine.State.Players {
		ctx.Engine.discardAllHandCards(playerID)
	}
	for playerID, count := range counts {
		if count > 0 {
			ctx.Engine.drawCards(playerID, count)
		}
	}
	return nil
}
