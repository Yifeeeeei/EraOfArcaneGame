package game

type Card1421014WindbreathMerchant struct{}

func (Card1421014WindbreathMerchant) ID() string   { return "1421014" }
func (Card1421014WindbreathMerchant) Name() string { return "风息谷旅商" }
func (Card1421014WindbreathMerchant) OnEnter(ctx *EffectContext) error {
	count := 0
	for _, unit := range ctx.Engine.getAllFieldCards(ctx.Engine.State.Players[ctx.PlayerID]) {
		if unit == ctx.Source || !unit.Card.IsCompanion() {
			continue
		}
		if isBeastPlantOrSpirit(unit) {
			count++
			if count >= 3 {
				break
			}
		}
	}
	if count == 0 {
		return nil
	}
	return DrawCards(count)(ctx)
}
