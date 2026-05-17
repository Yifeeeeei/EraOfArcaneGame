package game

import "eraofarcane/model"

type Card1621010DemonLord struct{ AlwaysActive }

func (Card1621010DemonLord) ID() string   { return "1621010" }
func (Card1621010DemonLord) Name() string { return "恶魔尊主" }

func (Card1621010DemonLord) DevourRequirement() map[string]int {
	return map[string]int{model.ElementShadow: 4}
}
