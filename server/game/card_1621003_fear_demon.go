package game

type Card1621003FearDemon struct{ AlwaysActive }

func (Card1621003FearDemon) ID() string { return "1621003" }

func (Card1621003FearDemon) Name() string { return "恐惧魔" }

func (Card1621003FearDemon) DevourRequirement() map[string]int {
	return map[string]int{DevourLife: 3}
}
