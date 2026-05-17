package game

type Card3021002Foresight struct{ AlwaysActive }

func (Card3021002Foresight) ID() string   { return "3021002" }
func (Card3021002Foresight) Name() string { return "预见" }

func (Card3021002Foresight) OnSpellCast(ctx *EffectContext) error {
	if !isFriendlySpellCast(ctx) || !isSpellBeingCast(ctx) {
		return nil
	}
	candidates := ctx.Engine.friendlyTopDeckCards(ctx.PlayerID, 3, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "foresight_to_bottom",
		"选择任意牌置于牌堆底,未选择的保持在牌堆顶", candidates, 0, len(candidates),
		func(selected []string) {
			for _, id := range selected {
				ctx.Engine.moveDeckCardToBottom(ctx.PlayerID, id)
			}
		})
	return nil
}
