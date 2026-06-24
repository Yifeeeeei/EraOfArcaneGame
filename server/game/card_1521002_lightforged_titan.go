package game

import "strings"

type Card1521002LightforgedTitan struct{ AlwaysActive }

func (Card1521002LightforgedTitan) ID() string   { return "1521002" }
func (Card1521002LightforgedTitan) Name() string { return "光铸泰坦" }
func (Card1521002LightforgedTitan) OnEnter(ctx *EffectContext) error {
	return DrawCards(2)(ctx)
}
func (Card1521002LightforgedTitan) OnDamaged(ctx *EffectContext) error {
	if ctx == nil || ctx.Source == nil || ctx.ExtraData == nil || ctx.ExtraData["damage_source"] != "spell" {
		return nil
	}
	skillNumber, _ := ctx.ExtraData["skill"].(string)
	skill := getCardDB()[skillNumber]
	if skill == nil {
		return nil
	}
	tag := skill.Tag
	if strings.Contains(tag, "驱动") || strings.Contains(tag, "神秘") || strings.Contains(tag, "聚能") {
		ctx.Source.CurrentLife--
		ctx.Source.DamageTakenThisTurn++
	}
	return nil
}
