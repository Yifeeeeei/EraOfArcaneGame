package game

import (
	"eraofarcane/model"
)

type Card1411003SandWitchSommer struct{ AlwaysActive }

func (Card1411003SandWitchSommer) ID() string { return "1411003" }

func (Card1411003SandWitchSommer) Name() string { return "沙之魔巫 梭默" }

func (Card1411003SandWitchSommer) ModifySpellArea(ctx *EffectContext, area *SpellArea) {
	if ctx.Source != nil && ctx.Source.Card.Category == model.ElementEarth && !isSorcerySkill(ctx.Source.Card) && *area == SpellAreaSingle {
		*area = SpellAreaSquare
	}
}
