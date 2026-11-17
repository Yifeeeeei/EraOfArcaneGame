package game

type Card1221012DragonPrinceDescendant struct{ AlwaysActive }

func (Card1221012DragonPrinceDescendant) ID() string { return "1221012" }

func (Card1221012DragonPrinceDescendant) Name() string { return "龙王子裔" }

func (Card1221012DragonPrinceDescendant) MasteryMax() int { return 2 }

func (Card1221012DragonPrinceDescendant) OnMastery(ctx *EffectContext, level int) error {
	if level != 2 {
		return nil
	}
	searchDeckToHandByPredicateWithResult(ctx, "dragon_prince_search", "检索1个水纹伙伴并使其入场费用-1水", isWaterCompanion, func(card *CardInstance) {
		card.Statuses["入场费用水-1"]++
	})
	return nil
}
