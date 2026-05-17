package game

const demonSummonerDeathReady = "death_search_ready"

type Card1621009DemonSummoner struct{ AlwaysActive }

func (Card1621009DemonSummoner) ID() string   { return "1621009" }
func (Card1621009DemonSummoner) Name() string { return "唤魔邪术士" }

func (Card1621009DemonSummoner) OnFriendlyDeath(ctx *EffectContext) error {
	ctx.Source.Statuses[demonSummonerDeathReady] = 1
	return nil
}

func (Card1621009DemonSummoner) OnPerTurn(ctx *EffectContext) error {
	if ctx.Source.Statuses[demonSummonerDeathReady] <= 0 {
		return nil
	}
	searchDeckToHandByPredicate(ctx, "demon_summoner_search",
		"检索1个暗影造物或恶魔", func(card *CardInstance) bool {
			return isConstructOrDemon(card)
		})
	ctx.Source.Statuses[demonSummonerDeathReady] = 0
	return nil
}
