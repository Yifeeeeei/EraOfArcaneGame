package game

import (
	"eraofarcane/model"
)

type Card1621108DemonChild struct{ AlwaysActive }

func (Card1621108DemonChild) ID() string { return "1621108" }

func (Card1621108DemonChild) Name() string { return "恶魔之子" }

func (Card1621108DemonChild) DevourCardRequirement() DevourCardRequirement {
	return DevourCardRequirement{Count: 1, Category: model.ElementShadow, CompanionOnly: true}
}
